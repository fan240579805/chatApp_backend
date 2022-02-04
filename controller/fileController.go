package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"path"
)




// ShowChatImage 展示图片接口
func ShowImage(c *gin.Context){
	imageName := c.Query("imageName")
	c.File(imageName)
}

// UploadImage 上传图片接口
func UploadImage(c *gin.Context) {
	_, fileHeader := SaveImage(c)
	// 存入数据库 待补充
	c.JSONP(200, gin.H{
		"Filename": fileHeader.Filename,
	})
}

// SaveImage 图片存储到服务器磁盘func
func SaveImage(c *gin.Context) (string, *multipart.FileHeader) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("解析formdata出错err：", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return "", nil
	}
	dst := path.Join("statics/images/", fileHeader.Filename) //上传文件保存路径
	fmt.Println(dst)

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
