/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign

import (
	"log"

	"github.com/flagship-io/abtasty-cli/internal/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetCampaign(id string) (feature_experimentation.CampaignFE, error) {
	body, err := httprequest.CampaignFERequester.HTTPGetCampaign(id)
	if err != nil {
		return feature_experimentation.CampaignFE{}, err
	}

	return body, nil
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [-i <campaign-id> | --id=<campaign-id>]",
	Short: "Get a campaign",
	Long:  `Get a campaign`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := GetCampaign(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "ProjectId", "Name", "Description", "Type", "Status"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {

	getCmd.Flags().StringVarP(&CampaignID, "id", "i", "", "id of the campaign you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	CampaignCmd.AddCommand(getCmd)
}
