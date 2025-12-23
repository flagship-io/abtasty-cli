package web_experimentation

import (
	"fmt"
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

func (f *AudienceRequester) HTTPDeleteAudience(id string) (string, error) {
	_, err := common.HTTPRequest[models.Audience](http.MethodDelete, utils.GetWebExperimentationHost()+"/v1/accounts/"+f.AccountID+"/audiences/"+id, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Audience %s deleted", id), nil
}
