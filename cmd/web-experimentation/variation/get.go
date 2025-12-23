/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation

import (
	"log"

	variation_global_code "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/variation-global-code"
	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetVariation(campaignID, variationID int) (variationResourceLoader web_experimentation.VariationResourceLoader, err error) {
	variation, err := httprequest.VariationWERequester.HTTPGetVariation(campaignID, variationID)
	if err != nil {
		return web_experimentation.VariationResourceLoader{}, err
	}

	variationResourceLoader = web_experimentation.VariationResourceLoader{Id: variation.Id, Name: variation.Name, Type: variation.Type, Description: variation.Description, Traffic: variation.Traffic}
	vgc, err := variation_global_code.GetVariationGlobalCode(variationID, campaignID)
	if err != nil {
		return web_experimentation.VariationResourceLoader{}, err
	}

	variationResourceLoader.Code = vgc
	return variationResourceLoader, nil
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [--test-id=<test-id>] [-i=<variation-id> | --id=<variation-id>]",
	Short: "Get a variation",
	Long:  `Get a variation`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := GetVariation(CampaignID, VariationID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		utils.FormatItem([]string{"Id", "Name", "Description", "Type", "Traffic"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().IntVarP(&VariationID, "id", "i", 0, "id of the variation you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	VariationCmd.AddCommand(getCmd)
}
