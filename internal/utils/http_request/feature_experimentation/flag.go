package feature_experimentation

import (
	"net/http"

	models "github.com/flagship-io/abtasty-cli/internal/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
)

type FlagRequester struct {
	*common.ResourceRequest
}

func (f *FlagRequester) HTTPListFlag() ([]models.Flag, error) {
	return common.HTTPGetAllPagesFE[models.Flag](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + f.AccountID + "/flags")
}

func (f *FlagRequester) HTTPGetFlag(id string) (models.Flag, error) {
	return common.HTTPGetItem[models.Flag](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + f.AccountID + "/flags/" + id)
}

func (f *FlagRequester) HTTPCreateFlag(data []byte) ([]byte, error) {
	return common.HTTPRequest[models.Flag](http.MethodPost, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+f.AccountID+"/flags", data)
}

func (f *FlagRequester) HTTPEditFlag(id string, data []byte) ([]byte, error) {
	return common.HTTPRequest[models.Flag](http.MethodPatch, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+f.AccountID+"/flags/"+id, data)
}

func (f *FlagRequester) HTTPDeleteFlag(id string) error {
	_, err := common.HTTPRequest[models.Flag](http.MethodDelete, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+f.AccountID+"/flags/"+id, nil)
	return err
}
