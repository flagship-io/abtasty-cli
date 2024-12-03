/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package audience

import (
	"log"

	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <audience-id> | --id <audience-id>]",
	Short: "Get an audience",
	Long:  `Get an audience`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.AudienceRequester.HTTPGetAudience(AudienceID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "Description", "Hidden", "Archive", "IsSegment"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().StringVarP(&AudienceID, "id", "i", "", "id of the audience you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	AudienceCmd.AddCommand(getCmd)
}
