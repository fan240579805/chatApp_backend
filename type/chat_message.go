package _type

import (
	"chatApp_backend/model"
)

// BePushedMsg 推送给制定用户id的新的消息的结构体
type BePushedMsg struct {
	DataType   string // 区分websocket要推送的数据是好友类型还是消息类型
	BePushedID string
	Message    model.Message
	UserInfo   *model.UserInfo
}

// ChatItem 聊天框item
type ChatItem struct {
	ChatID           string
	RecentMsg        model.Message
	ChatToNickName   string
	ChatToUserID     string
	ChatToUserAvatar string
	RecentTime       int64
}

type BePushedChat struct {
	DataType   string // 区分websocket要推送的数据是聊天框
	BePushedID string
	Chat       ChatItem
}

type WsMessageObj struct {
	DataType string // 区分websocket要推送的数据是好友类型还是消息类型
	ChatID   string
	Message  model.Message
}
