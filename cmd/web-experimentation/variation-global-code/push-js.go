/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation_global_code

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

var jsCode string
var jsFilePath string

// pushJsCmd represents push command
var pushJSCmd = &cobra.Command{
	Use:   "push-js [-i <variation-id> | --id <variation-id>] [--campaign-id <campaign-id>]",
	Short: "Push variation global js code",
	Long:  `Push variation global js code`,
	Run: func(cmd *cobra.Command, args []string) {
		var modificationId int
		var modificationValue string
		var codeByte []byte

		if !utils.CheckSingleFlag(jsFilePath != "", jsCode != "") {
			log.Fatalf("error occurred: %s", "1 flag is required. (file, code)")
		}

		campaignID, err := strconv.Atoi(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		variationID, err := strconv.Atoi(VariationID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		modifList, err := httprequest.ModificationRequester.HTTPListModification(campaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		for _, modification := range modifList {
			if modification.VariationID == variationID && modification.Type == "customScriptNew" && modification.Selector == "" {
				modificationId = modification.Id
				modificationValue = modification.Value
			}
		}

		if jsFilePath != "" {
			fileContent, err := os.ReadFile(jsFilePath)
			if err != nil {
				log.Fatalf("error occurred: %s", err)
			}

			codeByte = fileContent
		}

		if jsCode != "" {
			codeByte = []byte(jsCode)
		}

		if modificationId == 0 {
			modificationToPush := web_experimentation.ModificationCodeCreateStruct{
				InputType:   "modification",
				Name:        "",
				Value:       string(codeByte),
				Selector:    "",
				Type:        "customScriptNew",
				Engine:      string(codeByte),
				VariationID: variationID,
			}

			body, err := httprequest.ModificationRequester.HTTPCreateModification(campaignID, modificationToPush)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(body))
			return
		}

		modificationToPush := web_experimentation.ModificationCodeEditStruct{
			InputType: "modification",
			Value:     string(codeByte),
			Engine:    string(codeByte),
		}

		if !Override {
			apiHash := config.HashString(modificationValue)
			strHash := config.HashString(string(codeByte))
			if apiHash != strHash {
				log.Fatalf("error occurred: %s", utils.ERROR_LOCAL_CHANGED_FROM_REMOTE)
			}
		}

		body, err := httprequest.ModificationRequester.HTTPEditModification(campaignID, modificationId, modificationToPush)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(body))
	},
}

func init() {
	pushJSCmd.Flags().StringVarP(&CampaignID, "campaign-id", "", "", "id of the campaign")
	if err := pushJSCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	pushJSCmd.Flags().StringVarP(&VariationID, "id", "i", "", "id of variation")
	if err := pushJSCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	pushJSCmd.Flags().StringVarP(&jsCode, "code", "c", "", "new code to push in the variation")
	pushJSCmd.Flags().StringVarP(&jsFilePath, "file", "", "", "file that contains new code to push in the variation")

	pushJSCmd.Flags().BoolVarP(&Override, "override", "", false, "override remote variation global code js")

	VariationGlobalCodeCmd.AddCommand(pushJSCmd)
}
