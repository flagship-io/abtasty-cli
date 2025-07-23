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

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func CreateVariation(campaignID int, rawData []byte) []byte {

	var variationResourceLoader models.VariationResourceLoader

	err := json.Unmarshal(rawData, &variationResourceLoader)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	variationAPI := models.VariationWE{
		Name:        variationResourceLoader.Name,
		Description: variationResourceLoader.Description,
		Type:        "onthefly",
	}

	variationApiJSON, err := json.Marshal(variationAPI)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	variationHeader, err := httprequest.VariationWERequester.HTTPCreateVariationDataRaw(campaignID, variationApiJSON)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	parts := strings.Split(string(variationHeader), "/")
	variationID := parts[len(parts)-1]

	variationPatch := struct {
		Traffic int `json:"traffic,omitempty"`
	}{
		Traffic: variationResourceLoader.Traffic,
	}

	variationPatchJSON, err := json.Marshal(variationPatch)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	variationIDInt, err := strconv.Atoi(variationID)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	_, err = httprequest.VariationWERequester.HTTPEditVariation(campaignID, variationIDInt, variationPatchJSON)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if variationResourceLoader.Code.Js != "" {
		modificationToPush := web_experimentation.ModificationCodeCreateStruct{
			InputType:   "modification",
			Name:        "",
			Value:       string(variationResourceLoader.Code.Js),
			Selector:    "",
			Type:        "customScriptNew",
			Engine:      string(variationResourceLoader.Code.Js),
			VariationID: variationIDInt,
		}

		_, err = httprequest.ModificationRequester.HTTPCreateModification(campaignID, modificationToPush)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
	}

	if variationResourceLoader.Code.Css != "" {
		modificationToPush := web_experimentation.ModificationCodeCreateStruct{
			InputType:   "modification",
			Name:        "",
			Value:       string(variationResourceLoader.Code.Css),
			Selector:    "",
			Type:        "addCSS",
			Engine:      string(variationResourceLoader.Code.Css),
			VariationID: variationIDInt,
		}

		_, err := httprequest.ModificationRequester.HTTPCreateModification(campaignID, modificationToPush)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
	}

	body, err := httprequest.VariationWERequester.HTTPGetVariation(campaignID, variationIDInt)
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
