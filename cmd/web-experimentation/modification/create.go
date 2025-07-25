/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package modification

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func getModificationIDsFromURL(rawURL string) ([]int, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := u.Query()["ids"]
	ids := make([]int, 0, len(query))
	for _, s := range query {
		id, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func CreateModification(variationID int, modifResourceLoader web_experimentation.ModificationResourceLoader) ([]byte, error) {

	m := web_experimentation.ModificationCodeCreateStruct{
		InputType:   "modification",
		Name:        modifResourceLoader.Name,
		Value:       modifResourceLoader.Code,
		VariationID: variationID,
		Type:        "customScriptNew",
		Selector:    modifResourceLoader.Selector,
		Engine:      modifResourceLoader.Code,
	}

	if strings.ToLower(modifResourceLoader.Type) == "css" {
		m.Type = "addCSS"
	}

	dataRaw, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %v", err)
	}

	modificationHeader, err := httprequest.ModificationRequester.HTTPCreateModificationDataRaw(modifResourceLoader.CampaignID, dataRaw)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %v", err)
	}

	modificationIDs, err := getModificationIDsFromURL(string(modificationHeader))
	if err != nil {
		return nil, fmt.Errorf("error occurred: %v", err)
	}

	body, err := httprequest.ModificationRequester.HTTPGetModification(modifResourceLoader.CampaignID, modificationIDs[0])
	if err != nil {
		return nil, fmt.Errorf("error occurred: %v", err)
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %v", err)
	}

	return bodyByte, nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [--campaign-id=<campaign-id>] [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create a modification",
	Long:  `Create a modification`,
	Run: func(cmd *cobra.Command, args []string) {
		var modificationResourceLoader web_experimentation.ModificationResourceLoader
		err := json.Unmarshal([]byte(DataRaw), &modificationResourceLoader)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		resp, err := CreateModification(CampaignID, modificationResourceLoader)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", resp)
	},
}

func init() {

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your modification, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	ModificationCmd.AddCommand(createCmd)
}
