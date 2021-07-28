package internal

import (
	"testing"
)

func TestCreateGoup(t *testing.T) {
	tt := []struct {
		name   string
		expect error
	}{
		{
			name:   "test-group",
			expect: nil,
		},
		{
			name:   "",
			expect: ErrMissingValues,
		},
		{
			name:   "test group",
			expect: ErrInvalidGroupName,
		},
		{
			name:   "test@group",
			expect: ErrInvalidGroupNameSymbol,
		},
	}
	for _, tc := range tt {
		_, err := NewGroup(tc.name)
		if err != tc.expect {
			t.Fatalf("internal.NewGroup: want: %v, have: %v", tc.expect, err)
		}
	}
}

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

// This is a no-brainer test however it exists to ensure the function
// will not be changed unexpectedly
func TestFilterByTag(t *testing.T) {
	tt := []struct {
		account   Account
		filterTag string
		excpeted  bool
	}{
		{
			account:   Account{Tag: "tag_1"},
			filterTag: "tag_1",
			excpeted:  true,
		},
		{
			account:   Account{Tag: "tag_2"},
			filterTag: "tag_1",
			excpeted:  false,
		},
	}
	for _, tc := range tt {
		f := FilterByTag(tc.filterTag)
		if ok := f(&tc.account); ok != tc.excpeted {
			t.Fatalf("group.FilterByTag: want: %v, have: %v", tc.excpeted, ok)
		}
	}
}
