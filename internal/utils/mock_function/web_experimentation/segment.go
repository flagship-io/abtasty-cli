package web_experimentation

import (
	"net/http"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	mockfunction "github.com/flagship-io/abtasty-cli/utils/mock_function"
	"github.com/jarcoal/httpmock"
)

var TestSegment = models.Audience{
	Id:          "segment-id",
	Name:        "testSegmentName",
	Description: "testSegmentDescription",
	Hidden:      false,
	Archive:     false,
	IsSegment:   true,
}

var TestSegment1 = models.Audience{
	Id:          "segment-id-1",
	Name:        "testSegmentName1",
	Description: "testSegmentDescription1",
	Hidden:      false,
	Archive:     false,
	IsSegment:   true,
}

var TestSegmentList = []models.Audience{
	TestSegment,
	TestSegment1,
}

func APISegment() {

	respList := utils.HTTPListResponseWE[models.Audience]{
		Data: TestSegmentList,
		Pagination: utils.Pagination{
			Total:      1,
			Pages:      2,
			MaxPerPage: 10,
			Page:       1,
		},
	}

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/audiences/"+TestSegment.Id,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestSegment)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/audiences?type=segment&status=unarchive&_page=1&_max_per_page=100",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respList)
			return resp, nil
		},
	)

}
