// 群信息

package models

import "gorm.io/gorm"

type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Type    int
	Desc    string
}

// TableName 当你对GroupBasic结构体操作的时候，如创建、查询等，Gorm会使用TableName方法返回字符串作为数据库表的名称，这样就可以在代码中使用
// 结构体名称GroupBasic，而数据库中对应的表却是group_basic
func (table *GroupBasic) TableName() string {
	return "group_basic"
}
