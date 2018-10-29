package main

import (
	"sort"
	"testing"
)

var (
	testBlacklistPath = "./data/test-blacklist.txt"
)

func TestBlacklist_Size(t *testing.T) {
	expectedSize := 3
	blacklist := LoadBlacklistOrFail(testBlacklistPath)
	size := blacklist.Size()
	if size != expectedSize {
		t.Errorf("Size was incorrect, got: %d, expected: %d", size, expectedSize)
	}
}

func TestBlacklist_Contains_Present(t *testing.T) {
	domain := "a.example.com"
	blacklist := LoadBlacklistOrFail(testBlacklistPath)
	if !blacklist.Contains(domain) {
		t.Errorf("Blacklist should contain %s, but it was absent", domain)
	}
}

func TestBlacklist_Contains_Absent(t *testing.T) {
	domain := "absent.example.com"
	blacklist := LoadBlacklistOrFail(testBlacklistPath)
	if blacklist.Contains(domain) {
		t.Errorf("Blacklist should not contain %s, but it was present", domain)
	}
}

func TestBlacklist_Contains_UpperCase(t *testing.T) {
	domain := "A.example.COM"
	blacklist := LoadBlacklistOrFail(testBlacklistPath)
	if !blacklist.Contains(domain) {
		t.Errorf("Blacklist should contain %s, but it was absent", domain)
	}
}

func TestBlacklist_FileSorted(t *testing.T) {
	blacklist := LoadBlacklistOrFail(blacklistPath)
	sorted := sort.StringsAreSorted(blacklist.array)
	if !sorted {
		t.Errorf("The blacklist file at path %s must be sorted", blacklistPath)
	}
}
