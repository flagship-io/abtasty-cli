package web_experimentation

type Metric struct {
	Id           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	TestID       int    `json:"test_id"`
	AccountLevel bool   `json:"account_level"`
	Hidden       bool   `json:"hidden"`
	Context      any    `json:"context"`
}

type CustomTrackingMetricPayload struct {
	*Metric
	Value       string `json:"value"`
	CreatedFrom string `json:"created_from"`
}

type PageViewMetricPayload struct {
	*Metric
	PageConditions []ConditionItem `json:"page_conditions"`
}

type IndicatorMetricPayload struct {
	*Metric
	Category            string          `json:"category"`
	IndicatorConditions []ConditionItem `json:"indicator_conditions"`
}

type ConditionItem struct {
	Condition int    `json:"condition"`
	Value     string `json:"value"`
	Order     int    `json:"order"`
}

type MetricsData struct {
	Transactions    []TransactionMetric    `json:"transactions"`
	ActionTrackings []ActionTrackingMetric `json:"action_trackings"`
	WidgetTrackings []WidgetTrackingMetric `json:"widget_trackings"`
	CustomTrackings []CustomTrackingMetric `json:"custom_trackings"`
	PageViews       []PageViewMetric       `json:"page_views"`
	Indicators      []IndicatorMetric      `json:"indicators"`
}

type TransactionMetric struct {
	*Metric
	Affiliation string `json:"affiliation"`
	Favorite    bool   `json:"favorite"`
	UsedBy      []int  `json:"used_by"`
	Used        bool   `json:"used"`
}

type ActionTrackingMetric struct {
	*Metric
	UsedBy               []int  `json:"used_by"`
	Used                 bool   `json:"used"`
	GlobalModificationID int    `json:"global_modification_id"`
	Value                string `json:"value"`
	UsedByExtended       []any  `json:"used_by_extended"`
	UsedByMappings       []any  `json:"used_by_mappings"`
	CreatedAt            any    `json:"created_at"`
	Data                 any    `json:"data"`
	TestArchived         bool   `json:"test_archived"`
}

type WidgetTrackingMetric struct {
	*Metric
	UsedBy         []int  `json:"used_by"`
	Used           bool   `json:"used"`
	UsedByExtended []any  `json:"used_by_extended"`
	UsedByMappings []any  `json:"used_by_mappings"`
	CreatedAt      any    `json:"created_at"`
	PluginID       string `json:"plugin_id"`
	Siblings       int    `json:"siblings"`
	MainData       any    `json:"main_data"`
	Data           any    `json:"data"`
	TestArchived   bool   `json:"test_archived"`
}

type CustomTrackingMetric struct {
	*Metric
	Value          string `json:"value"`
	CreatedFrom    string `json:"created_from"`
	UsedByExtended []any  `json:"used_by_extended"`
	UsedByMappings []any  `json:"used_by_mappings"`
	Used           bool   `json:"used"`
}

type PageViewMetric struct {
	*Metric
	PageConditions []ConditionItem `json:"page_conditions"`
	Used           bool            `json:"used"`
	UsedBy         []int           `json:"used_by"`
	UsedByExtended []any           `json:"used_by_extended"`
	UsedByMappings []any           `json:"used_by_mappings"`
	CreatedAt      any             `json:"created_at"`
	TestArchived   bool            `json:"test_archived"`
}

type IndicatorMetric struct {
	*Metric
	Category string `json:"category"`
}
