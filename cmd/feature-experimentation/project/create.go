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

func CreateProject(dataRaw []byte) ([]byte, error) {
	body, err := httprequest.ProjectRequester.HTTPCreateProject(dataRaw)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [-n <name> | --name=<name>]",
	Short: "Create a project",
	Long:  `Create a project`,
	Run: func(cmd *cobra.Command, args []string) {
		projectRequest := models.Project{
			Name: ProjectName,
		}

		projectRequestJSON, err := json.Marshal(projectRequest)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		body, err := CreateProject(projectRequestJSON)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", body)
	},
}

func init() {

	createCmd.Flags().StringVarP(&ProjectName, "name", "n", "", "name of the project you want to create")

	if err := createCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	ProjectCmd.AddCommand(createCmd)
}
