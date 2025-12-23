package account_environment

import (
	"encoding/json"
	"testing"

	"github.com/flagship-io/abtasty-cli/internal/models"
	models_fe "github.com/flagship-io/abtasty-cli/internal/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	"github.com/flagship-io/abtasty-cli/internal/utils/http_request"
	mockfunction "github.com/flagship-io/abtasty-cli/internal/utils/mock_function"
	mockfunction_fe "github.com/flagship-io/abtasty-cli/internal/utils/mock_function/feature_experimentation"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	defer mockfunction_fe.InitMockAuth()

	mockfunction.SetMock(&http_request.ResourceRequester)

	mockfunction_fe.APIAccountEnvironment()

	m.Run()
}

var testAccount models.AccountJSON
var testAccountEnvironment models_fe.AccountEnvironmentFE
var testAccountEnvironmentList []models_fe.AccountEnvironmentFE

func TestAccountEnvironmentCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(AccountEnvironmentCmd)
	assert.Contains(t, output, "Manage your account environment\n\nUsage:\n  account-environment [use|list|current]")
}

func TestAccountEnvironmentHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(AccountEnvironmentCmd, "--help")
	assert.Contains(t, output, "Manage your account environment\n\nUsage:\n  account-environment [use|list|current]")
}

func TestAccountEnvironmentUseCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(AccountEnvironmentCmd, "use")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(AccountEnvironmentCmd, "use", "--id=account_environment_id")
	assert.Equal(t, "Account Environment ID set to: account_environment_id\n", successOutput)
}

func TestAccountEnvironmentCurrentCommand(t *testing.T) {

	utils.ExecuteCommand(AccountEnvironmentCmd, "use", "--id=account_environment_id")

	output, _ := utils.ExecuteCommand(AccountEnvironmentCmd, "current")

	err := json.Unmarshal([]byte(output), &testAccount)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_fe.TestAccountEnvironment.Id, testAccount.AccountEnvironmentID)
}

func TestAccountEnvironmentListCommand(t *testing.T) {

	output, _ := utils.ExecuteCommand(AccountEnvironmentCmd, "list")

	err := json.Unmarshal([]byte(output), &testAccountEnvironmentList)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_fe.TestAccountEnvironmentList, testAccountEnvironmentList)
}
