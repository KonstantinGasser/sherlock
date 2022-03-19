package cmd

import (
	"context"

	"github.com/KonstantinGasser/sherlock/core"
	"github.com/KonstantinGasser/sherlock/out"
	"github.com/spf13/cobra"
)

func cmdSetup(ctx context.Context, initer core.Initializer) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "setup allows to initially set-up a main password for your vault",
		Long:  "to encrypt and decrypt your vault you will need to set-up a main password",
		Run: func(cmd *cobra.Command, args []string) {

			// check if already setup ????

			out.Info("please provide a password to encrypt the default sherlock space...")
			passphrase, err := out.ReadPassword("(default) space password: ")
			if err != nil {
				out.Error(err.Error())
				return
			}

			if err := core.InitSherlock(passphrase, initer); err != nil {
				out.Error(err.Error())
				return
			}

			out.Success("successfully set-up sherlock!")
			out.Banner()
		},
	}
}
