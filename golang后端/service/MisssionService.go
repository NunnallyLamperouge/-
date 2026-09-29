package service

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"example/config"
	"example/entity"
	"example/utility"
)

// MissionCreateRequestBody 创建任务请求体
type MissionCreateRequestBody struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	StudentIDs  []uint `json:"studentIds" binding:"required"`
	EndTime     string `json:"endTime" binding:"required"`
	StartTime   string `json:"startTime" binding:"required"`
	Priority    int    `json:"priority"`
	TemplateId  uint   `json:"templateId" binding:"required"` // 关联的模板ID
}

// MissionSubmissionRequestBody 学生提交任务成果请求体
type MissionSubmissionRequestBody struct {
	Content    string `json:"content" binding:"required"` // 提交的内容
	TemplateID int    `json:"templateID"`
}

// CreateMission 创建任务
func CreateMission(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}

	// 获取当前教师ID
	teacherId := utility.GetUserIDFromContext(userId)

	// 解析请求体
	var body MissionCreateRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "提交了错误的JSON"))
		return
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", body.StartTime)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "起始时间格式错误"))
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", body.EndTime)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "截止时间格式错误"))
		return
	}

	// 检查时间有效性
	if endTime.Before(startTime) {
		c.JSON(200, entity.ErrorResponse[string](400, "截止时间不能早于起始时间"))
		return
	}

	// 检查模板是否存在
	var template entity.FakeTemplate
	if err := config.MysqlDataBase.First(&template, body.TemplateId).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "模板不存在"))
		return
	}

	// 获取教师姓名
	var teacherName string
	config.MysqlDataBase.Raw("SELECT username FROM users WHERE id = ?", teacherId).Scan(&teacherName)

	// 开启事务
	tx := config.MysqlDataBase.Begin()

	// 给每个学生发布任务
	successCount := 0
	for _, studentId := range body.StudentIDs {
		// 创建任务
		mission := entity.Mission{
			Name:          body.Name,
			Description:   body.Description,
			StudentUserId: studentId,
			TeacherUserId: uint(teacherId),
			Status:        entity.MissionStatusInProgress,
			StartTime:     startTime,
			EndTime:       endTime,
			Priority:      body.Priority,
			TemplateId:    body.TemplateId,
		}

		// 保存任务
		if err := tx.Create(&mission).Error; err != nil {
			tx.Rollback()
			c.JSON(200, entity.ErrorResponse[string](500, "创建任务失败"))
			return
		}

		// 创建站内通知
		notification := entity.Notification{
			Type:        entity.ThesisProgressNotification,
			Title:       "新任务已分配",
			Content:     "导师 " + teacherName + " 已为您分配新任务：" + body.Name + "，请及时查看并完成。",
			SenderID:    uint(teacherId),
			ReceiverID:  studentId,
			RelatedID:   mission.ID,
			RelatedType: "mission_assignment",
		}

		if err := sendNotificationInternal(notification); err != nil {
			log.Println("发送通知失败:", err)
			continue
		}

		// 获取学生邮箱
		var studentEmail string
		err := config.MysqlDataBase.Raw("SELECT email FROM users WHERE id = ?", studentId).Scan(&studentEmail).Error
		if err != nil || studentEmail == "" {
			successCount++
			continue
		}

		// 发送邮件
		emailContent := "亲爱的同学：<br><br>您的导师 " + teacherName + " 已为您分配新任务：'" + body.Name + "'。<br><br>" +
			"任务详情：<br>" +
			"- 描述：" + body.Description + "<br>" +
			"- 开始时间：" + body.StartTime + "<br>" +
			"- 截止时间：" + body.EndTime + "<br><br>" +
			"请登录系统查看更多详情并及时完成任务。"

		err = utility.SendMailFormatted(studentEmail, "新任务已分配："+body.Name, emailContent)
		if err != nil {
			log.Println("发送邮件失败:", err.Error())
		}

		successCount++
	}

	// 提交事务
	tx.Commit()

	c.JSON(200, entity.SuccessResponse(map[string]interface{}{
		"message": "成功发布任务",
		"success": successCount,
		"total":   len(body.StudentIDs),
	}))
}

// GetTeacherMissions 获取教师发布的所有任务
func GetTeacherMissions(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}

	// 获取当前教师ID
	teacherId := utility.GetUserIDFromContext(userId)

	var missions []entity.Mission
	if err := config.MysqlDataBase.Preload("Student").Where("teacher_user_id = ?", teacherId).Find(&missions).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取任务失败"))
		return
	}

	c.JSON(200, entity.SuccessResponse(missions))
}

