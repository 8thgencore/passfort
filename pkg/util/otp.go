package util

import (
	"fmt"
	"math/rand"
)

// GenerateOTP generates a six-digit one-time password (OTP) code as a string.
// The code consists of random digits between 000000 and 999999.
func GenerateOTP() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

// ValidateOTP validates if the provided OTP code is a six-digit numeric string.
func ValidateOTP(otp string) bool {
	// Use regular expression or any other suitable method to validate the OTP.
	// For simplicity, this example checks if the length is 6 and all characters are digits.
	if len(otp) != 6 {
		return false
	}

	for _, digit := range otp {
		if digit < '0' || digit > '9' {
			return false
		}
	}

	return true
}
