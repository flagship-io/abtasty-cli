/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation_global_code

import (
	"log"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// infoJSCmd represents info-js command
var infoJSCmd = &cobra.Command{
	Use:   "info-js [-i <variation-id> | --id <variation-id>] [--campaign-id <campaign-id>]",
	Short: "Get variation global js code info",
	Long:  `Get variation global js code info `,
	Run: func(cmd *cobra.Command, args []string) {
		var modif web_experimentation.Modification

		m, err := GetModification(VariationID, CampaignID, ModificationJS)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		modif = m

		utils.FormatItem([]string{"Id", "Name", "Type", "VariationID", "Selector", "Engine", "Value"}, modif, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	infoJSCmd.Flags().IntVarP(&CampaignID, "campaign-id", "", 0, "campaign id of the variation")

	if err := infoJSCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	infoJSCmd.Flags().IntVarP(&VariationID, "id", "i", 0, "id of the variation you want to display")

	if err := infoJSCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	VariationGlobalCodeCmd.AddCommand(infoJSCmd)
}
