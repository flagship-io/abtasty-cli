package web_experimentation

import (
	"net/http"

	models_ "github.com/flagship-io/abtasty-cli/models"
	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/jarcoal/httpmock"
)

var TestAccount = models_.AccountJSON{
	CurrentUsedCredential: "test_auth",
	AccountID:             "account_id",
	AccountEnvironmentID:  "account_environment_id",
}

var TestGlobalCode = models.GlobalCode_{
	OnDomReady: true,
	Value:      "console.log(\"test\")",
}

var accountID = "account_id"

var TestAccountGlobalCode = models.AccountWE{
	Id:         100000,
	Name:       "account_name",
	Identifier: "account_identifier",
	Role:       "account_role",
	GlobalCode: TestGlobalCode,
}

func APIAccount() {

	resp := utils.HTTPListResponseWE[models.AccountWE]{
		Data: []models.AccountWE{TestAccountGlobalCode},
		Pagination: utils.Pagination{
			Total:      1,
			Pages:      2,
			MaxPerPage: 10,
			Page:       1,
		},
	}

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+accountID,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestAccountGlobalCode)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, resp)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("PATCH", utils.GetWebExperimentationHost()+"/v1/accounts/"+accountID,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestAccountGlobalCode)
			return resp, nil
		},
	)
}
