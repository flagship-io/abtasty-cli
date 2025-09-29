/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign_targeting

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/flagship-io/abtasty-cli/utils/config"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

var createFile bool
var createSubFiles bool
var override bool

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <campaign-id> | --id <campaign-id>]",
	Short: "Get campaign targeting",
	Long:  `Get campaign targeting`,
	Run: func(cmd *cobra.Command, args []string) {
		targeting, err := httprequest.CampaignTargetingRequester.HTTPGetCampaignTargeting(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		bodyStringify, err := json.Marshal(targeting)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		if createFile && len(string(bodyStringify)) > 0 {
			_, err := config.CampaignTargetingDirectory(httprequest.CampaignTargetingRequester.WorkingDir, httprequest.CampaignTargetingRequester.AccountID, strconv.Itoa(CampaignID), string(bodyStringify), override)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}
			return
		}

		if len(string(bodyStringify)) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), string(bodyStringify))
			return
		}
	},
}

func init() {
	getCmd.Flags().IntVarP(&CampaignID, "id", "i", 0, "id of the campaign you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	getCmd.Flags().BoolVarP(&createFile, "create-file", "", false, "create a file that contains campaign targeting details")
	getCmd.Flags().BoolVarP(&override, "override", "", false, "override existing campaign targeting details file")

	CampaignTargetingCmd.AddCommand(getCmd)
}
