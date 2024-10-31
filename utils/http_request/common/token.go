package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/flagship-io/abtasty-cli/models"
	"github.com/flagship-io/abtasty-cli/utils"
)

func openLink(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported operating system")
	}
	err := cmd.Run()
	return err
}

func handleCallback(w http.ResponseWriter, r *http.Request, codeChan chan<- string) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No token found in URL", http.StatusBadRequest)
		os.Exit(0)
		return
	}

	codeChan <- code

	http.Redirect(w, r, utils.GetWebExperimentationBrowserAuthSuccess(), http.StatusSeeOther)

	go func() {
		time.Sleep(5 * time.Second)
		close(codeChan)
	}()
}

func HTTPRefreshTokenFE(client_id, refresh_token string) (models.TokenResponse, error) {
	var authenticationResponse models.TokenResponse
	authRequest := models.RefreshTokenRequestFE{
		ClientID:     client_id,
		GrantType:    "refresh_token",
		RefreshToken: refresh_token,
	}
	authRequestJSON, err := json.Marshal(authRequest)
	if err != nil {
		return models.TokenResponse{}, err
	}

	respBody, err := HTTPRequest[models.Token](http.MethodPost, utils.GetHostFeatureExperimentationAuth()+"/"+cred.AccountID+"/token", authRequestJSON)
	if err != nil {
		return models.TokenResponse{}, err
	}

	err = json.Unmarshal(respBody, &authenticationResponse)
	if err != nil {
		return models.TokenResponse{}, err
	}

	return authenticationResponse, err
}

func HTTPRefreshTokenWE(cred RequestConfig) (models.TokenResponse, error) {
	var authenticationResponse models.TokenResponse
	authRequest := models.RefreshTokenRequestWE{
		ClientID:     cred.ClientID,
		GrantType:    "refresh_token",
		RefreshToken: cred.RefreshToken,
		ClientSecret: cred.ClientSecret,
	}

	authRequestJSON, err := json.Marshal(authRequest)
	if err != nil {
		return models.TokenResponse{}, err
	}

	respBody, err := HTTPRequest[models.Token](http.MethodPost, utils.GetHostWebExperimentationAuth()+"/v1"+"/token", authRequestJSON)
	if err != nil {
		authResponse, err := InitiateBrowserAuth(cred.Username, cred.ClientID, cred.ClientSecret)
		if err != nil {
			return models.TokenResponse{}, err
		}

		return authResponse, nil
	}

	err = json.Unmarshal(respBody, &authenticationResponse)
	if err != nil {
		return models.TokenResponse{}, err
	}

	return authenticationResponse, err
}

func HTTPGetIdentifierWE() (models.UserMe, error) {
	var currentUser models.UserMe

	respBody, err := HTTPRequest[models.Token](http.MethodGet, utils.GetWebExperimentationHost()+"/v1"+"/users"+"/me", nil)
	if err != nil {
		return models.UserMe{}, err
	}

	err = json.Unmarshal(respBody, &currentUser)
	if err != nil {
		return models.UserMe{}, err
	}

	return currentUser, err
}

func InitiateBrowserAuth(username, clientID, clientSecret string) (models.TokenResponse, error) {
	if clientID == "" || clientSecret == "" {
		log.Fatal("Error while login, required fields (username, client ID, client secret)")
	}

	codeChan := make(chan string)
	var url = utils.GetWebExperimentationBrowserAuth(clientID, clientSecret)

	if err := openLink(url); err != nil {
		log.Fatalf("Error opening link: %s", err)
	}

	go func() {
		http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
			handleCallback(w, r, codeChan)
		})

		if err := http.ListenAndServe("127.0.0.1:8010", nil); err != nil {
			log.Fatalf("Error starting callback server: %s", err)
		}
	}()

	code := <-codeChan

	if code != "" {

		authenticationResponse, err := HTTPCreateTokenWEAuthorizationCode(clientID, clientSecret, code)
		if err != nil {
			return models.TokenResponse{}, err
		}

		if authenticationResponse.AccessToken == "" {
			return models.TokenResponse{}, errors.New("Credentials not valid.")
		}

		return authenticationResponse, nil
	}

	return models.TokenResponse{}, errors.New("Error occurred.")
}

func HTTPCreateTokenFE(clientId, clientSecret, accountId string) (models.TokenResponse, error) {
	var authenticationResponse models.TokenResponse
	authRequest := models.ClientCredentialsRequest{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scope:        "*",
		GrantType:    "client_credentials",
	}
	authRequestJSON, err := json.Marshal(authRequest)
	if err != nil {
		return models.TokenResponse{}, err
	}

	respBody, err := HTTPRequest[models.Token](http.MethodPost, utils.GetHostFeatureExperimentationAuth()+"/"+accountId+"/token?expires_in=43200", authRequestJSON)
	if err != nil {
		return models.TokenResponse{}, err
	}

	err = json.Unmarshal(respBody, &authenticationResponse)
	if err != nil {
		return models.TokenResponse{}, err
	}

	return authenticationResponse, err
}

