package space

import (
	"encoding/json"
	"fmt"

	"github.com/KonstantinGasser/sherlock/core/account"
	"github.com/KonstantinGasser/sherlock/security"
)

const (
	jsonKeyLogin          = "Logins"
	jsonKeyAwsConsole     = "AwsConsoles"
	jsonKeyAwsApiAccesses = "AwsApiAccesses"
)

// Space is an isolated directory to store accounts and files
//
// A space in sherlock describes a collection of accounts in
// the same context. Spaces under sherlock are stored as
// directory on the file system. The space is stored in JSON
// format where the meta data is plaintext while all secret information
// such as passwords, logins, API-Keys and file contents are only
// stored encrypted within the JSON-Object
type Space struct {
	// Key is a unique identifier across sherlock
	Key string
	// guardian allows the space to perform
	// the encryption and decryption of account and
	// file data
	guardian security.Encrypter
	// Accounts stores all accounts mapped to the space
	Accounts *Accounts
}

// accounts manages all available sherlock accounts
//
// info:
//		each account within its map is encrypted,
//		the meta-data for an account is human readable.
//		while serializing the space only the actual accounts (login,aws-*)
// 		must be encrypted.
type Accounts struct {
	// Logins holds all Login Types mapped to the space
	Logins map[string]*account.Login // type not there yet
	// AwsConsoles holds all AWS-Console login Types mapped
	// to the space
	AwsConsoles map[string]*account.AwsConsole
	// AwsApiAccesses holds all AWS-Programmatic-Access identities
	// mapped to the space
	AwsApiAccesses map[string]*account.AwsApiAccess
}

func New(key string) *Space {
	return &Space{
		Key:      key,
		guardian: security.Guard{},
		Accounts: &Accounts{
			Logins:         make(map[string]*account.Login),
			AwsConsoles:    make(map[string]*account.AwsConsole),
			AwsApiAccesses: make(map[string]*account.AwsApiAccess),
		},
	}
}

func (space Space) ToCipherSpace(passphrase string) (*CipherSpace, error) {

	var marshal = func(passphrase string, accs interface{}) (map[string][]byte, error) {
		// uff this feels like a hack...
		// issue is that I cannot pass in a map[string]account.Login when accs is a map[string]account.Serializer
		// even-though the account.Login for example implements the account.Serializer interface
		tmp, ok := accs.(map[string]account.Serializer)
		if !ok {
			return nil, fmt.Errorf("sherlock account type must implement the account.Serializer interface")
		}
		out := make(map[string][]byte)

		for key, acc := range tmp {
			b, err := acc.Serialize()
			if err != nil {
				return nil, err
			}

			encrypted, err := space.guardian.Encrypt(passphrase, b) // where do we get the passphrase from ??
			if err != nil {
				return nil, err
			}
			out[key] = encrypted
		}

		return out, nil
	}

	loginsEncry, err := marshal(passphrase, space.Accounts.Logins)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt Login Accounts: %v", err)
	}

	awsConsoleEncry, err := marshal(passphrase, space.Accounts.AwsConsoles)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt AWS-Console Accounts: %v", err)
	}

	awsApiAccessEncry, err := marshal(passphrase, space.Accounts.AwsApiAccesses)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt AWS-API Access Accounts: %v", err)
	}

	return &CipherSpace{
		Key: space.Key,
		Accounts: map[string]map[string][]byte{
			jsonKeyLogin:          loginsEncry,
			jsonKeyAwsConsole:     awsConsoleEncry,
			jsonKeyAwsApiAccesses: awsApiAccessEncry,
		},
	}, nil
}

// CipherSpace represents a space where all accounts
// are stored as bytes (encrypted). This type should be
// used when reading and unmarshaling a space from the
// filesystem and can be seen as a transfer struct to a Space
//
// CipherSpace implements the Serializer & Deserializer interface
type CipherSpace struct {
	Key      string
	guardian security.Decrypter

	// Example structure of Accounts
	// {
	// 	"Login": {
	// 		"login-1": encrypted-data,
	// 		"...": "..."
	// 	},
	// 	"AwsConsoles": {...},
	//	"AwsApiAccesses": {...}
	// }
	Accounts map[string]map[string][]byte
}

func NewCipher(key string) *CipherSpace {
	return &CipherSpace{
		Key:      key,
		guardian: security.Guard{},
		Accounts: make(map[string]map[string][]byte),
	}
}

func (cSpace CipherSpace) Serialize() ([]byte, error) {
	return json.MarshalIndent(cSpace, "", "\t")
}

func (cSpace *CipherSpace) Deserialize(v []byte) error {

	if err := json.Unmarshal(v, cSpace); err != nil {
		return err
	}

	cSpace.guardian = security.Guard{}
	return nil
}

// ToSpace asserts the CipherSpace to a regular Space
// All Accounts within the CipherSpace are decrypted
func (cSpace CipherSpace) ToSpace(passphrase string) (*Space, error) {

	var unmarshal = func(passphrase string, accs map[string][]byte, newFunc func() interface{}) (map[string]interface{}, error) {
		var out map[string]interface{}

		for key, acc := range accs {
			var newAccount = newFunc()
			if err := cSpace.guardian.Decrypt(passphrase, acc, &newAccount); err != nil {
				return nil, err
			}
			out[key] = newAccount
		}
		return out, nil
	}

	logins, err := unmarshal(passphrase, cSpace.Accounts[jsonKeyLogin], nil)
	if err != nil {
		return nil, err
	}

	awsConsoles, err := unmarshal(passphrase, cSpace.Accounts[jsonKeyAwsConsole], nil)
	if err != nil {
		return nil, err
	}

	AwsApiAccesses, err := unmarshal(passphrase, cSpace.Accounts[jsonKeyAwsApiAccesses], nil)
	if err != nil {
		return nil, err
	}

	return &Space{
		Key:      cSpace.Key,
		guardian: security.Guard{},
		Accounts: &Accounts{
			Logins:         logins,
			AwsConsoles:    awsConsoles,
			AwsApiAccesses: AwsApiAccesses,
		},
	}, nil

}
