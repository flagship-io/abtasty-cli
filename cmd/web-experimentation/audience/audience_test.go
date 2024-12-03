package audience

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
	mockfunction_we.APIAudience()

	m.Run()
}

var testAudience models.Audience
var testAudienceList []models.Audience

func TestAudienceCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(AudienceCmd)
	assert.Contains(t, output, "Manage your audiences")
}

func TestAudienceHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(AudienceCmd, "--help")
	assert.Contains(t, output, "Manage your audiences")
}

func TestAudienceGetCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(AudienceCmd, "get")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(AudienceCmd, "get", "--id="+"testAudienceId")

	err := json.Unmarshal([]byte(successOutput), &testAudience)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestAudience, testAudience)
}

func TestAudienceListCommand(t *testing.T) {

	output, _ := utils.ExecuteCommand(AudienceCmd, "list")

	err := json.Unmarshal([]byte(output), &testAudienceList)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestAudienceList, testAudienceList)
}
