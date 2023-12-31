package checker

import (
	"regexp"
	"strings"
)

// IsEmpty - check if string is empty
func IsEmpty(s string) bool {
	return strings.Trim(s, " ") == ""
}

// ValidateEmail - check whether the email format is compliant
func ValidateEmail(s string) bool {
	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return pattern.MatchString(s)
}
