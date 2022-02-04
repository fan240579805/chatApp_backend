package model

import "time"

/**
status：关系状态码
1: from主动加to，并成为好友
2：to主动加from，并成为好友
3：from删除to
4：to删除from
5: from拉黑了to
6: to拉黑了from
*/
type Relation struct {
	Status     int       `gorm:"column:status;not null"`
	From       string    `gorm:"column:from"` // 发起好友关系的人的唯一id
	To         string    `gorm:"column:to"`   // 接受好友关系的人的唯一id
	CreateTime time.Time `gorm:"column:createtime;default:null" json:"createtime"`
	UpdateTime time.Time `gorm:"column:updatetime;default:null" json:"updatetime"`
}
