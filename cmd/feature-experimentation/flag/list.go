/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package flag

import (
	"log"

	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all flags",
	Long:  `List all flags`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.FlagRequester.HTTPListFlag()
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "Type", "Description", "Source"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	FlagCmd.AddCommand(listCmd)
}
