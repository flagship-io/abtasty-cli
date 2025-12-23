/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package cmd

import (
	"os"

	feature_experimentation "github.com/flagship-io/abtasty-cli/cmd/feature-experimentation"
	web_experimentation "github.com/flagship-io/abtasty-cli/cmd/web-experimentation"

	"github.com/flagship-io/abtasty-cli/cmd/version"
	"github.com/flagship-io/abtasty-cli/internal/utils/config"
	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "abtasty-cli",
	Aliases: []string{"abt", "abtasty"},
	Short:   "abtasty-cli let you manage your campaigns, project, flags, etc... on both product web experimentation and feature experimentation",
	Long: `abtasty-cli is the main command, used to manage campaigns, projects, flags, etc... on both product web experimentation and feature experimentation
	
	- Web Experimentation is a customer experience optimization product that blends advanced testing with simple experience building to reach conversion goals confidently and quickly.

	- Feature Experimentation and Rollout is a feature flagging platform for modern developers. 
	Separate code deployments from feature releases to accelerate development cycles and mitigate risks.
	
	Complete documentation is available at https://docs.developers.flagship.io/docs/abtasty-cli-command-line-interface`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Root().Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubCommandPalettes() {
	RootCmd.AddCommand(version.VersionCmd)
	RootCmd.AddCommand(feature_experimentation.FeatureExperimentationCmd)
	RootCmd.AddCommand(web_experimentation.WebExperimentationCmd)
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&common.OutputFormat, "output-format", "", config.OutputFormat, "output format for the get and list subcommands for AB Tasty resources. Only 3 format are possible: table, json, json-pretty")
	RootCmd.PersistentFlags().StringVarP(&common.UserAgent, "user-agent", "", config.DefaultUserAgent, "custom user agent")

	viper.BindPFlag("output_format", RootCmd.PersistentFlags().Lookup("output-format"))

	addSubCommandPalettes()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Find home directory.
	_, err := config.CheckABTastyHomeDirectory()
	cobra.CheckErr(err)
}
