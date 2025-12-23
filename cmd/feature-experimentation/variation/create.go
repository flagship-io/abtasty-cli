/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func CreateVariation(campaignID, variationGroupID string, dataRaw []byte) ([]byte, error) {
	body, err := httprequest.VariationFERequester.HTTPCreateVariation(campaignID, variationGroupID, dataRaw)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [--campaign-id=<campaign-id>] [--variation-group-id=<variation-group-id>] [-d <data-raw> | --data-raw <data-raw>]",
	Short: "Create a variation",
	Long:  `Create a variation`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := CreateVariation(CampaignID, VariationGroupID, []byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", body)
	},
}

func init() {

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your variation, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	VariationCmd.AddCommand(createCmd)
}
