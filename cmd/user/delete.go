/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package user

import (
	"log"

	httprequest "github.com/flagship-io/flagship/utils/httpRequest"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "this delete user",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := httprequest.HTTPDeleteUsers(UserEmail)
		if err != nil {
			log.Fatalf("error occured: %v", err)
		}
		log.Println("Users deleted.")
	},
}

func init() {

	deleteCmd.Flags().StringVarP(&UserEmail, "email", "e", "", "the email")

	if err := deleteCmd.MarkFlagRequired("email"); err != nil {
		log.Fatalf("error occured: %v", err)
	}

	UserCmd.AddCommand(deleteCmd)
}
