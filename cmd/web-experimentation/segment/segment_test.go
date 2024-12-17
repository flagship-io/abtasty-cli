package segment

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
	mockfunction_we.APISegment()

	m.Run()
}

var testSegment models.Audience
var testSegmentList []models.Audience

func TestSegmentCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(SegmentCmd)
	assert.Contains(t, output, "Manage your segments")
}

func TestSegmentHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(SegmentCmd, "--help")
	assert.Contains(t, output, "Manage your segments")
}

func TestSegmentGetCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(SegmentCmd, "get")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(SegmentCmd, "get", "--id="+"segment-id")

	err := json.Unmarshal([]byte(successOutput), &testSegment)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestSegment, testSegment)
}

func TestSegmentListCommand(t *testing.T) {

	output, _ := utils.ExecuteCommand(SegmentCmd, "list")

	err := json.Unmarshal([]byte(output), &testSegmentList)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestSegmentList, testSegmentList)
}
