/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package metric

import (
	"encoding/json"
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
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

		projectJSON, err := json.Marshal(body)
		if err != nil {
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(projectJSON))
	},
}

func init() {
	MetricCmd.AddCommand(listCmd)
}
