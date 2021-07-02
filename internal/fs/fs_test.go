package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

// TestInitFs checks if after success of func all the directories and files
// have been created
func TestInitFs(t *testing.T) {
	f := Fs{
		mock: afero.NewMemMapFs(),
	}

	err := f.InitFs([]byte("init-default-vault-content"))
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
	if len(defaultVault) != len([]byte("init-default-vault-content")) {
		t.Fatalf("fs.InitFs: saved vault differs from input vault. want: %d, have: %d", len([]byte("init-default-vault-content")), len(defaultVault))
	}
}
