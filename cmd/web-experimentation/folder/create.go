/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package folder

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

func CreateFolder(dataRaw []byte) []byte {
	folderHeader, err := httprequest.FolderRequester.HTTPCreateFolder(dataRaw)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	parts := strings.Split(string(folderHeader), "/")
	folderID := parts[len(parts)-1]
	folderIDInt, err := strconv.Atoi(folderID)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	body, err := httprequest.FolderRequester.HTTPGetFolder(folderIDInt)
	if err != nil {
		log.Fatalf("error occurred: %s", err)
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("error occurred: %s", err)
	}

	return bodyByte
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data-raw=<data-raw>]",
	Short: "Create an folder",
	Long:  `Create an folder`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := CreateFolder([]byte(DataRaw))

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(resp))
	},
}

func init() {

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your folder, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	FolderCmd.AddCommand(createCmd)
}
