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

var cssCode string
var cssFilePath string

// pushCSSCmd represents push command
var pushCSSCmd = &cobra.Command{
	Use:   "push-css [-i <variation-id> | --id <variation-id>] [--campaign-id <campaign-id>]",
	Short: "Push variation global css code",
	Long:  `Push variation global css code`,
	Run: func(cmd *cobra.Command, args []string) {
		var codeByte []byte

		if !utils.CheckSingleFlag(cssFilePath != "", cssCode != "") {
			log.Fatalf("error occurred: %s", "1 flag is required. (file, code)")
		}

		m, err := GetVariationGlobalCodePerType(VariationID, CampaignID, ModificationCSS)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		if cssFilePath != "" {
			fileContent, err := os.ReadFile(cssFilePath)
			if err != nil {
				log.Fatalf("error occurred: %s", err)
			}

			codeByte = fileContent
		}

		if cssCode != "" {
			codeByte = []byte(cssCode)
		}

		if m.Id == 0 {
			modificationToPush := web_experimentation.ModificationCodeCreateStruct{
				InputType:   "modification",
				Name:        "",
				Value:       string(codeByte),
				Selector:    "",
				Type:        "addCSS",
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
	pushCSSCmd.Flags().IntVarP(&CampaignID, "campaign-id", "", 0, "id of the campaign")
	if err := pushCSSCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	pushCSSCmd.Flags().IntVarP(&VariationID, "id", "i", 0, "id of variation")
	if err := pushCSSCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	pushCSSCmd.Flags().StringVarP(&cssCode, "code", "c", "", "new code to push in the variation")
	pushCSSCmd.Flags().StringVarP(&cssFilePath, "file", "", "", "file that contains new code to push in the variation")

	pushCSSCmd.Flags().BoolVarP(&Override, "override", "", false, "override remote variation global code css")

	VariationGlobalCodeCmd.AddCommand(pushCSSCmd)
}
