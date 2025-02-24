package web_experimentation

import (
	"net/http"

	"github.com/flagship-io/abtasty-cli/models"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/jarcoal/httpmock"
)

var TestToken = models.Token{
	ClientID:  "client_id",
	AccountID: "account_id",
	ExpiresIn: 10,
	Scope:     "scope",
}

func APIToken() {

	token := "access_token"

	testAuthenticationResponse := models.TokenResponse{
		AccessToken: "testAccessToken",
	}

	httpmock.RegisterResponder("GET", utils.GetHostWebExperimentationAuth()+"/v1/token?access_token="+token,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestToken)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("POST", utils.GetHostWebExperimentationAuth()+"/v1/token",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, testAuthenticationResponse)
			return resp, nil
		},
	)
}
