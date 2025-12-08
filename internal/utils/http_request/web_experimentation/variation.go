package web_experimentation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
)

type VariationWERequester struct {
	*common.ResourceRequest
}

func (v *VariationWERequester) HTTPCreateVariation(campaignID int, variationData models.VariationWE) ([]byte, error) {
	data, err := json.Marshal(variationData)
	if err != nil {
		return nil, err
	}

	return common.HTTPRequest[models.VariationWE](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+v.AccountID+"/tests/"+strconv.Itoa(campaignID)+"/variations", data)
}

func (v *VariationWERequester) HTTPCreateVariationDataRaw(campaignID int, data []byte) ([]byte, error) {
	return common.HTTPRequest[models.VariationWE](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+v.AccountID+"/tests/"+strconv.Itoa(campaignID)+"/variations", data)
}

func (v *VariationWERequester) HTTPEditVariation(campaignID, variationID int, data []byte) ([]byte, error) {
	return common.HTTPRequest[models.VariationWE](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+v.AccountID+"/tests/"+strconv.Itoa(campaignID)+"/variations/"+strconv.Itoa(variationID), data)
}

func (v *VariationWERequester) HTTPGetVariation(campaignID, variationID int) (models.VariationWE, error) {
	return common.HTTPGetItem[models.VariationWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + v.AccountID + "/tests/" + strconv.Itoa(campaignID) + "/variations/" + strconv.Itoa(variationID))
}

func (v *VariationWERequester) HTTPDeleteVariation(campaignID, variationID int) (string, error) {
	_, err := common.HTTPRequest[models.VariationWE](http.MethodDelete, utils.GetWebExperimentationHost()+"/v1/accounts/"+v.AccountID+"/tests/"+strconv.Itoa(campaignID)+"/variations/"+strconv.Itoa(variationID), nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Variation %d deleted", variationID), nil
}
