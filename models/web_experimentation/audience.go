package web_experimentation

type Audience struct {
	Id              string       `json:"id,omitempty"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	Hidden          bool         `json:"hidden"`
	CreatedAt       DateTemplate `json:"created_at"`
	UpdatedAt       DateTemplate `json:"updated_at"`
	Archive         bool         `json:"archive"`
	IsSegment       bool         `json:"is_segment"`
	TestIDs         []int        `json:"test_ids"`
	LiveTestIDs     []int        `json:"live_test_ids"`
	ListTestsSource []struct {
		Id   int    `json:"id,omitempty"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"live_tests_source"`
	Groups []Group `json:"groups"`
}

type Group struct {
	Id         string        `json:"id,omitempty"`
	Targetings []TargetingWE `json:"targetings"`
}

type TargetingWE struct {
	Id               string `json:"id,omitempty"`
	Operator         string `json:"operator"`
	MutationObserver bool   `json:"mutation_observer"`
	Type             struct {
		Id   int    `json:"id,omitempty"`
		Name string `json:"name"`
	} `json:"type"`
	TimeFrame    int           `json:"timeframe,omitempty"`
	VisitedPages int           `json:"visited_pages,omitempty"`
	Conditions   []interface{} `json:"conditions"`
}
