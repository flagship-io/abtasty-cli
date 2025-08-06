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

func ListFolder() ([]web_experimentation.Folder, error) {
	return httprequest.FolderRequester.HTTPListFolder()
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all folders",
	Long:  `List all folders of an account`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := ListFolder()
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	FolderCmd.AddCommand(listCmd)
}
