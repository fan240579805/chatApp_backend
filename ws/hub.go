package ws

import (
	"chatApp_backend/common"
	"chatApp_backend/model"
	_type "chatApp_backend/type"
	"encoding/json"
	"log"
)

type ClientManger struct {
	Clients    map[string]*Client
	Register   chan *Client
	UnRegister chan *Client
	Broadcast  chan []byte
	ChatChan   chan []byte
	FriendChan chan []byte
}

var Manager = ClientManger{
	Clients:    make(map[string]*Client),
	Register:   make(chan *Client),
	UnRegister: make(chan *Client),
	Broadcast:  make(chan []byte), // 传递消息message的管道
	ChatChan:   make(chan []byte), // 推送chat聊天框给对方的管道
	FriendChan: make(chan []byte), // 推送firendReq给对方的管道
}

func (Manager *ClientManger) Start() {
	for {
		select {
		// 新用户加入
		case client := <-Manager.Register:
			log.Printf("用户加入:%s", client.ID)
			Manager.Clients[client.ID] = client
			jsonMessage, _ := json.Marshal(&model.Message{Content: "Successful connection to socket service"})
			client.Send <- jsonMessage
		//	旧用户注销
		case client := <-Manager.UnRegister:
			log.Printf("用户离开:%s", client.ID)

			if _, ok := Manager.Clients[client.ID]; ok {
				jsonMessage, _ := json.Marshal(&model.Message{Content: "A socket has disconnected"})
				client.Send <- jsonMessage
				close(client.Send)
				delete(Manager.Clients, client.ID)
			}

		case WsMessageObj := <-Manager.Broadcast:
			MessageChatStruct := _type.WsMessageObj{}
			json.Unmarshal(WsMessageObj, &MessageChatStruct)
			// 判断一下对方是否登录
			if !otherUserIsLogin(MessageChatStruct.Message.Recipient) {
				// 连接 Clients[MessageStruct.Recipient] 不存在，表示消息接收者没有登录 ，未读消息存入数据库
				// 并且给chat数据记录未读数量+1
				common.ModifyUnRead(MessageChatStruct.ChatID, true)
			}

			for id, conn := range Manager.Clients {
				// 如果该消息的接收者 id 等于该 client 的 uid ,则将该消息发往该 client
				if id != MessageChatStruct.Message.Recipient {
					// 直接跳转下一次循环
					continue
				}
				// 格式统一处理成bePushed格式
				userInfo, _ := model.SelectUser(MessageChatStruct.Message.Sender)
				bePushedMsg := _type.BePushedMsg{
					DataType:   "msg",
					BePushedID: MessageChatStruct.Message.Recipient,
					Message:    MessageChatStruct.Message,
					UserInfo: userInfo,
				}
				messageBodyByte, _ := json.Marshal(&bePushedMsg)
				select {
				// message 发给接收方; 如果能发送给接收方 ， 说明接收方 登陆了
				case conn.Send <- messageBodyByte:
					// 发送聊天图片，message的recipient会t改为sender自己，目的是为了推送给自己是自己展示
					// 这个if是为了过滤的这种情况
					if MessageChatStruct.Message.Recipient != MessageChatStruct.Message.Sender{
						// 登录了，也要chat unread++ ，因为前端需要全局小红点来提示已登录用户
						common.ModifyUnRead(MessageChatStruct.ChatID, true)
						// push给对方一个chat  ***!前端结合recentMsg是否是自己发的来确定是否展示小红点，以及是否清除小红点
						chatRoom, _ := model.SelectChatRecord(MessageChatStruct.ChatID)
						// 获取自身的简要信息,以便发给对方
						userProfile, _ := model.SelectUser(MessageChatStruct.Message.Recipient)
						var chatItem = &_type.ChatItem{
							ChatID:           chatRoom.ChatID,
							RecentMsg:        chatRoom.RecentMsg,
							ChatToNickName:   userProfile.Username,
							ChatToUserID:     userProfile.UserID,
							ChatToUserAvatar: userProfile.Avatar,
							RecentTime:       chatRoom.UpdatedAt.UnixMilli(),
						}
						chatItemByte, _ := json.Marshal(&chatItem)
						// push
						conn.Send <- chatItemByte
					}

				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}

		case newFriendDetail := <-Manager.FriendChan:
			newPushObject := _type.BePushedFriend{}
			json.Unmarshal(newFriendDetail, &newPushObject)
			for id, conn := range Manager.Clients {
				// 如果该消息的接收者 id 等于该 client 的 uid ,则将新的好友列表发往该 client
				// 连接 Clients[newPushObject.BePushedID] 不存在，表示被推送新好友列表的用户登录 ，进行相应的离线处理
				if id != newPushObject.BePushedID {
					// 暂存
					continue
				}
				select {
				// 新的好友列表 发给被推送方
				// 如果能发送给接收方 ， 说明接收方 登陆了
				case conn.Send <- newFriendDetail:
					// 发送完执行相应逻辑，也可以不执行

				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}

			}
		}
	}
}

// otherUserIsLogin 判断对方用户是否在线
func otherUserIsLogin(otherID string) bool {
	for uid, _ := range Manager.Clients {
		if otherID == uid {
			return true
		}
	}
	return false
}
