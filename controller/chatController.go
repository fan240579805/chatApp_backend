package controller

import (
	"chatApp/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func ModifyMsgReadState(c*gin.Context)  {
	msgFrom:=c.Query("msgFrom")
	msgTo:=c.Query("msgTo")
	err:=model.ModifyMsgState(msgFrom,msgTo)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,gin.H{
			"code":444,
			"state":"false",
		})
	}
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"state":"success",
	})
}

func GetChat(c * gin.Context)  {
	msgs,err:=model.GetChatContent()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity,gin.H{
			"code":444,
			"state":"false",
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"chatList":msgs,
	})
}

func asdSaveImg( c *gin.Context)  {
	f,err:=c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		c.JSON(200,gin.H{
			"status":"err",
		})
		return
	}
	dst:=path.Join("ttt/statics/",f.Filename)//上传文件保存路径
	err=c.SaveUploadedFile(f,dst)
	//url="./"+dst;
	model.ImgUrl="http://127.0.0.1:9998/api/showImg?imageName="+dst
	if err != nil {
		c.JSON(200,gin.H{"status":err.Error()})
		return
	}else {
		c.JSON(200,gin.H{
			"status":"upload ok!",
			"url":model.ImgUrl,
		})
	}
}


