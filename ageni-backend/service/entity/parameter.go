package entity

type Parameter struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Visible      bool   `json:"visible"`
	Required     bool   `json:"required"`
	Description  string `json:"description"`
	InputMethod  string `json:"input_method"`
	DefaultValue bool   `json:"default_value"`
}

type ParameterSlice []*Parameter
