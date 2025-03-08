package generate

import (
	"math/rand"
	"strings"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars    = "0123456789"
	symbolChars    = "!@#$%&*-_=+.?"
)

func Password(length int, includeLower, includeUpper, includeNumber, includeSymbol bool) string {
	var charSet strings.Builder
	if includeLower {
		charSet.WriteString(lowercaseChars)
	}
	if includeUpper {
		charSet.WriteString(uppercaseChars)
	}
	if includeNumber {
		charSet.WriteString(numberChars)
	}
	if includeSymbol {
		charSet.WriteString(symbolChars)
	}

	password := make([]rune, length)
	for i := range password {
		password[i] = rune(charSet.String()[rand.Intn(charSet.Len())])
	}

	return string(password)
}
