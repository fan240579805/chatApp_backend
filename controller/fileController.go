package controller

import (
	"chatApp_backend/_const"
	"chatApp_backend/common"
	"chatApp_backend/model"
	"chatApp_backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"time"
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
	_, fileHeader := SaveImageToDisk(c, _const.AVATAR_PATH)
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

// UploadChatImage 上传聊天图片接口, 并转发给被发送图片的用户
func UploadChatImage(c *gin.Context) {
	userid, _ := c.Get("userID")
	chatID := c.PostForm("chatID")
	recipient := c.PostForm("recipient")
	sender := c.PostForm("sender")

	chatImgPath, _ := SaveImageToDisk(c, _const.CHAT_IMG_PATH)
	chatImgUrl := _const.BASE_URL + "/api/showImg?imageName=" + chatImgPath
	//insertErr := model.InsertFile(chatImgUrl, userid.(string), "img", "图片")
	//if insertErr != nil {
	//	log.Println("file存入失败")
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": 2004,
	//		"msg":  "file存入失败",
	//	})
	//}
	message := model.Message{
		MsgID:     "msgID_" + utils.UniqueId(),
		Sender:    userid.(string),
		Recipient: recipient,
		Content:   chatImgUrl,
		SendTime:  time.Now().UnixMilli(),
		Type:      "img",
	}
	saveMsgErr := model.AddMessageRecord(message)
	if saveMsgErr != nil {
		log.Println("聊天图片存入数据库失败")
		c.JSON(http.StatusOK, gin.H{
			"code": 2004,
			"msg":  "聊天图片存入数据库失败",
		})
	}
	if saveMsgErr == nil {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "发送成功",
			"data": chatImgUrl,
		})
		// 更新最新消息
		common.ModifyRecentMsg(chatID, message)
		// 图片推送给自己
		PushChatMsg2User(chatID, sender, message)
		// 推送给别人
		PushChatMsg2User(chatID, recipient, message)
	}

}

// SaveImageToDisk 图片存储到服务器磁盘func
func SaveImageToDisk(c *gin.Context, diskPath string) (string, *multipart.FileHeader) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("解析formdata出错err：", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return "", nil
	}
	dst := path.Join(diskPath, fileHeader.Filename) //上传文件保存路径

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
