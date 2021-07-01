package cmd

import (
	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/spf13/cobra"
)

const (
	skippSetupFor = "setup"
)

func RootCmd(sherlock *internal.Sherlock) *cobra.Command {
	root := &cobra.Command{
		Use:           "sherlock",
		Short:         "sherlock a CLI password manager for the simple use",
		Version:       "not there yet",
		SilenceUsage:  true,
		SilenceErrors: true,
		// ensure that sherlock is properly set-up. This means that the default group
		// exists and that it holds an encrypted .vault file. "sherlock setup" is excluded from this check
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == skippSetupFor {
				return nil
			}
			return sherlock.IsSetUp()
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	root.AddCommand(cmdSetup(sherlock))
	root.AddCommand(cmdAddAccount(sherlock))

	return root
}
