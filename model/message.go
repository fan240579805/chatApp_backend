package model

import (
	"chatApp_backend/dao"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Message ORM映射结构体
type Message struct {
	MID       int       `gorm:"column:mid;primary_key;AUTO_INCREMENT"`
	MsgID     string    `gorm:"column:msgid"`
	Sender    string    `gorm:"column:sender" json:"sender"`       //  发送者唯一id
	Recipient string    `gorm:"column:recipient" json:"recipient"` //	接收者唯一id
	Content   string    `json:"content"`
	SendTime  int64     `json:"time"`
	Type      string    `json:"type"` //消息类型 img: 图片 text: 文本 audio: 音频
	CreatedAt time.Time `gorm:"column:createdat;default:null" json:"createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat;default:null" json:"updatedat"`
}

func (m Message) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan 实现方法
func (m *Message) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), m)
}

//	AddMessageRecord 往聊天记录表添加记录
func AddMessageRecord(msg Message) error {
	err := dao.DB.Debug().Create(&msg).Error
	if err != nil {
		return err
	}
	return nil
}


func SelectMessageRecord(msgID string) (Message, error) {
	var message Message
	err := dao.DB.Where("msgid=?", msgID).First(&message).Error
	if err != nil {
		return Message{}, err
	}
	return message, nil
}

func SelectMessages(mine string, other string) ([]Message, error) {
	// 发起该段好友relation的数组，from == userid
	var mineMessages []Message
	// 接受这段好友relation的数组  to == userid
	var otherMessages []Message
	// 2：此时请求friendList的用户在relation表中是from, 不需要from主动删除即 status = 2 的好友
	selectMineErr := dao.DB.Debug().Where("sender=? AND recipient=?", mine, other).Find(&mineMessages).Error
	if selectMineErr != nil {
		return mineMessages, selectMineErr
	}
	// 3：此时请求friendList的用户在relation表中是to, 不需要to主动删除即 status = 3 的好友
	selectOtherErr := dao.DB.Debug().Where("sender=? AND recipient=?", other, mine).Find(&otherMessages).Error
	if selectOtherErr != nil {
		return otherMessages, selectOtherErr
	}
	if len(mineMessages) > 0 || len(otherMessages) > 0 {
		return append(mineMessages, otherMessages...), nil
	} else {
		return []Message{}, nil
	}
}
