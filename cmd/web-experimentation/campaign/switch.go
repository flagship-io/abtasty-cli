/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package campaign

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

// SwitchCmd represents the Switch command
var SwitchCmd = &cobra.Command{
	Use:   "switch [-i <campaign-id> | --id=<campaign-id>] [-s <status> | --status=<status>]",
	Short: "Switch a campaign state",
	Long:  `Switch a campaign state`,
	Run: func(cmd *cobra.Command, args []string) {
		if !(Status == "active" || Status == "paused") {
			fmt.Fprintln(cmd.OutOrStdout(), "Status can only have 2 values: active or paused")
		} else {
			err := httprequest.CampaignWERequester.HTTPSwitchStateCampaign(CampaignID, Status)
			if err != nil {
				log.Fatalf("error occurred: %v", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "campaign status set to %s\n", Status)
		}

	},
}

func init() {

	SwitchCmd.Flags().IntVarP(&CampaignID, "id", "i", 0, "id of the campaign you want to switch state")
	SwitchCmd.Flags().StringVarP(&Status, "status", "s", "", "status you want set to the campaign. Only 2 values are possible: active and paused")

	if err := SwitchCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if err := SwitchCmd.MarkFlagRequired("status"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	CampaignCmd.AddCommand(SwitchCmd)
}
