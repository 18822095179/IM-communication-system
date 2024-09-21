// 人员关系

package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint // 谁的关系
	TargetId uint // 对应的谁
	Type     int  // 关系类型 1好友 2群 3...
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

// SearchFriend 查找好友
func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id=? and type=1", userId).Find(&contacts) // 从数据库的Contact表中查找userId的type=1(即好友类型)的所有记录，这些记录都是该userId的好友，并存储在contacts切片中。
	for _, v := range contacts {                                    // 从contacts切片中提取所有好友的ID存储到objIds中
		fmt.Println(">>>>>>>>>>>>>>", v)
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users) // 根据这些好友的userId查找到这些还有，存储在users切片中。
	return users
}
