package dao

import "github.com/jinzhu/gorm"


var(
	DB *gorm.DB
)


func InitMysql() (err error){
	// win密码6位 dsn:="root:123456@(127.0.0.1:3306)/chat?charset=utf8mb4&parseTime=True&loc=Local"
	// mac密码8位
	dsn:="root:12345678@(127.0.0.1:3306)/chat?charset=utf8mb4&parseTime=True&loc=Local"
	DB,err=gorm.Open("mysql",dsn)
	if err != nil {
		return
	}
	return DB.DB().Ping()
}

func Close(){
	DB.Close()
}
