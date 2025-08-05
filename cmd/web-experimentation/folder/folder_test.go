/*
Copyright © 2022 Flagship Team flagship@abtasty.com
*/
package folder

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
	mockfunction_we.APIFolder()

	m.Run()
}

var testFolder models.Folder
var testFolderList []models.Folder

func TestFolderCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(FolderCmd)
	assert.Contains(t, output, "Manage your folders")
}

func TestFolderHelpCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(FolderCmd, "--help")
	assert.Contains(t, output, "Manage your folders")
}

func TestFolderGetCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(FolderCmd, "get")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(FolderCmd, "get", "--id=900000")

	err := json.Unmarshal([]byte(successOutput), &testFolder)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestFolder, testFolder)
}

func TestFolderListCommand(t *testing.T) {

	output, _ := utils.ExecuteCommand(FolderCmd, "list")

	err := json.Unmarshal([]byte(output), &testFolderList)

	assert.Nil(t, err)

	assert.Equal(t, mockfunction_we.TestFolderList, testFolderList)
}

func TestFolderDeleteCommand(t *testing.T) {

	failOutput, _ := utils.ExecuteCommand(FolderCmd, "delete")
	assert.Contains(t, failOutput, "Error: required flag(s) \"id\" not set")

	successOutput, _ := utils.ExecuteCommand(FolderCmd, "delete", "--id=900000")
	assert.Equal(t, "Folder 900000 deleted\n", successOutput)
}
