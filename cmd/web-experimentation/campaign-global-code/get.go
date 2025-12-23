/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign_global_code

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

var createFile bool
var createSubFiles bool

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <campaign-id> | --id <campaign-id>]",
	Short: "Get campaign global code",
	Long:  `Get campaign global code`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.CampaignGlobalCodeRequester.HTTPGetCampaignGlobalCode(strconv.Itoa(CampaignID))
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		if createFile && len(body) > 0 {
			if !Override {
				jsFilePath, err := config.CampaignGlobalCodeFilePath(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, strconv.Itoa(CampaignID))
				if err != nil {
					log.Fatalf("error occurred: %v", err)
				}

				if _, err := os.Stat(jsFilePath); err == nil {
					fileHash, err := config.HashFile(jsFilePath)
					if err != nil {
						log.Fatalf("Error hashing file: %v", err)
					}

					strHash := config.HashString(body)
					if fileHash != strHash {
						log.Fatalf("error occurred: %s", utils.ERROR_REMOTE_CHANGED_FROM_LOCAL)
					}
				}
			}
			_, err := config.WriteCampaignGlobalCode(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, strconv.Itoa(CampaignID), body)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			return
		}

		if createSubFiles {
			if !Override {
				log.Fatalf("error occurred: %s", "You should run this command with the flag --override, this will automatically refresh your resources global code.")
			}

			if len(body) > 0 {
				_, err = config.WriteCampaignGlobalCode(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, strconv.Itoa(CampaignID), body)
				if err != nil {
					log.Fatalf("error occurred: %v", err)
				}
			}

			body, err := httprequest.ModificationRequester.HTTPListModification(CampaignID)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			for _, modification := range body {
				if modification.Type == "customScriptNew" && modification.Selector == "" {
					_, err := config.WriteVariationGlobalCodeJS(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, strconv.Itoa(CampaignID), strconv.Itoa(modification.VariationID), modification.Value)
					if err != nil {
						log.Fatalf("error occurred: %v", err)
					}
					continue
				}

				if modification.Type == "addCSS" && modification.Selector == "" {
					_, err := config.WriteVariationGlobalCodeCSS(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, strconv.Itoa(CampaignID), strconv.Itoa(modification.VariationID), modification.Value)
					if err != nil {
						log.Fatalf("error occurred: %v", err)
					}
					continue
				}

				pattern := `/\*\s*Selector: (.+)*\s*\*/`
				re := regexp.MustCompile(pattern)

				fileCode := config.AddHeaderSelectorComment(modification.Selector, []byte(modification.Value), re)
				config.WriteModificationCode(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, strconv.Itoa(CampaignID), strconv.Itoa(modification.VariationID), strconv.Itoa(modification.Id), fileCode)
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
	getCmd.Flags().IntVarP(&CampaignID, "id", "i", 0, "id of the campaign you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	getCmd.Flags().BoolVarP(&createFile, "create-file", "", false, "create a file that contains campaign global code")
	getCmd.Flags().BoolVarP(&createSubFiles, "create-subfiles", "", false, "create a file that contains campaign and variations global code")

	getCmd.Flags().BoolVarP(&Override, "override", "", false, "override local campaign global code file")

	CampaignGlobalCodeCmd.AddCommand(getCmd)
}
