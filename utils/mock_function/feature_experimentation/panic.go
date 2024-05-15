package feature_experimentation

import (
	"net/http"

	"github.com/flagship-io/abtasty-cli/utils"
	mockfunction "github.com/flagship-io/abtasty-cli/utils/mock_function"
	"github.com/jarcoal/httpmock"
)

func APIPanic() {

	httpmock.RegisterResponder("PATCH", utils.GetFeatureExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/account_environments/"+mockfunction.Auth.AccountEnvironmentID+"/panic",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, "")
			return resp, nil
		},
	)
}
