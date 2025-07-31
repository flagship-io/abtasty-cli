/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation_global_code

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

// getJsCmd represents get command
var getJSCmd = &cobra.Command{
	Use:   "get-js [-i <variation-id> | --id <variation-id>] [--campaign-id <campaign-id>]",
	Short: "Get variation global js code",
	Long:  `Get variation global js code`,
	Run: func(cmd *cobra.Command, args []string) {
		var jsCode string

		m, err := GetModification(VariationID, CampaignID, ModificationJS)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		jsCode = m.Value

		if CreateFile && len(jsCode) > 0 {
			if !Override {
				jsFilePath := config.VariationGlobalCodeJSFilePath(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, strconv.Itoa(CampaignID), strconv.Itoa(VariationID))
				if _, err := os.Stat(jsFilePath); err == nil {
					fileHash, err := config.HashFile(jsFilePath)
					if err != nil {
						log.Fatalf("Error hashing file: %v", err)
					}

					strHash := config.HashString(jsCode)
					if fileHash != strHash {
						log.Fatalf("error occurred: %s", utils.ERROR_REMOTE_CHANGED_FROM_LOCAL)
					}
				}
			}

			_, err := config.WriteVariationGlobalCodeJS(httprequest.ModificationRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, strconv.Itoa(CampaignID), strconv.Itoa(VariationID), jsCode)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			return
		}

		if len(jsCode) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), jsCode)
			return
		}
	},
}

func init() {
	getJSCmd.Flags().IntVarP(&CampaignID, "campaign-id", "", 0, "id of the campaign you want to display")

	if err := getJSCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getJSCmd.Flags().IntVarP(&VariationID, "id", "i", 0, "id of the variation you want to display")

	if err := getJSCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getJSCmd.Flags().BoolVarP(&CreateFile, "create-file", "", false, "create a file that contains variation global code")
	getJSCmd.Flags().BoolVarP(&Override, "override", "", false, "override local variation global code js file")

	VariationGlobalCodeCmd.AddCommand(getJSCmd)
}
