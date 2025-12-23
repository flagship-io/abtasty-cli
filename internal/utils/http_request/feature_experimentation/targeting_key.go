package feature_experimentation

import (
	"net/http"

	models "github.com/flagship-io/abtasty-cli/internal/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
)

type TargetingKeyRequester struct {
	*common.ResourceRequest
}

func (t *TargetingKeyRequester) HTTPListTargetingKey() ([]models.TargetingKey, error) {
	return common.HTTPGetAllPagesFE[models.TargetingKey](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + t.AccountID + "/targeting_keys")
}

func (t *TargetingKeyRequester) HTTPGetTargetingKey(id string) (models.TargetingKey, error) {
	return common.HTTPGetItem[models.TargetingKey](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + t.AccountID + "/targeting_keys/" + id)
}

func (t *TargetingKeyRequester) HTTPCreateTargetingKey(data []byte) ([]byte, error) {
	return common.HTTPRequest[models.TargetingKey](http.MethodPost, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+t.AccountID+"/targeting_keys", data)
}

func (t *TargetingKeyRequester) HTTPEditTargetingKey(id string, data []byte) ([]byte, error) {
	return common.HTTPRequest[models.TargetingKey](http.MethodPatch, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+t.AccountID+"/targeting_keys/"+id, data)
}

func (t *TargetingKeyRequester) HTTPDeleteTargetingKey(id string) error {
	_, err := common.HTTPRequest[models.TargetingKey](http.MethodDelete, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+t.AccountID+"/targeting_keys/"+id, nil)
	return err
}
