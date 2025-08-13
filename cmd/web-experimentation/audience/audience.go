/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package audience

import (
	"github.com/spf13/cobra"
)

var (
	AudienceID string
	DataRaw    string
)

// AudienceCmd represents the audience command
var AudienceCmd = &cobra.Command{
	Use:   "audience [list|get|create|delete]",
	Short: "Manage your audiences",
	Long:  `Manage your audiences`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
