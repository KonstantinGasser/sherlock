package account

import (
	"encoding/json"
	"fmt"

	"github.com/KonstantinGasser/sherlock/security"
)

type Logins map[string]*Login

func DefaultLogins() *Logins { return new(Logins) }

func (login Logins) Type() Type { return LoginT }

func (logins Logins) Find(key string, v interface{}) error {
	login, ok := logins[key]
	if !ok {
		return fmt.Errorf("could not find Login %q", key)
	}

	v = &Login{
		Key:         login.Key,
		AccountName: login.AccountName,
		Password:    login.Password,
	}

	return nil
}

func (logins Logins) Encrypt(encry security.Encrypter, passphrase string) (map[string][]byte, error) {
	out := make(map[string][]byte)

	for key, login := range logins {
		b, err := login.JSON()
		if err != nil {
			return nil, err
		}
		encryBytes, err := encry.Encrypt(passphrase, b)
		if err != nil {
			return nil, err
		}

		out[key] = encryBytes
	}
	return out, nil
}

func (logins Logins) Decrypt(decry security.Decrypter, v map[string][]byte, passphrase string) (Account, error) {

	for key, encryLogin := range v {
		var login Login

		if err := decry.Decrypt(passphrase, encryLogin, &login); err != nil {
			return nil, fmt.Errorf("could not decrypt account %q: %v", key, err)
		}
		logins[key] = &login
	}
	return logins, nil
}

type Login struct {
	Key         string
	AccountName string
	Password    string
}

func (login *Login) JSON() ([]byte, error) {
	return json.Marshal(login)
}
