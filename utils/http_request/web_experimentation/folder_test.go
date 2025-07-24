package web_experimentation

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/stretchr/testify/assert"
)

var folderRequester = FolderRequester{&common.ResourceRequest{AccountID: "account_id"}}

func TestHTTPGetFolder(t *testing.T) {

	respBody, err := folderRequester.HTTPGetFolder(100000)

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, 100000, respBody.Id)
	assert.Equal(t, "testFolderName", respBody.Name)
}

func TestHTTPListFolder(t *testing.T) {

	respBody, err := folderRequester.HTTPListFolder()

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, 100001, respBody[1].Id)
	assert.Equal(t, "testFolderName1", respBody[1].Name)
}
