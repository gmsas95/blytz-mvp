package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^\+?[\d\s\-\(\)]{10,15}$`)
	return phoneRegex.MatchString(phone)
}

// ValidateUsername validates username format
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// SanitizeString removes potentially harmful characters
func SanitizeString(input string) string {
	input = strings.TrimSpace(input)
	// Remove HTML tags
	input = regexp.MustCompile(`<[^>>]*>`).ReplaceAllString(input, "")
	return input
}
