/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package account

import (
	"fmt"
	"log"

	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
	"github.com/spf13/cobra"
)

// getCmd represents the list command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a specific account id",
	Long:  `Use a specific account id`,
	Run: func(cmd *cobra.Command, args []string) {
		if AccountID == "" {
			fmt.Fprintln(cmd.OutOrStderr(), "required flag account-id")
			return
		}

		err := config.SetAccountID(utils.FEATURE_EXPERIMENTATION, AccountID)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Account ID set to: "+AccountID)

	},
}

func init() {
	useCmd.Flags().StringVarP(&AccountID, "id", "i", "", "account id of the credentials")

	if err := useCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	AccountCmd.AddCommand(useCmd)
}
