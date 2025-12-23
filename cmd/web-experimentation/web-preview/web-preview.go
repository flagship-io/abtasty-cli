/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package web_preview

import (
	"github.com/spf13/cobra"
)

var CampaignID int
var VariationID int

// WebPreviewCmd represents the web preview command
var WebPreviewCmd = &cobra.Command{
	Use:   "web-preview [open]",
	Short: "Open web preview",
	Long:  `Open web preview`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
