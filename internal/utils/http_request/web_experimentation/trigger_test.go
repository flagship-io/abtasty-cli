package web_experimentation

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
	"github.com/stretchr/testify/assert"
)

var triggerRequester = TriggerRequester{&common.ResourceRequest{AccountID: "account_id"}}

func TestHTTPGetTrigger(t *testing.T) {

	respBody, err := triggerRequester.HTTPGetTrigger("trigger-id")

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "trigger-id", respBody.Id)
	assert.Equal(t, "testTriggerName", respBody.Name)
	assert.Equal(t, false, respBody.IsSegment)

}

func TestHTTPListTrigger(t *testing.T) {

	respBody, err := triggerRequester.HTTPListTrigger()

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "trigger-id-1", respBody[1].Id)
	assert.Equal(t, "testTriggerName1", respBody[1].Name)
	assert.Equal(t, false, respBody[1].IsSegment)

}
