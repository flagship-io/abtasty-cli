package mockfunction

import (
	"net/http"
	"strconv"

	"github.com/flagship-io/flagship/models"
	"github.com/flagship-io/flagship/utils"
	"github.com/flagship-io/flagship/utils/config"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
)

var TestToken = models.Token{
	ClientID:  "client_id",
	AccountID: "account_id",
	ExpiresIn: 0,
	Scope:     "*",
}

func APIToken() {
	config.SetViperMock()

	token := "token"
	tokenExpiration := 0

	testAuthenticationResponse := models.AuthenticationResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	httpmock.RegisterResponder("GET", utils.GetHostAuth()+"/token?access_token="+token,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestToken)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("POST", utils.GetHostAuth()+"/"+viper.GetString("account_id")+"/token?expires_in="+strconv.Itoa(tokenExpiration),
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, testAuthenticationResponse)
			return resp, nil
		},
	)
}
