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

// getCSSCmd represents get command
var getCSSCmd = &cobra.Command{
	Use:   "get-css [-i <variation-id> | --id <variation-id>] [--campaign-id <campaign-id>]",
	Short: "Get variation global css code",
	Long:  `Get variation global css code`,
	Run: func(cmd *cobra.Command, args []string) {
		var cssCode string

		campaignID, err := strconv.Atoi(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		variationID, err := strconv.Atoi(VariationID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		body, err := httprequest.ModificationRequester.HTTPListModification(campaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		for _, modification := range body {
			if modification.VariationID == variationID && modification.Type == "addCSS" && modification.Selector == "" {
				cssCode = modification.Value
			}
		}

		if CreateFile && len(cssCode) > 0 {
			if !Override {
				cssFilePath := config.VariationGlobalCodeCSSFilePath(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, CampaignID, VariationID)
				if _, err := os.Stat(cssFilePath); err == nil {
					fileHash, err := config.HashFile(cssFilePath)
					if err != nil {
						log.Fatalf("Error hashing file: %v", err)
					}

					strHash := config.HashString(cssCode)
					if fileHash != strHash {
						log.Fatalf("error occurred: %s", utils.ERROR_REMOTE_CHANGED_FROM_LOCAL)
					}
				}
			}

			_, err := config.WriteVariationGlobalCodeCSS(httprequest.ModificationRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, CampaignID, VariationID, cssCode)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			return
		}

		if len(cssCode) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), cssCode)
			return
		}
	},
}

func init() {
	getCSSCmd.Flags().StringVarP(&CampaignID, "campaign-id", "", "", "campaign id of the variation")

	if err := getCSSCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getCSSCmd.Flags().StringVarP(&VariationID, "id", "i", "", "variation id")

	if err := getCSSCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getCSSCmd.Flags().BoolVarP(&CreateFile, "create-file", "", false, "create a file that contains variation global code")

	getCSSCmd.Flags().BoolVarP(&Override, "override", "", false, "override local variation global code css file")

	VariationGlobalCodeCmd.AddCommand(getCSSCmd)
}
