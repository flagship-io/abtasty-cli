/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package metric

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

func CreateMetric(dataRaw []byte) ([]byte, error) {
	metricHeader, err := httprequest.MetricRequester.HTTPCreateMetric(dataRaw)
	if err != nil {
		return nil, err
	}

	// Parse URL to extract the ids query parameter
	parts := strings.Split(string(metricHeader), "?")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid response format: %s", string(metricHeader))
	}

	queryParams, err := url.ParseQuery(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse query parameters: %w", err)
	}

	metricID := queryParams.Get("ids")
	if metricID == "" {
		return nil, fmt.Errorf("ids parameter not found in response: %s", string(metricHeader))
	}

	body, err := httprequest.MetricRequester.HTTPGetMetric(metricID)
	if err != nil {
		return nil, err
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bodyByte, nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create a metric",
	Long:  `Create a metric`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := CreateMetric([]byte(DataRaw))
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
