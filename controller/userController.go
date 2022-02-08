package controller

import (
	"chatApp_backend/_const"
	"chatApp_backend/dao"
	"chatApp_backend/middle"
	"chatApp_backend/model"
	"chatApp_backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostRegister(c *gin.Context) {
	userid := "userid_" + utils.UniqueId()
	username := c.PostForm("UserName")
	email := c.PostForm("Email")
	nickname := c.PostForm("NickName")
	password := c.PostForm("PassWord")
	avatarPath, _ := SaveImage(c)
	var u1 = &model.User{
		UserID:   userid,
		Username: username,
		NickName: nickname,
		Password: password,
		Avatar:   _const.BASE_URL+"/api/showImg?imageName=" + avatarPath,
		Email:    email,
	}
	fmt.Println(u1)
	err := model.AddUser(*u1)
	if err != nil {
		fmt.Println(err.Error())
		//注册失败，用户名或密码或电话号码重复
		//返回错误代码用这个http.StatusUnprocessableEntity=422才能使得前端axios catch到err
		c.JSON(http.StatusOK, gin.H{
			"code": "422",
			"msg":  err.Error(),
		})
		return
	} else {
		//注册成功
		fmt.Printf("注册成功: %s", u1.Username)

		token, err := middle.CreateToken(u1.Username, u1.UserID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 444,
				"msg":  "create token failed",
			})
		}
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "注册成功",
			"data": gin.H{
				"username": u1.Username,
				"nickname": u1.NickName,
				"userID":   u1.UserID,
				"token":    token,
			},
		})

	}
}

func PostLogin(c *gin.Context) {
	var u1 model.User
	c.ShouldBind(&u1)
	err := dao.DB.Debug().Where("username=? AND password=?", u1.Username, u1.Password).First(&u1).Error
	if err != nil {
		//登录失败，用户名或密码或电话号码重复或用户不存在
		//返回错误代码用这个http.StatusUnprocessableEntity=422才能使得前端axios catch到err
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"Err":  "该账号不存在"})
		return
	} else {
		//登录成功
		fmt.Printf("登录成功: %s", u1.Username)

		token, err := middle.CreateToken(u1.Username, u1.UserID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 444,
				"msg":  "create token failed",
			})
		}
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "登录成功",
			"data": gin.H{
				"username": u1.Username,
				"userID":   u1.UserID,
				"token":    token,
			},
		})

	}
}

// GetUserInfo 获取用户基本信息：昵称，头像等
func GetUserInfo(c *gin.Context) {
	id, _ := c.Get("userID")
	userInfo, err := model.SelectUser(id.(string))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "user dismiss",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": userInfo,
		})
	}

}

// UpdateUserInfo 更新用户的info基本信息，
func UpdateUserInfo(c *gin.Context) {
	// 获取id并转化为int
	//userid, _ := c.Params.Get("userid")
	ids, _ := c.Get("userID")
	var action model.ModifyAction
	c.ShouldBindJSON(&action)

	// 更新数据库
	newUserInfo, err := model.ModifyChatUserInfo(ids.(string), &action)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": "更新失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新成功",
			"data": gin.H{
				"Username": newUserInfo.Username,
				"NickName": newUserInfo.NickName,
				"Avatar":   newUserInfo.Avatar,
				"Email":    newUserInfo.Email,
			},
		})
	}
}

// DeleteUser 根据id删除用户
func DeleteUser(c *gin.Context) {
	// 获取id并转化为int
	idS, _ := c.Params.Get("userid")

	err := model.DeleteUser(idS)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 404,
			"msg":  "删除失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// GetChatList 根据userid获取对应的聊天会话列表
func GetChatList(c *gin.Context) {
	idS, _ := c.Params.Get("userid")
	chatList, err := model.GetUserChatList(idS)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 404,
			"msg":  "请求聊天列表失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"chatList": chatList,
		},
	})
}
