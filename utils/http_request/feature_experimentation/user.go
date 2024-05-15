package feature_experimentation

import (
	"net/http"
	"net/url"

	models "github.com/flagship-io/abtasty-cli/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
)

type UserRequester struct {
	*common.ResourceRequest
}

func (u *UserRequester) HTTPListUsers() ([]models.User, error) {
	return common.HTTPGetAllPagesFE[models.User](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + u.AccountID + "/account_environments/" + u.AccountEnvironmentID + "/users")
}

func (u *UserRequester) HTTPBatchUpdateUsers(data string) ([]byte, error) {
	return common.HTTPRequest[models.User](http.MethodPut, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+u.AccountID+"/account_environments/"+u.AccountEnvironmentID+"/users", []byte(data))
}

func (u *UserRequester) HTTPDeleteUsers(email string) error {
	_, err := common.HTTPRequest[models.User](http.MethodDelete, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+u.AccountID+"/account_environments/"+u.AccountEnvironmentID+"/users?emails[]="+url.QueryEscape(email), nil)
	return err
}
