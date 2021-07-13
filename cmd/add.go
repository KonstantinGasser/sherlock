package cmd

import (
	"context"

	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/terminal"
	"github.com/spf13/cobra"
)

func cmdAdd(ctx context.Context, sherlock *internal.Sherlock) *cobra.Command {
	add := &cobra.Command{
		Use:   "add",
		Short: "add an group or account to sherlock",
		Long:  "add either a new group to sherlock or an account to a sherlock-group",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	add.AddCommand(cmdAddGroup(ctx, sherlock))
	add.AddCommand(cmdAddAccount(ctx, sherlock))

	return add
}

type addGroupOptions struct {
	gid      string
	insecure bool
}

func cmdAddGroup(ctx context.Context, sherlock *internal.Sherlock) *cobra.Command {
	var opts addGroupOptions
	addGroup := &cobra.Command{
		Use:   "group",
		Short: "add a group to sherlock",
		Long:  "add a new group for accounts to sherlock",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.gid == "" {
				terminal.Error("group name not set (use --name)")
				return
			}
			groupKey, err := terminal.ReadPassword("(%s) password: ", opts.gid)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			if err := sherlock.SetupGroup(opts.gid, groupKey); err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.Success("group %q added to sherlock", opts.gid)
		},
	}
	addGroup.Flags().StringVarP(&opts.gid, "gid", "G", "", "name for the sherlock group")
	addGroup.Flags().BoolVarP(&opts.insecure, "insecure", "i", false, "allow insecure group password")

	return addGroup
}

type addAccountOptions struct {
	name     string
	gid      string
	tag      string
	insecure bool
}

func cmdAddAccount(ctx context.Context, sherlock *internal.Sherlock) *cobra.Command {
	var opts addAccountOptions
	addGroup := &cobra.Command{
		Use:   "account",
		Short: "add an account to a sherlock group",
		Long:  "add a new account to a sherlock group",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.name == "" {
				terminal.Error("account name required (use --name)")
				return
			}
			groupKey, err := terminal.ReadPassword("(%s) password: ", opts.gid)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			password, err := terminal.ReadPassword("(%s) password: ", opts.name)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			account, err := internal.NewAccount(opts.name, password, opts.tag, opts.insecure)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			if err := sherlock.AddAccount(ctx, account, groupKey, opts.gid); err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.Success("account %q successfully added to %q", account.Name, opts.gid)
		},
	}

	addGroup.Flags().StringVarP(&opts.name, "name", "n", "", "name for the account")
	addGroup.Flags().StringVarP(&opts.gid, "gid", "G", "default", "group name where to add the account")
	addGroup.Flags().StringVarP(&opts.tag, "tag", "t", "", "optional tag for this account")
	addGroup.Flags().BoolVarP(&opts.insecure, "insecure", "i", false, "allow insecure group password")

	return addGroup
}
