package feature_experimentation

import (
	"net/http"

	models "github.com/flagship-io/abtasty-cli/internal/models/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
)

type GoalRequester struct {
	*common.ResourceRequest
}

func (g *GoalRequester) HTTPListGoal() ([]models.Goal, error) {
	return common.HTTPGetAllPagesFE[models.Goal](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + g.AccountID + "/account_environments/" + g.AccountEnvironmentID + "/goals")
}

func (g *GoalRequester) HTTPGetGoal(id string) (models.Goal, error) {
	return common.HTTPGetItem[models.Goal](utils.GetFeatureExperimentationHost() + "/v1/accounts/" + g.AccountID + "/account_environments/" + g.AccountEnvironmentID + "/goals/" + id)
}

func (g *GoalRequester) HTTPCreateGoal(data []byte) ([]byte, error) {
	return common.HTTPRequest[models.Flag](http.MethodPost, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+g.AccountID+"/account_environments/"+g.AccountEnvironmentID+"/goals", data)
}

func (g *GoalRequester) HTTPEditGoal(id string, data []byte) ([]byte, error) {
	return common.HTTPRequest[models.Flag](http.MethodPatch, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+g.AccountID+"/account_environments/"+g.AccountEnvironmentID+"/goals/"+id, data)
}

func (g *GoalRequester) HTTPDeleteGoal(id string) error {
	_, err := common.HTTPRequest[models.Flag](http.MethodDelete, utils.GetFeatureExperimentationHost()+"/v1/accounts/"+g.AccountID+"/account_environments/"+g.AccountEnvironmentID+"/goals/"+id, nil)
	return err
}
