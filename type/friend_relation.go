package _type

import (
	"chatApp_backend/model"
	"time"
)

// BePushedFriend 推送给制定用户id的新的好友的结构体
type BePushedFriend struct {
	BePushedID string
	Friend Friend
}


// Friend 好友列表的每一个好友信息
type Friend struct {
	FriendProfile *model.UserInfo
	AddTime       time.Time
	Status        int
}