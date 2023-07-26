package str

import (
	"strconv"
	"strings"
)

func StrTrim(in ...string) []string {
	out := make([]string, len(in))
	for k, v := range in {
		out[k] = strings.TrimSpace(v)
	}

	return out
}

func StrIsNumber(str string) bool {
	if str == "" {
		return false
	}

	if strings.Contains(str, ".") {
		_, err := strconv.ParseFloat(str, 64)
		return err == nil
	}

	_, err := strconv.Atoi(str)
	return err == nil
}

func StrTrimSpecialCharacter(str string) string {
	specialCharacters := []string{"&&", "||", "!", "(", ")", "{", "}", "[", "]", "^", "\"", "~", "*", "?", ":"}
	for _, v := range specialCharacters {
		str = strings.Replace(str, v, "", -1)
	}
	return str
}

func StrCleanString(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}
