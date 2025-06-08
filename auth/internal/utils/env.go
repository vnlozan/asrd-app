package utils

import (
	"os"
	"strings"
)

func String(key string, def string) string {
	result := os.Getenv(key)
	if result == "" {
		return def
	}
	return strings.TrimSpace(result)
}
