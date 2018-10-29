package main

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	cacheHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "cache",
			Name:      "operation_duration_seconds",
			Help:      "Duration of an operation on the cache.",
			Buckets:   []float64{1e-6, 1.75e-6, 2.5e-6, 3.75e-6, 5e-6, 6.25e-6, 7.5e-6, 8.75e-6, 1e-5},
		},
		[]string{"operation"},
	)
)

// Cache represents a cache for the DNS queries already resolved.
type Cache struct {
	c *cache.Cache
}

// NewCache initializes and return a new cache.
func NewCache() Cache {
	return Cache{c: cache.New(1*time.Minute, 30*time.Second)}
}

// Get an item from the cache. Returns the item or nil,
// and a bool indicating whether the key was found.
func (cache Cache) Get(question *dns.Question) (*dns.Msg, bool) {

	// collect metrics
	timer := prometheus.NewTimer(makeObserver("get"))
	defer timer.ObserveDuration()

	// get the value from the cache
	key := serializeQuestion(question)
	value, found := cache.c.Get(key)
	if found {
		return value.(*dns.Msg), true
	}
	return nil, false
}

// Set an item in the cache, replacing any existing one with the same key.
func (cache Cache) Set(question *dns.Question, message *dns.Msg, duration time.Duration) {

	// collect metrics
	timer := prometheus.NewTimer(makeObserver("set"))
	defer timer.ObserveDuration()

	// add the value in the cache
	key := serializeQuestion(question)
	cache.c.Set(key, message, duration)
}

// serializeQuestion maps a DNS question to a string that represents it.
// This method is useful to use the question as a key in a cache of type map[string]something.
func serializeQuestion(question *dns.Question) string {
	return fmt.Sprintf("%d.%d.%s", question.Qtype, question.Qclass, question.Name)
}

// makeObserver creates a prometheus Observer to measure the duration of
// the operations on the cache.
func makeObserver(operation string) prometheus.Observer {
	return prometheus.ObserverFunc(func(value float64) {
		cacheHistogram.WithLabelValues(operation).Observe(value)
	})
}
