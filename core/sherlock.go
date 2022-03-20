package core

import (
	"fmt"

	"github.com/KonstantinGasser/sherlock/core/account"
	"github.com/KonstantinGasser/sherlock/core/space"
	"github.com/KonstantinGasser/sherlock/fs"
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

// Init initializes the sherlock cli
//
// Init requires that there is a directory
// in the users home directory with a `.sherlock` folder.
// Further, Init ensures that a default space.CipherSpace
// is present - if not it creates one
func (sh Sherlock) Init(passphrase string, overwrite bool) error {
	defaultSpace, err := space.New(defaultNameSpace).ToCipherSpace(passphrase)
	if err != nil {
		return err
	}

	return sh.fs.Init(defaultSpace.Key, overwrite, defaultSpace)
}

// IsSetup verifies the state of sherlock
//
// if IsSetup returns nil the required directories under .sherlock
// exists and the default space is not corrupted.
// Else returns the according error
func (sh Sherlock) IsSetup() error {
	return sh.fs.IsSetup()
}

func (sh Sherlock) HasSpace(key string) bool {
	return sh.fs.SpaceExists(key)
}

// AddSpace adds a new space to sherlock
//
// AddSpace ignores spaces which already exists
func (sh Sherlock) AddSpace(passphrase string, key string) error {

	s, err := space.New(key).ToCipherSpace(passphrase)
	if err != nil {
		return fmt.Errorf("could not create space: %v", err)
	}

	if err := sh.fs.Mkdir(key, fs.IrgnoreExisiting); err != nil {
		return fmt.Errorf("did not create space: %v", err)
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
