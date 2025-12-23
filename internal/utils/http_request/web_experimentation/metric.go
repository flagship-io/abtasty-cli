package web_experimentation

import (
	"fmt"
	"net/http"
	"strconv"

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

func (a *MetricRequester) HTTPGetMetric(id int) (interface{}, error) {
	b, err := common.HTTPGetAllPagesWEMetric(utils.GetWebExperimentationHost() + "/v1/accounts/" + a.AccountID + "/metrics?ids=" + strconv.Itoa(id) + "&")
	if err != nil {
		return nil, err
	}

	// Return the first metric type that has len != 0
	if len(b.Transactions) != 0 {
		return b.Transactions[0], nil
	}
	if len(b.ActionTrackings) != 0 {
		return b.ActionTrackings[0], nil
	}
	if len(b.WidgetTrackings) != 0 {
		return b.WidgetTrackings[0], nil
	}
	if len(b.CustomTrackings) != 0 {
		return b.CustomTrackings[0], nil
	}
	if len(b.PageViews) != 0 {
		return b.PageViews[0], nil
	}
	if len(b.Indicators) != 0 {
		return b.Indicators[0], nil
	}

	return nil, nil
}

func (a *MetricRequester) HTTPCreateMetric(data []byte) ([]byte, error) {
	return common.HTTPRequest[models.MetricsData](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+a.AccountID+"/metrics", data)
}

func (t *MetricRequester) HTTPEditMetric(id int, data []byte) ([]byte, error) {
	return common.HTTPRequest[models.MetricsData](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+t.AccountID+"/metrics/"+strconv.Itoa(id), data)
}

func (f *MetricRequester) HTTPDeleteMetric(id int) (string, error) {
	_, err := common.HTTPRequest[models.MetricsData](http.MethodDelete, utils.GetWebExperimentationHost()+"/v1/accounts/"+f.AccountID+"/metrics/"+strconv.Itoa(id), nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Metric %d deleted", id), nil
}
