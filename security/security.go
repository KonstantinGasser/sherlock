package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

type Encrypter interface {
	Encrypt(passphrase string, in []byte) ([]byte, error)
}

type Decrypter interface {
	Decrypt(passphrase string, in []byte, out interface{}) error
}

type EncryptDecrypter interface {
	Encrypter
	Decrypter
}

type Guard struct{}

// Encrypt encrypts the passed in bytes using the passphrase with the AES-256 encryption standard
func (guard Guard) Encrypt(passphrase string, in []byte) ([]byte, error) {

	aesKey := sha265Hash(passphrase)

	block, err := aes.NewCipher(aesKey[:16])
	if err != nil {
		return nil, fmt.Errorf("could not create AES cipher: %v", err)
	}

	out := make([]byte, aes.BlockSize+len(in))
	iv := out[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(out[aes.BlockSize:], in)

	return out, nil
}

func sha265Hash(s string) []byte {
	b := sha256.Sum256([]byte(s))
	hexb := hex.EncodeToString(b[:])
	return []byte(hexb)
}

// Decrypt encrypts the data using the key - WOW what a comment (jokes on me; it was me - To me: FIX IT)
func (guard Guard) Decrypt(passphrase string, in []byte, out interface{}) error {
	aesKey := sha265Hash(passphrase)

	block, err := aes.NewCipher(aesKey[:16])
	if err != nil {
		return err
	}

	decrypted := in[aes.BlockSize:]

	iv := in[:aes.BlockSize]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decrypted, decrypted)

	if err := json.Unmarshal(decrypted, &out); err != nil {
		return err
	}
	return nil
}
