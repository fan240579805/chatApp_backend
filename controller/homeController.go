package controller

import (
	"chatApp/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserState struct {
	State bool  `json:"state"`
}
func Showhello(c * gin.Context){
	c.JSON(200,gin.H{
		"hello":"ni好",
	})
}

func GetUserList(c * gin.Context)  {
	pageS:=c.Query("page")
	contentSizeS:=c.Query("pageSize")
	// 搜索框中的内容
	search := c.Query("search")

	page,_ := strconv.Atoi(pageS)
	contentSize,_ := strconv.Atoi(contentSizeS)

	// 分页查询
	users,count,err:=model.AllUsers(page,contentSize,search)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,"错误")
	}

	c.JSON(200,gin.H{
		"contentSize":contentSize,
		"data":gin.H{
			// 用户总数
			"contentCount":count,

			// 用户列表
			"userlist":users,
		},
	})
}

func UpdateUserState( c * gin.Context)  {
	// 获取id并转化为int
	idS,_:=c.Params.Get("id")
	id,_ := strconv.Atoi(idS)

	// 获取state并转化为bool
	var s1 UserState
	c.ShouldBind(&s1)


	// 更新数据库
	user,err:=model.DbUpdateUserState(id,s1.State)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,gin.H{"err":"更新失败"})
	}else {
		c.JSON(http.StatusOK,gin.H{
			"code":200,
			"data":user,
		})
	}
}
