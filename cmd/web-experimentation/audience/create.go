/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package audience

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func CreateAudience(dataRaw []byte) []byte {
	audienceHeader, err := httprequest.AudienceRequester.HTTPCreateAudience(dataRaw)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	parts := strings.Split(string(audienceHeader), "/")
	audienceID := parts[len(parts)-1]
	body, err := httprequest.AudienceRequester.HTTPGetAudience(audienceID)
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
	Use:   "create [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create an audience",
	Long:  `Create an audience`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := CreateAudience([]byte(DataRaw))

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
