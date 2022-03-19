/*
Copyright © 2022 Konstantin Gasser konstantin.gasser@me.com

*/
package cmd

import (
	"context"

	"github.com/KonstantinGasser/sherlock/core"
	"github.com/spf13/cobra"
)

const (
	cmdSetupName = "setup"
)

func RootCommand(sh *core.Sherlock) *cobra.Command {
	root := &cobra.Command{
		Use:           "sherlock",
		Short:         "simple to use encrypted password and file CLI tool",
		Version:       "not-there-yet",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == cmdSetupName {
				return nil
			}

			if err := sh.IsSetup(); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	ctx := context.Background()
	root.AddCommand(cmdSetup(ctx, sh))
	root.AddCommand(cmdAdd(ctx, sh))
	return root
}
