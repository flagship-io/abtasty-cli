package feature_experimentation

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/utils"
	mockfunction_fe "github.com/flagship-io/abtasty-cli/utils/mock_function/feature_experimentation"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	defer mockfunction_fe.InitMockAuth()

	m.Run()
}

func TestFECommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(FeatureExperimentationCmd)
	assert.Contains(t, output, "Manage resources related to the feature experimentation product")
}

func TestFEHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(FeatureExperimentationCmd, "--help")
	assert.Contains(t, output, "Manage resources related to the feature experimentation product")
}
