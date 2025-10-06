/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation_group

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func EditVariationGroup(campaignID, id string, dataRaw []byte) ([]byte, error) {
	body, err := httprequest.VariationGroupRequester.HTTPEditVariationGroup(campaignID, id, dataRaw)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [--campaign-id=<campaign-id>] [-i <variation-group-id> | --id <variation-group-id>] [-d <data-raw> | --data-raw <data-raw>]",
	Short: "Edit a variation group",
	Long:  `Edit a variation group`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := EditVariationGroup(CampaignID, VariationGroupID, []byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", body)
	},
}

func init() {

	editCmd.Flags().StringVarP(&VariationGroupID, "id", "i", "", "the variation group id of your variation")
	editCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to edit your variation group, check the doc for details")

	if err := editCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if err := editCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	VariationGroupCmd.AddCommand(editCmd)
}
