package account

import "encoding/json"

type AwsConsole struct {
	AccountID   [12]byte
	AccountName string
	Password    string
}

func DefaultAwsConsole() Account { return new(AwsConsole) }

func NewAwsConsole(id [12]byte, name, password string) *AwsConsole {
	return &AwsConsole{
		AccountID:   id,
		AccountName: name,
		Password:    password,
	}
}

func (awsC AwsConsole) Serialize() ([]byte, error) {
	return json.Marshal(awsC)
}
