/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package folder

import (
	"github.com/spf13/cobra"
)

var (
	FolderID int
	DataRaw  string
)

// FolderCmd represents the folder command
var FolderCmd = &cobra.Command{
	Use:   "folder [create|edit|list|get|delete]",
	Short: "Manage your folders",
	Long:  `Manage your folders`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
