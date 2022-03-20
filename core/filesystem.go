package core

import (
	"github.com/KonstantinGasser/sherlock/fs"
)

// Writer interface allows the any resource to
// persist data in the filesystem under a given key.
// However, the resource from which the content is provided
// must implement the fs.Serializer interface
type Writer interface {
	Write(key string, s fs.Serializer) error
}

// Reader interface describes a resource which
// requests information from the filesystem and
// receives the byte content of the requested information
type Reader interface {
	Read(key string) ([]byte, error)
}

// Mkdirer interface describes a resource which is able
// to create a directory
type Mkdirer interface {
	Mkdir(spaceKey string, opts ...func(path string) error) error
}

type SherlockFS interface {
	Writer
	Reader
	Mkdirer
	Init(key string, overwrite bool, s fs.Serializer) error
	IsSetup() error
	SpaceExists(key string) bool
}
