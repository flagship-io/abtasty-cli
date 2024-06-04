/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation_group

import (
	"log"

	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [--campaign-id=<campaign-id>] [-i <variation-group-id> | --id <variation-group-id>]",
	Short: "Get a variation group",
	Long:  `Get a variation group`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.VariationGroupRequester.HTTPGetVariationGroup(CampaignID, VariationGroupID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().StringVarP(&VariationGroupID, "id", "i", "", "the variation group id of your variation")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	VariationGroupCmd.AddCommand(getCmd)
}
