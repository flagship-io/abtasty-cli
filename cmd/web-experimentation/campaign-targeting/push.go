/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign_targeting

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	we "github.com/flagship-io/abtasty-cli/utils/http_request/web_experimentation"
	"github.com/spf13/cobra"
)

var dataRaw string
var filePath string

// pushCmd represents push command
var pushCmd = &cobra.Command{
	Use:   "push [-i <campaign-id> | --id <campaign-id>]",
	Short: "Push campaign targeting",
	Long:  `Push campaign targeting`,
	Run: func(cmd *cobra.Command, args []string) {
		var codeByte []byte
		var jsonModel web_experimentation.TargetingCampaignModelJSON

		if !utils.CheckSingleFlag(filePath != "", dataRaw != "") {
			log.Fatalf("error occurred: %s", "1 flag is required. (file, code)")
		}

		if filePath != "" {
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatalf("error occurred: %s", err)
			}

			codeByte = fileContent
		}

		if dataRaw != "" {
			codeByte = []byte(dataRaw)
		}

		err := json.Unmarshal(codeByte, &jsonModel)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		model := we.JsonModelToModel(jsonModel)

		parsedModel, err := json.Marshal(model)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		body, err := httprequest.CampaignTargetingRequester.HTTPPushCampaignTargeting(CampaignID, parsedModel)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(body))
	},
}

func init() {
	pushCmd.Flags().IntVarP(&CampaignID, "id", "i", 0, "id of the campaign")
	if err := pushCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	pushCmd.Flags().StringVarP(&dataRaw, "data-raw", "d", "", "new targeting json to push in the campaign")
	pushCmd.Flags().StringVarP(&filePath, "file", "", "", "file that contains new targeting json to push in the campaign")

	CampaignTargetingCmd.AddCommand(pushCmd)
}
