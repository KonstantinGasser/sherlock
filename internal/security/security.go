package security

import (
	"fmt"
	"regexp"
)

const (
	// LowSecurity if password rating in rage(0,49)
	Low int = iota
	// MidSecurity if password rating in rage(50, 74)
	Satifsfied
	// HighSecurity if password rating in rage(75, 100)
	High
)

// EncryptVault encrypts the data using the key
func EncryptVault(b []byte, key string) error {
	return fmt.Errorf("security.EncryptVault: not implemented")
}

// DecryptVault encrypts the data using the key
func DecryptVault(b []byte, key string) error {
	return fmt.Errorf("security.DecryptVault: not implemented")
}

// PasswordStrength evaluates how strong the password is based on
// the variety and diversity of the chosen characters
func PasswordStrength(password string) int {

	eval := func() int {
		var strength int
		regN := regexp.MustCompile(`[0-9]`)
		numbers := regN.FindAllString(password, -1)
		strength += len(numbers) * 4

		regC := regexp.MustCompile(`[A-Z]`)
		caper := regC.FindAllString(password, -1)
		strength += (len(password) - len(caper)) * 2

		regL := regexp.MustCompile(`[a-z]`)
		lower := regL.FindAllString(password, -1)
		strength += (len(password) - len(lower)) * 2

		regS := regexp.MustCompile(`[$#_-]`)
		specials := regS.FindAllString(password, -1)
		strength += len(specials) * 6
		return strength
	}
	switch strength := eval(); {
	case (strength >= 75):
		return High
	case (strength >= 45 && strength < 74):
		return Satifsfied
	default:
		return Low
	}
}
