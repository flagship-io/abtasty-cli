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
		mockfunction_fe.APIVariationGroup()
	}
}

type VariationGroupRequester struct {
	*common.ResourceRequest
}

func (vg *VariationGroupRequester) HTTPListVariationGroup(campaignID string) ([]models.VariationGroup, error) {
	return common.HTTPGetAllPagesFE[models.VariationGroup](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + vg.AccountID + "/account_environments/" + vg.AccountEnvironmentID + "/campaigns/" + campaignID + "/variation_groups")
}

func (vg *VariationGroupRequester) HTTPGetVariationGroup(campaignID, id string) (models.VariationGroup, error) {
	return common.HTTPGetItem[models.VariationGroup](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + vg.AccountID + "/account_environments/" + vg.AccountEnvironmentID + "/campaigns/" + campaignID + "/variation_groups/" + id)
}

func (vg *VariationGroupRequester) HTTPCreateVariationGroup(campaignID, data string) ([]byte, error) {
	return common.HTTPRequest[models.VariationGroup](http.MethodPost, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+vg.AccountID+"/account_environments/"+vg.AccountEnvironmentID+"/campaigns/"+campaignID+"/variation_groups", []byte(data))
}

func (vg *VariationGroupRequester) HTTPEditVariationGroup(campaignID, id, data string) ([]byte, error) {
	return common.HTTPRequest[models.VariationGroup](http.MethodPatch, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+vg.AccountID+"/account_environments/"+vg.AccountEnvironmentID+"/campaigns/"+campaignID+"/variation_groups/"+id, []byte(data))
}

func (vg *VariationGroupRequester) HTTPDeleteVariationGroup(campaignID, id string) error {
	_, err := common.HTTPRequest[models.VariationGroup](http.MethodDelete, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+vg.AccountID+"/account_environments/"+vg.AccountEnvironmentID+"/campaigns/"+campaignID+"/variation_groups/"+id, nil)
	return err
}
