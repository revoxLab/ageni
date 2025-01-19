package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/readonme/open-studio/common/fatal"
	"github.com/readonme/open-studio/common/httpserver"
	"github.com/readonme/open-studio/conf"
	"github.com/readonme/open-studio/router/middleware"
)

func StartHttp(config *conf.Config) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		defer fatal.RecoverFromPanic()
		c.Next()
	})
	r.Use(middleware.CORS())
	r.Use(middleware.PrePack())

	webGroup := r.Group("web")
	webGroup.POST("wallet_login", webWalletLogin)

	r = publicHttp(r)

	r = studioHttp(r)

	config.HttpServer.Handler = r
	httpserver.Run(config.HttpServer)
}

func publicHttp(r *gin.Engine) *gin.Engine {
	publicGroup := r.Group("")
	{
		publicGroup.POST("/v1/studio/bot_detail", studioBotDetail)
	}
	return r
}

func studioHttp(r *gin.Engine) *gin.Engine {
	studioGroup := r.Group("/v1/studio")
	{
		studioGroup.Use(middleware.AuthMiddleware())
		// /v1/wallet/agent/send/message
		studioGroup.POST("/bot_list", botList)
		//studioGroup.POST("/bot_detail", studioBotDetail)
		studioGroup.POST("/send_message", studioSendMessage)
		studioGroup.POST("/get_message", studioGetMessage)
		studioGroup.POST("/create_conversation", studioCreateConversation)
		studioGroup.POST("/conversation_history", studioConversationHistory)
		studioGroup.POST("/conversation_list", studioConversationList)
		studioGroup.POST("/conversation_delete", studioConversationDelete)
		studioGroup.POST("/plugin_list", studioPluginList)
		studioGroup.POST("/plugin_detail", studioPluginDetail)
		studioGroup.GET("/plugin_tabs", studioPluginTabs)
		studioGroup.GET("/bot_tabs", studioBotTabs)
		studioGroup.POST("/create_bot", createBot)
		studioGroup.POST("/update_bot_info", updateBotInfo)
		studioGroup.POST("/publish_bot", publishBot)
		studioGroup.POST("/save_bot_draft", saveBotDraft)
		studioGroup.POST("/get_bot_draft", getBotDraft)
		studioGroup.POST("/user_bot_list", userBotList)
	}
	return r
}
