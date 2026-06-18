package util

import "strings"

func Required(values ...string) bool {
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			return false
		}
	}
	return true
}

func ValidPriority(v string) bool {
	return v == "low" || v == "medium" || v == "high"
}
func ValidStatus(v string) bool {
	return v == "pending" || v == "completed"
}
