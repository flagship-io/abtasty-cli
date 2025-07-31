/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation_global_code

import (
	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

var WorkingDir string
var CampaignID int
var VariationID int
var CreateFile bool
var Override bool

type ModificationType string

const (
	ModificationJS  ModificationType = "customScriptNew"
	ModificationCSS ModificationType = "addCSS"
)

func GetModification(variationID, campaignID int, mType ModificationType) (m web_experimentation.Modification, err error) {
	body, err := httprequest.ModificationRequester.HTTPListModification(campaignID)
	if err != nil {
		return web_experimentation.Modification{}, err
	}

	for _, modification := range body {
		if modification.VariationID == variationID && modification.Type == string(mType) && modification.Selector == "" {
			m = modification
		}
	}

	return m, nil
}

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
