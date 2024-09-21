package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"` //在数据库中的别名
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

// GetUserList 用于获取用户列表，它通过调用utils.DB.Find()方法查询数据库中的所有用户，并将结果存储在data切片中
func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data) // 调用DB.Find来查询数据库下的所有用户并储存在创建好的data切片中
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

// FindUserByNameAndPwd 登录
// 根据用户名和密码查询用户，它通过调用utils.DB.Where().First()方法根据给定的用户名和密码查询用户，并将结果存储在user变量中。然后，获取当
// 前时间戳并加密生成身份验证令牌，最后通过调用utils.DB.Model().Where().Update()方法更新用户记录中的身份验证令牌字段的值
func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name=? and pass_word=?", name, password).First(&user) // 根据name和password查询用户并储存到user中

	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())                             // 获取时间戳并转换为字符串格式，存储在str中
	temp := utils.MD5Encode(str)                                            // 使用MD5对str进行加密
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("Identity", temp) // 更新用户身份验证令牌Identity的值为temp
	return user
}

// FineUserByName 通过名字查询用户
func FineUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

// FineUserByPhone 通过电话号查询用户
func FineUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", phone).First(&user)
	return user
}

// FineUserByEmail 通过Email查询用户
func FineUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

// CreateUser 新增用户
func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

// DeleteUser 删除用户
func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

// UpdateUser 更新用户信息
func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}
