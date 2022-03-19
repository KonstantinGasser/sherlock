package account

import (
	"encoding/json"

	"github.com/KonstantinGasser/sherlock/security"
)

type AwsApiAccesses map[string]*AwsApiAccess

func DefaultAwsAccesses() *AwsApiAccesses { return new(AwsApiAccesses) }

func (consoles AwsApiAccesses) Type() Type { return AwsApiAccessT }

func (accesses AwsApiAccesses) Encrypt(encry security.Encrypter, passpharase string) (map[string][]byte, error) {
	out := make(map[string][]byte)

	for key, access := range accesses {
		b, err := access.JSON()
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

func (accesses AwsApiAccesses) Decrypt(decry security.Decrypter, v map[string][]byte, passphrase string) (Account, error) {
	for key, encryAccess := range v {
		var access AwsApiAccess

		if err := decry.Decrypt(passphrase, encryAccess, &access); err != nil {
			return nil, err
		}
		accesses[key] = &access
	}
	return accesses, nil
}

func (accesses AwsApiAccesses) Find(key string, v interface{}) error {
	return nil
}

type AwsApiAccess struct {
	AccessID  string
	AccessKey string
	Region    string
}

func (awsA AwsApiAccess) JSON() ([]byte, error) {
	return json.Marshal(awsA)
}
