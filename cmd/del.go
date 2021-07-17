package cmd

import (
	"context"

	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/terminal"
	"github.com/spf13/cobra"
)

func cmdDel(ctx context.Context, sherlock *internal.Sherlock) *cobra.Command {
	del := &cobra.Command{
		Use:   "del",
		Short: "delete a group or account from sherlock",
		Long:  "delete a group or account from sherlock",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	del.AddCommand(cmdDelAccount(ctx, sherlock))

	return del
}

type delAccOptions struct {
	force bool
}

func cmdDelAccount(ctx context.Context, sherlock *internal.Sherlock) *cobra.Command {
	var opts delAccOptions
	del := &cobra.Command{
		Use:   "account",
		Short: "delete an account from a group",
		Long:  "delete an account from a group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				terminal.Error("account key required (group@account)")
				return
			}

			groupKey, err := terminal.ReadPassword("(%s) password: ", args[0])
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

			if err := sherlock.UpdateState(ctx, args[0], groupKey, internal.OptAccDelete()); err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.Success("account %q successfully deleted", args[0])
		},
	}
	del.Flags().BoolVarP(&opts.force, "force", "f", false, "will bypass [y/N] prompt")

	return del
}
