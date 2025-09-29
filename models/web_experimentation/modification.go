package web_experimentation

type Modification struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Value       string       `json:"value"`
	VariationID int          `json:"variation_id"`
	Selector    string       `json:"selector"`
	Engine      string       `json:"engine"`
	UpdatedBy   UpdatedBy_   `json:"updated_by"`
	UpdatedAt   DateTemplate `json:"updated_at"`
}

type UpdatedBy_ struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type ModificationWE struct {
	GlobalModifications []Modification `json:"global_modifications"`
	Modifications       []Modification `json:"modifications"`
}

type ModificationDataWE struct {
	Data ModificationWE `json:"_data"`
}

type ModificationResourceLoader struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	CampaignID  int    `json:"campaign_id,omitempty"`
	VariationID int    `json:"variation_id,omitempty"`
	Selector    string `json:"selector,omitempty"`
	Code        string `json:"code,omitempty"`
}
