/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package modification_code

import (
	"fmt"
	"log"
	"strconv"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils/config"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var override bool

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <modification-id> | --id <modification-id>] [--campaign-id <campaign-id>]",
	Short: "Get modification code",
	Long:  `Get modification code`,
	Run: func(cmd *cobra.Command, args []string) {
		var modif *web_experimentation.Modification

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

		for _, modification := range body {
			if modification.Type == "customScriptNew" && modification.Selector != "" {
				modif = &modification
			}
		}

		if modif == nil {
			log.Fatalf("error occurred: no modification found")
		}

		if CreateFile {
			fileCode := config.AddHeaderSelectorComment(modif.Selector, modif.Value)
			modificationCodeDir, err := config.ModificationCodeDirectory(viper.GetString("working_dir"), httprequest.CampaignGlobalCodeRequester.AccountID, CampaignID, strconv.Itoa(modif.VariationID), ModificationID, modif.Selector, fileCode, override)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "modification code file generated successfully: ", modificationCodeDir)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(modif.Value))
	},
}

func init() {
	getCmd.Flags().StringVarP(&CampaignID, "campaign-id", "", "", "id of the campaign")

	if err := getCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getCmd.Flags().StringVarP(&ModificationID, "id", "i", "", "id of the  modification code you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getCmd.Flags().BoolVarP(&CreateFile, "create-file", "", false, "create a file that contains modification global code")
	getCmd.Flags().BoolVarP(&override, "override", "", false, "override existing modification code file")

	ModificationCodeCmd.AddCommand(getCmd)
}
