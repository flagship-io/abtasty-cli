/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package account

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
	Short: "List all account id associated to your AB Tasty account",
	Long:  `List all account id associated to your AB Tasty account`,
	Run: func(cmd *cobra.Command, args []string) {

		body, err := httprequest.AccountWERequester.HTTPListAccount()
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "Identifier", "Role"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {

	AccountCmd.AddCommand(listCmd)
}
