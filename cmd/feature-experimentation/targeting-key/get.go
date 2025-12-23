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

func GetTargetingKey(id string) (feature_experimentation.TargetingKey, error) {
	body, err := httprequest.TargetingKeyRequester.HTTPGetTargetingKey(id)
	if err != nil {
		return feature_experimentation.TargetingKey{}, err
	}
	return body, nil
}

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <targeting-key-id> | --id=<targeting-key-id>]",
	Short: "Get a targeting key",
	Long:  `Get a targeting key`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := GetTargetingKey(TargetingKeyID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "Type", "Description"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().StringVarP(&TargetingKeyID, "id", "i", "", "id of the targeting key you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	TargetingKeyCmd.AddCommand(getCmd)
}
