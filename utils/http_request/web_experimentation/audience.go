package web_experimentation

import (
	"os"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	mockfunction_we "github.com/flagship-io/abtasty-cli/utils/mock_function/web_experimentation"
	"github.com/jarcoal/httpmock"
)

func init() {
	if os.Getenv("ABT_ENV") == "MOCK" {
		httpmock.Activate()
		mockfunction_we.APIAudience()
	}
}

type AudienceRequester struct {
	*common.ResourceRequest
}

func (a *AudienceRequester) HTTPListAudience() ([]models.Audience, error) {
	return common.HTTPGetAllPagesWE[models.Audience](utils.GetWebExperimentationHost() + "/v1/accounts/" + a.AccountID + "/audiences?status=unarchive&")
}

func (a *AudienceRequester) HTTPGetAudience(id string) (models.Audience, error) {
	return common.HTTPGetItem[models.Audience](utils.GetWebExperimentationHost() + "/v1/accounts/" + a.AccountID + "/audiences/" + id)
}
