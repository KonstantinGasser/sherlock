package cmd

import (
	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/internal/terminal"
	"github.com/spf13/cobra"
)

func cmdSetup(sherlock *internal.Sherlock) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "setup allows to initially set-up a main password for your vault",
		Long:  "to encrypt and decrypt your vault you will need to set-up a main password",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sherlock.IsSetUp(); err == nil {
				terminal.Error("sherlock is already set-up")
				return
			}
			terminal.Success("sherlock has a default group for accounts not mapped to any group.\nPlease provide a group password for the default group.")

			groupKey, err := terminal.ReadPassword("(default) group password: ")
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			if err := sherlock.Setup(groupKey); err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.Success("sherlock successfully set-up")
		},
	}
}
