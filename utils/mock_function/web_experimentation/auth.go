package web_experimentation

import (
	"log"
	"os"

	"github.com/flagship-io/abtasty-cli/models"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
)

var TestAuth = models.Auth{
	Username:     "test_auth",
	ClientID:     "testAuthClientID",
	ClientSecret: "testAuthClientSecret",
	Token:        "testAccessToken",
}

func InitMockAuth() {
	credPath, err := config.CredentialPath(utils.WEB_EXPERIMENTATION, "test_auth")
	if err != nil {
		log.Fatalf("error occurred: %s", err)
	}

	os.Remove(credPath)
}
