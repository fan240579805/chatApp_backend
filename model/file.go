package model

import (
	"chatApp_backend/dao"
	"chatApp_backend/utils"
	"time"
)

/**
处理文件  图片 语音model
*/
type File struct {
	ID        int       `gorm:"column:id;unique;not null;primary_key;AUTO_INCREMENT"`
	FileId    string    `gorm:"column:fileid;unique;not null;"`
	Owner     string    `gorm:"column:owner;not null"`    // 文件的拥有者用户的唯一id
	Type      string    `gorm:"column:typename;not null"` // 文件类型 audio语音 image图片
	FileInfo  string    `gorm:"column:fileinfo;not null"` // 文件命，文件概况
	FileUrl   string    `gorm:"column:url;not null"`  // 文件命，文件概况
	CreatedAt time.Time `gorm:"column:createdat;default:null" json:"createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat;default:null" json:"updatedat"`
}

func InsertFile(fileUrl string, userid string, typeName string, fileInfo string) error {
	file := &File{
		FileId:   "fileID_" + utils.UniqueId(),
		Owner:    userid,
		Type:     typeName,
		FileInfo: fileInfo,
		FileUrl:  fileUrl,
	}
	dbErr := dao.DB.Debug().Create(&file).Error
	if dbErr != nil {
		return dbErr
	}
	return nil
}
