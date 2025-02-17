// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameBotPlugin = "bot_plugin"

// BotPlugin mapped from table <bot_plugin>
type BotPlugin struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	BotID     int64     `gorm:"column:bot_id;not null" json:"bot_id"`
	PluginID  int64     `gorm:"column:plugin_id;not null" json:"plugin_id"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName BotPlugin's table name
func (*BotPlugin) TableName() string {
	return TableNameBotPlugin
}
