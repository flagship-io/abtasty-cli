/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package account

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
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Get current running account ID",
	Long:  `Get current running account ID`,
	Run: func(cmd *cobra.Command, args []string) {

		var accountYaml models.AccountYaml
		var account models.AccountJSON

		credPath, err := config.CredentialPath(utils.FEATURE_EXPERIMENTATION, utils.HOME_CLI)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		yamlFile, err := os.ReadFile(credPath)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		// Unmarshal the YAML data into the struct
		err = yaml.Unmarshal(yamlFile, &accountYaml)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		account.CurrentUsedCredential = accountYaml.CurrentUsedCredential
		account.AccountID = accountYaml.AccountID

		utils.FormatItem([]string{"CurrentUsedCredential", "AccountID"}, account, viper.GetString("output_format"), cmd.OutOrStdout())

	},
}

func init() {
	AccountCmd.AddCommand(currentCmd)
}
