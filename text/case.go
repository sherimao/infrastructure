package text

import (
	"bytes"
	"strings"
)

// CamelCase 转大驼峰
func CamelCase(snake string) string {
	newStr := make([]rune, 0)
	upNextChar := true
	snake = strings.ToLower(snake)
	for _, chr := range snake {
		switch {
		case upNextChar:
			upNextChar = false
			if 'a' <= chr && chr <= 'z' {
				chr -= 'a' - 'A'
			}
		case chr == '_':
			upNextChar = true
			continue
		default:
		}
		if chr != '_' {
			newStr = append(newStr, chr)
		}
	}
	return string(newStr)
}

// SnakeCase 转下划线
func SnakeCase(camel string) string {
	var buf bytes.Buffer
	for _, c := range camel {
		if 'A' <= c && c <= 'Z' {
			if buf.Len() > 0 {
				buf.WriteRune('_')
			}
			buf.WriteRune(c - 'A' + 'a')
		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
