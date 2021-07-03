package cmd

import (
	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/terminal"
	"github.com/spf13/cobra"
)

type listOptions struct {
	filterByTag string
}

func cmdList(sherlock *internal.Sherlock) *cobra.Command {
	var opts listOptions

	list := &cobra.Command{
		Use:   "list",
		Short: "setup allows to initially set-up a main password for your vault",
		Long:  "to encrypt and decrypt your vault you will need to set-up a main password",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var gid = "default"
			if len(args) > 0 {
				gid = args[0]
			}
			groupKey, err := terminal.ReadPassword("group (%s) password: ", gid)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			group, err := sherlock.LoadGroup(gid, groupKey)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.ToTable(
				[]string{"Group", "Account", "#Tag", "Created On", "Updated On"},
				group.Table(
					internal.FilterByTag(opts.filterByTag),
				),
				terminal.TableWithCellMerge(0),
			)
		},
	}
	list.Flags().StringVarP(&opts.filterByTag, "tag", "t", "", "filter accounts by tag name")

	return list
}
