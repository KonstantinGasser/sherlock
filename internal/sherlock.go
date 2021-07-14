package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/KonstantinGasser/sherlock/security"
)

const (
	// querySplitPoint refers to the command line argument coming from the user
	// in the form of group@account and the separator used for it
	querySplitPoint = "@"
)

var (
	ErrNotSetup     = fmt.Errorf("sherlock needs to bee set-up first (use sherlock setup)")
	ErrNoSuchGroup  = fmt.Errorf("provided group cannot be found (use sherlock add --group)")
	ErrWrongKey     = fmt.Errorf("wrong group key")
	ErrInvalidQuery = fmt.Errorf("invalid query. Query should be %q", "group@account")
)

type StateOption func(g *Group, gid, acc string) error

// OptAccPassword returns a StateOption to change an account password
func OptAccPassword(password string, insecure bool) StateOption {
	return func(g *Group, gid, acc string) error {
		account, err := g.lookup(acc)
		if err != nil {
			return err
		}
		if err := account.updatePassword(password, insecure); err != nil {
			return err
		}
		return nil
	}
}

// OptAccName returns a StateOption to change an account name
func OptAccName(name string) StateOption {
	return func(g *Group, gid, acc string) error {
		if ok := g.exists(name); ok {
			return ErrAccountExists
		}
		account, err := g.lookup(acc)
		if err != nil {
			return err
		}
		account.updateName(name)
		return nil
	}
}

// FileSystem declares the functions sherlock requires to
// interact with the underlying file system
type FileSystem interface {
	InitFs(initVault []byte) error
	CreateGroup(name string, initVault []byte) error
	GroupExists(name string) error
	VaultExists(group string) error
	ReadGroupVault(group string) ([]byte, error)
	Write(ctx context.Context, gid string, data []byte) error
}

type Sherlock struct {
	fileSystem FileSystem
}

// New return new Sherlock instance
func NewSherlock(fs FileSystem) *Sherlock {
	return &Sherlock{
		fileSystem: fs,
	}
}

func (sh Sherlock) IsSetUp() error {
	if err := sh.fileSystem.GroupExists("default"); err == nil { // default group does not exists
		return ErrNotSetup
	}
	if err := sh.fileSystem.VaultExists("default"); err == nil {
		return ErrNotSetup
	}
	return nil
}

// Setup checks if a main password for the vault has already been
// set which is required for every further command. Setup will create required directories
// if those are missing
func (sh *Sherlock) Setup(groupKey string) error {
	vault, err := security.InitWithDefault(groupKey, Group{
		GID:      "default",
		Accounts: make([]*Account, 0),
	})
	if err != nil {
		return err
	}

	if err := sh.fileSystem.InitFs(vault); err != nil {
		return err
	}
	return nil
}

// SetupGroup creates the group in the file system
// if the group does not already exists
func (sh Sherlock) SetupGroup(name string, groupKey string, insecure bool) error {
	if err := sh.GroupExists(name); err != nil {
		return err
	}
	group, err := NewGroup(name)
	if err != nil {
		return err
	}
	if !insecure {
		// check password strength for group key
		if err := group.secure(groupKey); err != nil {
			return err
		}
	}
	vault, err := security.InitWithDefault(groupKey, group)
	if err != nil {
		return err
	}
	return sh.fileSystem.CreateGroup(name, vault)
}

func (sh Sherlock) GroupExists(name string) error {
	return sh.fileSystem.GroupExists(name)
}

// AddAccount looks up the group-vault appending its accounts slice with the new account if the account does not
// yet exists
func (sh *Sherlock) AddAccount(ctx context.Context, account *Account, groupKey string, gid string) error {
	bytes, err := sh.fileSystem.ReadGroupVault(gid)
	if err != nil {
		return err
	}
	var group Group
	if err := security.DecryptVault(bytes, groupKey, &group); err != nil {
		return ErrWrongKey
	}
	if err := group.append(account); err != nil {
		return err
	}
	return sh.WriteGroup(ctx, gid, groupKey, &group)
}

// GetAccount looks up the requested account
// to locate an account the query needs to include the group
// like so group@account
func (sh Sherlock) GetAccount(query string, groupKey string) (*Account, error) {
	keySet, err := splitQuery(query)
	if err != nil {
		return nil, err
	}

	group, err := sh.LoadGroup(keySet[0], groupKey)
	if err != nil {
		return nil, err
	}
	return group.lookup(keySet[1])
}

// UpdateState executes the passed in StateOption to perform state changes on a group
func (sh Sherlock) UpdateState(ctx context.Context, query, groupKey string, opt StateOption) error {
	keySet, err := splitQuery(query)
	if err != nil {
		return err
	}

	group, err := sh.LoadGroup(keySet[0], groupKey)
	if err != nil {
		return err
	}
	if err := opt(group, keySet[0], keySet[1]); err != nil {
		return err
	}
	return sh.WriteGroup(ctx, keySet[0], groupKey, group)
}

// DeleteAccount deletes an account mapped to a group. If it is the last account in the group
// the group remains and will not get deleted
func (sh Sherlock) DeleteAccount(ctx context.Context, gid, account string, groupKey string) error {
	bytes, err := sh.fileSystem.ReadGroupVault(gid)
	if err != nil {
		return err
	}

	var g Group
	if err := security.DecryptVault(bytes, groupKey, &g); err != nil {
		return err
	}
	if err := g.delete(account); err != nil {
		return err
	}

	return sh.WriteGroup(ctx, gid, groupKey, &g)
}

// LoadGroup loads and decrypts the group vault
func (sh Sherlock) LoadGroup(gid string, groupKey string) (*Group, error) {
	bytes, err := sh.fileSystem.ReadGroupVault(gid)
	if err != nil {
		return nil, err
	}
	var group Group
	if err := security.DecryptVault(bytes, groupKey, &group); err != nil {
		return nil, ErrWrongKey
	}
	return &group, nil
}

// WriteGroup encrypts and write the group vault
func (sh Sherlock) WriteGroup(ctx context.Context, gid string, groupKey string, group *Group) error {
	serialized, err := group.serizalize()
	if err != nil {
		return err
	}
	encrypted, err := security.EncryptVault(serialized, groupKey)
	if err != nil {
		return err
	}
	return sh.fileSystem.Write(ctx, gid, encrypted)
}

// splitQuery verifies that a query (for get,update command) are in the correct
// format: group@account
func splitQuery(query string) ([]string, error) {
	set := strings.Split(query, querySplitPoint)
	if len(set) != 2 {
		return nil, ErrInvalidQuery
	}
	return set, nil
}
