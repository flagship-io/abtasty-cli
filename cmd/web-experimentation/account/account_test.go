package account

import (
	"encoding/json"
	"testing"

	"github.com/flagship-io/abtasty-cli/models"
	models_we "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	mockfunction "github.com/flagship-io/abtasty-cli/utils/mock_function"
	"github.com/flagship-io/abtasty-cli/utils/mock_function/web_experimentation"
	mockfunction_we "github.com/flagship-io/abtasty-cli/utils/mock_function/web_experimentation"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	defer mockfunction_we.InitMockAuth()

	mockfunction.SetMock(&http_request.ResourceRequester)

	mockfunction_we.APIAccount()

	m.Run()
}

var testAccount models.AccountJSON
var testAccounts []models_we.AccountWE

func TestAccountCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(AccountCmd)
	assert.Contains(t, output, "Manage your CLI authentication\n\nUsage:\n  account [use|list|current]")
}

func TestAccountHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(AccountCmd, "--help")
	assert.Contains(t, output, "Manage your CLI authentication\n\nUsage:\n  account [use|list|current]")
}

func TestAccountUseCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(AccountCmd, "use")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(AccountCmd, "use", "-i=account_id")
	assert.Equal(t, "Account ID set to : account_id\n", successOutput)
}

func TestAccountCurrentCommand(t *testing.T) {

	output, _ := utils.ExecuteCommand(AccountCmd, "current")

	err := json.Unmarshal([]byte(output), &testAccount)

	assert.Nil(t, err)

	assert.Equal(t, web_experimentation.TestAccount.AccountID, testAccount.AccountID)
}

func TestAccountListCommand(t *testing.T) {

	output, _ := utils.ExecuteCommand(AccountCmd, "list")

	err := json.Unmarshal([]byte(output), &testAccounts)

	assert.Nil(t, err)

	assert.Equal(t, web_experimentation.TestAccountGlobalCode.Id, testAccounts[0].Id)
}
