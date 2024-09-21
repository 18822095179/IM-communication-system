package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) { // 获取所有用户列表
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0,
		"message": "用户名已注册！",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	//user.Name = c.Query("name") // 从HTTP请求的查询参数中获取名为name的值。
	//password := c.Query("password")
	//repassword := c.Query("repassword")
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("Identity")
	fmt.Println(user.Name, ">>>>>>>>>>", password, repassword)
	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FineUserByName(user.Name)
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(200, gin.H{
			"code":    -1, // 0成功 -1失败
			"message": "用户名或密码不能为空！",
			"data":    user,
		})
		return
	}
	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "用户名已注册！",
			"data":    user,
		})
		return
	}

	if password != repassword {
		c.JSON(200, gin.H{
			"code":    -1, // 0成功 -1失败
			"message": "两次密码不一致",
			"data":    user,
		})
		return
	}
	//user.PassWord = password
	user.PassWord = utils.MakePassword(password, salt) // md5加密的密码
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0, // 0成功 -1失败
		"message": "新增用户成功！",
		"data":    user,
	})
}

// FindUserByNameAndPwd
// @Summary 登录
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	//name := c.Query("name")
	//password := c.Query("password")
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	user := models.FineUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, // 0成功 -1失败
			"message": "该用户不存在！",
			"data":    data,
		})
		return
	}
	fmt.Println(user)
	flag := utils.ValidPassword(password, user.Salt, user.PassWord) // 验证密码，即将password加密，然后对比与数据库中加密后的密码是否一样
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, // 0成功 -1失败
			"message": "密码不正确！",
			"data":    data,
		})
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(200, gin.H{
		"code":    0, // 0成功 -1失败
		"message": "登录成功",
		"data":    data,
	})

}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.Name = c.Query("name")
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0, // 0成功 -1失败
		"message": "删除用户成功！",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id")) // 从HTTP的POST表单中获取名为id的值
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	fmt.Println("update:", user)

	_, err := govalidator.ValidateStruct(user) // 使用go中的govalidator库验证结构体中的字段，比如email格式是否正确
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1, // 0成功 -1失败
			"message": "修改参数不匹配",
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0, // 0成功 -1失败
			"message": "修改用户成功！",
			"data":    user,
		})
	}
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// SendMsg 处理WebSocket连接的请求并发送信息
func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil) // 使用Upgrade方法将HTTP请求升级为WebSocket连接。ws为WebSocket的连接对象
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) { // 关闭WebSocket连接
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c) // 发送消息
}

// MsgHandler 从一个消息发布服务订阅消息，并通过WebSocket连接发送给客户端
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for { // 使用for循环，持续接收信息并通过WebSocket进行发送
		msg, err := utils.Subscribe(c, utils.PublishKey) // 调用utils包下的Subscribe函数订阅信息
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("发送消息：", msg) // 打印一下要发送的信息
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg) // 将消息和时间组合在一起，形成一个完整的消息字符串
		err = ws.WriteMessage(1, []byte(m))      // 通过WebSocket连接发送消息字符串给客户端，这里发送的是文本消息。
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

// SearchFriends 查找好友
func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriend(uint(id))

	//c.JSON(200, gin.H{
	//	"code":    0,
	//	"message": "查询好友列表成功！",
	//	"data":    users,
	//})
	utils.RespOKList(c.Writer, users, len(users))
}
