/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
				err := sherlock.SetupGroup(opts.name)
				if err != nil {
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
			password, err := terminal.ReadPassword()
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			account, err := internal.NewAccount(opts.gid, opts.name, password, opts.desc, opts.insecure)
			if err != nil {
				terminal.Error(err.Error())
				return
			}
			if err := sherlock.AddAccount(context.Background(), account); err != nil {
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
