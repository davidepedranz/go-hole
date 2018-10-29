package main

import "os"

// getEnvOrDefault reads the value of an environment variable,
// if the variable is not set, it returns a default value.
func getEnvOrDefault(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
