package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

var (
	defaultInitVault = []byte("init-default-vault-content")
)

// TestInitFs checks if after success of func all the directories and files
// have been created
func TestInitFs(t *testing.T) {
	f := Fs{
		mock: afero.NewMemMapFs(),
	}

	err := f.InitFs(defaultInitVault)
	if err != nil {
		t.Fatalf("Fs.InitFs: want: nil, have: %v", err)
	}

	// check if all exists
	_, err = f.mock.Stat(filepath.Join(homepath(), sherlockRoot, groupsDir, defaultGroup))
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("fs.InitFs: default group dir not created")
		}
	}
	defaultVault, err := afero.ReadFile(f.mock, buildVaultPath(defaultGroup))
	if err != nil {
		t.Fatalf("fs.InitFs: could not open default group vault: %v", err)
	}
	if len(defaultVault) != len(defaultInitVault) {
		t.Fatalf("fs.InitFs: saved vault differs from input vault. want: %d, have: %d", len(defaultInitVault), len(defaultVault))
	}
}

func TestCreateGroup(t *testing.T) {
	var testGroup string = "test-group"

	f := Fs{
		mock: afero.NewMemMapFs(),
	}

	err := f.CreateGroup(testGroup, defaultInitVault)
	if err != nil {
		t.Fatalf("fs.CreateGroup: want: nil, have: %v", err)
	}

	// check if exists
	vault, err := afero.ReadFile(f.mock, buildVaultPath(testGroup))
	if err != nil {
		t.Fatalf("fs.CreateGroup: could not open test group vault: %v", err)
	}
	if len(vault) != len(defaultInitVault) {
		t.Fatalf("fs.CreateGroup: saved vault differs from input vault. want: %d, have: %d", len(defaultInitVault), len(vault))
	}
}
