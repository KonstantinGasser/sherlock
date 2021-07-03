package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
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

func hash(key string) []byte {
	b := sha256.Sum256([]byte(key))
	hexB := hex.EncodeToString(b[:])
	return []byte(hexB)
}

// InitWithDefault encrypts and empty map[string]interface with a
// provided key
func InitWithDefault(key string, defaultVault interface{}) ([]byte, error) {
	byteVault, err := json.Marshal(defaultVault)
	if err != nil {
		return nil, err
	}
	aesKey := hash(key)
	block, err := aes.NewCipher(aesKey[:16])
	if err != nil {
		return nil, err
	}
	encypted := make([]byte, aes.BlockSize+len(byteVault))
	iv := encypted[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encypted[aes.BlockSize:], byteVault)

	return encypted, err
}

// EncryptVault encrypts the data using the key
func EncryptVault(b []byte, key string) ([]byte, error) {
	aeskey := hash(key)

	block, err := aes.NewCipher(aeskey[:16])
	if err != nil {
		return nil, err
	}
	encrypted := make([]byte, aes.BlockSize+len(b))

	iv := encrypted[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	stream.XORKeyStream(encrypted[aes.BlockSize:], b)

	return encrypted, err
}

// DecryptVault encrypts the data using the key
func DecryptVault(b []byte, key string, v interface{}) error {
	aesKey := hash(key)

	block, err := aes.NewCipher(aesKey[:16])
	if err != nil {
		return err
	}

	decrypted := b[aes.BlockSize:]

	iv := b[:aes.BlockSize]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decrypted, decrypted)

	if err := json.Unmarshal(decrypted, &v); err != nil {
		return err
	}
	return nil
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
