/*
Copyright © 2022 Konstantin Gasser konstantin.gasser@me.com

*/
package main

import (
	"os"

	"github.com/KonstantinGasser/sherlock/cmd"
)

type (
	Aer interface {
		A()
	}
	B struct{}
)

func (b B) A() {}

func main() {

	if err := cmd.RootCommand(nil).Execute(); err != nil {
		os.Exit(1)
	}
}
