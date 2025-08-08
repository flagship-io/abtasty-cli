/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign

import (
	"log"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetCampaign(campaignID int) (web_experimentation.CampaignWE, error) {
	body, err := httprequest.CampaignWERequester.HTTPGetCampaign(campaignID)
	if err != nil {
		return web_experimentation.CampaignWE{}, err
	}

	return body, nil
}

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <campaign-id> | --id <campaign-id>]",
	Short: "Get a campaign",
	Long:  `Get a campaign`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := GetCampaign(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "Description", "Type", "State", "Url"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().IntVarP(&CampaignID, "id", "i", 0, "id of the campaign you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	CampaignCmd.AddCommand(getCmd)
}