// GetStudentMissions 获取学生的所有任务
func GetStudentMissions(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}

	// 获取当前学生ID
	studentId := utility.GetUserIDFromContext(userId)

	var missions []entity.Mission
	if err := config.MysqlDataBase.Preload("Teacher").Preload("Template").Where("student_user_id = ?", studentId).Find(&missions).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取任务失败"))
		return
	}

	c.JSON(200, entity.SuccessResponse(missions))
}

// UpdateMissionStatus 更新任务状态
func UpdateMissionStatus(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}

	// 获取参数
	missionId := c.Param("id")
	var body struct {
		Status          entity.MissionStatus `json:"status"`
		Feedback        string               `json:"feedback"`
		RejectionReason string               `json:"rejectionReason"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "提交了错误的JSON"))
		return
	}

	// 检查任务是否存在
	var mission entity.Mission
	if err := config.MysqlDataBase.First(&mission, missionId).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "任务不存在"))
		return
	}

	currentUserId := utility.GetUserIDFromContext(userId)
	updates := make(map[string]interface{})
	now := time.Now()

	// 检查用户权限
	if mission.TeacherUserId == uint(currentUserId) {
		// 教师可以更新状态和反馈
		if body.Status < entity.MissionStatusDraft || body.Status > entity.MissionStatusCompleted {
			c.JSON(200, entity.ErrorResponse[string](400, "无效的状态值"))
			return
		}

		// 如果当前状态是待审核，教师只能将其改为已完成或驳回
		if mission.Status == entity.MissionStatusSubmitted &&
			body.Status != entity.MissionStatusRejected &&
			body.Status != entity.MissionStatusCompleted {
			c.JSON(200, entity.ErrorResponse[string](400, "当前只能将任务标记为已完成或驳回"))
			return
		}

		updates["status"] = body.Status
		updates["feedback"] = body.Feedback
		updates["review_time"] = now

		if body.Status == entity.MissionStatusRejected {
			updates["rejection_reason"] = body.RejectionReason
		}
	} else if mission.StudentUserId == uint(currentUserId) {
		// 学生只能将任务状态改为已提交
		if body.Status != entity.MissionStatusSubmitted {
			c.JSON(200, entity.ErrorResponse[string](400, "学生只能提交任务"))
			return
		}

		if mission.Status == entity.MissionStatusSubmitted ||
			mission.Status == entity.MissionStatusCompleted {
			c.JSON(200, entity.ErrorResponse[string](400, "当前状态下不能提交任务"))
			return
		}

		updates["status"] = entity.MissionStatusSubmitted
		updates["submission_time"] = now
	} else {
		c.JSON(200, entity.ErrorResponse[string](403, "无权更新此任务"))
		return
	}

	// 更新任务
	if err := config.MysqlDataBase.Model(&mission).Updates(updates).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "更新任务失败"))
		return
	}

	// 如果是教师提供反馈，发送通知给学生
	if mission.TeacherUserId == uint(currentUserId) && (body.Feedback != "" || body.Status != mission.Status) {
		// 获取教师姓名
		var teacherName string
		config.MysqlDataBase.Raw("SELECT username FROM users WHERE id = ?", currentUserId).Scan(&teacherName)

		// 构建通知标题和内容
		var title, content string
		switch body.Status {
		case entity.MissionStatusRejected:
			title = "任务被驳回"
			content = "导师 " + teacherName + " 已驳回您的任务：" + mission.Name
			if body.RejectionReason != "" {
				content += "，原因：" + body.RejectionReason
			}
		case entity.MissionStatusCompleted:
			title = "任务已完成"
			content = "导师 " + teacherName + " 已通过您的任务：" + mission.Name
		default:
			title = "任务状态更新"
			content = "导师 " + teacherName + " 更新了任务：" + mission.Name + " 的状态"
		}

		if body.Feedback != "" {
			content += "，并提供了新的反馈"
		}

		// 创建站内通知
		notification := entity.Notification{
			Type:        entity.ReviewFeedbackNotification,
			Title:       title,
			Content:     content,
			SenderID:    uint(currentUserId),
			ReceiverID:  mission.StudentUserId,
			RelatedID:   mission.ID,
			RelatedType: "mission_feedback",
		}

		sendNotificationInternal(notification)

		// 发送邮件
		var studentEmail string
		if err := config.MysqlDataBase.Raw("SELECT email FROM users WHERE id = ?", mission.StudentUserId).Scan(&studentEmail).Error; err == nil && studentEmail != "" {
			var statusText string
			switch body.Status {
			case entity.MissionStatusRejected:
				statusText = "已驳回"
			case entity.MissionStatusCompleted:
				statusText = "已完成"
			default:
				statusText = "状态已更新"
			}

			emailContent := "亲爱的同学：<br><br>您的任务 '" + mission.Name + "' " + statusText + "。<br><br>"
			if body.RejectionReason != "" {
				emailContent += "驳回原因：<br>" + body.RejectionReason + "<br><br>"
			}
			if body.Feedback != "" {
				emailContent += "导师反馈：<br>" + body.Feedback + "<br><br>"
			}
			emailContent += "请登录系统查看详情。"

			utility.SendMailFormatted(studentEmail, title+"："+mission.Name, emailContent)
		}
	}

	c.JSON(200, entity.SuccessResponse(map[string]string{
		"message": "更新任务成功",
	}))
}

// DeleteMission 删除任务
func DeleteMission(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}

	// 获取任务ID
	missionId := c.Param("id")

	// 检查任务是否存在
	var mission entity.Mission
	if err := config.MysqlDataBase.First(&mission, missionId).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "任务不存在"))
		return
	}

	// 检查权限（只有教师可以删除）
	currentUserId := utility.GetUserIDFromContext(userId)
	if mission.TeacherUserId != uint(currentUserId) {
		c.JSON(200, entity.ErrorResponse[string](403, "无权删除此任务"))
		return
	}

	// 删除任务
	if err := config.MysqlDataBase.Delete(&mission).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "删除任务失败"))
		return
	}

	c.JSON(200, entity.SuccessResponse(map[string]string{
		"message": "删除任务成功",
	}))
}

// SubmitMission 学生提交任务成果
func SubmitMission(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}

	studentId := utility.GetUserIDFromContext(userId)
	missionId := c.Param("id")

	// 检查任务是否存在且属于该学生
	var mission entity.Mission
	if err := config.MysqlDataBase.First(&mission, missionId).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "任务不存在"))
		return
	}

	if mission.StudentUserId != uint(studentId) {
		c.JSON(200, entity.ErrorResponse[string](403, "无权提交此任务"))
		return
	}

	// 解析请求体
	var body MissionSubmissionRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "提交了错误的JSON"))
		return
	}
	log.Println(body.TemplateID)
	now := time.Now()
	// 更新任务
	updates := map[string]interface{}{
		"content":         body.Content,
		"status":          entity.MissionStatusSubmitted,
		"template_id":     body.TemplateID,
		"submission_time": now,
	}

	if err := config.MysqlDataBase.Model(&mission).Updates(updates).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "提交任务失败"))
		return
	}

	// 获取学生姓名
	var studentName string
	config.MysqlDataBase.Raw("SELECT username FROM users WHERE id = ?", studentId).Scan(&studentName)

	// 创建站内通知
	notification := entity.Notification{
		Type:        entity.ThesisProgressNotification,
		Title:       "任务提交待审核",
		Content:     "学生 " + studentName + " 已提交任务：" + mission.Name + "，请及时审核。",
		SenderID:    uint(studentId),
		ReceiverID:  mission.TeacherUserId,
		RelatedID:   mission.ID,
		RelatedType: "mission_submission",
	}

	if err := sendNotificationInternal(notification); err != nil {
		log.Println("发送通知失败:", err)
	}

	// 获取教师邮箱并发送邮件
	var teacherEmail string
	if err := config.MysqlDataBase.Raw("SELECT email FROM users WHERE id = ?", mission.TeacherUserId).Scan(&teacherEmail).Error; err == nil && teacherEmail != "" {
		emailContent := "尊敬的导师：<br><br>您的学生 " + studentName + " 已提交任务：'" + mission.Name + "'。<br><br>" +
			"请登录系统查看详情并进行审核。"

		err = utility.SendMailFormatted(teacherEmail, "任务提交待审核："+mission.Name, emailContent)
		if err != nil {
			log.Println("发送邮件失败:", err.Error())
		}
	}

	c.JSON(200, entity.SuccessResponse(map[string]string{
		"message": "提交任务成功",
	}))
}

// 内部函数：发送通知
func sendNotificationInternal(notification entity.Notification) error {
	return config.MysqlDataBase.Create(&notification).Error
}

// CheckDeadlines 检查任务截止时间并发送提醒
// 建议通过定时任务调用此方法
func CheckDeadlines() {
	var missions []entity.Mission
	now := time.Now()

	// 查找未完成且未发送提醒的任务
	err := config.MysqlDataBase.Where(
		"status != ? AND reminder_sent = ? AND end_time > ? AND end_time <= ?",
		entity.MissionStatusCompleted,
		false,
		now,
		now.Add(24*time.Hour),
	).Find(&missions).Error

	if err != nil {
		log.Println("查询待提醒任务失败:", err)
		return
	}

	for _, mission := range missions {
		// 发送站内通知
		notification := entity.Notification{
			Type:        entity.ThesisProgressNotification,
			Title:       "任务即将截止提醒",
			Content:     fmt.Sprintf("您的任务 '%s' 将在24小时内截止，请及时完成。", mission.Name),
			ReceiverID:  mission.StudentUserId,
			RelatedID:   mission.ID,
			RelatedType: "mission",
		}

		if err := sendNotificationInternal(notification); err != nil {
			log.Printf("发送任务截止提醒通知失败 (任务ID: %d): %v\n", mission.ID, err)
			continue
		}

		// 发送邮件提醒
		var studentEmail string
		if err := config.MysqlDataBase.Raw("SELECT email FROM users WHERE id = ?", mission.StudentUserId).Scan(&studentEmail).Error; err == nil && studentEmail != "" {
			emailContent := fmt.Sprintf(
				"亲爱的同学：<br><br>您的任务 '%s' 将在24小时内截止。<br><br>请及时登录系统完成任务。",
				mission.Name,
			)

			err = utility.SendEmail(studentEmail, "任务截止提醒", emailContent)
			if err != nil {
				log.Printf("发送任务截止提醒邮件失败 (任务ID: %d): %v\n", mission.ID, err)
			}
		}

		// 更新提醒状态
		mission.ReminderSent = true
		mission.LastReminderTime = &now
		if err := config.MysqlDataBase.Save(&mission).Error; err != nil {
			log.Printf("更新任务提醒状态失败 (任务ID: %d): %v\n", mission.ID, err)
		}
	}
}

// GetMissionHistory 获取任务修改历史
func GetMissionHistory(c *gin.Context) {
	missionId := c.Param("id")

	var mission entity.Mission
	if err := config.MysqlDataBase.First(&mission, missionId).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "任务不存在"))
		return
	}

	// 获取模板历史记录
	var templateHistories []entity.TemplateHistory
	if err := config.MysqlDataBase.Where("template_id = ?", mission.TemplateId).Order("created_at DESC").Find(&templateHistories).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取历史记录失败"))
		return
	}

	c.JSON(200, entity.SuccessResponse(gin.H{
		"mission": mission,
		"history": templateHistories,
	}))
}

// GetMissionStatistics 获取任务统计信息
func GetMissionStatistics(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}

	var stats struct {
		Total      int64 `json:"total"`
		InProgress int64 `json:"inProgress"`
		Submitted  int64 `json:"submitted"`
		Completed  int64 `json:"completed"`
		Rejected   int64 `json:"rejected"`
	}

	db := config.MysqlDataBase.Model(&entity.Mission{})

	// 根据用户ID查询是否为教师
	var teacherCount int64
	config.MysqlDataBase.Model(&entity.Guide{}).Where("teacher_user_id = ?", userId).Count(&teacherCount)

	if teacherCount > 0 {
		// 是教师
		db = db.Where("teacher_user_id = ?", userId)
	} else {
		// 是学生
		db = db.Where("student_user_id = ?", userId)
	}

	// 统计总数
	db.Count(&stats.Total)

	// 统计各状态数量
	db.Where("status = ?", entity.MissionStatusInProgress).Count(&stats.InProgress)
	db.Where("status = ?", entity.MissionStatusSubmitted).Count(&stats.Submitted)
	db.Where("status = ?", entity.MissionStatusCompleted).Count(&stats.Completed)
	db.Where("status = ?", entity.MissionStatusRejected).Count(&stats.Rejected)

	c.JSON(200, entity.SuccessResponse(stats))
}
