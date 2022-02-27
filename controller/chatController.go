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

func ModifyMsgReadState(c *gin.Context) {
	msgFrom := c.Query("msgFrom")
	msgTo := c.Query("msgTo")
	err := model.ModifyMsgState(msgFrom, msgTo)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":  444,
			"state": "false",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"state": "success",
	})
}

func GetChat(c *gin.Context) {
	msgs, err := model.GetChatContent()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":  444,
			"state": "false",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"chatList": msgs,
	})
}



// PushChatMsg2User 将聊天内容通过ws推送给在线用户
func PushChatMsg2User(pushedUserid string, messageData model.Message) {
	pushedObj := &_type.BePushedMsg{
		DataType:   "msg",
		BePushedID: pushedUserid,
		Message:    messageData,
	}

	msgByte, err := json.Marshal(pushedObj)
	if err != nil {
		log.Println("解析推送消息出错")
	}
	ws.Manager.Broadcast <- msgByte
}
