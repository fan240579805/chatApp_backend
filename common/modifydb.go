package common

import (
	"chatApp_backend/model"
	"log"
)

// ModifyRecentMsg 更新最近消息
func ModifyRecentMsg(chatID string, newRecentMsg model.Message) {
	err := model.UpdateRecentMsg(chatID, newRecentMsg)
	if err != nil {
		log.Println(err.Error())
	}
}

// ModifyUnRead 更新聊天框未读数量  isAdd = true 表示在原来基础上 + 1
func ModifyUnRead(chatID string, isAdd bool) {
	err := model.UpdateUnRead(chatID, isAdd)
	if err != nil {
		log.Println(err.Error())
	}
}