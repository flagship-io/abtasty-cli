/*
Copyright © 2022 Flagship Team flagship@abtasty.com

*/

package working_directory

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
	"github.com/spf13/cobra"
)

var mu sync.Mutex

// SetCmd represents the working-dir command
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set working directory",
	Long:  `Set working directory to pull code`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := config.CheckWorkingDirectory(Path)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		absPath, err := filepath.Abs(Path)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}

		err = config.SetWorkingDir(utils.WEB_EXPERIMENTATION, absPath)
		if err != nil {
			log.Fatalf("error occurred: %s", err)
		}
	},
}

func init() {
	WorkingDirectoryCmd.AddCommand(SetCmd)
}
