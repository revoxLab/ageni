package model

import (
	"encoding/json"
)

func (b *BotDraft) GetPluginIds() ([]int64, error) {
	var ids []int64
	if len(b.Plugins) == 0 {
		return ids, nil
	}
	err := json.Unmarshal([]byte(b.Plugins), &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}
