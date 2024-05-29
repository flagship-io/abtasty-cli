/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package account_environment

import "github.com/spf13/cobra"

var (
	Username             string
	AccountEnvironmentID string
)

// AccountEnvironmentCmd represents the account environment command
var AccountEnvironmentCmd = &cobra.Command{
	Use:   "account-environment [use|list|current]",
	Short: "Manage your account environment",
	Long:  `Manage your account environment`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
