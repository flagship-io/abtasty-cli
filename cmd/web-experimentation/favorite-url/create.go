/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package favorite_url

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func CreateFavoriteURL(dataRaw []byte) []byte {
	favoriteUrlHeader, err := httprequest.FavoriteUrlRequester.HTTPCreateFavoriteUrl(dataRaw)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	parts := strings.Split(string(favoriteUrlHeader), "/")
	favoriteUrlID := parts[len(parts)-1]
	body, err := httprequest.FavoriteUrlRequester.HTTPGetFavoriteUrl(favoriteUrlID)
	if err != nil {
		log.Fatalf("error occurred: %s", err)
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("error occurred: %s", err)
	}

	return bodyByte
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create a favorite URL",
	Long:  `Create a favorite URL`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := CreateFavoriteURL([]byte(DataRaw))

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(resp))
	},
}

func init() {

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your favorite URL, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	FavoriteUrlCmd.AddCommand(createCmd)
}
