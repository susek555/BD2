package enums

import (
	"errors"
	"strings"
)

var ErrStringConversion = errors.New("fetched database value couldn't be mapped to string")

func convertAppFormatToDBFormat(s string) string {
	sepWords := strings.Split(s, " ")
	joinedWords := strings.Join(sepWords, "_")
	return strings.ToLower(joinedWords)
}

func convertDBFormatToAppFormat(s string, wholeUpper bool) string {
	sepWords := strings.Split(s, "_")
	joinedWords := strings.Join(sepWords, " ")
	if wholeUpper {
		return strings.ToUpper(joinedWords)
	} else {
		return capitalize(joinedWords)
	}
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
