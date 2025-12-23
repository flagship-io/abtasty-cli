/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package modification_code

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	"github.com/flagship-io/abtasty-cli/internal/utils/config"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <modification-id> | --id <modification-id>] [--campaign-id <campaign-id>]",
	Short: "Get modification code",
	Long:  `Get modification code`,
	Run: func(cmd *cobra.Command, args []string) {
		campaignID, err := strconv.Atoi(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		modificationID, err := strconv.Atoi(ModificationID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		body, err := httprequest.ModificationRequester.HTTPGetModification(campaignID, modificationID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		if body == (web_experimentation.Modification{}) {
			log.Fatalf("error occurred: no modification found")
		}

		if CreateFile {
			if !Override {
				jsFilePath, err := config.ModificationCodeFilePath(httprequest.CampaignGlobalCodeRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, CampaignID, strconv.Itoa(body.VariationID), ModificationID)
				if err != nil {
					log.Fatalf("error occurred: %v", err)
				}

				if _, err := os.Stat(jsFilePath); err == nil {
					fileHash, err := config.HashFile(jsFilePath)
					if err != nil {
						log.Fatalf("Error hashing file: %v", err)
					}

					strHash := config.HashString(body.Value)
					if fileHash != strHash {
						log.Fatalf("error occurred: %s", utils.ERROR_REMOTE_CHANGED_FROM_LOCAL)
					}
				}
			}

			pattern := `/\*\s*Selector: (.+)*\s*\*/`
			re := regexp.MustCompile(pattern)

			fileCode := config.AddHeaderSelectorComment(body.Selector, []byte(body.Value), re)
			_, err := config.WriteModificationCode(httprequest.ModificationRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, CampaignID, strconv.Itoa(body.VariationID), ModificationID, fileCode)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(body.Value))
	},
}

func init() {
	getCmd.Flags().StringVarP(&CampaignID, "campaign-id", "", "", "campaign id of the modification")

	if err := getCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getCmd.Flags().StringVarP(&ModificationID, "id", "i", "", "id of the modification you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getCmd.Flags().BoolVarP(&CreateFile, "create-file", "", false, "create a file that contains modification global code")
	getCmd.Flags().BoolVarP(&Override, "override", "", false, "override local modification code file")

	ModificationCodeCmd.AddCommand(getCmd)
}
