/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation_group

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func DeleteVariationGroup(campaignId, id string) (string, error) {
	err := httprequest.VariationGroupRequester.HTTPDeleteVariationGroup(campaignId, id)
	if err != nil {
		return "", err
	}
	return "Variation group deleted", nil
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [--campaign-id=<campaign-id>] [-i <variation-group-id> | --id <variation-group-id>]",
	Short: "Delete a variation group",
	Long:  `Delete a variation group`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := DeleteVariationGroup(CampaignID, VariationGroupID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), resp)
	},
}

func init() {

	deleteCmd.Flags().StringVarP(&VariationGroupID, "id", "i", "", "the variation group id of your variation")

	if err := deleteCmd.MarkFlagRequired("id"); err != nil {
		fmt.Fprintf(deleteCmd.OutOrStderr(), "error occurred: %s", err)
	}
	VariationGroupCmd.AddCommand(deleteCmd)
}
