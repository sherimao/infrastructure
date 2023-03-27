package db

import (
	"fmt"
	"strings"
)

// EscapeLike 处理like
func EscapeLike(s string) string {
	s = strings.Replace(s, "%", "\\%", -1)
	s = strings.Replace(s, "_", "\\_", -1)
	return fmt.Sprintf("%%%s%%", s)
}
