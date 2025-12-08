package web_experimentation

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/internal/utils"
	mockfunction_we "github.com/flagship-io/abtasty-cli/internal/utils/mock_function/web_experimentation"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	defer mockfunction_we.InitMockAuth()

	m.Run()
}

func TestWECommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(WebExperimentationCmd)
	assert.Contains(t, output, "Manage resources related to the web experimentation product")
}

func TestWEHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(WebExperimentationCmd, "--help")
	assert.Contains(t, output, "Manage resources related to the web experimentation product")
}
