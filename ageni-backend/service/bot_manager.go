package service

import (
	"context"
	"github.com/readonme/open-studio/dal"
	"github.com/readonme/open-studio/dal/query"
	"github.com/readonme/open-studio/service/bot"
	"github.com/readonme/open-studio/service/bot/config"
)

func GetBot(botId int64) (bot.Bot, error) {
	b := query.Use(dal.StudioDB).Bot
	botModel, err := b.WithContext(context.TODO()).Where(b.ID.Eq(botId)).First()
	if err != nil {
		return nil, err
	}
	return &config.ConfigBot{Bot: botModel}, nil
}
