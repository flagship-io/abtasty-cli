/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package auth

import (
	"log"
	"os"

	"github.com/flagship-io/abtasty-cli/models"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// getCmd represents the list command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get an auth credential for web experimentation",
	Long:  `get an auth credential for web experimentation`,
	Run: func(cmd *cobra.Command, args []string) {

		var authYaml models.AuthYaml
		var auth models.Auth

		credPath, err := config.CredentialPath(utils.WEB_EXPERIMENTATION, Username)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		yamlFile, err := os.ReadFile(credPath)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		// Unmarshal the YAML data into the struct
		err = yaml.Unmarshal(yamlFile, &authYaml)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		auth.Username = authYaml.Username
		auth.ClientID = authYaml.ClientID
		auth.ClientSecret = authYaml.ClientSecret
		auth.Token = authYaml.Token

		utils.FormatItem([]string{"Username", "Token"}, auth, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	getCmd.Flags().StringVarP(&Username, "username", "u", "", "username of the credentials you want to display")

	if err := getCmd.MarkFlagRequired("username"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	AuthCmd.AddCommand(getCmd)
}
