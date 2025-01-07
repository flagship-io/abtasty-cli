package feature_experimentation

import (
	"net/http"
	"os"

	models "github.com/flagship-io/abtasty-cli/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	mockfunction_fe "github.com/flagship-io/abtasty-cli/utils/mock_function/feature_experimentation"
	"github.com/jarcoal/httpmock"
)

func init() {
	if os.Getenv("ABT_ENV") == "MOCK" {
		httpmock.Activate()
		mockfunction_fe.APITargetingKey()
	}
}

type TargetingKeyRequester struct {
	*common.ResourceRequest
}

func (t *TargetingKeyRequester) HTTPListTargetingKey() ([]models.TargetingKey, error) {
	return common.HTTPGetAllPagesFE[models.TargetingKey](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + t.AccountID + "/targeting_keys")
}

func (t *TargetingKeyRequester) HTTPGetTargetingKey(id string) (models.TargetingKey, error) {
	return common.HTTPGetItem[models.TargetingKey](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + t.AccountID + "/targeting_keys/" + id)
}

func (t *TargetingKeyRequester) HTTPCreateTargetingKey(data string) ([]byte, error) {
	return common.HTTPRequest[models.TargetingKey](http.MethodPost, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+t.AccountID+"/targeting_keys", []byte(data))
}

func (t *TargetingKeyRequester) HTTPEditTargetingKey(id, data string) ([]byte, error) {
	return common.HTTPRequest[models.TargetingKey](http.MethodPatch, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+t.AccountID+"/targeting_keys/"+id, []byte(data))
}

func (t *TargetingKeyRequester) HTTPDeleteTargetingKey(id string) error {
	_, err := common.HTTPRequest[models.TargetingKey](http.MethodDelete, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+t.AccountID+"/targeting_keys/"+id, nil)
	return err
}
