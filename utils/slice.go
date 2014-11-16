package utils

import (
	"strings"
)

func IsSliceContainsStr(slice []string, str string) bool {
	str = strings.ToLower(str)
	for _, s := range slice {
		if strings.ToLower(s) == str {
			return true
		}
	}
	return false
}
