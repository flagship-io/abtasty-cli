/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/
package variation

import (
	"log"

	"github.com/flagship-io/flagship/utils"
	httprequest "github.com/flagship-io/flagship/utils/httpRequest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [--campaign-id=<campaign-id>] [--variation-group-id=<variation-group-id>] [-i <variation-id> | --id=<variation-id>]",
	Short: "Get a variation",
	Long:  `Get a variation in your variation group`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.HTTPGetVariation(CampaignID, VariationGroupID, VariationID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "Reference", "Allocation"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().StringVarP(&VariationID, "id", "i", "", "id of the variation you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	VariationCmd.AddCommand(getCmd)
}
