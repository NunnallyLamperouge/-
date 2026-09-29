package service

import (
	"example/config"
	"example/entity"
	"example/utility"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 申请导师指导请求体
type ApplyGuideRequestBody struct {
	TeacherID int    `json:"teacherId" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

// 更新指导状态请求体
type UpdateGuideStatusRequestBody struct {
	Status int `json:"status" binding:"required"`
}

// ApplyForGuide 学生申请导师指导
func ApplyForGuide(c *gin.Context) {
	var body ApplyGuideRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "提交了错误的JSON"))
		return
	}

	// 获取当前学生用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}
	studentId := utility.GetUserIDFromContext(userId)

	// 检查是否已经申请过该导师
	var existingGuide entity.Guide
	result := config.MysqlDataBase.Where("student_user_id = ? AND teacher_user_id = ?", studentId, body.TeacherID).First(&existingGuide)
	if result.RowsAffected > 0 {
		c.JSON(200, entity.ErrorResponse[string](400, "已经申请过该导师指导，请勿重复申请"))
		return
	}

	// 创建指导申请记录
	guide := entity.Guide{
		StudentUserId: uint(studentId),
		TeacherUserId: uint(body.TeacherID),
		Status:        0,
		Message:       body.Message,
	}

	if err := config.MysqlDataBase.Create(&guide).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "申请导师指导失败: "+err.Error()))
		return
	}

	// 获取学生姓名
	var studentName string
	config.MysqlDataBase.Raw("SELECT username FROM users WHERE id = ?", studentId).Scan(&studentName)

	// 发送通知给导师
	notification := entity.Notification{
		Type:        entity.ThesisProgressNotification,
		Title:       "新的论文指导申请",
		Content:     "学生 " + studentName + " 向您申请论文指导",
		SenderID:    uint(studentId),
		ReceiverID:  uint(body.TeacherID),
		RelatedID:   guide.ID,
		RelatedType: "guide_application",
	}

	if err := sendNotificationInternal(notification); err != nil {
		// 通知发送失败，记录日志但不影响业务流程
		// 可以记录日志
	}

	c.JSON(200, entity.SuccessResponse(guide))
}

// GetStudentGuideList 获取学生的导师指导列表
func GetStudentGuideList(c *gin.Context) {
	// 获取当前学生用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}
	studentId := utility.GetUserIDFromContext(userId)

	var guides []entity.Guide
	if err := config.MysqlDataBase.Where("student_user_id = ?", studentId).Preload("Teacher").Find(&guides).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取指导关系列表失败: "+err.Error()))
		return
	}

	c.JSON(200, entity.SuccessResponse(guides))
}

// GetTeacherGuideList 获取教师的指导学生列表
func GetTeacherGuideList(c *gin.Context) {
	// 获取当前教师用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}
	teacherId := utility.GetUserIDFromContext(userId)

	// 状态参数，可选
	status := c.Query("status")

	var guides []entity.Guide
	query := config.MysqlDataBase.Where("teacher_user_id = ? AND status = 1", teacherId).Preload("Student")

	// 如果提供了状态参数，则按状态筛选
	if status != "" {
		statusInt, err := strconv.Atoi(status)
		if err == nil {
			query = query.Where("status = ?", statusInt)
		}
	}

	if err := query.Find(&guides).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取指导学生列表失败: "+err.Error()))
		return
	}

	c.JSON(200, entity.SuccessResponse(guides))
}

// GetTeacherPendingGuideList 获取教师待确认的指导申请列表
func GetTeacherPendingGuideList(c *gin.Context) {
	// 获取当前教师用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}
	teacherId := utility.GetUserIDFromContext(userId)

	var guides []entity.Guide
	if err := config.MysqlDataBase.Where("teacher_user_id = ? AND status = 0", teacherId).Preload("Student").Find(&guides).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取待确认指导申请列表失败: "+err.Error()))
		return
	}

	c.JSON(200, entity.SuccessResponse(guides))
}

// UpdateGuideStatus 教师更新指导关系状态（确认或拒绝）
func UpdateGuideStatus(c *gin.Context) {
	id := c.Param("id")
	guideId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "无效的指导关系ID"))
		return
	}

	var body UpdateGuideStatusRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "提交了错误的JSON"))
		return
	}

	// 获取当前教师用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}
	teacherId := utility.GetUserIDFromContext(userId)

	// 查找指导关系
	var guide entity.Guide
	if err := config.MysqlDataBase.First(&guide, guideId).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "指导关系不存在"))
		return
	}

	// 检查是否是该指导关系的导师
	if guide.TeacherUserId != uint(teacherId) {
		c.JSON(200, entity.ErrorResponse[string](403, "无权修改此指导关系"))
		return
	}

	// 检查状态值是否有效（0或1）
	if body.Status != 0 && body.Status != 1 {
		c.JSON(200, entity.ErrorResponse[string](400, "无效的状态值，只能为0（拒绝）或1（确认）"))
		return
	}

	// 更新指导关系状态
	if err := config.MysqlDataBase.Model(&guide).Update("status", body.Status).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "更新指导关系状态失败: "+err.Error()))
		return
	}

	// 获取更新后的指导关系
	if err := config.MysqlDataBase.First(&guide, guideId).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取更新后的指导关系失败"))
		return
	}

	// 获取教师姓名
	var teacherName string
	config.MysqlDataBase.Raw("SELECT username FROM users WHERE id = ?", teacherId).Scan(&teacherName)

	// 发送通知给学生
	notificationType := entity.ReviewFeedbackNotification
	notificationTitle := "论文指导申请已确认"
	notificationContent := "导师 " + teacherName + " 已同意您的论文指导申请"

	if body.Status == 0 {
		notificationType = entity.SystemNotification
		notificationTitle = "论文指导申请被拒绝"
		notificationContent = "导师 " + teacherName + " 拒绝了您的论文指导申请"
	}

	notification := entity.Notification{
		Type:        notificationType,
		Title:       notificationTitle,
		Content:     notificationContent,
		SenderID:    uint(teacherId),
		ReceiverID:  uint(guide.StudentUserId),
		RelatedID:   uint(guideId),
		RelatedType: "guide_status_update",
	}

	if err := sendNotificationInternal(notification); err != nil {
		// 通知发送失败，记录日志但不影响业务流程
	}

	c.JSON(200, entity.SuccessResponse(guide))
}

// DeleteGuideRelation 删除指导关系
func DeleteGuideRelation(c *gin.Context) {
	id := c.Param("id")
	guideId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "无效的指导关系ID"))
		return
	}

	// 获取当前用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}
	currentUserId := utility.GetUserIDFromContext(userId)

	// 查找指导关系
	var guide entity.Guide
	if err := config.MysqlDataBase.First(&guide, guideId).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "指导关系不存在"))
		return
	}

	// 检查是否有权限删除（学生或老师都可以删除）
	if guide.StudentUserId != uint(currentUserId) && guide.TeacherUserId != uint(currentUserId) {
		c.JSON(200, entity.ErrorResponse[string](403, "无权删除此指导关系"))
		return
	}

	// 获取当前用户姓名
	var userName string
	config.MysqlDataBase.Raw("SELECT username FROM users WHERE id = ?", currentUserId).Scan(&userName)

	// 确定接收者ID
	var receiverId int
	var userRole string
	if uint(currentUserId) == guide.StudentUserId {
		receiverId = int(guide.TeacherUserId)
		userRole = "学生"
	} else {
		receiverId = int(guide.StudentUserId)
		userRole = "导师"
	}

	// 删除指导关系
	if err := config.MysqlDataBase.Delete(&guide).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "删除指导关系失败: "+err.Error()))
		return
	}

	// 发送通知给对方
	notification := entity.Notification{
		Type:        entity.SystemNotification,
		Title:       "论文指导关系已终止",
		Content:     userRole + " " + userName + " 已终止与您的论文指导关系",
		SenderID:    uint(currentUserId),
		ReceiverID:  uint(receiverId),
		RelatedID:   uint(guideId),
		RelatedType: "guide_relation_deleted",
	}

	if err := sendNotificationInternal(notification); err != nil {
		// 通知发送失败，记录日志但不影响业务流程
	}

	c.JSON(200, entity.SuccessResponse("指导关系删除成功"))
}

type SendTemplateOutRequestBody struct {
	Code       string `json:"code" binding:"required"`
	StudentIDs []int  `json:"studentIds" binding:"required"`
}

func SendTemplateOut(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}

	// 获取当前教师ID
	teacherId := utility.GetUserIDFromContext(userId)

	// 解析请求体
	var body SendTemplateOutRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "提交了错误的JSON"))
		return
	}

	// 验证模板码是否存在
	var template entity.FakeTemplate
	if err := config.MysqlDataBase.Where("code = ?", body.Code).First(&template).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "模板不存在"))
		return
	}

	// 获取教师姓名
	var teacherName string
	config.MysqlDataBase.Raw("SELECT username FROM users WHERE id = ?", teacherId).Scan(&teacherName)

	// 给每个学生发送通知
	successCount := 0
	for _, studentId := range body.StudentIDs {
		// 创建站内通知
		notification := entity.Notification{
			Type:        entity.FormatTemplateUpdate,
			Title:       "论文模板已发放",
			Content:     "导师 " + teacherName + " 已向您发放论文模板，模板码为：" + body.Code + ",请及时查看",
			SenderID:    uint(teacherId),
			ReceiverID:  uint(studentId),
			RelatedID:   template.ID,
			RelatedType: "template_distribution",
		}

		// 保存通知
		if err := sendNotificationInternal(notification); err != nil {
			// 通知发送失败，继续处理下一个学生
			continue
		}

		// 获取学生邮箱
		var studentEmail string
		err := config.MysqlDataBase.Raw("SELECT email FROM users WHERE id = ?", studentId).Scan(&studentEmail).Error
		if err != nil || studentEmail == "" {
			// 学生邮箱不存在，跳过邮件发送
			successCount++
			continue
		}

		// 发送邮件
		emailContent := "亲爱的同学：<br><br>您的导师 " + teacherName + " 已为您发放论文模板，模板码为：" + body.Code + "。<br>请登录系统查看详情并及时完成论文撰写任务。"

		err = utility.SendMailFormatted(studentEmail, "论文模板已发放", emailContent)
		if err != nil {
			log.Println(err.Error())
			// 邮件发送失败，记录日志但不影响业务流程
			// 可以记录日志
		}

		successCount++
	}

	c.JSON(200, entity.SuccessResponse(map[string]interface{}{
		"message": "成功发放模板",
		"success": successCount,
		"total":   len(body.StudentIDs),
	}))
}
