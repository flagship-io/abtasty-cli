/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign_global_code

import (
	"github.com/spf13/cobra"
)

var CampaignID string
var Override bool

// CampaignGlobalCodeCmd represents the campaign global code command
var CampaignGlobalCodeCmd = &cobra.Command{
	Use:     "campaign-global-code [get|push]",
	Short:   "Manage campaign global code",
	Aliases: []string{"cgc"},
	Long:    `Manage campaign global code`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
