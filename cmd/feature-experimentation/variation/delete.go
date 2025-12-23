/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func DeleteVariation(campaignId, variationGroupId, id string) (string, error) {
	err := httprequest.VariationFERequester.HTTPDeleteVariation(campaignId, variationGroupId, id)
	if err != nil {
		return "", err
	}
	return "Variation deleted", nil
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [--campaign-id=<campaign-id>] [--variation-group-id=<variation-group-id>] [-i <variation-id> | --id=<variation-id>]",
	Short: "Delete a variation",
	Long:  `Delete a variation`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := DeleteVariation(CampaignID, VariationGroupID, VariationID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), resp)
	},
}

func init() {

	deleteCmd.Flags().StringVarP(&VariationID, "id", "i", "", "id of the variation you want to delete")

	if err := deleteCmd.MarkFlagRequired("id"); err != nil {
		fmt.Fprintf(deleteCmd.OutOrStderr(), "error occurred: %s", err)
	}
	VariationCmd.AddCommand(deleteCmd)
}
