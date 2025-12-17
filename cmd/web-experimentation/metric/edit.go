/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package metric

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [-i <metric-id> | --id=<metric-id>][-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Edit a metric",
	Long:  `Edit a metric`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := httprequest.MetricRequester.HTTPEditMetric(MetricID, []byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(resp))
	},
}

func init() {
	editCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to edit your metric, check the doc for details")

	if err := editCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	editCmd.Flags().StringVarP(&MetricID, "metric-id", "m", "", "ID of the metric to edit")

	if err := editCmd.MarkFlagRequired("metric-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	MetricCmd.AddCommand(editCmd)
}
