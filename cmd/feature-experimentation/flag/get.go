/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package flag

import (
	"log"

	"github.com/flagship-io/abtasty-cli/internal/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetFlag(id string) (feature_experimentation.Flag, error) {
	body, err := httprequest.FlagRequester.HTTPGetFlag(id)
	if err != nil {
		return feature_experimentation.Flag{}, err
	}

	return body, nil
}

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <flag-id> | --id=<flag-id>]",
	Short: "Get a flag",
	Long:  `Get a flag`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := GetFlag(FlagID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "Type", "Description", "Source"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().StringVarP(&FlagID, "id", "i", "", "id of the flag you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	FlagCmd.AddCommand(getCmd)
}
