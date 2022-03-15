package account

import "github.com/KonstantinGasser/sherlock/security"

type Account interface {
	Finder
	EncryptDecrypter
}

type EncryptDecrypter interface {
	Encrypter
	Decrypter
}

type Decrypter interface {
	Decrypt(decry security.Decrypter, v map[string][]byte, passphrase string) (Account, error)
}

type Encrypter interface {
	Encrypt(ency security.Encrypter, passphrase string) (map[string][]byte, error)
}

type Finder interface {
	// Find allows to query an Account for a specific
	// account within its structure
	Find(key string, v interface{}) error
}

type Serializer interface {
	Serialize() ([]byte, error)
}
