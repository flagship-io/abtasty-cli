/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign_global_code

import (
	"fmt"
	"log"
	"regexp"
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
	Short: "Get campaign global code",
	Long:  `Get campaign global code`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.CampaignGlobalCodeRequester.HTTPGetCampaignGlobalCode(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		if createFile && len(body) > 0 {
			_, err := config.CampaignGlobalCodeDirectory(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, CampaignID, body, override)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}
			return
		}

		if createSubFiles {
			campaignID, err := strconv.Atoi(CampaignID)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			if len(body) > 0 {
				_, err = config.CampaignGlobalCodeDirectory(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, CampaignID, body, override)
				if err != nil {
					log.Fatalf("error occurred: %v", err)
				}
			}

			body, err := httprequest.ModificationRequester.HTTPListModification(campaignID)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			for _, modification := range body {
				if modification.Type == "customScriptNew" && modification.Selector == "" {
					_, err := config.VariationGlobalCodeDirectoryJS(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, CampaignID, strconv.Itoa(modification.VariationID), modification.Value, override)
					if err != nil {
						log.Fatalf("error occurred: %v", err)
					}
					continue
				}

				if modification.Type == "addCSS" && modification.Selector == "" {
					_, err := config.VariationGlobalCodeDirectoryCSS(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, CampaignID, strconv.Itoa(modification.VariationID), modification.Value, override)
					if err != nil {
						log.Fatalf("error occurred: %v", err)
					}
					continue
				}
				pattern := `/\*\s*Selector: (.+)*\s*\*/`
				re := regexp.MustCompile(pattern)

				fileCode := config.AddHeaderSelectorComment(modification.Selector, []byte(modification.Value), re)
				config.ModificationCodeDirectory(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, CampaignID, strconv.Itoa(modification.VariationID), strconv.Itoa(modification.Id), modification.Selector, fileCode, override)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Sub files code generated successfully: "+httprequest.CampaignGlobalCodeRequester.WorkingDir+"/.abtasty")
			return
		}

		if len(body) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), body)
			return
		}
	},
}

func init() {
	getCmd.Flags().StringVarP(&CampaignID, "id", "i", "", "id of the campaign you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	getCmd.Flags().BoolVarP(&createFile, "create-file", "", false, "create a file that contains campaign global code")
	getCmd.Flags().BoolVarP(&createSubFiles, "create-subfiles", "", false, "create a file that contains campaign and variations global code")

	getCmd.Flags().BoolVarP(&override, "override", "", false, "override existing campaign global code file")

	CampaignGlobalCodeCmd.AddCommand(getCmd)
}
