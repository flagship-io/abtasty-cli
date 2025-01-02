/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package account_global_code

import (
	"github.com/spf13/cobra"
)

var AccountID string
var Override bool

// AccountGlobalCodeCmd represents the account global code command
var AccountGlobalCodeCmd = &cobra.Command{
	Use:     "account-global-code [get|push]",
	Short:   "Manage account global code",
	Aliases: []string{"agc"},
	Long:    `Manage account global code`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
