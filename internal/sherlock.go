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
	ErrNoSuchGroup  = fmt.Errorf("provided group cannot be found (use sherlock add group)")
	ErrWrongKey     = fmt.Errorf("wrong group key")
	ErrInvalidQuery = fmt.Errorf("invalid query. Query should be %q", "group@account")
)

// StateOption describes a function with can alter the state of
// a group.
type StateOption func(g *group, acc string) error

// OptAddAccount returns a StateOption allowing to append
// an account to an group
func OptAddAccount(account *account) StateOption {
	return func(g *group, acc string) error {
		return g.append(account)
	}
}

// OptAccPassword returns a StateOption to change
// an account password
func OptAccPassword(password string, insecure bool) StateOption {
	return func(g *group, acc string) error {
		account, err := g.lookup(acc)
		if err != nil {
			return err
		}
		if err := account.update(updateFieldPassword(password, insecure)); err != nil {
			return err
		}
		return nil
	}
}

// OptAccName returns a StateOption to change
// an account name
func OptAccName(name string) StateOption {
	return func(g *group, acc string) error {
		if ok := g.exists(name); ok {
			return ErrAccountExists
		}
		account, err := g.lookup(acc)
		if err != nil {
			return err
		}
		if err := account.update(updateFieldName(name)); err != nil {
			return err
		}
		return nil
	}
}

// OptsAccTag returns a StateOption with allows
// to change the tag field of an account
func OptsAccTag(tag string) StateOption {
	return func(g *group, acc string) error {
		account, err := g.lookup(acc)
		if err != nil {
			return err
		}
		if err := account.update(updateFieldTag(tag)); err != nil {
			return err
		}
		return nil
	}
}

// OptAccDelete returns a StateOption deleting
// an account if it exists
func OptAccDelete() StateOption {
	return func(g *group, acc string) error {
		return g.delete(acc)
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
	Delete(ctx context.Context, gid string) error
	Write(ctx context.Context, gid string, data []byte) error
	ReadRegisteredGroups() ([]string, error)
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

// IsSetUp verifies that sherlock is operational
//
// to be operational there must be a folder $HOME/.sherlock/group
// with the default group and an encrypted default vault for which
// the user has set a group password.
func (sh Sherlock) IsSetUp() error {
	if err := sh.fileSystem.GroupExists("default"); err == nil { // default group does not exists
		return ErrNotSetup
	}
	if err := sh.fileSystem.VaultExists("default"); err == nil {
		return ErrNotSetup
	}
	return nil
}

// Setup sets the sherlock environment up
//
// the env requires to have an default group with an encrypted default vault
// sitting in $HOME/.sherlock/group.
func (sh *Sherlock) Setup(groupKey string) error {
	vault, err := security.InitWithDefault(groupKey, newDefaultGroup())
	if err != nil {
		return err
	}

	if err := sh.fileSystem.InitFs(vault); err != nil {
		return err
	}
	return nil
}

// DeleteGroup irreversible deletes a group from sherlock
// and the underlying file-system
func (sh *Sherlock) DeleteGroup(ctx context.Context, gid string) error {
	return sh.fileSystem.Delete(ctx, gid)
}

// SetupGroup creates a new group in sherlock
//
// a group creation will be rejected if the GID already
// exits, or the groupKey is to weak (if !insecure). The created group
// will be initialized with an encrypted default vault.
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

// GroupExists looks up if a group exists within sherlock
func (sh Sherlock) GroupExists(name string) error {
	return sh.fileSystem.GroupExists(name)
}

// CheckGroupKey performs a pre-check to check if groupKey is correct
//
// it should only be used to verify that an inputed groupKey by the user
// is matching the group the user is trying to access.
func (sh *Sherlock) CheckGroupKey(ctx context.Context, query, groupKey string) error {
	gid, _, err := SplitQuery(query)
	if err != nil {
		return err
	}
	bytes, err := sh.fileSystem.ReadGroupVault(gid)
	if err != nil {
		return err
	}
	var g group
	if err := security.Decrypt(bytes, groupKey, &g); err != nil {
		return ErrWrongKey
	}
	return nil
}

// GetAccount looks up the requested account
//
// the lookup is performed through the query (group@account).
func (sh Sherlock) GetAccount(query string, groupKey string) (*account, error) {
	gid, name, err := SplitQuery(query)
	if err != nil {
		return nil, err
	}

	group, err := sh.LoadGroup(gid, groupKey)
	if err != nil {
		return nil, err
	}
	return group.lookup(name)
}

// UpdateState executes the passed in StateOption to perform state changes on a group
//
// it allows to modify a group/account (adding accounts, changing account) through the passed StateOption.
func (sh Sherlock) UpdateState(ctx context.Context, query, groupKey string, opt StateOption) error {
	gid, name, err := SplitQuery(query)
	if err != nil {
		return err
	}

	group, err := sh.LoadGroup(gid, groupKey)
	if err != nil {
		return err
	}
	if err := opt(group, name); err != nil {
		return err
	}
	return sh.writeGroup(ctx, gid, groupKey, group)
}

// LoadGroup loads a group
//
// it wraps the reading of the group and the decryption
// functions together
func (sh Sherlock) LoadGroup(gid string, groupKey string) (*group, error) {
	bytes, err := sh.fileSystem.ReadGroupVault(gid)
	if err != nil {
		return nil, err
	}
	var g group
	if err := security.Decrypt(bytes, groupKey, &g); err != nil {
		return nil, ErrWrongKey
	}
	return &g, nil
}

// writeGroup saves a group in sherlock
//
// it wraps the encryption and the writing of a group together
// in one function
func (sh Sherlock) writeGroup(ctx context.Context, gid string, groupKey string, group *group) error {
	serialized, err := group.serizalize()
	if err != nil {
		return err
	}
	encrypted, err := security.Encrypt(serialized, groupKey)
	if err != nil {
		return err
	}
	return sh.fileSystem.Write(ctx, gid, encrypted)
}

// SplitQuery separates the user query into it pieces (group, account)
//
// quires not following the format will result in a ErrInvalidQuery error
// format: group@account
func SplitQuery(query string) (string, string, error) {
	set := strings.Split(query, querySplitPoint)
	if len(set) != 2 {
		return "", "", ErrInvalidQuery
	}
	return set[0], set[1], nil
}

func NameValidation(name string) bool {
	return !strings.Contains(name, querySplitPoint)
}

// ReadRegisteredGroups loads saved groups
func (sh Sherlock) ReadRegisteredGroups() ([]string, error) {
	groups, err := sh.fileSystem.ReadRegisteredGroups()
	if err != nil {
		return nil, err
	}
	return groups, nil
}
