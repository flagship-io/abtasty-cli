/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package modification

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func EditModification(modificationID int, rawData []byte) ([]byte, error) {
	var modifResourceLoader web_experimentation.ModificationResourceLoader
	err := json.Unmarshal(rawData, &modifResourceLoader)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %v", err)
	}

	if modifResourceLoader.CampaignID == 0 {
		return nil, fmt.Errorf("error occurred: missing property %s", "campaign_id")
	}

	m := web_experimentation.ModificationCodeEditStruct{
		Name:     modifResourceLoader.Name,
		Value:    modifResourceLoader.Code,
		Selector: modifResourceLoader.Selector,
		Engine:   modifResourceLoader.Code,
	}

	dataRaw, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %v", err)
	}

	modificationPatched, err := httprequest.ModificationRequester.HTTPEditModificationDataRaw(modifResourceLoader.CampaignID, modificationID, dataRaw)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %v", err)
	}

	return modificationPatched, nil
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [-i <modification-id> | --id=<modification-id>] [--campaign-id=<campaign-id>] [ -d <data-raw> | --data-raw=<data-raw>]",
	Short: "Edit a modification",
	Long:  `Edit a modification`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := EditModification(ModificationID, []byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", body)
	},
}

func init() {
	editCmd.Flags().IntVarP(&ModificationID, "id", "i", 0, "id of the modification you want to edit")
	editCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to edit your modification, check the doc for details")

	if err := editCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	if err := editCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	ModificationCmd.AddCommand(editCmd)
}
