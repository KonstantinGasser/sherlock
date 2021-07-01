package sherlock

import "context"

// FileSystem declares the functions sherlock requires to
// interact with the underlying file system
type FileSystem interface{}

type Sherlock struct {
	fileSystem FileSystem
}

// New return new Sherlock instance
func New(fs FileSystem) *Sherlock {
	return &Sherlock{
		fileSystem: fs,
	}
}

func (sh *Sherlock) Add(ctx context.Context) error {

}
