package internal

import (
	"testing"
	"time"
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
			name:     "gr@up@testaccount",
			password: "helloworld",
			tag:      "testing",
			insecure: true,
			created:  false,
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

func TestAccountExpiratin(t *testing.T) {

	tt := []struct {
		name     string
		a        account
		expected string
	}{
		{
			name: "expiration time 6 month from time.Now",
			a: account{
				UpdatedOn: time.Now(),
			},
			expected: "valid for 180 days 22 hours",
		},
		{
			name: "expired 6 month ago",
			a: account{
				UpdatedOn: time.Now().AddDate(-1, 0, 0),
			},
			expected: "expired 185 days 2 hours ago",
		},
	}

	for _, tc := range tt {
		expText := tc.a.Expiration()

		if expText != tc.expected {
			t.Fatalf("account.Expiration: want=%q, have: %q", tc.expected, expText)
		}
	}
}
