package web_experimentation

type AccountWE struct {
	Id         int         `json:"id,omitempty"`
	Name       string      `json:"name"`
	Identifier string      `json:"identifier"`
	Role       string      `json:"role"`
	GlobalCode GlobalCode_ `json:"global_code"`
}

type CurrentAccountWE struct {
	Id        int    `json:"id,omitempty"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type GlobalCode_ struct {
	OnDomReady bool   `json:"on_dom_ready"`
	Value      string `json:"value"`
}
