package web_experimentation

import (
	"fmt"
	"log"
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

func (c *CampaignTargetingRequester) HTTPGetCampaignTargeting(id string) (models.TargetingCampaignModelJSON, error) {
	resp, err := common.HTTPGetItem[models.CampaignWE](utils.GetWebExperimentationHost() + "/v1/accounts/" + c.AccountID + "/tests/" + id)

	var segmentIds []string = []string{}
	var triggerIds []string = []string{}
	var urlScopes []models.UrlScopesCampaignModelJSON = []models.UrlScopesCampaignModelJSON{}
	var selectorScopes []models.SelectorScopesCampaignModelJSON = []models.SelectorScopesCampaignModelJSON{}
	var codeScope models.CodeScopesCampaign = models.CodeScopesCampaign{}
	var favoriteUrlScopes []models.FavoriteUrlScopesCampaign = []models.FavoriteUrlScopesCampaign{}
	var displayFrequencyType = resp.DisplayFrequency.Type
	var displayFrequencyUnit = resp.DisplayFrequency.Unit
	var displayFrequencyValue = resp.DisplayFrequency.Value
	var elementAppearsAfterPageLoad = resp.MutationObserver

	for _, audience := range resp.Audiences {
		if audience.IsSegment {
			segmentIds = append(segmentIds, audience.Id)
		} else {
			triggerIds = append(triggerIds, audience.Id)
		}
	}

	for _, urlScope := range resp.UrlScopes {
		condition, err := urlScopeConditionTransformer(urlScope.Condition, urlScope.Include)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		urlScopes = append(urlScopes, models.UrlScopesCampaignModelJSON{
			Value:     urlScope.Value,
			Condition: condition,
		})
	}

	for _, selectorScope := range resp.SelectorScopes {
		condition, err := selectorScopeConditionTransformer(selectorScope.Condition, selectorScope.Include)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		selectorScopes = append(selectorScopes, models.SelectorScopesCampaignModelJSON{
			Value:     selectorScope.Value,
			Condition: condition,
		})
	}

	if len(resp.CodeScopes) > 0 {
		codeScope = resp.CodeScopes[0]
	}

	for _, favoriteUrlScope := range resp.FavoriteUrlScopes {
		favoriteUrlScopes = append(favoriteUrlScopes, models.FavoriteUrlScopesCampaign{
			FavoriteUrlID: favoriteUrlScope.FavoriteUrlID,
			Include:       favoriteUrlScope.Include,
		})
	}

	targetingCampaign := models.TargetingCampaignModelJSON{
		SegmentIDs:                  segmentIds,
		UrlScopes:                   urlScopes,
		FavoriteUrlScopes:           favoriteUrlScopes,
		SelectorScopes:              selectorScopes,
		CodeScope:                   codeScope,
		ElementAppearsAfterPageLoad: elementAppearsAfterPageLoad,
		TriggerIDs:                  triggerIds,
		TargetingFrequency:          models.TargetingFrequency{Type: displayFrequencyType, Unit: displayFrequencyUnit, Value: displayFrequencyValue},
	}

	return targetingCampaign, err
}

func JsonModelToModel(campaignTargetingJSON models.TargetingCampaignModelJSON) models.TargetingCampaign {
	var audienceIds []string = []string{}
	var urlScopes []models.UrlScopesCampaign = []models.UrlScopesCampaign{}
	var selectorScopes []models.SelectorScopesCampaign = []models.SelectorScopesCampaign{}
	var codeScopes []models.CodeScopesCampaign = []models.CodeScopesCampaign{}
	var favoriteUrlScopes []models.FavoriteUrlScopesCampaign = []models.FavoriteUrlScopesCampaign{}
	var displayFrequencyType = campaignTargetingJSON.TargetingFrequency
	var elementAppearsAfterPageLoad = campaignTargetingJSON.ElementAppearsAfterPageLoad

	for _, segmentID := range campaignTargetingJSON.SegmentIDs {
		audienceIds = append(audienceIds, segmentID)
	}

	for _, triggerID := range campaignTargetingJSON.TriggerIDs {
		audienceIds = append(audienceIds, triggerID)
	}

	for _, urlScope := range campaignTargetingJSON.UrlScopes {
		include, err := getIncludeFromUrlScope(urlScope)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		condition, err := getConditionFromUrlScope(urlScope)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		urlScopes = append(urlScopes, models.UrlScopesCampaign{
			Value:     urlScope.Value,
			Include:   include,
			Condition: condition,
		})
	}

	for _, selectorScope := range campaignTargetingJSON.SelectorScopes {
		include, err := getIncludeFromSelectorScope(selectorScope)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		condition, err := getConditionFromSelectorScope(selectorScope)
		if err != nil {
			log.Fatalf("error occurred: %v", err)
		}

		selectorScopes = append(selectorScopes, models.SelectorScopesCampaign{
			Value:     selectorScope.Value,
			Include:   include,
			Condition: condition,
		})
	}

	if campaignTargetingJSON.CodeScope != (models.CodeScopesCampaign{}) {
		codeScopes = append(codeScopes, models.CodeScopesCampaign{
			Value: campaignTargetingJSON.CodeScope.Value,
		})
	}

	for _, favoriteUrlScope := range campaignTargetingJSON.FavoriteUrlScopes {
		favoriteUrlScopes = append(favoriteUrlScopes, models.FavoriteUrlScopesCampaign{
			FavoriteUrlID: favoriteUrlScope.FavoriteUrlID,
			Include:       favoriteUrlScope.Include,
		})
	}

	targetingCampaign := models.TargetingCampaign{
		AudienceIDs:           audienceIds,
		SelectorScopes:        selectorScopes,
		UrlScopes:             urlScopes,
		CodeScopes:            codeScopes,
		DisplayFrequencyType:  displayFrequencyType.Type,
		DisplayFrequencyUnit:  displayFrequencyType.Unit,
		DisplayFrequencyValue: displayFrequencyType.Value,
		FavoriteUrlScope:      favoriteUrlScopes,
		MutationObserver:      elementAppearsAfterPageLoad,
	}

	return targetingCampaign
}

func urlScopeConditionTransformer(condition int, include bool) (string, error) {
	switch condition {
	case 1:
		if include {
			return IS_EXACTLY, nil
		}
		return IS_NOT_EXACTLY, nil

	case 10:
		if include {
			return CONTAINS, nil
		}
		return DOES_NOT_CONTAINS, nil

	case 11:
		if include {
			return IS_REGULAR_EXPRESSION, nil
		}
		return IS_NOT_REGULAR_EXPRESSION, nil

	case 40:
		if include {
			return IS, nil
		}
		return IS_NOT, nil

	case 50:
		if include {
			return IS_SAVED_PAGE, nil
		}
		return IS_NOT_SAVED_PAGE, nil

	default:
		return "", fmt.Errorf(`Condition "%d" not found !`, condition)
	}
}

func selectorScopeConditionTransformer(condition int, include bool) (string, error) {
	switch condition {
	case 43:
		if include {
			return "is selector id", nil
		}
		return "is not selector id", nil

	case 44:
		if include {
			return "is selector class", nil
		}
		return "is not selector class", nil

	case 45:
		if include {
			return "is selector custom", nil
		}
		return "is not selector custom", nil

	default:
		return "", fmt.Errorf(`Condition "%d" not found !`, condition)
	}
}

func getConditionFromUrlScope(urlScope models.UrlScopesCampaignModelJSON) (int, error) {
	switch urlScope.Condition {
	case IS_EXACTLY, IS_NOT_EXACTLY:
		return 1, nil

	case CONTAINS, DOES_NOT_CONTAINS:
		return 10, nil

	case IS_REGULAR_EXPRESSION, IS_NOT_REGULAR_EXPRESSION:
		return 11, nil

	case IS, IS_NOT:
		return 40, nil

	case IS_SAVED_PAGE, IS_NOT_SAVED_PAGE:
		return 50, nil

	default:
		return 0, fmt.Errorf(`Condition "%s" not found !`, urlScope.Condition)
	}
}

func getIncludeFromUrlScope(urlScope models.UrlScopesCampaignModelJSON) (bool, error) {
	if urlScope.Condition == IS_EXACTLY || urlScope.Condition == CONTAINS ||
		urlScope.Condition == IS_REGULAR_EXPRESSION || urlScope.Condition == IS ||
		urlScope.Condition == IS_SAVED_PAGE {
		return true, nil
	}

	if urlScope.Condition == IS_NOT_EXACTLY || urlScope.Condition == DOES_NOT_CONTAINS ||
		urlScope.Condition == IS_NOT_REGULAR_EXPRESSION || urlScope.Condition == IS_NOT ||
		urlScope.Condition == IS_NOT_SAVED_PAGE {
		return false, nil
	}

	return false, fmt.Errorf(`Condition "%s" not found !`, urlScope.Condition)
}

func getConditionFromSelectorScope(selectorScope models.SelectorScopesCampaignModelJSON) (int, error) {
	switch selectorScope.Condition {
	case IS_SELECTOR_ID, IS_NOT_SELECTOR_ID:
		return 43, nil

	case IS_SELECTOR_CLASS, IS_NOT_SELECTOR_CLASS:
		return 44, nil

	case IS_SELECTOR_CUSTOM, IS_NOT_SELECTOR_CUSTOM:
		return 45, nil

	default:
		return 0, fmt.Errorf(`Condition "%s" not found !`, selectorScope.Condition)
	}
}

func getIncludeFromSelectorScope(selectorScope models.SelectorScopesCampaignModelJSON) (bool, error) {
	if selectorScope.Condition == IS_SELECTOR_ID || selectorScope.Condition == IS_SELECTOR_CLASS ||
		selectorScope.Condition == IS_NOT_SELECTOR_CUSTOM {
		return true, nil
	}

	if selectorScope.Condition == IS_NOT_SELECTOR_ID || selectorScope.Condition == IS_NOT_SELECTOR_CLASS ||
		selectorScope.Condition == IS_NOT_SELECTOR_CUSTOM {
		return false, nil
	}

	return false, fmt.Errorf(`Condition "%s" not found !`, selectorScope.Condition)
}
