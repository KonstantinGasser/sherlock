/*
Copyright © 2022 Konstantin Gasser konstantin.gasser@me.com

*/
package cmd

import (
	"context"

	"github.com/KonstantinGasser/sherlock/core"
	"github.com/KonstantinGasser/sherlock/fs"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RootCommand(sh *core.Sherlock) *cobra.Command {
	root := &cobra.Command{
		Use:           "sherlock",
		Short:         "simple to use encrypted password and file CLI tool",
		Version:       "not-there-yet",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	ctx := context.Background()
	root.AddCommand(cmdSetup(ctx, fs.New(afero.NewOsFs())))
	root.AddCommand(cmdAdd(ctx, sh))
	return root
}
