package feature_experimentation

import (
	"encoding/json"
	"net/http"

	models "github.com/flagship-io/flagship/models/feature_experimentation"
	"github.com/flagship-io/flagship/utils"
	"github.com/flagship-io/flagship/utils/http_request/common"
)

type CampaignRequester struct {
	*common.ResourceRequest
}

func (c *CampaignRequester) HTTPListCampaign() ([]models.Campaign, error) {
	return common.HTTPGetAllPagesFE[models.Campaign](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + c.AccountID + "/account_environments/" + c.AccountEnvironmentID + "/campaigns")
}

func (c *CampaignRequester) HTTPGetCampaign(id string) (models.Campaign, error) {
	return common.HTTPGetItem[models.Campaign](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + c.AccountID + "/account_environments/" + c.AccountEnvironmentID + "/campaigns/" + id)
}

func (c *CampaignRequester) HTTPCreateCampaign(data string) ([]byte, error) {
	return common.HTTPRequest[models.Campaign](http.MethodPost, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+c.AccountID+"/account_environments/"+c.AccountEnvironmentID+"/campaigns", []byte(data))
}

func (c *CampaignRequester) HTTPEditCampaign(id, data string) ([]byte, error) {
	return common.HTTPRequest[models.Campaign](http.MethodPatch, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+c.AccountID+"/account_environments/"+c.AccountEnvironmentID+"/campaigns/"+id, []byte(data))
}

func (c *CampaignRequester) HTTPSwitchCampaign(id, state string) error {
	campaignSwitchRequest := models.CampaignSwitchRequest{
		State: state,
	}

	campaignSwitchRequestJSON, err := json.Marshal(campaignSwitchRequest)
	if err != nil {
		return err
	}

	_, err = common.HTTPRequest[models.Campaign](http.MethodPatch, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+c.AccountID+"/account_environments/"+c.AccountEnvironmentID+"/campaigns/"+id+"/toggle", campaignSwitchRequestJSON)
	return err
}

func (c *CampaignRequester) HTTPDeleteCampaign(id string) error {
	_, err := common.HTTPRequest[models.Campaign](http.MethodDelete, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+c.AccountID+"/account_environments/"+c.AccountEnvironmentID+"/campaigns/"+id, nil)
	return err
}
