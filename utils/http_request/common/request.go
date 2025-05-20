package common

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/flagship-io/abtasty-cli/models"
	"github.com/flagship-io/abtasty-cli/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/config"
)

var UserAgent string
var OutputFormat string

var c = http.Client{Timeout: time.Duration(10) * time.Second}
var counter = false

type PageResultFE struct {
	Items      json.RawMessage `json:"items"`
	TotalCount int             `json:"total_count"`
}

type ResourceRequest struct {
	AccountID            string `mapstructure:"account_id"`
	AccountEnvironmentID string `mapstructure:"account_environment_id"`
	WorkingDir           string `mapstructure:"working_dir"`
	Identifier           string `mapstructure:"identifier"`
	Email                string `mapstructure:"email"`
}

func (c *ResourceRequest) Init(cL *RequestConfig) {
	c.AccountEnvironmentID = cL.AccountEnvironmentID
	c.AccountID = cL.AccountID
	c.WorkingDir = cL.WorkingDirectory
}

type PageResultWE struct {
	Data       json.RawMessage `json:"_data"`
	Pagination PaginationWE    `json:"_pagination"`
}

type PaginationWE struct {
	Total      int `json:"_total"`
	Pages      int `json:"_pages"`
	Page       int `json:"_page"`
	MaxPerPage int `json:"_max_per_page"`
}

type RequestConfig struct {
	Product               string
	Username              string `mapstructure:"username"`
	ClientID              string `mapstructure:"client_id"`
	ClientSecret          string `mapstructure:"client_secret"`
	AccountID             string `mapstructure:"account_id"`
	AccountEnvironmentID  string `mapstructure:"account_environment_id"`
	Token                 string `mapstructure:"token"`
	Scope                 string `mapstructure:"scope"`
	RefreshToken          string `mapstructure:"refresh_token"`
	CurrentUsedCredential string `mapstructure:"current_used_credential"`
	OutputFormat          string `mapstructure:"output_format"`
	WorkingDirectory      string `mapstructure:"working_dir"`
	Identifier            string `mapstructure:"identifier"`
	Email                 string `mapstructure:"email"`
}

type HitRequest struct {
	DS             string                `json:"ds"`
	ClientID       string                `json:"cid"`
	VisitorID      string                `json:"vid"`
	Type           string                `json:"t"`
	CustomVariable CustomVariableRequest `json:"cv"`
}

type CustomVariableRequest struct {
	Version        string `json:"version"`
	Timestamp      string `json:"timestamp"`
	StackType      string `json:"stack.type"`
	UserAgent      string `json:"user.agent"`
	OutputFormat   string `json:"output.format"`
	EnvironmentId  string `json:"envId"`
	AccountId      string `json:"accountId"`
	HttpMethod     string `json:"http.method"`
	HttpURL        string `json:"http.url"`
	ABTastyProduct string `json:"abtasty.product"`
	Identifier     string `json:"identifier"`
}

var cred RequestConfig

func Init(credL RequestConfig) {
	cred = credL
}

func regenerateToken(product, configName string) {
	var authenticationResponse models.TokenResponse
	var err error

	if product == utils.FEATURE_EXPERIMENTATION {
		authenticationResponse, err = HTTPCreateTokenFE(cred.ClientID, cred.ClientSecret, cred.AccountID)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
	} else {
		authenticationResponse, err = HTTPRefreshTokenWE(cred)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}
	}

	if authenticationResponse.AccessToken == "" {
		log.Fatal("client_id or client_secret not valid")
	}

	cred.RefreshToken = authenticationResponse.RefreshToken
	cred.Token = authenticationResponse.AccessToken
	err = config.RewriteToken(product, configName, authenticationResponse)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

}

