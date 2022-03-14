package account

import "encoding/json"

type Login struct {
	Key         string
	AccountName string
	Password    string
}

func NewLogin(key, accountName, password string) *Login {
	return &Login{
		Key:         key,
		AccountName: accountName,
		Password:    password,
	}
}

func (login Login) Serialize() ([]byte, error) {
	return json.Marshal(login)
}
