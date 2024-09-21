package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

// InitConfig 初始化配置文件
func InitConfig() {
	viper.SetConfigName("app")    // 设置配置文件的名称叫app，通常为app.yaml或app.json
	viper.AddConfigPath("config") // 告诉viper在config目录下查找app.yaml或app.json配置文件
	err := viper.ReadInConfig()   // 尝试读取配置文件
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited。。。。")

}

// InitRedis 初始化Redis配置
func InitRedis() {
	Red = redis.NewClient(&redis.Options{ // 创建一个新的Redis客户端实例，并赋值给Red
		Addr:         viper.GetString("redis.addr"),      // 配置了redis客户端连接地址
		Password:     viper.GetString("redis.password"),  // 配置了redis客户端连接密码
		DB:           viper.GetInt("redis.DB"),           // 配置了redis客户端使用的数据库编号
		PoolSize:     viper.GetInt("redis.poolSize"),     // 配置了redis客户端连接池大小
		MinIdleConns: viper.GetInt("redis.minIdleConns"), // 配置了redis客户端最小空闲连接数
	})
}

// InitMySQL 初始化MySQL数据库配置（使用gorm连接数据库）和日志记录器
func InitMySQL() {
	// 自定义日志模板 打印SQL语句
	newLogger := logger.New( // 创建一个日志记录器实例
		log.New(os.Stdout, "\r\n", log.LstdFlags), //使用标准输出（os.Stdout）作为日志的输出目标，并设置日志消息之间的分隔符为\r\n
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值，这里是1秒
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		})
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger}) // 使用viper读取配置文件连接数据库
	fmt.Println("MySQL inited。。。。")
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println(user)
}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到Redis。(发布消息到Redis指定频道)
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publish 。。。。", msg)           // 打印一下要发布的消息看一眼
	err = Red.Publish(ctx, channel, ctx).Err() // 使用Redis实例Red在指定频道上发布消息，ctx为上下文对象，用于控制操作的超时和取消
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅Redis消息。(订阅Redis指定频道并接收消息)
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)  // 创建一个接收指定频道消息的对象sub
	fmt.Println("Subscribe 。。。。", ctx)  // 打印订阅的消息
	msg, err := sub.ReceiveMessage(ctx) // 用sub订阅这个频道，等待接收消息
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe 。。。。", msg.Payload) // 打印消息msg的Payload字段（即实际的数据内容）
	return msg.Payload, err                    // 返回接收到的消息和可能发生的错误
}
