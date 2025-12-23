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

	model "github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	we "github.com/flagship-io/abtasty-cli/internal/utils/http_request/web_experimentation"

	"github.com/spf13/cobra"
)

func CreateCampaign(rawData []byte) ([]byte, error) {
	var campaignModel model.CampaignWEResourceLoader

	err := json.Unmarshal(rawData, &campaignModel)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	campaignCommon, err := json.Marshal(struct {
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
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	campaignHeader, err := httprequest.CampaignWERequester.HTTPCreateCampaign(campaignCommon)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	parts := strings.Split(string(campaignHeader), "/")
	campaignID := parts[len(parts)-1]

	campaignIDInt, err := strconv.Atoi(campaignID)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	if campaignModel.Traffic != 0 || campaignModel.GlobalCodeCampaign != "" {

		campaignPatch, err := json.Marshal(struct {
			Traffic            int      `json:"traffic,omitempty"`
			GlobalCodeCampaign string   `json:"global_code,omitempty"`
			SourceCode         string   `json:"source_code,omitempty"`
			Labels             []string `json:"labels,omitempty"`
			Status             string   `json:"status,omitempty"`
		}{
			Traffic:            campaignModel.Traffic,
			GlobalCodeCampaign: campaignModel.GlobalCodeCampaign,
			SourceCode:         campaignModel.SourceCode,
			Labels:             campaignModel.Labels,
			Status:             campaignModel.Status,
		})
		if err != nil {
			return nil, fmt.Errorf("error occurred: %s", err)
		}

		_, err = httprequest.CampaignWERequester.HTTPEditCampaign(campaignIDInt, campaignPatch)
		if err != nil {
			return nil, fmt.Errorf("error occurred: %s", err)
		}

	}

	if campaignModel.CampaignTargeting != nil {
		model, err := we.JsonModelToModel(*campaignModel.CampaignTargeting)
		if err != nil {
			return nil, fmt.Errorf("error occurred: %v", err)
		}

		parsedModel, err := json.Marshal(model)
		if err != nil {
			return nil, fmt.Errorf("error occurred: %s", err)
		}

		_, err = httprequest.CampaignTargetingRequester.HTTPPushCampaignTargeting(campaignIDInt, parsedModel)
		if err != nil {
			return nil, fmt.Errorf("error occurred: %s", err)
		}
	}

	body, err := httprequest.CampaignWERequester.HTTPGetCampaign(campaignIDInt)
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
			fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
			return
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
