package feature_experimentation

import (
	"net/http"
	"os"

	models "github.com/flagship-io/abtasty-cli/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	mockfunction_fe "github.com/flagship-io/abtasty-cli/utils/mock_function/feature_experimentation"
	"github.com/jarcoal/httpmock"
)

func init() {
	if os.Getenv("ABT_ENV") == "MOCK" {
		httpmock.Activate()
		mockfunction_fe.APIVariation()
	}
}

type VariationFERequester struct {
	*common.ResourceRequest
}

func (v *VariationFERequester) HTTPListVariation(campaignID, variationGroupID string) ([]models.VariationFE, error) {
	return common.HTTPGetAllPagesFE[models.VariationFE](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + v.AccountID + "/account_environments/" + v.AccountEnvironmentID + "/campaigns/" + campaignID + "/variation_groups/" + variationGroupID + "/variations")
}

func (v *VariationFERequester) HTTPGetVariation(campaignID, variationGroupID, id string) (models.VariationFE, error) {
	return common.HTTPGetItem[models.VariationFE](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + v.AccountID + "/account_environments/" + v.AccountEnvironmentID + "/campaigns/" + campaignID + "/variation_groups/" + variationGroupID + "/variations/" + id)
}

func (v *VariationFERequester) HTTPCreateVariation(campaignID, variationGroupID, data string) ([]byte, error) {
	return common.HTTPRequest[models.VariationFE](http.MethodPost, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+v.AccountID+"/account_environments/"+v.AccountEnvironmentID+"/campaigns/"+campaignID+"/variation_groups/"+variationGroupID+"/variations", []byte(data))
}

func (v *VariationFERequester) HTTPEditVariation(campaignID, variationGroupID, id, data string) ([]byte, error) {
	return common.HTTPRequest[models.VariationFE](http.MethodPatch, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+v.AccountID+"/account_environments/"+v.AccountEnvironmentID+"/campaigns/"+campaignID+"/variation_groups/"+variationGroupID+"/variations/"+id, []byte(data))
}

func (v *VariationFERequester) HTTPDeleteVariation(campaignID, variationGroupID, id string) error {
	_, err := common.HTTPRequest[models.VariationFE](http.MethodDelete, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+v.AccountID+"/account_environments/"+v.AccountEnvironmentID+"/campaigns/"+campaignID+"/variation_groups/"+variationGroupID+"/variations/"+id, nil)
	return err
}
