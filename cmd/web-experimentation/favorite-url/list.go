/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package favorite_url

import (
	"log"

	"github.com/flagship-io/abtasty-cli/internal/utils"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all favorite URLs",
	Long:  `List all favorite URLs of an account`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.FavoriteUrlRequester.HTTPListFavoriteUrl()
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "AllPositiveConditions", "AllNegativeConditions", "CssCode"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	FavoriteUrlCmd.AddCommand(listCmd)
}
