package utils

import (
	"strings"
	"unicode"
)

func NormalizeTag(tagName string) string {
	tagName = strings.ToLower(tagName)

	tagName = strings.TrimSpace(tagName)

	var sb strings.Builder
	for _, char := range tagName {
		if !unicode.IsSpace(char) {
			sb.WriteRune(char)
		}
	}
	return sb.String()
}
