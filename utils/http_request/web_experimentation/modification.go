package web_experimentation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
)

type ModificationRequester struct {
	*common.ResourceRequest
}

func (m *ModificationRequester) HTTPListModification(campaignID int) ([]models.Modification, error) {
	resp, err := common.HTTPGetItem[models.ModificationDataWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + m.AccountID + "/tests/" + strconv.Itoa(campaignID) + "/modifications")
	return resp.Data.Modifications, err
}

func (m *ModificationRequester) HTTPGetModification(campaignID int, id int) ([]models.Modification, error) {
	resp, err := common.HTTPGetItem[models.ModificationDataWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + m.AccountID + "/tests/" + strconv.Itoa(campaignID) + "/modifications?ids=" + strconv.Itoa(id))
	return resp.Data.Modifications, err
}

func (m *ModificationRequester) HTTPEditModification(campaignID int, id int, modificationData web_experimentation.ModificationCodeEditStruct) ([]byte, error) {
	data, err := json.Marshal(modificationData)
	if err != nil {
		return nil, err
	}

	return common.HTTPRequest[models.ModificationDataWE](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+m.AccountID+"/tests/"+strconv.Itoa(campaignID)+"/modifications/"+strconv.Itoa(id), data)
}

func (m *ModificationRequester) HTTPCreateModification(campaignID int, modificationData web_experimentation.ModificationCodeCreateStruct) ([]byte, error) {
	data, err := json.Marshal(modificationData)
	if err != nil {
		return nil, err
	}

	return common.HTTPRequest[models.ModificationDataWE](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+m.AccountID+"/tests/"+strconv.Itoa(campaignID)+"/modifications", data)
}

func (m *ModificationRequester) HTTPEditModificationDataRaw(campaignID int, id int, data string) ([]byte, error) {
	return common.HTTPRequest[models.ModificationDataWE](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+m.AccountID+"/tests/"+strconv.Itoa(campaignID)+"/modifications/"+strconv.Itoa(id), []byte(data))
}

func (m *ModificationRequester) HTTPCreateModificationDataRaw(campaignID int, data []byte) ([]byte, error) {
	return common.HTTPRequest[models.ModificationDataWE](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+m.AccountID+"/tests/"+strconv.Itoa(campaignID)+"/modifications", data)
}

func (m *ModificationRequester) HTTPDeleteModification(campaignID int, id int) error {
	_, err := common.HTTPRequest[models.ModificationDataWE](http.MethodDelete, utils.GetWebExperimentationHost()+"/v1/accounts/"+m.AccountID+"/tests/"+strconv.Itoa(campaignID)+"/modifications/"+strconv.Itoa(id)+"?input_type=modification", nil)
	return err
}
