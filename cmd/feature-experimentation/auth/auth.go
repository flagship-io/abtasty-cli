/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package auth

import "github.com/spf13/cobra"

var (
	Username     string
	ClientID     string
	ClientSecret string
	AccountId    string
	AccountEnvID string
)

// AuthCmd represents the auth command
var AuthCmd = &cobra.Command{
	Use:     "authentication [login|get|list|delete]",
	Aliases: []string{"auth"},
	Short:   "Manage authentication for feature experimentation",
	Long:    `Manage authentication for feature experimentation`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
