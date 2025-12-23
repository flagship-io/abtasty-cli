package web_preview

import (
	"testing"

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
	defer mockfunction_we.InitMockAuth()

	mockfunction.SetMock(&http_request.ResourceRequester)

	mockfunction_we.APICampaign()

	m.Run()
}

func TestCampaignTargetingCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(WebPreviewCmd)
	assert.Contains(t, output, "Open web preview")
}

func TestCampaignTargetingHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(WebPreviewCmd, "--help")
	assert.Contains(t, output, "Open web preview")
}

func TestCampaignTargetingOpenCommand(t *testing.T) {
	failOutput, _ := utils.ExecuteCommand(WebPreviewCmd, "open")
	assert.Contains(t, failOutput, "Error: required flag(s) \"campaign-id\", \"variation-id\" not set")

}
