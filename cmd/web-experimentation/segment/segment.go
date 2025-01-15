/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package segment

import (
	"github.com/spf13/cobra"
)

var (
	SegmentID string
	DataRaw   string
)

// segmentCmd represents the segment command
var SegmentCmd = &cobra.Command{
	Use:   "segment [list|get]",
	Short: "Manage your segments",
	Long:  `Manage your segments`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
