package web_experimentation

import (
	"os"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	mockfunction_we "github.com/flagship-io/abtasty-cli/utils/mock_function/web_experimentation"
	"github.com/jarcoal/httpmock"
)

func init() {
	if os.Getenv("ABT_ENV") == "MOCK" {
		httpmock.Activate()
		mockfunction_we.APIFavoriteUrl()
	}
}

type FavoriteUrlRequester struct {
	*common.ResourceRequest
}

func (f *FavoriteUrlRequester) HTTPListFavoriteUrl() ([]models.FavoriteURL, error) {
	return common.HTTPGetAllPagesWE[models.FavoriteURL](utils.GetWebExperimentationHost() + "/v1/accounts/" + f.AccountID + "/favorite-urls?")
}

func (f *FavoriteUrlRequester) HTTPGetFavoriteUrl(id string) (models.FavoriteURL, error) {
	return common.HTTPGetItem[models.FavoriteURL](utils.GetWebExperimentationHost() + "/v1/accounts/" + f.AccountID + "/favorite-urls/" + id)
}
