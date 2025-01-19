package models

import (
	"github.com/golang/protobuf/ptypes/timestamp"
)

type BotListReq struct {
	Page     int32  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize int32  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Tab      string `protobuf:"bytes,3,opt,name=tab,proto3" json:"tab,omitempty"`
	Keywords string `protobuf:"bytes,4,opt,name=keywords,proto3" json:"keywords,omitempty"`
	PickType int32  `protobuf:"varint,5,opt,name=pick_type,json=pickType,proto3" json:"pick_type,omitempty"`
}

type BotListResp struct {
	Page     int32       `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize int32       `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Total    int32       `protobuf:"varint,3,opt,name=total,proto3" json:"total,omitempty"`
	Bots     []*AgentBot `protobuf:"bytes,4,rep,name=bots,proto3" json:"bots,omitempty"`
}

type AgentBot struct {
	Id            int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Image         string               `protobuf:"bytes,3,opt,name=image,proto3" json:"image,omitempty"`
	Desc          string               `protobuf:"bytes,4,opt,name=desc,proto3" json:"desc,omitempty"`
	Users         int32                `json:"users"`
	Conversations int32                `json:"conversations"`
	CreatedAt     *timestamp.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	WorkModes     []*WorkMode          `protobuf:"bytes,9,rep,name=work_modes,json=workModes,proto3" json:"work_modes,omitempty"`
	Configuration []*Configuration     `protobuf:"bytes,10,rep,name=configuration,proto3" json:"configuration,omitempty"`
	Creator       *Creator             `protobuf:"bytes,11,opt,name=creator,proto3" json:"creator,omitempty"`
	CreatorId     int64                `protobuf:"varint,12,opt,name=creator_id,json=creatorId,proto3" json:"creator_id,omitempty"`
	WelcomeMsg    string               `protobuf:"bytes,13,opt,name=welcome_msg,json=welcomeMsg,proto3" json:"welcome_msg,omitempty"`
	GuideInfo     []string             `protobuf:"bytes,14,rep,name=guide_info,json=guideInfo,proto3" json:"guide_info,omitempty"`
	Tab           string               `protobuf:"bytes,15,opt,name=tab,proto3" json:"tab,omitempty"`
	LinkedPlugin  []*Plugin            `protobuf:"bytes,16,rep,name=linked_plugin,json=linkedPlugin,proto3" json:"linked_plugin,omitempty"`
	Status        int32                `protobuf:"varint,17,opt,name=status,proto3" json:"status,omitempty"`
}
type WorkMode struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

type Configuration struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

type Creator struct {
	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	HeadPic string `protobuf:"bytes,2,opt,name=head_pic,json=headPic,proto3" json:"head_pic,omitempty"`
}

type Plugin struct {
	Id            int64                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Image         string               `protobuf:"bytes,3,opt,name=image,proto3" json:"image,omitempty"`
	Desc          string               `protobuf:"bytes,4,opt,name=desc,proto3" json:"desc,omitempty"`
	Tab           string               `protobuf:"bytes,5,opt,name=tab,proto3" json:"tab,omitempty"`
	Methods       []*Method            `protobuf:"bytes,6,rep,name=methods,proto3" json:"methods,omitempty"`
	Used          int32                `protobuf:"varint,7,opt,name=used,proto3" json:"used,omitempty"`
	LinkedAgent   []*AgentBot          `protobuf:"bytes,8,rep,name=linked_agent,json=linkedAgent,proto3" json:"linked_agent,omitempty"`
	CreatedAt     *timestamp.Timestamp `protobuf:"bytes,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Creator       *Creator             `protobuf:"bytes,11,opt,name=creator,proto3" json:"creator,omitempty"`
	CreatorId     int32                `protobuf:"varint,12,opt,name=creator_id,json=creatorId,proto3" json:"creator_id,omitempty"`
	DependPlugins []int64              `protobuf:"varint,13,rep,packed,name=depend_plugins,json=dependPlugins,proto3" json:"depend_plugins,omitempty"`
}

type Method struct {
	Id             int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	PluginId       int64  `protobuf:"varint,2,opt,name=plugin_id,json=pluginId,proto3" json:"plugin_id,omitempty"`
	Name           string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description    string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	HttpSubPath    string `protobuf:"bytes,5,opt,name=http_sub_path,json=httpSubPath,proto3" json:"http_sub_path,omitempty"`
	HttpMethod     string `protobuf:"bytes,6,opt,name=http_method,json=httpMethod,proto3" json:"http_method,omitempty"`
	MethodCallName string `protobuf:"bytes,7,opt,name=method_call_name,json=methodCallName,proto3" json:"method_call_name,omitempty"`
	InputSchema    string `protobuf:"bytes,8,opt,name=input_schema,json=inputSchema,proto3" json:"input_schema,omitempty"`
	InputExample   string `protobuf:"bytes,9,opt,name=input_example,json=inputExample,proto3" json:"input_example,omitempty"`
	OutputSchema   string `protobuf:"bytes,10,opt,name=output_schema,json=outputSchema,proto3" json:"output_schema,omitempty"`
	OutputExample  string `protobuf:"bytes,11,opt,name=output_example,json=outputExample,proto3" json:"output_example,omitempty"`
	Status         int32  `protobuf:"varint,12,opt,name=status,proto3" json:"status,omitempty"`
}

type UserBotListReq struct {
	Page    int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PerPage int32 `protobuf:"varint,2,opt,name=per_page,json=perPage,proto3" json:"per_page,omitempty"`
	UserId  int64 `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type CreateConversationReq struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	BotId  int64  `protobuf:"varint,2,opt,name=bot_id,json=botId,proto3" json:"bot_id,omitempty"`
	Title  string `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Type   string `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"` //normal or debug
}

