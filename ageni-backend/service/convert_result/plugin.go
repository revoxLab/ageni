package convert_result

import (
	"encoding/json"
	"github.com/readonme/open-studio/common/log"
	"github.com/readonme/open-studio/dal/model"
	entity2 "github.com/readonme/open-studio/service/entity"
	"github.com/readonme/open-studio/service/models"
	user2 "github.com/readonme/open-studio/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PluginListResult(in entity2.PluginAggregateSlice) ([]*models.Plugin, error) {
	if len(in) == 0 {
		return make([]*models.Plugin, 0), nil
	}

	creatorIDs := make([]int64, 0, len(in))
	for _, item := range in {
		if item.Plugin.CreatorName == "" || item.Plugin.CreatorPic == "" {
			creatorIDs = append(creatorIDs, int64(item.Plugin.CreatorID))
		}
	}

	userMap := make(map[int64]*model.User)

	if len(creatorIDs) > 0 {
		users := user2.GetUserByIds(creatorIDs)
		for _, u := range users {
			userMap[u.ID] = u
		}
	}

	out := make([]*models.Plugin, 0, len(in))
	for _, item := range in {
		plugin := item.Plugin
		methods, err := MethodListResult(item.Methods)
		if err != nil {
			return nil, err
		}

		creator := &models.Creator{}

		if plugin.CreatorPic != "" {
			creator.HeadPic = plugin.CreatorPic
		} else if userInfo, ok := userMap[int64(plugin.CreatorID)]; ok {
			creator.HeadPic = userInfo.Headimgurl
		}

		if plugin.CreatorName != "" {
			creator.Name = plugin.CreatorName
		} else if userInfo, ok := userMap[int64(plugin.CreatorID)]; ok {
			creator.Name = userInfo.WalletAddress
		}

		result := &models.Plugin{
			Id:            plugin.ID,
			Tab:           plugin.Tab,
			Name:          plugin.Name,
			Creator:       creator,
			Methods:       methods,
			Image:         plugin.ImageURL,
			Desc:          plugin.Description,
			CreatorId:     plugin.CreatorID,
			CreatedAt:     timestamppb.New(plugin.CreatedAt),
			DependPlugins: plugin.GetDependPluginIds(),
		}

		if item.Bots != nil {
			result.LinkedAgent = make([]*models.AgentBot, 0, len(item.Bots))
			for _, bot := range item.Bots {
				result.LinkedAgent = append(result.LinkedAgent, BotResult(bot, nil, false))
			}
		}

		out = append(out, result)
	}

	return out, nil
}

func PluginModelListResult(in []*model.Plugin) ([]*models.Plugin, error) {
	if len(in) == 0 {
		return make([]*models.Plugin, 0), nil
	}

	creatorIDs := make([]int64, 0, len(in))
	for _, plugin := range in {
		if plugin.CreatorName == "" || plugin.CreatorPic == "" {
			creatorIDs = append(creatorIDs, int64(plugin.CreatorID))
		}
	}

	userMap := make(map[int64]*model.User)

	if len(creatorIDs) > 0 {
		users := user2.GetUserByIds(creatorIDs)
		for _, u := range users {
			userMap[u.ID] = u
		}
	}
	out := make([]*models.Plugin, 0, len(in))
	for _, plugin := range in {
		creator := &models.Creator{}

		if plugin.CreatorPic != "" {
			creator.HeadPic = plugin.CreatorPic
		} else if userInfo, ok := userMap[int64(plugin.CreatorID)]; ok {
			creator.HeadPic = userInfo.Headimgurl
		}

		if plugin.CreatorName != "" {
			creator.Name = plugin.CreatorName
		} else if userInfo, ok := userMap[int64(plugin.CreatorID)]; ok {
			creator.Name = userInfo.WalletAddress
		}

		result := &models.Plugin{
			Id:            plugin.ID,
			Tab:           plugin.Tab,
			Name:          plugin.Name,
			Creator:       creator,
			Image:         plugin.ImageURL,
			Desc:          plugin.Description,
			CreatorId:     plugin.CreatorID,
			CreatedAt:     timestamppb.New(plugin.CreatedAt),
			DependPlugins: plugin.GetDependPluginIds(),
		}

		out = append(out, result)
	}

	return out, nil
}

func PluginResult(in *model.Plugin, bots []*model.Bot, methodIn []*model.Method) (*models.Plugin, error) {
	out := &models.Plugin{}
	methods, err := MethodListResult(methodIn)
	if err != nil {
		return nil, err
	}

	creator := &models.Creator{}
	userInfo, _ := user2.GetUser(int64(in.CreatorID))

	if in.CreatorPic != "" {
		creator.HeadPic = in.CreatorPic
	} else {
		if err != nil {
			log.Errorf("PluginResult GetUser err:%v", err)
		} else {
			creator.HeadPic = userInfo.Headimgurl
		}
	}

	if in.CreatorName != "" {
		creator.Name = in.CreatorName
	} else {
		creator.Name = userInfo.WalletAddress
	}

	if bots != nil {
		out.LinkedAgent = make([]*models.AgentBot, 0, len(bots))
		for _, item := range bots {
			out.LinkedAgent = append(out.LinkedAgent, BotResult(item, nil, false))
		}
	}

	out.Id = in.ID
	out.Tab = in.Tab
	out.Name = in.Name
	out.Creator = creator
	out.Methods = methods
	out.Image = in.ImageURL
	out.Desc = in.Description
	out.CreatorId = in.CreatorID
	out.DependPlugins = in.GetDependPluginIds()
	out.CreatedAt = timestamppb.New(in.CreatedAt)
	return out, nil
}
func MethodListResult(in []*model.Method) ([]*models.Method, error) {
	if in == nil {
		return nil, nil
	}
	out := make([]*models.Method, 0, len(in))
	for _, item := range in {
		_out, err := MethodResult(item)
		if err != nil {
			return nil, err
		}
		out = append(out, _out)
	}
	return out, nil
}

func MethodResult(in *model.Method) (*models.Method, error) {
	if in == nil {
		return nil, nil
	}
	inputSchema := &model.InputSchema{}
	if err := json.Unmarshal([]byte(in.InputSchema), inputSchema); err != nil {
		return nil, err
	}
	out := &models.Method{Id: in.ID}
	out.Name = in.Name
	out.Description = in.Description
	out.PluginId = in.PluginID
	out.HttpSubPath = in.HTTPSubPath
	out.HttpMethod = in.HTTPMethod
	out.MethodCallName = in.MethodCallName
	out.OutputSchema = in.OutputSchema
	out.OutputExample = in.OutputExample
	out.InputSchema = in.InputSchema
	out.InputExample = in.InputExample
	out.Status = in.Status
	return out, nil
}
