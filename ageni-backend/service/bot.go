package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/readonme/open-studio/dal"
	"github.com/readonme/open-studio/dal/model"
	"github.com/readonme/open-studio/dal/query"
	"github.com/readonme/open-studio/service/entity"
	"github.com/readonme/open-studio/service/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const BotCreate = 1
const BotPublish = 2
const BotDelete = 3

const DraftNormal = 0
const DraftDelete = 1

func GetBotConversations(botId int64) int32 {
	c := query.Use(dal.StudioDB).Conversation
	count, _ := c.WithContext(context.TODO()).Where(c.BotID.Eq(botId)).Count()
	return int32(count)
}

func GetBotUsers(botId int64) int32 {
	c := query.Use(dal.StudioDB).Conversation
	count, _ := c.WithContext(context.TODO()).Where(c.BotID.Eq(botId)).Distinct(c.UserID).Count()
	return int32(count)
}
func BotDetail(ctx context.Context, in *models.BotDetailReq) (*entity.BotAggregate, error) {
	b := query.Use(dal.StudioDB).Bot
	bot, err := b.WithContext(ctx).Where(b.ID.Eq(in.BotId)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	bot.Users = GetBotUsers(bot.ID)
	bot.Conversations = GetBotConversations(bot.ID)
	out := &entity.BotAggregate{
		Bot: bot,
	}
	pluginList, err := GetPluginListByBotId(bot.ID)
	if err != nil {
		return nil, err
	}
	out.Plugins = pluginList
	return out, nil
}

func BotList(ctx context.Context, in *models.BotListReq) (entity.BotAggregateSlice, error) {
	if in.PageSize == 0 {
		in.PageSize = 10
	}
	offset := (in.Page - 1) * in.PageSize
	c := query.Use(dal.StudioDB).Bot
	q := c.WithContext(ctx)
	if len(in.Tab) == 0 || in.Tab == "Most used" {
		q = q.Order(c.Users.Desc())
	} else {
		q = q.Where(c.Tab.Eq(in.Tab)).Order(c.ID.Desc())
	}
	if len(in.Keywords) > 0 {
		q = q.Where(c.Name.Like(fmt.Sprintf("%s%s%s", "%", in.Keywords, "%")))
	}
	if in.PickType == 1 {
		q = q.Where(c.PickType.Eq(in.PickType))
	}
	botList, err := q.Where(c.Status.Eq(BotPublish)).Limit(int(in.PageSize)).Offset(int(offset)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	out := make(entity.BotAggregateSlice, 0, len(botList))
	pluginMap, err := GetPluginListByBotIds(botList)
	for _, item := range botList {
		item.Users = GetBotUsers(item.ID)
		item.Conversations = GetBotConversations(item.ID)
		tmp := &entity.BotAggregate{Bot: item}
		if _v, ok := pluginMap[item.ID]; ok {
			tmp.Plugins = _v
		}
		out = append(out, tmp)
	}
	return out, nil
}

func UserBotList(ctx context.Context, in *models.UserBotListReq) (entity.BotAggregateSlice, error) {
	if in.PerPage == 0 {
		in.PerPage = 10
	}
	offset := (in.Page - 1) * in.PerPage
	c := query.Use(dal.StudioDB).Bot
	q := c.WithContext(ctx)
	q = q.Order(c.UpdatedAt.Desc())
	botList, err := q.Where(c.CreatorID.Eq(in.UserId)).Limit(int(in.PerPage)).Offset(int(offset)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	out := make(entity.BotAggregateSlice, 0, len(botList))
	pluginMap, err := GetPluginListByBotIds(botList)
	for _, item := range botList {
		item.Users = GetBotUsers(item.ID)
		item.Conversations = GetBotConversations(item.ID)
		tmp := &entity.BotAggregate{Bot: item}
		if _v, ok := pluginMap[item.ID]; ok {
			tmp.Plugins = _v
		}
		out = append(out, tmp)
	}
	return out, nil
}

func CreateBot(ctx context.Context, in *models.CreateBotReq) (*model.Bot, error) {
	c := query.Use(dal.StudioDB).Bot
	info := in.Info
	bot := &model.Bot{
		Status:      BotCreate,
		Name:        info.Name,
		Description: info.Description,
		Image:       info.Image,
		CreatorID:   in.UserId,
		Config:      "{}",
		Tab:         info.Type,
	}
	err := c.WithContext(ctx).Create(bot)
	if err != nil {
		return nil, err
	}
	return bot, nil
}
func PublishBot(ctx context.Context, in *models.PublishBotReq) (*model.Bot, error) {
	c := query.Use(dal.StudioDB)

	botContent := in.Content
	err := c.Transaction(func(tx *query.Query) error {
		_, err := CheckBotUser(tx, botContent.BotId, in.UserId)
		if err != nil {
			return err
		}

		config := map[string]interface{}{
			"prompt":         botContent.Prompt,
			"plugin_ids":     botContent.Plugins,
			"welcome_msg":    botContent.WelcomeMsg,
			"guide_info":     botContent.GuideInfo,
			"model_settings": botContent.ModelSettings,
		}
		configJSON, err := json.Marshal(config)
		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"config": string(configJSON),
			"status": BotPublish,
		}

		_, err = tx.Bot.WithContext(context.TODO()).Where(tx.Bot.ID.Eq(botContent.BotId)).Updates(updates)
		if err != nil {
			return err
		}

		_, err = tx.BotPlugin.WithContext(context.TODO()).Where(tx.BotPlugin.BotID.Eq(botContent.BotId)).Delete()
		if err != nil {
			return err
		}

		for _, pluginID := range botContent.Plugins {
			err = tx.BotPlugin.WithContext(context.TODO()).Create(&model.BotPlugin{
				BotID:    botContent.BotId,
				PluginID: pluginID,
			})
			if err != nil {
				return err
			}
		}

		_, err = tx.BotDraft.WithContext(context.TODO()).Where(tx.BotDraft.BotID.Eq(botContent.BotId)).Update(tx.BotDraft.Status, DraftDelete)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	updatedBot, err := c.Bot.WithContext(context.TODO()).Where(c.Bot.ID.Eq(botContent.BotId)).First()
	if err != nil {
		return nil, err
	}

	return updatedBot, nil
}

func UpdateBotInfo(ctx context.Context, in *models.UpdateBotInfoReq) (*model.Bot, error) {
	c := query.Use(dal.StudioDB)

	var updatedBot *model.Bot
	info := in.Info
	err := c.Transaction(func(tx *query.Query) error {
		bot, err := CheckBotUser(tx, in.BotId, in.UserId)
		if err != nil {
			return err
		}

		bot.Name = info.Name
		bot.Description = info.Description
		bot.Image = info.Image
		bot.Tab = info.Type

		_, err = tx.Bot.WithContext(context.TODO()).Where(tx.Bot.ID.Eq(in.BotId)).Updates(bot)
		if err != nil {
			return err
		}

		updatedBot = bot
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedBot, nil
}

func CheckBotUser(tx *query.Query, botId, userId int64) (*model.Bot, error) {
	bot, err := tx.Bot.WithContext(context.TODO()).Where(tx.Bot.ID.Eq(botId)).First()
	if err != nil {
		return nil, err
	}
	if bot.CreatorID != userId {
		return nil, errors.New("permission denied")
	}
	return bot, nil
}

func SaveBotDraft(ctx context.Context, in *models.SaveBotDraftReq) (int64, error) {
	c := query.Use(dal.StudioDB)

	var draftID int64
	draft := in.Draft
	err := c.Transaction(func(tx *query.Query) error {
		_, err := CheckBotUser(tx, draft.BotId, in.UserId)
		if err != nil {
			return err
		}
		plugins, err := json.Marshal(draft.Plugins)
		if err != nil {
			return err
		}

		GuideInfo, err := json.Marshal(draft.GuideInfo)
		if err != nil {
			return err
		}
		var modelSettings []byte
		if draft.ModelSettings != nil {
			modelSettings, err = json.Marshal(draft.ModelSettings)
			if err != nil {
				return err
			}
		}

		existingDraft, err := tx.WithContext(context.TODO()).BotDraft.Where(tx.BotDraft.BotID.Eq(draft.BotId), tx.BotDraft.Status.Eq(DraftNormal)).First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if existingDraft != nil {
			existingDraft.Prompt = draft.Prompt
			existingDraft.Plugins = string(plugins)
			existingDraft.WelcomeMsg = draft.WelcomeMsg
			existingDraft.GuideInfo = string(GuideInfo)
			existingDraft.ModelSettings = string(modelSettings)

			_, err = tx.WithContext(context.TODO()).BotDraft.Updates(existingDraft)
			if err != nil {
				return err
			}

			draftID = existingDraft.ID
		} else {
			newDraft := &model.BotDraft{
				BotID:         draft.BotId,
				CreatorID:     in.UserId,
				Prompt:        draft.Prompt,
				Plugins:       string(plugins),
				WelcomeMsg:    draft.WelcomeMsg,
				GuideInfo:     string(GuideInfo),
				ModelSettings: string(modelSettings),
			}

			err = tx.WithContext(context.TODO()).BotDraft.Create(newDraft)
			if err != nil {
				return err
			}

			draftID = newDraft.ID
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return draftID, nil
}

func getOnlineBotContent(botId int64) (*models.BotDetailContent, error) {
	c := query.Use(dal.StudioDB).Bot
	bot, err := c.WithContext(context.TODO()).Where(c.ID.Eq(botId)).First()
	if err != nil {
		return nil, err
	}
	pluginIds, err := GetPluginIdsByBotId(bot.ID)
	if err != nil {
		return nil, err
	}
	content := &models.BotDetailContent{
		BotId:         bot.ID,
		CreatorId:     bot.CreatorID,
		Prompt:        bot.GetPrompt(),
		Plugins:       pluginIds,
		WelcomeMsg:    bot.GetWelcomeMsg(),
		GuideInfo:     bot.GetGuideInfo(),
		ModelSettings: bot.GetModelSettings(),
		CreatedAt:     timestamppb.New(bot.CreatedAt),
		UpdatedAt:     timestamppb.New(bot.CreatedAt),
	}
	return content, nil
}

func GetBotDraft(ctx context.Context, in *models.GetBotDraftReq) (*models.BotDetailContent, error) {
	c := query.Use(dal.StudioDB).BotDraft

	draft, err := c.WithContext(ctx).Where(c.BotID.Eq(in.BotId), c.Status.Eq(DraftNormal)).Last()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			botDetail, err := getOnlineBotContent(in.BotId)
			if err != nil {
				return nil, err
			}
			return botDetail, nil
		}
		return nil, status.Errorf(codes.Internal, "Failed to fetch bot draft: %v", err)
	}

	var plugins []int64
	var GuideInfo []string
	var ModelSettings *models.ModelSettings

	if len(draft.Plugins) != 0 {
		if err := json.Unmarshal([]byte(draft.Plugins), &plugins); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to unmarshal plugins: %v", err)
		}
	}
	if len(draft.GuideInfo) != 0 {
		if err := json.Unmarshal([]byte(draft.GuideInfo), &GuideInfo); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to unmarshal nav questions: %v", err)
		}
	}

	if len(draft.ModelSettings) != 0 {
		if err := json.Unmarshal([]byte(draft.ModelSettings), &ModelSettings); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to unmarshal nav questions: %v", err)
		}
	}

	resp := &models.BotDetailContent{
		Id:                  draft.ID,
		BotId:               draft.BotID,
		CreatorId:           draft.CreatorID,
		Prompt:              draft.Prompt,
		Plugins:             plugins,
		WelcomeMsg:          draft.WelcomeMsg,
		GuideInfo:           GuideInfo,
		ModelSettings:       ModelSettings,
		DebugConversationId: draft.DebugConversationID,
		CreatedAt:           timestamppb.New(draft.CreatedAt),
		UpdatedAt:           timestamppb.New(draft.UpdatedAt),
	}
	return resp, nil
}
