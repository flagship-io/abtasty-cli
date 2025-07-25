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

	model "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	we "github.com/flagship-io/abtasty-cli/utils/http_request/web_experimentation"

	"github.com/spf13/cobra"
)

func CreateCampaign(rawData []byte) ([]byte, error) {
	var campaignModel model.CampaignWEResourceLoader

	err := json.Unmarshal(rawData, &campaignModel)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
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
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	if campaignModel.Traffic != 0 || campaignModel.GlobalCodeCampaign != "" {

		campaignPatch, _ := json.Marshal(struct {
			Traffic            int    `json:"traffic,omitempty"`
			GlobalCodeCampaign string `json:"global_code,omitempty"`
		}{
			Traffic:            campaignModel.Traffic,
			GlobalCodeCampaign: campaignModel.GlobalCodeCampaign,
		})

		_, err = httprequest.CampaignWERequester.HTTPEditCampaign(campaignIDInt, campaignPatch)
		if err != nil {
			return nil, fmt.Errorf("error occurred: %s", err)
		}

	}

	if campaignModel.CampaignTargeting != nil {
		model := we.JsonModelToModel(*campaignModel.CampaignTargeting)

		parsedModel, err := json.Marshal(model)
		if err != nil {
			return nil, fmt.Errorf("error occurred: %s", err)
		}

		_, err = httprequest.CampaignTargetingRequester.HTTPPushCampaignTargeting(campaignID, parsedModel)
		if err != nil {
			return nil, fmt.Errorf("error occurred: %s", err)
		}
	}

	body, err := httprequest.CampaignWERequester.HTTPGetCampaign(campaignID)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	return bodyByte, nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create a campaign",
	Long:  `Create a campaign`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := CreateCampaign([]byte(DataRaw))
		if err != nil {
			log.Fatalf("%v", err)
		}

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
