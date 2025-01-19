package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/readonme/open-studio/caller"
	"github.com/readonme/open-studio/common/log"
	"github.com/readonme/open-studio/dal"
	"github.com/readonme/open-studio/dal/model"
	"github.com/readonme/open-studio/dal/query"
	"github.com/readonme/open-studio/plugin_utils"
	"github.com/readonme/open-studio/service/bot"
	"github.com/readonme/open-studio/service/models"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const ConversationNormal = 1
const ConversationDebug = 2

type ConfigBot struct {
	*model.Bot
}

func (impl *ConfigBot) getOnlineConf() (conf *BotConfig, err error) {
	conf = &BotConfig{}
	conf.Prompt, err = impl.GetConfigString("prompt")
	if err != nil {
		return nil, err
	}
	conf.PluginIds, err = GetPluginIdList(impl.ID)
	if err != nil {
		return nil, err
	}
	conf.Settings = impl.GetModelSettings()
	return
}

type BotConfig struct {
	Prompt    string
	PluginIds []int64
	Settings  *models.ModelSettings
}

func (impl *ConfigBot) getDebugConf() (conf *BotConfig, err error) {
	d := query.Use(dal.StudioDB).BotDraft
	draft, err := d.WithContext(context.TODO()).Where(d.BotID.Eq(impl.ID)).Last()
	if err != nil {
		return nil, err
	}
	conf = &BotConfig{}
	conf.Prompt = draft.Prompt
	conf.PluginIds, err = draft.GetPluginIds()
	if err != nil {
		return nil, err
	}
	if len(draft.ModelSettings) != 0 {
		var settings *models.ModelSettings
		err := json.Unmarshal([]byte(draft.ModelSettings), &settings)
		if err != nil {
			return nil, err
		}
		conf.Settings = settings
	}

	return conf, nil
}

func (impl *ConfigBot) ProcessConversation(processContext *bot.ProcessContext) (reply *bot.ProcessedMsg, err error) {
	reply = new(bot.ProcessedMsg)
	var conf *BotConfig
	if processContext.Conversation.Type == ConversationDebug {
		conf, err = impl.getDebugConf()
	} else {
		conf, err = impl.getOnlineConf()
	}

	pluginMethodModels, err := plugin_utils.GetPluginMethodModels(conf.PluginIds)
	if err != nil {
		return nil, err
	}
	var finalDialog = []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: conf.Prompt},
	}
	finalDialog = append(finalDialog, processContext.Messages...)

	var tools []openai.Tool
	var fs []*openai.FunctionDefinition
	if pluginMethodModels != nil {
		for _, methodModel := range pluginMethodModels.Models {
			jsonSchema := plugin_utils.ConvertInputSchemaToJsonSchema(methodModel.Method.InputSchema)
			if err != nil {
				log.Errorf("Failed to convert input schema for method %d: %v", methodModel.Method.ID, err)
				continue
			}

			function := &openai.FunctionDefinition{
				Name:        methodModel.Method.Name,
				Description: methodModel.Method.Description,
			}
			if jsonSchema != nil {
				function.Parameters = jsonSchema
				function.Strict = true
			}
			fs = append(fs, function)
			tools = append(tools, openai.Tool{
				Type:     openai.ToolTypeFunction,
				Function: function,
			})
		}
	}

	settings := conf.Settings
	resp, err := caller.OpenAIClient.CreateChatCompletion(context.Background(),
		openai.ChatCompletionRequest{
			Model:       settings.Model,
			TopP:        settings.TopP,
			Temperature: settings.Temperature,
			MaxTokens:   int(settings.MaxLength),
			Messages:    finalDialog,
			Tools:       tools,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no response from AI")
	}

	msg := resp.Choices[0].Message
	localToolCalls := make([]bot.StudioToolCall, 0)
	backendToolCalls := make([]bot.StudioToolCall, 0)
	backendToolResult := make([]openai.ChatCompletionMessage, 0)
	for _, t := range msg.ToolCalls {
		methodModel, err := pluginMethodModels.GetByMethodName(t.Function.Name)
		if err != nil {
			log.Errorf("Failed to get method model for tool call %s: %v", t.Function.Name, err)
			return nil, err
		}
		t.Function.Name = methodModel.Method.Name
		if methodModel.Plugin.PluginType != "http" {
			localToolCalls = append(localToolCalls, bot.StudioToolCall{
				PluginId:  methodModel.Plugin.ID,
				IsBackEnd: false,
				ToolCall:  t,
			})
		} else {
			backendToolCalls = append(backendToolCalls, bot.StudioToolCall{
				PluginId:  methodModel.Plugin.ID,
				IsBackEnd: true,
				ToolCall:  t,
			})
			httpResp, httpErr := ProcessHttpCall(methodModel, t)
			if httpErr != nil {
				httpResp = openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					ToolCallID: t.ID,
					Content:    fmt.Sprintf("call function error, err=%v", httpErr),
				}
				log.Errorf("call function error, err=%v", httpErr)
			}
			backendToolResult = append(backendToolResult, httpResp)
		}
	}
	reply.Content = msg.Content
	reply.ToolCalls = localToolCalls
	reply.BackendToolCalls = backendToolCalls
	reply.BackendToolResult = backendToolResult
	reply.OriginMessage = msg
	return reply, nil
}

