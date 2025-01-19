package convert_result

import (
	"github.com/readonme/open-studio/common/log"
	"github.com/readonme/open-studio/dal/model"
	"github.com/readonme/open-studio/service/entity"
	"github.com/readonme/open-studio/service/models"
	user2 "github.com/readonme/open-studio/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BotListResult(in []*entity.BotAggregate) ([]*models.AgentBot, error) {
	creatorIDs := make([]int64, 0, len(in))
	for _, item := range in {
		if item.Bot.CreatorName == "" || item.Bot.CreatorPic == "" {
			creatorIDs = append(creatorIDs, int64(item.Bot.CreatorID))
		}
	}

	userMap := make(map[int64]*model.User)

	if len(creatorIDs) > 0 {
		userResp := user2.GetUserByIds(creatorIDs)
		for _, u := range userResp {
			userMap[u.ID] = u
		}
	}

	out := make([]*models.AgentBot, 0, len(in))
	for _, item := range in {
		bot := item.Bot
		creator := &models.Creator{}

		if bot.CreatorPic != "" {
			creator.HeadPic = bot.CreatorPic
		} else if userInfo, ok := userMap[int64(bot.CreatorID)]; ok {
			creator.HeadPic = userInfo.Headimgurl
		}

		if bot.CreatorName != "" {
			creator.Name = bot.CreatorName
		} else if userInfo, ok := userMap[int64(bot.CreatorID)]; ok {
			creator.Name = userInfo.WalletAddress
		}

		result := &models.AgentBot{
			Id:            int32(bot.ID),
			Status:        bot.Status,
			Name:          bot.Name,
			Image:         bot.Image,
			Desc:          bot.Description,
			CreatorId:     bot.CreatorID,
			Users:         bot.Users,
			Conversations: bot.Conversations,
			Creator:       creator,
			WelcomeMsg:    bot.GetWelcomeMsg(),
			GuideInfo:     bot.GetGuideInfo(),
			CreatedAt:     timestamppb.New(bot.CreatedAt),
			Tab:           bot.Tab,
		}

		if item.Plugins != nil {
			pluginsPb, err := PluginModelListResult(item.Plugins)
			if err != nil {
				log.Errorf("BotList2Pb PluginModelList2Pb err:%v", err)
			} else {
				result.LinkedPlugin = pluginsPb
			}
		}

		out = append(out, result)
	}

	return out, nil
}

func BotResult(in *model.Bot, plugins []*model.Plugin, needCreatorInfo bool) *models.AgentBot {
	creator := &models.Creator{}
	if needCreatorInfo {
		if in.CreatorPic == "" || in.CreatorName == "" {
			userInfo, _ := user2.GetUser(in.CreatorID)
			if in.CreatorPic == "" {
				creator.HeadPic = userInfo.Headimgurl
			}

			if in.CreatorName == "" {
				creator.Name = userInfo.WalletAddress
			}
		}

		if in.CreatorPic != "" {
			creator.HeadPic = in.CreatorPic
		}
		if in.CreatorName != "" {
			creator.Name = in.CreatorName
		}
	}

	result := &models.AgentBot{
		Id:            int32(in.ID),
		Status:        in.Status,
		Name:          in.Name,
		Image:         in.Image,
		Desc:          in.Description,
		CreatorId:     in.CreatorID,
		Creator:       creator,
		Users:         in.Users,
		Conversations: in.Conversations,
		WelcomeMsg:    in.GetWelcomeMsg(),
		GuideInfo:     in.GetGuideInfo(),
		CreatedAt:     timestamppb.New(in.CreatedAt),
		Tab:           in.Tab,
	}

	if plugins != nil {
		pluginsPb, err := PluginModelListResult(plugins)
		if err != nil {
			log.Errorf("Bot2Pb err:%v", err)
		} else {
			result.LinkedPlugin = pluginsPb
		}
	}
	return result
}
