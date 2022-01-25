package main

import (
	"chatApp/controller"
	"chatApp/dao"
	"chatApp/middle"
	"chatApp/model"
	"chatApp/ws"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main()  {

	//连接初始化数据库
	err:=dao.InitMysql()
	if err != nil {
		panic(err)
		return
	}
	//延迟关闭
	defer dao.Close()

	//数据库自迁移结构体
	dao.DB.AutoMigrate(&model.User{},&model.Message{})

	r:=gin.Default()
	//中间件解决跨域
	r.Use(middle.Cors1())

	go ws.Manager.Start()

	apiGroup:=r.Group("api")
	{
		apiGroup.POST("/Registerauth",controller.PostRegister)
		apiGroup.POST("/login",controller.PostLogin)
		apiGroup.GET("/userInfo",middle.JWTAuthMiddleware(),controller.GetUserInfo)
		apiGroup.GET("/menuList",middle.JWTAuthMiddleware(),controller.GetMenuList)
		apiGroup.GET("/Userlist",middle.JWTAuthMiddleware(),controller.GetUserList)
		apiGroup.PUT("/updateUserState/:id",middle.JWTAuthMiddleware(),controller.UpdateUserState)
		apiGroup.GET("/GetEditUser/:id",middle.JWTAuthMiddleware(),controller.GetEditUser)
		apiGroup.PUT("/updateUser/:id",middle.JWTAuthMiddleware(),controller.UpdateUser)
		apiGroup.DELETE("/deleteUser/:id",middle.JWTAuthMiddleware(),controller.DeleteUser)
		apiGroup.GET("/userChatlist",middle.JWTAuthMiddleware(),controller.GetChatList)
		apiGroup.GET("/ws",ws.WsHandler)
		apiGroup.PUT("/modifyMsgState",middle.JWTAuthMiddleware(),controller.ModifyMsgReadState)
		apiGroup.GET("/getChatList",middle.JWTAuthMiddleware(),controller.GetChat)
		apiGroup.POST("/saveImg",middle.JWTAuthMiddleware(),controller.SaveImg)
		apiGroup.GET("/showImg",controller.ShowChatImage)
		apiGroup.GET("/hello",controller.Showhello)
		apiGroup.POST("/AJAX/:id",testAJAX)
		apiGroup.GET("/JSONP",testJsonp)
	}

	r.Run(":9998")
}
func testAJAX(context *gin.Context) {
	id, _ :=context.Params.Get("id")
	fmt.Println(id);
	context.JSON(400,gin.H{
		"msg":123,
	})
}
func testJsonp(c *gin.Context) {
	fmt.Println("jsonp")
	c.JSONP(200,gin.H{
		"msg":"jsonp",
	})
}