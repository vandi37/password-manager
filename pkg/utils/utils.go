package utils

import "strings"

func SplitString(s string, num int) ([]string, string) {
	var result []string
	for len(s) > num {
		result = append(result, s[:num])
		s = s[num:]
	}
	return result, strings.TrimSpace(s)
}
