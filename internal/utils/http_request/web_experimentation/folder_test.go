package web_experimentation

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/stretchr/testify/assert"
)

var folderRequester = FolderRequester{&common.ResourceRequest{AccountID: "account_id"}}

func TestHTTPCreateFolder(t *testing.T) {
	data := "{\"name\":\"testFolderName\""

	respBody, err := folderRequester.HTTPCreateFolder([]byte(data))

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "{\"id\":900000,\"name\":\"testFolderName\"}", string(respBody))
}

func TestHTTPGetFolder(t *testing.T) {
	respBody, err := folderRequester.HTTPGetFolder(900000)

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, 900000, respBody.Id)
	assert.Equal(t, "testFolderName", respBody.Name)
}

func TestHTTPListFolder(t *testing.T) {
	respBody, err := folderRequester.HTTPListFolder()

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, 900001, respBody[1].Id)
	assert.Equal(t, "testFolderName1", respBody[1].Name)
}

func TestHTTPDeleteFolder(t *testing.T) {
	resp, err := folderRequester.HTTPDeleteFolder(900000)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
