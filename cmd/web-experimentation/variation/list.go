/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation

import (
	"log"

	variation_global_code "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/variation-global-code"
	"github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListVariations(campaignID int) (variations []web_experimentation.VariationResourceLoader, err error) {
	campaign, err := httprequest.CampaignWERequester.HTTPGetCampaign(campaignID)
	if err != nil {
		return []web_experimentation.VariationResourceLoader{}, err
	}

	for _, variation := range campaign.Variations {
		variationResourceLoader := web_experimentation.VariationResourceLoader{Id: variation.Id, Name: variation.Name, Type: variation.Type, Description: variation.Description, Traffic: variation.Traffic}
		vgc, err := variation_global_code.GetVariationGlobalCode(variation.Id, campaignID)
		if err != nil {
			return []web_experimentation.VariationResourceLoader{}, err
		}

		variationResourceLoader.Code = vgc
		variations = append(variations, variationResourceLoader)
	}

	return variations, nil
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [--campaign-id=<campaign-id>]",
	Short: "List variations of a campaign",
	Long:  `List variations of a campaign`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := ListVariations(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		utils.FormatItem([]string{"Id", "Name", "Description", "Type", "Traffic"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	VariationCmd.AddCommand(listCmd)
}