type CreateConversationResp struct {
	BotId          int64                `protobuf:"varint,1,opt,name=bot_id,json=botId,proto3" json:"bot_id,omitempty"`
	ConversationId int64                `protobuf:"varint,2,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	CreatedAt      *timestamp.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	WelcomeMsg     string               `protobuf:"bytes,4,opt,name=welcome_msg,json=welcomeMsg,proto3" json:"welcome_msg,omitempty"`
}

type ConversationHistoryReq struct {
	ConversationId int64 `protobuf:"varint,1,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	Page           int64 `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	PageSize       int64 `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	UserId         int64 `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type ConversationHistoryResp struct {
	ConversationId int64           `protobuf:"varint,1,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	CreatedHuman   string          `protobuf:"bytes,2,opt,name=created_human,json=createdHuman,proto3" json:"created_human,omitempty"`
	Messages       []*AgentMessage `protobuf:"bytes,3,rep,name=messages,proto3" json:"messages,omitempty"`
}

type AgentMessage struct {
	Id             int64                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ConversationId int64                `protobuf:"varint,2,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	MessageType    string               `protobuf:"bytes,3,opt,name=message_type,json=messageType,proto3" json:"message_type,omitempty"`
	Content        string               `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	ToolCalls      []*ToolCall          `protobuf:"bytes,5,rep,name=tool_calls,json=toolCalls,proto3" json:"tool_calls,omitempty"`
	CreatedAt      *timestamp.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Status         int32                `protobuf:"varint,7,opt,name=status,proto3" json:"status,omitempty"`
	NextMessageIds []int64              `protobuf:"varint,8,rep,packed,name=next_message_ids,json=nextMessageIds,proto3" json:"next_message_ids,omitempty"`
	ToolCallId     string               `protobuf:"bytes,9,opt,name=tool_call_id,json=toolCallId,proto3" json:"tool_call_id,omitempty"`
}

type ToolCall struct {
	Id        string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type      string            `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Function  *ToolCallFunction `protobuf:"bytes,3,opt,name=function,proto3" json:"function,omitempty"`
	PluginId  int64             `protobuf:"varint,4,opt,name=plugin_id,json=pluginId,proto3" json:"plugin_id,omitempty"`
	IsBackend bool              `protobuf:"varint,5,opt,name=is_backend,json=isBackend,proto3" json:"is_backend,omitempty"`
}

type ToolCallFunction struct {
	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Arguments string `protobuf:"bytes,2,opt,name=arguments,proto3" json:"arguments,omitempty"`
}

