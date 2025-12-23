/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package feature_experimentation

import (
	"log"
	"os"

	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/account"
	account_environment "github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/account-environment"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/analyze"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/auth"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/campaign"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/flag"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/goal"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/panic"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/project"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/resource"
	targeting_key "github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/targeting-key"
	t "github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/token"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/user"
	"github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/variation"
	variation_group "github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/variation-group"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	mockfunction_ "github.com/flagship-io/abtasty-cli/utils/mock_function"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accountID, accountEnvID, token string

// FeatureExperimentationCmd represents the feature experimentation command
var FeatureExperimentationCmd = &cobra.Command{
	Use:     "feature-experimentation [auth|account|account-environment|project|campaign|flag|goal|targeting-key|variation-group|variation]",
	Aliases: []string{"feature-experimentation", "feature-exp", "fe", "feat-exp", "feature"},
	Short:   "Manage resources related to the feature experimentation product",
	Long:    `Manage resources related to the feature experimentation product`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func addSubCommandPalettes() {
	FeatureExperimentationCmd.AddCommand(campaign.CampaignCmd)
	FeatureExperimentationCmd.AddCommand(project.ProjectCmd)
	FeatureExperimentationCmd.AddCommand(panic.PanicCmd)
	FeatureExperimentationCmd.AddCommand(user.UserCmd)
	FeatureExperimentationCmd.AddCommand(variation_group.VariationGroupCmd)
	FeatureExperimentationCmd.AddCommand(variation.VariationCmd)
	FeatureExperimentationCmd.AddCommand(flag.FlagCmd)
	FeatureExperimentationCmd.AddCommand(goal.GoalCmd)
	FeatureExperimentationCmd.AddCommand(targeting_key.TargetingKeyCmd)
	FeatureExperimentationCmd.AddCommand(analyze.AnalyzeCmd)
	FeatureExperimentationCmd.AddCommand(resource.ResourceCmd)
	FeatureExperimentationCmd.AddCommand(auth.AuthCmd)
	FeatureExperimentationCmd.AddCommand(account.AccountCmd)
	FeatureExperimentationCmd.AddCommand(t.TokenCmd)
	FeatureExperimentationCmd.AddCommand(account_environment.AccountEnvironmentCmd)
}

func init() {
	addSubCommandPalettes()

	FeatureExperimentationCmd.PersistentFlags().StringVarP(&accountID, "account-id", "", "", "override account ID in command")
	FeatureExperimentationCmd.PersistentFlags().StringVarP(&accountEnvID, "account-env-id", "", "", "override account environment ID in command")
	FeatureExperimentationCmd.PersistentFlags().StringVarP(&token, "token", "", "", "override token in command")

}

func initConfig() {
	v := viper.New()
	homeDir, _ := os.UserHomeDir()
	var requestConfig = common.RequestConfig{Product: utils.FEATURE_EXPERIMENTATION}

	v.SetConfigFile(homeDir + "/.abtasty/credentials/" + utils.FEATURE_EXPERIMENTATION + "/.cli.yaml")
	v.MergeInConfig()
	if v.GetString("current_used_credential") != "" {
		vL, err := config.ReadAuth(utils.FEATURE_EXPERIMENTATION, v.GetString("current_used_credential"))
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
		v.MergeConfigMap(vL.AllSettings())
	}

	if accountID != "" {
		v.Set("username", "no-username")
		v.Set("account_id", accountID)
	}

	if accountEnvID != "" {
		v.Set("username", "no-username")
		v.Set("account_environment_id", accountEnvID)
	}

	if token != "" {
		v.Set("username", "no-username")
		v.Set("token", token)
	}

	v.Unmarshal(&requestConfig)
	common.Init(requestConfig)
	viper.MergeConfigMap(v.AllSettings())

	r := &http_request.ResourceRequester

	if os.Getenv("ABT_ENV") == "MOCK" {
		r.Init(&mockfunction_.Auth)
		return
	}

	r.Init(&requestConfig)
}
