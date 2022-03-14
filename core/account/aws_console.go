package account

import "encoding/json"

type AwsConsole struct {
	AccountID   [12]byte
	AccountName string
	Password    string
}

func (awsC AwsConsole) Serialize() ([]byte, error) {
	return json.Marshal(awsC)
}
