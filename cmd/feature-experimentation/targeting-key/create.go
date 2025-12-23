/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package targeting_key

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

func CreateTargetingKey(dataRaw []byte) ([]byte, error) {
	body, err := httprequest.TargetingKeyRequester.HTTPCreateTargetingKey(dataRaw)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data-raw <data-raw>]",
	Short: "Create a targeting key",
	Long:  `Create a targeting key`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := CreateTargetingKey([]byte(DataRaw))
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", body)
	},
}

func init() {

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your targeting key, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	TargetingKeyCmd.AddCommand(createCmd)
}
