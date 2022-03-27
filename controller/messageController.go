package controller

import (
	"chatApp_backend/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

type chatBetweenType struct {
	Myself    string
	Other     string
	PageIndex int
}

type msgArray []model.Message

func (array msgArray) Len() int {
	return len(array)
}

func (array msgArray) Less(i, j int) bool {
	return array[i].SendTime > array[j].SendTime // 若为大于号，则从大到小
}

func (array msgArray) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
}

func reserveCurMsgList(m []model.Message) []model.Message {
	for i, j := 0, len(m)-1; i < j; i, j = i+1, j-1 {
		m[i], m[j] = m[j], m[i]
	}
	return m
}

// GetMsgList 获取聊天list
func GetMsgList(c *gin.Context) {
	var params chatBetweenType
	c.ShouldBindJSON(&params)
	var onePageCnt = 20
	messageList, selectErr := model.SelectMessages(params.Myself, params.Other)
	if selectErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 444,
			"msg":  "获取聊天记录失败",
			"data": messageList,
		})
	} else {
		sort.Sort(msgArray(messageList))
		page := params.PageIndex
		var startIndex = page * onePageCnt
		var endIndex = (page + 1) * onePageCnt
		if startIndex > len(messageList)-1 {
			// 没有聊天记录数据了
			c.JSON(http.StatusOK, gin.H{
				"code": 2333,
				"msg":  "没有聊天记录数据了",
				"data": "",
			})
			return
		}

		if endIndex >= len(messageList)-1 {
			endIndex = len(messageList) - 1
		}
		// 分页获取聊天记录，每页20
		curPageMsgList := messageList[startIndex:endIndex]
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "",
			"data": reserveCurMsgList(curPageMsgList), // 因为排序是根据sendTime从大到小，而前端需要大的在最底部，所以reserve一下
		})
	}
}
