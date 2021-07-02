package cmd

import (
	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/internal/terminal"
	"github.com/spf13/cobra"
)

func cmdList(sherlock *internal.Sherlock) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "setup allows to initially set-up a main password for your vault",
		Long:  "to encrypt and decrypt your vault you will need to set-up a main password",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			groupKey, err := terminal.ReadPassword("group (%s) password: ", args[0])
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			group, err := sherlock.LoadGroup(args[0], groupKey)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.ToTable(
				[]string{"Group", "Account", "#Tag", "Created On", "Updated On"},
				group.Table(),
				terminal.TableWithCellMerge(0),
			)
		},
	}
}
