package model

import (
	"chatApp/dao"
	"chatApp/utils"
	"time"
)

/**
处理文件  图片 语音model
*/
type File struct {
	FileId     string    `gorm:"column:fileid;not null"`
	Owner      string    `gorm:"column:owner;not null"`    // 文件的拥有者用户的唯一id
	Type       string    `gorm:"column:type;not null"`     // 文件类型 audio语音 image图片
	FileInfo   string    `gorm:"column:fileinfo;not null"` // 文件命，文件概况
	CreateTime time.Time `gorm:"column:createtime;default:null" json:"createtime"`
	UpdateTime time.Time `gorm:"column:updatetime;default:null" json:"updatetime"`
}

func InsertFile(filePath string, userid string, typeName string, fileInfo string) error {
	file := &File{
		FileId:   "fileid_" + utils.UniqueId(),
		Owner:    userid,
		Type:     typeName,
		FileInfo: fileInfo,
	}
	dbErr := dao.DB.Debug().Create(&file).Error
	if dbErr != nil {
		return dbErr
	}
	return nil
}
