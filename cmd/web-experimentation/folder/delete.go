/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package folder

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

// deleteCmd represents delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [-i <folder-id> | --id=<folder-id>]",
	Short: "Delete a folder",
	Long:  `Delete a folder`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := httprequest.FolderRequester.HTTPDeleteFolder(FolderID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), resp)
	},
}

func init() {
	deleteCmd.Flags().IntVarP(&FolderID, "id", "i", 0, "id of the folder you want to delete")

	if err := deleteCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	FolderCmd.AddCommand(deleteCmd)
}
