package entity

import "github.com/readonme/open-studio/dal/model"

type BotAggregate struct {
	Bot     *model.Bot
	Plugins []*model.Plugin
}

type BotAggregateSlice []*BotAggregate

type PluginAggregate struct {
	Bots    []*model.Bot
	Methods []*model.Method
	Plugin  *model.Plugin
}

type PluginAggregateSlice []*PluginAggregate
