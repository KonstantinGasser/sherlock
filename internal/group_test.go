package internal

import (
	"testing"
)

func TestGroupAppend(t *testing.T) {

	tt := []struct {
		group    Group
		account  Account
		expected error
	}{
		{
			group: Group{
				GID: "test-group",
				// accounts in this case only have the account name since the rest
				// is not important for this test
				Accounts: []*Account{
					{
						Name: "test",
					},
				},
			},
			account: Account{
				Name: "some-other-account",
			},
			expected: nil,
		},
		{
			group: Group{
				GID: "test-group",
				// accounts in this case only have the account name since the rest
				// is not important for this test
				Accounts: []*Account{
					{
						Name: "same-account",
					},
				},
			},
			account: Account{
				Name: "same-account",
			},
			expected: ErrAccountExists,
		},
	}

	for _, tc := range tt {
		err := tc.group.append(&tc.account)
		if err != tc.expected {
			t.Fatalf("Group.append: want: %v, have: %v", tc.expected, err)
		}
	}
}
