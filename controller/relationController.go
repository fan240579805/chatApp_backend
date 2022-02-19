package controller

import (
	"chatApp_backend/model"
	_type "chatApp_backend/type"
	"chatApp_backend/ws"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type addFriendReqParams struct {
	Username string
	Fromid   string
	Toid     string
}

// AddFriendReq from用户像to用户发起好友请求（to还没同意）
// 此时的 from to是明确的
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
		c.JSON(http.StatusOK, gin.H{
			"code": 2003,
			"msg":  addFriendERR.Error(),
		})
	} else {
		rightRelation, _ := model.GetRightRelationRecord(params.Fromid, params.Toid)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "添加请求成功",
			"data": gin.H{
				"status": rightRelation.Status,
			},
		})
		// 执行推送给to用户好友请求逻辑 pushToUser
		PushFriendReq2user(rightRelation, rightRelation.From, rightRelation.To)
	}
}

// AcceptFriendReq to接受from用户发起的好友请求，将 status 设为1
// 此时的 from to 是明确的
func AcceptFriendReq(c *gin.Context) {
	var relation model.Relation
	c.ShouldBindJSON(&relation)

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
			"data": FormatFriendList(relation.To, c),
		})
		rightRelation, _ := model.GetRightRelationRecord(relation.From, relation.To)
		// pushFromUser 将新的好友列表推送给发起好友请求的FromUser
		PushFriendReq2user(rightRelation, relation.To, relation.From)
	}
}

// RejectFriendReq to拒绝from用户发起的好友请求，将该关系记录直接删除，保证下次添加是新的fromuser
// 此时的 from to 是明确的
func RejectFriendReq(c *gin.Context) {
	var relation model.Relation
	c.ShouldBindJSON(&relation)
	rightRelation, _ := model.GetRightRelationRecord(relation.From, relation.To)
	fmt.Println("拒绝rigth", rightRelation)

	modifyErr := model.DeleteRelation(rightRelation)
	if modifyErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 2003,
			"msg":  "拒绝失败",
		})
	} else {
		fmt.Println("拒绝", rightRelation.To)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "拒绝成功",
			"data": FormatFriendList(rightRelation.To, c),
		})
		// pushFromUser 将新的好友列表推送给发起好友请求的FromUser
		PushFriendReq2user(rightRelation, rightRelation.To, rightRelation.From)
	}
}

// GetFriendList 获取好友列表
func GetFriendList(c *gin.Context) {
	userid, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 2003,
			"msg":  "auth中间件获取userid失败",
		})
	} else {
		// 处理一下friendList并返回给前端
		friendList := FormatFriendList(userid.(string), c)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取好友列表成功",
			"data": friendList,
		})
	}
}

// DeleteFriendReq 某用户发起的删除好友请求，
// 根据发起删除的用户是to还是from来为 status 设 2：from删除to 还是 3：to删除from
func DeleteFriendReq(c *gin.Context) {
	// 删除好友时并不知道到底 发起删除的用户是from还是被删除人是from
	var unknownRecord model.Relation
	c.ShouldBindJSON(&unknownRecord)
	// 先判断一下， unknownRecord.From: 主动发起删除的用户
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
				"data": FormatFriendList(rightRelation.From, c),
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
				"data": FormatFriendList(rightRelation.To, c),
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

// TakeBlackReq 某用户发起的删除好友请求，
// 根据发起拉黑的用户是to还是from来为 status 设 4：from拉黑了to  还是 5：to拉黑了from
func TakeBlackReq(c *gin.Context) {
	// 拉黑时并不知道到底 发起拉黑的用户是from还是被拉黑人是from
	var unknownRecord model.Relation
	c.ShouldBindJSON(&unknownRecord)
	// 先判断一下
	rightRelation, _ := model.GetRightRelationRecord(unknownRecord.From, unknownRecord.To)
	if unknownRecord.From == rightRelation.From {
		// 发起拉黑的用户就是  发起好友请求的人（from） status = 4
		modifyErr := model.ModifyStatus(rightRelation.From, rightRelation.To, 4)
		if modifyErr != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 2003,
				"msg":  "拉黑失败",
			})
		} else {
			// 推送给发起拉黑与被拉黑用户双方，因为前端双方都需要这个status进行拦截黑名单发送聊天操作
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "拉黑成功",
				"data": FormatFriendList(rightRelation.From, c),
			})
			// 推送给被拉黑用户to
			PushFriendReq2user(rightRelation, rightRelation.From, rightRelation.To)
		}
	} else if unknownRecord.From == rightRelation.To {
		// 发起拉黑的用户是  一开始接受好友请求的人（to） status = 5
		modifyErr := model.ModifyStatus(rightRelation.From, rightRelation.To, 5)
		if modifyErr != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 2003,
				"msg":  "拉黑失败",
			})
		} else {
			// 推送给发起拉黑与被拉黑用户双方，因为前端双方都需要这个status进行拦截黑名单发送聊天操作
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "拉黑成功",
				"data": FormatFriendList(rightRelation.To, c),
			})
			// 推送给被拉黑用户from
			PushFriendReq2user(rightRelation, rightRelation.To, rightRelation.From)
		}
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 2004,
			"msg":  "查询记录失败，不存在该关系记录",
		})
	}
}

