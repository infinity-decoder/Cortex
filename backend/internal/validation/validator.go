package validation

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// ValidateDomain validates a domain name format
func ValidateDomain(domain string) error {
	if domain == "" {
		return errors.New("domain is required")
	}
	
	// Check length (max 253 characters for FQDN)
	if len(domain) > 253 {
		return errors.New("domain too long (maximum 253 characters)")
	}
	
	// Check format - basic domain validation
	// Allows: alphanumeric, hyphens, dots
	// Must start and end with alphanumeric
	matched, err := regexp.MatchString(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`, domain)
	if err != nil {
		return errors.New("invalid domain format")
	}
	
	if !matched {
		return errors.New("invalid domain format")
	}
	
	return nil
}

// ValidateEmail validates an email address format
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	
	// Basic email validation regex
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return errors.New("invalid email format")
	}
	
	if !matched {
		return errors.New("invalid email format")
	}
	
	// Check length
	if len(email) > 254 {
		return errors.New("email too long (maximum 254 characters)")
	}
	
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	
	// Minimum length
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters long")
	}
	
	// Maximum length
	if len(password) > 128 {
		return errors.New("password too long (maximum 128 characters)")
	}
	
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	var missing []string
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasNumber {
		missing = append(missing, "number")
	}
	if !hasSpecial {
		missing = append(missing, "special character")
	}
	
	if len(missing) > 0 {
		return errors.New("password must contain at least one " + strings.Join(missing, ", "))
	}
	
	return nil
}

// ValidateOrgName validates organization name
func ValidateOrgName(name string) error {
	if name == "" {
		return errors.New("organization name is required")
	}
	
	if len(name) < 2 {
		return errors.New("organization name must be at least 2 characters")
	}
	
	if len(name) > 100 {
		return errors.New("organization name too long (maximum 100 characters)")
	}
	
	// Allow alphanumeric, spaces, hyphens, underscores
	matched, err := regexp.MatchString(`^[a-zA-Z0-9\s\-_]+$`, name)
	if err != nil {
		return errors.New("invalid organization name format")
	}
	
	if !matched {
		return errors.New("organization name can only contain letters, numbers, spaces, hyphens, and underscores")
	}
	
	return nil
}

// ValidateFullName validates full name
func ValidateFullName(name string) error {
	if name == "" {
		return errors.New("full name is required")
	}
	
	if len(name) < 2 {
		return errors.New("full name must be at least 2 characters")
	}
	
	if len(name) > 100 {
		return errors.New("full name too long (maximum 100 characters)")
	}
	
	return nil
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
	// Remove null bytes and control characters
	input = strings.ReplaceAll(input, "\x00", "")
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, "\n", "")
	return strings.TrimSpace(input)
}
