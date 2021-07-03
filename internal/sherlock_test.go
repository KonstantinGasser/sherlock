package internal

import (
	"testing"

	"github.com/KonstantinGasser/sherlock/fs"
	"github.com/spf13/afero"
)

func memLock() *Sherlock {
	return &Sherlock{
		fileSystem: fs.New(afero.NewMemMapFs()),
	}
}

// TestSetup testis if the in-mem fs is setup (which will not be the case)
// and then sets up sherlock
func TestSetup(t *testing.T) {
	sh := memLock()
	if err := sh.IsSetUp(); err == nil {
		t.Fatalf("sherlock.IsSetup: want: nil (not-setup), have: %v", err)
	}

	err := sh.Setup("default_group_key")
	if err != nil {
		t.Fatalf("sherlock.Setup: want: nil, have: %v", err)
	}
}
