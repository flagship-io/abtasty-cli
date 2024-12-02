package web_experimentation

import (
	"fmt"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
)

type FavoriteUrlRequester struct {
	*common.ResourceRequest
}

func (f *FavoriteUrlRequester) HTTPListFavoriteUrl() ([]models.FavoriteURL, error) {
	fmt.Println(utils.GetWebExperimentationHost() + "/v1/accounts/" + f.AccountID + "/favorite-urls")
	return common.HTTPGetAllPagesWE[models.FavoriteURL](utils.GetWebExperimentationHost() + "/v1/accounts/" + f.AccountID + "/favorite-urls?")
}

func (f *FavoriteUrlRequester) HTTPGetFavoriteUrl(id string) (models.FavoriteURL, error) {
	return common.HTTPGetItem[models.FavoriteURL](utils.GetWebExperimentationHost() + "/v1/accounts/" + f.AccountID + "/favorite-urls/" + id)
}
