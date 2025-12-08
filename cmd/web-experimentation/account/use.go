/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package account

import (
	"fmt"
	"log"

	"github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	"github.com/flagship-io/abtasty-cli/internal/utils/config"
	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
	"github.com/spf13/cobra"
)

// getCmd represents the list command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a specific account id",
	Long:  `Use a specific account id`,
	Run: func(cmd *cobra.Command, args []string) {
		if AccountID == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "required flag account-id")
			return
		}

		err := config.SetAccountID(utils.WEB_EXPERIMENTATION, AccountID)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		dir, derr := utils.DefaultGlobalCodeWorkingDir()
		if derr != nil {
			log.Fatalf("error occurred: %s", derr)
		}

		err = config.SetWorkingDir(utils.WEB_EXPERIMENTATION, dir)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		currentUser, err := common.HTTPGetIdentifierWE()
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		if currentUser.LastAccount != (web_experimentation.AccountWE{}) {
			err := config.SetIdentifier(utils.WEB_EXPERIMENTATION, currentUser.LastAccount.Identifier)
			if err != nil {
				log.Fatalf("error occurred: %s", err)
			}

			err = config.SetEmail(utils.WEB_EXPERIMENTATION, currentUser.Email)
			if err != nil {
				log.Fatalf("error occurred: %s", err)
			}
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Account ID set to : "+AccountID)

	},
}

func init() {
	useCmd.Flags().StringVarP(&AccountID, "id", "i", "", "account id of the credentials you want to manage")

	if err := useCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	AccountCmd.AddCommand(useCmd)
}
