/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package accountenvironment

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
	Short: "use a specific account environment id",
	Long:  `use a specific account environment id`,
	Run: func(cmd *cobra.Command, args []string) {
		if AccountEnvironmentID == "" {
			fmt.Fprintln(cmd.OutOrStderr(), "required flag account-id or account-environment-id")
			return
		}

		err := config.SetAccountEnvID(utils.FEATURE_EXPERIMENTATION, AccountEnvironmentID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Account Environment ID set to : "+AccountEnvironmentID)

	},
}

func init() {
	useCmd.Flags().StringVarP(&AccountEnvironmentID, "id", "i", "", "account env id of the credentials you want to manage")

	if err := useCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	AccountEnvironmentCmd.AddCommand(useCmd)
}
