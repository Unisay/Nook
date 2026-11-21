// Package display prepares document text for terminal layout.
package display

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Text replaces each control or malformed byte with one visible cell. Each
// replacement has the original byte length so markdown source spans remain valid.
func Text(s string) string {
	var b strings.Builder
	start := 0
	for i, r := range s {
		if !unicode.IsControl(r) && r != utf8.RuneError {
			continue
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size > 1 {
			continue
		}
		b.WriteString(s[start:i])
		switch {
		case r == '\t':
			b.WriteByte(' ')
		case size == 2:
			b.WriteString("¡")
		default:
			b.WriteByte('?')
		}
		start = i + size
	}
	if start == 0 {
		return s
	}
	b.WriteString(s[start:])
	return b.String()
}
