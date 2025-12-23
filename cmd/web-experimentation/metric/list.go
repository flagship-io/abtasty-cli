/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package metric

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	web_experimentation "github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

var metricType string

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

		// Filter by type if specified
		if metricType != "" {
			body, err := filterMetricsByType(body, metricType)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}

			metricsJSON, err := json.Marshal(body)
			if err != nil {
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(metricsJSON))
			return
		}

		metricsJSON, err := json.Marshal(body)
		if err != nil {
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(metricsJSON))
	},
}

func filterMetricsByType(data web_experimentation.MetricsData, typeFilter string) (any, error) {
	if len(data.ActionTrackings) == 0 && len(data.Transactions) == 0 && len(data.WidgetTrackings) == 0 && len(data.CustomTrackings) == 0 && len(data.PageViews) == 0 && len(data.Indicators) == 0 {
		return data, nil
	}

	normalizedType := strings.ToLower(strings.ReplaceAll(typeFilter, " ", "_"))

	filtered := any(nil)

	switch normalizedType {
	case "action_tracking", "actiontracking", "at":
		if data.ActionTrackings == nil {
			return []web_experimentation.ActionTrackingMetric{}, nil
		}

		filtered = data.ActionTrackings
	case "transaction", "transactions", "t":
		if data.Transactions == nil {
			return []web_experimentation.TransactionMetric{}, nil
		}

		filtered = data.Transactions
	case "widget_tracking", "widgettracking", "wt":
		if data.WidgetTrackings == nil {
			return []web_experimentation.WidgetTrackingMetric{}, nil
		}

		filtered = data.WidgetTrackings
	case "custom_tracking", "customtracking", "ct":
		if data.CustomTrackings == nil {
			return []web_experimentation.CustomTrackingMetric{}, nil
		}

		filtered = data.CustomTrackings
	case "page_view", "pageview", "pv":
		if data.PageViews == nil {
			return []web_experimentation.PageViewMetric{}, nil
		}

		filtered = data.PageViews
	case "indicator", "indicators", "i":
		if data.Indicators == nil {
			return []web_experimentation.IndicatorMetric{}, nil
		}

		filtered = data.Indicators
	default:
		return nil, fmt.Errorf("invalid metric type: %s. Valid types are: action_tracking, transaction, widget_tracking, custom_tracking, page_view, indicator", typeFilter)
	}

	return filtered, nil
}

func init() {
	MetricCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&metricType, "type", "t", "", "Filter metrics by type (action_tracking, transaction, widget_tracking, custom_tracking, page_view, indicator)")
}
