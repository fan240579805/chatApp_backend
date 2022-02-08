package controller

import (
	"chatApp_backend/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type addFriendReqParams struct {
	Username string
	Fromid   string
	Toid     string
}

// AddFriendReq from用户像to用户发起好友请求（to还没同意）
func AddFriendReq(c *gin.Context) {
	var params addFriendReqParams
	bindErr := c.ShouldBindJSON(&params)
	if bindErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg":  "服务器获取请求失败",
			"code": 2004,
		})
	}
	addFriendERR := model.AddFriendRecord(params.Username, params.Fromid, params.Toid)
	if addFriendERR != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 2003,
			"msg":  "添加请求失败",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "添加请求成功",
		})
		// 执行推送给to用户好友请求逻辑 pushToUser
	}
}

// AcceptFriendReq to接受from用户发起的好友请求，将 status 设为1
func AcceptFriendReq(c *gin.Context) {
	var relation model.Relation
	c.ShouldBindJSON(&relation)
	fmt.Println("getright", relation)

	modifyErr := model.ModifyStatus(relation.From, relation.To, 1)
	if modifyErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 2003,
			"msg":  "接受失败",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "接受成功",
			"data": "新的好友列表给接受请求的ToUser",
		})
		// pushFromUser 将新的好友列表推送给发起好友请求的FromUser

	}
}

// DeleteFriendReq 某用户发起的删除好友请求，
// 根据发起删除的用户是to还是from来为 status 设 2：from删除to 还是 3：to删除from
func DeleteFriendReq(c *gin.Context) {
	// 删除好友时并不知道到底 发起删除的用户是from还是被删除人是from
	var unknownRecord model.Relation
	c.ShouldBindJSON(&unknownRecord)
	fmt.Println("getright", unknownRecord)
	// 先判断一下
	rightRelation, _ := model.GetRightRelationRecord(unknownRecord.From, unknownRecord.To)
	if unknownRecord.From == rightRelation.From {
		// 发起删除的用户就是  发起好友请求的人（from） status = 2
		modifyErr := model.ModifyStatus(rightRelation.From, rightRelation.To, 2)
		if modifyErr != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 2003,
				"msg":  "删除失败",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "删除成功",
				"data": "新的好友列表给发起删除的FromUser",
			})
			// 单方面删除，不需要给被删除的to推送好友列表
		}
	} else if unknownRecord.From == rightRelation.To {
		// 发起删除的用户是  一开始接受好友请求的人（to） status = 3
		modifyErr := model.ModifyStatus(rightRelation.From, rightRelation.To, 3)
		if modifyErr != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 2003,
				"msg":  "删除失败",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "删除成功",
				"data": "新的好友列表发起删除的ToUser",
			})
			// 单方面删除，不需要给被删除的from推送好友列表

		}
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 2004,
			"msg":  "查询记录失败，不存在该关系记录",
		})
	}

}
