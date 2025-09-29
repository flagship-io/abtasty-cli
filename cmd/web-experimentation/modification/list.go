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

func ListModifications(campaignID int) (modificationsRL []web_experimentation.ModificationResourceLoader, err error) {
	modifications, err := httprequest.ModificationRequester.HTTPListModification(campaignID)
	if err != nil {
		return []web_experimentation.ModificationResourceLoader{}, err
	}

	for _, modification := range modifications {
		modificationsRL = append(modificationsRL, web_experimentation.ModificationResourceLoader{Id: modification.Id, Name: modification.Name, Type: getTypeFromModificationAPI(modification.Type), CampaignID: campaignID, Selector: modification.Selector, Code: modification.Value, VariationID: modification.VariationID})
	}

	return modificationsRL, nil
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [--campaign-id <campaign-id>]",
	Short: "List all modifications",
	Long:  `List all modifications of a campaign`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := ListModifications(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		utils.FormatItem([]string{"Id", "Name", "Type", "VariationID", "CampaignID", "Selector", "Code"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	listCmd.Flags().IntVarP(&CampaignID, "campaign-id", "", 0, "campaign id of your modification")
	if err := listCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	ModificationCmd.AddCommand(listCmd)
}
