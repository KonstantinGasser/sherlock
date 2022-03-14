package core

import (
	"fmt"

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

	return nil, fmt.Errorf("not implemented yet")
}

// InitSherlock initializes the sherlock cli
//
// Initializing requires that there is a directory
// in the users home directory with a `.sherlock` folder.
// Further, InitSherlock ensures that a default space.Space
// is present - if not it creates one
func InitSherlock(passphrase string, initer Initializer) error {
	defaultSpace, err := space.New(defaultNameSpace).ToCipherSpace(passphrase)
	if err != nil {
		return err
	}

	b, err := defaultSpace.Serialize()
	if err != nil {
		return fmt.Errorf("could not serialize Space: %v", err)
	}
	return initer.Initialize(defaultSpace.Key, b)
}

func (sh Sherlock) CreateSpace(passphrase string, key string) error {

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
