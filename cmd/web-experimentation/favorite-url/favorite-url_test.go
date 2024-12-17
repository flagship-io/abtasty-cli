package favorite_url

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
	mockfunction_we.APIFavoriteUrl()

	m.Run()
}

var testFavoriteUrl models.FavoriteURL
var testFavoriteUrlList []models.FavoriteURL

func TestFavoriteUrlCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(FavoriteUrlCmd)
	assert.Contains(t, output, "Manage your favorite url")
}

func TestFavoriteUrlHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(FavoriteUrlCmd, "--help")
	assert.Contains(t, output, "Manage your favorite url")
}

func TestFavoriteUrlGetCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(FavoriteUrlCmd, "get")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(FavoriteUrlCmd, "get", "--id="+"testFavoriteUrlId")

	err := json.Unmarshal([]byte(successOutput), &testFavoriteUrl)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestFavoriteUrl, testFavoriteUrl)
}

func TestFavoriteUrlListCommand(t *testing.T) {

	output, _ := utils.ExecuteCommand(FavoriteUrlCmd, "list")

	err := json.Unmarshal([]byte(output), &testFavoriteUrlList)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestFavoriteUrlList, testFavoriteUrlList)
}
