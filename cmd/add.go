package cmd

import (
	"context"
	"strconv"

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
	insecure bool
}

func cmdAddGroup(ctx context.Context, sherlock *internal.Sherlock) *cobra.Command {
	var opts addGroupOptions
	addGroup := &cobra.Command{
		Use:   "group",
		Short: "add a group to sherlock",
		Long:  "add a new group for accounts to sherlock",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				terminal.Error("group name not set (sherlock add group [group-name])")
				return
			}
			groupKey, err := terminal.ReadPassword("(%s) password: ", args[0])
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			if err := sherlock.SetupGroup(args[0], groupKey, opts.insecure); err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.Success("group %q added to sherlock", args[0])
		},
	}
	addGroup.Flags().BoolVarP(&opts.insecure, "insecure", "i", false, "allow insecure group password")

	return addGroup
}

type addAccountOptions struct {
	tag      string
	insecure bool
	gen      string
}

func cmdAddAccount(ctx context.Context, sherlock *internal.Sherlock) *cobra.Command {
	var opts addAccountOptions
	addGroup := &cobra.Command{
		Use:   "account",
		Short: "add an account to a sherlock group",
		Long:  "add a new account to a sherlock group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				terminal.Error("account name not set (sherlock add account [account-name])")
				return
			}
			// check if the group exists
			gid, _, err := internal.SplitQuery(args[0])
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			if err := sherlock.GroupExists(gid); err == nil {
				terminal.Error("group does not exist")
				return
			}

			groupKey, err := terminal.ReadPassword("(%s) password: ", args[0])
			if err != nil {
				terminal.Error(err.Error())
				return
			}

			// validate the password/key
			err = sherlock.CheckGroupKey(ctx, args[0], groupKey)
			if err != nil {
				terminal.Error(err.Error())
				return
			}

			// figure out password: either auto gen password or read from stdin
			var password string
			if opts.gen != "" { // generate password
				passwdLen, err := strconv.Atoi(opts.gen)
				if err != nil || passwdLen < 10 {
					terminal.Error("invalid length number for auto generated password (must be number grater then 10")
					return
				}
				password, err = internal.AutoGeneratePassword(passwdLen)
				if err != nil {
					terminal.Error(err.Error())
					return
				}
				terminal.Info("generated password : %s", password)
			} else {
				password, err = terminal.ReadPassword("(%s) password: ", args[0])
				if err != nil {
					terminal.Error(err.Error())
					return
				}
			}
			// create/store new Account
			account, err := internal.NewAccount(args[0], password, opts.tag, opts.insecure)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			if err := sherlock.UpdateState(ctx, args[0], groupKey, internal.OptAddAccount(account)); err != nil {
				terminal.Error(err.Error())
				return
			}
			terminal.Success("account %q successfully added to %q", account.Name, args[0])
		},
	}

	addGroup.Flags().StringVarP(&opts.tag, "tag", "t", "", "optional tag for this account")
	addGroup.Flags().BoolVarP(&opts.insecure, "insecure", "i", false, "allow insecure group password")

	// I set this to string to make input validation checking easier if the input data is not a valid number
	addGroup.Flags().StringVarP(&opts.gen, "gen", "e", "", "length for auto-generate secure password. Create your own password when not set")

	return addGroup
}
