package web_experimentation

import (
	"net/http"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
)

type AudienceRequester struct {
	*common.ResourceRequest
}

func (a *AudienceRequester) HTTPListAudience() ([]models.Audience, error) {
	return common.HTTPGetAllPagesWE[models.Audience](utils.GetWebExperimentationHost() + "/v1/accounts/" + a.AccountID + "/audiences?status=unarchive&")
}

func (a *AudienceRequester) HTTPGetAudience(id string) (models.Audience, error) {
	return common.HTTPGetItem[models.Audience](utils.GetWebExperimentationHost() + "/v1/accounts/" + a.AccountID + "/audiences/" + id)
}

func (a *AudienceRequester) HTTPCreateAudience(data []byte) ([]byte, error) {
	return common.HTTPRequest[models.Audience](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+a.AccountID+"/audiences", data)
}
