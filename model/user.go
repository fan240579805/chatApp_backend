package model

import (
	"chatApp_backend/_const"
	"chatApp_backend/dao"
	"encoding/json"
	"errors"
	"log"
	"time"
)

var ImgUrl string

type User struct {
	ID        int       `gorm:"column:id;unique;not null;primary_key;AUTO_INCREMENT"`
	UserID    string    `gorm:"column:userid;unique;not null"`
	Username  string    `gorm:"column:username;unique"`
	Password  string    `gorm:"column:password;not null"`
	Avatar    string    `gorm:"column:avatar;default:null"`
	Email     string    `gorm:"column:email;default:null"`
	NickName  string    `gorm:"column:nickname"`
	ChatList  string    `gorm:"column:chatlist"`
	CreatedAt time.Time `gorm:"column:createdat;default:null" json:"createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat;default:null" json:"updatedat"`
}
type UserInfo struct {
	UserID   string
	Username string
	Avatar   string
	Email    string
	NickName string
}

type ChatList []string

type ModifyAction struct {
	InfoAttr  string
	Playloads string
}

func checkUser(user *User) error {
	var compareUserName User
	var compareEmail User
	dao.DB.Where("email=?", user.Email).First(&compareEmail)
	dao.DB.Where("username=?", user.Username).First(&compareUserName)
	if compareUserName.UserID == "" && compareEmail.UserID == "" {
		return nil
	} else if compareUserName.UserID == "" && compareEmail.UserID != "" {
		return errors.New("邮箱已被注册")
	} else if compareUserName.UserID != "" && compareEmail.UserID == "" {
		return errors.New("账号已被注册")
	} else {
		return errors.New("账号和邮箱都已被注册")
	}
}

// AddUser 插入新用户到数据库,对userName，email去重
func AddUser(user User) error {
	checkErr := checkUser(&user)
	if checkErr != nil {
		return checkErr
	}
	err := dao.DB.Debug().Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// SelectUser 根据id查询用户信息
func SelectUser(id string) (*UserInfo, error) {
	var userInfo *UserInfo
	var user User
	err := dao.DB.Debug().Where("userid=?", id).Select("id,userid,username,nickname,avatar,email").First(&user).Error
	if err != nil {
		return userInfo, err
	}
	userInfo = &UserInfo{
		UserID:   user.UserID,
		Username: user.Username,
		Avatar:   user.Avatar,
		Email:    user.Email,
		NickName: user.NickName,
	}
	return userInfo, nil
}

// ModifyChatUserInfo 根据 id 查询出用户并更新相应信息; modifyAction-要更新的数据库字段命及参数
// 更新chatList字段是前端回传json格式的数组就行
func ModifyChatUserInfo(id string, action *ModifyAction) (User, error) {
	log.Println("asdasdasad",id,action)
	var user User
	err := dao.DB.Where("userid=?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	if action.InfoAttr == "avatar" {
		var avatarFileName = action.Playloads
		action.Playloads = _const.BASE_URL + "/api/showImg?imageName=" + _const.AVATAR_PATH + avatarFileName
	}
	log.Println("asdasdasad",id,action)
	err = dao.DB.Model(&user).Update(action.InfoAttr, action.Playloads).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// DeleteUser 根据id删除用户
func DeleteUser(id string) error {
	err := dao.DB.Where("userid=?", id).Delete(&User{}).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

// GetUserChatList 获取聊天会话列表
func GetUserChatList(id string) ([]Chat, error) {
	var chatListJSON string
	var chatList []Chat
	err := dao.DB.Debug().Where("userid=", id).Select("chatlist").First(&chatListJSON).Error
	// 将存入数据库的chatList JSON字符串解析为字节流
	chatListBtyes, _ := json.Marshal(chatListJSON)
	// 再将字节流反序列化成为chatList数据结构
	json.Unmarshal(chatListBtyes, &chatList)
	if err != nil {
		return chatList, err
	}
	return chatList, nil
}

// FindUser 根据搜索结果查询用户
func FindUser(search string) (*UserInfo, error) {
	// 搜索数据库
	var userInfo *UserInfo
	var user User
	err := dao.DB.Debug().Where("username=?", search).Select("id,userid,username,nickname,avatar").
		First(&user).Error
	if err != nil {
		return userInfo, err
	} else {
		userInfo = &UserInfo{
			UserID:   user.UserID,
			Username: user.Username,
			Avatar:   user.Avatar,
			Email:    user.Email,
			NickName: user.NickName,
		}
		return userInfo, nil
	}
}
