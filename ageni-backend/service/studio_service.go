package service

import (
	"context"
	"github.com/readonme/open-studio/common/log"
	"github.com/readonme/open-studio/service/convert_result"
	"github.com/readonme/open-studio/service/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func BotListService(ctx context.Context, in *models.BotListReq) (*models.BotListResp, error) {
	list, err := BotList(ctx, in)
	if err != nil {
		return nil, err
	}
	m, e := convert_result.BotListResult(list)
	if e != nil {
		return nil, e
	}
	return &models.BotListResp{
		Bots: m,
	}, nil
}

func UserBotListService(ctx context.Context, in *models.UserBotListReq) (*models.BotListResp, error) {
	list, err := UserBotList(ctx, in)
	if err != nil {
		return nil, err
	}
	m, e := convert_result.BotListResult(list)
	if e != nil {
		return nil, e
	}
	return &models.BotListResp{
		Bots: m,
	}, nil
}

func CreateConversation(ctx context.Context, in *models.CreateConversationReq) (*models.CreateConversationResp, error) {
	cid, err := ConversationCreate(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &models.CreateConversationResp{
		BotId:          in.BotId,
		ConversationId: cid,
		CreatedAt:      timestamppb.New(time.Now()),
	}
	return resp, nil
}

func ConversationHistoryService(ctx context.Context, in *models.ConversationHistoryReq) (*models.ConversationHistoryResp, error) {
	list, err := ConversationHistory(ctx, in)
	if err != nil {
		return nil, err
	}
	m, e := convert_result.MessageListResult(list)
	if e != nil {
		return nil, e
	}
	return &models.ConversationHistoryResp{
		ConversationId: in.ConversationId,
		Messages:       m,
	}, nil
}

func MessageDetail(ctx context.Context, in *models.MessageDetailReq) (*models.MessageDetailResp, error) {
	resp := &models.MessageDetailResp{}
	node, err := MessageInfo(in.MessageId)
	if err != nil {
		return nil, err
	}
	resp.Data, err = convert_result.MessageResult(node)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SendMessageService(ctx context.Context, in *models.SendMessageReq) (*models.SendMessageResp, error) {
	_, id, err := SendMessage(ctx, in)
	if err != nil {
		log.Errorf("SendMessage err. c_id=%v, content=%v, err=%v", in.ConversationId, in.Content, err)
		return nil, err
	}
	//service.SetBotConversations(botId)
	return &models.SendMessageResp{
		RespMessageId: id,
	}, nil
}

func ConversationListService(ctx context.Context, in *models.ConversationListReq) (*models.ConversationListResp, error) {
	list, err := ConversationList(ctx, in)
	if err != nil {
		return nil, err
	}
	return &models.ConversationListResp{
		Data: convert_result.ConversationListResult(list),
	}, nil
}

func BotDetailService(ctx context.Context, in *models.BotDetailReq) (*models.BotDetailResp, error) {
	bot, err := BotDetail(ctx, in)
	if err != nil {
		return nil, err
	}
	return &models.BotDetailResp{
		Bot: convert_result.BotResult(bot.Bot, bot.Plugins, true),
	}, nil
}

func PluginListService(ctx context.Context, in *models.PluginListReq) (*models.PluginListResp, error) {
	plugins, err := PluginList(in.Tab, in.Keywords, in.Page, in.PageSize, in.Ids)
	if err != nil {
		return nil, err
	}
	ps, err := convert_result.PluginListResult(plugins)
	if err != nil {
		return nil, err
	}
	return &models.PluginListResp{
		List: ps,
	}, nil
}

func PluginDetailService(ctx context.Context, in *models.PluginDetailReq) (*models.PluginDetailResp, error) {
	plugin, err := PluginDetail(in.Id)
	if err != nil {
		return nil, err
	}
	pluginPb, err := convert_result.PluginResult(plugin.Plugin, plugin.Bots, plugin.Methods)
	if err != nil {
		return nil, err
	}
	return &models.PluginDetailResp{
		Plugin: pluginPb,
	}, nil
}
