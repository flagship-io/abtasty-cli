package common

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/flagship-io/flagship/utils"
	"github.com/flagship-io/flagship/utils/config"
)

var UserAgent string

var c = http.Client{Timeout: time.Duration(10) * time.Second}
var counter = false

type PageResult struct {
	Items      json.RawMessage `json:"items"`
	TotalCount int             `json:"total_count"`
}

type ResourceRequest struct {
	AccountID            string `mapstructure:"account_id"`
	AccountEnvironmentID string `mapstructure:"account_environment_id"`
}

func (c *ResourceRequest) Init(cL *RequestConfig) {
	c.AccountEnvironmentID = cL.AccountEnvironmentID
	c.AccountID = cL.AccountID
}

type PageResultWE struct {
	Data       json.RawMessage `json:"_data"`
	Pagination Pagination      `json:"_pagination"`
}

type Pagination struct {
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
	RefreshToken          string `mapstructure:"refresh_token"`
	CurrentUsedCredential string `mapstructure:"current_used_credential"`
}

var cred RequestConfig

func Init(credL RequestConfig) {
	cred = credL
}

func regenerateToken(product, configName string) {
	authenticationResponse, err := HTTPRefreshToken(cred.ClientID, cred.RefreshToken)

	if err != nil {
		log.Fatalf("%s", err)
	}

	if authenticationResponse.AccessToken == "" {
		log.Fatal("client_id or client_secret not valid")
	} else {
		cred.RefreshToken = authenticationResponse.RefreshToken
		cred.Token = authenticationResponse.AccessToken
		config.WriteToken(product, configName, authenticationResponse)
	}
}

func HTTPRequest[T any](method string, url string, body []byte) ([]byte, error) {
	var bodyIO io.Reader = nil
	if body != nil {
		bodyIO = bytes.NewBuffer(body)
	}

	//fmt.Println(method, url)
	var resource T

	resourceType := reflect.TypeOf(resource).String()

	if resourceType == "feature_experimentation.Campaign" || resourceType == "feature_experimentation.Goal" {
		if cred.AccountID == "" || cred.AccountEnvironmentID == "" {
			log.Fatalf("account_id or account_environment_id required, Please configure your CLI")
		}
	}

	req, err := http.NewRequest(method, url, bodyIO)
	if err != nil {
		log.Panicf("error occurred on request creation: %v", err)
	}

	if cred.Product == utils.FEATURE_EXPERIMENTATION {
		if cred.AccountID == "" {
			log.Fatalf("account_id required, Please configure your CLI")
		}
		// for resource loader
		if resourceType == "resource.ResourceData" && !strings.Contains(url, "token") && (cred.AccountID == "" || cred.AccountEnvironmentID == "") {
			log.Fatalf("account_id or account_environment_id required, Please configure your CLI")
		}

		if strings.Contains(url, "token") && cred.ClientID == "" && cred.ClientSecret == "" {
			log.Fatalf("client_id or client_secret required, Please configure your CLI")
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

	if body != nil {
		req.Header.Add("Content-Type", `application/json`)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

func HTTPGetAllPages[T any](resource string) ([]T, error) {
	currentPage := 1
	results := []T{}
	for {
		respBody, err := HTTPRequest[T](http.MethodGet, fmt.Sprintf("%s?_page=%d&_max_per_page=100", resource, currentPage), nil)
		if err != nil {
			return nil, err
		}
		pageResult := &PageResult{}
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

func HTTPGetAllPagesWe[T any](resource string) ([]T, error) {
	currentPage := 1
	results := []T{}
	for {
		respBody, err := HTTPRequest[T](http.MethodGet, fmt.Sprintf("%s?_page=%d&_max_per_page=100", resource, currentPage), nil)
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
