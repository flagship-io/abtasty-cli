/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package account

import "github.com/spf13/cobra"

var (
	Username  string
	AccountID string
)

// AccountCmd represents the account command
var AccountCmd = &cobra.Command{
	Use:   "account [use|list|current]",
	Short: "Manage your account",
	Long:  `Manage your account`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
