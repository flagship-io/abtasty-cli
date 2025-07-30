package web_experimentation

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"

	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
)

type CampaignWERequester struct {
	*common.ResourceRequest
}

func (t *CampaignWERequester) HTTPListCampaign() ([]models.CampaignWE, error) {
	return common.HTTPGetAllPagesWE[models.CampaignWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + t.AccountID + "/tests?state=play,pause,~unarchive&")
}

func (t *CampaignWERequester) HTTPCreateCampaign(data []byte) ([]byte, error) {
	return common.HTTPRequest[models.CampaignWECommon](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/tests", data)
}

func (t *CampaignWERequester) HTTPEditCampaign(id int, data []byte) ([]byte, error) {
	return common.HTTPRequest[models.CampaignWE](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/tests/"+strconv.Itoa(id), data)
}

func (t *CampaignWERequester) HTTPGetCampaign(id int) (models.CampaignWE, error) {
	return common.HTTPGetItem[models.CampaignWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + t.AccountID + "/tests/" + strconv.Itoa(id))
}

func (t *CampaignWERequester) HTTPDeleteCampaign(id int) error {
	_, err := common.HTTPRequest[models.CampaignWE](http.MethodDelete, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/tests/"+strconv.Itoa(id), nil)

	return err
}

func (t *CampaignWERequester) HTTPSwitchStateCampaign(id int, state string) error {
	status := "play"
	if state == "paused" {
		status = "pause"
	}

	campaignSwitchRequest := models.CampaignState{
		Status: status,
	}

	campaignSwitchRequestJSON, err := json.Marshal(campaignSwitchRequest)
	if err != nil {
		return err
	}

	_, err = common.HTTPRequest[models.CampaignWE](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/tests/"+strconv.Itoa(id), campaignSwitchRequestJSON)
	return err
}
