package web_experimentation

import (
	"net/http"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	mockfunction "github.com/flagship-io/abtasty-cli/utils/mock_function"
	"github.com/jarcoal/httpmock"
)

var TestFavoriteUrl = models.FavoriteURL{
	Id:      "testFavoriteUrlId",
	Name:    "testFavoriteUrlName",
	CssCode: ".id{color: 'blue'}",
}

var TestFavoriteUrl1 = models.FavoriteURL{
	Id:      "testFavoriteUrlId1",
	Name:    "testFavoriteUrlName1",
	CssCode: ".id{color: 'red'}",
}

var TestFavoriteUrlList = []models.FavoriteURL{
	TestFavoriteUrl,
	TestFavoriteUrl1,
}

func APIFavoriteUrl() {

	respList := utils.HTTPListResponseWE[models.FavoriteURL]{
		Data: TestFavoriteUrlList,
		Pagination: utils.Pagination{
			Total:      1,
			Pages:      2,
			MaxPerPage: 10,
			Page:       1,
		},
	}

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/favorite-urls/"+TestFavoriteUrl.Id,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestFavoriteUrl)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/favorite-urls",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respList)
			return resp, nil
		},
	)

}
