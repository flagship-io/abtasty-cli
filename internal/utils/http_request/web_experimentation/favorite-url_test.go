package web_experimentation

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
	"github.com/stretchr/testify/assert"
)

var favoriteUrlRequester = FavoriteUrlRequester{&common.ResourceRequest{AccountID: "account_id"}}

func TestHTTPGetFavoriteUrl(t *testing.T) {

	respBody, err := favoriteUrlRequester.HTTPGetFavoriteUrl("testFavoriteUrlId")

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "testFavoriteUrlId", respBody.Id)
	assert.Equal(t, "testFavoriteUrlName", respBody.Name)
	assert.Equal(t, ".id{color: 'blue'}", respBody.CssCode)

}

func TestHTTPListFavoriteUrl(t *testing.T) {

	respBody, err := favoriteUrlRequester.HTTPListFavoriteUrl()

	assert.NotNil(t, respBody)
	assert.Nil(t, err)

	assert.Equal(t, "testFavoriteUrlId1", respBody[1].Id)
	assert.Equal(t, "testFavoriteUrlName1", respBody[1].Name)
	assert.Equal(t, ".id{color: 'red'}", respBody[1].CssCode)

}
