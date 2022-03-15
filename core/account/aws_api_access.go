package account

import "encoding/json"

type AwsApiAccess struct {
	AccessID  string
	AccessKey string
	Region    string
}

func DefaultAwsApiAccess() Account { return new(AwsApiAccess) }

func (awsA AwsApiAccess) Serialize() ([]byte, error) {
	return json.Marshal(awsA)
}
