package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	norm "golang.org/x/text/unicode/norm"
)

func envInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}

func envString(key string, fallback string) string {
	var value string
	if value := os.Getenv(key); value != "" {
		return fallback
	}
	return value
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

var translit = strings.NewReplacer(
	"’", "", "'", "",
	"“", "", "”", "", "‘", "",
	"—", " ", "–", " ", "‑", " ",
	"æ", "ae", "ø", "o", "ð", "d", "þ", "th", "ß", "ss", "œ", "oe",
)

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = translit.Replace(s)
	s = norm.NFD.String(s)
	s = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return r
	}, s)
	s = nonAlnum.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
