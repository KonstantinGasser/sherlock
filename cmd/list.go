package cmd

import (
	"context"

	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/terminal"
	"github.com/enescakir/emoji"
	"github.com/spf13/cobra"
)

type listOptions struct {
	filterByTag string
	all         bool
	verbose     bool
}

func cmdList(ctx context.Context, sherlock *internal.Sherlock) *cobra.Command {
	var opts listOptions

	list := &cobra.Command{
		Use:   "list",
		Short: "list all accounts mapped to a given group",
		Long:  "with the list command you can inspect all accounts mapped to a given group",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var gid = "default"
			if opts.all {
				groupList, err := sherlock.ReadRegisteredGroups()
				if err != nil {
					terminal.Error(err.Error())
					return
				}
				terminal.Info("Registered Groups : ")
				for _, group := range groupList {
					terminal.SingleRow(emoji.RadioButton, group)
				}
				return
			} else if len(args) > 0 {
				gid = args[0]
			}
			groupKey, err := terminal.ReadPassword("(%s) password: ", gid)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			group, err := sherlock.LoadGroup(gid, groupKey)
			if err != nil {
				terminal.Error(err.Error())
				return
			}

			headers := []string{"Group", "Account", "#Tag", "Created On"}
			if opts.verbose {
				headers = append(headers, "Updated On", "Expires In")
			}

			terminal.ToTable(
				headers,
				group.Table(
					opts.verbose,
					internal.FilterByTag(opts.filterByTag),
				),
				terminal.TableWithCellMerge(0),
			)
		},
	}
	list.Flags().StringVarP(&opts.filterByTag, "tag", "t", "", "filter accounts by tag name")
	list.Flags().BoolVarP(&opts.all, "all", "a", false, "show all registered groups")
	list.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "display additional information")

	return list
}
