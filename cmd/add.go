/*
Copyright © 2022 Konstantin Gasser konstantin.gasser@me.com

*/
package cmd

import (
	"context"

	"github.com/KonstantinGasser/sherlock/core"
	"github.com/KonstantinGasser/sherlock/out"
	"github.com/spf13/cobra"
)

func cmdAdd(ctx context.Context, sh *core.Sherlock) *cobra.Command {
	add := &cobra.Command{
		Use:   "add",
		Short: "add allows to add resource to sherlock",
		Long:  `With the add command resources such as space, accounts and files can be added to sherlock`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	add.AddCommand(cmdAddSpace(ctx, sh))
	return add
}

func cmdAddSpace(ctx context.Context, sh *core.Sherlock) *cobra.Command {
	return &cobra.Command{
		Use:   "space",
		Short: "space adds a new space to sherlock",
		Args:  cobra.ExactArgs(1),
		Long:  `With the space command a new space gets added to sherlock under which accounts and files can be stored in a secure way`,
		Run: func(cmd *cobra.Command, args []string) {

			spaceName := args[0]

			passphrase, err := out.ReadPassword("set passphrase for space %q: ", spaceName)
			if err != nil {
				out.Error(err.Error())
				return
			}

			if err := sh.CreateSpace(passphrase, spaceName); err != nil {
				out.Error(err.Error())
				return
			}
		},
	}
}
