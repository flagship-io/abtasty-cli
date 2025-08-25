/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package web_experimentation

import (
	"log"
	"os"

	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/account"
	account_global_code "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/account-global-code"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/audience"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/auth"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/campaign"
	campaign_global_code "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/campaign-global-code"
	campaign_targeting "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/campaign-targeting"
	favorite_url "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/favorite-url"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/folder"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/modification"
	modification_code "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/modification-code"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/resource"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/segment"
	tag_rebuild "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/tag-rebuild"
	t "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/token"

	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/trigger"
	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/variation"
	variation_global_code "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/variation-global-code"
	web_preview "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/web-preview"
	working_directory "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/working-directory"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	mockfunction_ "github.com/flagship-io/abtasty-cli/utils/mock_function"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accountID, token string

// WebExperimentationCmd represents the web experimentation command
var WebExperimentationCmd = &cobra.Command{
	Use:     "web-experimentation [auth|account|campaign|global-code|variation]",
	Aliases: []string{"web-experimentation", "web-exp", "we", "web"},
	Short:   "Manage resources related to the web experimentation product",
	Long:    `Manage resources related to the web experimentation product`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func addSubCommandPalettes() {
	WebExperimentationCmd.AddCommand(campaign.CampaignCmd)
	WebExperimentationCmd.AddCommand(variation.VariationCmd)
	WebExperimentationCmd.AddCommand(auth.AuthCmd)
	WebExperimentationCmd.AddCommand(account.AccountCmd)
	WebExperimentationCmd.AddCommand(campaign_global_code.CampaignGlobalCodeCmd)
	WebExperimentationCmd.AddCommand(account_global_code.AccountGlobalCodeCmd)
	WebExperimentationCmd.AddCommand(variation_global_code.VariationGlobalCodeCmd)
	WebExperimentationCmd.AddCommand(modification_code.ModificationCodeCmd)
	WebExperimentationCmd.AddCommand(t.TokenCmd)
	WebExperimentationCmd.AddCommand(modification.ModificationCmd)
	WebExperimentationCmd.AddCommand(working_directory.WorkingDirectoryCmd)
	WebExperimentationCmd.AddCommand(tag_rebuild.RebuildTagCmd)
	WebExperimentationCmd.AddCommand(audience.AudienceCmd)
	WebExperimentationCmd.AddCommand(segment.SegmentCmd)
	WebExperimentationCmd.AddCommand(trigger.TriggerCmd)
	WebExperimentationCmd.AddCommand(favorite_url.FavoriteUrlCmd)
	WebExperimentationCmd.AddCommand(campaign_targeting.CampaignTargetingCmd)
	WebExperimentationCmd.AddCommand(web_preview.WebPreviewCmd)
	WebExperimentationCmd.AddCommand(resource.ResourceCmd)
	WebExperimentationCmd.AddCommand(folder.FolderCmd)
}

func init() {
	addSubCommandPalettes()

	WebExperimentationCmd.PersistentFlags().StringVarP(&accountID, "account-id", "i", "", "override account ID in command")
	WebExperimentationCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "override token in command")

}

func initConfig() {
	v := viper.New()
	homeDir, _ := os.UserHomeDir()
	var requestConfig = common.RequestConfig{Product: utils.WEB_EXPERIMENTATION}

	v.SetConfigFile(homeDir + "/.abtasty/credentials/" + utils.WEB_EXPERIMENTATION + "/.cli.yaml")
	v.MergeInConfig()
	if v.GetString("current_used_credential") != "" {
		vL, err := config.ReadAuth(utils.WEB_EXPERIMENTATION, v.GetString("current_used_credential"))
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		if accountID != "" {
			v.Set("account_id", accountID)
		}

		if token != "" {
			v.Set("token", token)
		}

		v.MergeConfigMap(vL.AllSettings())
	}

	v.Unmarshal(&requestConfig)
	common.Init(requestConfig)

	r := &http_request.ResourceRequester

	if os.Getenv("ABT_ENV") == "MOCK" {
		r.Init(&mockfunction_.Auth)
		return
	}

	r.Init(&requestConfig)
}
