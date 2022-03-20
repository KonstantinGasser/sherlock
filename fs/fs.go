package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

const (
	// basepath is the path in which all
	// sherlock related information is written/read to/from
	basepath = ".sherlock"
	// spacespath is the path in which all user created
	// spaces are stored, including the default space
	spacespath = "spaces"
	// spacefile is the file name in which a space is stored
	spacefile = ".space"
	// filespath is the path in which encrypted files within
	// a space will be stored
	filespaths = "files"

	// defaultSpace is the space which is created automatically on set-up
	defaultSpace = "default"
)

var (
	ErrNotSetup       = fmt.Errorf("file system with %s directory not setup", basepath)
	ErrNoSpaceFound   = fmt.Errorf("missing %s in space", spacefile)
	ErrCorruptedSpace = fmt.Errorf(`default space looks corrupted...
Overwrite the .space file at %s with this:

{
	"Key": "default",
	"Accounts": {
		"Logins": {},
		"AwsConsoles": {},
		"AwsApiAccesses": {}
	}
}

Or execute sherlock setup --overwrite
			`, filepath.Join("~/.sherlock/spaces/default", spacefile))
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

// Init initializes the required folder structure for sherlock
//
// under the sherlock-root `.sherlock` the folder `spaces` with a `default`
// space will be created. With the overwrite flag an already setup sherlock
// instance can be overwritten
func (fs Filesystem) Init(key string, overwrite bool, s Serializer) error {

	defaultPath, err := spacepath(key)
	if err != nil {
		return err
	}

	// do nothing if sherlock is already initialized
	if exists(defaultPath) && !overwrite {
		return nil
	}
	// ensure sherlock folder structure is created
	// 0770 only current user can read/write/execute
	if err := fs.fs.MkdirAll(defaultPath, 0770); err != nil {
		return fmt.Errorf("could not create sherlock folder structure: %v", err)
	}

	// write space in default namespace
	return fs.Write(key, s)
}

// IsSetup verifies the state of sherlock
//
// if all is ok (all directories and the default space exists in a ok condition)
// the function returns nil. In other cases IsSetup returns an error that directories
// are missing or that the default space is corrupted
func (fs Filesystem) IsSetup() error {
	defaultPath, err := spacepath(defaultSpace)
	if err != nil {
		return err
	}
	// make sure directories required exists
	if !exists(defaultPath) {
		return ErrNotSetup
	}

	// make sure default space is not empty
	// and content is valid JSON
	b, err := afero.ReadFile(fs.fs, filepath.Join(defaultPath, spacefile))
	if err != nil {
		return ErrNoSpaceFound
	}

	var tmp map[string]interface{} // used to check if file has valid JSON
	if err := json.Unmarshal(b, &tmp); err != nil {
		return ErrCorruptedSpace
	}
	return nil
}

func IrgnoreExisiting(path string) error {

	if _, err := os.Stat(path); err != nil {
		return nil
	}

	return fmt.Errorf("path already exists")
}

// Mkdir creates a new space directory with the permissions
// of 0700
func (fs Filesystem) Mkdir(space string, opts ...func(path string) error) error {
	path, err := spacepath(space)
	if err != nil {
		return err
	}

	for _, opt := range opts {
		if err := opt(path); err != nil {
			return fmt.Errorf("did not create %s: %v", path, err)
		}
	}

	return fs.fs.Mkdir(path, 0700)
}

func (fs Filesystem) Write(key string, s Serializer) error {
	path, err := spacepath(key)
	if err != nil {
		return err
	}
	data, err := s.Serialize()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(path, spacefile), data, 0700)
}

func (fs Filesystem) Read(key string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// spacepath builds the path in which a space
// can be found.
// While the base of the path is always the same
// the actual space directory depends on the space.Key
// example: /Users/username/.sherlock/spaces/<key>
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

func (fs Filesystem) SpaceExists(key string) bool {
	path, err := spacepath(key)
	if err != nil {
		return true // not good but ok for now; what does an error here mean?
	}
	return exists(path)
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
	}
	return true
}
