package core

import "github.com/KonstantinGasser/sherlock/fs"

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

type SherlockFS interface {
	Writer
	Reader
}

type Initializer interface {
	Initialize(key string, sapce []byte) error
}
