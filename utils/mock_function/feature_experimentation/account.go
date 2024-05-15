package feature_experimentation

import (
	"github.com/flagship-io/abtasty-cli/models"
)

var TestAccount = models.AccountJSON{
	CurrentUsedCredential: "test_auth",
	AccountID:             "account_id",
	AccountEnvironmentID:  "account_environment_id",
}
