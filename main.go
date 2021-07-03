package main

import (
	"github.com/KonstantinGasser/sherlock/cmd"
	"github.com/KonstantinGasser/sherlock/fs"
	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/terminal"
	"github.com/spf13/afero"
)

func main() {
	fileSystem := fs.New(afero.NewOsFs())
	sherlock := internal.NewSherlock(fileSystem)

	if err := cmd.RootCmd(sherlock).Execute(); err != nil {
		terminal.Error("%s", err)
	}
}
