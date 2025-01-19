package service

import (
	"context"
	"errors"
	xecode "github.com/readonme/open-studio/common"
	"github.com/readonme/open-studio/dal"
	"github.com/readonme/open-studio/dal/model"
	"github.com/readonme/open-studio/dal/query"
	"github.com/readonme/open-studio/service/bot/config"
	"github.com/readonme/open-studio/service/models"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"unicode/utf8"
)

func ConversationCreate(ctx context.Context, in *models.CreateConversationReq) (int64, error) {
	if in.BotId <= 0 || len(in.Title) == 0 {
		return 0, xecode.ErrParamCheck
	}
	bot, err := GetBotById(in.BotId)
	if err != nil {
		return 0, err
	}
	t := query.Use(dal.StudioDB)
	var conversationType = int32(1)
	if in.Type == "draft" {
		_, err = CheckBotUser(t, in.BotId, in.UserId)
		if err != nil {
			return 0, err
		}
		conversationType = 2
	}

	title := in.Title
	if utf8.RuneCountInString(in.Title) > 25 {
		title = string([]rune(in.Title)[:25])
	}
	cm := &model.Conversation{
		BotID:  in.BotId,
		Title:  title,
		UserID: in.UserId,
		Type:   conversationType,
	}
	err = t.Transaction(func(tx *query.Query) error {
		if err = tx.Conversation.WithContext(context.Background()).Create(cm); err != nil {
			return err
		}
		mm := &model.Message{
			BotID:          in.BotId,
			UserID:         in.UserId,
			ConversationID: cm.ID,
			MessageType:    openai.ChatMessageRoleSystem,
			Content:        bot.GetWelcomeMsg(),
			Status:         MessageStatusSuccess,
		}
		if err = tx.Message.WithContext(context.Background()).Create(mm); err != nil {
			return err
		}
		if in.Type == "draft" {
			_, err = tx.BotDraft.WithContext(context.Background()).Where(tx.BotDraft.BotID.Eq(in.BotId)).Update(tx.BotDraft.DebugConversationID, cm.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return cm.ID, err
}

func ConversationHistory(ctx context.Context, in *models.ConversationHistoryReq) ([]*model.Message, error) {
	if in.PageSize == 0 {
		in.PageSize = 10
	}
	offset := (in.Page - 1) * in.PageSize
	c := query.Use(dal.StudioDB).Message
	list, err := c.WithContext(ctx).Where(c.ConversationID.Eq(in.ConversationId)).
		Limit(int(in.PageSize)).
		Order(c.ID.Desc()).
		Offset(int(offset)).
		Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return list, err
	}
	return list, err
}

func ConversationList(ctx context.Context, in *models.ConversationListReq) ([]*model.Conversation, error) {
	if in.PageSize == 0 {
		in.PageSize = 10
	}
	offset := (in.Page - 1) * in.PageSize
	c := query.Use(dal.StudioDB).Conversation
	list, err := c.WithContext(ctx).Where(c.BotID.Eq(in.BotId), c.UserID.Eq(in.UserId), c.Status.Eq(1), c.Type.Neq(config.ConversationDebug)).
		Limit(int(in.PageSize)).
		Order(c.ID.Desc()).
		Offset(int(offset)).
		Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return list, err
	}
	return list, err
}

func GetBotById(botId int64) (*model.Bot, error) {
	c := query.Use(dal.StudioDB).Bot
	return c.WithContext(context.Background()).Where(c.ID.Eq(botId)).First()
}

func ConversationDelete(ctx context.Context, in *models.ConversationDeleteReq) error {
	c := query.Use(dal.StudioDB).Conversation
	_, err := c.WithContext(ctx).Where(c.ID.Eq(in.Id), c.UserID.Eq(in.UserId)).Update(c.Status, 2)
	return err
}
