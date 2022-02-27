package _type

import "chatApp_backend/model"

// BePushedMsg 推送给制定用户id的新的好友的结构体
type BePushedMsg struct {
	DataType   string // 区分websocket要推送的数据是好友类型还是消息类型
	BePushedID string
	Message     model.Message
}