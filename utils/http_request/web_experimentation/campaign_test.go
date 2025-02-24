package web_experimentation

import (
	"strconv"
	"testing"

	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/stretchr/testify/assert"
)

var campaignRequester = CampaignWERequester{&common.ResourceRequest{AccountID: "account_id"}}

func TestHTTPCreateCampaign(t *testing.T) {

	data := "{\"name\":\"testCampaignName\",\"type\":\"ab\",\"url\":\"https://abtasty.com\",\"description\":\"testCampaignDescription\",\"global_code\":\"console.log(\"Hello World!\")\"}"

	respBody, err := campaignRequester.HTTPCreateCampaign(data)

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "{\"id\":100000,\"name\":\"testCampaignName\",\"url\":\"https://abtasty.com\",\"description\":\"testCampaignDescription\",\"type\":\"ab\",\"sub_type\":\"\",\"master\":null,\"state\":\"\",\"traffic\":null,\"variations\":null,\"sub_tests\":null,\"created_at\":{\"readable_date\":\"\",\"timestamp\":0,\"pattern\":\"\"},\"labels\":null,\"last_play\":{\"readable_date\":\"\",\"timestamp\":0,\"pattern\":\"\"},\"last_pause\":{\"readable_date\":\"\",\"timestamp\":0,\"pattern\":\"\"},\"global_code\":\"console.log(\\\"Hello World!\\\")\",\"source_code\":\"\",\"audiences\":null,\"selector_scopes\":null,\"code_scopes\":null,\"favorite_url_scopes\":null,\"url_scopes\":[{\"condition\":40,\"include\":true,\"value\":\"https://abtasty.com\"},{\"condition\":40,\"include\":false,\"value\":\"https://abtasty.com\"}],\"mutation_observer\":false,\"display_frequency\":{\"type\":\"\",\"unit\":\"\",\"value\":0},\"report\":{\"token\":\"\",\"comment\":\"\"}}", string(respBody))
}

func TestHTTPEditCampaign(t *testing.T) {

	data := "{\"name\":\"testCampaignName\",\"type\":\"ab\",\"url\":\"https://abtasty.com\",\"description\":\"testCampaignDescription\",\"global_code\":\"console.log(\"Hello World!\")\"}"

	respBody, err := campaignRequester.HTTPEditCampaign(strconv.Itoa(100000), data)

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "{\"id\":100000,\"name\":\"testCampaignName\",\"url\":\"https://abtasty.com\",\"description\":\"testCampaignDescription\",\"type\":\"ab\",\"sub_type\":\"\",\"master\":null,\"state\":\"\",\"traffic\":null,\"variations\":null,\"sub_tests\":null,\"created_at\":{\"readable_date\":\"\",\"timestamp\":0,\"pattern\":\"\"},\"labels\":null,\"last_play\":{\"readable_date\":\"\",\"timestamp\":0,\"pattern\":\"\"},\"last_pause\":{\"readable_date\":\"\",\"timestamp\":0,\"pattern\":\"\"},\"global_code\":\"console.log(\\\"Hello World!\\\")\",\"source_code\":\"\",\"audiences\":null,\"selector_scopes\":null,\"code_scopes\":null,\"favorite_url_scopes\":null,\"url_scopes\":[{\"condition\":40,\"include\":true,\"value\":\"https://abtasty.com\"},{\"condition\":40,\"include\":false,\"value\":\"https://abtasty.com\"}],\"mutation_observer\":false,\"display_frequency\":{\"type\":\"\",\"unit\":\"\",\"value\":0},\"report\":{\"token\":\"\",\"comment\":\"\"}}", string(respBody))
}

func TestHTTPGetCampaign(t *testing.T) {

	respBody, err := campaignRequester.HTTPGetCampaign("100000")

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, 100000, respBody.Id)
	assert.Equal(t, "testCampaignName", respBody.Name)
	assert.Equal(t, "console.log(\"Hello World!\")", respBody.GlobalCodeCampaign)
	assert.Equal(t, "testCampaignDescription", respBody.Description)
	assert.Equal(t, "ab", respBody.Type)

}

func TestHTTPListCampaign(t *testing.T) {

	respBody, err := campaignRequester.HTTPListCampaign()

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, 100000, respBody[0].Id)
	assert.Equal(t, "testCampaignName", respBody[0].Name)
	assert.Equal(t, "console.log(\"Hello World!\")", respBody[0].GlobalCodeCampaign)
	assert.Equal(t, "testCampaignDescription", respBody[0].Description)
	assert.Equal(t, "ab", respBody[0].Type)

	assert.Equal(t, 100001, respBody[1].Id)
	assert.Equal(t, "testCampaignName1", respBody[1].Name)
	assert.Equal(t, "console.log(\"Hello Earth!\")", respBody[1].GlobalCodeCampaign)
	assert.Equal(t, "testCampaignDescription1", respBody[1].Description)
	assert.Equal(t, "ab", respBody[1].Type)

}

func TestHTTPDeleteCampaign(t *testing.T) {

	err := campaignRequester.HTTPDeleteCampaign("100000")

	assert.Nil(t, err)
}

func TestHTTPSwitchCampaign(t *testing.T) {

	err := campaignRequester.HTTPSwitchStateCampaign("100000", "active")

	assert.Nil(t, err)
}
