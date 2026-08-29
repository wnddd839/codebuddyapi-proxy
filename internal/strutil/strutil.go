package strutil

import (
	"cmp"
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// First 返回第一个非空（trim 后）字符串，基于 cmp.Or（Go 1.22+）。
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

// Compact 去除首尾空白。
func Compact(value string) string {
	return strings.TrimSpace(value)
}

// Truncate 返回 trim 后最多 n 字节的子串。
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

// RandomHex 生成 n 字节随机数的十六进制串（2n 个字符）。
func RandomHex(n int) string {
	if n < 1 {
		n = 1
	}
	buf := make([]byte, n)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

// MaskSecret 遮蔽密钥中间段，两侧各保留 visible 个字符。
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
