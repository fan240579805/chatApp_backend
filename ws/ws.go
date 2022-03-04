package ws

import (
	"chatApp_backend/controller"
	"chatApp_backend/model"
	_type "chatApp_backend/type"
	"chatApp_backend/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

// Client is a websocket client
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
	//User   model.User
	mutex sync.Mutex
}

func WsHandler(c *gin.Context) {
	uid := c.Query("userid")
	fmt.Println(uid)
	//to_uid:=c.Query("to_uid")

	// ws升级器
	ws, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}
	// 获取用户 登陆状态
	//id,_ := strconv.Atoi(uid)
	//var u1 model.User
	//dao.DB.Where("id=?",id).First(&u1)
	client := &Client{
		ID:     uid,
		Socket: ws,
		Send:   make(chan []byte, 256),
		//UserInfo: u1,
	}

	//申明定时器15s，设置心跳时间为15s
	//ticker := time.NewTicker(time.Second * 15)
	//go client.timeWriter(ticker)
	// 将该登陆用户存入 Register Map中
	Manager.Register <- client
	go client.Read()
	go client.Write()
}

func (c *Client) Read() {
	// 延迟关闭客户端连接
	defer func() {
		Manager.UnRegister <- c
		c.Socket.Close()
	}()
	// 不断从连接中读取message
	for {
		_, messageChatObj, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.UnRegister <- c
			c.Socket.Close()
			break
		}
		log.Printf("读取到客户端的信息:%s", string(messageChatObj))

		// 将发送时间赋给message
		var wsMsgObj _type.WsMessageObj
		json.Unmarshal(messageChatObj, &wsMsgObj)
		wsMsgObj.Message.SendTime = time.Now()
		wsMsgObj.Message.MsgID = "MsgID_" + utils.UniqueId()
		//if m1.Type == "img" {
		//	m1.Image = model.ImgUrl
		//}

		messageBody, _ := json.Marshal(&wsMsgObj.Message)

		// 发给自身
		c.Send <- messageBody

		// 存入数据库
		msg := model.Message{}
		json.Unmarshal(messageBody, &msg)
		AddMsgErr := model.AddMessageRecord(msg)
		if AddMsgErr != nil {
			log.Println(err)
		}
		// 更新最近消息
		controller.ModifyRecentMsg(wsMsgObj.ChatID, string(messageBody))

		Manager.Broadcast <- messageChatObj
	}
}

// 将数据流写回给前端
func (c *Client) Write() {
	// 延迟关闭客户端连接
	defer func() {
		Manager.UnRegister <- c
		c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("发送到到客户端的信息:%s", string(message))

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
