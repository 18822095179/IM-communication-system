package router

import (
	"ginchat/docs"
	"ginchat/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {

	r := gin.Default()

	//swagger
	docs.SwaggerInfo.BasePath = ""                                       // 设置swagger的基本路径为空字符串
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // 设置swagger的url为根路径"/swagger/*any",这样用户可以通过"/swagger"访问Swagger UI

	//静态资源
	r.Static("/asset", "asset/") // 设置静态资源目录为"/asset",路径为"asset/"
	r.LoadHTMLGlob("views/**/*") //设置HTML模板目录为"views/**/*"，路径为"views/" ("**表示匹配任意数量的子目录，而*表示匹配当前目录下的任意文件名)

	//首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.GET("/chat", service.Chat)
	r.POST("/searchFriends", service.SearchFriends)

	//用户模块
	r.POST("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("user/findUserByNameAndPwd", service.FindUserByNameAndPwd)

	// 发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	return r

}
