package cmd

import (
	"context"

	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/internal/terminal"
	"github.com/spf13/cobra"
)

type delOptions struct {
	gid     string
	account string
	force   bool
}

func cmdDel(sherlock *internal.Sherlock) *cobra.Command {
	var opts delOptions
	del := &cobra.Command{
		Use:   "del",
		Short: "delete an account from a group",
		Long:  "delete an account from a group",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if opts.account == "" || opts.gid == "" {
				terminal.Error("account name and group are required (--account, --gid)")
				return
			}

			groupKey, err := terminal.ReadPassword("group (%s) password: ", opts.gid)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			if !opts.force {
				confirm := terminal.YesNo("delete account [y/N]: ")
				if !confirm {
					return
				}
			}

			if err := sherlock.DeleteAccount(context.Background(), opts.gid, opts.account, groupKey); err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.Success("account %q successfully deleted from %q", opts.account, opts.gid)
		},
	}
	del.Flags().StringVarP(&opts.gid, "gid", "g", "default", "group from which to delete the account")
	del.Flags().StringVarP(&opts.account, "account", "a", "", "account name to delete")
	del.Flags().BoolVarP(&opts.force, "force", "f", false, "will bypass [y/N] prompt")

	return del
}
