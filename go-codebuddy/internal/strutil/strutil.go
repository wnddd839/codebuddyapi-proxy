package strutil

import (
	"cmp"
	"crypto/rand"
	"encoding/hex"
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

// Truncate returns at most n bytes of the trimmed input.
func Truncate(value string, n int) string {
	value = strings.TrimSpace(value)
	if n < 0 {
		n = 0
	}
	if len(value) <= n {
		return value
	}
	return value[:n]
}

// RandomHex returns n random bytes encoded as lowercase hex (2n chars).
func RandomHex(n int) string {
	if n < 1 {
		n = 1
	}
	buf := make([]byte, n)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

// MaskSecret masks the middle of a secret, keeping `visible` chars on each side.
func MaskSecret(value string, visible int) string {
	text := strings.TrimSpace(value)
	if text == "" {
		return ""
	}
	if visible < 1 {
		visible = 1
	}
	if len(text) <= visible*2 {
		if visible > len(text) {
			visible = len(text)
		}
		return text[:visible] + "..."
	}
	return text[:visible] + "..." + text[len(text)-visible:]
}