/* func HTTPCreateTokenWE(clientId, clientSecret, accountId string) (models.TokenResponse, error) {
	var authenticationResponse models.TokenResponse
	authRequest := models.ClientCredentialsRequest{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	authRequestJSON, err := json.Marshal(authRequest)
	if err != nil {
		return models.TokenResponse{}, err
	}

	respBody, err := HTTPRequest[models.Token](http.MethodPost, utils.GetHostWebExperimentationAuth()+"/v1/token", authRequestJSON)
	if err != nil {
		return models.TokenResponse{}, err
	}

	err = json.Unmarshal(respBody, &authenticationResponse)
	if err != nil {
		return models.TokenResponse{}, err
	}

	return authenticationResponse, err
} */

func HTTPCreateTokenWEAuthorizationCode(client_id, client_secret, code string) (models.TokenResponse, error) {
	var authenticationResponse models.TokenResponse
	authRequest := models.AuthorizationCodeRequest{
		ClientID:     client_id,
		ClientSecret: client_secret,
		GrantType:    "authorization_code",
		Code:         code,
	}
	authRequestJSON, err := json.Marshal(authRequest)
	if err != nil {
		return models.TokenResponse{}, err
	}

	respBody, err := HTTPRequest[models.Token](http.MethodPost, utils.GetHostWebExperimentationAuth()+"/v1/token", authRequestJSON)
	if err != nil {
		return models.TokenResponse{}, err
	}

	err = json.Unmarshal(respBody, &authenticationResponse)
	if err != nil {
		return models.TokenResponse{}, err
	}

	return authenticationResponse, err
}

/* func HTTPCreateTokenWEPassword(client_id, client_secret, username, password, mfaCode string) (models.TokenResponse, error) {
	var authenticationResponse models.TokenResponse
	var mfaResponse models.MfaRequestWE
	var mfmResponse models.MfaRequestWE

	authRequest := models.PasswordRequest{
		ClientID:     client_id,
		ClientSecret: client_secret,
		GrantType:    "password",
		Username:     username,
		Password:     password,
	}
	authRequestJSON, err := json.Marshal(authRequest)
	if err != nil {
		return models.TokenResponse{}, err
	}

	mfaRespBody, err := HTTPRequest[models.MfaRequestWE](http.MethodPost, utils.GetHostWebExperimentationAuth()+"/v1/token", authRequestJSON)
	if err != nil {
		return models.TokenResponse{}, err
	}

	err = json.Unmarshal(mfaRespBody, &mfaResponse)
	if err != nil {
		return models.TokenResponse{}, err
	}

	mfmRequest := models.MultiFactorMethodRequestWE{
		GrantType: "multi_factor_methods",
		MfaToken:  mfaResponse.MfaToken,
		MfaMethod: "totp",
	}

	mfmRequestJSON, err := json.Marshal(mfmRequest)
	if err != nil {
		return models.TokenResponse{}, err
	}

	mfmRespBody, err := HTTPRequest[models.MfaRequestWE](http.MethodPost, utils.GetHostWebExperimentationAuth()+"/v1/token", mfmRequestJSON)
	if err != nil {
		return models.TokenResponse{}, err
	}

	err = json.Unmarshal(mfmRespBody, &mfmResponse)
	if err != nil {
		return models.TokenResponse{}, err
	}

	mfRequest := models.MultiFactorRequestWE{
		GrantType: "multi_factor",
		MfaToken:  mfmResponse.MfaToken,
		MfaMethod: "totp",
		Code:      mfaCode,
	}

	mfRequestJSON, err := json.Marshal(mfRequest)
	if err != nil {
		return models.TokenResponse{}, err
	}

	respBody, err := HTTPRequest[models.MfaRequestWE](http.MethodPost, utils.GetHostWebExperimentationAuth()+"/v1/token", mfRequestJSON)
	if err != nil {
		return models.TokenResponse{}, err
	}

	err = json.Unmarshal(respBody, &authenticationResponse)
	if err != nil {
		return models.TokenResponse{}, err
	}

	return authenticationResponse, err
}
*/

func HTTPCheckTokenFE() (models.Token, error) {
	return HTTPGetItem[models.Token](utils.GetHostFeatureExperimentationAuth() + "/token?access_token=" + cred.Token)
}
