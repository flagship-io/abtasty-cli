/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package audience

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func CreateAudience(dataRaw []byte) ([]byte, error) {
	var audiencePayload models.AudiencePayload
	err := json.Unmarshal(dataRaw, &audiencePayload)
	if err != nil {
		return nil, err
	}

	audienceHeader, err := httprequest.AudienceRequester.HTTPCreateAudience(dataRaw)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(audienceHeader), "/")
	audienceID := parts[len(parts)-1]
	body, err := httprequest.AudienceRequester.HTTPGetAudience(audienceID)
	if err != nil {
		return nil, err
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bodyByte, nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create an audience",
	Long:  `Create an audience`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := CreateAudience([]byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %s", err)

		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(resp))
	},
}

func init() {

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your audience, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	AudienceCmd.AddCommand(createCmd)
}
