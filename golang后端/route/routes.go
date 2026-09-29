package route

import (
	"example/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("login", service.Login)
		authGroup.POST("sendCode", service.SendVerifyCode)
		authGroup.POST("register", service.Register)
	}
	userGroup := r.Group("/api/user", preHandler())
	{
		userGroup.GET("me", service.GetMyInfo)
		userGroup.GET("all", service.GetAllUsers)
		userGroup.GET("all-teacher", service.GetAllTeachers)
		userGroup.GET("all-student", service.GetAllStudents)
		userGroup.POST("uploadImage", service.UploadImage)
	}
	// 新增群组聊天服务路由组
	chatGroup := r.Group("/api/chat", preHandler())
	{
		// 群组管理
		chatGroup.POST("/groups", service.CreateChatGroup)
		chatGroup.GET("/groups", service.GetUserChatGroups)
		chatGroup.GET("/groups/:id", service.GetChatGroupDetail)
		chatGroup.PUT("/groups/:id", service.UpdateChatGroup)
		chatGroup.DELETE("/groups/:id", service.DeleteChatGroup)

		// 成员管理
		chatGroup.POST("/groups/:id/members", service.AddGroupMember)
		chatGroup.DELETE("/groups/:id/members/:userId", service.RemoveGroupMember)
		chatGroup.POST("/groups/:id/leave", service.LeaveGroup)
		chatGroup.GET("/groups/:id/members", service.GetGroupMembers)
		chatGroup.PUT("/groups/:id/members/:userId", service.UpdateGroupMemberInfo)
		chatGroup.POST("/groups/:id/members/:userId/mute", service.MuteGroupMember)
		chatGroup.POST("/groups/:id/members/:userId/admin", service.SetGroupAdmin)

		// 消息管理
		chatGroup.GET("/groups/:id/messages", service.GetGroupMessages)
	}

	// 模板相关路由
	templateGroup := r.Group("/api/template", preHandler())
	{
		templateGroup.POST("/create", service.FakeTemplateCreate)
		templateGroup.GET("/list", service.FakeTemplateGetList)
		templateGroup.GET("/list-student", service.FakeTemplateGetListStudent)
		templateGroup.GET("/detail/:id", service.FakeTemplateGetById)
		templateGroup.POST("/update/:id", service.FakeTemplateUpdate)
		templateGroup.POST("/delete/:id", service.FakeTemplateDelete)
		templateGroup.POST("/send-template-out", service.SendTemplateOut)
		templateGroup.POST("/fork", service.FakeTemplateFork)
		templateGroup.POST("/adverise", service.GetAdverise)
		templateGroup.POST("/adveriseUniversal", service.GetAdveriseUniversal)
		templateGroup.GET("/history/:id/list", service.GetTemplateHistortyIDs)
		templateGroup.GET("/history/:id/content", service.GetTemplateHistoryContent)
	}

	// 论文指导相关路由
	guideGroup := r.Group("/api/guide", preHandler())
	{
		// 学生视角
		guideGroup.POST("/apply", service.ApplyForGuide)
		guideGroup.GET("/student/list", service.GetStudentGuideList)

		// 导师视角
		guideGroup.GET("/teacher/list", service.GetTeacherGuideList)
		guideGroup.GET("/teacher/pending", service.GetTeacherPendingGuideList)
		guideGroup.PUT("/:id/status", service.UpdateGuideStatus)

		// 共用
		guideGroup.DELETE("/:id", service.DeleteGuideRelation)
	}

	// 仪表盘数据相关路由
	dashboardGroup := r.Group("/api/dashboard", preHandler())
	{
		// 学生仪表盘数据
		dashboardGroup.GET("/student", service.GetStudentDashboardData)
	}

	// 论文任务相关路由
	missionGroup := r.Group("/api/mission", preHandler())
	{
		// 任务管理
		missionGroup.POST("/create", service.CreateMission) // 创建任务
		missionGroup.DELETE("/:id", service.DeleteMission)  // 删除任务

		// 任务查询
		missionGroup.GET("/teacher", service.GetTeacherMissions)      // 获取教师发布的任务
		missionGroup.GET("/student", service.GetStudentMissions)      // 获取学生的任务
		missionGroup.GET("/statistics", service.GetMissionStatistics) // 获取任务统计信息
		missionGroup.GET("/:id/history", service.GetMissionHistory)   // 获取任务历史记录

		// 任务状态更新
		missionGroup.PUT("/:id/status", service.UpdateMissionStatus) // 更新任务状态
		missionGroup.POST("/:id/submit", service.SubmitMission)      // 提交任务
	}

	// 新增通知服务路由组
	notificationGroup := r.Group("/api/notifications", preHandler())
	{
		notificationGroup.GET("", service.GetUserNotifications)
		notificationGroup.GET("/:id", service.GetNotificationDetail)
		notificationGroup.POST("/:id/read", service.MarkNotificationAsRead)
		notificationGroup.POST("/read-all", service.MarkAllNotificationsAsRead)
		notificationGroup.DELETE("/:id", service.DeleteNotification)
		notificationGroup.DELETE("", service.DeleteAllNotifications)
		notificationGroup.GET("/unread-count", service.GetUnreadNotificationCount)
		notificationGroup.GET("/settings", service.GetNotificationSettings)
		notificationGroup.POST("/settings", service.UpdateNotificationSettings)
	}

	websocketGroup := r.Group("/api/ws")
	{
		websocketGroup.GET("groupChat/:id", service.HandleGroupChat)
		websocketGroup.GET("projectSuggest", service.StreamChatCompletion)

	}

	virtualGroup := r.Group("/api/virtual", preHandler())
	{
		virtualGroup.POST("/question", service.GetVirtualQuestion)
		virtualGroup.POST("/audio", service.GetVirtualQuestionAudio) // 音频生成方式
		virtualGroup.POST("/submit", service.SubmitVirtualDefense)
		virtualGroup.GET("/get", service.GetVirtualDefense)
		virtualGroup.GET("/list", service.GetVirtualList)
		virtualGroup.POST("/score", service.GetScore)
	}

	// 新增RAG服务相关路由组
	ragGroup := r.Group("/api/rag", preHandler())
	{
		// 工作区管理
		ragGroup.POST("/workspace", service.CreateUserWorkspace)
		ragGroup.GET("/workspace", service.GetUserWorkspace)

		// 文档管理
		ragGroup.POST("/documents", service.UploadRAGDocument)        // 上传文档
		ragGroup.GET("/documents", service.GetUserRAGDocuments)       // 获取用户文档列表
		ragGroup.POST("/documents/add", service.AddRAGDocument)       // 添加文档到工作区
		ragGroup.POST("/documents/delete", service.DeleteRAGDocument) // 从工作区删除文档
		ragGroup.POST("/documents/status", service.SwitchFilePinStatus)
		ragGroup.POST("chatComplention", service.RAGParagraphAdvancedRequest)
	}
}
