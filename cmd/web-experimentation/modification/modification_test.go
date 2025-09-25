package modification

import (
	"encoding/json"
	"testing"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
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

var testModification models.ModificationResourceLoader
var testModificationList []models.ModificationResourceLoader

func TestModificationCommand(t *testing.T) {
	successOutput, _ := utils.ExecuteCommand(ModificationCmd)
	assert.Contains(t, successOutput, "Manage your modifications")
}

func TestModificationHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(ModificationCmd, "--help")
	assert.Contains(t, output, "Manage your modifications")
}

func TestModificationGetCommand(t *testing.T) {
	failOutput, err := utils.ExecuteCommand(ModificationCmd, "get", "--campaign-id=100000")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(ModificationCmd, "get", "--id=120003", "--campaign-id=100000")
	err = json.Unmarshal([]byte(successOutput), &testModification)

	assert.Nil(t, err)

	modification := mockfunction_we.TestModification.Data.Modifications[0]
	_modif := web_experimentation.ModificationResourceLoader{Id: modification.Id, Name: modification.Name, Type: getTypeFromModificationAPI(modification.Type), Selector: modification.Selector, Code: modification.Value, VariationID: modification.VariationID, CampaignID: CampaignID}

	assert.Equal(t, _modif, testModification)
}

func TestModificationListCommand(t *testing.T) {
	var modificationsRL []web_experimentation.ModificationResourceLoader
	output, err := utils.ExecuteCommand(ModificationCmd, "list", "--campaign-id=100000")
	err = json.Unmarshal([]byte(output), &testModificationList)

	assert.Nil(t, err)

	for _, modification := range []models.Modification{mockfunction_we.TestModificationsJS, mockfunction_we.TestModificationsCSS} {
		modificationsRL = append(modificationsRL, web_experimentation.ModificationResourceLoader{Id: modification.Id, Name: modification.Name, Type: getTypeFromModificationAPI(modification.Type), Selector: modification.Selector, Code: modification.Value, VariationID: modification.VariationID, CampaignID: CampaignID})
	}

	assert.Equal(t, modificationsRL, testModificationList)
}

func TestModificationDeleteCommand(t *testing.T) {
	failOutput, _ := utils.ExecuteCommand(ModificationCmd, "delete", "--campaign-id=100000")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(ModificationCmd, "delete", "--campaign-id=100000", "--id=120003")
	assert.Equal(t, "Modification 120003 deleted\n", successOutput)
}
