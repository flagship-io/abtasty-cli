/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation_global_code

import (
	"github.com/spf13/cobra"
)

var WorkingDir string
var CampaignID string
var VariationID string
var CreateFile bool
var Override bool

// VariationGlobalCodeCmd represents the variation global code command
var VariationGlobalCodeCmd = &cobra.Command{
	Use:     "variation-global-code [get-js|get-css|push-js|push-css|info-js|info-css]",
	Short:   "Manage variation global code",
	Aliases: []string{"vgc"},
	Long:    `Manage variation global code`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
