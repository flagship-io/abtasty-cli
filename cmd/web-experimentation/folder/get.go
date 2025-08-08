/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package folder

import (
	"log"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetFolder(folderID int) (web_experimentation.Folder, error) {
	body, err := httprequest.FolderRequester.HTTPGetFolder(folderID)
	if err != nil {
		return web_experimentation.Folder{}, err
	}

	return body, nil
}

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <folder-id> | --id <folder-id>]",
	Short: "Get an folder",
	Long:  `Get an folder`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := GetFolder(FolderID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		utils.FormatItem([]string{"Id", "Name"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	getCmd.Flags().IntVarP(&FolderID, "id", "i", 0, "id of the folder you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	FolderCmd.AddCommand(getCmd)
}
