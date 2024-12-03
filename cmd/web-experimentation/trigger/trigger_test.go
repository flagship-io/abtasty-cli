package trigger

import (
	"encoding/json"
	"testing"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	mockfunction "github.com/flagship-io/abtasty-cli/utils/mock_function"
	mockfunction_we "github.com/flagship-io/abtasty-cli/utils/mock_function/web_experimentation"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockfunction.SetMock(&http_request.ResourceRequester)
	mockfunction_we.APITrigger()

	m.Run()
}

var testTrigger models.Audience
var testTriggerList []models.Audience

func TestTriggerCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(TriggerCmd)
	assert.Contains(t, output, "Manage your triggers")
}

func TestTriggerHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(TriggerCmd, "--help")
	assert.Contains(t, output, "Manage your triggers")
}

func TestTriggerGetCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(TriggerCmd, "get")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(TriggerCmd, "get", "--id="+"trigger-id")

	err := json.Unmarshal([]byte(successOutput), &testTrigger)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestTrigger, testTrigger)
}

func TestTriggerListCommand(t *testing.T) {

	output, _ := utils.ExecuteCommand(TriggerCmd, "list")

	err := json.Unmarshal([]byte(output), &testTriggerList)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestTriggerList, testTriggerList)
}
