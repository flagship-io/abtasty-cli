/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package flag

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func DeleteFlag(id string) error {
	return httprequest.FlagRequester.HTTPDeleteFlag(id)
}

// deleteCmd represents delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [-i <flag-id> | --id=<flag-id>]",
	Short: "Delete a flag",
	Long:  `Delete a flag`,
	Run: func(cmd *cobra.Command, args []string) {
		err := DeleteFlag(FlagID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Flag deleted")

	},
}

func init() {
	deleteCmd.Flags().StringVarP(&FlagID, "id", "i", "", "id of the flag you want to delete")

	if err := deleteCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	FlagCmd.AddCommand(deleteCmd)
}
