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

// deleteCmd represents delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [-i <metric-id> | --id=<metric-id>]",
	Short: "Delete a metric",
	Long:  `Delete a metric`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := httprequest.MetricRequester.HTTPDeleteMetric(MetricID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), resp)
	},
}

func init() {
	deleteCmd.Flags().IntVarP(&MetricID, "id", "i", 0, "id of the metric you want to delete")

	if err := deleteCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	MetricCmd.AddCommand(deleteCmd)
}
