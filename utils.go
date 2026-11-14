package main

import (
	"regexp"
	"strings"
	"unicode"

	norm "golang.org/x/text/unicode/norm"
)

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
