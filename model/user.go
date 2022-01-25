package model

import (
	"chatApp/dao"
)

var ImgUrl  string

type User struct{
	ID  int                 `form:"id"             gorm:"unique;not null"`
	Name string				`form:"name"           gorm:"unique"`
	Email string			`form:"email"          gorm:"unique;not null"`
	Password string			`form:"password"       gorm:"not null"`
	Telephone string		`form:"telephone"      gorm:"unique;not null" `
	Role string
	State bool
	IsTip bool
}

type EditUser struct {
	Name string
	Telephone string
	Email string
}

// 查询所有用户以及根据搜索结果查询用户
func AllUsers(page int , contentSize int,search string) ([]User , int ,error) {

	var users []User
	var count int
	var of int

	// 偏移量= （页数-1）* 每一页的数据量
	of = (page-1)*contentSize
	// 搜索数据库
	if search != "" {
		err:=dao.DB.Debug().Where("name=?",search).Select("id,name, email,telephone,role,state").
			Find(&users).Count(&count).Error
		if err != nil {
			return users , count ,err
		}else {
			return users , count ,nil
		}
	}else { // 搜索框没有东西

		// 总用户数量
		dao.DB.Find(&users).Count(&count)

		err:=dao.DB.Debug().Select("id,name, email,telephone,role,state").
			Offset(of).Limit(contentSize).Find(&users).Error
		if err != nil {
			return users , 0 , err
		}else {
			return users , count , nil
		}
	}

}

// 更新用户状态
func DbUpdateUserState(id int,state bool) (User , error) {
	var user User
	err:=dao.DB.Where("id=?",id).First(&user).Error
	if err != nil {
		return user,err
	}else {
		err:=dao.DB.Debug().Model(&user).Update("state",state).Error
		if err != nil {
			return user,err
		}
		return user,nil
	}
}

// 根据id查询套更新的用户信息
func SelectEditUser(id int) (User,error)  {
	var user User
	err:=dao.DB.Where("id=?",id).Select("id,name,email,telephone").First(&user).Error
	if err != nil {
		return user,err
	}
	return user,nil
}

// 根据id更新要更新的用户
func DbUpdateUser(id int, editUser EditUser) (User , error) {

	var u User
	err:=dao.DB.Where("id=?",id).First(&u).Error
	if err != nil {
		return u,err
	}else {
		err:=dao.DB.Debug().Model(&u).
			Updates(map[string]interface{}{"email":editUser.Email,"telephone":editUser.Telephone}).Error
		if err != nil {
			return u,err
		}
		return u,nil
	}
}

// 根据id删除用户
func DeleteUser(id int) error {
	err:=dao.DB.Where("id=?",id).Delete(&User{}).Error
	if err != nil {
		return err
	}else {
		return nil
	}
}


// 请求聊天用户列表
func GetUserChatList() ([]User,error) {
	var u1 []User
	err:=dao.DB.Select("id,name,state,is_tip").Find(&u1).Error
	if err != nil {
		return u1,err
	}
	return u1,nil
}

// 根据 id 查询 更新聊天用户的 已读未读状态
func ModifyChatUserState(id int, tipState bool) ( error ) {
	var user User
	err:=dao.DB.Where("id=?",id).First(&user).Error
	if err != nil {
		return err
	}
	err=dao.DB.Model(&user).Update("is_tip",tipState).Error
	if err != nil {
		return err
	}
	return nil
}