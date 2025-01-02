package campaign_targeting

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
	output, _ := utils.ExecuteCommand(CampaignTargetingCmd)
	assert.Contains(t, output, "Get campaign targeting")
}

func TestCampaignTargetingHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(CampaignTargetingCmd, "--help")
	assert.Contains(t, output, "Get campaign targeting")
}

func TestCampaignTargetingGetCommand(t *testing.T) {
	failOutput, _ := utils.ExecuteCommand(CampaignTargetingCmd, "get")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(CampaignTargetingCmd, "get", "-i=100000")
	assert.Equal(t, "{\"segment_ids\":[],\"url_scopes\":[{\"condition\":\"is\",\"value\":\"https://abtasty.com\"},{\"condition\":\"is not\",\"value\":\"https://abtasty.com\"}],\"favorite_url_scopes\":[],\"selector_scopes\":[],\"code_scope\":{\"value\":\"\"},\"element_appears_after_page_load\":false,\"triggers_ids\":[],\"targeting_frequency\":{}}\n", successOutput)
}
