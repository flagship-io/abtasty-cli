/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation

import (
	"log"

	"github.com/flagship-io/abtasty-cli/internal/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListVariations(campaignID, variationGroupID string) ([]feature_experimentation.VariationFE, error) {
	return httprequest.VariationFERequester.HTTPListVariation(campaignID, variationGroupID)
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [--campaign-id=<campaign-id>] [--variation-group-id=<variation-group-id>]",
	Short: "List all variations",
	Long:  `List all variations`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := ListVariations(CampaignID, VariationGroupID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "Reference", "Allocation"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	VariationCmd.AddCommand(listCmd)
}
