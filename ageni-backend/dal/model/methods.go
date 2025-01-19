package model

import (
	"encoding/json"
	"time"
)

type MethodModel struct {
	ID             int64        `json:"id"`
	PluginID       int64        `json:"plugin_id"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	HTTPSubPath    string       `json:"http_sub_path"`
	HTTPMethod     string       `json:"http_method"`
	MethodCallName string       `json:"method_call_name"`
	InputSchema    InputSchema  `json:"input_schema"`
	OutputSchema   OutputSchema `json:"output_schema"`
	Status         int32        `json:"status"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type InputSchema struct {
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Type         string        `json:"type"`
	InputMethod  string        `json:"input_method"`
	Required     bool          `json:"required"`
	DefaultValue interface{}   `json:"default_value"`
	Visible      bool          `json:"visible"`
	Enum         []interface{} `json:"enum,omitempty"`
	Properties   []Parameter   `json:"properties,omitempty"`
}

type OutputSchema struct {
	Parameters []OutputField `json:"parameters"`
}

type OutputField struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Type         string        `json:"type"`
	DefaultValue interface{}   `json:"default_value,omitempty"`
	Example      interface{}   `json:"example,omitempty"`
	Properties   []OutputField `json:"properties,omitempty"`
	Items        *OutputField  `json:"items,omitempty"`
	Format       string        `json:"format,omitempty"`
}

func (m *Method) ToModel() (*MethodModel, error) {
	model := &MethodModel{
		ID:             m.ID,
		PluginID:       m.PluginID,
		Name:           m.Name,
		Description:    m.Description,
		HTTPSubPath:    m.HTTPSubPath,
		HTTPMethod:     m.HTTPMethod,
		MethodCallName: m.MethodCallName,
		Status:         m.Status,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}

	if len(m.InputSchema) != 0 {
		var inputSchema InputSchema
		if err := json.Unmarshal([]byte(m.InputSchema), &inputSchema); err != nil {
			return nil, err
		}
		model.InputSchema = inputSchema
	}

	if len(m.OutputSchema) != 0 {
		var outputSchema OutputSchema
		if err := json.Unmarshal([]byte(m.OutputSchema), &outputSchema); err != nil {
			return nil, err
		}
		model.OutputSchema = outputSchema
	}
	return model, nil
}
