package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/readonme/open-studio/common/log"
	"github.com/readonme/open-studio/service/models"
	"github.com/sashabaranov/go-openai"
)

func (b *Bot) GetBot() *Bot {
	return b
}

func (b *Bot) GetConfig(key string) (interface{}, error) {
	if b.Config == "" {
		return nil, errors.New("config is empty")
	}

	var configMap map[string]interface{}
	err := json.Unmarshal([]byte(b.Config), &configMap)
	if err != nil {
		return nil, err
	}

	value, exists := configMap[key]
	if !exists {
		return nil, errors.New(fmt.Sprintf("key:%s not found in config", key))
	}

	return value, nil
}

func (b *Bot) GetConfigString(key string) (string, error) {
	i, e := b.GetConfig(key)
	if e != nil {
		return "", e
	}
	msg, ok := i.(string)
	if !ok {
		log.Error("GetConfigString assertion failed")
		return "", errors.New("assertion failed")
	}
	return msg, nil
}

func (b *Bot) GetWelcomeMsg() string {
	i, e := b.GetConfigString("welcome_msg")
	if e != nil {
		return ""
	}
	return i
}

func (b *Bot) GetPrompt() string {
	i, e := b.GetConfigString("prompt")
	if e != nil {
		return ""
	}
	return i
}

func (b *Bot) GetConfigIdList(key string) ([]int64, error) {
	if b.Config == "" {
		return nil, errors.New("config is empty")
	}

	var configMap map[string]interface{}
	err := json.Unmarshal([]byte(b.Config), &configMap)
	if err != nil {
		return nil, err
	}

	value, exists := configMap[key]
	if !exists {
		return nil, errors.New(fmt.Sprintf("key:%s not found in config", key))
	}

	ids, ok := value.([]interface{})
	if !ok {
		return nil, errors.New("invalid ids config")
	}
	return convertFloatToInt(ids), nil
}

func convertFloatToInt(floatSlice []interface{}) []int64 {
	intSlice := make([]int64, len(floatSlice))
	for i, v := range floatSlice {
		intSlice[i] = int64(v.(float64))
	}
	return intSlice
}

func (b *Bot) GetGuideInfo() []string {
	i, e := b.GetConfig("guide_info")
	if e != nil {
		return nil
	}
	msg, ok := i.([]interface{})
	if !ok {
		//log.Error("GetGuideInfo assertion failed")
	}
	ms, err := interfaceSliceToStringSlice(msg)
	if err != nil {
		return nil
	}
	return ms
}

func (b *Bot) GetModelSettings() *models.ModelSettings {
	defaultSettings := &models.ModelSettings{
		Model:       openai.GPT4o,
		Temperature: 1,
		TopP:        1,
		Rounds:      10,
		MaxLength:   4095,
	}
	i, e := b.GetConfig("model_settings")
	if e != nil {
		return defaultSettings
	}

	settingsMap, ok := i.(map[string]interface{})
	if !ok {
		log.Errorf("GetModelSettings: unexpected type for model_settings")
		return defaultSettings
	}

	jsonData, err := json.Marshal(settingsMap)
	if err != nil {
		log.Errorf("GetModelSettings: failed to marshal settings: %v", err)
		return defaultSettings
	}

	settings := &models.ModelSettings{}
	err = json.Unmarshal(jsonData, settings)
	if err != nil {
		log.Errorf("GetModelSettings: failed to unmarshal settings: %v", err)
		return defaultSettings
	}

	return settings
}

func interfaceSliceToStringSlice(slice []interface{}) ([]string, error) {
	result := make([]string, len(slice))
	for i, v := range slice {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element at index %d is not a string", i)
		}
		result[i] = str
	}
	return result, nil
}
