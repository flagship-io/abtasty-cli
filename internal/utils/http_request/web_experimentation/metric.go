package web_experimentation

import (
	"fmt"
	"net/http"

	models "github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
)

type MetricRequester struct {
	*common.ResourceRequest
}

func (a *MetricRequester) HTTPListMetrics() (models.MetricsData, error) {

	b, err := common.HTTPGetAllPagesWEMetric(utils.GetWebExperimentationHost() + "/v1/accounts/" + a.AccountID + "/metrics?")
	if err != nil {
		return models.MetricsData{}, err
	}

	return b, nil
}

func (a *MetricRequester) HTTPCreateMetric(data []byte) ([]byte, error) {
	return common.HTTPRequest[models.MetricsData](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+a.AccountID+"/metrics", data)
}

func (t *MetricRequester) HTTPEditMetric(id string, data []byte) ([]byte, error) {
	return common.HTTPRequest[models.MetricsData](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/metrics/"+id, data)
}

func (f *MetricRequester) HTTPDeleteMetric(id string) (string, error) {
	_, err := common.HTTPRequest[models.MetricsData](http.MethodDelete, utils.GetWebExperimentationHost()+"/v1/accounts/"+f.AccountID+"/metrics/"+id, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Metric %s deleted", id), nil
}
