/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package account_global_code

import (
	"fmt"
	"log"

	"github.com/flagship-io/abtasty-cli/utils/config"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

var createFile bool
var override bool

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <account-id> | --id <account-id>]",
	Short: "Get account global code",
	Long:  `Get account global code`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.AccountGlobalCodeRequester.HTTPGetAccountGlobalCode(AccountID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		if createFile && len(body) > 0 {
			_, err := config.AccountGlobalCodeDirectory(httprequest.AccountGlobalCodeRequester.WorkingDir, AccountID, body, override)
			if err != nil {
				log.Fatalf("error occurred: %s", err)
			}

			return
		}

		if len(body) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), body)
			return
		}
	},
}

func init() {
	getCmd.Flags().StringVarP(&AccountID, "id", "i", "", "id of the account you want to display")
	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	getCmd.Flags().BoolVarP(&createFile, "create-file", "", false, "create a file that contains account global code")
	getCmd.Flags().BoolVarP(&override, "override", "", false, "override existing account global code file")

	AccountGlobalCodeCmd.AddCommand(getCmd)
}
