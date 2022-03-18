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
	// Logins map[string]account.Account // type not there yet
	Logins account.Account
	// AwsConsoles holds all AWS-Console login Types mapped
	// to the space
	AwsConsoles account.Account
	// AwsApiAccesses holds all AWS-Programmatic-Access identities
	// mapped to the space
	// I am not a fan of the map[string]account.Account type of thing..
	// However, go maps with custom interfaces ignore that concrete types implementing
	// the interface. It would be nice if Logins,Aws* could keep their type..
	// Think: re-implement logic to store accounts
	AwsApiAccesses account.Account
}

func New(key string) *Space {
	return &Space{
		Key:      key,
		guardian: &security.Guard{},
		Accounts: &Accounts{
			Logins:         new(account.Logins),
			AwsConsoles:    make(account.AwsConsoles),
			AwsApiAccesses: make(account.AwsApiAccesses),
		},
	}
}

func (space Space) ToCipherSpace(passphrase string) (*CipherSpace, error) {

	loginsEncry, err := space.Accounts.Logins.Encrypt(space.guardian, passphrase)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt Login Accounts: %v", err)
	}

	awsConsoleEncry, err := space.Accounts.AwsConsoles.Encrypt(space.guardian, passphrase)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt AWS-Console Accounts: %v", err)
	}

	awsApiAccessEncry, err := space.Accounts.AwsApiAccesses.Encrypt(space.guardian, passphrase)
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
		guardian: &security.Guard{},
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

	cSpace.guardian = &security.Guard{}
	return nil
}

// ToSpace asserts the CipherSpace to a regular Space
// All Accounts within the CipherSpace are decrypted
func (cSpace CipherSpace) ToSpace(passphrase string) (*Space, error) {

	logins, err := account.DefaultLogins().Decrypt(cSpace.guardian, cSpace.Accounts[jsonKeyLogin], passphrase)
	if err != nil {
		return nil, err
	}

	awsConsoles, err := account.DefaultAwsConsoles().Decrypt(cSpace.guardian, cSpace.Accounts[jsonKeyAwsConsole], passphrase)
	if err != nil {
		return nil, err
	}

	AwsApiAccesses, err := account.DefaultAwsAccesses().Decrypt(cSpace.guardian, cSpace.Accounts[jsonKeyAwsApiAccesses], passphrase)
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
