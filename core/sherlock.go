package core

import (
	"fmt"

	"github.com/KonstantinGasser/sherlock/core/account"
	"github.com/KonstantinGasser/sherlock/core/space"
)

const (
	defaultNameSpace = "default"
)

type Sherlock struct {
	// fs allows the sherlock type to
	// communicate with the underlaying filesystem
	fs SherlockFS
}

// NewSherlock returns a new Sherlock instance.
func NewSherlock(fs SherlockFS) (*Sherlock, error) {

	return &Sherlock{
		fs: fs,
	}, nil
}

// InitSherlock initializes the sherlock cli
//
// Initializing requires that there is a directory
// in the users home directory with a `.sherlock` folder.
// Further, InitSherlock ensures that a default space.Space
// is present - if not it creates one
func (sh Sherlock) Init(passphrase string) error {
	defaultSpace, err := space.New(defaultNameSpace).ToCipherSpace(passphrase)
	if err != nil {
		return err
	}

	return sh.fs.Init(defaultSpace.Key, defaultSpace)
}

func (sh Sherlock) IsSetup() error {
	return sh.fs.IsSetup()
}

func (sh Sherlock) AddSpace(passphrase string, key string) error {

	s, err := space.New(key).ToCipherSpace(passphrase)
	if err != nil {
		return fmt.Errorf("could not create space: %v", err)
	}

	return sh.fs.Write(key, s)
}

func (sh Sherlock) GetSpace(passphrase string, key string) (*space.Space, error) {

	b, err := sh.fs.Read(key)
	if err != nil {
		return nil, fmt.Errorf("could not read sherlock space: %v", err)
	}
	var cSpace = space.NewCipher(key)
	if err := cSpace.Deserialize(b); err != nil {
		return nil, fmt.Errorf("could not deserialize sherlock space: %v", err)
	}

	return cSpace.ToSpace(passphrase)
}

func (sh Sherlock) AddAccount(spaceKey string, passphrase string, acc account.Account) error {

	b, err := sh.fs.Read(spaceKey)
	if err != nil {
		return err
	}

	var cSpace = space.NewCipher(spaceKey)
	if err := cSpace.Deserialize(b); err != nil {
		return fmt.Errorf("could not deserialize space: %v", err)
	}

	space, err := cSpace.ToSpace(passphrase)
	if err != nil {
		return fmt.Errorf("could not decrypt space: %v", err)
	}
	space.Add(acc)

	return nil
}
