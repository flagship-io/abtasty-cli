/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package modification_code

import (
	"github.com/spf13/cobra"
)

var CampaignID string
var ModificationID string
var CreateFile bool
var Override bool

// ModificationCodeCmd represents the variation global code command
var ModificationCodeCmd = &cobra.Command{
	Use:     "modification-code [get|push]",
	Short:   "Manage modification code",
	Aliases: []string{"mc"},
	Long:    `Manage modification code`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
