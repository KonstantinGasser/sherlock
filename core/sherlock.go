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
func InitSherlock(initer Initializer) error {
	defaultSpace := space.New(defaultNameSpace)

	b, err := defaultSpace.Serialize()
	if err != nil {
		return fmt.Errorf("could not serialize Space: %v", err)
	}
	return initer.Initialize(defaultSpace.Key, b)
}
