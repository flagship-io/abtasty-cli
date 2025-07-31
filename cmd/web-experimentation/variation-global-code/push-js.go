/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation_global_code

import (
	"fmt"
	"log"
	"os"

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
		var codeByte []byte

		if !utils.CheckSingleFlag(jsFilePath != "", jsCode != "") {
			log.Fatalf("error occurred: %s", "1 flag is required. (file, code)")
		}

		m, err := GetModification(VariationID, CampaignID, ModificationJS)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
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

		if m.Id == 0 {
			modificationToPush := web_experimentation.ModificationCodeCreateStruct{
				InputType:   "modification",
				Name:        "",
				Value:       string(codeByte),
				Selector:    "",
				Type:        "customScriptNew",
				Engine:      string(codeByte),
				VariationID: VariationID,
			}

			body, err := httprequest.ModificationRequester.HTTPCreateModification(CampaignID, modificationToPush)
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
			apiHash := config.HashString(m.Value)
			strHash := config.HashString(string(codeByte))
			if apiHash != strHash {
				log.Fatalf("error occurred: %s", utils.ERROR_LOCAL_CHANGED_FROM_REMOTE)
			}
		}

		body, err := httprequest.ModificationRequester.HTTPEditModification(CampaignID, m.Id, modificationToPush)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(body))
	},
}

func init() {
	pushJSCmd.Flags().IntVarP(&CampaignID, "campaign-id", "", 0, "id of the campaign")
	if err := pushJSCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	pushJSCmd.Flags().IntVarP(&VariationID, "id", "i", 0, "id of variation")
	if err := pushJSCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	pushJSCmd.Flags().StringVarP(&jsCode, "code", "c", "", "new code to push in the variation")
	pushJSCmd.Flags().StringVarP(&jsFilePath, "file", "", "", "file that contains new code to push in the variation")

	pushJSCmd.Flags().BoolVarP(&Override, "override", "", false, "override remote variation global code js")

	VariationGlobalCodeCmd.AddCommand(pushJSCmd)
}
