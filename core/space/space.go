package space

import (
	"encoding/json"
	"fmt"

	"github.com/KonstantinGasser/sherlock/core/account"
	"github.com/KonstantinGasser/sherlock/security"
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
	guardian security.EncryptDecrypter
	// Accounts stores all accounts mapped to the space
	Accounts *accounts
}

// accounts manages all available sherlock accounts
//
// info:
//		each account within its map is encrypted,
//		the meta-data for an account is human readable.
//		while serializing the space only the actual accounts (login,aws-*)
// 		must be encrypted.
type accounts struct {
	// Logins holds all Login Types mapped to the space
	Logins map[string]interface{} // type not there yet
	// AwsConsoles holds all AWS-Console login Types mapped
	// to the space
	AwsConsoles map[string]interface{} // type not there yet
	// AwsApiAccesses holds all AWS-Programmatic-Access identities
	// mapped to the space
	AwsApiAccesses map[string]interface{} // type not there yet
}

func New(key string) *Space {
	return &Space{
		Key:      key,
		guardian: security.Guard{},
		Accounts: &accounts{
			Logins:         make(map[string]interface{}),
			AwsConsoles:    make(map[string]interface{}),
			AwsApiAccesses: make(map[string]interface{}),
		},
	}
}

// Serialize marshal's the content of a space
func (space Space) Serialize() ([]byte, error) {

	loginsEncry, err := space.marshalAccount(space.Accounts.Logins) // this is ok. fixing in post, once account types exist
	if err != nil {
		return nil, fmt.Errorf("could not encrypt Login Accounts: %v", err)
	}

	awsConsoleEncry, err := space.marshalAccount(space.Accounts.AwsConsoles)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt AWS-Console Accounts: %v", err)
	}

	awsApiAccessEncry, err := space.marshalAccount(space.Accounts.AwsApiAccesses)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt AWS-API Access Accounts: %v", err)
	}

	return json.MarshalIndent(map[string]interface{}{
		"Key": space.Key,
		"Accounts": map[string]interface{}{
			"Logins":         loginsEncry,
			"AwsConsoles":    awsConsoleEncry,
			"AwsApiAccesses": awsApiAccessEncry,
		},
	}, "", "\t")
}

// marshalAccount
func (space Space) marshalAccount(accs map[string]account.Serializer) ([]byte, error) {
	out := make(map[string]interface{})

	for key, acc := range accs {
		b, err := acc.Serialize()
		if err != nil {
			return nil, err
		}

		encrypted, err := space.guardian.Encrypt("", b) // where do we get the passphrase ??
		if err != nil {
			return nil, err
		}
		out[key] = encrypted
	}

	return json.MarshalIndent(out, "", "\t")
}
