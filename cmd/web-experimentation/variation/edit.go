/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package variation

import (
	"encoding/json"
	"fmt"
	"log"

	vgc "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/variation-global-code"
	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func createOrEditVariationGlobalCode(variationID, campaignID int, code string, modifType vgc.ModificationType) error {
	modif, err := vgc.GetModification(variationID, campaignID, modifType)
	if err != nil {
		return fmt.Errorf("error occurred: %v", err)
	}

	if code != "" {
		if modif == (web_experimentation.Modification{}) {
			modificationToPush := web_experimentation.ModificationCodeCreateStruct{
				InputType:   "modification",
				Name:        "",
				Value:       code,
				Selector:    "",
				Type:        string(modifType),
				Engine:      code,
				VariationID: variationID,
			}

			_, err = httprequest.ModificationRequester.HTTPCreateModification(campaignID, modificationToPush)
			if err != nil {
				return fmt.Errorf("error occurred: %v", err)
			}
		} else {
			modificationToPush := web_experimentation.ModificationCodeEditStruct{
				Value:  string(code),
				Engine: string(code),
			}

			_, err = httprequest.ModificationRequester.HTTPEditModification(campaignID, modif.Id, modificationToPush)
			if err != nil {
				return fmt.Errorf("error occurred: %v", err)
			}
		}
	}

	return nil
}

func EditVariation(variationID, campaignID int, rawData []byte) ([]byte, error) {

	var variationResourceLoader models.VariationResourceLoader

	err := json.Unmarshal(rawData, &variationResourceLoader)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %v", err)
	}

	if variationResourceLoader.Code.Js != "" {
		err = createOrEditVariationGlobalCode(variationID, campaignID, variationResourceLoader.Code.Js, vgc.ModificationJS)
		if err != nil {
			return nil, fmt.Errorf("error occurred: %v", err)
		}
	}

	if variationResourceLoader.Code.Css != "" {
		err = createOrEditVariationGlobalCode(variationID, campaignID, variationResourceLoader.Code.Css, vgc.ModificationCSS)
		if err != nil {
			return nil, fmt.Errorf("error occurred: %v", err)
		}
	}

	variationAPI := models.VariationWE{
		Name:        variationResourceLoader.Name,
		Description: variationResourceLoader.Description,
		Traffic:     variationResourceLoader.Traffic,
	}

	variationApiJSON, err := json.Marshal(variationAPI)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %v", err)
	}

	variationPatch, err := httprequest.VariationWERequester.HTTPEditVariation(campaignID, variationID, variationApiJSON)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %v", err)
	}

	return variationPatch, nil
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [-i <variation-id> | --id=<variation-id>] [--campaign-id=<campaign-id>] [ -d <data-raw> | --data-raw=<data-raw>]",
	Short: "Edit a variation",
	Long:  `Edit a variation`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := EditVariation(VariationID, CampaignID, []byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", body)
	},
}

func init() {

	editCmd.Flags().IntVarP(&VariationID, "id", "i", 0, "id of the variation you want to edit")
	editCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to edit your variation, check the doc for details")

	if err := editCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if err := editCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	VariationCmd.AddCommand(editCmd)
}
