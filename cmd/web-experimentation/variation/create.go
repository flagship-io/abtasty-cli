/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func CreateVariation(campaignId int, rawData []byte) []byte {
	variationHeader, err := httprequest.VariationWERequester.HTTPCreateVariationDataRaw(fmt.Sprint(CampaignID), rawData)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	parts := strings.Split(string(variationHeader), "/")
	variationID := parts[len(parts)-1]
	variationIDInt, err := strconv.Atoi(variationID)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	body, err := httprequest.VariationWERequester.HTTPGetVariation(campaignId, variationIDInt)
	if err != nil {
		log.Fatalf("error occurred: %s", err)
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("error occurred: %s", err)
	}

	return bodyByte
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [--campaign-id=<campaign-id>] [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create a variation",
	Long:  `Create a variation`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := CreateVariation(CampaignID, []byte(DataRaw))
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", resp)
	},
}

func init() {

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your variation, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	VariationCmd.AddCommand(createCmd)
}
