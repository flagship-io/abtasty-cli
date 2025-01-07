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
		mockfunction_fe.APIGoal()
	}
}

type GoalRequester struct {
	*common.ResourceRequest
}

func (g *GoalRequester) HTTPListGoal() ([]models.Goal, error) {
	return common.HTTPGetAllPagesFE[models.Goal](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + g.AccountID + "/account_environments/" + g.AccountEnvironmentID + "/goals")
}

func (g *GoalRequester) HTTPGetGoal(id string) (models.Goal, error) {
	return common.HTTPGetItem[models.Goal](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + g.AccountID + "/account_environments/" + g.AccountEnvironmentID + "/goals/" + id)
}

func (g *GoalRequester) HTTPCreateGoal(data string) ([]byte, error) {
	return common.HTTPRequest[models.Flag](http.MethodPost, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+g.AccountID+"/account_environments/"+g.AccountEnvironmentID+"/goals", []byte(data))
}

func (g *GoalRequester) HTTPEditGoal(id, data string) ([]byte, error) {
	return common.HTTPRequest[models.Flag](http.MethodPatch, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+g.AccountID+"/account_environments/"+g.AccountEnvironmentID+"/goals/"+id, []byte(data))
}

func (g *GoalRequester) HTTPDeleteGoal(id string) error {
	_, err := common.HTTPRequest[models.Flag](http.MethodDelete, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+g.AccountID+"/account_environments/"+g.AccountEnvironmentID+"/goals/"+id, nil)
	return err
}
