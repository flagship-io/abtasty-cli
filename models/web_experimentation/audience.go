package web_experimentation

import "encoding/json"

// Audience
type Audience struct {
	Id              string       `json:"id,omitempty"`
	Name            string       `json:"name"`
	Description     string       `json:"description,omitempty"`
	Hidden          bool         `json:"hidden,omitempty"`
	CreatedAt       DateTemplate `json:"created_at,omitempty"`
	UpdatedAt       DateTemplate `json:"updated_at,omitempty"`
	Archive         bool         `json:"archive,omitempty"`
	IsSegment       bool         `json:"is_segment,omitempty"`
	TestIDs         []int        `json:"test_ids,omitempty"`
	LiveTestIDs     []int        `json:"live_test_ids,omitempty"`
	ListTestsSource []struct {
		Id   int    `json:"id,omitempty"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"live_tests_source,omitempty"`
	Groups                []Group `json:"groups"`
	AllowDuplicatedConfig bool    `json:"allow_duplicated_config"`
}

type Group struct {
	Id         string           `json:"id,omitempty"`
	Targetings []TargetingGroup `json:"targetings"`
}

type TargetingGroup struct {
	Id               string `json:"id,omitempty"`
	Operator         string `json:"operator"`
	MutationObserver bool   `json:"mutation_observer"`
	TimeFrame        int    `json:"timeframe,omitempty"` // for PAGES_INTEREST
	VisitedPages     int    `json:"visited_pages,omitempty"`
	Conditions       []any  `json:"conditions"`
	Count            int    `json:"count,omitempty"`              // for PAGE_VIEW
	Type             any    `json:"type"`                         // struct{id int, name string}, or string/int
	CheckMode        string `json:"check_mode,omitempty"`         // Defined only for DATALAYER_TYPE and JS_VARIABLE targetings Possible values: loading, periodic, custom
	CheckModeLatency int    `json:"check_mode_latency,omitempty"` // for DATALAYER_TYPE, JS_VARIABLE
}

type AudiencePayload struct {
	Id                    string             `json:"id,omitempty"`
	Name                  string             `json:"name"`
	Description           string             `json:"description,omitempty"`
	Hidden                bool               `json:"hidden,omitempty"`
	Groups                [][]TargetingGroup `json:"groups"`
	AllowDuplicatedConfig bool               `json:"allow_duplicated_config"`
}

type Device struct {
	Include       bool `json:"include"`
	Value         int  `json:"value"`
	IsSegmentType bool `json:"is_segment_type,omitempty"`
}

type IPRange struct {
	Include     bool   `json:"include"`
	Range       bool   `json:"range"`
	From        string `json:"from"`
	To          string `json:"to,omitempty"`
	Description string `json:"description,omitempty"`
}

type DataLayer struct {
	Key       string   `json:"key"`
	Value     []string `json:"value"`
	Condition int      `json:"condition"`
}

type Code = IncludeValueString

type GeoLocation struct {
	Include          bool   `json:"include"`
	Country          string `json:"country"`
	City             string `json:"city,omitempty"`
	LeastSubdivision string `json:"least_subdivision,omitempty"`
	MostSubdivision  string `json:"most_subdivision,omitempty"`
}

type UrlParameter = IncludeNameValue
type Cookie = IncludeNameValue

type IncludeNameValue struct {
	Include bool   `json:"include"`
	Name    string `json:"name"`
	Value   string `json:"value,omitempty"`
}

type Selector struct {
	Include   bool   `json:"include"`
	Value     string `json:"value"`
	Condition int    `json:"condition"`
}

type Provider struct {
	Include        bool   `json:"include"`
	SegmentName    string `json:"segment_name"`
	Value          string `json:"value"`
	SecondaryValue string `json:"secondary_value,omitempty"`
	Condition      *int   `json:"condition,omitempty"`
	ProviderID     int    `json:"provider_id"`
}

type LandingPage = StringMatchPayload

type NewOrReturningVisitorPayload struct {
	NewVisitor bool `json:"new_visitor"`
}

// Non Supported
type Operator string

const (
	OpAnd  Operator = "and"
	OpOr   Operator = "or"
	OpAuto Operator = "auto"
)

type Condition struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type IncludeSegmentName struct {
	Include     bool   `json:"include"`
	SegmentName string `json:"segment_name"`
}

type IncludeSegmentIDString struct {
	Include   bool   `json:"include"`
	SegmentID string `json:"segment_id"`
}

type IncludeSegmentIDIntCond struct {
	Include      bool `json:"include"`
	Condition    int  `json:"condition,omitempty"`
	SegmentID    int  `json:"segment_id"`
	SegmentValue *int `json:"segment_value,omitempty"`
}

type IncludeSegTypeNameCond struct {
	Include        bool   `json:"include"`
	SegmentType    int    `json:"segment_type"`
	SegmentName    string `json:"segment_name"`
	Condition      int    `json:"condition"`
	AttributeValue string `json:"attribute_value,omitempty"`
	SegmentValue   *int   `json:"segment_value,omitempty"`
}

type IncludeSegTypeName struct {
	Include     bool   `json:"include"`
	SegmentName string `json:"segment_name"`
	SegmentType int    `json:"segment_type"`
}

type LangPayload struct {
	Include bool   `json:"include"`
	Lang    string `json:"lang"`
}

type StringMatchPayload struct {
	Include   bool   `json:"include"`
	Condition int    `json:"condition"`
	Value     string `json:"value"`
}

type BrowserPayload struct {
	Include     bool   `json:"include"`
	Value       int    `json:"value"`
	ValueCustom string `json:"value_custom,omitempty"`
	Version     string `json:"version,omitempty"`
}

type SourceTypePayload struct {
	Include bool `json:"include"`
	Type    int  `json:"type"`
}

type ScreenSizePayload struct {
	Include bool `json:"include"`
	Min     *int `json:"min,omitempty"`
	Max     *int `json:"max,omitempty"`
}

type JSVarPayload struct {
	Include   bool   `json:"include"`
	Name      string `json:"name"`
	Condition int    `json:"condition"`
	Value     string `json:"value,omitempty"`
}

type CampaignExpoPayload struct {
	Include     bool `json:"include"`
	Type        int  `json:"type"`
	CampaignID  int  `json:"campaign_id"`
	VariationID *int `json:"variation_id,omitempty"`
}

type IncludeValueString struct {
	Include bool   `json:"include"`
	Value   string `json:"value"`
}

type NumericWithConditionPayload struct {
	Include        bool `json:"include"`
	Value          int  `json:"value"`
	Condition      int  `json:"condition"`
	SecondaryValue *int `json:"secondary_value,omitempty"`
}

type ValueNonNegative struct {
	Include bool `json:"include"`
	Value   int  `json:"value"`
}

type WeatherPayload struct {
	Include   bool     `json:"include"`
	Condition int      `json:"condition"`
	Min       float64  `json:"min"`
	Max       *float64 `json:"max,omitempty"`
}

type EcommercePayload struct {
	Include bool   `json:"include"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

type CustomVarPayload struct {
	Include  bool   `json:"include"`
	Category string `json:"category"`
	Action   string `json:"action,omitempty"`
}

type LastPurchasePayload struct {
	Include    bool       `json:"include"`
	Condition  int        `json:"condition"`
	MinDate    string     `json:"min_date,omitempty"`
	MaxDate    string     `json:"max_date,omitempty"`
	MinDay     *int       `json:"min_day,omitempty"`
	MaxDay     *int       `json:"max_day,omitempty"`
	Properties []Property `json:"properties"`
}

type PurchaseFrequencyPayload struct {
	Include    bool       `json:"include"`
	Condition  int        `json:"condition"`
	Duration   int        `json:"duration"`
	Value      int        `json:"value"`
	MinDate    string     `json:"min_date,omitempty"`
	MaxDate    string     `json:"max_date,omitempty"`
	MinDay     *int       `json:"min_day,omitempty"`
	MaxDay     *int       `json:"max_day,omitempty"`
	Properties []Property `json:"properties"`
}

type KeywordPayload struct {
	Include bool     `json:"include"`
	Values  []string `json:"values"`
}

type EngagementPayload struct {
	Segment int `json:"segment"`
}

type ContentInterestPayload struct {
	Operator    string       `json:"operator"`
	Times       int          `json:"times"`
	Occurrences []Occurrence `json:"occurrences"`
}

type AbandonedCartPayload struct {
	Include    bool                    `json:"include"`
	Force      bool                    `json:"force,omitempty"`
	Properties []AbandonedCartProperty `json:"properties"`
}

type ActionTrackingPayload = IncludeValueString

type NumberPageViewPayload = ValueNonNegative

type DaySincePayload = ValueNonNegative

type IncludeOnly struct {
	Include bool `json:"include"`
}

type SourcePayload = IncludeValueString

type PageMatchPayload = StringMatchPayload

type IncludeSegmentIDInt struct {
	Include   bool `json:"include"`
	SegmentID int  `json:"segment_id"`
}

type Property struct {
	Property  int      `json:"property"`
	Values    []string `json:"values"`
	Condition int      `json:"condition"`
}

type Occurrence struct {
	Keyword string `json:"keyword"`
	Times   int    `json:"times"`
}

type AbandonedCartProperty struct {
	Type           string `json:"type"`
	Operator       string `json:"operator"`
	MainValue      any    `json:"main_value"`
	SecondaryValue any    `json:"secondary_value,omitempty"`
}

type CSATPayload struct {
	Feedback    string `json:"feedback"`
	Type        string `json:"type"`
	CampaignID  int    `json:"campaign_id"`
	VariationID *int   `json:"variation_id,omitempty"`
}

type NPSPayload struct {
	Condition      string `json:"condition"`
	Grade          int    `json:"grade"`
	SecondaryGrade *int   `json:"secondary_grade,omitempty"`
	Type           string `json:"type"`
	CampaignID     int    `json:"campaign_id"`
	VariationID    *int   `json:"variation_id,omitempty"`
}
