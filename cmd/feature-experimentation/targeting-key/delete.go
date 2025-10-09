/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package targeting_key

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func DeleteTargetingKey(id string) (string, error) {
	err := httprequest.TargetingKeyRequester.HTTPDeleteTargetingKey(id)
	if err != nil {
		return "", err
	}
	return "Targeting key deleted", nil
}

// deleteCmd represents delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [-i <targeting-key-id> | --id=<targeting-key-id>]",
	Short: "Delete a targeting key",
	Long:  `Delete a targeting key`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := DeleteTargetingKey(TargetingKeyID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), resp)

	},
}

func init() {
	deleteCmd.Flags().StringVarP(&TargetingKeyID, "id", "i", "", "id of the targeting key you want to delete")

	if err := deleteCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	TargetingKeyCmd.AddCommand(deleteCmd)
}
