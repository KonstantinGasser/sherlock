package internal

import (
	"testing"

	"github.com/KonstantinGasser/sherlock/fs"
	"github.com/spf13/afero"
)

func memLock() *Sherlock {
	return &Sherlock{
		fileSystem: fs.New(afero.NewMemMapFs()),
	}
}

// TestSetup testis if the in-mem fs is setup (which will not be the case)
// and then sets up sherlock
func TestSetup(t *testing.T) {
	sh := memLock()
	if err := sh.IsSetUp(); err == nil {
		t.Fatalf("sherlock.IsSetup: want: nil (not-setup), have: %v", err)
	}

	err := sh.Setup("default_group_key")
	if err != nil {
		t.Fatalf("sherlock.Setup: want: nil, have: %v", err)
	}
}

func TestOptAccPassword(t *testing.T) {
	tt := []struct {
		g        Group
		accName  string
		newPass  string
		insecure bool
		ok       bool
	}{
		{
			g: Group{
				GID: "test1",
				Accounts: []*Account{
					{
						Name:     "test-acc1",
						Password: "hello-world",
					},
				},
			},
			accName:  "test-acc1",
			newPass:  "helloworld",
			insecure: false,
			ok:       false,
		},
		{
			g: Group{
				GID: "test2",
				Accounts: []*Account{
					{
						Name:     "test-acc2",
						Password: "hello-world",
					},
				},
			},
			accName:  "test-acc2",
			newPass:  "insecure-password",
			insecure: true,
			ok:       true,
		},
		{
			g: Group{
				GID: "test3",
				Accounts: []*Account{
					{
						Name:     "test-acc3",
						Password: "hello-world",
					},
				},
			},
			accName:  "test-acc3",
			newPass:  "$wsert-2w345_2@34#!0?",
			insecure: false,
			ok:       true,
		},
	}

	for _, tc := range tt {
		err := OptAccPassword(tc.newPass, tc.insecure)(&tc.g, tc.accName)
		if (err != nil && tc.ok) || (err == nil && !tc.ok) {
			t.Fatalf("internal.OptAccPassword: want:updated==%v, have:err==%v", tc.ok, err)
		}
		if tc.ok && tc.newPass != tc.g.Accounts[0].Password {
			t.Fatalf("internal.OptAccPassword: want: %s, have: %s", tc.newPass, tc.g.Accounts[0].Password)
		}
	}
}

func TestOptAccName(t *testing.T) {
	tt := []struct {
		g       Group
		accName string
		newName string
		err     error
	}{
		{
			g: Group{
				GID: "test1",
				Accounts: []*Account{
					{
						Name: "test-acc1",
					},
				},
			},
			accName: "test-acc1",
			newName: "test-acc1_1",
			err:     nil,
		},
		{
			g: Group{
				GID: "test2",
				Accounts: []*Account{
					{
						Name: "test-acc2",
					},
				},
			},
			accName: "test-acc2_not_found",
			newName: "test-acc2_2",
			err:     ErrNoSuchAccount,
		},
	}

	for _, tc := range tt {
		err := OptAccName(tc.newName)(&tc.g, tc.accName)
		if err != tc.err {
			t.Fatalf("internal.OptAccName: want: %s, have: %s", tc.err, err)
		}
		if err == nil && (tc.newName != tc.g.Accounts[0].Name) {
			t.Fatalf("internal.OptAccName: want: %s, have: %s", tc.newName, tc.g.Accounts[0].Name)
		}
	}
}

func TestOptAccDelete(t *testing.T) {
	tt := []struct {
		g           Group
		toBeDeleted string
		err         error
	}{
		{
			g: Group{
				GID: "test2",
				Accounts: []*Account{
					{
						Name: "test1",
					},
				},
			},
			toBeDeleted: "test1",
			err:         nil,
		},
	}

	for _, tc := range tt {
		err := OptAccDelete()(&tc.g, tc.toBeDeleted)
		if err != tc.err {
			t.Fatalf("internal.OptAccDelete: want: %v, have: %v", tc.err, err)
		}
		if tc.err != nil {
			if ok := tc.g.exists(tc.toBeDeleted); !ok {
				t.Fatalf("internal.OptAccDelete: account deleted but had an error: want:delete==%v, have:delete==%v", tc.err == nil, ok)
			}
		}
		if tc.err == nil {
			if ok := tc.g.exists(tc.toBeDeleted); ok {
				t.Fatalf("internal.OptAccDelete: account not deleted: want:delete==%v, have:delete==%v", tc.err == nil, ok)
			}
		}

	}
}
