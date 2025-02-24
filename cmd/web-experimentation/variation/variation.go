/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/

package variation

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	CampaignID  int
	VariationID int
	DataRaw     string
)

// VariationCmd represents the variation command
var VariationCmd = &cobra.Command{
	Use:   "variation [get|delete]",
	Short: "Manage your variations",
	Long:  `Manage your variations`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	VariationCmd.PersistentFlags().IntVarP(&CampaignID, "campaign-id", "", 0, "campaign id of your variation")

	if err := VariationCmd.MarkPersistentFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
}