func ProcessHttpCall(method *plugin_utils.PluginMethodModel, t openai.ToolCall) (response openai.ChatCompletionMessage, err error) {

	result, e := processHttpCall(method, t)
	if e != nil {
		return response, e
	}
	return openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		ToolCallID: t.ID,
		Content:    result,
	}, nil
}

type AuthInfo struct {
	Location string            `json:"location"`
	Params   map[string]string `json:"params"`
}

type RequestParams struct {
	URLPath      string
	QueryParams  url.Values
	BodyParams   map[string]interface{}
	HeaderParams http.Header
}

func processHttpCall(method *plugin_utils.PluginMethodModel, t openai.ToolCall) (response string, err error) {
	var arguments map[string]interface{}
	err = json.Unmarshal([]byte(t.Function.Arguments), &arguments)
	if err != nil {
		return "", fmt.Errorf("parse argumnets faield: %v", err)
	}

	params := RequestParams{
		URLPath:      method.Plugin.HTTPPath + method.Method.HTTPSubPath,
		QueryParams:  url.Values{},
		BodyParams:   make(map[string]interface{}),
		HeaderParams: make(http.Header),
	}

	err = processParameters(method.Method.InputSchema.Parameters, arguments, &params)
	if err != nil {
		return "", err
	}

	err = addAuthentication(method.Plugin.AuthInfo, &params)
	if err != nil {
		return "", fmt.Errorf("add authentication failed: %v", err)
	}

	fullURL, err := url.Parse(params.URLPath)
	if err != nil {
		return "", fmt.Errorf("parse url path failed: %v", err)
	}
	fullURL.RawQuery = params.QueryParams.Encode()

	var reqBody []byte
	if method.Method.HTTPMethod != "GET" && len(params.BodyParams) > 0 {
		reqBody, err = json.Marshal(params.BodyParams)
		if err != nil {
			return "", fmt.Errorf("build body failed: %v", err)
		}
	}

	req, err := http.NewRequest(method.Method.HTTPMethod, fullURL.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("create request failed: %v", err)
	}

	for _, header := range method.Plugin.Headers {
		req.Header.Set(header.Name, header.Value)
	}

	for k, v := range params.HeaderParams {
		req.Header[k] = v
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response failed: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("request failed, code: %d, response: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

func addAuthentication(authInfo model.AuthInfo, params *RequestParams) error {
	switch authInfo.AuthType {
	case "service":
		s := authInfo.ServiceConfig
		for _, p := range s.Params {
			switch p.In {
			case "header":
				params.HeaderParams.Set(p.Name, p.Value)
			case "query":
				params.QueryParams.Set(p.Name, p.Value)
			case "body":
				params.BodyParams[p.Name] = p.Value
			default:
				return fmt.Errorf("unsupported auth location: %s", p.In)
			}
		}
	default:
		return fmt.Errorf("unsupported auth type: %s", authInfo.AuthType)
	}
	return nil
}

func GetPluginIdList(botId int64) ([]int64, error) {
	s := query.Use(dal.StudioDB).BotPlugin
	ps, err := s.WithContext(context.TODO()).Where(s.BotID.Eq(botId)).Select(s.PluginID).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	idSlice := make([]int64, len(ps))
	for _, p := range ps {
		idSlice = append(idSlice, p.PluginID)
	}
	return idSlice, nil
}

func processParameters(params []model.Parameter, arguments map[string]interface{}, reqParams *RequestParams) error {
	for _, param := range params {
		value, ok := arguments[param.Name]
		if !ok {
			if param.Required {
				return fmt.Errorf("need parameter: %s", param.Name)
			}
			continue
		}

		err := validateParameter(param, value)
		if err != nil {
			return err
		}

		switch param.InputMethod {
		case "path":
			reqParams.URLPath = strings.ReplaceAll(reqParams.URLPath, "{"+param.Name+"}", fmt.Sprint(value))
		case "query":
			if value != nil {
				reqParams.QueryParams.Set(param.Name, fmt.Sprint(value))
			}
		case "header":
			reqParams.HeaderParams.Set(param.Name, fmt.Sprint(value))
		case "body":
			if param.Type == "object" && param.Properties != nil {
				nestedObj := make(map[string]interface{})
				err := processNestedObject(param, value, nestedObj)
				if err != nil {
					return err
				}
				for k, v := range nestedObj {
					reqParams.BodyParams[k] = v
				}
			} else {
				reqParams.BodyParams[param.Name] = value
			}
		}
	}
	return nil
}
func validateParameter(param model.Parameter, value interface{}) error {
	if !param.Required && value == nil {
		return nil
	}

	switch param.Type {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("param %s should be string", param.Name)
		}
	case "integer":
		switch v := value.(type) {
		case float64:
			if v != float64(int(v)) {
				return fmt.Errorf("param %s should be integer", param.Name)
			}
		case int:
		default:
			return fmt.Errorf("param %s should be integer", param.Name)
		}
	case "number":
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("param %s should be number", param.Name)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("param %s should be boolean", param.Name)
		}
	case "object", "array":
	default:
		return fmt.Errorf("unknown parameter type: %s", param.Type)
	}
	return nil
}

func processNestedObject(param model.Parameter, value interface{}, parentObj map[string]interface{}) error {
	if param.Type == "object" {
		objValue, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("param %s should be object", param.Name)
		}
		for _, prop := range param.Properties {
			propValue, exists := objValue[prop.Name]
			if !exists {
				if prop.Required {
					if prop.DefaultValue != nil {
						propValue = prop.DefaultValue
					} else {
						return fmt.Errorf("missing required nested parameter: %s.%s", param.Name, prop.Name)
					}
				} else {
					continue
				}
			}
			err := validateParameter(prop, propValue)
			if err != nil {
				return err
			}
			if prop.Type == "object" {
				nestedObj := make(map[string]interface{})
				err = processNestedObject(prop, propValue, nestedObj)
				if err != nil {
					return err
				}
				objValue[prop.Name] = nestedObj
			} else {
				objValue[prop.Name] = propValue
			}
		}
		parentObj[param.Name] = objValue
	} else if param.Type == "array" {
		arrayValue, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("param %s should be array", param.Name)
		}
		for i, item := range arrayValue {
			err := validateParameter(model.Parameter{Type: param.Properties[0].Type, Name: fmt.Sprintf("%s[%d]", param.Name, i)}, item)
			if err != nil {
				return err
			}
		}
		parentObj[param.Name] = arrayValue
	}
	return nil
}
