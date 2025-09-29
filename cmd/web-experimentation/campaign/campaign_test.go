package campaign

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
	mockfunction_we.APICampaign()

	m.Run()
}

var testCampaign models.CampaignWE
var testCampaignList []models.CampaignWE

func TestCampaignCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(CampaignCmd)
	assert.Contains(t, output, "Manage your campaigns")
}

func TestCampaignHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(CampaignCmd, "--help")
	assert.Contains(t, output, "Manage your campaigns")
}

func TestCampaignEditCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(CampaignCmd, "edit")
	assert.Contains(t, failOutput, "Error: required flag(s) \"data-raw\", \"id\" not set")

	output, _ := utils.ExecuteCommand(CampaignCmd, "edit", "--id=100000", "--data-raw={\"name\":\"testCampaignName1\",\"type\":\"ab\",\"url\":\"https://abtasty1.com\",\"description\":\"testCampaignDescription1\"}")
	err := json.Unmarshal([]byte(output), &testCampaign)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestCampaign, testCampaign)
}

func TestCampaignGetCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(CampaignCmd, "get")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(CampaignCmd, "get", "--id="+strconv.Itoa(100000))

	err := json.Unmarshal([]byte(successOutput), &testCampaign)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestCampaign, testCampaign)
}

func TestCampaignListCommand(t *testing.T) {

	output, _ := utils.ExecuteCommand(CampaignCmd, "list")

	err := json.Unmarshal([]byte(output), &testCampaignList)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestCampaignList, testCampaignList)
}

func TestCampaignDeleteCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(CampaignCmd, "delete")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(CampaignCmd, "delete", "--id=100000")
	assert.Equal(t, "Campaign 100000 deleted\n", successOutput)
}

func TestCampaignSwitchCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(CampaignCmd, "switch")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\", \"status\" not set")

	failOutput1, _ := utils.ExecuteCommand(CampaignCmd, "switch", "--id=100000", "--status=notKnown")
	assert.Equal(t, "Status can only have 2 values: active or paused\n", failOutput1)

	successOutput, _ := utils.ExecuteCommand(CampaignCmd, "switch", "--id=100000", "--status=active")
	assert.Equal(t, "campaign status set to active\n", successOutput)
}
