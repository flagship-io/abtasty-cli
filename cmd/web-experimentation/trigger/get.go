/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package trigger

import (
	"log"

	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <trigger-id> | --id <trigger-id>]",
	Short: "Get a trigger",
	Long:  `Get a trigger`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.TriggerRequester.HTTPGetTrigger(TriggerID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name", "Description", "Hidden", "Archive", "IsSegment"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().StringVarP(&TriggerID, "id", "i", "", "id of the trigger you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	TriggerCmd.AddCommand(getCmd)
}
