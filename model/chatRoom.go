package model

import (
	"chatApp_backend/dao"
	"chatApp_backend/utils"
	"time"
)

// 每一个会话的 model
type Chat struct {
	ID        int       `gorm:"column:id;unique;not null;primary_key;AUTO_INCREMENT"`
	ChatID    string    `gorm:"column:chatid;unique;not null"`
	Owner     string    `gorm:"column:ownerid"`          // 改对话框所有者, 发起聊天的用户的id
	Other     string    `gorm:"column:otherid"`          // 被发起聊天的Other其他人id
	Unread    int       `gorm:"column:unread;default:0"` // 会话未读数量
	RecentMsg string    `gorm:"column:recentmsg"`        // 最近一条聊天记录用于展示聊天列表, Message类型的json字符串
	CreatedAt time.Time `gorm:"column:createdat;default:null" json:"createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat;default:null" json:"updatedat"`
}

// CreateChat 创建一条对话记录
func CreateChat(sender string, recipient string) (*Chat, error) {
	chat := &Chat{
		ChatID:    "chatID_" + utils.UniqueId(),
		Owner:     sender,
		Other:     recipient,
		Unread:    0,
		RecentMsg: "",
	}
	err := dao.DB.Create(&chat).Error
	if err != nil {
		return &Chat{}, err
	}
	return chat, nil
}

func CheckChatExist(sender string, recipient string) (Chat, error) {
	var chat Chat
	finderr1 := dao.DB.Debug().Where("ownerid = ? and otherid= ?", sender, recipient).First(&chat).Error
	if finderr1 != nil {
		// 说明没查到用例1，继续查
		finderr2 := dao.DB.Debug().Where("ownerid = ? and otherid= ?", recipient, sender).First(&chat).Error
		if finderr2 != nil {
			return chat, finderr2
		} else {
			return chat, nil
		}
		return Chat{}, finderr1
	}
	return chat, nil
}

// UpdateUnRead 更新未读消息数量
func UpdateUnRead(chatID string, isAdd bool) error {

	chat, selectErr := SelectChatRecord(chatID)
	if selectErr != nil {
		return selectErr
	}
	var newReadCnt = 0
	if isAdd {
		newReadCnt = chat.Unread + 1
	}
	err := dao.DB.Model(&chat).Update("unread", newReadCnt).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateRecentMsg 更新最近消息
func UpdateRecentMsg(chatID string, newMsg string) error {
	chat, selectErr := SelectChatRecord(chatID)
	if selectErr != nil {
		return selectErr
	}
	err := dao.DB.Model(&chat).Update("recentmsg", newMsg).Error
	if err != nil {
		return err
	}
	return nil
}

// SelectChatRecord 查询某 聊天框
func SelectChatRecord(chatID string) (Chat, error) {
	var chat Chat
	err := dao.DB.Where("chatid=?", chatID).First(&chat).Error
	if err != nil {
		return chat, err
	}
	return chat, nil
}

// SelectChatList 从数据库筛选出聊天列表
func SelectChatList(userid string) ([]Chat, error) {
	var mineChatList []Chat
	var otherChatList []Chat
	// 查询自己主动发起的chat列表
	selectMineErr := dao.DB.Debug().Where("owner=?", userid).Find(&mineChatList).Error
	if selectMineErr != nil {
		return mineChatList, selectMineErr
	}
	// 查询自己被动发起的chat列表
	selectOtherErr := dao.DB.Debug().Where("other=?", userid).Find(&otherChatList).Error
	if selectOtherErr != nil {
		return otherChatList, selectOtherErr
	}
	if len(mineChatList) > 0 || len(otherChatList) > 0 {
		return append(mineChatList, otherChatList...), nil
	} else {
		return []Chat{}, nil
	}
}
