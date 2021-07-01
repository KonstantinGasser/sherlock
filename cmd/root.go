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
	"github.com/KonstantinGasser/sherlock/internal"
	"github.com/spf13/cobra"
)

const (
	skippSetupFor = "setup"
)

func RootCmd(sherlock *internal.Sherlock) *cobra.Command {
	root := &cobra.Command{
		Use:           "sherlock",
		Short:         "sherlock a CLI password manager for the simple use",
		Version:       "not there yet",
		SilenceUsage:  true,
		SilenceErrors: true,
		// ensure that sherlock is properly set-up. This means that the default group
		// exists and that it holds an encrypted .vault file. "sherlock setup" is excluded from this check
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == skippSetupFor {
				return nil
			}
			return sherlock.IsSetUp()
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	root.AddCommand(cmdSetup(sherlock))
	root.AddCommand(cmdAddAccount(sherlock))

	return root
}
