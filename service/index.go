package service

import (
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html") // 使用template.ParseFiles解析两个html文件，解析后的模板对象储存在ind变量中
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index") // 使用ind.Execute方法来渲染模板，并将渲染结果写入c.Writer,"index"是传递给模板的数据
	//c.JSON(200, gin.H{
	//	"message": "welcome!!!",
	//})
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/main.html")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId")) // 从HTTP上下文中获取userId
	token := c.Query("token")                    // 从HTTP上下文中获取token
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	//fmt.Println("ToChat>>>>>>", user)
	ind.Execute(c.Writer, user) // 调用ind.Execute用于执行模板，并将结果昔日c.Writer，即HTTP响应中。
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
