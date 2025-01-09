package http_request

import (
	"os"

	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/flagship-io/abtasty-cli/utils/http_request/feature_experimentation"
	"github.com/flagship-io/abtasty-cli/utils/http_request/web_experimentation"
	mockfunction_fe "github.com/flagship-io/abtasty-cli/utils/mock_function/feature_experimentation"
	mockfunction_we "github.com/flagship-io/abtasty-cli/utils/mock_function/web_experimentation"
	"github.com/jarcoal/httpmock"
)

type HTTPResource interface {
	Init(*common.RequestConfig)
}

var ResourceRequester common.ResourceRequest

var HTTPResources = []HTTPResource{&ResourceRequester}

// feature experimentation
var CampaignFERequester feature_experimentation.CampaignFERequester = feature_experimentation.CampaignFERequester{ResourceRequest: &ResourceRequester}
var AccountEnvironmentFERequester feature_experimentation.AccountEnvironmentFERequester = feature_experimentation.AccountEnvironmentFERequester{ResourceRequest: &ResourceRequester}
var FlagRequester feature_experimentation.FlagRequester = feature_experimentation.FlagRequester{ResourceRequest: &ResourceRequester}
var GoalRequester feature_experimentation.GoalRequester = feature_experimentation.GoalRequester{ResourceRequest: &ResourceRequester}
var ProjectRequester feature_experimentation.ProjectRequester = feature_experimentation.ProjectRequester{ResourceRequest: &ResourceRequester}
var UserRequester feature_experimentation.UserRequester = feature_experimentation.UserRequester{ResourceRequest: &ResourceRequester}
var TargetingKeyRequester feature_experimentation.TargetingKeyRequester = feature_experimentation.TargetingKeyRequester{ResourceRequest: &ResourceRequester}
var VariationGroupRequester feature_experimentation.VariationGroupRequester = feature_experimentation.VariationGroupRequester{ResourceRequest: &ResourceRequester}
var VariationFERequester feature_experimentation.VariationFERequester = feature_experimentation.VariationFERequester{ResourceRequest: &ResourceRequester}
var PanicRequester feature_experimentation.PanicRequester = feature_experimentation.PanicRequester{ResourceRequest: &ResourceRequester}

// web experimentation
var CampaignWERequester web_experimentation.CampaignWERequester = web_experimentation.CampaignWERequester{ResourceRequest: &ResourceRequester}
var AccountWERequester web_experimentation.AccountWERequester = web_experimentation.AccountWERequester{ResourceRequest: &ResourceRequester}
var VariationWERequester web_experimentation.VariationWERequester = web_experimentation.VariationWERequester{ResourceRequest: &ResourceRequester}
var CampaignGlobalCodeRequester web_experimentation.CampaignGlobalCodeRequester = web_experimentation.CampaignGlobalCodeRequester{ResourceRequest: &ResourceRequester}
var AccountGlobalCodeRequester web_experimentation.AccountGlobalCodeRequester = web_experimentation.AccountGlobalCodeRequester{ResourceRequest: &ResourceRequester}
var ModificationRequester web_experimentation.ModificationRequester = web_experimentation.ModificationRequester{ResourceRequest: &ResourceRequester}
var AudienceRequester web_experimentation.AudienceRequester = web_experimentation.AudienceRequester{ResourceRequest: &ResourceRequester}
var SegmentRequester web_experimentation.SegmentRequester = web_experimentation.SegmentRequester{ResourceRequest: &ResourceRequester}
var TriggerRequester web_experimentation.TriggerRequester = web_experimentation.TriggerRequester{ResourceRequest: &ResourceRequester}
var FavoriteUrlRequester web_experimentation.FavoriteUrlRequester = web_experimentation.FavoriteUrlRequester{ResourceRequest: &ResourceRequester}
var CampaignTargetingRequester web_experimentation.CampaignTargetingRequester = web_experimentation.CampaignTargetingRequester{ResourceRequest: &ResourceRequester}

func init() {
	if os.Getenv("ABT_ENV") == "MOCK" {
		httpmock.Activate()

		// FE&R Resources
		mockfunction_fe.APIProject()
		mockfunction_fe.APICampaign()
		mockfunction_fe.APIAccountEnvironment()
		mockfunction_fe.APIFlag()
		mockfunction_fe.APIGoal()
		mockfunction_fe.APITargetingKey()
		mockfunction_fe.APIVariationGroup()
		mockfunction_fe.APIVariation()
		mockfunction_fe.APIUser()
		mockfunction_fe.APIPanic()
		mockfunction_fe.APIToken()

		// WE&P Resources
		mockfunction_we.APIAccount()
		mockfunction_we.APIAudience()
		mockfunction_we.APICampaign()
		mockfunction_we.APIFavoriteUrl()
		mockfunction_we.APIModification()
		mockfunction_we.APISegment()
		mockfunction_we.APIToken()
		mockfunction_we.APITrigger()
		mockfunction_we.APIVariation()
	}
}
