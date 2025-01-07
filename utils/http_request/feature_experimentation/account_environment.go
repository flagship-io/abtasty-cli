package feature_experimentation

import (
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
		mockfunction_fe.APIAccountEnvironment()
	}
}

type AccountEnvironmentFERequester struct {
	*common.ResourceRequest
}

func (a *AccountEnvironmentFERequester) HTTPListAccountEnvironment(accountID string) ([]models.AccountEnvironmentFE, error) {
	if accountID == "" {
		accountID = a.AccountID
	}

	return common.HTTPGetAllPagesFE[models.AccountEnvironmentFE](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + accountID + "/account_environments")
}

func (a *AccountEnvironmentFERequester) HTTPGetAccountEnvironment(id string) (models.AccountEnvironmentFE, error) {
	return common.HTTPGetItem[models.AccountEnvironmentFE](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + a.AccountID + "/account_environments/" + id)
}
