package _type

import (
	"chatApp_backend/model"
	"time"
)

// BePushedMsg 推送给制定用户id的新的消息的结构体
type BePushedMsg struct {
	DataType   string // 区分websocket要推送的数据是好友类型还是消息类型
	BePushedID string
	Message    model.Message
}

// ChatItem 聊天框item
type ChatItem struct {
	ChatID           string
	RecentMsg        string
	ChatToUserName   string
	ChatToUserID     string
	ChatToUserAvatar string
	RecentTime       time.Time
}

type WsMessageObj struct {
	ChatID           string
	Message    model.Message
}