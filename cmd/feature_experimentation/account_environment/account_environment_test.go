package accountenvironment

import (
	"encoding/json"
	"testing"

	"github.com/flagship-io/flagship/models"
	"github.com/flagship-io/flagship/utils"
	"github.com/flagship-io/flagship/utils/http_request"
	mockfunction "github.com/flagship-io/flagship/utils/mock_function"
	mockfunction_fe "github.com/flagship-io/flagship/utils/mock_function/feature_experimentation"
	"github.com/flagship-io/flagship/utils/mock_function/web_experimentation"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	defer mockfunction_fe.InitMockAuth()

	mockfunction.SetMock(&http_request.ResourceRequester)

	mockfunction_fe.APIToken()

	m.Run()
}

var testAccount models.AccountJSON

func TestAccountEnvironmentCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(AccountEnvironmentCmd)
	assert.Contains(t, output, "Manage your CLI authentication\n\nUsage:\n  account-environment [use|list|current]")
}

func TestAccountEnvironmentHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(AccountEnvironmentCmd, "--help")
	assert.Contains(t, output, "Manage your CLI authentication\n\nUsage:\n  account-environment [use|list|current]")
}

func TestAccountEnvironmentUseCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(AccountEnvironmentCmd, "use")
	assert.Contains(t, failOutput, "Error: required flag(s) \"account-environment-id\" not set")

	successOutput, _ := utils.ExecuteCommand(AccountEnvironmentCmd, "use", "--account-environment-id=account_environment_id")
	assert.Equal(t, "Account Environment ID set to : account_environment_id\n", successOutput)
}

func TestAccountEnvironmentCurrentCommand(t *testing.T) {

	output, _ := utils.ExecuteCommand(AccountEnvironmentCmd, "current")

	err := json.Unmarshal([]byte(output), &testAccount)

	assert.Nil(t, err)

	assert.Equal(t, web_experimentation.TestAccount.AccountEnvironmentID, testAccount.AccountEnvironmentID)
}