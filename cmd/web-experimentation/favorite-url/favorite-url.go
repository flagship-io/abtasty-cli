/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package favorite_url

import (
	"github.com/spf13/cobra"
)

var (
	FavoriteUrlID string
	DataRaw       string
)

// FavoriteUrlCmd represents the favorite url command
var FavoriteUrlCmd = &cobra.Command{
	Use:   "favorite-url [list|get]",
	Short: "Manage your favorite urls",
	Long:  `Manage your favorite url`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
