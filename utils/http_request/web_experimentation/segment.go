package web_experimentation

import (
	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
)

type SegmentRequester struct {
	*common.ResourceRequest
}

func (t *SegmentRequester) HTTPListSegment() ([]models.Audience, error) {
	return common.HTTPGetAllPagesWE[models.Audience](utils.GetWebExperimentationHost() + "/v1/accounts/" + t.AccountID + "/audiences?type=segment&status=unarchive&")
}

func (t *SegmentRequester) HTTPGetSegment(id string) (models.Audience, error) {
	return common.HTTPGetItem[models.Audience](utils.GetWebExperimentationHost() + "/v1/accounts/" + t.AccountID + "/audiences/" + id)
}
