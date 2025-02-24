package feature_experimentation

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/stretchr/testify/assert"
)

var campaignRequester = CampaignFERequester{&common.ResourceRequest{AccountID: "account_id", AccountEnvironmentID: "account_environment_id"}}

func TestHTTPGetCampaign(t *testing.T) {
	respBody, err := campaignRequester.HTTPGetCampaign("testCampaignID")

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "testCampaignID", respBody.Id)
	assert.Equal(t, "testCampaignName", respBody.Name)
	assert.Equal(t, "testProjectID", respBody.ProjectId)
	assert.Equal(t, "testCampaignDescription", respBody.Description)
	assert.Equal(t, "toggle", respBody.Type)
}

func TestHTTPListCampaign(t *testing.T) {
	respBody, err := campaignRequester.HTTPListCampaign()

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "testCampaignID", respBody[0].Id)
	assert.Equal(t, "testCampaignName", respBody[0].Name)
	assert.Equal(t, "testProjectID", respBody[0].ProjectId)
	assert.Equal(t, "testCampaignDescription", respBody[0].Description)
	assert.Equal(t, "toggle", respBody[0].Type)

	assert.Equal(t, "testCampaignID1", respBody[1].Id)
	assert.Equal(t, "testCampaignName1", respBody[1].Name)
	assert.Equal(t, "testProjectID1", respBody[1].ProjectId)
	assert.Equal(t, "testCampaignDescription1", respBody[1].Description)
	assert.Equal(t, "toggle", respBody[1].Type)
}

func TestHTTPCreateCampaign(t *testing.T) {

	dataCampaign := "{\"project_id\":\"testProjectID\",\"name\":\"testCampaignName\",\"description\":\"testCampaignDescription\",\"type\":\"toggle\",\"variation_groups\":[{\"name\":\"variationGroupName\",\"variations\":[{\"name\":\"My variation 1\",\"allocation\":50,\"reference\":true,\"modifications\":{\"value\":{\"color\":\"blue\"}}},{\"name\":\"My variation 2\",\"allocation\":50,\"reference\":false,\"modifications\":{\"value\":{\"color\":\"red\"}}}],\"targeting\":{\"targeting_groups\":[{\"targetings\":[{\"operator\":\"CONTAINS\",\"key\":\"isVIP\",\"value\":\"true\"}]}]}}],\"scheduler\":{\"start_date\":\"2022-02-01 10:00:00\",\"stop_date\":\"2022-02-02 08:00:00\",\"timezone\":\"Europe/Paris\"}}"
	respBody, err := campaignRequester.HTTPCreateCampaign(dataCampaign)

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, []byte("{\"id\":\"testCampaignID\",\"project_id\":\"testProjectID\",\"name\":\"testCampaignName\",\"description\":\"testCampaignDescription\",\"type\":\"toggle\",\"variation_groups\":[{\"name\":\"variationGroupName\",\"variations\":[{\"name\":\"My variation 1\",\"reference\":true,\"allocation\":50,\"modifications\":{\"type\":\"string\",\"value\":{\"color\":\"blue\"}}},{\"name\":\"My variation 2\",\"reference\":false,\"allocation\":50,\"modifications\":{\"type\":\"string\",\"value\":{\"color\":\"red\"}}}],\"targeting\":{\"targeting_groups\":[{\"targetings\":[{\"key\":\"isVIP\",\"operator\":\"CONTAINS\",\"value\":true}]}]}}],\"scheduler\":{\"start_date\":\"2022-02-01 10:00:00\",\"stop_date\":\"2022-02-02 08:00:00\",\"timezone\":\"Europe/Paris\"}}"), respBody)
}

func TestHTTPEditCampaign(t *testing.T) {

	dataCampaign := "{\"project_id\":\"testProjectID1\",\"name\":\"testCampaignName1\",\"description\":\"testCampaignDescription1\",\"type\":\"toggle\",\"variation_groups\":[{\"name\":\"variationGroupName\",\"variations\":[{\"name\":\"My variation 1\",\"allocation\":50,\"reference\":true,\"modifications\":{\"value\":{\"color\":\"blue\"}}},{\"name\":\"My variation 2\",\"allocation\":50,\"reference\":false,\"modifications\":{\"value\":{\"color\":\"red\"}}}],\"targeting\":{\"targeting_groups\":[{\"targetings\":[{\"operator\":\"CONTAINS\",\"key\":\"isVIP\",\"value\":\"true\"}]}]}}],\"scheduler\":{\"start_date\":\"2022-02-01 10:00:00\",\"stop_date\":\"2022-02-02 08:00:00\",\"timezone\":\"Europe/Paris\"}}"

	respBody, err := campaignRequester.HTTPEditCampaign("testCampaignID", dataCampaign)

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, []byte("{\"id\":\"testCampaignID\",\"project_id\":\"testProjectID1\",\"name\":\"testCampaignName1\",\"description\":\"testCampaignDescription1\",\"type\":\"toggle\",\"variation_groups\":[{\"name\":\"variationGroupName\",\"variations\":[{\"name\":\"My variation 1\",\"reference\":true,\"allocation\":50,\"modifications\":{\"type\":\"string\",\"value\":{\"color\":\"blue\"}}},{\"name\":\"My variation 2\",\"reference\":false,\"allocation\":50,\"modifications\":{\"type\":\"string\",\"value\":{\"color\":\"red\"}}}],\"targeting\":{\"targeting_groups\":[{\"targetings\":[{\"key\":\"isVIP\",\"operator\":\"CONTAINS\",\"value\":true}]}]}}],\"scheduler\":{\"start_date\":\"2022-02-01 10:00:00\",\"stop_date\":\"2022-02-02 08:00:00\",\"timezone\":\"Europe/Paris\"}}"), respBody)
}

func TestHTTPDeleteCampaign(t *testing.T) {

	err := campaignRequester.HTTPDeleteCampaign("testCampaignID")

	assert.Nil(t, err)
}

func TestHTTPSwitchCampaign(t *testing.T) {

	err := campaignRequester.HTTPSwitchStateCampaign("testCampaignID", "active")

	assert.Nil(t, err)
}
