/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package modification

import (
	"fmt"
	"log"
	"strconv"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils/config"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

// deleteCmd represents delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [-i <modification-id> | --id=<modification-id>] [--campaign-id <campaign-id>]",
	Short: "Delete a modification",
	Long:  `Delete a modification`,
	Run: func(cmd *cobra.Command, args []string) {
		var modif *web_experimentation.Modification
		body, err := httprequest.ModificationRequester.HTTPGetModification(CampaignID, ModificationID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		resp, err := httprequest.ModificationRequester.HTTPDeleteModification(CampaignID, ModificationID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), resp)

		if body != (web_experimentation.Modification{}) {
			if body.Type == "customScriptNew" && body.Selector != "" {
				modif = &body
			}
		}
		config.DeleteModificationCodeDirectory(httprequest.CampaignWERequester.WorkingDir, httprequest.CampaignWERequester.AccountID, strconv.Itoa(CampaignID), strconv.Itoa(modif.VariationID), strconv.Itoa(ModificationID))
	},
}

func init() {
	deleteCmd.Flags().IntVarP(&ModificationID, "id", "i", 0, "id of the modification you want to delete")

	if err := deleteCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	deleteCmd.Flags().IntVarP(&CampaignID, "campaign-id", "", 0, "campaign id of your modification")
	if err := deleteCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	ModificationCmd.AddCommand(deleteCmd)
}
