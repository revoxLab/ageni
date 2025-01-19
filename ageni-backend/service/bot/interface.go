package bot

import (
	"encoding/json"
	"github.com/readonme/open-studio/dal/model"
	"github.com/sashabaranov/go-openai"
	"time"
)

type Bot interface {
	ProcessConversation(c *ProcessContext) (reply *ProcessedMsg, err error)
	GetBot() *model.Bot
}

type ProcessContext struct {
	Conversation *model.Conversation
	Messages     []openai.ChatCompletionMessage
}

type StudioToolCall struct {
	PluginId  int64 `json:"plugin_id"`
	IsBackEnd bool  `json:"is_backend"`
	openai.ToolCall
}

type ProcessedMsg struct {
	Id                int64                          `json:"id"`
	Content           string                         `json:"content"`
	ToolCalls         []StudioToolCall               `json:"tool_calls,omitempty"`
	BackendToolCalls  []StudioToolCall               `json:"backend_tool_calls,omitempty"`
	BackendToolResult []openai.ChatCompletionMessage `json:"backend_tool_result,omitempty"`
	OriginMessage     openai.ChatCompletionMessage   `json:"origin_message"`
	CreatedAt         time.Time                      `json:"created_at"`
}

func (m *ProcessedMsg) ToolCallStr() string {
	if len(m.ToolCalls) == 0 {
		return ""
	}
	//return json
	s, e := json.Marshal(m.ToolCalls)
	if e != nil {
		return ""
	}
	return string(s)

}

func (m *ProcessedMsg) BackendToolCallStr() string {
	if len(m.BackendToolCalls) == 0 {
		return ""
	}
	//return json
	s, e := json.Marshal(m.BackendToolCalls)
	if e != nil {
		return ""
	}
	return string(s)
}
