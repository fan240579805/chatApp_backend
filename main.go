package main

import (
	"chatApp_backend/controller"
	"chatApp_backend/dao"
	"chatApp_backend/middle"
	"chatApp_backend/model"
	"chatApp_backend/ws"
	"fmt"
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
		apiGroup.POST("/AuthToken", middle.JWTAuthMiddleware(), controller.GetUserInfo)
		apiGroup.GET("/userInfo", middle.JWTAuthMiddleware(), controller.GetUserInfo)
		apiGroup.POST("/updateUserInfo", middle.JWTAuthMiddleware(), controller.UpdateUserInfo)
		apiGroup.POST("/addFriendReq", middle.JWTAuthMiddleware(), controller.AddFriendReq)
		apiGroup.POST("/acceptFriendReq", middle.JWTAuthMiddleware(), controller.AcceptFriendReq)
		apiGroup.POST("/deleteFriendReq", middle.JWTAuthMiddleware(), controller.DeleteFriendReq)
		apiGroup.GET("/getFriendList", middle.JWTAuthMiddleware(), controller.GetFriendList)
		apiGroup.POST("/takeBlack", middle.JWTAuthMiddleware(), controller.TakeBlackReq)
		//apiGroup.GET("/userChatlist",middle.JWTAuthMiddleware(),controller.GetChatList)
		//apiGroup.PUT("/modifyMsgState",middle.JWTAuthMiddleware(),controller.ModifyMsgReadState)
		//apiGroup.GET("/getChatList",middle.JWTAuthMiddleware(),controller.GetChat)
		//apiGroup.POST("/saveImg",middle.JWTAuthMiddleware(),controller.SaveImg)
		apiGroup.GET("/showImg", controller.ShowImage)
		apiGroup.GET("/AJAX/:id", testAJAX)
		apiGroup.POST("/uploadFile", controller.UploadImage)
	}

	r.Run(":9998")
}
func testAJAX(context *gin.Context) {
	id, _ := context.Params.Get("id")
	fmt.Println(id)
	context.JSON(400, gin.H{
		"msg": 123,
	})
}

type file struct {
	userid string
	path   string
}
