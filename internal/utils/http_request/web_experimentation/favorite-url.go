package web_experimentation

import (
	"net/http"

	models "github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
)

type FavoriteUrlRequester struct {
	*common.ResourceRequest
}

func (f *FavoriteUrlRequester) HTTPListFavoriteUrl() ([]models.FavoriteURL, error) {
	return common.HTTPGetAllPagesWE[models.FavoriteURL](utils.GetWebExperimentationHost() + "/v1/accounts/" + f.AccountID + "/favorite-urls?")
}

func (f *FavoriteUrlRequester) HTTPGetFavoriteUrl(id string) (models.FavoriteURL, error) {
	return common.HTTPGetItem[models.FavoriteURL](utils.GetWebExperimentationHost() + "/v1/accounts/" + f.AccountID + "/favorite-urls/" + id)
}

func (a *FavoriteUrlRequester) HTTPCreateFavoriteUrl(data []byte) ([]byte, error) {
	return common.HTTPRequest[models.FavoriteURL](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+a.AccountID+"/favorite-urls", data)
}
