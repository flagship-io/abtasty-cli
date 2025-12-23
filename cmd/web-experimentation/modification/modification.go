/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package modification

import (
	"github.com/spf13/cobra"
)

var (
	CampaignID     int
	VariationID    int
	ModificationID int
	Status         string
	DataRaw        string
)

// modificationCmd represents the modification command
var ModificationCmd = &cobra.Command{
	Use:   "modification [create|edit|get|list|delete]",
	Short: "Manage your modifications",
	Long:  `Manage your modifications`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func getTypeFromModificationAPI(t string) string {
	if t == "addCSS" {
		return "css"
	}
	return "js"
}
