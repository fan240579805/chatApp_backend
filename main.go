package main

import (
	"chatApp_backend/controller"
	"chatApp_backend/dao"
	"chatApp_backend/middle"
	"chatApp_backend/model"
	"chatApp_backend/ws"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	//连接初始化数据库
	err := dao.InitMysql()
	if err != nil {
		panic(err)
		return
	}
	//延迟关闭
	defer dao.Close()

	//数据库自迁移结构体
	dao.DB.AutoMigrate(&model.User{}, &model.Message{}, &model.Relation{}, &model.File{}, &model.Chat{})

	r := gin.Default()
	////中间件解决跨域
	//r.Use(middle.Cors1())
	go ws.Manager.Start()
	apiGroup := r.Group("api")
	{
		apiGroup.GET("/ws", ws.WsHandler)
		apiGroup.POST("/Registerauth", controller.PostRegister)
		apiGroup.POST("/login", controller.PostLogin)
		apiGroup.POST("/logout", middle.JWTAuthMiddleware(), controller.LogOut)
		apiGroup.POST("/AuthToken", middle.JWTAuthMiddleware(), controller.GetUserInfo)
		apiGroup.POST("/resetPassword", controller.ResetPassword)

		apiGroup.GET("/userInfo", middle.JWTAuthMiddleware(), controller.GetUserInfo)
		apiGroup.POST("/updateUserInfo", middle.JWTAuthMiddleware(), controller.UpdateUserInfo)
		apiGroup.POST("/addFriendReq", middle.JWTAuthMiddleware(), controller.AddFriendReq)
		apiGroup.POST("/acceptFriendReq", middle.JWTAuthMiddleware(), controller.AcceptFriendReq)
		apiGroup.POST("/rejectFriendReq", middle.JWTAuthMiddleware(), controller.RejectFriendReq)
		apiGroup.POST("/deleteFriendReq", middle.JWTAuthMiddleware(), controller.DeleteFriendReq)
		apiGroup.POST("/bothDelFriend", middle.JWTAuthMiddleware(), controller.BothDelFriend)

		apiGroup.GET("/getFriendList", middle.JWTAuthMiddleware(), controller.GetFriendList)
		apiGroup.POST("/takeBlack", middle.JWTAuthMiddleware(), controller.TakeBlackReq)
		apiGroup.POST("/cancelBlack", middle.JWTAuthMiddleware(), controller.CancelBlack)
		apiGroup.GET("/getBlackStatus", middle.JWTAuthMiddleware(), controller.GetBlackStatus)
		apiGroup.GET("/getBlackList", middle.JWTAuthMiddleware(), controller.GetBlackList)
		apiGroup.POST("/canIChat", middle.JWTAuthMiddleware(), controller.CanIChat)
		apiGroup.POST("/resetUnread", middle.JWTAuthMiddleware(), controller.ResetUnread)

		apiGroup.POST("/searchUser/:username", middle.JWTAuthMiddleware(), controller.SearchUser)

		apiGroup.GET("/getChatList", middle.JWTAuthMiddleware(), controller.GetMineChatList)
		apiGroup.POST("/getMsgList", middle.JWTAuthMiddleware(), controller.GetMsgList)
		apiGroup.POST("/uploadChatImg", middle.JWTAuthMiddleware(), controller.UploadChatImage)
		apiGroup.POST("/uploadAudio", middle.JWTAuthMiddleware(), controller.UploadAudio)
		apiGroup.POST("/bothDelMsg", middle.JWTAuthMiddleware(), controller.BothDelMsg)

		apiGroup.POST("/makeChat", middle.JWTAuthMiddleware(), controller.MakeChat)
		apiGroup.POST("/toggleInTheChat", middle.JWTAuthMiddleware(), ws.JoinChatRoom)
		apiGroup.POST("/exitChat", middle.JWTAuthMiddleware(), ws.ExitChatRoom)

		//apiGroup.GET("/getChatList",middle.JWTAuthMiddleware(),controller.GetChat)
		apiGroup.POST("/modifyAvatar", middle.JWTAuthMiddleware(), controller.ModifyAvatar)
		apiGroup.GET("/showImg", controller.ShowImage)
		apiGroup.GET("/showAudio", controller.ShowAudio)

		apiGroup.POST("/sendEmailVcode", controller.SendVcode)

	}

	r.Run(":9998")
}
