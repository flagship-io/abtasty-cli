/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	model "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	we "github.com/flagship-io/abtasty-cli/utils/http_request/web_experimentation"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create a campaign",
	Long:  `Create a campaign`,
	Run: func(cmd *cobra.Command, args []string) {
		var campaignModel model.CampaignWEResourceLoader

		err := json.Unmarshal([]byte(DataRaw), &campaignModel)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		campaignCommon, _ := json.Marshal(struct {
			Name        string `json:"name"`
			Url         string `json:"url"`
			Description string `json:"description"`
			Type        string `json:"type"`
		}{
			Name:        campaignModel.Name,
			Url:         campaignModel.Url,
			Description: campaignModel.Description,
			Type:        campaignModel.Type,
		})

		campaignPatch, _ := json.Marshal(struct {
			Traffic            int    `json:"traffic,omitempty"`
			GlobalCodeCampaign string `json:"global_code,omitempty"`
		}{
			Traffic:            campaignModel.Traffic,
			GlobalCodeCampaign: campaignModel.GlobalCodeCampaign,
		})

		campaignHeader, err := httprequest.CampaignWERequester.HTTPCreateCampaign(campaignCommon)
		parts := strings.Split(string(campaignHeader), "/")
		campaignID := parts[len(parts)-1]
		fmt.Println(campaignID)

		model := we.JsonModelToModel(*campaignModel.CampaignTargeting)

		parsedModel, err := json.Marshal(model)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		_, err = httprequest.CampaignTargetingRequester.HTTPPushCampaignTargeting(campaignID, parsedModel)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		campaignEdit, err := httprequest.CampaignWERequester.HTTPEditCampaign(campaignID, campaignPatch)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(campaignEdit))

		for _, v := range campaignModel.Variations {
			variation := struct {
				Name        string `json:"name,omitempty"`
				Description string `json:"description,omitempty"`
				Type        string `json:"type,omitempty"`
			}{
				Name:        v.Name,
				Description: v.Description,
				Type:        v.Type,
			}

			variationJSON, err := json.Marshal(variation)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			variationHeader, err := httprequest.VariationWERequester.HTTPCreateVariationDataRaw(campaignID, variationJSON)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			parts := strings.Split(string(variationHeader), "/")
			variationID := parts[len(parts)-1]
			fmt.Println(variationID)

			variationPatch := struct {
				Traffic int `json:"traffic,omitempty"`
			}{
				Traffic: v.Traffic,
			}

			variationPatchJSON, err := json.Marshal(variationPatch)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			_, err = httprequest.VariationWERequester.HTTPEditVariation(campaignID, variationID, variationPatchJSON)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}
		}

		//err = httprequest.CampaignWERequester.HTTPDeleteCampaign(campaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
	},
}

func init() {

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your campaign, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	CampaignCmd.AddCommand(createCmd)
}
