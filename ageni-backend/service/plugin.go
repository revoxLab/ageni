package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/readonme/open-studio/dal"
	"github.com/readonme/open-studio/dal/model"
	"github.com/readonme/open-studio/dal/query"
	entity2 "github.com/readonme/open-studio/service/entity"
	"gorm.io/gorm"
)

func PluginDetail(id int64) (*entity2.PluginAggregate, error) {
	p := query.Use(dal.StudioDB).Plugin
	plugin, err := p.WithContext(context.Background()).Where(p.ID.Eq(id)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	m := query.Use(dal.StudioDB).Method
	methods, err := m.WithContext(context.Background()).Where(m.PluginID.Eq(id)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	bots, err := GetBotListByPluginId(plugin.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	out := &entity2.PluginAggregate{
		Bots:    bots,
		Plugin:  plugin,
		Methods: methods,
	}
	return out, nil
}

func PluginList(tab, keywords string, page, pageSize int32, pluginIds []int64) (entity2.PluginAggregateSlice, error) {
	if pageSize == 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	p := query.Use(dal.StudioDB).Plugin
	x := p.WithContext(context.Background())
	if len(tab) == 0 || tab == "Most used" {
		x = x.Order(p.Bots.Desc())
	} else {
		x = x.Where(p.Tab.Eq(tab)).Order(p.ID.Desc())
	}

	if len(pluginIds) > 0 {
		x = x.Where(p.ID.In(pluginIds...))
	} else {
		x = x.Limit(int(pageSize)).Offset(int(offset))
	}
	if len(keywords) > 0 {
		x = x.Where(p.Name.Like(fmt.Sprintf("%s%s%s", "%", keywords, "%")))
	}
	pluginList, err := x.Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	pIds := make([]int64, 0, len(pluginList))
	for _, item := range pluginList {
		pIds = append(pIds, item.ID)
	}
	botMap, err := GetBotListByPluginIds(pIds)
	if err != nil {
		return nil, err
	}
	methodMap, err := GetMethodMapByPluginIds(pIds)
	if err != nil {
		return nil, err
	}
	out := make(entity2.PluginAggregateSlice, 0, len(pluginList))
	for _, item := range pluginList {
		tmp := &entity2.PluginAggregate{
			Plugin: item,
		}
		if _v, ok := botMap[item.ID]; ok {
			tmp.Bots = _v
		}
		if _v, ok := methodMap[item.ID]; ok {
			tmp.Methods = _v
		}
		out = append(out, tmp)
	}
	return out, err
}

func GetBotListByPluginId(pluginId int64) ([]*model.Bot, error) {
	b := query.Use(dal.StudioDB).BotPlugin
	res, err := b.WithContext(context.Background()).Where(b.PluginID.Eq(pluginId)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	bIds := make([]int64, 0, len(res))
	for _, item := range res {
		bIds = append(bIds, item.BotID)
	}
	b2 := query.Use(dal.StudioDB).Bot
	list, err := b2.WithContext(context.Background()).Where(b2.ID.In(bIds...), b2.Status.Eq(BotPublish)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return list, nil
}

func GetPluginIdsByBotId(botId int64) ([]int64, error) {
	b := query.Use(dal.StudioDB).BotPlugin
	res, err := b.WithContext(context.Background()).Where(b.BotID.Eq(botId)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	pIds := make([]int64, 0, len(res))
	for _, item := range res {
		pIds = append(pIds, item.PluginID)
	}
	return pIds, nil
}

func GetPluginListByBotId(botId int64) ([]*model.Plugin, error) {
	b := query.Use(dal.StudioDB).BotPlugin
	res, err := b.WithContext(context.Background()).Where(b.BotID.Eq(botId)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	pIds := make([]int64, 0, len(res))
	for _, item := range res {
		pIds = append(pIds, item.PluginID)
	}
	p := query.Use(dal.StudioDB).Plugin
	list, err := p.WithContext(context.Background()).Where(p.ID.In(pIds...)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return list, nil
}

func GetPluginListByBotIds(botList []*model.Bot) (map[int64][]*model.Plugin, error) {
	botIds := make([]int64, 0, len(botList))
	for _, item := range botList {
		botIds = append(botIds, item.ID)
	}
	b := query.Use(dal.StudioDB).BotPlugin
	res, err := b.WithContext(context.Background()).Where(b.BotID.In(botIds...)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	pIds := make([]int64, 0, len(res))
	for _, item := range res {
		pIds = append(pIds, item.PluginID)
	}
	p := query.Use(dal.StudioDB).Plugin
	list, err := p.WithContext(context.Background()).Where(p.ID.In(pIds...)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	pluginIdMap := map[int64]*model.Plugin{}
	for _, item := range list {
		pluginIdMap[item.ID] = item
	}
	out := make(map[int64][]*model.Plugin)
	for _, item := range res {
		if out[item.BotID] == nil {
			out[item.BotID] = make([]*model.Plugin, 0)
		}
		if _v, ok := pluginIdMap[item.PluginID]; ok {
			out[item.BotID] = append(out[item.BotID], _v)
		}
	}
	return out, nil
}

func GetBotListByPluginIds(pluginIds []int64) (map[int64][]*model.Bot, error) {
	b := query.Use(dal.StudioDB).BotPlugin
	res, err := b.WithContext(context.Background()).Where(b.PluginID.In(pluginIds...)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	bIds := make([]int64, 0, len(res))
	for _, item := range res {
		bIds = append(bIds, item.BotID)
	}
	p := query.Use(dal.StudioDB).Bot
	list, err := p.WithContext(context.Background()).Where(p.ID.In(bIds...), p.Status.Eq(BotPublish)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	botIdMap := map[int64]*model.Bot{}
	for _, item := range list {
		botIdMap[item.ID] = item
	}
	out := make(map[int64][]*model.Bot)
	for _, item := range res {
		if out[item.PluginID] == nil {
			out[item.PluginID] = make([]*model.Bot, 0)
		}
		if _v, ok := botIdMap[item.BotID]; ok {
			out[item.PluginID] = append(out[item.PluginID], _v)
		}
	}
	return out, nil
}

func GetMethodMapByPluginIds(pluginIds []int64) (map[int64][]*model.Method, error) {
	m := query.Use(dal.StudioDB).Method
	methods, err := m.WithContext(context.Background()).Where(m.PluginID.In(pluginIds...)).Find()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	out := make(map[int64][]*model.Method)
	for _, item := range methods {
		if out[item.PluginID] == nil {
			out[item.PluginID] = make([]*model.Method, 0)
		}
		out[item.PluginID] = append(out[item.PluginID], item)
	}
	return out, nil
}
