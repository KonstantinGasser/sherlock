package account

import (
	"encoding/json"

	"github.com/KonstantinGasser/sherlock/security"
)

type AwsConsoles map[string]*AwsConsole

func DefaultAwsConsoles() *AwsConsoles { return new(AwsConsoles) }

func (consoles AwsConsoles) Type() Type { return AwsConsoleT }

func (consoles AwsConsoles) Encrypt(encry security.Encrypter, passpharase string) (map[string][]byte, error) {
	out := make(map[string][]byte)

	for key, console := range consoles {
		b, err := console.JSON()
		if err != nil {
			return nil, err
		}
		encryBytes, err := encry.Encrypt(passpharase, b)
		if err != nil {
			return nil, err
		}

		out[key] = encryBytes
	}
	return out, nil
}

func (consoles AwsConsoles) Decrypt(decry security.Decrypter, v map[string][]byte, passphrase string) (Account, error) {
	for key, encryAccess := range v {
		var console AwsConsole

		if err := decry.Decrypt(passphrase, encryAccess, &console); err != nil {
			return nil, err
		}
		consoles[key] = &console
	}
	return consoles, nil
}

func (consoles AwsConsoles) Find(key string, v interface{}) error {
	return nil
}

type AwsConsole struct {
	AccountID   [12]byte
	AccountName string
	Password    string
}

func (console AwsConsole) JSON() ([]byte, error) {
	return json.Marshal(console)
}
