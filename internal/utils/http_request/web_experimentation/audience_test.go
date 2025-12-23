package web_experimentation

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
	"github.com/stretchr/testify/assert"
)

var audienceRequester = AudienceRequester{&common.ResourceRequest{AccountID: "account_id"}}

func TestHTTPGetAudience(t *testing.T) {

	respBody, err := audienceRequester.HTTPGetAudience("testAudienceId")

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "testAudienceId", respBody.Id)
	assert.Equal(t, "testAudienceName", respBody.Name)
	assert.Equal(t, false, respBody.IsSegment)

}

func TestHTTPListAudience(t *testing.T) {

	respBody, err := audienceRequester.HTTPListAudiences()

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "testAudienceId1", respBody[1].Id)
	assert.Equal(t, "testAudienceName1", respBody[1].Name)
	assert.Equal(t, true, respBody[1].IsSegment)

}
