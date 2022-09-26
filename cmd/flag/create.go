/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/
package flag

import (
	"log"

	httprequest "github.com/flagship-io/flagship/utils/httpRequest"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data-raw <data-raw>]",
	Short: "Create a flag",
	Long:  `Create a flag in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.HTTPCreateFlag(DataRaw)
		if err != nil {
			log.Fatalf("error occured: %v", err)
		}
		log.Printf("flag created: %s", body)
	},
}

func init() {

	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your flag, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occured: %v", err)
	}

	FlagCmd.AddCommand(createCmd)
}
