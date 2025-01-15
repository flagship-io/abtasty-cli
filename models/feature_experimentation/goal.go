package feature_experimentation

type Goal struct {
	Id       string   `json:"id,omitempty"`
	Label    string   `json:"label,omitempty"`
	Type     string   `json:"type,omitempty"`
	Operator string   `json:"operator,omitempty"`
	Value    string   `json:"value,omitempty"`
	Metrics  []string `json:"metrics,omitempty"`
}
