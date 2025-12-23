/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package targeting_key

import (
	"log"

	"github.com/flagship-io/abtasty-cli/internal/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListTargetingKeys() ([]feature_experimentation.TargetingKey, error) {
	return httprequest.TargetingKeyRequester.HTTPListTargetingKey()
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all targeting keys",
	Long:  `List all targeting keys`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := ListTargetingKeys()
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "Type", "Description"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	TargetingKeyCmd.AddCommand(listCmd)
}
