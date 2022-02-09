package controller

import (
	"chatApp_backend/_const"
	"chatApp_backend/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"path"
)

// ShowImage 展示图片接口
func ShowImage(c *gin.Context) {
	imageName := c.Query("imageName")
	c.File(imageName)
}

// ModifyAvatar 更新并保存头像路径
func ModifyAvatar(c *gin.Context) {
	uid := c.PostForm("userid")
	InfoAttr := c.PostForm("InfoAttr")
	_, fileHeader := SaveAvatarImage(c)
	newUserInfo, err := model.ModifyChatUserInfo(uid, &model.ModifyAction{
		InfoAttr:  InfoAttr,
		Playloads: fileHeader.Filename,
	})
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

// UploadImage 上传图片接口
func UploadImage(c *gin.Context) {
	//_, fileHeader := SaveImage(c)
	// 存入数据库 待补充
	//c.JSONP(200, gin.H{
	//	"Filename": fileHeader.Filename,
	//})
}

// SaveAvatarImage 图片存储到服务器磁盘func
func SaveAvatarImage(c *gin.Context) (string, *multipart.FileHeader) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("解析formdata出错err：", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return "", nil
	}
	dst := path.Join(_const.AVATAR_PATH, fileHeader.Filename) //上传文件保存路径

	saveError := c.SaveUploadedFile(fileHeader, dst)
	if saveError != nil {
		fmt.Printf("存储服务器时出错err：", saveError.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": saveError.Error(),
		})
		return "", nil
	}
	return dst, fileHeader
}
