package ws

import (
	"chatApp_backend/model"
	"encoding/json"
	"log"
)

type ClientManger struct {
	Clients  map[string]*Client
	Register chan *Client
	UnRegister chan *Client
	Broadcast chan []byte
	PushChat chan []byte
}

var Manager = ClientManger{
	Clients:    make(map[string]*Client),
	Register:   make(chan *Client),
	UnRegister: make(chan *Client),
	// 传递消息message的管道
	Broadcast:  make(chan []byte),
	// 推送chat聊天框给对方的管道
	PushChat: make(chan []byte),
}

func createId(uid string) string  {
	return uid
}

func (Manager *ClientManger) Start()  {
	for  {
		select {
		// 新用户加入
		case client:= <- Manager.Register:
			log.Printf("用户加入:%s", client.ID)
			Manager.Clients[client.ID] = client
			jsonMessage, _ := json.Marshal(&model.Message{Content: "Successful connection to socket service"})
			client.Send <- jsonMessage
		//	旧用户注销
		case client:= <- Manager.UnRegister:
			log.Printf("用户离开:%s", client.ID)

			if _, ok := Manager.Clients[client.ID]; ok {
				jsonMessage, _ := json.Marshal(&model.Message{Content: "A socket has disconnected"})
				client.Send <- jsonMessage
				close(client.Send)
				delete(Manager.Clients, client.ID)
			}
		case message:=<-Manager.Broadcast:
			MessageStruct :=model.Message{}
			json.Unmarshal(message, &MessageStruct)
			for id,conn	:=	range Manager.Clients{
				// id = createId(uid,to_id)
				// 如果该消息的接收者 id 等于该 client 的 uid ,则将该消息发往该 client
				// 连接 Clients[MessageStruct.Recipient] 不存在，表示消息接收者没有 登录 ，未读消息存入数据库
				if id != createId(MessageStruct.Recipient) {
					continue
				}
				select {
				// message 发给接收方
				// 如果能发送给接收方 ， 说明接收方 登陆了
				case conn.Send <- message:
					// 只要能够互相通过 conn.Send 发送的话，就直接 将消息置为 已读 ，这样的话不会出现
					// 两个人正在互相聊天中可互相的消息都是未读的情况，但是有点冗余
					_ =model.ModifyMsgState(MessageStruct.Sender,MessageStruct.Recipient)


				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}

			}
		}
	}
}