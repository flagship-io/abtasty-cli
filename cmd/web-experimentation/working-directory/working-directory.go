/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/

package working_directory

import (
	"log"

	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/spf13/cobra"
)

var Path string

// WorkingDirectoryCmd represents the working-dir command
var WorkingDirectoryCmd = &cobra.Command{
	Use:     "working-directory [set]",
	Aliases: []string{"working-dir"},
	Short:   "Manage your working directory",
	Long:    `Manage your working directory to pull code`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	WorkingDirectoryCmd.PersistentFlags().StringVarP(&Path, "path", "", utils.DefaultGlobalCodeWorkingDir(), "path to set for working dir")

	if err := WorkingDirectoryCmd.MarkPersistentFlagRequired("path"); err != nil {
		log.Fatalf("error occurred: %v", err)
	}
}
