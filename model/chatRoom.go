package model

import (
	"time"
)

// 每一个会话的 model
type Chat struct {
	ID         int       `gorm:"column:id;unique;not null;primary_key;AUTO_INCREMENT"`
	ChatID     string    `gorm:"column:chatid;unique;not null;"`
	Owner      string    `gorm:"column:owner;not null"`   // 改对话框所有者, 发起聊天的用户的id
	Other      string    `gorm:"column:other;not null"`   // 被发起聊天的Other其他人id
	Unread     int       `gorm:"column:unread;default:0"` // 会话未读数量
	RecentMsg  string    `gorm:"column:recentmsg"`        // 最近一条聊天记录用于展示聊天列表, Message类型的json字符串
	CreateTime time.Time `gorm:"column:createtime;default:null" json:"createtime"`
	UpdateTime time.Time `gorm:"column:updatetime;default:null" json:"updatetime"`
}
