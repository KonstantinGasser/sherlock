package core

import "github.com/KonstantinGasser/sherlock/core/space"

type Writer interface {
	Write(key string, space []byte) error
}

type Reader interface {
	// Read reads space-data into the passed in space.Space.
	Read(space *space.Space) error
}

type SherlockFS interface {
	Writer
	Reader
}

type Initializer interface {
	Initialize(key string, sapce []byte) error
}
