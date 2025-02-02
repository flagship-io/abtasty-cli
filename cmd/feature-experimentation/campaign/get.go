/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign

import (
	"log"

	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [-i <campaign-id> | --id=<campaign-id>]",
	Short: "Get a campaign",
	Long:  `Get a campaign`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.CampaignFERequester.HTTPGetCampaign(CampaignID)
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
