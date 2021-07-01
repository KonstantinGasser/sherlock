package cmd

import (
	"context"

	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/internal/terminal"
	"github.com/spf13/cobra"
)

type addOptions struct {
	isGroup  bool
	gid      string
	name     string
	desc     string
	insecure bool
}

func cmdAddAccount(sherlock *internal.Sherlock) *cobra.Command {
	var opts addOptions

	add := &cobra.Command{
		Use:   "add",
		Short: "add account to sherlock",
		Long:  "add and configure a new account you want to access in a secure manner",
		Run: func(cmd *cobra.Command, args []string) {

			// creation of a group
			if opts.isGroup {
				if opts.name == "" {
					terminal.Error("group name required (--name)")
					return
				}
				partionKey, err := terminal.ReadPassword("partition password: ")
				if err != nil {
					terminal.Error(err.Error())
					return
				}
				if err := sherlock.SetupGroup(opts.name, partionKey); err != nil {
					terminal.Error(err.Error())
					return
				}
				terminal.Success("Group %q added to sherlock", opts.name)
				return
			}

			if opts.name == "" {
				terminal.Error("account name required (--name)")
				return
			}
			partionKey, err := terminal.ReadPassword("partition password: ")
			if err != nil {
				terminal.Error(err.Error())
				return
			}

			password, err := terminal.ReadPassword("account password: ")
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			account, err := internal.NewAccount(opts.name, password, opts.desc, opts.insecure)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			if err := sherlock.AddAccount(context.Background(), account, partionKey, opts.gid); err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.Success("Account successfully added")
		},
	}

	add.Flags().StringVarP(&opts.gid, "gid", "g", "default", "map account to existing group")
	add.Flags().StringVarP(&opts.name, "name", "n", "", "name of the account/group")
	add.Flags().StringVarP(&opts.desc, "desc", "d", "", "description for the account")
	add.Flags().BoolVarP(&opts.insecure, "insecure", "i", false, "allow insecure password for account")
	add.Flags().BoolVarP(&opts.isGroup, "group", "G", false, "add a group to organize accounts")

	return add
}
