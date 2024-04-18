/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package account_global_code

import (
	"fmt"
	"log"

	"github.com/flagship-io/flagship/utils/config"
	httprequest "github.com/flagship-io/flagship/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accountID string
var createFile bool

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <account-id> | --id <account-id>]",
	Short: "Get account global code",
	Long:  `Get account global code`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.AccountGlobalCodeRequester.HTTPGetAccountGlobalCode(accountID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		if createFile {
			accountCodeDir := config.AccountGlobalCodeDirectory(viper.GetString("working_dir"), accountID, body)
			fmt.Fprintln(cmd.OutOrStdout(), "Account code file generated successfully: ", accountCodeDir)
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), body)
	},
}

func init() {
	getCmd.Flags().StringVarP(&accountID, "id", "i", "", "id of the global code account you want to display")
	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	getCmd.Flags().BoolVarP(&createFile, "create-file", "", false, "create a file that contains account global code")

	AccountGlobalCodeCmd.AddCommand(getCmd)
}
