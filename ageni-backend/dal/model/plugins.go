package model

import (
	"encoding/json"
	"errors"
	"time"
)

type PluginModel struct {
	ID              int64     `json:"id"`
	PluginKey       string    `json:"plugin_key"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ImageURL        string    `json:"image_url"`
	Status          int32     `json:"status"`
	PluginType      string    `json:"plugin_type"`
	HTTPPath        string    `json:"http_path"`
	RPCPath         string    `json:"rpc_path"`
	ContractAddress string    `json:"contract_address"`
	ContractAbi     string    `json:"contract_abi"`
	Headers         []Header  `json:"headers"`
	AuthInfo        AuthInfo  `json:"auth_info"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AuthInfo struct {
	AuthType      string             `json:"auth_type"`
	ServiceConfig *ServiceAuthConfig `json:"service_config,omitempty"`
	OAuth2Config  *OAuth2AuthConfig  `json:"oauth2_config,omitempty"`
}

type ServiceAuthConfig struct {
	Params   []AuthParam `json:"params"`
	Location string      `json:"location"`
}

type OAuth2AuthConfig struct {
	ClientID       string             `json:"client_id"`
	ClientSecret   string             `json:"client_secret"`
	AuthURL        string             `json:"auth_url"`
	TokenURL       string             `json:"token_url"`
	Scope          string             `json:"scope"`
	RedirectURI    string             `json:"redirect_uri"`
	TokenPlacement TokenPlacementInfo `json:"token_placement"`
}

type AuthParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	In    string `json:"in"`
}

type TokenPlacementInfo struct {
	Name   string `json:"name"`
	In     string `json:"in"`
	Prefix string `json:"prefix"`
}

func (p *Plugin) ToModel() (*PluginModel, error) {
	model := &PluginModel{
		ID:          p.ID,
		PluginKey:   p.PluginKey,
		Name:        p.Name,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		Status:      p.Status,
		PluginType:  p.PluginType,
		HTTPPath:    p.HTTPPath,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	var headers []Header
	if len(p.Headers) != 0 {
		if err := json.Unmarshal([]byte(p.Headers), &headers); err != nil {
			return nil, err
		}
		model.Headers = headers

	}
	if len(p.AuthInfo) != 0 {
		var authInfo AuthInfo
		if err := json.Unmarshal([]byte(p.AuthInfo), &authInfo); err != nil {
			return nil, err
		}
		model.AuthInfo = authInfo
	}

	return model, nil
}

func (p *Plugin) GetConfig(key string) (interface{}, error) {
	if p.Config == "" {
		return nil, errors.New("config is empty")
	}

	var configMap map[string]interface{}
	err := json.Unmarshal([]byte(p.Config), &configMap)
	if err != nil {
		return nil, err
	}

	value, exists := configMap[key]
	if !exists {
		return nil, errors.New("key not found in config")
	}

	return value, nil
}

func (p *Plugin) GetDependPluginIds() []int64 {
	dependIds, e := p.GetConfigIdList("depend_plugins")
	if e != nil {
		//log.Errorf("GetDependPluginIds err:%v", e)
		return nil
	}
	return dependIds
}

func (p *Plugin) GetConfigIdList(key string) ([]int64, error) {
	if p.Config == "" {
		return nil, errors.New("config is empty")
	}

	var configMap map[string]interface{}
	err := json.Unmarshal([]byte(p.Config), &configMap)
	if err != nil {
		return nil, err
	}

	value, exists := configMap[key]
	if !exists {
		return nil, errors.New("key not found in config")
	}

	ids, ok := value.([]interface{})
	if !ok {
		return nil, errors.New("invalid ids config")
	}
	return convertFloatToInt(ids), nil
}
