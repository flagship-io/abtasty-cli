package web_experimentation

import (
	"net/http"

	models "github.com/flagship-io/flagship/models/web_experimentation"
	"github.com/flagship-io/flagship/utils"
	"github.com/flagship-io/flagship/utils/http_request/common"
)

type CampaignWERequester struct {
	*common.ResourceRequest
}

func (t *CampaignWERequester) HTTPListCampaign() ([]models.CampaignWE, error) {
	return common.HTTPGetAllPagesWE[models.CampaignWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + t.AccountID + "/tests")
}

func (t *CampaignWERequester) HTTPGetCampaign(id string) (models.CampaignWE, error) {
	return common.HTTPGetItem[models.CampaignWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + t.AccountID + "/tests/" + id)
}

func (t *CampaignWERequester) HTTPCreateCampaign(data string) ([]byte, error) {
	return common.HTTPRequest[models.CampaignWE](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/tests", []byte(data))
}

func (t *CampaignWERequester) HTTPEditCampaign(id, data string) ([]byte, error) {
	return common.HTTPRequest[models.CampaignWE](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/tests/"+id, []byte(data))
}

/* func HTTPSwitchCampaign(id, state string) error {
	campaignSwitchRequest := models.CampaignSwitchRequest{
		State: state,
	}

	campaignSwitchRequestJSON, err := json.Marshal(campaignSwitchRequest)
	if err != nil {
		return err
	}

	_, err = common.HTTPRequest(http.MethodPatch, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+viper.GetString("account_id")+"/account_environments/"+viper.GetString("account_environment_id")+"/campaigns/"+id+"/toggle", campaignSwitchRequestJSON)
	return err
} */

func (t *CampaignWERequester) HTTPDeleteCampaign(id string) error {
	_, err := common.HTTPRequest[models.CampaignWE](http.MethodDelete, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+t.AccountID+"/tests/"+id, nil)
	return err
}