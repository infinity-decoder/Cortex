package validation

import (
	"unicode"
)

// PasswordStrength represents password strength levels
type PasswordStrength int

const (
	PasswordStrengthWeak PasswordStrength = iota
	PasswordStrengthFair
	PasswordStrengthGood
	PasswordStrengthStrong
)

// CheckPasswordStrength checks password strength and returns level
func CheckPasswordStrength(password string) PasswordStrength {
	if len(password) < 8 {
		return PasswordStrengthWeak
	}
	
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
		length     = len(password)
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
	
	score := 0
	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSpecial {
		score++
	}
	if length >= 12 {
		score++
	}
	if length >= 16 {
		score++
	}
	
	switch {
	case score >= 6:
		return PasswordStrengthStrong
	case score >= 4:
		return PasswordStrengthGood
	case score >= 2:
		return PasswordStrengthFair
	default:
		return PasswordStrengthWeak
	}
}

// GetPasswordRequirements returns a list of password requirements
func GetPasswordRequirements() []string {
	return []string{
		"At least 12 characters long",
		"Contains at least one uppercase letter",
		"Contains at least one lowercase letter",
		"Contains at least one number",
		"Contains at least one special character",
	}
}

// ValidatePasswordWithDetails returns detailed validation errors
func ValidatePasswordWithDetails(password string) (bool, []string) {
	var errors []string
	
	if len(password) < 12 {
		errors = append(errors, "Password must be at least 12 characters long")
	}
	
	if len(password) > 128 {
		errors = append(errors, "Password too long (maximum 128 characters)")
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
	
	if !hasUpper {
		errors = append(errors, "Password must contain at least one uppercase letter")
	}
	if !hasLower {
		errors = append(errors, "Password must contain at least one lowercase letter")
	}
	if !hasNumber {
		errors = append(errors, "Password must contain at least one number")
	}
	if !hasSpecial {
		errors = append(errors, "Password must contain at least one special character")
	}
	
	return len(errors) == 0, errors
}