// CancelBlack 取消拉黑
func CancelBlack(c *gin.Context) {
	// 拉黑时并不知道到底 发起拉黑的用户是from还是被拉黑人是from
	var unknownRecord model.Relation
	c.ShouldBindJSON(&unknownRecord)
	// 先判断一下
	rightRelation, _ := model.GetRightRelationRecord(unknownRecord.From, unknownRecord.To)
	if unknownRecord.From == rightRelation.From {
		// 发起取消拉黑的用户就是  发起好友请求的人（from） status = 1
		modifyErr := model.ModifyStatus(rightRelation.From, rightRelation.To, 1)
		if modifyErr != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 2003,
				"msg":  "拉黑失败",
			})
		} else {
			// 推送给发起拉黑与被拉黑用户双方，因为前端双方都需要这个status进行拦截黑名单发送聊天操作
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "拉黑成功",
				"data": FormatFriendList(rightRelation.From, c),
			})
			// 推送给被拉黑用户to
			PushFriendReq2user(rightRelation, rightRelation.From, rightRelation.To)
		}
	} else if unknownRecord.From == rightRelation.To {
		// 发起拉黑的用户是  一开始接受好友请求的人（to） status = 1
		modifyErr := model.ModifyStatus(rightRelation.From, rightRelation.To, 1)
		if modifyErr != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 2003,
				"msg":  "取消拉黑失败",
			})
		} else {
			// 推送给发起拉黑与被拉黑用户双方，因为前端双方都需要这个status进行拦截黑名单发送聊天操作
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "取消拉黑成功",
				"data": FormatFriendList(rightRelation.To, c),
			})
			// 推送给被拉黑用户from
			PushFriendReq2user(rightRelation, rightRelation.To, rightRelation.From)
		}
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 2004,
			"msg":  "查询记录失败，不存在该关系记录",
		})
	}
}

// GetBlackStatus 获取当前用户是否拉黑对方
func GetBlackStatus(c *gin.Context) {
	FromID:=c.Query("from")
	ToID:=c.Query("to")

	rightRelation, _ := model.GetRightRelationRecord(FromID,ToID)
	log.Println("getblack",FromID)
	if FromID == rightRelation.From {
		if rightRelation.Status == 4 {
			// 当前用户确实拉黑了对方
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "",
				"data": gin.H{
					"isblack": true,
				},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "",
				"data": gin.H{
					"isblack": false,
				},
			})
		}
	} else {
		// 此时查询是否拉黑对方的user在relation中是接受方，所以查询是否为5
		if rightRelation.Status == 5 {
			// 当前用户确实拉黑了对方
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "",
				"data": gin.H{
					"isblack": true,
				},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "",
				"data": gin.H{
					"isblack": false,
				},
			})
		}
	}
}

// FormatFriendList 格式化获得 friendList ;  userid为当前登录用户
func FormatFriendList(userid string, c *gin.Context) []*_type.Friend {
	relationList, err := model.SelectFriends(userid)
	var friendList = make([]*_type.Friend, len(relationList))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 2003,
			"msg":  "获取好友列表失败",
		})
	} else {
		for i, relation := range relationList {
			if relation.From == userid {
				// 此时自己是发起添加请求的
				// 不要把自己的信息录入 friendList
				userProfile, _ := model.SelectUser(relation.To)
				friendList[i] = &_type.Friend{
					FriendProfile: userProfile,
					AddTime:       relation.CreatedAt,
					Status:        relation.Status,
					IsMaster:      true,
				}
			} else if relation.To == userid {
				// 此时自己是被添加的
				// 不要把自己的信息录入 friendList
				userProfile, _ := model.SelectUser(relation.From)
				friendList[i] = &_type.Friend{
					FriendProfile: userProfile,
					AddTime:       relation.CreatedAt,
					Status:        relation.Status,
					IsMaster:      false,
				}
			}
		}
	}

	return friendList
}

// PushFriendReq2user 将要推送的friend数据流传入chan给manager处理
func PushFriendReq2user(rightRelation model.Relation, newFriendUserid string, pushedUserid string) {
	userProfile, _ := model.SelectUser(newFriendUserid)
	var isMaster bool
	if rightRelation.From == pushedUserid {
		// 此时被推送用户是主动发起好友的人
		isMaster = true
	} else {
		isMaster = false
	}
	friend := &_type.Friend{
		FriendProfile: userProfile,
		AddTime:       rightRelation.CreatedAt,
		Status:        rightRelation.Status,
		IsMaster:      isMaster,
	}
	pushedObj := &_type.BePushedFriend{
		DataType:   "friend",
		BePushedID: pushedUserid,
		Friend:     *friend,
	}

	FriendChan, err := json.Marshal(pushedObj)
	if err != nil {
		log.Println("解析推送好友列表出错")
	}
	ws.Manager.FriendChan <- FriendChan
}
