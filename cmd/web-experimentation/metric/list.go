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
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all metrics",
	Long:  `List all metrics of an account`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.MetricRequester.HTTPListMetrics()
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Transactions", "ActionTrackings", "WidgetTrackings", "CustomTrackings", "PageViews", "Indicators"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	MetricCmd.AddCommand(listCmd)
}
