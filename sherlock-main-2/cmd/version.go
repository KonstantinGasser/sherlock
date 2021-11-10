package cmd

import (
	"github.com/KonstantinGasser/sherlock/terminal"
	"github.com/spf13/cobra"
)

var Version = "dev"

func cmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "display sherlock version",
		Long:  "display sherlock version",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			terminal.Version(Version)
		},
	}
}
