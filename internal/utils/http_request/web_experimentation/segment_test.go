package web_experimentation

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/stretchr/testify/assert"
)

var segmentRequester = SegmentRequester{&common.ResourceRequest{AccountID: "account_id"}}

func TestHTTPGetSegment(t *testing.T) {

	respBody, err := segmentRequester.HTTPGetSegment("segment-id")

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "segment-id", respBody.Id)
	assert.Equal(t, "testSegmentName", respBody.Name)
	assert.Equal(t, true, respBody.IsSegment)

}

func TestHTTPListSegment(t *testing.T) {

	respBody, err := segmentRequester.HTTPListSegment()

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "segment-id-1", respBody[1].Id)
	assert.Equal(t, "testSegmentName1", respBody[1].Name)
	assert.Equal(t, true, respBody[1].IsSegment)

}
