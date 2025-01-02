/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign_global_code

import (
	"fmt"
	"log"
	"os"

	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

var code string
var filePath string

// pushCmd represents push command
var pushCmd = &cobra.Command{
	Use:   "push [-i <campaign-id> | --id <campaign-id>]",
	Short: "Push campaign global code",
	Long:  `Push campaign global code`,
	Run: func(cmd *cobra.Command, args []string) {
		var codeByte []byte

		if !utils.CheckSingleFlag(filePath != "", code != "") {
			log.Fatalf("error occurred: %s", "1 flag is required. (file, code)")
		}

		if filePath != "" {
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatalf("error occurred: %s", err)
			}

			codeByte = fileContent
		}

		if code != "" {
			codeByte = []byte(code)
		}

		apiCampaignGlobalCode, err := httprequest.CampaignGlobalCodeRequester.HTTPGetCampaignGlobalCode(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		if !Override {
			apiHash := config.HashString(apiCampaignGlobalCode)
			strHash := config.HashString(string(codeByte))
			if apiHash != strHash {
				log.Fatalf("error occurred: %s", utils.ERROR_LOCAL_CHANGED_FROM_REMOTE)
			}
		}

		body, err := httprequest.CampaignGlobalCodeRequester.HTTPPushCampaignGlobalCode(CampaignID, codeByte)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(body))
	},
}

func init() {
	pushCmd.Flags().StringVarP(&CampaignID, "id", "i", "", "id of the campaign")
	if err := pushCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	pushCmd.Flags().StringVarP(&code, "code", "c", "", "new code to push in the campaign")
	pushCmd.Flags().StringVarP(&filePath, "file", "", "", "file that contains new code to push in the campaign")

	pushCmd.Flags().BoolVarP(&Override, "override", "", false, "override remote campaign global code")

	CampaignGlobalCodeCmd.AddCommand(pushCmd)
}
