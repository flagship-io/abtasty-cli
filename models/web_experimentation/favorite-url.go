package web_experimentation

type FavoriteURL struct {
	Id                    string        `json:"id,omitempty"`
	Name                  string        `json:"name"`
	AllPositiveConditions bool          `json:"all_positive_conditions"`
	AllNegativeConditions bool          `json:"all_negative_conditions"`
	CssSelectorDisplayed  bool          `json:"css_selector_displayed"`
	CssCode               string        `json:"css_code"`
	CreatedAt             DateTemplate  `json:"created_at"`
	UpdatedAt             DateTemplate  `json:"updated_at"`
	Conditions            []interface{} `json:"conditions"`
	DatalayerConditions   []interface{} `json:"datalayer_conditions"`
	CssSelectorConditions []interface{} `json:"css_selector_conditions"`
}
