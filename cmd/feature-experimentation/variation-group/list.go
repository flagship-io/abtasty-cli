/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation_group

import (
	"log"

	"github.com/flagship-io/abtasty-cli/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListVariationGroups(campaignID string) ([]feature_experimentation.VariationGroup, error) {
	return httprequest.VariationGroupRequester.HTTPListVariationGroup(campaignID)
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [--campaign-id=<campaign-id>]",
	Short: "List all variation groups",
	Long:  `List all variation groups`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := ListVariationGroups(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	VariationGroupCmd.AddCommand(listCmd)
}
