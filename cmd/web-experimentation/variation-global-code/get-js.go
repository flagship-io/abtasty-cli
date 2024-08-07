/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation_global_code

import (
	"fmt"
	"log"
	"strconv"

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
			if modification.VariationID == variationID && modification.Type == "customScriptNew" && modification.Selector == "" {
				jsCode = modification.Value
			}
		}

		if CreateFile && len(jsCode) > 0 {
			_, err := config.VariationGlobalCodeDirectoryJS(httprequest.ModificationRequester.WorkingDir, httprequest.CampaignGlobalCodeRequester.AccountID, CampaignID, VariationID, jsCode, Override)
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
	getJSCmd.Flags().StringVarP(&CampaignID, "campaign-id", "", "", "id of the campaign you want to display")

	if err := getJSCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getJSCmd.Flags().StringVarP(&VariationID, "id", "i", "", "id of the variation you want to display")

	if err := getJSCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getJSCmd.Flags().BoolVarP(&CreateFile, "create-file", "", false, "create a file that contains variation global code")
	getJSCmd.Flags().BoolVarP(&Override, "override", "", false, "override existing variation global code file")

	VariationGlobalCodeCmd.AddCommand(getJSCmd)
}
