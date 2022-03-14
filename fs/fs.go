package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

const (
	// basepath is the path in which all
	// sherlock related information is written/read to/from
	basepath   = ".sherlock"
	spacespath = "spaces"
)

type Serializer interface {
	Serialize() ([]byte, error)
}

type Deserializer interface {
	Deserialize(v []byte) error
}

type Filesystem struct {
	// using afero's implementation of
	// the filesystems allows for easier and
	// cleaner testing
	fs afero.Fs
}

// New returns the filesystem implementation required by
// sherlock.
//
// The Filesystem implements the Initializer, Writer and Reader
// interfaces
// For testing proposes the filesystem can be created with
// an in-memory filesystem using the afero.NewMemMapFs
func New(fs afero.Fs) *Filesystem {
	return &Filesystem{
		fs: fs,
	}
}

// Initialize initializes the required folder structure for sherlock
//
// under the sherlock-root `.sherlock` the folder `spaces` with a `default`
// space will be created. If
func (fs Filesystem) Initialize(key string, space []byte) error {

	defaultPath, err := spacepath(key)
	if err != nil {
		return err
	}

	// do nothing if sherlock is already initialized
	if exists(defaultPath) {
		return nil
	}
	// ensure sherlock folder structure is created
	// 0770 only current user can read/write/execute
	if err := os.MkdirAll(defaultPath, 0770); err != nil {
		return fmt.Errorf("could not create sherlock folder structure: %v", err)
	}

	// write space in default namespace
	return fs.Write(key, space)
}

func (fs Filesystem) Write(key string, space []byte) error {
	return fmt.Errorf("not implemented yet")
}

func (fs Filesystem) Read(key string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// spacepath builds the path in which a space
// can be found.
// While the base of the path is always the same
// the actual space directory depends on the space.Key
func spacepath(key string) (string, error) {
	home, err := userhome()
	if err != nil {
		return "", fmt.Errorf("could not locate user home direcotry: %v", err)
	}
	return filepath.Join(home, basepath, spacespath, key), nil
}

func userhome() (string, error) {
	// return os.UserHomeDir()
	return "./", nil // for testing
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
	}
	return true
}
