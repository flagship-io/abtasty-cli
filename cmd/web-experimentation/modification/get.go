/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package modification

import (
	"log"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetModification(campaignID, modificationID int) (web_experimentation.ModificationResourceLoader, error) {
	modification, err := httprequest.ModificationRequester.HTTPGetModification(campaignID, modificationID)
	if err != nil {
		return web_experimentation.ModificationResourceLoader{}, err
	}

	return web_experimentation.ModificationResourceLoader{Id: modification.Id, Name: modification.Name, Type: getTypeFromModificationAPI(modification.Type), CampaignID: campaignID, Selector: modification.Selector, Code: modification.Value, VariationID: modification.VariationID}, nil
}

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <modification-id> | --id <modification-id>] [--campaign-id <campaign-id>]",
	Short: "Get a modification",
	Long:  `Get a modification`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := GetModification(CampaignID, ModificationID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		utils.FormatItem([]string{"Id", "Name", "Type", "VariationID", "Selector", "Code"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	getCmd.Flags().IntVarP(&ModificationID, "id", "i", 0, "id of the modification you want to display")
	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getCmd.Flags().IntVarP(&CampaignID, "campaign-id", "", 0, "campaign id of your modification")
	if err := getCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	ModificationCmd.AddCommand(getCmd)
}
