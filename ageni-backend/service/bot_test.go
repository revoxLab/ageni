package service

import (
	"context"
	"github.com/BurntSushi/toml"
	"github.com/readonme/open-studio/conf"
	"github.com/readonme/open-studio/dal"
	"github.com/readonme/open-studio/service/models"
	"testing"
)

func init() {
	confPath := "../conf/local.toml"
	if _, err := toml.DecodeFile(confPath, &conf.Conf); err != nil {
		panic(err)
	}
	dal.Init(conf.Conf.StudioDB)
}

func TestConversationCreate(t *testing.T) {
	ctx := context.Background()
	in := &models.CreateConversationReq{UserId: 179, BotId: 1, Title: "eamon test"}
	cid, err := ConversationCreate(ctx, in)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success cid:", cid)
}

func TestBotList(t *testing.T) {
	result, err := BotList(context.Background(), &models.BotListReq{Keywords: "cross", Tab: "test"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
