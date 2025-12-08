package web_experimentation

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/stretchr/testify/assert"
)

var campaignTargetingRequester = CampaignTargetingRequester{&common.ResourceRequest{AccountID: "account_id"}}

func TestHTTPCampaignTargeting(t *testing.T) {

	respBody, err := campaignTargetingRequester.HTTPGetCampaignTargeting(100000)

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "IS", respBody.UrlScopes[0].Condition)
	assert.Equal(t, "https://abtasty.com", respBody.UrlScopes[0].Value)

}
