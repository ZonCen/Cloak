package helpers

import "strings"

func RemoveSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

func CheckSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}
