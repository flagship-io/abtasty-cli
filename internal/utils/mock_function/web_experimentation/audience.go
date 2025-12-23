package web_experimentation

import (
	"net/http"

	models "github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	mockfunction "github.com/flagship-io/abtasty-cli/internal/utils/mock_function"
	"github.com/jarcoal/httpmock"
)

var TestAudience = models.Audience{
	Id:          "testAudienceId",
	Name:        "testAudienceName",
	Description: "testAudienceDesc",
	Hidden:      false,
	Archive:     false,
	IsSegment:   false,
}

var TestAudience1 = models.Audience{
	Id:          "testAudienceId1",
	Name:        "testAudienceName1",
	Description: "testAudienceDesc1",
	Hidden:      false,
	Archive:     false,
	IsSegment:   true,
}

var TestAudienceList = []models.Audience{
	TestAudience,
	TestAudience1,
}

func APIAudience() {

	respList := utils.HTTPListResponseWE[models.Audience]{
		Data: TestAudienceList,
		Pagination: utils.Pagination{
			Total:      1,
			Pages:      2,
			MaxPerPage: 10,
			Page:       1,
		},
	}

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/audiences/"+TestAudience.Id,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestAudience)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/audiences?status=unarchive&_page=1&_max_per_page=100",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respList)
			return resp, nil
		},
	)

}
