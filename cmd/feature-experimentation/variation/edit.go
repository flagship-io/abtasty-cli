/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

func EditVariation(campaignID, variationGroupID, id string, dataRaw []byte) ([]byte, error) {
	body, err := httprequest.VariationFERequester.HTTPEditVariation(campaignID, variationGroupID, id, dataRaw)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [--campaign-id=<campaign-id>] [--variation-group-id=<variation-group-id>] [-i <variation-id> | --id=<variation-id>] [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Edit a variation",
	Long:  `Edit a variation`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := EditVariation(CampaignID, VariationGroupID, VariationID, []byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", body)
	},
}

func init() {

	editCmd.Flags().StringVarP(&VariationID, "id", "i", "", "id of the variation you want to edit")
	editCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to edit your variation, check the doc for details")

	if err := editCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if err := editCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	VariationCmd.AddCommand(editCmd)
}
