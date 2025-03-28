/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package user

import (
	"encoding/json"
	"fmt"
	"log"

	models "github.com/flagship-io/abtasty-cli/models/feature_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
)

// createCmd represents the add command
var createCmd = &cobra.Command{
	Use:   "create [-d <data-raw> | --data <data-raw>]",
	Short: "Create a user with right",
	Long:  `Create a user with right`,
	Run: func(cmd *cobra.Command, args []string) {
		var userToAdd []models.User

		err := json.Unmarshal([]byte(DataRaw), &userToAdd)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		for _, user := range userToAdd {
			if user.Role == "SUPER_ADMIN" || user.Role == "PROJECT_MANAGER" || user.Role == "MEMBER" || user.Role == "GUEST" {
				continue
			}

			log.Fatalf("error occurred: Role %s for the user %s is not supported, we only support: SUPER_ADMIN, PROJECT_MANAGER, MEMBER, GUEST", user.Role, user.Email)
		}

		_, err = httprequest.UserRequester.HTTPBatchUpdateUsers(DataRaw)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), "users created")
	},
}

func init() {
	createCmd.Flags().StringVarP(&DataRaw, "data-raw", "d", "", "raw data contains all the info to create your user with right, check the doc for details")

	if err := createCmd.MarkFlagRequired("data-raw"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	UserCmd.AddCommand(createCmd)
}
