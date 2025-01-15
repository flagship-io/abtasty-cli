/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package trigger

import (
	"github.com/spf13/cobra"
)

var (
	TriggerID string
	DataRaw   string
)

// triggerCmd represents the trigger command
var TriggerCmd = &cobra.Command{
	Use:   "trigger [list|get]",
	Short: "Manage your triggers",
	Long:  `Manage your triggers`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
