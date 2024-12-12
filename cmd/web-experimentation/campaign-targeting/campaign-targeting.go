/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign_targeting

import (
	"github.com/spf13/cobra"
)

var CampaignID string

// CampaignTargetingCmd represents the campaign targeting command
var CampaignTargetingCmd = &cobra.Command{
	Use:   "campaign-targeting [get|push]",
	Short: "Manage campaign targeting",
	Long:  `Manage campaign targeting`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
