package web_experimentation

import (
	"testing"

	"github.com/flagship-io/flagship/utils"
	mockfunction_fe "github.com/flagship-io/flagship/utils/mock_function/feature_experimentation"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	defer mockfunction_fe.InitMockAuth()

	m.Run()
}

func TestAccountCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(WebExperimentationCmd)
	assert.Contains(t, output, "Manage resources related to the web experimentation product")
}

func TestAccountHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(WebExperimentationCmd, "--help")
	assert.Contains(t, output, "Manage resources related to the web experimentation product")
}
