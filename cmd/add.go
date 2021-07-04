package cmd

import (
	"context"

	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/terminal"
	"github.com/spf13/cobra"
)

type addOptions struct {
	group    string
	gid      string
	name     string
	tag      string
	insecure bool
}

func cmdAddAccount(sherlock *internal.Sherlock) *cobra.Command {
	var opts addOptions

	add := &cobra.Command{
		Use:   "add",
		Short: "add an account to sherlock",
		Long:  "add and configure a new account you want to access in a secure manner",
		Run: func(cmd *cobra.Command, args []string) {
			// creation of a group
			if opts.group != "default" {
				groupKey, err := terminal.ReadPassword("(%s) password: ", opts.group)
				if err != nil {
					terminal.Error(err.Error())
					return
				}
				if err := sherlock.SetupGroup(opts.group, groupKey); err != nil {
					terminal.Error(err.Error())
					return
				}
				terminal.Success("Group %q added to sherlock", opts.group)
				return
			}

			if opts.name == "" {
				terminal.Error("account name required (--name)")
				return
			}
			groupKey, err := terminal.ReadPassword("(%s) password: ", opts.gid)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			password, err := terminal.ReadPassword("account (%s) password: ", opts.name)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			account, err := internal.NewAccount(opts.name, password, opts.tag, opts.insecure)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			if err := sherlock.AddAccount(context.Background(), account, groupKey, opts.gid); err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.Success("Account %q successfully added to %q", account.Name, opts.gid)
		},
	}

	add.Flags().StringVarP(&opts.gid, "gid", "g", "default", "map account to existing group")
	add.Flags().StringVarP(&opts.name, "name", "n", "", "name of the account/group")
	add.Flags().StringVarP(&opts.tag, "tag", "t", "", "a tag to give some more meaning")
	add.Flags().BoolVarP(&opts.insecure, "insecure", "i", false, "allow insecure password for account")
	add.Flags().StringVarP(&opts.group, "group", "G", "default", "add a group to organize accounts")

	return add
}
