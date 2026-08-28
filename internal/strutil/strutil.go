package strutil

import (
	"cmp"
	"strings"
)

// First returns the first non-empty trimmed string using cmp.Or (Go 1.22+).
func First(values ...string) string {
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || value == "<nil>" {
			continue
		}
		cleaned = append(cleaned, value)
	}
	if len(cleaned) == 0 {
		return ""
	}
	return cmp.Or(cleaned...)
}

// Compact trims surrounding whitespace.
func Compact(value string) string {
	return strings.TrimSpace(value)
}
