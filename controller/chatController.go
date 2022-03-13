package controller

import (
	"chatApp_backend/model"
	_type "chatApp_backend/type"
	"chatApp_backend/ws"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type chatType struct {
	Sender    string
	Recipient string
}

// MakeChat 发起聊天
func MakeChat(c *gin.Context) {
	var chatParams chatType
	c.ShouldBindJSON(&chatParams)
	existChat, checkErr := model.CheckChatExist(chatParams.Sender, chatParams.Recipient)
	if checkErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 444,
			"msg":  "发起聊天失败",
		})
	} else {
		var userProfile *model.UserInfo
		// 已经存在一条聊天框，需要对发起makeChat请求的用户进行判断，看他是发起聊天的用户还是被发起聊天的用户
		if existChat.Owner == chatParams.Sender {
			// 此时 发起makeChat 的就是聊天框的拥有者，所以这时候需要把对方other的信息查询出来
			userProfile, _ = model.SelectUser(existChat.Other)
		} else {
			// 此时 发起makeChat 的就是聊天框的被发起者，所以这时候需要把owner的信息查询出来
			userProfile, _ = model.SelectUser(existChat.Owner)
		}
		chatItem := &_type.ChatItem{
			ChatID:           existChat.ChatID,
			RecentMsg:        existChat.RecentMsg,
			ChatToNickName:   userProfile.NickName,
			ChatToUserID:     userProfile.UserID,
			ChatToUserAvatar: userProfile.Avatar,
			RecentTime:       existChat.UpdatedAt.UnixMilli(),
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "发起聊天成功",
			"data": chatItem,
		})
	}
	if existChat.ChatID == "" {
		chat, err := model.CreateChat(chatParams.Sender, chatParams.Recipient)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 444,
				"msg":  "发起聊天失败",
			})
		} else {
			userProfile, _ := model.SelectUser(chat.Other)
			chatItem := &_type.ChatItem{
				ChatID:           chat.ChatID,
				RecentMsg:        chat.RecentMsg,
				ChatToNickName:   userProfile.NickName,
				ChatToUserID:     userProfile.UserID,
				ChatToUserAvatar: userProfile.Avatar,
				RecentTime:       chat.UpdatedAt.UnixMilli(),
			}
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "发起聊天成功",
				"data": chatItem,
			})
		}
	}
}

// GetMineChatList 获取聊天列表
func GetMineChatList(c *gin.Context) {
	userid, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 2003,
			"msg":  "auth中间件获取userid失败",
		})
	} else {
		// 处理一下ChatList并返回给前端
		chatList := FormatChatList(userid.(string), c)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取聊天列表成功",
			"data": chatList,
		})
	}
}

// FormatChatList 格式化获得 chatList ;  userid为当前登录用户
func FormatChatList(userid string, c *gin.Context) []*_type.ChatItem {
	chatRoomList, err := model.SelectChatList(userid)
	var chatList = make([]*_type.ChatItem, len(chatRoomList))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 2003,
			"msg":  "获取聊天列表失败",
		})
	} else {
		for i, chatRoom := range chatRoomList {
			if chatRoom.Owner == userid {
				// 此时自己是发起添聊天的
				// 不要把自己的信息录入 chatList
				userProfile, _ := model.SelectUser(chatRoom.Other)
				chatList[i] = &_type.ChatItem{
					ChatID:           chatRoom.ChatID,
					RecentMsg:        chatRoom.RecentMsg,
					ChatToNickName:   userProfile.NickName,
					ChatToUserID:     userProfile.UserID,
					ChatToUserAvatar: userProfile.Avatar,
					RecentTime:       chatRoom.UpdatedAt.UnixMilli(),
				}
			} else if chatRoom.Other == userid {
				// 此时自己是被发起聊天的
				// 不要把自己的信息录入 chatList
				userProfile, _ := model.SelectUser(chatRoom.Owner)
				chatList[i] = &_type.ChatItem{
					ChatID:           chatRoom.ChatID,
					RecentMsg:        chatRoom.RecentMsg,
					ChatToNickName:   userProfile.NickName,
					ChatToUserID:     userProfile.UserID,
					ChatToUserAvatar: userProfile.Avatar,
					RecentTime:       chatRoom.UpdatedAt.UnixMilli(),
				}
			}
		}
	}

	return chatList
}

// PushChatMsg2User 将聊天内容通过ws推送给在线用户
func PushChatMsg2User(chatID string, bePushedID string, messageData model.Message) {

	var wsMsgObj _type.WsMessageObj
	wsMsgObj.DataType = "msg"
	wsMsgObj.Message = messageData
	wsMsgObj.ChatID = chatID
	// 修改消息的recipient，因为websocket-hub只认 Recipient 进行推送，此时修改Recipient也不担心数据库字段被更改，因为以及存入数据库了
	wsMsgObj.Message.Recipient = bePushedID

	msgByte, err := json.Marshal(wsMsgObj)
	if err != nil {
		log.Println("解析推送消息出错")
	}
	ws.Manager.Broadcast <- msgByte
}
