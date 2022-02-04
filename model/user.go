package model

import (
	"chatApp/dao"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

var ImgUrl string

type User struct {
	ID         int       `gorm:"column:id;unique;not null;primary_key;AUTO_INCREMENT"`
	UserID     string    `gorm:"column:userid;unique;not null"`
	Username   string    `gorm:"column:username;unique"`
	Password   string    `gorm:"column:password;not null"`
	Avatar     string    `gorm:"column:avatar;default:null"`
	Email      string    `gorm:"column:email;default:null"`
	NickName   string    `gorm:"column:nickname"`
	ChatList   string    `gorm:"column:chatlist"`
	CreateTime time.Time `gorm:"column:createtime;default:null" json:"createtime"`
	UpdateTime time.Time `gorm:"column:updatetime;default:null" json:"updatetime"`
}

type ChatList []string

type ModifyAction struct {
	InfoAttr  string
	Playloads string
}

func MakeChatList(db *gorm.DB) (string, error) {
	var chats = &ChatList{"123", "123asdasds"}
	bs, err := json.Marshal(chats)
	user := &User{
		ID:       0,
		UserID:   "12312313asda",
		Username: "asdasd",
		Password: "123",
		Avatar:   "asd",
		NickName: "aaa",
		ChatList: string(bs),
	}
	_ = db.Debug().Create(user).Error
	fmt.Println("create user success")
	if err != nil {
		return "出错啦", err
	} else {
		return string(bs), nil
	}
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
func SelectUser(id string) (User, error) {
	var user User
	err := dao.DB.Debug().Where("userid=?", id).Select("id,userid,username,nickname,avatar,email").First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// ModifyChatUserInfo 根据 id 查询出用户并更新相应信息; modifyAction-要更新的数据库字段命及参数
// 更新chatList字段是前端回传json格式的数组就行
func ModifyChatUserInfo(id string, action *ModifyAction) (User, error) {
	var user User
	err := dao.DB.Where("userid=?", id).First(&user).Error
	if err != nil {
		return user, err
	}
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

// findAllUsers 查询所有用户以及根据搜索结果查询用户
func findAllUsers(page int, contentSize int, search string) ([]User, int, error) {
	var users []User
	var count int
	var of int

	// 偏移量= （页数-1）* 每一页的数据量
	of = (page - 1) * contentSize
	// 搜索数据库
	if search != "" {
		err := dao.DB.Debug().Where("username=?", search).Select("id,userid,username,nickname,avatar").
			Find(&users).Count(&count).Error
		if err != nil {
			return users, count, err
		} else {
			return users, count, nil
		}
	} else { // 搜索框没有东西

		// 总用户数量
		dao.DB.Find(&users).Count(&count)

		err := dao.DB.Debug().Select("id,name, email,telephone,role,state").
			Offset(of).Limit(contentSize).Find(&users).Error
		if err != nil {
			return users, 0, err
		} else {
			return users, count, nil
		}
	}

}
