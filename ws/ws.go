package ws

import (
	"chatApp/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// Client is a websocket client
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
	//UserInfo model.User
}

func WsHandler(c *gin.Context) {
	uid := c.Query("uid")
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
		ID:     createId(uid),
		Socket: ws,
		Send:   make(chan []byte, 256),
		//UserInfo: u1,
	}
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
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.UnRegister <- c
			c.Socket.Close()
			break
		}
		log.Printf("读取到客户端的信息:%s", string(message))

		// 将发送时间赋给message
		var m1 model.Message
		json.Unmarshal(message, &m1)
		m1.SendTime = time.Now().Unix()

		if m1.Type == "img" {
			m1.Image = model.ImgUrl
		}

		message1, _ := json.Marshal(&m1)

		// 发给自身
		c.Send <- message1

		// 存入数据库
		msg := model.Message{}
		json.Unmarshal(message1, &msg)
		err = model.AddMessageRecord(msg)
		if err != nil {
			log.Println(err)
		}

		Manager.Broadcast <- message1
	}
}

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
