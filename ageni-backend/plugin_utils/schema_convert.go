package plugin_utils

import (
	"context"
	"fmt"
	"github.com/readonme/open-studio/common/log"
	"github.com/readonme/open-studio/dal"
	"github.com/readonme/open-studio/dal/model"
	"github.com/readonme/open-studio/dal/query"
	"strconv"
	"strings"
)

func ConvertInputSchemaToJsonSchema(input model.InputSchema) map[string]interface{} {
	if len(input.Parameters) == 0 {
		return nil
	}
	jsonSchema := map[string]interface{}{
		"type":                 "object",
		"properties":           map[string]interface{}{},
		"required":             []string{},
		"additionalProperties": false,
	}
	for _, param := range input.Parameters {
		if !param.Visible {
			continue
		}
		jsonSchema["properties"].(map[string]interface{})[param.Name] = convertParameterToJsonSchema(param)
		jsonSchema["required"] = append(jsonSchema["required"].([]string), param.Name)
	}

	return jsonSchema
}

func convertParameterToJsonSchema(param model.Parameter) map[string]interface{} {
	property := map[string]interface{}{
		"description": param.Description,
		"type":        param.Type,
	}

	if len(param.Enum) > 0 {
		property["enum"] = param.Enum
	}
	if !param.Required {
		property["type"] = []string{param.Type, "null"}
	}

	switch param.Type {
	case "array":
		if len(param.Properties) > 0 {
			property["items"] = convertParameterToJsonSchema(param.Properties[0])
		}
	case "object":
		property["properties"] = map[string]interface{}{}
		required := []string{}
		for _, prop := range param.Properties {
			if !prop.Visible {
				continue
			}
			property["properties"].(map[string]interface{})[prop.Name] = convertParameterToJsonSchema(prop)
			required = append(required, prop.Name)
		}
		if len(required) > 0 {
			property["required"] = required
		}
	}

	return property
}
func main() {
	inputSchema := model.InputSchema{
		Parameters: []model.Parameter{
			{
				Name:        "username",
				Description: "username",
				Type:        "string",
				Required:    true,
			},
			{
				Name:         "age",
				Description:  "age",
				Type:         "integer",
				DefaultValue: 18,
			},
			{
				Name:        "preferences",
				Description: "preferences",
				Type:        "object",
				Properties: []model.Parameter{
					{
						Name:        "theme",
						Description: "theme",
						Type:        "string",
						Enum:        []interface{}{"light", "dark"},
					},
					{
						Name:         "notifications",
						Description:  "notifications",
						Type:         "boolean",
						DefaultValue: true,
					},
				},
			},
			{
				Name:        "tags",
				Description: "tags",
				Type:        "array",
				Properties: []model.Parameter{
					{
						Type: "string",
					},
				},
			},
		},
	}

	jsonSchema := ConvertInputSchemaToJsonSchema(inputSchema)
	prettyPrintJsonSchema(jsonSchema)
}

func prettyPrintJsonSchema(schema map[string]interface{}) {
	prettyPrint(schema, 0)
}

func prettyPrint(v interface{}, indent int) {
	switch val := v.(type) {
	case map[string]interface{}:
		fmt.Println("{")
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		for i, k := range keys {
			fmt.Printf("%s\"%s\": ", spaces(indent+2), k)
			prettyPrint(val[k], indent+2)
			if i < len(keys)-1 {
				fmt.Println(",")
			} else {
				fmt.Println()
			}
		}
		fmt.Printf("%s}", spaces(indent))
	case []interface{}:
		fmt.Println("[")
		for i, item := range val {
			fmt.Print(spaces(indent + 2))
			prettyPrint(item, indent+2)
			if i < len(val)-1 {
				fmt.Println(",")
			} else {
				fmt.Println()
			}
		}
		fmt.Printf("%s]", spaces(indent))
	default:
		fmt.Printf("%v", val)
	}
}

func spaces(n int) string {
	return fmt.Sprintf("%*s", n, "")
}

type PluginMethodModel struct {
	Plugin *model.PluginModel `json:"plugin"`
	Method *model.MethodModel `json:"method"`
}

type PluginMethodModelList struct {
	Models []PluginMethodModel
}

func (list *PluginMethodModelList) GetByMethodID(methodIDStr string) (*PluginMethodModel, error) {
	methodID, e := strconv.ParseInt(methodIDStr, 10, 64)
	if e != nil {
		return nil, e
	}
	for _, model := range list.Models {
		if model.Method.ID == methodID {
			return &model, nil
		}
	}
	return nil, fmt.Errorf("method with ID %d not found", methodID)
}

func (list *PluginMethodModelList) GetByMethodKey(methodKey string) (*PluginMethodModel, error) {
	names := strings.Split(methodKey, "-")
	if len(names) < 2 {
		return nil, fmt.Errorf("invalid method key: %s", methodKey)
	}
	pluginKey, methodName := names[0], names[1]
	for _, m := range list.Models {
		if m.Plugin.PluginKey == pluginKey && m.Method.Name == methodName {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("method with key %s not found", methodKey)
}

func (list *PluginMethodModelList) GetByMethodName(methodKey string) (*PluginMethodModel, error) {
	for _, m := range list.Models {
		if m.Method.Name == methodKey {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("method with key %s not found", methodKey)
}

func GetPluginMethodModels(pluginIds []int64) (*PluginMethodModelList, error) {
	if len(pluginIds) == 0 {
		return nil, nil
	}
	var result PluginMethodModelList

	p := query.Use(dal.StudioDB).Plugin
	plugins, err := p.WithContext(context.TODO()).Where(p.ID.In(pluginIds...)).Find()
	if err != nil {
		log.Errorf("get plugins err, pluginIds:%v, err:%v", pluginIds, err)
		return nil, err
	}

	m := query.Use(dal.StudioDB).Method
	methods, err := m.WithContext(context.TODO()).Where(m.PluginID.In(pluginIds...)).Find()
	if err != nil {
		log.Errorf("get plugin methods err, pluginIds:%v, err:%v", pluginIds, err)
		return nil, err
	}

	pluginMap := make(map[int64]*model.PluginModel)
	for _, plugin := range plugins {
		pluginModel, err := plugin.ToModel()
		if err != nil {
			log.Errorf("plugin to model err, pluginId:%v, err:%v", plugin.ID, err)
			continue
		}
		pluginMap[plugin.ID] = pluginModel
	}

	for _, method := range methods {
		methodModel, err := method.ToModel()
		if err != nil {
			log.Errorf("method to model err, methodId:%v, err:%v", method.ID, err)
			continue
		}

		pluginModel, ok := pluginMap[method.PluginID]
		if !ok {
			log.Errorf("plugin not found for method, methodId:%v, pluginId:%v", method.ID, method.PluginID)
			continue
		}

		pluginMethodModel := PluginMethodModel{
			Plugin: pluginModel,
			Method: methodModel,
		}
		result.Models = append(result.Models, pluginMethodModel)
	}

	return &result, nil
}
