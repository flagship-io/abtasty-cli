/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package metric

import (
	"log"

	"github.com/flagship-io/abtasty-cli/internal/utils"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var getCmd = &cobra.Command{
	Use:   "get [-i <metric-id> | --id <metric-id>]",
	Short: "Get a metric",
	Long:  `Get a metric by its ID`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.MetricRequester.HTTPGetMetric(MetricID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "Type", "Hidden", "TestID", "AccountLevel"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().StringVarP(&MetricID, "id", "i", "", "id of the audience you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	MetricCmd.AddCommand(getCmd)
}
