package model

import (
	"encoding/json"
	"github.com/sashabaranov/go-openai"
)

func (m *Message) GetToolCalls() []openai.ToolCall {
	if len(m.ToolCall) == 0 {
		return nil
	}
	var toolCalls []openai.ToolCall
	err := json.Unmarshal([]byte(m.ToolCall), &toolCalls)
	if err != nil {
		return nil
	}
	return toolCalls
}

func (m *Message) GetBackendToolCalls() []openai.ToolCall {
	if len(m.BackendToolCall) == 0 {
		return nil
	}
	var toolCalls []openai.ToolCall
	err := json.Unmarshal([]byte(m.BackendToolCall), &toolCalls)
	if err != nil {
		return nil
	}
	return toolCalls
}

func (m *Message) GetAllToolCalls() []openai.ToolCall {
	localToolCalls := m.GetToolCalls()
	backendToolCalls := m.GetBackendToolCalls()
	return append(localToolCalls, backendToolCalls...)
}

func (m *Message) SetNextMsgIds(ids []int64) {
	if len(ids) == 0 {
		return
	}
	idsStr, err := json.Marshal(ids)
	if err != nil {
		return
	}
	m.NextMsgIds = string(idsStr)
}

func (m *Message) GetNextMsgIds() []int64 {
	if len(m.NextMsgIds) == 0 {
		return nil
	}
	var ids []int64
	err := json.Unmarshal([]byte(m.NextMsgIds), &ids)
	if err != nil {
		return nil
	}
	return ids
}
