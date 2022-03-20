/*
Copyright © 2022 Konstantin Gasser konstantin.gasser@me.com

*/
package main

import (
	"os"

	"github.com/KonstantinGasser/sherlock/cmd"
	"github.com/KonstantinGasser/sherlock/core"
	"github.com/KonstantinGasser/sherlock/fs"
	"github.com/KonstantinGasser/sherlock/out"
	"github.com/spf13/afero"
)

func main() {

	shfs := fs.New(afero.NewOsFs())
	sh, err := core.NewSherlock(shfs)
	if err != nil {
		out.Error("could not create sherlock: %v", err)
		os.Exit(1)
	}

	if err := cmd.RootCommand(sh).Execute(); err != nil {
		out.Error(err.Error())
		os.Exit(1)
	}
}
