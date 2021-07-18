package internal

import (
	"testing"
)

func TestNewAccount(t *testing.T) {

	tt := []struct {
		name     string
		password string
		tag      string
		insecure bool
		created  bool
	}{
		{
			name:     "group@testaccount",
			password: "fsdf$35dfg0-43563sdf34",
			tag:      "testing",
			insecure: false,
			created:  true,
		},
		{
			name:     "group@testaccount",
			password: "helloworld",
			tag:      "testing",
			insecure: false,
			created:  false,
		},
		{
			name:     "group@test account",
			password: "helloworld",
			tag:      "testing",
			insecure: false,
			created:  false,
		},
		{
			name:     "group@testaccount",
			password: "helloworld",
			tag:      "testing",
			insecure: true,
			created:  true,
		},
		{
			name:     "",
			password: "",
			tag:      "testing",
			insecure: false,
			created:  false,
		},
	}

	for _, tc := range tt {
		a, err := NewAccount(tc.name, tc.password, tc.tag, tc.insecure)
		if (tc.created && a == nil) || (!tc.created && a != nil) {
			t.Fatalf("internal.NewAccount: want:created==%v, have: error==%v", tc.created, err)
		}
	}
}
