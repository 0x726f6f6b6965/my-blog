package utils

import (
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
)

var commonWords = map[string]struct{}{
	"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
	"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
}

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func filter(tokens []string) []string {
	r := make([]string, 0)
	for _, token := range tokens {
		// get lowercase
		lower := strings.ToLower(token)
		// filter common words
		if _, ok := commonWords[lower]; !ok {
			// add base form
			r = append(r, english.Stem(lower, false))
		}
	}
	return r
}

func GetTokens(text string) []string {
	tokens := tokenize(text)
	tokens = filter(tokens)
	return tokens
}
