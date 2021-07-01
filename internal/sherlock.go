package internal

import (
	"context"
	"fmt"
)

const (
	skippSetupFor = "setup"
)

var (
	ErrNotSetup    = fmt.Errorf("sherlock needs to set-up (use sherlock setup)")
	ErrNoSuchGroup = fmt.Errorf("provided group cannot be found (use sherlock add --group)")
)

// FileSystem declares the functions sherlock requires to
// interact with the underlying file system
type FileSystem interface {
	InitFs() error
	CreateGroup(name string) error
	GroupExists(name string) error
	ReadGroupVault(group string) ([]byte, error)
	WriteAccount(account *Account) error
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

// Setup checks if a main password for the vault has already been
// set which is required for every further command. Setup will create required directories
// if those are missing
func (sh *Sherlock) Setup(calledCmd string) error {
	if err := sh.fileSystem.InitFs(); err != nil {
		return err
	}
	return nil
}

func (sh Sherlock) SetupGroup(name string) error {
	if err := sh.GroupExists(name); err != nil {
		return err
	}
	return sh.fileSystem.CreateGroup(name)
}

func (sh Sherlock) GroupExists(name string) error {
	return sh.fileSystem.GroupExists(name)
}

func (sh *Sherlock) AddAccount(ctx context.Context, account *Account) error {
	return nil
}
