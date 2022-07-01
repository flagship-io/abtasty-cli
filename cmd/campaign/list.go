/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package campaign

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listCampaign() string {
	return "get All Campaign"
}

// createCmd represents the create command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "this list campaign",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(listCampaign())
	},
}

func init() {

	// Here you will define your flags and configuration settings.
	CampaignCmd.AddCommand(listCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}