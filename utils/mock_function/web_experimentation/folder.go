package web_experimentation

import (
	"net/http"
	"strconv"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	mockfunction "github.com/flagship-io/abtasty-cli/utils/mock_function"
	"github.com/jarcoal/httpmock"
)

var TestFolder = models.Folder{
	Id:   900000,
	Name: "testFolderName",
}

var TestFolder1 = models.Folder{
	Id:   900001,
	Name: "testFolderName1",
}

var TestFolderList = []models.Folder{
	TestFolder,
	TestFolder1,
}

func APIFolder() {
	respList := utils.HTTPListResponseWE[models.Folder]{
		Data: TestFolderList,
		Pagination: utils.Pagination{
			Total:      1,
			Pages:      2,
			MaxPerPage: 10,
			Page:       1,
		},
	}

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/folders/"+strconv.Itoa(TestFolder.Id),
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestFolder)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/folders?_page=1&_max_per_page=100",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respList)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("POST", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/folders",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestFolder)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("PATCH", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/folders/"+strconv.Itoa(TestFolder.Id),
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestFolder)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("DELETE", utils.GetWebExperimentationHost()+"/v1/accounts/"+mockfunction.Auth.AccountID+"/folders/"+strconv.Itoa(TestFolder.Id),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		},
	)

}
