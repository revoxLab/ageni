package service

import (
	"context"
	"errors"
	xstudio "github.com/readonme/open-studio/common"
	"github.com/readonme/open-studio/common/fatal"
	"github.com/readonme/open-studio/common/log"
	"github.com/readonme/open-studio/dal"
	"github.com/readonme/open-studio/dal/model"
	"github.com/readonme/open-studio/dal/query"
	bot2 "github.com/readonme/open-studio/service/bot"
	"github.com/readonme/open-studio/service/models"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

const MessageStatusUnprocessed = 0
const MessageStatusSuccess = 1
const MessageStatusFailed = 2

func InsertMessage(msg *model.Message) (*model.Message, error) {
	db := query.Use(dal.StudioDB)

	err := db.Message.WithContext(context.TODO()).Create(msg)
	if err != nil {
		log.Errorf("CreateMessageErr err:%v", err)
		return nil, err
	}

	return msg, nil
}

func GetConversationMessage(conversationId int64) ([]*model.Message, error) {
	m := query.Use(dal.StudioDB).Message
	msgs, err := m.WithContext(context.TODO()).Where(m.ConversationID.Eq(conversationId), m.MessageType.Neq(openai.ChatMessageRoleSystem)).Order(m.ID.Asc()).Find()
	return msgs, err
}

func GetConversation(conversationId int64) (*model.Conversation, error) {
	c := query.Use(dal.StudioDB).Conversation
	conv, err := c.WithContext(context.TODO()).Where(c.ID.Eq(conversationId)).First()
	if err != nil {
		return nil, err
	}
	return conv, nil
}

func stringDifference(list1, list2 []string) []string {
	set := make(map[string]struct{})
	for _, item := range list2 {
		set[item] = struct{}{}
	}

	var result []string

	for _, item := range list1 {
		if _, exists := set[item]; !exists {
			result = append(result, item)
		}
	}

	return result
}
func checkAndFixConversations(botId int64, req *models.SendMessageReq, messages []*model.Message) []*model.Message {
	needToolCallIds := make([]string, 0, len(messages))
	haveToolCallIds := make([]string, 0, len(messages))
	for _, m := range messages {
		tcs := m.GetToolCalls()
		if len(tcs) != 0 {
			for _, t := range tcs {
				needToolCallIds = append(needToolCallIds, t.ID)
			}
		}
		if len(m.ToolCallID) != 0 {
			haveToolCallIds = append(haveToolCallIds, m.ToolCallID)
		}
	}
	diffToolCallIds := stringDifference(needToolCallIds, haveToolCallIds)
	if len(diffToolCallIds) == 0 {
		return messages
	}
	newMessages := make([]*model.Message, 0, len(diffToolCallIds))
	for _, d := range diffToolCallIds {

		newMessages = append(newMessages, &model.Message{
			BotID:          botId,
			UserID:         req.UserId,
			ConversationID: req.ConversationId,
			MessageType:    openai.ChatMessageRoleTool,
			Content:        "failed",
			ToolCallID:     d,
			Status:         MessageStatusSuccess,
		})
	}
	db := query.Use(dal.StudioDB)

	err := db.Message.WithContext(context.TODO()).Create(newMessages...)
	if err != nil {
		log.Errorf("Create miss ToolCall MessageErr err:%v", err)
		return messages
	}
	messages = append(messages, newMessages...)
	return messages
}

func SendMessage(ctx context.Context, req *models.SendMessageReq) (int64, int64, error) {
	conversation, e := GetConversation(req.ConversationId)
	botId := conversation.BotID
	if e != nil {
		log.Errorf("SendMessage GetConversation err=%v, id=%v", e, req.ConversationId)
		return 0, 0, e
	}
	messages, e := GetConversationMessage(req.ConversationId)
	if e != nil {
		log.Errorf("SendMessage GetConversation err=%v, id=%v", e, req.ConversationId)
		return 0, 0, e
	}

	for _, r := range req.ToolResults {
		msg := &model.Message{
			BotID:          botId,
			UserID:         req.UserId,
			ConversationID: req.ConversationId,
			MessageType:    openai.ChatMessageRoleTool,
			Content:        r.Content,
			ToolCallID:     r.ToolCallId,
			Status:         MessageStatusSuccess,
		}
		_, e = InsertMessage(msg)
		if e != nil {
			log.Errorf("InsertMessage err:%v", e)
			return botId, 0, e
		}
		messages = append(messages, msg)
	}
	if len(req.Content) != 0 {

		messages = checkAndFixConversations(botId, req, messages)

		msg := &model.Message{
			BotID:          botId,
			UserID:         req.UserId,
			ConversationID: req.ConversationId,
			MessageType:    openai.ChatMessageRoleUser,
			Content:        req.Content,
			Status:         MessageStatusSuccess,
		}
		_, e = InsertMessage(msg)
		if e != nil {
			log.Errorf("InsertMessage err:%v", e)
			return botId, 0, e
		}

		messages = append(messages, msg)
	}
	AIMessage, err := InsertMessage(&model.Message{
		BotID:          botId,
		UserID:         req.UserId,
		ConversationID: req.ConversationId,
		MessageType:    openai.ChatMessageRoleAssistant,
		Status:         MessageStatusUnprocessed,
	})
	if err != nil {
		log.Errorf("InsertMessage err:%v", e)
		return botId, 0, e
	}
	// Start asynchronous processing
	go func() {
		defer fatal.RecoverFromPanic()
		bot, err := GetBot(botId)
		if err != nil {
			log.Errorf("GetBot err:%v", err)
			updateMessageStatus(context.TODO(), AIMessage.ID, MessageStatusFailed)
			return
		}
		openAIMsgs := make([]openai.ChatCompletionMessage, len(messages))
		for i, m := range messages {
			om := openai.ChatCompletionMessage{}
			om.Role = m.MessageType
			om.Content = m.Content
			om.ToolCallID = m.ToolCallID
			om.ToolCalls = m.GetAllToolCalls()
			openAIMsgs[i] = om
		}
		processMsg, botErr := bot.ProcessConversation(&bot2.ProcessContext{
			Conversation: conversation,
			Messages:     openAIMsgs,
		})
		if botErr != nil {
			log.Errorf("ProcessConversation err:%v conversationId:%v", botErr, req.ConversationId)
			updateMessageStatus(context.TODO(), AIMessage.ID, MessageStatusFailed)
			return
		}
		AIMessage.Content = processMsg.Content
		AIMessage.ToolCall = processMsg.ToolCallStr()
		AIMessage.BackendToolCall = processMsg.BackendToolCallStr()
		NextMessageIds := make([]int64, 0)
		for _, r := range processMsg.BackendToolResult {
			msg := &model.Message{
				BotID:          botId,
				UserID:         req.UserId,
				ConversationID: req.ConversationId,
				MessageType:    openai.ChatMessageRoleTool,
				Content:        r.Content,
				ToolCallID:     r.ToolCallID,
				Status:         MessageStatusSuccess,
			}
			m, e := InsertMessage(msg)
			NextMessageIds = append(NextMessageIds, m.ID)
			if e != nil {
				log.Errorf("backend tool InsertMessage err:%v", e)
				updateMessageStatus(context.TODO(), AIMessage.ID, MessageStatusFailed)
				return
			}
		}
		db := query.Use(dal.StudioDB).Message

		if len(processMsg.BackendToolResult) != 0 && len(processMsg.ToolCalls) == 0 {
			openAIMsgs = append(openAIMsgs, processMsg.OriginMessage)
			openAIMsgs = append(openAIMsgs, processMsg.BackendToolResult...)
			finalMessage, botErr := bot.ProcessConversation(&bot2.ProcessContext{
				Conversation: conversation,
				Messages:     openAIMsgs,
			})
			if botErr != nil {
				log.Errorf("ProcessConversation err:%v conversationId:%v", botErr, req.ConversationId)
				updateMessageStatus(context.TODO(), AIMessage.ID, MessageStatusFailed)
				return
			}
			Message, err := InsertMessage(&model.Message{
				BotID:          botId,
				UserID:         req.UserId,
				ConversationID: req.ConversationId,
				MessageType:    openai.ChatMessageRoleAssistant,
				Content:        finalMessage.Content,
				Status:         MessageStatusSuccess,
			})
			if err != nil {
				log.Errorf("InsertMessage err:%v", e)
				updateMessageStatus(context.TODO(), AIMessage.ID, MessageStatusFailed)
				return
			}
			NextMessageIds = append(NextMessageIds, Message.ID)
		}

		AIMessage.Status = MessageStatusSuccess
		AIMessage.SetNextMsgIds(NextMessageIds)
		_, err = db.WithContext(context.TODO()).Updates(AIMessage)
		if err != nil {
			log.Errorf("Failed to update Message status: %v", err)
			return
		}
	}()

	// Return the Message ID immediately
	return botId, AIMessage.ID, nil
}

// updateMessageStatus updates the status of a Message
func updateMessageStatus(ctx context.Context, MessageID int64, status int32) {
	db := query.Use(dal.StudioDB)
	_, err := db.Message.WithContext(ctx).Where(db.Message.ID.Eq(MessageID)).Update(db.Message.Status, status)
	if err != nil {
		log.Errorf("Failed to update Message status: %v", err)
	}
}

func MessageInfo(messageId int64) (*model.Message, error) {
	q := query.Use(dal.StudioDB).Message
	n, e := q.WithContext(context.Background()).Where(q.ID.Eq(messageId)).First()
	if e != nil && errors.Is(e, gorm.ErrRecordNotFound) {
		return nil, e
	}
	if n.Status == MessageStatusUnprocessed {
		return n, xstudio.ErrGenerating
	}
	if n.Status == MessageStatusFailed {
		return n, xstudio.ErrGenerateFailed
	}
	return n, nil
}
