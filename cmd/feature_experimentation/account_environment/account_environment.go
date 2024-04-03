/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package accountenvironment

import "github.com/spf13/cobra"

var (
	Username             string
	AccountEnvironmentID string
)

// ConfigurationCmd represents the configuration command
var AccountEnvironmentCmd = &cobra.Command{
	Use:   "account-environment [use|list|current]",
	Short: "Manage your CLI authentication",
	Long:  `Manage your CLI authentication`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
