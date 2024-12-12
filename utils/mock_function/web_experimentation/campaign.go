package web_experimentation

import (
	"net/http"
	"strconv"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	mockfunction "github.com/flagship-io/abtasty-cli/utils/mock_function"
	"github.com/jarcoal/httpmock"
)

var TestCampaign = models.CampaignWE{
	Id:                 100000,
	Name:               "testCampaignName",
	Description:        "testCampaignDescription",
	Type:               "ab",
	GlobalCodeCampaign: "console.log(\"Hello World!\")",
	Url:                "https://abtasty.com",
	UrlScopes: []models.UrlScopesCampaign{
		{
			Condition: 40,
			Include:   true,
			Value:     "https://abtasty.com",
		},
		{
			Condition: 41,
			Include:   false,
			Value:     "https://abtasty.com",
		},
	},
}

var TestCampaign1 = models.CampaignWE{
	Id:                 100001,
	Name:               "testCampaignName1",
	Description:        "testCampaignDescription1",
	Type:               "ab",
	GlobalCodeCampaign: "console.log(\"Hello Earth!\")",
	Url:                "https://abtasty.com",
}

var TestCampaignWithVariation = models.CampaignWE{
	Id:                 100002,
	Name:               "testCampaignName2",
	Description:        "testCampaignDescription2",
	Type:               "ab",
	GlobalCodeCampaign: "console.log(\"Hello World2!\")",
	Url:                "https://abtasty.com",
	UrlScopes: []models.UrlScopesCampaign{
		{
			Condition: 40,
			Include:   true,
			Value:     "https://abtasty.com",
		},
		{
			Condition: 41,
			Include:   false,
			Value:     "https://abtasty.com",
		},
	},
	Variations: []models.VariationWE{TestVariation},
}

var TestCampaignList = []models.CampaignWE{
	TestCampaign,
	TestCampaign1,
}

func APICampaign() {

	respList := utils.HTTPListResponseWE[models.CampaignWE]{
		Data: TestCampaignList,
		Pagination: utils.Pagination{
			Total:      1,
			Pages:      2,
			MaxPerPage: 10,
			Page:       1,
		},
	}

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/tests/"+strconv.Itoa(TestCampaign.Id),
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestCampaign)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/tests/"+strconv.Itoa(TestCampaignWithVariation.Id),
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestCampaignWithVariation)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/tests",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respList)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("POST", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/tests",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestCampaign)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("PATCH", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/tests/"+strconv.Itoa(TestCampaign.Id),
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestCampaign)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("DELETE", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/tests/"+strconv.Itoa(TestCampaign.Id),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		},
	)
}
