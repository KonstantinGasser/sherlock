package space

import "encoding/json"

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
	// Accounts stores all accounts mapped to the space
	Accounts *accounts
}

type accounts struct {
	// Logins holds all Login Types mapped to the space
	Logins map[string]interface{} // type not there yet
	// AwsConsoles holds all AWS-Console login Types mapped
	// to the space
	AwsConsoles map[string]interface{} // type not there yet
	// AwsApiKeys holds all AWS-Programmatic-Access identities
	// mapped to the space
	AwsApiKeys map[string]interface{} // type not there yet
}

func New(key string) *Space {
	return &Space{
		Key: key,
		Accounts: &accounts{
			Logins:      make(map[string]interface{}),
			AwsConsoles: make(map[string]interface{}),
			AwsApiKeys:  make(map[string]interface{}),
		},
	}
}

// Serialize marshales the content of a space
func (space Space) Serialize() ([]byte, error) {
	return json.MarshalIndent(space, "", "\t")
}
