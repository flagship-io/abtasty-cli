/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign

import (
	"encoding/json"
	"fmt"
	"log"

	model "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	we "github.com/flagship-io/abtasty-cli/utils/http_request/web_experimentation"
	"github.com/spf13/cobra"
)

func EditCampaign(campaignID int, rawData []byte) ([]byte, error) {
	var campaignModel model.CampaignWEResourceLoader

	err := json.Unmarshal(rawData, &campaignModel)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
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

	campaignCommon, err := json.Marshal(struct {
		Name               string `json:"name,omitempty"`
		Url                string `json:"url,omitempty"`
		Description        string `json:"description,omitempty"`
		Type               string `json:"type,omitempty"`
		Traffic            int    `json:"traffic,omitempty"`
		GlobalCodeCampaign string `json:"global_code,omitempty"`
	}{
		Name:               campaignModel.Name,
		Url:                campaignModel.Url,
		Description:        campaignModel.Description,
		Type:               campaignModel.Type,
		Traffic:            campaignModel.Traffic,
		GlobalCodeCampaign: campaignModel.GlobalCodeCampaign,
	})
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	campaignEdited, err := httprequest.CampaignWERequester.HTTPEditCampaign(campaignID, campaignCommon)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	return campaignEdited, nil
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [-i <campaign-id> | --id=<campaign-id>] [ -d <data-raw> | --data-raw=<data-raw>]",
	Short: "Edit a campaign",
	Long:  `Edit a campaign`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := EditCampaign(CampaignID, []byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", body)
	},
}

func init() {

	editCmd.Flags().IntVarP(&CampaignID, "id", "i", 0, "id of the campaign you want to edit")
	editCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to edit your campaign, check the doc for details")

	if err := editCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if err := editCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	CampaignCmd.AddCommand(editCmd)
}
