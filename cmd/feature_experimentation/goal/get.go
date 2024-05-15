/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package goal

import (
	"log"

	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents get command
var getCmd = &cobra.Command{
	Use:   "get [-i <goal-id> | --id=<goal-id>]",
	Short: "Get a goal",
	Long:  `Get a goal`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := httprequest.GoalRequester.HTTPGetGoal(GoalID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Label", "Type", "Operator", "Value"}, body, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().StringVarP(&GoalID, "id", "i", "", "id of the goal you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	GoalCmd.AddCommand(getCmd)
}
