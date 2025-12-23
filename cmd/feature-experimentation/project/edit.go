/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package project

import (
	"encoding/json"
	"fmt"
	"log"

	models "github.com/flagship-io/abtasty-cli/internal/models/feature_experimentation"
	httprequest "github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	"github.com/spf13/cobra"
)

func EditProject(id string, dataRaw []byte) ([]byte, error) {
	body, err := httprequest.ProjectRequester.HTTPEditProject(id, dataRaw)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [-i <project-id> | --id=<project-id>] [-n <name> | --name=<name>]",
	Short: "Edit a project",
	Long:  `Edit a project`,
	Run: func(cmd *cobra.Command, args []string) {
		projectRequest := models.Project{
			Name: ProjectName,
		}

		projectRequestJSON, err := json.Marshal(projectRequest)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		body, err := EditProject(ProjectId, projectRequestJSON)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", body)
	},
}

func init() {

	editCmd.Flags().StringVarP(&ProjectId, "id", "i", "", "id of the project you want to edit")

	editCmd.Flags().StringVarP(&ProjectName, "name", "n", "", "name you want to set for the project")

	if err := editCmd.MarkFlagRequired("id"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if err := editCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	ProjectCmd.AddCommand(editCmd)

}
