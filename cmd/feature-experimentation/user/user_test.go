package user

import (
	"encoding/json"
	"testing"

	models "github.com/flagship-io/abtasty-cli/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	mockfunction "github.com/flagship-io/abtasty-cli/utils/mock_function"
	mockfunction_fe "github.com/flagship-io/abtasty-cli/utils/mock_function/feature_experimentation"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockfunction.SetMock(&http_request.ResourceRequester)
	mockfunction_fe.APIUser()

	m.Run()
}

var testUserList []models.User

func TestUserCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(UserCmd)
	assert.Contains(t, output, "Manage your users")
}

func TestUserHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(UserCmd, "--help")
	assert.Contains(t, output, "Manage your users")
}

func TestUserListCommand(t *testing.T) {

	output, _ := utils.ExecuteCommand(UserCmd, "list")

	err := json.Unmarshal([]byte(output), &testUserList)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_fe.TestUserList, testUserList)
}

func TestUserCreateCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(UserCmd, "create")
	assert.Contains(t, failOutput, "Error: required flag(s) \"data-raw\" not set")

	successOutput, _ := utils.ExecuteCommand(UserCmd, "create", `--data-raw=[{"email":"example@abtasty.com","role":"PROJECT_MANAGER"},{"email":"example1@abtasty.com","role":"SUPER_ADMIN"}]`)
	assert.Contains(t, "users created\n", successOutput)
}

func TestUserEditCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(UserCmd, "edit")
	assert.Contains(t, failOutput, "Error: required flag(s) \"data-raw\" not set")

	successOutput, _ := utils.ExecuteCommand(UserCmd, "edit", "--data-raw=[{\"email\":\"example@abtasty.com\",\"role\":\"PROJECT_MANAGER\"},{\"email\":\"example1@abtasty.com\",\"role\":\"SUPER_ADMIN\"}]")
	assert.Contains(t, "users edited\n", successOutput)
}

func TestUserDeleteCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(UserCmd, "delete")
	assert.Contains(t, failOutput, "Error: required flag(s) \"email\" not set")

	successOutput, _ := utils.ExecuteCommand(UserCmd, "delete", "--email=example@abtasty.com")
	assert.Equal(t, "Email deleted\n", successOutput)
}
