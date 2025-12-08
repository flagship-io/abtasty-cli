/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package project

import (
	"log"

	"github.com/flagship-io/abtasty-cli/internal/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetProject(id string) (feature_experimentation.Project, error) {
	body, err := httprequest.ProjectRequester.HTTPGetProject(id)
	if err != nil {
		return feature_experimentation.Project{}, err
	}
	return body, nil
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [-i <project-id> | --id=<project-id>]",
	Short: "Get a project",
	Long:  `Get a project`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := GetProject(ProjectId)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		utils.FormatItem([]string{"Id", "Name"}, body, viper.GetString("output_format"), cmd.OutOrStdout())
	},
}

func init() {

	getCmd.Flags().StringVarP(&ProjectId, "id", "i", "", "id of the project you want to display")

	if err := getCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
	ProjectCmd.AddCommand(getCmd)

}
