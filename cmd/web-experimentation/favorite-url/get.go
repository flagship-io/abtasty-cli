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

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <favorite-url-id> | --id <favorite-url-id>]",
	Short: "Get a favorite url",
	Long:  `Get a favorite url`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.FavoriteUrlRequester.HTTPGetFavoriteUrl(FavoriteUrlID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "AllPositiveConditions", "AllNegativeConditions", "CssCode"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().StringVarP(&FavoriteUrlID, "id", "i", "", "id of the favorite url you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	FavoriteUrlCmd.AddCommand(getCmd)
}
