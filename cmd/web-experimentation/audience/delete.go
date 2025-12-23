/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package audience

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

// deleteCmd represents delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [-i <audience-id> | --id=<audience-id>]",
	Short: "Delete an audience",
	Long:  `Delete an audience`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := httprequest.AudienceRequester.HTTPDeleteAudience(AudienceID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), resp)
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&AudienceID, "id", "i", "", "id of the audience you want to delete")

	if err := deleteCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	AudienceCmd.AddCommand(deleteCmd)
}
