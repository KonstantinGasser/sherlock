package cmd

import (
	"context"

	sherlockErrors "github.com/KonstantinGasser/sherlock/errors"
	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/terminal"
	"github.com/spf13/cobra"
)

type listOptions struct {
	filterByTag string
}

func cmdList(ctx context.Context, sherlock *internal.Sherlock) *cobra.Command {
	var opts listOptions

	list := &cobra.Command{
		Use:   "list",
		Short: "list all accounts mapped to a given group",
		Long:  "with the list command you can inspect all accounts mapped to a given group",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var gid = "default"
			if len(args) > 0 {
				gid = args[0]
			}
			// load group to check whether the group is there
			err := sherlock.GroupExists(gid)
			if err != nil {
				shErr, ok := err.(*sherlockErrors.SherlockErr)
				if ok && shErr.SherlockErrTemplate == sherlockErrors.ErrGroupNotFound {
					err = shErr
				}
				terminal.Error(err.Error())
				return
			}
			groupKey, err := terminal.ReadPassword("(%s) password: ", gid)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			group, err := sherlock.LoadGroup(gid, groupKey)
			if err != nil {
				shErr, ok := err.(*sherlockErrors.SherlockErr)
				if ok && shErr.SherlockErrTemplate == sherlockErrors.ErrGroupNotFound {
					err = shErr
				}
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