type MessageDetailReq struct {
	MessageId int64 `protobuf:"varint,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
}

type MessageDetailResp struct {
	Data *AgentMessage `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

type SendMessageReq struct {
	UserId         int64         `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ConversationId int64         `protobuf:"varint,2,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	Content        string        `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	ToolResults    []*ToolResult `protobuf:"bytes,4,rep,name=tool_results,json=toolResults,proto3" json:"tool_results,omitempty"`
}

type ToolResult struct {
	ToolCallId string `protobuf:"bytes,1,opt,name=tool_call_id,json=toolCallId,proto3" json:"tool_call_id,omitempty"`
	Content    string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

type SendMessageResp struct {
	MessageId      int64 `protobuf:"varint,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	RespMessageId  int64 `protobuf:"varint,2,opt,name=resp_message_id,json=respMessageId,proto3" json:"resp_message_id,omitempty"`
	ConversationId int64 `protobuf:"varint,3,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
}

type ConversationListReq struct {
	Page     int64 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize int64 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	BotId    int64 `protobuf:"varint,3,opt,name=bot_id,json=botId,proto3" json:"bot_id,omitempty"`
	UserId   int64 `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type ConversationListResp struct {
	Data []*Conversation `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

type Conversation struct {
	Id    int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
}

type BotDetailReq struct {
	BotId int64 `protobuf:"varint,1,opt,name=bot_id,json=botId,proto3" json:"bot_id,omitempty"`
}

type BotDetailResp struct {
	Bot *AgentBot `protobuf:"bytes,1,opt,name=bot,proto3" json:"bot,omitempty"`
}

type ConversationDeleteReq struct {
	Id     int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId int64 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type PluginListReq struct {
	Tab      string  `protobuf:"bytes,1,opt,name=tab,proto3" json:"tab,omitempty"`
	Page     int32   `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	PageSize int32   `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Keywords string  `protobuf:"bytes,4,opt,name=keywords,proto3" json:"keywords,omitempty"`
	Ids      []int64 `protobuf:"varint,5,rep,packed,name=ids,proto3" json:"ids,omitempty"`
}
type PluginListResp struct {
	Page     int32     `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize int32     `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Total    int32     `protobuf:"varint,3,opt,name=total,proto3" json:"total,omitempty"`
	List     []*Plugin `protobuf:"bytes,4,rep,name=list,proto3" json:"list,omitempty"`
}

type PluginDetailReq struct {
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

type PluginDetailResp struct {
	Plugin *Plugin `protobuf:"bytes,1,opt,name=plugin,proto3" json:"plugin,omitempty"`
}

type CreateBotReq struct {
	UserId int64    `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Info   *BotInfo `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
}

type BotInfo struct {
	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Image       string `protobuf:"bytes,3,opt,name=image,proto3" json:"image,omitempty"`
	Type        string `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
}

type PublishBotReq struct {
	UserId  int64             `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Content *BotDetailContent `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

type BotDetailContent struct {
	Id                  int64                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	BotId               int64                `protobuf:"varint,2,opt,name=bot_id,json=botId,proto3" json:"bot_id,omitempty"`
	CreatorId           int64                `protobuf:"varint,3,opt,name=creator_id,json=creatorId,proto3" json:"creator_id,omitempty"`
	Prompt              string               `protobuf:"bytes,4,opt,name=prompt,proto3" json:"prompt,omitempty"`
	Plugins             []int64              `protobuf:"varint,5,rep,packed,name=plugins,proto3" json:"plugins,omitempty"`
	WelcomeMsg          string               `protobuf:"bytes,6,opt,name=welcome_msg,json=welcomeMsg,proto3" json:"welcome_msg,omitempty"`
	GuideInfo           []string             `protobuf:"bytes,7,rep,name=guide_info,json=guideInfo,proto3" json:"guide_info,omitempty"`
	DebugConversationId int64                `protobuf:"varint,8,opt,name=debug_conversation_id,json=debugConversationId,proto3" json:"debug_conversation_id,omitempty"`
	CreatedAt           *timestamp.Timestamp `protobuf:"bytes,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt           *timestamp.Timestamp `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	ModelSettings       *ModelSettings       `protobuf:"bytes,11,opt,name=model_settings,json=modelSettings,proto3" json:"model_settings,omitempty"`
}

type UpdateBotInfoReq struct {
	BotId  int64    `protobuf:"varint,1,opt,name=bot_id,json=botId,proto3" json:"bot_id,omitempty"`
	UserId int64    `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Info   *BotInfo `protobuf:"bytes,3,opt,name=info,proto3" json:"info,omitempty"`
}

type SaveBotDraftReq struct {
	Draft  *BotDetailContent `protobuf:"bytes,2,opt,name=draft,proto3" json:"draft,omitempty"`
	UserId int64             `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type GetBotDraftReq struct {
	BotId  int64 `protobuf:"varint,1,opt,name=bot_id,json=botId,proto3" json:"bot_id,omitempty"`
	UserId int64 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type ModelSettings struct {
	Model       string  `protobuf:"bytes,1,opt,name=model,proto3" json:"model,omitempty"`
	Temperature float32 `protobuf:"fixed32,2,opt,name=temperature,proto3" json:"temperature,omitempty"`
	TopP        float32 `protobuf:"fixed32,3,opt,name=top_p,json=topP,proto3" json:"top_p,omitempty"`
	Rounds      int32   `protobuf:"varint,4,opt,name=rounds,proto3" json:"rounds,omitempty"`
	MaxLength   int32   `protobuf:"varint,5,opt,name=max_length,json=maxLength,proto3" json:"max_length,omitempty"`
}
