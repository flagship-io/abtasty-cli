/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	model "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	we "github.com/flagship-io/abtasty-cli/utils/http_request/web_experimentation"

	"github.com/spf13/cobra"
)

func CreateCampaign(rawData []byte) []byte {
	var campaignModel model.CampaignWEResourceLoader

	err := json.Unmarshal(rawData, &campaignModel)
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

	campaignHeader, err := httprequest.CampaignWERequester.HTTPCreateCampaign(campaignCommon)
	parts := strings.Split(string(campaignHeader), "/")
	campaignID := parts[len(parts)-1]

	campaignIDInt, err := strconv.Atoi(campaignID)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if campaignModel.Traffic != 0 || campaignModel.GlobalCodeCampaign != "" {

		campaignPatch, _ := json.Marshal(struct {
			Traffic            int    `json:"traffic,omitempty"`
			GlobalCodeCampaign string `json:"global_code,omitempty"`
		}{
			Traffic:            campaignModel.Traffic,
			GlobalCodeCampaign: campaignModel.GlobalCodeCampaign,
		})

		_, err = httprequest.CampaignWERequester.HTTPEditCampaign(campaignID, campaignPatch)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

	}

	if campaignModel.CampaignTargeting != nil {
		model := we.JsonModelToModel(*campaignModel.CampaignTargeting)

		parsedModel, err := json.Marshal(model)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		_, err = httprequest.CampaignTargetingRequester.HTTPPushCampaignTargeting(campaignID, parsedModel)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
	}

	if len(campaignModel.Variations) != 0 {
		for _, v := range campaignModel.Variations {
			varGlobalCode := v.GlobalCodeVariation
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

			variationIDInt, err := strconv.Atoi(variationID)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			if varGlobalCode.Js != "" {
				modificationToPush := web_experimentation.ModificationCodeCreateStruct{
					InputType:   "modification",
					Name:        "",
					Value:       string(varGlobalCode.Js),
					Selector:    "",
					Type:        "customScriptNew",
					Engine:      string(varGlobalCode.Js),
					VariationID: variationIDInt,
				}

				_, err = httprequest.ModificationRequester.HTTPCreateModification(campaignIDInt, modificationToPush)
				if err != nil {
					log.Fatalf("error occurred: %v", err)
				}
			}

			if varGlobalCode.Css != "" {
				modificationToPush := web_experimentation.ModificationCodeCreateStruct{
					InputType:   "modification",
					Name:        "",
					Value:       string(varGlobalCode.Css),
					Selector:    "",
					Type:        "addCSS",
					Engine:      string(varGlobalCode.Css),
					VariationID: variationIDInt,
				}

				_, err := httprequest.ModificationRequester.HTTPCreateModification(campaignIDInt, modificationToPush)
				if err != nil {
					log.Fatalf("error occurred: %v", err)
				}
			}
		}
	}

	//delete origin variation
	body, err := httprequest.CampaignWERequester.HTTPGetCampaign(campaignID)
	if err != nil {
		log.Fatalf("error occurred: %s", err)
	}

	if body.Variations[0].Id != 0 {
		err = httprequest.VariationWERequester.HTTPDeleteVariation(campaignIDInt, body.Variations[0].Id)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("error occurred: %s", err)
	}

	return bodyByte
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create a campaign",
	Long:  `Create a campaign`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := CreateCampaign([]byte(DataRaw))
		fmt.Fprintln(cmd.OutOrStdout(), string(resp))
	},
}

func init() {

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your campaign, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	CampaignCmd.AddCommand(createCmd)
}
