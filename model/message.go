package model

import (
	"chatApp/dao"
)

// Message is return msg
type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	Unread    bool   `json:"unread"`
	FromUser  string `json:"fromUser"`
	ToUser    string `json:"toUser"`
	SendTime  int64  `json:"time"`
	Image     string `json:"image"`
	Type      string `json:"type"`
}

func AddMessageRecord(msg Message) error {

	err:=dao.DB.Create(&msg).Error
	if err != nil {
		return err
	}
	return nil
}

func ModifyMsgState(msgFrom string, msgTo string) error {

	// 1. 先查询出 该条要变为 已读的消息
	//var m1 []Message
	//dao.DB.Debug().Where("sender=? AND recipient =? ",msgFrom,msgTo).Find(&m1)
	//fmt.Println("m1",m1)

	err:=dao.DB.Debug().Model(&Message{}).Where("sender=? AND recipient =? ",msgFrom,msgTo).Update("unread",true).Error
	if err != nil {
		return err
	}
	return nil
}

func GetChatContent() ([]Message,error) {
	var m1  []Message
	err:=dao.DB.Find(&m1).Error
	if err != nil {
		return m1,err
	}

	return m1,nil
}