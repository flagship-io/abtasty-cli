package web_experimentation

type CampaignWE struct {
	Id                 int                         `json:"id,omitempty"`
	Name               string                      `json:"name"`
	Url                string                      `json:"url"`
	Description        string                      `json:"description"`
	Type               string                      `json:"type"`
	SubType            string                      `json:"sub_type"`
	Master             *CampaignWE                 `json:"master"`
	State              string                      `json:"state"`
	Traffic            *Traffic                    `json:"traffic"`
	Variations         []VariationWE               `json:"variations"`
	SubTests           []CampaignWE                `json:"sub_tests"`
	CreatingDate       DateTemplate                `json:"created_at"`
	Labels             []string                    `json:"labels"`
	LastPlayTimestamp  DateTemplate                `json:"last_play"`
	LastPauseTimestamp DateTemplate                `json:"last_pause"`
	GlobalCodeCampaign string                      `json:"global_code"`
	SourceCode         string                      `json:"source_code"`
	Audiences          []AudienceCampaign          `json:"audiences"`
	SelectorScopes     []SelectorScopesCampaign    `json:"selector_scopes"`
	CodeScopes         []CodeScopesCampaign        `json:"code_scopes"`
	FavoriteUrlScopes  []FavoriteUrlScopesCampaign `json:"favorite_url_scopes"`
	UrlScopes          []UrlScopesCampaign         `json:"url_scopes"`
	MutationObserver   bool                        `json:"mutation_observer"`
	DisplayFrequency   DisplayFrequencyCampaign    `json:"display_frequency"`
	Report             CampaignReport              `json:"report"`
}

type Traffic struct {
	Value                int    `json:"value"`
	LastIncreasedTraffic string `json:"last_increased_traffic"`
	Visitors             int    `json:"visitors"`
	OriginalVisitors     int    `json:"original_visitors"`
	VisitorsLimit        int    `json:"visitors_limit"`
}

type DateTemplate struct {
	ReadableDate string `json:"readable_date"`
	Timestamp    int    `json:"timestamp"`
	Pattern      string `json:"pattern"`
}

type CampaignState struct {
	Active bool `json:"active"`
}

type AudienceCampaign struct {
	Id               string       `json:"id,omitempty"`
	Name             string       `json:"name"`
	AudienceOriginID int          `json:"audience_origin_id"`
	Hidden           bool         `json:"hidden"`
	CreatedAt        DateTemplate `json:"created_at"`
	UpdatedAt        DateTemplate `json:"updated_at"`
	Archive          bool         `json:"archive"`
	IsSegment        bool         `json:"is_segment"`
	AccountID        int          `json:"account_id"`
}

type SelectorScopesCampaign struct {
	Id        int    `json:"id,omitempty"`
	Condition int    `json:"condition"`
	Include   bool   `json:"include"`
	Value     string `json:"value"`
}

type SelectorScopesCampaignModelJSON struct {
	Id        int    `json:"id,omitempty"`
	Condition string `json:"condition"`
	Value     string `json:"value"`
}

type CodeScopesCampaign struct {
	Id    int    `json:"id,omitempty"`
	Value string `json:"value"`
}

type FavoriteUrlScopesCampaign struct {
	Id            int    `json:"id,omitempty"`
	Include       bool   `json:"include"`
	FavoriteUrlID string `json:"favorite_url_id"`
	Name          string `json:"name,omitempty"`
}

type UrlScopesCampaign struct {
	Id        int    `json:"id,omitempty"`
	Condition int    `json:"condition"`
	Include   bool   `json:"include"`
	Value     string `json:"value"`
}

type UrlScopesCampaignModelJSON struct {
	Id        int    `json:"id,omitempty"`
	Condition string `json:"condition"`
	Value     string `json:"value"`
}

type DisplayFrequencyCampaign struct {
	Type  string `json:"type"`
	Unit  string `json:"unit"`
	Value int    `json:"value"`
}

type TargetingCampaign struct {
	AudienceIDs           []string                    `json:"audience_ids"`
	UrlScopes             []UrlScopesCampaign         `json:"url_scopes"`
	SelectorScopes        []SelectorScopesCampaign    `json:"selector_scopes"`
	CodeScopes            []CodeScopesCampaign        `json:"code_scopes"`
	FavoriteUrlScope      []FavoriteUrlScopesCampaign `json:"favorite_url_scopes"`
	DisplayFrequencyType  string                      `json:"display_frequency_type"`
	DisplayFrequencyUnit  string                      `json:"display_frequency_unit,omitempty"`
	DisplayFrequencyValue int                         `json:"display_frequency_value,omitempty"`
	MutationObserver      bool                        `json:"mutation_observer"`
}

type TargetingCampaignModelJSON struct {
	SegmentIDs []string `json:"segment_ids"`

	UrlScopes                   []UrlScopesCampaignModelJSON      `json:"url_scopes"`
	FavoriteUrlScopes           []FavoriteUrlScopesCampaign       `json:"favorite_url_scopes"`
	SelectorScopes              []SelectorScopesCampaignModelJSON `json:"selector_scopes"`
	CodeScope                   CodeScopesCampaign                `json:"code_scope"`
	ElementAppearsAfterPageLoad bool                              `json:"element_appears_after_page_load"`

	TriggerIDs []string `json:"triggers_ids"`

	TargetingFrequency TargetingFrequency `json:"targeting_frequency"`
}

type CampaignReport struct {
	Token   string `json:"token"`
	Comment string `json:"comment"`
}

type TargetingFrequency struct {
	Type  string `json:"type,omitempty"`
	Unit  string `json:"unit,omitempty"`
	Value int    `json:"value,omitempty"`
}
