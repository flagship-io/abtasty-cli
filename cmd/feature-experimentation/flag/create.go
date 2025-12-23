/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package flag

import (
	"encoding/json"
	"fmt"
	"log"

	models "github.com/flagship-io/abtasty-cli/internal/models/feature_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

func CreateFlag(dataRaw []byte) ([]byte, error) {
	body, err := httprequest.FlagRequester.HTTPCreateFlag(dataRaw)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [--name <flag-name> --type <flag-type> --default-value <default-value> --description <description> --predefined-values <predefined-values> | -d <data-raw> | --data-raw <data-raw>]",
	Short: "Create a flag",
	Long:  `Create a flag`,
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte
		var predefinedValues_ []string

		if DataRaw != "" {
			data = []byte(DataRaw)
		} else {
			if FlagPredefinedValues != "" {
				err := json.Unmarshal([]byte(FlagPredefinedValues), &predefinedValues_)
				if err != nil {
					log.Fatalf("error occurred: %s", err)
				}
			}

			data_, err := json.Marshal(models.Flag{
				Name:             FlagName,
				Type:             FlagType,
				Description:      FlagDescription,
				DefaultValue:     FlagDefaultValue,
				Source:           "cli",
				PredefinedValues: predefinedValues_,
			})

			if err != nil {
				log.Fatalf("error occurred: %s", err)
			}

			data = data_
		}

		body, err := httprequest.FlagRequester.HTTPCreateFlag(data)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", body)
	},
}

func init() {

	createCmd.Flags().StringVarP(&FlagName, "name", "", "", "name of the flag")
	createCmd.Flags().StringVarP(&FlagType, "type", "", "", "type of the flag")
	createCmd.Flags().StringVarP(&FlagDescription, "description", "", "Flag created from the CLI", "description of the flag")
	createCmd.Flags().StringVarP(&FlagDefaultValue, "default-value", "", "", "default value of the flag")
	createCmd.Flags().StringVarP(&FlagPredefinedValues, "predefined-values", "", "", "predefined valued for the flag")

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your flag, check the doc for details")

	FlagCmd.AddCommand(createCmd)
}
