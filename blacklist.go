package main

import (
	"bufio"
	"bytes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/willf/bloom"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	blacklistPath = "./data/blacklist.txt"

	blacklistHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "blacklist",
			Name:      "lookup_duration_seconds",
			Help:      "Duration of a domain lookup in the blacklist.",
			Buckets:   []float64{1e-6, 1.75e-6, 2.5e-6, 3.75e-6, 5e-6, 6.25e-6, 7.5e-6, 8.75e-6, 1e-5},
		},
		[]string{"bloom_filter", "array"},
	)
)

// Blacklist represents a set of domains to block.
// Blocked domains serve ads, tracking, malware, etc.
type Blacklist struct {
	filter *bloom.BloomFilter
	array  []string
}

// Size returns the number of domains in the blacklist.
func (blacklist *Blacklist) Size() int {
	return len(blacklist.array)
}

// Contains checks if the given domain belongs to the blacklist:
// the method returns true if the domain is present, false otherwise.
func (blacklist *Blacklist) Contains(domain string) bool {
	start := time.Now()

	// check the bloom filter first: it either says "definitely no present" or "maybe present"
	lower := strings.ToLower(domain)
	possiblyPresent := blacklist.filter.TestString(lower)
	if possiblyPresent {

		// the domain might be present... we need to manually check the list
		index := sort.SearchStrings(blacklist.array, lower)
		present := index < len(blacklist.array) && blacklist.array[index] == lower

		// collect metrics
		duration := time.Since(start).Seconds()
		if present {
			blacklistHistogram.WithLabelValues("maybe", "present").Observe(duration)
		} else {
			blacklistHistogram.WithLabelValues("maybe", "absent").Observe(duration)
		}

		return present
	}

	// collect metrics
	duration := time.Since(start).Seconds()
	blacklistHistogram.WithLabelValues("absent", "absent").Observe(duration)

	// if here, the domain is not present at all
	return false
}

// LoadBlacklistOrFail loads the blacklist from the given file
// and panics if there are errors with loading the list.
func LoadBlacklistOrFail(path string) *Blacklist {
	domains, err := LoadBlacklist(path)
	if err != nil {
		log.Panic(err)
	}
	return domains
}

// LoadBlacklist loads the blacklist from the given file.
func LoadBlacklist(path string) (*Blacklist, error) {

	// open the file
	file, errCount := os.Open(path)
	if errCount != nil {
		return nil, errCount
	}
	defer file.Close()

	// count the number of blocked domains
	lines, errCount := countLines(file)
	if errCount != nil {
		return nil, errCount
	}

	// allocate the data structure of optimal Size
	blacklist := Blacklist{
		filter: bloom.New(uint(lines)*10, 5),
		array:  make([]string, lines),
	}

	// start again from the beginning of the file
	_, errSeek := file.Seek(0, 0)
	if errSeek != nil {
		return nil, errSeek
	}

	// read the file line-by-line
	// NB: the file MUST be ORDERED and all domains LOWER CASE!
	reader := bufio.NewReader(file)
	for i := 0; i < lines; i++ {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// NB: we trim spaces to eliminate the '\n' at the end
		trimmed := strings.TrimSpace(line)

		// add the domain to both the bloom filter and the ordered array
		blacklist.filter.AddString(trimmed)
		blacklist.array[i] = trimmed
	}

	return &blacklist, nil
}

// countLines count the number of lines in the given file.
// Please note that it will move to the end of the file, so rewind is needed.
func countLines(file io.Reader) (int, error) {
	size := 2 * 1024
	buffer := make([]byte, size)
	count := 0
	for {
		c, err := file.Read(buffer)
		count += bytes.Count(buffer[:c], []byte{'\n'})

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
