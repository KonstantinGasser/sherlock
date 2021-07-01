package cmd

import (
	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/internal/terminal"
	"github.com/spf13/cobra"
)

type listOptions struct {
	gid string
}

func cmdList(sherlock *internal.Sherlock) *cobra.Command {
	var opts listOptions
	list := &cobra.Command{
		Use:   "list",
		Short: "setup allows to initially set-up a main password for your vault",
		Long:  "to encrypt and decrypt your vault you will need to set-up a main password",
		Run: func(cmd *cobra.Command, args []string) {
			partitionKey, err := terminal.ReadPassword("partition password: ")
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			group, err := sherlock.LoadGroup(opts.gid, partitionKey)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.ToTable([]string{"Group", "Account", "Desc", "Created On"}, group.Table())
		},
	}
	list.Flags().StringVarP(&opts.gid, "group", "G", "default", "list accounts mapped to group")

	return list
}
