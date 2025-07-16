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

func CreateModification(campaignId int, rawData []byte) []byte {
	modificationHeader, err := httprequest.ModificationRequester.HTTPCreateModificationDataRaw(CampaignID, DataRaw)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	modificationIDs, err := getModificationIDsFromURL(string(modificationHeader))
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	body, err := httprequest.ModificationRequester.HTTPGetModification(campaignId, modificationIDs[0])
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
	Short: "Create a modification",
	Long:  `Create a modification`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := CreateModification(CampaignID, []byte(DataRaw))
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
