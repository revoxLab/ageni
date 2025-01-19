package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/readonme/open-studio/common/log"
	"github.com/readonme/open-studio/common/response"
	"github.com/readonme/open-studio/conf"
	"github.com/readonme/open-studio/lib"
	"github.com/readonme/open-studio/service"
	"github.com/readonme/open-studio/service/models"
	user2 "github.com/readonme/open-studio/user"
	"time"
)

func webWalletLogin(c *gin.Context) {
	req := &user2.WalletLoginRequest{}
	if err := c.ShouldBind(req); err != nil {
		response.JSONFail(c, err, "")
		return
	}
	_, token, err := user2.WebWalletLogin(lib.RequestContext(c), req)
	if err != nil {
		response.JSONFail(c, err, "")
		return
	}
	response.AbortWithJSONSuccess(c, map[string]string{
		"token": token,
	})
}

func studioSendMessage(ctx *gin.Context) {
	req := &models.SendMessageReq{}
	if err := ctx.ShouldBind(req); err != nil {
		response.JSONFail(ctx, err, nil)
		return
	}
	if req.ConversationId == 0 {
		response.JSONFail(ctx, errors.New("params error"), nil)
		return
	}

	res, err := service.SendMessageService(context.Background(), req)
	if err != nil || res == nil {
		log.Errorf("lense sendMessage err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	response.JSONSuccess(ctx, res)
}

type GetStudioMessageReq struct {
	MessageId int64 `json:"message_id"`
}

type CreateConversationReq struct {
	BotId string `json:"bot_id"`
}

type CreateConversationResp struct {
	Id        string    `json:"id"`
	BotId     string    `json:"bot_id"`
	CreatedAt time.Time `json:"created_at"`
}

func studioCreateConversation(ctx *gin.Context) {
	req := &models.CreateConversationReq{}
	if err := ctx.ShouldBind(req); err != nil {
		response.JSONFail(ctx, err, nil)
		return
	}
	req.UserId = lib.GetContextUid(ctx)
	resp, err := service.CreateConversation(context.Background(), req)
	if err != nil || resp == nil {
		response.JSONFail(ctx, err, nil)
		return
	}
	response.JSONSuccess(ctx, resp)
}

type ConversationHistoryReq struct {
	ConversationId string `json:"conversation_id"`
}

type AgentMessage struct {
	Id         string    `json:"id"`
	SenderType string    `json:"sender_type"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type ConversationHistoryResp struct {
	ConversationId string          `json:"conversation_id"`
	Messages       []*AgentMessage `json:"messages"`
}

type BotListReq struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Tab      string `json:"tab"`
	Keywords string `json:"keywords"`
	PickType int32  `json:"pick_type"`
}

func botList(c *gin.Context) {
	req := new(BotListReq)
	if err := c.ShouldBind(req); err != nil {
		log.InfofWithContext(c, "Fail to parse request, err=%s", err)
		response.JSONFail(c, fmt.Errorf("fail to parse request, err=%s", err), nil)
		return
	}
	resp, err := service.BotListService(context.Background(), &models.BotListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Tab:      req.Tab,
		Keywords: req.Keywords,
		PickType: req.PickType,
	})
	if err != nil {
		response.JSONFail(c, err, nil)
		return
	}
	response.JSONSuccess(c, resp)
}

func studioGetMessage(ctx *gin.Context) {
	req := &GetStudioMessageReq{}
	if err := ctx.ShouldBind(req); err != nil {
		response.JSONFail(ctx, err, nil)
		return
	}

	resp, err := service.MessageDetail(context.Background(), &models.MessageDetailReq{MessageId: req.MessageId})
	if err != nil {
		response.JSONFail(ctx, err, nil)
		return
	}
	response.JSONSuccess(ctx, resp.Data)
}

func studioConversationHistory(ctx *gin.Context) {
	req := &models.ConversationHistoryReq{}
	if err := ctx.ShouldBind(req); err != nil {
		log.Errorf("studioConversationHistory err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	req.UserId = lib.GetContextUid(ctx)
	resp, err := service.ConversationHistoryService(context.Background(), req)
	if err != nil || resp == nil {
		log.Errorf("studioConversationHistory err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	response.JSONSuccess(ctx, resp)
}

func studioConversationDelete(ctx *gin.Context) {
	req := &models.ConversationDeleteReq{}
	if err := ctx.ShouldBind(req); err != nil {
		log.Errorf("studioConversationDelete err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	req.UserId = lib.GetContextUid(ctx)
	err := service.ConversationDelete(lib.RequestContext(ctx), req)
	if err != nil {
		log.Errorf("studioConversationDelete err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	response.JSONSuccess(ctx, nil)
}

func studioConversationList(ctx *gin.Context) {
	req := &models.ConversationListReq{}
	if err := ctx.ShouldBind(req); err != nil {
		log.Errorf("studioConversationList err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	req.UserId = lib.GetContextUid(ctx)
	resp, err := service.ConversationListService(lib.RequestContext(ctx), req)
	if err != nil || resp == nil {
		log.Errorf("studioConversationList err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	response.JSONSuccess(ctx, resp.Data)
}

func studioBotDetail(ctx *gin.Context) {
	req := &models.BotDetailReq{}
	if err := ctx.ShouldBind(req); err != nil {
		log.Errorf("studioBotDetail err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	resp, err := service.BotDetailService(lib.RequestContext(ctx), req)
	if err != nil || resp == nil {
		log.Errorf("studioBotDetail err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}

	response.JSONSuccess(ctx, resp.Bot)
}

func studioPluginList(ctx *gin.Context) {
	req := &models.PluginListReq{}
	if err := ctx.ShouldBind(req); err != nil {
		log.Errorf("studioPluginList err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	resp, err := service.PluginListService(lib.RequestContext(ctx), req)
	if err != nil || resp == nil {
		log.Errorf("studioPluginList err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	response.JSONSuccess(ctx, resp)
}

func studioPluginDetail(ctx *gin.Context) {
	req := &models.PluginDetailReq{}
	if err := ctx.ShouldBind(req); err != nil {
		log.Errorf("studioPluginDetail err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	resp, err := service.PluginDetailService(lib.RequestContext(ctx), req)
	if err != nil || resp == nil {
		log.Errorf("studioPluginDetail err:%v", err)
		response.JSONFail(ctx, err, nil)
		return
	}
	response.JSONSuccess(ctx, resp)
}

func studioPluginTabs(ctx *gin.Context) {
	response.JSONSuccess(ctx, conf.Conf.PluginTabs)
}

func studioBotTabs(ctx *gin.Context) {
	response.JSONSuccess(ctx, conf.Conf.BotTabs)
}

type CreateBotReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Type        string `json:"type"`
}

func createBot(c *gin.Context) {
	req := new(CreateBotReq)
	if err := c.ShouldBind(req); err != nil {
		log.InfofWithContext(c, "Fail to parse request, err=%s", err)
		response.JSONFail(c, fmt.Errorf("fail to parse request, err=%s", err), nil)
		return
	}
	userId := lib.GetContextUid(c)
	resp, err := service.CreateBot(context.Background(), &models.CreateBotReq{
		UserId: userId,
		Info: &models.BotInfo{
			Name:        req.Name,
			Description: req.Description,
			Image:       req.Image,
			Type:        req.Type,
		},
	})
	if err != nil {
		response.JSONFail(c, err, nil)
		return
	}
	response.JSONSuccess(c, resp)
}

type PublishBotReq struct {
	*BotDetailContent
}

func publishBot(c *gin.Context) {
	req := new(PublishBotReq)
	if err := c.ShouldBind(req); err != nil {
		log.InfofWithContext(c, "Fail to parse request, err=%s", err)
		response.JSONFail(c, fmt.Errorf("fail to parse request, err=%s", err), nil)
		return
	}
	userId := lib.GetContextUid(c)
	resp, err := service.PublishBot(context.Background(), &models.PublishBotReq{
		UserId:  userId,
		Content: convertToBotDetailContent(req.BotDetailContent),
	})
	if err != nil {
		response.JSONFail(c, err, nil)
		return
	}
	response.JSONSuccess(c, resp)
}

type UpdateBotInfoReq struct {
	BotID       int64  `json:"bot_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Type        string `json:"type"`
}

func updateBotInfo(c *gin.Context) {
	req := new(UpdateBotInfoReq)
	if err := c.ShouldBind(req); err != nil {
		log.InfofWithContext(c, "Fail to parse request, err=%s", err)
		response.JSONFail(c, fmt.Errorf("fail to parse request, err=%s", err), nil)
		return
	}
	userId := lib.GetContextUid(c)
	resp, err := service.UpdateBotInfo(context.Background(), &models.UpdateBotInfoReq{
		BotId:  req.BotID,
		UserId: userId,
		Info: &models.BotInfo{
			Name:        req.Name,
			Description: req.Description,
			Image:       req.Image,
			Type:        req.Type,
		},
	})
	if err != nil {
		response.JSONFail(c, err, nil)
		return
	}
	response.JSONSuccess(c, resp)
}

type SaveBotDraftReq struct {
	*BotDetailContent
}

type BotDetailContent struct {
	BotId               int64          `json:"bot_id"`
	Prompt              string         `json:"prompt"`
	Plugins             []int64        `json:"plugins"`
	WelcomeMsg          string         `json:"welcome_msg"`
	GuideInfo           []string       `json:"guide_info"`
	DebugConversationId int64          `json:"debug_conversation_id"`
	ModelSettings       *ModelSettings `json:"model_settings"`
	Tab                 string         `json:"tab"`
}

type ModelSettings struct {
	Model       string  `json:"model,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	TopP        float32 `json:"top_p,omitempty"`
	Rounds      int32   `json:"rounds,omitempty"`
	MaxLength   int32   `json:"max_length,omitempty"`
}

func saveBotDraft(c *gin.Context) {
	req := new(SaveBotDraftReq)
	if err := c.ShouldBind(req); err != nil {
		log.InfofWithContext(c, "Fail to parse request, err=%s", err)
		response.JSONFail(c, fmt.Errorf("fail to parse request, err=%s", err), nil)
		return
	}
	userId := lib.GetContextUid(c)
	resp, err := service.SaveBotDraft(context.Background(), &models.SaveBotDraftReq{
		Draft:  convertToBotDetailContent(req.BotDetailContent),
		UserId: userId,
	})
	if err != nil {
		response.JSONFail(c, err, nil)
		return
	}
	response.JSONSuccess(c, resp)
}

type GetBotDraftReq struct {
	BotID int64 `json:"bot_id"`
}

func getBotDraft(c *gin.Context) {
	req := new(GetBotDraftReq)
	if err := c.ShouldBind(req); err != nil {
		log.InfofWithContext(c, "Fail to parse request, err=%s", err)
		response.JSONFail(c, fmt.Errorf("fail to parse request, err=%s", err), nil)
		return
	}
	userId := lib.GetContextUid(c)
	resp, err := service.GetBotDraft(context.Background(), &models.GetBotDraftReq{
		BotId:  req.BotID,
		UserId: userId,
	})
	if err != nil {
		response.JSONFail(c, err, nil)
		return
	}
	response.JSONSuccess(c, map[string]interface{}{
		"draft": resp,
	})
}

type UserBotListReq struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
	UserID   int64 `json:"user_id"`
}

func userBotList(c *gin.Context) {
	req := new(UserBotListReq)
	if err := c.ShouldBind(req); err != nil {
		log.InfofWithContext(c, "Fail to parse request, err=%s", err)
		response.JSONFail(c, fmt.Errorf("fail to parse request, err=%s", err), nil)
		return
	}
	userId := lib.GetContextUid(c)
	resp, err := service.UserBotListService(context.Background(), &models.UserBotListReq{
		Page:    req.Page,
		PerPage: req.PageSize,
		UserId:  userId,
	})
	if err != nil {
		response.JSONFail(c, err, nil)
		return
	}
	response.JSONSuccess(c, resp)
}

func convertToBotDetailContent(content *BotDetailContent) *models.BotDetailContent {
	return &models.BotDetailContent{
		BotId:         content.BotId,
		Prompt:        content.Prompt,
		Plugins:       content.Plugins,
		WelcomeMsg:    content.WelcomeMsg,
		GuideInfo:     content.GuideInfo,
		ModelSettings: convertToModelSettings(content.ModelSettings),
	}
}

func convertToModelSettings(settings *ModelSettings) *models.ModelSettings {
	return &models.ModelSettings{
		Model:       settings.Model,
		Temperature: settings.Temperature,
		TopP:        settings.TopP,
		Rounds:      settings.Rounds,
		MaxLength:   settings.MaxLength,
	}
}
