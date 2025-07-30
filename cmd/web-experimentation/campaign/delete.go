/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign

import (
	"fmt"
	"log"
	"strconv"

	"github.com/flagship-io/abtasty-cli/utils/config"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

// deleteCmd represents delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [-i <campaign-id> | --id=<campaign-id>]",
	Short: "Delete a campaign",
	Long:  `Delete a campaign`,
	Run: func(cmd *cobra.Command, args []string) {
		err := httprequest.CampaignWERequester.HTTPDeleteCampaign(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Campaign deleted")

		config.DeleteCampaignGlobalCodeDirectory(httprequest.CampaignWERequester.WorkingDir, httprequest.CampaignWERequester.AccountID, strconv.Itoa(CampaignID))
	},
}

func init() {
	deleteCmd.Flags().IntVarP(&CampaignID, "id", "i", 0, "id of the campaign you want to delete")

	if err := deleteCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	CampaignCmd.AddCommand(deleteCmd)
}
