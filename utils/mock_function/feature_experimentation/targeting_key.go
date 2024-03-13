package feature_experimentation

import (
	"net/http"

	models "github.com/flagship-io/flagship/models/feature_experimentation"
	"github.com/flagship-io/flagship/utils"
	"github.com/flagship-io/flagship/utils/config"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
)

var TestTargetingKey = models.TargetingKey{
	Id:          "testTargetingKeyID",
	Name:        "testTargetingKeyName",
	Type:        "string",
	Description: "testTargetingKeyDescription",
}

var TestTargetingKey1 = models.TargetingKey{
	Id:          "testTargetingKeyID1",
	Name:        "testTargetingKeyName1",
	Type:        "string",
	Description: "testTargetingKeyDescription1",
}

var TestTargetingKeyEdit = models.TargetingKey{
	Id:          "testTargetingKeyID",
	Name:        "testTargetingKeyName1",
	Type:        "string",
	Description: "testTargetingKeyDescription1",
}

var TestTargetingKeyList = []models.TargetingKey{
	TestTargetingKey,
	TestTargetingKey1,
}

func APITargetingKey() {

	config.SetViperMock()

	resp := utils.HTTPListResponse[models.TargetingKey]{
		Items:             TestTargetingKeyList,
		CurrentItemsCount: 2,
		CurrentPage:       1,
		TotalCount:        2,
		ItemsPerPage:      10,
		LastPage:          1,
	}

	httpmock.RegisterResponder("GET", utils.GetFeatureExperimentationHost()+"/v1/accounts/"+viper.GetString("account_id")+"/targeting_keys/"+TestTargetingKey.Id,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestTargetingKey)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", utils.GetFeatureExperimentationHost()+"/v1/accounts/"+viper.GetString("account_id")+"/targeting_keys",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, resp)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("POST", utils.GetFeatureExperimentationHost()+"/v1/accounts/"+viper.GetString("account_id")+"/targeting_keys",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestTargetingKey)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("PATCH", utils.GetFeatureExperimentationHost()+"/v1/accounts/"+viper.GetString("account_id")+"/targeting_keys/"+TestTargetingKey.Id,
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, TestTargetingKeyEdit)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("DELETE", utils.GetFeatureExperimentationHost()+"/v1/accounts/"+viper.GetString("account_id")+"/targeting_keys/"+TestTargetingKey.Id,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		},
	)
}
