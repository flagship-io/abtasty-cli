package web_experimentation

import (
	"testing"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/flagship-io/abtasty-cli/utils/mock_function/web_experimentation"
	"github.com/stretchr/testify/assert"
)

var modificationRequester = ModificationRequester{&common.ResourceRequest{AccountID: "account_id"}}

func TestHTTPCreateModification(t *testing.T) {
	data := models.ModificationCodeCreateStruct{Name: "modification", Type: "customScriptNew", VariationID: 110000, Value: "console.log(\"Hello World!\")"}

	respBody, err := modificationRequester.HTTPCreateModification(100000, data)

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, []byte("{\"id\":120003,\"name\":\"modification\",\"type\":\"customScriptNew\",\"value\":\"console.log(\\\"test modification\\\")\",\"variation_id\":110000,\"selector\":\"document.querySelector()\",\"engine\":\"\",\"updated_by\":{\"id\":0,\"email\":\"\"},\"updated_at\":{\"readable_date\":\"\",\"timestamp\":0,\"pattern\":\"\"}}"), respBody)
}

func TestHTTPEditModification(t *testing.T) {
	data := models.ModificationCodeEditStruct{Name: "modification", Type: "customScriptNew", Value: "console.log(\"Hello World!\")"}

	respBody, err := modificationRequester.HTTPEditModification(100000, 120001, data)

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, []byte("{\"id\":120001,\"name\":\"modification\",\"type\":\"customScriptNew\",\"value\":\"console.log(\\\"test modification\\\")\",\"variation_id\":110000,\"selector\":\"\",\"engine\":\"\",\"updated_by\":{\"id\":0,\"email\":\"\"},\"updated_at\":{\"readable_date\":\"\",\"timestamp\":0,\"pattern\":\"\"}}"), respBody)
}

func TestHTTPListModification(t *testing.T) {

	respBody, err := modificationRequester.HTTPListModification(100000)

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, web_experimentation.TestModifications.Data.Modifications, respBody)

}

func TestHTTPGetModification(t *testing.T) {

	respBody, err := modificationRequester.HTTPGetModification(100000, 120003)

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, []models.Modification{web_experimentation.TestModification.Data.Modifications[0]}, respBody)

}

func TestHTTPDeleteModification(t *testing.T) {

	resp, err := modificationRequester.HTTPDeleteModification(100000, 120003)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
