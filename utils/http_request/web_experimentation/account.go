package web_experimentation

import (
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
