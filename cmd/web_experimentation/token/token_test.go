package token

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/models"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	mockfunction "github.com/flagship-io/abtasty-cli/utils/mock_function"
	mockfunction_we "github.com/flagship-io/abtasty-cli/utils/mock_function/web_experimentation"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var testToken models.Token

func TestMain(m *testing.M) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockfunction.SetMock(&http_request.ResourceRequester)
	mockfunction_we.APIToken()
	m.Run()
}

func TestTokenCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(TokenCmd)
	assert.Contains(t, output, "Manage your token\n")
}

func TestTokenHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(TokenCmd, "--help")
	assert.Contains(t, output, "Manage your token\n")
}
