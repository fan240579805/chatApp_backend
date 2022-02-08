package model

import (
	"chatApp_backend/dao"
	"time"
)

// Message is return msg
type Message struct {ID        int       `gorm:"column:id;unique;not null;primary_key;AUTO_INCREMENT"`

	MsgID     int       `gorm:"column:msgid;unique;not null"`
	Sender    string    `gorm:"column:sender;unique;not null" json:"sender"`       //  发送者唯一id
	Recipient string    `gorm:"column:recipient;unique;not null" json:"recipient"` //	接收者唯一id
	Content   string    `json:"content"`
	SendTime  int64     `json:"time"`
	Type      string    `json:"type"` //消息类型 img: 图片 text: 文本 audio: 音频
	CreatedAt time.Time `gorm:"column:createdat;default:null" json:"createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat;default:null" json:"updatedat"`
}

//	AddMessageRecord 往聊天记录表添加记录
func AddMessageRecord(msg Message) error {
	err := dao.DB.Create(&msg).Error
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

	err := dao.DB.Debug().Model(&Message{}).Where("sender=? AND recipient =? ", msgFrom, msgTo).Update("unread", true).Error
	if err != nil {
		return err
	}
	return nil
}

func GetChatContent() ([]Message, error) {
	var m1 []Message
	err := dao.DB.Find(&m1).Error
	if err != nil {
		return m1, err
	}

	return m1, nil
}
