/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package folder

import (
	"encoding/json"
	"fmt"
	"log"

	model "github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

func EditFolder(folderID int, rawData []byte) ([]byte, error) {
	var folderModel model.Folder

	err := json.Unmarshal(rawData, &folderModel)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	folderPatch, err := json.Marshal(struct {
		Name string `json:"name,omitempty"`
	}{
		Name: folderModel.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	folderEdited, err := httprequest.FolderRequester.HTTPEditFolder(folderID, folderPatch)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err)
	}

	return folderEdited, nil
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [-i <folder-id> | --id=<folder-id>] [ -d <data-raw> | --data-raw=<data-raw>]",
	Short: "Edit a folder",
	Long:  `Edit a folder`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := EditFolder(FolderID, []byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", body)
	},
}

func init() {
	editCmd.Flags().IntVarP(&FolderID, "id", "i", 0, "id of the folder you want to edit")
	editCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to edit your folder, check the doc for details")

	if err := editCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if err := editCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	FolderCmd.AddCommand(editCmd)
}
