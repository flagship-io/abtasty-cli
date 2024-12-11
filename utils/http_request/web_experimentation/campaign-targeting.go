package web_experimentation

import (
	"net/http"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
)

type CampaignTargetingRequester struct {
	*common.ResourceRequest
}

func (c *CampaignTargetingRequester) HTTPPushCampaignTargeting(id string, code []byte) ([]byte, error) {
	return common.HTTPRequest[models.CampaignWE](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+c.AccountID+"/tests/"+id, []byte(code))
}

func (c *CampaignTargetingRequester) HTTPGetCampaignTargeting(id string) (models.TargetingCampaign, error) {
	resp, err := common.HTTPGetItem[models.CampaignWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + c.AccountID + "/tests/" + id)

	var audienceIds []string = []string{}
	var urlScopes []models.UrlScopesCampaign = []models.UrlScopesCampaign{}
	var selectorScopes []models.SelectorScopesCampaign = []models.SelectorScopesCampaign{}
	var codeScopes []models.CodeScopesCampaign = []models.CodeScopesCampaign{}
	var favoriteUrlScopes []models.FavoriteUrlScopesCampaign = []models.FavoriteUrlScopesCampaign{}
	var displayFrequencyType = resp.DisplayFrequency.Type
	var mutationObserver = resp.MutationObserver

	for _, audience := range resp.Audiences {
		audienceIds = append(audienceIds, audience.Id)
	}

	for _, urlScope := range resp.UrlScopes {
		urlScopes = append(urlScopes, models.UrlScopesCampaign{
			Value:     urlScope.Value,
			Include:   urlScope.Include,
			Condition: urlScope.Condition,
		})
	}

	for _, selectorScope := range resp.SelectorScopes {
		selectorScopes = append(selectorScopes, models.SelectorScopesCampaign{
			Value:     selectorScope.Value,
			Include:   selectorScope.Include,
			Condition: selectorScope.Condition,
		})
	}

	for _, codeScope := range resp.CodeScopes {
		codeScopes = append(codeScopes, models.CodeScopesCampaign{
			Value: codeScope.Value,
		})
	}

	for _, favoriteUrlScope := range resp.FavoriteUrlScopes {
		favoriteUrlScopes = append(favoriteUrlScopes, models.FavoriteUrlScopesCampaign{
			FavoriteUrlID: favoriteUrlScope.FavoriteUrlID,
			Include:       favoriteUrlScope.Include,
		})
	}

	targetingCampaign := models.TargetingCampaign{
		AudienceIDs:          audienceIds,
		SelectorScopes:       selectorScopes,
		UrlScopes:            urlScopes,
		CodeScopes:           codeScopes,
		DisplayFrequencyType: displayFrequencyType,
		FavoriteUrlScope:     favoriteUrlScopes,
		MutationObserver:     mutationObserver,
	}

	return targetingCampaign, err
}
