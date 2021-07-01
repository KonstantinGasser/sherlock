package fs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	sherlockRoot  = ".sherlock"
	groupsDir     = "groups"
	defaultGroup  = "default"
	vaultFileName = ".vault"
)

var (
	ErrNoSuchGroup = fmt.Errorf("group not found in sherlock")
	ErrNoSuchVault = fmt.Errorf("vault for group not found in sherlock")
	ErrGroupExists = fmt.Errorf("group already exists")
)

type Fs struct{}

func New() *Fs {
	return &Fs{}
}

// ReadVault reads the stored .vault file
func (fs Fs) ReadGroupVault(group string) ([]byte, error) {
	return ioutil.ReadFile(buildVaultPath(group))
}

// InitFs creates all directories required to be setup to use
// sherlock. If the directory exists nothing happens
func (fs Fs) InitFs(initVault []byte) error {
	if err := os.MkdirAll(filepath.Join(homepath(), sherlockRoot, groupsDir, defaultGroup), 0777); err != nil {
		return err
	}

	f, err := os.OpenFile(buildVaultPath(defaultGroup), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, bytes.NewReader(initVault)); err != nil {
		return err
	}
	return nil
}

func (fs Fs) CreateGroup(name string, initVault []byte) error {
	if err := os.MkdirAll(filepath.Join(homepath(), sherlockRoot, groupsDir, name), 0777); err != nil {
		return err
	}
	f, err := os.OpenFile(buildVaultPath(name), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, bytes.NewReader(initVault)); err != nil {
		return err
	}
	return nil
}

func (fs Fs) GroupExists(name string) error {
	_, err := os.Stat(buildGroupPath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return ErrGroupExists
}

func (fs Fs) VaultExists(group string) error {
	_, err := os.Stat(buildVaultPath(group))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return ErrNoSuchVault
}

func (fs Fs) Write(ctx context.Context, gid string, data []byte) error {
	if err := ioutil.WriteFile(buildVaultPath(gid), data, os.ModeAppend); err != nil {
		return err
	}
	return nil
}

func buildGroupPath(group string) string {
	return filepath.Join(homepath(), sherlockRoot, groupsDir, group)
}

// buildVaultPath creates a file path like
// => $HOME/.sherlock/groups/{group}/.vault
func buildVaultPath(group string) string {
	return filepath.Join(homepath(), sherlockRoot, groupsDir, group, vaultFileName)
}

func homepath() string {
	home, _ := os.UserHomeDir()
	return home
}
