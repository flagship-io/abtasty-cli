package web_experimentation

import (
	"encoding/json"
	"net/http"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
)

type CampaignWERequester struct {
	*common.ResourceRequest
}

func (t *CampaignWERequester) HTTPListCampaign() ([]models.CampaignWE, error) {
	return common.HTTPGetAllPagesWE[models.CampaignWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + t.AccountID + "/tests")
}

func (t *CampaignWERequester) HTTPCreateCampaign(data string) ([]byte, error) {
	return common.HTTPRequest[models.CampaignWE](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/tests", []byte(data))
}

func (t *CampaignWERequester) HTTPEditCampaign(id, data string) ([]byte, error) {
	return common.HTTPRequest[models.CampaignWE](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/tests/"+id, []byte(data))
}

func (t *CampaignWERequester) HTTPGetCampaign(id string) (models.CampaignWE, error) {
	return common.HTTPGetItem[models.CampaignWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + t.AccountID + "/tests/" + id)
}

func (t *CampaignWERequester) HTTPDeleteCampaign(id string) error {
	_, err := common.HTTPRequest[models.CampaignWE](http.MethodDelete, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/tests/"+id, nil)

	return err
}

func (t *CampaignWERequester) HTTPSwitchStateCampaign(id, state string) error {
	var active bool
	if state == "active" {
		active = true
	}

	campaignSwitchRequest := models.CampaignState{
		Active: active,
	}

	campaignSwitchRequestJSON, err := json.Marshal(campaignSwitchRequest)
	if err != nil {
		return err
	}

	_, err = common.HTTPRequest[models.CampaignWE](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/tests/"+id, campaignSwitchRequestJSON)
	return err
}
