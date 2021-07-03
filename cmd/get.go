package cmd

import (
	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/KonstantinGasser/sherlock/terminal"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

type getOptions struct {
	verbose bool
}

func cmdGet(sherlock *internal.Sherlock) *cobra.Command {
	var opts getOptions
	get := &cobra.Command{
		Use:   "get",
		Short: "setup allows to initially set-up a main password for your vault",
		Long:  "to encrypt and decrypt your vault you will need to set-up a main password",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			groupKey, err := terminal.ReadPassword("group (%s) password: ", args[0])
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			account, err := sherlock.GetAccount(args[0], groupKey)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			if opts.verbose {
				terminal.Info(account.Password)
			}
			clipboard.WriteAll(account.Password)
		},
	}
	get.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "print plain password to cli")

	return get
}
