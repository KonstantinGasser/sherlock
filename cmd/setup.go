package cmd

import (
	"context"
	"errors"

	"github.com/KonstantinGasser/sherlock/core"
	"github.com/KonstantinGasser/sherlock/fs"
	"github.com/KonstantinGasser/sherlock/out"
	"github.com/spf13/cobra"
)

func cmdSetup(ctx context.Context, sh *core.Sherlock) *cobra.Command {
	var overwrite bool
	setup := &cobra.Command{
		Use:   "setup",
		Short: "setup allows to initially set-up a main password for your vault",
		Long:  "to encrypt and decrypt your vault you will need to set-up a main password",
		Run: func(cmd *cobra.Command, args []string) {

			// if overwrite is set disregard the fact that sherlock
			// might be setup and has accounts in the default space
			if overwrite {
				out.Warning("using the --overwrite flag will cause the default space to be deleted\n and initialized with an empty space and a new passphrase")

				if ok := out.YesNo("overwrite current default space [y/N]: "); !ok {
					out.Info("canceled overwrite of default space. sherlock was not changed")
					return
				}
				passphrase, err := out.ReadPassword("(default) space password: ")
				if err != nil {
					out.Error(err.Error())
					return
				}
				if err := sh.Init(passphrase, overwrite); err != nil {
					out.Error(err.Error())
					return
				}

				out.Success("successfully overwritten default space")
				return
			}

			setupErr := sh.IsSetup()

			// ignore request if sherlock is already setup
			if setupErr == nil {
				out.Warning("sherlock is already setup")
				return
			}

			if errors.Is(setupErr, fs.ErrNoSpaceFound) || errors.Is(setupErr, fs.ErrCorruptedSpace) {
				out.Warning("while checking if sherlock is setup we found the following issue:")
				out.Error(setupErr.Error())
				out.Info("to fix the issue use `sherlock setup --overwrite` to ensure the default space is working correctly")
				return
			}

			out.Info("please provide a password to encrypt the default sherlock space...")
			passphrase, err := out.ReadPassword("(default) space password: ")
			if err != nil {
				out.Error(err.Error())
				return
			}

			if err := sh.Init(passphrase, overwrite); err != nil {
				out.Error(err.Error())
				return
			}

			out.Success("successfully set-up sherlock!")
			out.Banner()
		},
	}

	setup.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "if set will delete the current default space and creates a new empty one")
	return setup
}
