/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package goal

import (
	"fmt"
	"log"

	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

func DeleteGoal(id string) (string, error) {
	err := httprequest.GoalRequester.HTTPDeleteGoal(id)
	if err != nil {
		return "", err
	}
	return "Goal deleted", nil
}

// deleteCmd represents delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [-i <goal-id> | --id=<goal-id>]",
	Short: "Delete a goal",
	Long:  `Delete a goal`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := DeleteGoal(GoalID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), resp)

	},
}

func init() {
	deleteCmd.Flags().StringVarP(&GoalID, "id", "i", "", "id of the goal you want to delete")

	if err := deleteCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	GoalCmd.AddCommand(deleteCmd)
}
