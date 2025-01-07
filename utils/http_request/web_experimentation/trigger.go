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
		mockfunction_we.APITrigger()
	}
}

type TriggerRequester struct {
	*common.ResourceRequest
}

func (t *TriggerRequester) HTTPListTrigger() ([]models.Audience, error) {
	return common.HTTPGetAllPagesWE[models.Audience](utils.GetWebExperimentationHost() + "/v1/accounts/" + t.AccountID + "/audiences?type=trigger&status=unarchive&")
}

func (t *TriggerRequester) HTTPGetTrigger(id string) (models.Audience, error) {
	return common.HTTPGetItem[models.Audience](utils.GetWebExperimentationHost() + "/v1/accounts/" + t.AccountID + "/audiences/" + id)
}
