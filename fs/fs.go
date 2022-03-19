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
func (fs Filesystem) Init(key string, s Serializer) error {

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
	return fs.Write(key, s)
}

func (fs Filesystem) IsSetup() error {
	defaultPath, err := spacepath(defaultSpace)
	if err != nil {
		return err
	}
	// make sure directories required exists
	if !exists(defaultPath) {
		return fmt.Errorf("default space not found")
	}

	// make sure default space is not empty (TODO: check if default space is ok and not corrupted)
	b, err := os.ReadFile(filepath.Join(defaultPath, spacefile))
	if err != nil {
		return fmt.Errorf("could not read default space")
	}

	var tmp map[string]interface{} // used to check if file has valid JSON
	if err := json.Unmarshal(b, &tmp); err != nil {
		return fmt.Errorf(`
default space looks corrupted...
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
		`, filepath.Join(defaultPath, spacefile))

	}
	return nil
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
