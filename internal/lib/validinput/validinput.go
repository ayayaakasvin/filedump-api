package validinput

import (
	"regexp"
	"strings"
)

var (
	uppercase = regexp.MustCompile(`[A-Z]`)
	lowercase = regexp.MustCompile(`[a-z]`)
	digit     = regexp.MustCompile(`[0-9]`)

	minLengthPassword = 8
	minLengthUsername = 3
)

func IsValidPassword(password string) bool {
	if len(password) < minLengthPassword {
		return false
	}

	if !uppercase.MatchString(password) {
		return false
	}

	if !lowercase.MatchString(password) {
		return false
	}

	if !digit.MatchString(password) {
		return false
	}

	return true
}

func IsValidUsername(username string) bool {
	if len(username) < minLengthUsername {
		return false
	}

	if !(lowercase.MatchString(username) || uppercase.MatchString(username)) {
		return false
	}

	return true
}

func IsValidFileName(filename string) bool {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return false
	}
	// Disallow common illegal filename characters
	illegal := regexp.MustCompile(`[\\/:\*\?"<>\|]`)
	return !illegal.MatchString(filename)
}