func HTTPRequest[T any](method string, url string, body []byte) ([]byte, error) {
	if (cred.Product == utils.WEB_EXPERIMENTATION && !strings.Contains(cred.Email, "@abtasty.com")) || (cred.Product == utils.FEATURE_EXPERIMENTATION && cred.AccountEnvironmentID != "") {
		sendAnalyticHit(method, url)
	}

	var bodyIO io.Reader = nil
	if body != nil {
		bodyIO = bytes.NewBuffer(body)
	}

	var resource T

	resourceType := reflect.TypeOf(resource)

	if resourceType == reflect.TypeOf(feature_experimentation.Goal{}) || resourceType == reflect.TypeOf(feature_experimentation.CampaignFE{}) {
		if cred.AccountID == "" || cred.AccountEnvironmentID == "" {
			log.Fatalf("account_id or account_environment_id required, Please authenticate your CLI")
		}
	}

	req, err := http.NewRequest(method, url, bodyIO)
	if err != nil {
		log.Panicf("error occurred on request creation: %v", err)
	}

	if cred.Product == utils.FEATURE_EXPERIMENTATION {
		if (cred.Username == "" || cred.AccountID == "") && resourceType != reflect.TypeOf(models.Token{}) {
			log.Fatalf("username and account_id required, Please authenticate your CLI")
		}
		// for resource loader
		if resourceType.String() == "resource.ResourceData" && !strings.Contains(url, "token") && (cred.AccountID == "" || cred.AccountEnvironmentID == "") {
			log.Fatalf("account_id or account_environment_id required, Please authenticate your CLI")
		}

		/* 		if strings.Contains(url, "token") && cred.ClientID == "" && cred.ClientSecret == "" {
			log.Fatalf("client_id or client_secret required, Please authenticate your CLI")
		} */
	}

	if cred.Product == utils.WEB_EXPERIMENTATION {
		if resourceType != reflect.TypeOf(web_experimentation.AccountWE{}) && resourceType != reflect.TypeOf(web_experimentation.CurrentAccountWE{}) && !strings.Contains(url, "token") && !strings.Contains(url, "/users/me") && cred.AccountID == "" {
			log.Fatalf("username, account_id required, Please use the account command to select your account")
		}
	}

	if !strings.Contains(url, "token") && cred.Token == "" {
		regenerateToken(cred.Product, cred.CurrentUsedCredential)
	}

	req.Header.Add("Accept", `*/*`)
	req.Header.Add("Authorization", "Bearer "+cred.Token)
	req.Header.Add("Accept-Encoding", `gzip, deflate, br`)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return []byte{}, nil
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	respBody, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if (resp.StatusCode == 401) && !counter {
		counter = true
		regenerateToken(cred.Product, cred.CurrentUsedCredential)
		return HTTPRequest[T](method, url, body)
	}

	match, _ := regexp.MatchString("4\\d\\d|5\\d\\d", resp.Status)
	if match {
		err := errors.New(string(respBody))
		return nil, err
	}

	if cred.Product == utils.WEB_EXPERIMENTATION && method == "POST" && (resourceType == reflect.TypeOf(web_experimentation.CampaignWECommon{}) || resourceType == reflect.TypeOf(web_experimentation.VariationWE{}) || resourceType == reflect.TypeOf(web_experimentation.Audience{})) {
		return []byte(resp.Header.Get("location")), err
	}

	return respBody, err
}

func HTTPGetItem[T any](resource string) (T, error) {
	var result T
	respBody, err := HTTPRequest[T](http.MethodGet, resource, nil)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(respBody, &result)
	return result, err
}

func HTTPGetAllPagesFE[T any](resource string) ([]T, error) {
	currentPage := 1
	results := []T{}
	for {
		respBody, err := HTTPRequest[T](http.MethodGet, fmt.Sprintf("%s?_page=%d&_max_per_page=100", resource, currentPage), nil)
		if err != nil {
			return nil, err
		}
		pageResult := &PageResultFE{}
		err = json.Unmarshal(respBody, pageResult)
		if err != nil {
			return nil, err
		}

		typedItems := []T{}
		err = json.Unmarshal(pageResult.Items, &typedItems)
		if err != nil {
			return nil, err
		}
		results = append(results, typedItems...)

		if len(results) >= pageResult.TotalCount || len(pageResult.Items) == 0 {
			break
		}
		currentPage++
	}
	return results, nil
}

func HTTPGetAllPagesWE[T any](resource string) ([]T, error) {
	currentPage := 1
	results := []T{}
	for {
		respBody, err := HTTPRequest[T](http.MethodGet, fmt.Sprintf("%s_page=%d&_max_per_page=100", resource, currentPage), nil)
		if err != nil {
			return nil, err
		}
		pageResult := &PageResultWE{}
		err = json.Unmarshal(respBody, pageResult)
		if err != nil {
			return nil, err
		}

		typedItems := []T{}
		err = json.Unmarshal(pageResult.Data, &typedItems)
		if err != nil {
			return nil, err
		}
		results = append(results, typedItems...)

		if len(results) >= pageResult.Pagination.Total || len(pageResult.Data) == 0 {
			break
		}
		currentPage++
	}
	return results, nil
}

func sendAnalyticHit(method string, url string) (int, error) {
	var bodyIO io.Reader = nil
	var clientID = ""

	if cred.Product == utils.FEATURE_EXPERIMENTATION {
		clientID = cred.AccountEnvironmentID
	}

	if cred.Product == utils.WEB_EXPERIMENTATION {
		clientID = cred.Identifier
	}

	var customVariable = CustomVariableRequest{
		Version:        "1",
		Timestamp:      time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		StackType:      "Tools",
		UserAgent:      UserAgent,
		ABTastyProduct: cred.Product,
		EnvironmentId:  cred.AccountEnvironmentID,
		Identifier:     cred.Identifier,
		AccountId:      cred.AccountID,
		HttpMethod:     method,
		HttpURL:        url,
		OutputFormat:   OutputFormat,
	}

	var hit = HitRequest{
		DS:             "APP",
		ClientID:       clientID,
		VisitorID:      cred.AccountID,
		Type:           "USAGE",
		CustomVariable: customVariable,
	}

	body, err := json.Marshal(hit)
	if err != nil {
		log.Fatalf("error occurred: %v", err)
	}

	if body != nil {
		bodyIO = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(http.MethodPost, utils.HIT_ANALYTICS_URL, bodyIO)
	if err != nil {
		log.Panicf("error occurred on request creation: %v", err)
	}

	req.Header.Add("Accept", `*/*`)
	req.Header.Add("Accept-Encoding", `gzip, deflate, br`)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
