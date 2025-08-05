/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation

import (
	"log"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListVariations(campaignID int) ([]web_experimentation.VariationWE, error) {
	body, err := httprequest.CampaignWERequester.HTTPGetCampaign(campaignID)
	if err != nil {
		return []web_experimentation.VariationWE{}, err
	}

	return body.Variations, nil
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
