package service

import (
	"context"
	"github.com/BurntSushi/toml"
	"github.com/readonme/open-studio/common/log"
	"github.com/readonme/open-studio/conf"
	"github.com/readonme/open-studio/dal"
	"github.com/readonme/open-studio/service/models"
	"testing"
)

func init() {
	confPath := "../../conf/local.toml"
	_, _ = toml.DecodeFile(confPath, &conf.Conf)

	dal.Init(conf.Conf.StudioDB)

}

func TestCreateBot(t *testing.T) {
	ctx := context.Background()

	req := &models.CreateBotReq{
		UserId: 123,
		Info: &models.BotInfo{
			Name:        "TestBot",
			Description: "A test bot",
			Image:       "http://example.com/image.jpg",
		},
	}

	resp, err := CreateBot(ctx, req)
	log.Infof("resp:%v, err:%v", resp, err)
}

func TestPublishBot(t *testing.T) {
	ctx := context.Background()

	req := &models.PublishBotReq{
		UserId: 123,
		Content: &models.BotDetailContent{
			BotId:   9,
			Prompt:  "This is a aaatest prompt",
			Plugins: []int64{1, 2},
		},
	}

	resp, err := PublishBot(ctx, req)
	log.Infof("resp:%v, err:%v", resp, err)

}

func TestUpdateBotInfo(t *testing.T) {
	ctx := context.Background()

	req := &models.UpdateBotInfoReq{
		BotId:  9,
		UserId: 123,
		Info: &models.BotInfo{
			Name:        "UpdatedBot",
			Description: "An updated test bot",
			Image:       "https://example.com/updated_image.jpg",
		},
	}

	resp, err := UpdateBotInfo(ctx, req)
	log.Infof("resp:%v, err:%v", resp, err)

}

func TestSaveBotDraft(t *testing.T) {
	ctx := context.Background()

	req := &models.SaveBotDraftReq{
		UserId: 123,
		Draft: &models.BotDetailContent{
			BotId:      9,
			Prompt:     "This is a draft prompt",
			WelcomeMsg: "draft welcome",
			Plugins:    []int64{5, 6},
		},
	}

	resp, err := SaveBotDraft(ctx, req)
	log.Infof("resp:%v, err:%v", resp, err)

}

func TestGetBotDraft(t *testing.T) {
	ctx := context.Background()

	req := &models.GetBotDraftReq{
		BotId:  1,
		UserId: 1524,
	}

	resp, err := GetBotDraft(ctx, req)
	log.Infof("resp:%v, err:%v", resp, err)
}

func TestUserBotList(t *testing.T) {
	ctx := context.Background()

	req := &models.UserBotListReq{
		UserId:  123,
		Page:    1,
		PerPage: 10,
	}

	resp, err := UserBotList(ctx, req)
	log.Infof("resp:%v, err:%v", resp, err)

}

func TestStudioServiceImpl_BotDetail(t *testing.T) {
	ctx := context.Background()

	req := &models.BotDetailReq{
		BotId: 10,
	}

	resp, err := BotDetail(ctx, req)
	log.Infof("resp:%v, err:%v", resp, err)
}
