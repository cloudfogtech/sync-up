package utils

import (
	"strings"
	"unicode"
)

func GetPrintable(raw string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, raw)
}
