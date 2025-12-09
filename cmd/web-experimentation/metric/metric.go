/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package metric

import (
	"github.com/spf13/cobra"
)

var (
	MetricID string
	DataRaw  string
)

// MetricCmd represents the metric command
var MetricCmd = &cobra.Command{
	Use:   "metric [list|create|delete]",
	Short: "Manage your metrics",
	Long:  `Manage your metrics`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
