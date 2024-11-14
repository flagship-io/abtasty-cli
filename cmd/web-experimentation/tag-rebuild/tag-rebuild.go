/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package tag_rebuild

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

var (
	CampaignID     int
	ModificationID int
	Status         string
	DataRaw        string
)

// modificationCmd represents the modification command
var RebuildTagCmd = &cobra.Command{
	Use:   "tag-rebuild",
	Short: "Rebuild your tag",
	Long:  `Rebuild the tag of your account`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := httprequest.AccountWERequester.HTTPRebuildTag(); err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Tag rebuild request sent...")
	},
}
