package web_experimentation

import (
	"net/http"
	"os"
	"strconv"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	mockfunction_we "github.com/flagship-io/abtasty-cli/utils/mock_function/web_experimentation"
	"github.com/jarcoal/httpmock"
)

func init() {
	if os.Getenv("ABT_ENV") == "MOCK" {
		httpmock.Activate()
		mockfunction_we.APIVariation()
	}
}

type VariationWERequester struct {
	*common.ResourceRequest
}

func (v *VariationWERequester) HTTPGetVariation(testID, id int) (models.VariationWE, error) {
	return common.HTTPGetItem[models.VariationWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + v.AccountID + "/tests/" + strconv.Itoa(testID) + "/variations/" + strconv.Itoa(id))
}

func (v *VariationWERequester) HTTPDeleteVariation(testID, id int) error {
	_, err := common.HTTPRequest[models.VariationWE](http.MethodDelete, utils.GetWebExperimentationHost()+"/v1/accounts/"+v.AccountID+"/tests/"+strconv.Itoa(testID)+"/variations/"+strconv.Itoa(id), nil)
	return err
}
