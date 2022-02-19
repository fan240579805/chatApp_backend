package model

import (
	"chatApp_backend/dao"
	"errors"
	"time"
)

// Relation
//	status：关系状态码
//  -1：from主动加to，to用户还没同意，挂起
//	1: from主动加to，to用户同意了，并成为好友
//	2：from删除to
//	3：to删除from
//	4: from拉黑了to
//	5: to拉黑了from
type Relation struct {
	ID        int       `gorm:"column:id;unique;not null;primary_key;AUTO_INCREMENT"`
	Status    int       `gorm:"column:status;not null"`
	From      string    `gorm:"column:fromid;not null"` // 发起好友关系的人的唯一id
	To        string    `gorm:"column:toid;not null"`   // 接受好友关系的人的唯一id
	CreatedAt time.Time `gorm:"column:createdat;default:null" json:"createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat;default:null" json:"updatedat"`
}

// AddFriendRecord 发起添加好友请求
func AddFriendRecord(username string, fromUserid string, ToUserid string) error {
	var user User
	// 查询一下是否已经添加
	rightRelation, _ := GetRightRelationRecord(fromUserid, ToUserid)
	if rightRelation.From == "" {
		// 通过账号得到to用户的userid
		finderr := dao.DB.Where("username=?", username).Select("userid").Find(&user).Error
		if finderr != nil {
			return finderr
		}
		relation := &Relation{
			Status: -1,
			From:   fromUserid,
			To:     user.UserID,
		}
		err := dao.DB.Debug().Create(&relation).Error
		if err != nil {
			return errors.New("用户入库失败")
		}
		return nil
	} else {
		return errors.New("已经添加过该用户")
	}
}

// ModifyStatus 根据正确的关系record更新status
func ModifyStatus(userid1 string, userid2 string, status int) error {
	rightRelation, getRightErr := GetRightRelationRecord(userid1, userid2)
	if getRightErr != nil {
		return getRightErr
	}
	if rightRelation.From != "" {
		updateErr := dao.DB.Model(&rightRelation).Debug().Update("status", status).Error
		if updateErr != nil {
			return updateErr
		}
	} else {
		return getRightErr
	}
	return nil
}

// GetRightRelationRecord 在不知道到底是谁发起的好友关系时
// 检查这两位用户分别在数据库中的记录，谁是发起好友请求的人，谁是接受的人
// 返回一个正确的关系记录以便后续操作
func GetRightRelationRecord(fromUserid string, toUserid string) (Relation, error) {
	var relation Relation
	finderr1 := dao.DB.Debug().Where(&Relation{From: fromUserid, To: toUserid}).First(&relation).Error
	if finderr1 != nil {
		// 说明没查到用例1，继续查
		finderr2 := dao.DB.Debug().Where(&Relation{From: toUserid, To: fromUserid}).First(&relation).Error
		if finderr2 != nil {
			return relation, finderr2
		} else {
			return relation, nil
		}
		return Relation{}, finderr1
	}
	return relation, nil
}

// SelectFriends 选出userid用户的好友数组
func SelectFriends(userid string) ([]Relation, error) {
	// 发起该段好友relation的数组，from == userid
	var mineRealtions []Relation
	// 接受这段好友relation的数组  to == userid
	var otherRealtions []Relation
	// 2：此时请求friendList的用户在relation表中是from, 不需要from主动删除即 status = 2 的好友
	selectMineErr := dao.DB.Debug().Where("fromid=? AND status!=?", userid, 2).Find(&mineRealtions).Error
	if selectMineErr != nil {
		return mineRealtions, selectMineErr
	}
	// 3：此时请求friendList的用户在relation表中是to, 不需要to主动删除即 status = 3 的好友
	selectOtherErr := dao.DB.Debug().Where("toid=? AND status!=?", userid, 3).Find(&otherRealtions).Error
	if selectOtherErr != nil {
		return otherRealtions, selectOtherErr
	}
	if len(mineRealtions) > 0 || len(otherRealtions) > 0 {
		return append(mineRealtions, otherRealtions...), nil
	} else {
		return []Relation{}, nil
	}
}

func DeleteRelation(relation Relation) error {
	err:= dao.DB.Debug().Where("fromid=? AND toid=?", relation.From,relation.To).Delete(&Relation{}).Error
	if err != nil {
		return err
	}
	return nil
}