package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/m1/go-generate-password/generator"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const (
	// minStrength is the lower limit a password has to be secure
	minStrength = 60
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

// Encrypt encrypts the data using the key
func Encrypt(b []byte, key string) ([]byte, error) {
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

	return encrypted, nil
}

// Decrypt encrypts the data using the key
func Decrypt(b []byte, key string, v interface{}) error {
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
func PasswordStrength(password string) error {
	return passwordvalidator.Validate(password, minStrength)
}

func GenPassword(len int) (string, error) {
	gen, err := generator.New(&generator.Config{
		Length:                     len,
		IncludeSymbols:             true,
		IncludeNumbers:             true,
		IncludeLowercaseLetters:    true,
		IncludeUppercaseLetters:    true,
		ExcludeSimilarCharacters:   true,
		ExcludeAmbiguousCharacters: true,
	})
	if err != nil {
		return "", err
	}
	password, err := gen.Generate()
	return *password, err
}
