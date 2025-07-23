package web_experimentation

type VariationWE struct {
	Id           int         `json:"id,omitempty"`
	Name         string      `json:"name"`
	Description  string      `json:"description,omitempty"`
	Type         string      `json:"type,omitempty"`
	Traffic      int         `json:"traffic,omitempty"`
	VisualEditor bool        `json:"visual_editor,omitempty"`
	CodeEditor   bool        `json:"code_editor,omitempty"`
	Components   []Component `json:"components,omitempty"`
}

type VariationGlobalCode struct {
	Js  string `json:"js,omitempty"`
	Css string `json:"css,omitempty"`
}

type VariationResourceLoader struct {
	Name        string              `json:"name,omitempty"`
	Type        string              `json:"type,omitempty"`
	Description string              `json:"description,omitempty"`
	Traffic     int                 `json:"traffic,omitempty"`
	Code        VariationGlobalCode `json:"code,omitempty"`
}

type Component struct {
	Id          int      `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Js          string   `json:"js"`
	Css         string   `json:"css"`
	Html        string   `json:"html"`
	Form        string   `json:"form"`
	Options     string   `json:"options"`
}

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
