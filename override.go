package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

const overridePath = "./data/override.txt"

func LoadOverrideListOrFail(path string) map[string]string {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Printf("No domain override file found (%s). Continuing with empty set.", path)
			return map[string]string{}
		} else {
			log.Panic(err)
		}
	}

	// open the file
	file, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	overrideMap := map[string]string{}

	// read the file line-by-line
	// NB: the file MUST be ORDERED and all domains LOWER CASE!
	reader := bufio.NewReader(file)
	i := 0
	for ; ; i++ {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}

		// NOTE: expecting "domain  ipv4"
		tuple := strings.Fields(strings.TrimSpace(line))
		overrideMap[tuple[0]] = tuple[1]
	}
	log.Printf("Loaded %d entries in override map from %s.", i, path)
	return overrideMap
}
