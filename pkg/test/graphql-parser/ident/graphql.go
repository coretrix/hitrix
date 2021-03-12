package ident

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func ParseMixedCaps(name string) Name {
	var words Name

	runes := []rune(name)
	w, i := 0, 0
	for i+1 <= len(runes) {
		eow := false
		if i+1 == len(runes) {
			eow = true
		} else if unicode.IsLower(runes[i]) && unicode.IsUpper(runes[i+1]) {
			eow = true
		} else if i+2 < len(runes) && unicode.IsUpper(runes[i]) && unicode.IsUpper(runes[i+1]) && unicode.IsLower(runes[i+2]) {
			eow = true

			if string(runes[i:i+3]) == "IDs" {
				eow = false
			}
		}
		i++
		if !eow {
			continue
		}

		word := string(runes[w:i])
		if initialism, ok := isInitialism(word); ok {
			words = append(words, initialism)
		} else if i1, i2, ok := isTwoInitialisms(word); ok {
			words = append(words, i1, i2)
		} else {
			words = append(words, word)
		}
		w = i
	}
	return words
}

type Name []string

func (n Name) ToLowerCamelCase() string {
	for i, word := range n {
		if i == 0 {
			n[i] = strings.ToLower(word)
			continue
		}
		r, size := utf8.DecodeRuneInString(word)
		n[i] = string(unicode.ToUpper(r)) + strings.ToLower(word[size:])
	}
	return strings.Join(n, "")
}

func isInitialism(word string) (string, bool) {
	initialism := strings.ToUpper(word)
	_, ok := initialisms[initialism]
	return initialism, ok
}

func isTwoInitialisms(word string) (string, string, bool) {
	word = strings.ToUpper(word)
	for i := 2; i <= len(word)-2; i++ {
		_, ok1 := initialisms[word[:i]]
		_, ok2 := initialisms[word[i:]]
		if ok1 && ok2 {
			return word[:i], word[i:], true
		}
	}
	return "", "", false
}

var initialisms = map[string]struct{}{
	"ACL":   {},
	"API":   {},
	"ASCII": {},
	"CPU":   {},
	"CSS":   {},
	"DNS":   {},
	"EOF":   {},
	"GUID":  {},
	"HTML":  {},
	"HTTP":  {},
	"HTTPS": {},
	"ID":    {},
	"IP":    {},
	"JSON":  {},
	"LHS":   {},
	"QPS":   {},
	"RAM":   {},
	"RHS":   {},
	"RPC":   {},
	"SLA":   {},
	"SMTP":  {},
	"SQL":   {},
	"SSH":   {},
	"TCP":   {},
	"TLS":   {},
	"TTL":   {},
	"UDP":   {},
	"UI":    {},
	"UID":   {},
	"UUID":  {},
	"URI":   {},
	"URL":   {},
	"UTF8":  {},
	"VM":    {},
	"XML":   {},
	"XMPP":  {},
	"XSRF":  {},
	"XSS":   {},
	"RSS":   {},
}
