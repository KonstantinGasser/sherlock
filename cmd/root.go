/*
Copyright © 2022 Konstantin Gasser konstantin.gasser@me.com

*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func RootCommand(notThereYet interface{}) *cobra.Command {
	root := &cobra.Command{
		Use:           "sherlock",
		Short:         "simple to use encrypted password and file CLI tool",
		Version:       "not-there-yet",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not yet implemented")
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	ctx := context.Background()
	root.AddCommand(cmdSetup(ctx, nil))
	root.AddCommand(cmdAdd(ctx, nil))
	return root
}
