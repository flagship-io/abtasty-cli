package auth

import (
	"encoding/json"
	"testing"

	"github.com/flagship-io/abtasty-cli/models"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	mockfunction "github.com/flagship-io/abtasty-cli/utils/mock_function"
	"github.com/flagship-io/abtasty-cli/utils/mock_function/feature_experimentation"
	mockfunction_fe "github.com/flagship-io/abtasty-cli/utils/mock_function/feature_experimentation"
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

var testAuth models.Auth
var testAuthList []models.Auth

func TestAuthCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(AuthCmd)
	assert.Contains(t, output, "Manage authentication for feature experimentation\n\nUsage:\n  authentication [login|get|list|delete]")
}

func TestAuthHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(AuthCmd, "--help")
	assert.Contains(t, output, "Manage authentication for feature experimentation\n\nUsage:\n  authentication [login|get|list|delete]")
}

/* func TestAuthLoginCommand(t *testing.T) {
	successOutput, _ := utils.ExecuteCommand(AuthCmd, "login", "-u=test_auth", "-i=testAuthClientID", "-s=testAuthClientSecret", "-a=account_id")
	assert.Equal(t, "Credential created successfully\n", successOutput)
} */

func TestAuthListCommand(t *testing.T) {

	config.CreateAuthFile(utils.FEATURE_EXPERIMENTATION, "test_auth", "testAuthClientID", "testAuthClientSecret", models.TokenResponse{AccessToken: "testAccessToken", RefreshToken: "testRefreshToken"})

	output, _ := utils.ExecuteCommand(AuthCmd, "list")

	err := json.Unmarshal([]byte(output), &testAuthList)

	assert.Nil(t, err)

	byt, err := json.Marshal(feature_experimentation.TestAuth)

	assert.Nil(t, err)

	assert.Contains(t, output, string(byt))
}

func TestAuthGetCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(AuthCmd, "get")
	assert.Contains(t, failOutput, "Error: required flag(s) \"username\" not set")

	successOutput, _ := utils.ExecuteCommand(AuthCmd, "get", "--username=test_auth")
	err := json.Unmarshal([]byte(successOutput), &testAuth)

	assert.Nil(t, err)

	assert.Equal(t, feature_experimentation.TestAuth, testAuth)
}

func TestAuthDeleteCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(AuthCmd, "delete")
	assert.Contains(t, failOutput, "Error: required flag(s) \"username\" not set")

	output, _ := utils.ExecuteCommand(AuthCmd, "delete", "--username=test_auth")

	assert.Contains(t, output, "Credential deleted successfully")
}
