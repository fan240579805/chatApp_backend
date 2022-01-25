package controller

import (
	"chatApp/dao"
	"chatApp/middle"
	"chatApp/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Pos1tRegister(c *gin.Context) {

}

func PostRegister(c *gin.Context) {
	var u1 model.User
	c.ShouldBind(&u1)
	if u1.Name == "admin" {
		u1.Role = "超级管理员"
		u1.State = true
	} else {
		u1.Role = "普通用户"
		u1.State = false
	}
	err := dao.DB.Debug().Create(&u1).Error
	if err != nil {
		//注册失败，用户名或密码或电话号码重复
		//返回错误代码用这个http.StatusUnprocessableEntity=422才能使得前端axios catch到err
		c.JSON(http.StatusUnprocessableEntity, gin.H{"Err": "用户名或邮箱重复"})
		return
	} else {
		//注册成功
		fmt.Println(u1)

		token, err := middle.CreateToken(u1.Name, u1.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 444,
				"msg":  "create token failed",
			})
		}
		c.JSON(200, gin.H{
			"code": 2002,
			"msg":  "register success",
			"data": gin.H{
				"username": u1.Name,
				"token":    token,
			},
		})

	}
}

func PostLogin(c *gin.Context) {
	var u1 model.User
	c.ShouldBind(&u1)
	err := dao.DB.Where("email=? AND password=?", u1.Email, u1.Password).First(&u1).Error
	if err != nil {
		//登录失败，用户名或密码或电话号码重复或用户不存在
		//返回错误代码用这个http.StatusUnprocessableEntity=422才能使得前端axios catch到err
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"Err": "该账号不存在"})
		return
	} else {
		//登录成功
		fmt.Println(u1)

		token, err := middle.CreateToken(u1.Name, u1.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 444,
				"msg":  "create token failed",
			})
		}
		c.JSON(200, gin.H{
			"code": 2002,
			"msg":  "login success",
			"data": gin.H{
				"user":  u1.ID,
				"token": token,
			},
		})

	}
}

func GetUserInfo(c *gin.Context) {
	userName, _ := c.Get("username")
	id, _ := c.Get("userID")

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"userName": userName,
		"uid":      id,
	})
}

func GetMenuList(c *gin.Context) {
	//menuList:=make(map[string]interface{},10)

	m := model.NewMenu()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": m,
	})
}

func GetEditUser(c *gin.Context) {
	// 获取id并转化为int
	idS, _ := c.Params.Get("id")
	id, _ := strconv.Atoi(idS)

	user, err := model.SelectEditUser(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"Err": "查找不到该用户"})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"name":      user.Name,
			"email":     user.Email,
			"telephone": user.Telephone,
		},
	})
}

func UpdateUser(c *gin.Context) {
	// 获取id并转化为int
	idS, _ := c.Params.Get("id")
	id, _ := strconv.Atoi(idS)

	var u1 model.EditUser
	c.ShouldBind(&u1)
	//state,_=strconv.ParseBool(s)

	// 更新数据库
	user, err := model.DbUpdateUser(id, u1)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": "更新失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": user,
		})
	}
}

func DeleteUser(c *gin.Context) {
	// 获取id并转化为int
	idS, _ := c.Params.Get("id")
	id, _ := strconv.Atoi(idS)

	err := model.DeleteUser(id)
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

func GetChatList(c *gin.Context) {

	users, err := model.GetUserChatList()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 404,
			"msg":  "请求聊天用户失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"users": users,
	})
}

func TurnState(c *gin.Context) {
	// 获取id并转化为int
	idS, _ := c.Params.Get("id")
	id, _ := strconv.Atoi(idS)

	err := model.ModifyChatUserState(id, false)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 404,
			"msg":  err.Error(),
		})
	}

	users, err := model.GetUserChatList()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 404,
			"msg":  "请求聊天用户失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"users": users,
	})
}
