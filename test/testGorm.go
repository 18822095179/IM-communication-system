package main

import (
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:12345@tcp(localhost:3307)/ginchat1?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 表示没有这个表，将会自己创建
	//db.AutoMigrate(&models.UserBasic{})
	//db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.GroupBasic{})
	db.AutoMigrate(&models.Contact{})

	//// 创建
	//user := &models.UserBasic{}
	//user.Name = "申专"
	//db.Create(user)
	//
	//// 读取
	//fmt.Println(db.First(user, 1)) // 根据整型主键查找
	//
	//// 修改
	//db.Model(user).Update("PassWord", "1234")
	//
	//// 删除

}
