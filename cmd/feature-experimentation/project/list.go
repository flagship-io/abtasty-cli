/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package project

import (
	"log"

	"github.com/flagship-io/abtasty-cli/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	httprequest "github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListProjects() ([]feature_experimentation.Project, error) {
	return httprequest.ProjectRequester.HTTPListProject()
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  `List all projects`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := ListProjects()
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {
	ProjectCmd.AddCommand(listCmd)
}
