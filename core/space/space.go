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
	guardian security.Encrypter
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

// To think about::
// 	Should the space be serializable, or should the space have a function which returns
// 	a CipherSpace which then can be (de-)serialized?
// 	This way one could work with a Space with no issues and would not need to worry about encrypting/decrypting?
// 		Process::
//			Write after modification
//			-> create Space
//				-> do something with it
//					-> call ToCipherSpace() (encrypt accounts) // implements Serializer interface
//						-> write CipherSpace to filesystem
//
//			Read for modification
//			-> read space bytes from filesystem
//				-> unmarshal into CipherSpace
//					-> Deserialize to Space (decrypt account(s))
//						-> do something with the Space

// Serialize marshal's the content of a space
//
// While the Space meta-data such as the space.Key will not be encrypted
// all accounts listed in the account type will be. However, the keys of each account
// will remain un-encrypted and human readable.
// As a result Serialize returns a byte slice of JSON-Space Object
func (space Space) Serialize() ([]byte, error) {

	/*
		encrypt all accounts within the space
	*/

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

// marshalAccount returns a serialized JSON-Object where each account is encrypted using the
// space passphrase
func (space Space) marshalAccount(passphrase string, accs map[string]account.Serializer) ([]byte, error) {
	out := make(map[string]interface{})

	for key, acc := range accs {
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

	return json.MarshalIndent(out, "", "\t")
}

// CipherSpace represents a space where all accounts
// are stored as bytes (encrypted). This type should be
// used when reading and unmarshaling a space from the
// filesystem
type CipherSpace struct {
	Key      string
	guardian security.Decrypter
	Accounts map[string]map[string][]byte
}

// ToSpace asserts the CipherSpace to a regular Space
// All Accounts within the CipherSpace are decrypted
func (cSpace CipherSpace) ToSpace(passphrase string) (*Space, error) {

	logins, err := cSpace.unmarshalAccount(passphrase, cSpace.Accounts["Logins"], nil)
	if err != nil {
		return nil, err
	}

	awsConsoles, err := cSpace.unmarshalAccount(passphrase, cSpace.Accounts["AwsConsoles"], nil)
	if err != nil {
		return nil, err
	}

	AwsApiAccesses, err := cSpace.unmarshalAccount(passphrase, cSpace.Accounts["AwsApiAccesses"], nil)
	if err != nil {
		return nil, err
	}

	return &Space{
		Key:      cSpace.Key,
		guardian: security.Guard{},
		Accounts: &accounts{
			Logins:         logins,
			AwsConsoles:    awsConsoles,
			AwsApiAccesses: AwsApiAccesses,
		},
	}, nil

}

func (cSpace CipherSpace) unmarshalAccount(passphrase string, accs map[string][]byte, newFunc func() interface{}) (map[string]interface{}, error) {

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
