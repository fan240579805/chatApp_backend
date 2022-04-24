package controller

import (
	"chatApp_backend/_const"
	"chatApp_backend/dao"
	"chatApp_backend/middle"
	"chatApp_backend/model"
	"chatApp_backend/utils"
	"chatApp_backend/ws"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"

	"gopkg.in/gomail.v2"
)

func PostRegister(c *gin.Context) {
	userid := "userid_" + utils.UniqueId()
	username := c.PostForm("UserName")
	email := c.PostForm("Email")
	nickname := c.PostForm("NickName")
	password := c.PostForm("PassWord")
	avatarPath, _ := SaveImageToDisk(c, _const.AVATAR_PATH)
	var u1 = &model.User{
		UserID:   userid,
		Username: username,
		NickName: nickname,
		Password: password,
		Avatar:   _const.BASE_URL + "/api/showImg?imageName=" + avatarPath,
		Email:    email,
	}
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
				"avatar":   u1.Avatar,
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
				"avatar":   u1.Avatar,
			},
		})

	}
}

// ResetPassword 忘记密码
func ResetPassword(c *gin.Context) {
	type ResetPwdParams struct {
		UserName string
		EmailCode string
		PassWord string
	}
	var params ResetPwdParams
	c.ShouldBindJSON(&params)
	username := params.UserName
	newPassword := params.PassWord
	emailVcode := params.EmailCode
	var user model.User

	selectErr := dao.DB.Where("username=?", username).Find(&user).Error
	if selectErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "没有该用户",
			"data": false,
		})
		return
	}

	if CheckVcode(user.UserID, emailVcode) {
		var modifyAction model.ModifyAction
		modifyAction.InfoAttr = "password"
		modifyAction.Playloads = newPassword
		_, modifyErr := model.ModifyChatUserInfo(user.UserID, &modifyAction)
		if modifyErr != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "更新失败，请重试",
				"data": false,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "更新成功",
			"data": true,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "验证码错误",
			"data": false,
		})
	}
}

func LogOut(c *gin.Context) {
	id, _ := c.Get("userID")
	ws.UserExit(id.(string))
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登出成功",
		"data": "",
	})
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
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "删除成功",
		})
	}
}

// SearchUser 根据账号搜索用户
func SearchUser(c *gin.Context) {
	searchUsername, _ := c.Params.Get("username")
	userid, _ := c.Get("userID")
	userInfo, err := model.FindUser(searchUsername)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
	}
	if userid.(string) == userInfo.UserID {
		// 用户搜索自己，返回自己
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "搜索成功",
			"data": gin.H{
				"userInfo": userInfo,
				"Status":   -99,
				"isMaster": true,
			},
		})
		return
	}
	// 先查一下，看看用户是否已经添加搜索的好友
	rightRelation, rightRelationErr := model.GetRightRelationRecord(userid.(string), userInfo.UserID)
	var isMaster bool
	if rightRelation.From == userid {
		// 此时发起搜索的用户是主动发起好友的人
		isMaster = true
	} else {
		isMaster = false
	}
	if rightRelationErr == nil && rightRelation.Status == -1 {
		fmt.Println("searchL", rightRelation)
		// 搜索出来的好友已经添加过了，额外传多一个status
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "已添加该用户",
			"data": gin.H{
				"userInfo": userInfo,
				"Status":   rightRelation.Status,
				"isMaster": isMaster,
			},
		})
	} else {
		var Status int
		if rightRelation.From == "" {
			// 两个用户从未建立关系
			Status = 0
		} else {
			Status = rightRelation.Status
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "搜索成功",
			"data": gin.H{
				"userInfo": userInfo,
				"Status":   Status,
				"isMaster": isMaster,
			},
		})
	}
}


// MailboxConf 邮箱配置
type MailboxConf struct {
	// 邮件标题
	Title string
	// 邮件内容
	Body string
	// 收件人列表
	RecipientList []string
	// 发件人账号
	Sender string
	// 发件人密码，QQ邮箱这里配置授权码
	SPassword string
	// SMTP 服务器地址， QQ邮箱是smtp.qq.com
	SMTPAddr string
	// SMTP端口 QQ邮箱是25
	SMTPPort int
}

type MailType struct {
	Vcode        string
	LastSendTime int64
}

var ClientMailMap = make(map[string]*MailType)

var maxTime = 60

type MailParams struct {
	UserName string
	Email string
}

func SendVcode(c *gin.Context) {
	var emailParams MailParams
	c.ShouldBindJSON(&emailParams)
	var user model.User
	selectErr := dao.DB.Debug().Where("username=?", emailParams.UserName).Find(&user).Error
	if selectErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "没有该用户",
			"data": false,
		})
		return
	}
	var userid = user.UserID

	if ClientMailMap[userid] == nil {
		ClientMailMap[userid] = &MailType{
			Vcode:        "",
			LastSendTime: 0,
		}
	}

	var curTime = time.Now().Unix()
	if curTime-ClientMailMap[userid].LastSendTime < int64(maxTime) {
		// 60s 只发送一次验证码
		return
	}
	var mailConf MailboxConf
	mailConf.Title = "验证"
	// 这里就是我们发送的邮箱内容，但是也可以通过下面的html代码作为邮件内容
	// mailConf.Body = "坚持才是胜利，奥里给"

	//这里支持群发，只需填写多个人的邮箱即可，我这里发送人使用的是QQ邮箱，所以接收人也必须都要是
	//QQ邮箱
	mailConf.RecipientList = []string{emailParams.Email}
	mailConf.Sender = `2953336033@qq.com`

	//这里QQ邮箱要填写授权码，网易邮箱则直接填写自己的邮箱密码，授权码获得方法在下面
	mailConf.SPassword = "qscozlfhzrtidhdc"

	//下面是官方邮箱提供的SMTP服务地址和端口
	// QQ邮箱：SMTP服务器地址：smtp.qq.com（端口：587）
	// 雅虎邮箱: SMTP服务器地址：smtp.yahoo.com（端口：587）
	// 163邮箱：SMTP服务器地址：smtp.163.com（端口：25）
	// 126邮箱: SMTP服务器地址：smtp.126.com（端口：25）
	// 新浪邮箱: SMTP服务器地址：smtp.sina.com（端口：25）

	mailConf.SMTPAddr = `smtp.qq.com`
	mailConf.SMTPPort = 587

	//产生六位数验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	ClientMailMap[userid].Vcode = vcode
	ClientMailMap[userid].LastSendTime = time.Now().Unix()

	//发送的内容
	html := fmt.Sprintf(`<div>
        <div>
            尊敬的用户，您好！
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>你本次的验证码为%s,为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
        </div>
        <div>
            <p>此邮箱为系统邮箱，请勿回复。</p>
        </div>    
    </div>`, vcode)

	m := gomail.NewMessage()

	// 第三个参数是我们发送者的名称，但是如果对方有发送者的好友，优先显示对方好友备注名
	m.SetHeader(`From`, mailConf.Sender, "chatApp官方")
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		log.Fatalf("Send Email Fail, %s", err.Error())
		return
	}
	log.Printf("Send Email Success")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码发送成功",
		"data": "",
	})
}

type MailCheckParams struct {
	InputCode string
}

func CheckVcode(userid string, InputCode string) bool {
	//uid, _ := c.Get("userID")
	//var userid = uid.(string)

	//var mailCheckStruct MailCheckParams
	//c.ShouldBindJSON(&mailCheckStruct)

	if ClientMailMap[userid].Vcode == InputCode {
		return true
	} else {
		return false
	}

}



