package modification

import (
	"encoding/json"
	"strconv"
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
	mockfunction_we.APIModification()

	m.Run()
}

var testModification models.Modification
var testModification_ []models.Modification
var testModificationList []models.Modification

func TestModificationCommand(t *testing.T) {
	successOutput, _ := utils.ExecuteCommand(ModificationCmd, "--campaign-id="+strconv.Itoa(100000))
	assert.Contains(t, successOutput, "Manage your modifications")
}

func TestModificationHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(ModificationCmd, "--help")
	assert.Contains(t, output, "Manage your modifications")
}

func TestModificationEditCommand(t *testing.T) {
	failOutput, _ := utils.ExecuteCommand(ModificationCmd, "edit", "--campaign-id=100000")
	assert.Contains(t, failOutput, "Error: required flag(s) \"data-raw\", \"id\" not set")

	output, _ := utils.ExecuteCommand(ModificationCmd, "edit", "--campaign-id=100000", "--id=120003", "--data-raw='{\"name\":\"testCampaignName1\",\"type\":\"ab\",\"url\":\"https://abtasty1.com\",\"description\":\"testCampaignDescription1\",\"global_code\":\"console.log(\"Hello World!\")\"}'")

	err := json.Unmarshal([]byte(output), &testModification)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestElementModification, testModification)
}

func TestModificationGetCommand(t *testing.T) {
	failOutput, err := utils.ExecuteCommand(ModificationCmd, "get", "--campaign-id=100000")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(ModificationCmd, "get", "--id=120003", "--campaign-id=100000")
	err = json.Unmarshal([]byte(successOutput), &testModification)

	assert.Nil(t, err)
	assert.Equal(t, mockfunction_we.TestModification.Data.Modifications[0], testModification)
}

func TestModificationListCommand(t *testing.T) {
	output, err := utils.ExecuteCommand(ModificationCmd, "list", "--campaign-id=100000")
	err = json.Unmarshal([]byte(output), &testModificationList)

	assert.Nil(t, err)
	assert.Equal(t, []models.Modification{mockfunction_we.TestModificationsJS, mockfunction_we.TestModificationsCSS}, testModificationList)
}

func TestModificationDeleteCommand(t *testing.T) {
	failOutput, _ := utils.ExecuteCommand(ModificationCmd, "delete", "--campaign-id=100000")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(ModificationCmd, "delete", "--campaign-id=100000", "--id=120003")
	assert.Equal(t, "Modification 120003 deleted\n", successOutput)
}
