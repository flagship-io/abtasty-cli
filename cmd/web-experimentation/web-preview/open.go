/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package web_preview

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type WebPreview struct {
	CampaignID  int    `json:"campaign_id"`
	VariationID int    `json:"variation_id"`
	Url         string `json:"url"`
}

var isVariation bool

// openCmd represents open command
var openCmd = &cobra.Command{
	Use:   "open [--campaign-id <campaign-id>] [--variation-id <variation-id>]",
	Short: "Open web preview on variation",
	Long:  `Open web preview on variation`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.CampaignWERequester.HTTPGetCampaign(CampaignID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		for _, v := range body.Variations {
			if v.Id == VariationID {
				isVariation = true
			}
		}

		if isVariation {
			url := fmt.Sprintf(`%s/%s/?ab_project=preview&testId=%d&variationId=%d&t=%s`, body.Url, viper.GetString("identifier"), body.Id, VariationID, body.Report.Token)
			if err := common.OpenLink(url); err != nil {
				log.Fatalf("Error opening link: %s", err)
			}

			return
		}

		log.Fatalln("error occurred: no campaign or variation found !")
	},
}

func init() {
	openCmd.Flags().IntVarP(&CampaignID, "campaign-id", "", 0, "id of the campaign you want to display")
	if err := openCmd.MarkFlagRequired("campaign-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	openCmd.Flags().IntVarP(&VariationID, "variation-id", "", 0, "id of the variation you want to display")
	if err := openCmd.MarkFlagRequired("variation-id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	WebPreviewCmd.AddCommand(openCmd)
}
