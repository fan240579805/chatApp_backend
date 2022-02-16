package _type

import (
	"chatApp_backend/model"
	"time"
)

// BePushedFriend 推送给制定用户id的新的好友的结构体
type BePushedFriend struct {
	DataType   string   // 区分websocket要推送的数据是好友类型还是消息类型
	BePushedID string
	Friend     Friend
}

// Friend 好友列表的每一个好友信息
type Friend struct {
	FriendProfile *model.UserInfo
	AddTime       time.Time
	Status        int
}
