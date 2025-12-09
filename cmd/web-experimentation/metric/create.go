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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create a metric",
	Long:  `Create a metric`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := httprequest.MetricRequester.HTTPCreateMetric([]byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(resp))
	},
}

func init() {
	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your metric, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	MetricCmd.AddCommand(createCmd)
}
