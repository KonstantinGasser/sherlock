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
		expected error
	}{
		{
			name:     "testaccount",
			password: "fsdf$35dfg0-43563sdf34",
			tag:      "testing",
			insecure: false,
			expected: nil,
		},
		{
			name:     "testaccount",
			password: "helloworld",
			tag:      "testing",
			insecure: false,
			expected: ErrInsecurePassword,
		},
		{
			name:     "test account",
			password: "helloworld",
			tag:      "testing",
			insecure: false,
			expected: ErrInvalidAccountName,
		},
		{
			name:     "testaccount",
			password: "helloworld",
			tag:      "testing",
			insecure: true,
			expected: nil,
		},
		{
			name:     "",
			password: "",
			tag:      "testing",
			insecure: false,
			expected: ErrMissingValues,
		},
	}

	for _, tc := range tt {
		_, err := NewAccount(tc.name, tc.password, tc.tag, tc.insecure)
		if err != tc.expected {
			t.Fatalf("internal.NewAccount: want: %v, have: %v", tc.expected, err)
		}
	}
}
