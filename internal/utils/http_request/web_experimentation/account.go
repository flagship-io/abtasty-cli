package web_experimentation

import (
	"net/http"

	"github.com/flagship-io/abtasty-cli/models"
	models_ "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
)

type AccountWERequester struct {
	*common.ResourceRequest
}

func (a *AccountWERequester) HTTPListAccount() ([]models_.AccountWE, error) {
	return common.HTTPGetAllPagesWE[models_.AccountWE](utils.GetWebExperimentationHost() + "/v1/accounts?")
}

func HTTPUserMe() (models.UserMe, error) {
	return common.HTTPGetItem[models.UserMe](utils.GetWebExperimentationHost() + "/v1/users/me")
}

func (a *AccountWERequester) HTTPRebuildTag() error {
	_, err := common.HTTPRequest[interface{}](http.MethodPatch, utils.GetWebExperimentationBackEndHost()+"/accounts/"+a.AccountID+"/tag-rebuild", nil)
	return err
}
