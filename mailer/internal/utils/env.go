package utils

import (
	"os"
	"strconv"
	"strings"
)

func String(key string, def string) string {
	result := os.Getenv(key)
	if result == "" {
		return def
	}
	return strings.TrimSpace(result)
}

func Int(key string, def int) int {
	result := os.Getenv(key)
	if result == "" {
		return def
	}
	value, err := strconv.Atoi(strings.TrimSpace(result))
	if err != nil {
		return def
	}
	return value
}
