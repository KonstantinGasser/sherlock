package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/KonstantinGasser/sherlock/security"
)

const (
	querySplitPoint = "@"
)

var (
	ErrNotSetup     = fmt.Errorf("sherlock needs to bee set-up first (use sherlock setup)")
	ErrNoSuchGroup  = fmt.Errorf("provided group cannot be found (use sherlock add --group)")
	ErrWrongKey     = fmt.Errorf("wrong group key")
	ErrInvalidQuery = fmt.Errorf("invalid query. Query should be %q", "group@account")
)

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
func (sh Sherlock) SetupGroup(name string, groupKey string) error {
	if err := sh.GroupExists(name); err != nil {
		return err
	}
	group, err := NewGroup(name)
	if err != nil {
		return err
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
	serialized, err := group.serizalize()
	if err != nil {
		return err
	}
	encrypted, err := security.EncryptVault(serialized, groupKey)
	if err != nil {
		return err
	}
	if err := sh.fileSystem.Write(ctx, gid, encrypted); err != nil {
		return err
	}
	return nil
}

func (sh Sherlock) GetAccount(query string, groupKey string) (*Account, error) {
	keySet, err := sh.splitQuery(query)
	if err != nil {
		return nil, err
	}

	group, err := sh.LoadGroup(keySet[0], groupKey)
	if err != nil {
		return nil, err
	}
	return group.lookup(keySet[1])
}

func (sh Sherlock) DeleteAccount(ctx context.Context, group, account string, groupKey string) error {
	bytes, err := sh.fileSystem.ReadGroupVault(group)
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

	serialized, err := g.serizalize()
	if err != nil {
		return err
	}
	encrypted, err := security.EncryptVault(serialized, groupKey)
	if err != nil {
		return err
	}
	if err := sh.fileSystem.Write(ctx, group, encrypted); err != nil {
		return err
	}
	return nil
}

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

func (sh Sherlock) splitQuery(query string) ([]string, error) {
	set := strings.Split(query, querySplitPoint)
	if len(set) != 2 {
		return nil, ErrInvalidQuery
	}
	return set, nil
}
