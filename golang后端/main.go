package main

import (
	"example/config"
	"example/entity"
	"example/route"
	"example/service"
	"example/utility"
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ReadConfig()
	service.InitSpeechServices()
	// 初始化mysql
	config.InitMysqlDataBase()
	// 初始化redis
	config.InitRedis()

	// 禁用外键约束，避免迁移顺序问题
	config.MysqlDataBase.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// 自动迁移数据库
	config.MysqlDataBase.AutoMigrate(
		&entity.User{},
		&service.ImageUpload{},
		&entity.FakeTemplate{},
		&entity.TemplateHistory{},
		&entity.Guide{},
		&entity.Notification{},
		&entity.ChatGroup{},
		&entity.Message{},
		&entity.GroupMember{},
		&entity.NotificationSetting{},
		&entity.Mission{},
		&entity.Virtual{},
		&entity.VirtualQuestion{},
	)

	// 恢复外键约束
	//config.MysqlDataBase.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// 创建必要的目录
	createRequiredDirectories()
	utility.InitAnythingLLM(config.Config.AnythingLLM.BaseURL, config.Config.AnythingLLM.APIKey)
	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.Config.Gin.CorsAllowOrigins, // 允许的前端来源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "x-requested-with"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           6 * time.Hour, // 预检请求的缓存时间
	}))
	// app.Static("/api/uploads", "./uploads") 静态文件访问路径以及目录
	app.Static("/api/uploads", "./uploads")
	app.Static("/api/audio", "./audio")
	app.Static("/api/user_documents", "./user_documents") // 用户文档静态服务
	route.RegisterRoutes(app)
	app.Run(config.Config.Gin.Port)
}

// createRequiredDirectories 创建应用所需的目录
func createRequiredDirectories() {
	directories := []string{
		"uploads",
		"audio",
		"user_documents",
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("创建目录 %s 失败: %v\n", dir, err)
		}
	}
}
