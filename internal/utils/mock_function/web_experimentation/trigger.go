package web_experimentation

import (
	"net/http"

	models "github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	mockfunction "github.com/flagship-io/abtasty-cli/internal/utils/mock_function"
	"github.com/jarcoal/httpmock"
)

var TestTrigger = models.Audience{
	Id:          "trigger-id",
	Name:        "testTriggerName",
	Description: "testTriggerDescription",
	Hidden:      false,
	Archive:     false,
	IsSegment:   false,
}

var TestTrigger1 = models.Audience{
	Id:          "trigger-id-1",
	Name:        "testTriggerName1",
	Description: "testTriggerDescription1",
	Hidden:      false,
	Archive:     false,
	IsSegment:   false,
}

var TestTriggerList = []models.Audience{
	TestTrigger,
	TestTrigger1,
}

func APITrigger() {

	respList := utils.HTTPListResponseWE[models.Audience]{
		Data: TestTriggerList,
		Pagination: utils.Pagination{
			Total:      1,
			Pages:      2,
			MaxPerPage: 10,
			Page:       1,
		},
	}

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/audiences/"+TestTrigger.Id,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestTrigger)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/audiences?type=trigger&status=unarchive&_page=1&_max_per_page=100",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respList)
			return resp, nil
		},
	)

}
