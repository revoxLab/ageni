package convert_result

import (
	"encoding/json"
	"github.com/readonme/open-studio/dal/model"
	"github.com/readonme/open-studio/service/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MessageListResult(in []*model.Message) ([]*models.AgentMessage, error) {
	out := make([]*models.AgentMessage, 0, len(in))
	for _, item := range in {
		pb, err := MessageResult(item)
		if err != nil {
			return nil, err
		}
		out = append(out, pb)
	}
	return out, nil
}

func TooCallStrResult(in string) ([]*models.ToolCall, error) {
	if len(in) == 0 {
		return nil, nil
	}
	var out []*models.ToolCall
	err := json.Unmarshal([]byte(in), &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func MessageResult(in *model.Message) (*models.AgentMessage, error) {
	toolCalls, err := TooCallStrResult(in.ToolCall)
	if err != nil {
		return nil, err
	}
	return &models.AgentMessage{
		Id:             in.ID,
		ConversationId: in.ConversationID,
		Content:        in.Content,
		MessageType:    in.MessageType,
		ToolCallId:     in.ToolCallID,
		ToolCalls:      toolCalls,
		CreatedAt:      timestamppb.New(in.CreatedAt),
		Status:         in.Status,
		NextMessageIds: in.GetNextMsgIds(),
	}, nil
}

func ConversationListResult(in []*model.Conversation) []*models.Conversation {
	out := make([]*models.Conversation, 0, len(in))
	for _, item := range in {
		out = append(out, ConversationResult(item))
	}
	return out
}

func ConversationResult(in *model.Conversation) *models.Conversation {
	return &models.Conversation{
		Id:    in.ID,
		Title: in.Title,
	}
}
