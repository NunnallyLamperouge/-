package service

import (
	"example/config"
	"example/entity"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserNotifications 获取用户的通知列表
func GetUserNotifications(c *gin.Context) {
	userId, _ := c.Get("userId")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	notificationType, _ := strconv.Atoi(c.DefaultQuery("type", "0")) // 0表示全部类型
	isRead := c.DefaultQuery("isRead", "")                           // 空字符串表示全部，"true"表示已读，"false"表示未读

	offset := (pageNum - 1) * pageSize
	var notifications []entity.Notification
	var total int64

	db := config.MysqlDataBase.Model(&entity.Notification{}).Where("receiver_id = ?", userId)

	// 根据类型筛选
	if notificationType > 0 {
		db = db.Where("type = ?", notificationType)
	}

	// 根据已读状态筛选
	if isRead == "true" {
		db = db.Where("is_read = ?", true)
	} else if isRead == "false" {
		db = db.Where("is_read = ?", false)
	}

	// 获取总数
	db.Count(&total)

	// 获取分页数据
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "获取通知列表失败"+err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(gin.H{
		"total":         total,
		"notifications": notifications,
		"pageNum":       pageNum,
		"pageSize":      pageSize,
	}))
}

// GetNotificationDetail 获取通知详情
func GetNotificationDetail(c *gin.Context) {
	userId, _ := c.Get("userId")
	notificationId, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var notification entity.Notification
	err := config.MysqlDataBase.Where("id = ? AND receiver_id = ?", notificationId, userId).First(&notification).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "通知不存在或无权查看"))
		return
	}

	// 如果通知未读，则标记为已读
	if !notification.IsRead {
		now := time.Now()
		notification.IsRead = true
		notification.ReadTime = &now
		config.MysqlDataBase.Save(&notification)
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(notification))
}

// MarkNotificationAsRead 标记通知为已读
func MarkNotificationAsRead(c *gin.Context) {
	userId, _ := c.Get("userId")
	notificationId, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var notification entity.Notification
	err := config.MysqlDataBase.Where("id = ? AND receiver_id = ?", notificationId, userId).First(&notification).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "通知不存在或无权操作"))
		return
	}

	if notification.IsRead {
		c.JSON(http.StatusOK, entity.SuccessResponse("通知已经是已读状态"))
		return
	}

	now := time.Now()
	notification.IsRead = true
	notification.ReadTime = &now
	err = config.MysqlDataBase.Save(&notification).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "标记已读失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse("标记已读成功"))
}

// MarkAllNotificationsAsRead 标记所有通知为已读
func MarkAllNotificationsAsRead(c *gin.Context) {
	userId, _ := c.Get("userId")
	notificationType, _ := strconv.Atoi(c.DefaultQuery("type", "0")) // 0表示全部类型

	db := config.MysqlDataBase.Model(&entity.Notification{}).Where("receiver_id = ? AND is_read = ?", userId, false)

	// 根据类型筛选
	if notificationType > 0 {
		db = db.Where("type = ?", notificationType)
	}

	now := time.Now()
	err := db.Updates(map[string]interface{}{
		"is_read":   true,
		"read_time": now,
	}).Error

	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "标记全部已读失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse("标记全部已读成功"))
}

// DeleteNotification 删除通知
func DeleteNotification(c *gin.Context) {
	userId, _ := c.Get("userId")
	notificationId, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var notification entity.Notification
	err := config.MysqlDataBase.Where("id = ? AND receiver_id = ?", notificationId, userId).First(&notification).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "通知不存在或无权操作"))
		return
	}

	err = config.MysqlDataBase.Delete(&notification).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "删除通知失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse("删除通知成功"))
}

// DeleteAllNotifications 删除所有通知
func DeleteAllNotifications(c *gin.Context) {
	userId, _ := c.Get("userId")
	notificationType, _ := strconv.Atoi(c.DefaultQuery("type", "0")) // 0表示全部类型
	isRead := c.DefaultQuery("isRead", "")                           // 空字符串表示全部，"true"表示已读，"false"表示未读

	db := config.MysqlDataBase.Where("receiver_id = ?", userId)

	// 根据类型筛选
	if notificationType > 0 {
		db = db.Where("type = ?", notificationType)
	}

	// 根据已读状态筛选
	if isRead == "true" {
		db = db.Where("is_read = ?", true)
	} else if isRead == "false" {
		db = db.Where("is_read = ?", false)
	}

	err := db.Delete(&entity.Notification{}).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "删除通知失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse("删除通知成功"))
}

// GetUnreadNotificationCount 获取未读通知数量
func GetUnreadNotificationCount(c *gin.Context) {
	userId, _ := c.Get("userId")

	var counts struct {
		Total                 int64 `json:"total"`
		System                int64 `json:"system"`
		ThesisProgress        int64 `json:"thesisProgress"`
		ReviewFeedback        int64 `json:"reviewFeedback"`
		Milestone             int64 `json:"milestone"`
		CoauthoringInvitation int64 `json:"coauthoringInvitation"`
		FormatTemplateUpdate  int64 `json:"formatTemplateUpdate"`
		DefenseSchedule       int64 `json:"defenseSchedule"`
		ThesisSubmission      int64 `json:"thesisSubmission"`
	}

	// 获取总未读数
	config.MysqlDataBase.Model(&entity.Notification{}).
		Where("receiver_id = ? AND is_read = ?", userId, false).
		Count(&counts.Total)

	// 获取各类型未读数
	config.MysqlDataBase.Model(&entity.Notification{}).
		Where("receiver_id = ? AND is_read = ? AND type = ?", userId, false, entity.SystemNotification).
		Count(&counts.System)

	config.MysqlDataBase.Model(&entity.Notification{}).
		Where("receiver_id = ? AND is_read = ? AND type = ?", userId, false, entity.ThesisProgressNotification).
		Count(&counts.ThesisProgress)

	config.MysqlDataBase.Model(&entity.Notification{}).
		Where("receiver_id = ? AND is_read = ? AND type = ?", userId, false, entity.ReviewFeedbackNotification).
		Count(&counts.ReviewFeedback)

	config.MysqlDataBase.Model(&entity.Notification{}).
		Where("receiver_id = ? AND is_read = ? AND type = ?", userId, false, entity.MilestoneNotification).
		Count(&counts.Milestone)

	config.MysqlDataBase.Model(&entity.Notification{}).
		Where("receiver_id = ? AND is_read = ? AND type = ?", userId, false, entity.CoauthoringInvitation).
		Count(&counts.CoauthoringInvitation)

	config.MysqlDataBase.Model(&entity.Notification{}).
		Where("receiver_id = ? AND is_read = ? AND type = ?", userId, false, entity.FormatTemplateUpdate).
		Count(&counts.FormatTemplateUpdate)

	config.MysqlDataBase.Model(&entity.Notification{}).
		Where("receiver_id = ? AND is_read = ? AND type = ?", userId, false, entity.DefenseScheduleNotification).
		Count(&counts.DefenseSchedule)

	config.MysqlDataBase.Model(&entity.Notification{}).
		Where("receiver_id = ? AND is_read = ? AND type = ?", userId, false, entity.ThesisSubmissionNotification).
		Count(&counts.ThesisSubmission)

	c.JSON(http.StatusOK, entity.SuccessResponse(counts))
}

// UpdateNotificationSettings 更新通知设置
func UpdateNotificationSettings(c *gin.Context) {
	userId, _ := c.Get("userId")

	var settings entity.NotificationSetting
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "请求参数错误"))
		return
	}

	// 确保设置的是当前用户的设置
	settings.UserID = uint(userId.(int))

	// 查找是否已有设置
	var existingSetting entity.NotificationSetting
	result := config.MysqlDataBase.Where("user_id = ?", userId).First(&existingSetting)

	if result.Error == nil {
		// 已有设置，更新
		existingSetting.SystemNotification = settings.SystemNotification
		existingSetting.ThesisProgressNotification = settings.ThesisProgressNotification
		existingSetting.ReviewFeedbackNotification = settings.ReviewFeedbackNotification
		existingSetting.MilestoneNotification = settings.MilestoneNotification
		existingSetting.CoauthoringInvitation = settings.CoauthoringInvitation
		existingSetting.FormatTemplateUpdate = settings.FormatTemplateUpdate
		existingSetting.DefenseScheduleNotification = settings.DefenseScheduleNotification
		existingSetting.ThesisSubmissionNotification = settings.ThesisSubmissionNotification
		existingSetting.EmailNotification = settings.EmailNotification
		existingSetting.PushNotification = settings.PushNotification

		err := config.MysqlDataBase.Save(&existingSetting).Error
		if err != nil {
			c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "更新通知设置失败"))
			return
		}
	} else {
		// 没有设置，创建新的
		err := config.MysqlDataBase.Create(&settings).Error
		if err != nil {
			c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "创建通知设置失败"))
			return
		}
	}

	c.JSON(http.StatusOK, entity.SuccessResponse("更新通知设置成功"))
}

// GetNotificationSettings 获取通知设置
func GetNotificationSettings(c *gin.Context) {
	userId, _ := c.Get("userId")

	var settings entity.NotificationSetting
	result := config.MysqlDataBase.Where("user_id = ?", userId).First(&settings)

	if result.Error != nil {
		// 没有找到设置，返回默认设置
		settings = entity.NotificationSetting{
			UserID:                       uint(userId.(int)),
			SystemNotification:           true,
			ThesisProgressNotification:   true,
			ReviewFeedbackNotification:   true,
			MilestoneNotification:        true,
			CoauthoringInvitation:        true,
			FormatTemplateUpdate:         true,
			DefenseScheduleNotification:  true,
			ThesisSubmissionNotification: true,
			EmailNotification:            false,
			PushNotification:             true,
		}
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(settings))
}

// SendNotification 发送通知（内部方法，可供其他服务调用）
func SendNotification(notification entity.Notification) error {
	return config.MysqlDataBase.Create(&notification).Error
}

// SendSystemNotification 发送系统通知
func SendSystemNotification(c *gin.Context) {
	// 仅管理员可调用
	userId, _ := c.Get("userId")

	// 检查是否为管理员
	var isAdmin bool
	err := config.MysqlDataBase.Raw("SELECT is_admin FROM user WHERE id = ?", userId).Scan(&isAdmin).Error
	if err != nil || !isAdmin {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "无权操作，仅管理员可发送系统通知"))
		return
	}

	var req struct {
		Title       string `json:"title" binding:"required"`
		Content     string `json:"content" binding:"required"`
		UserIDs     []uint `json:"userIds"` // 为空则发送给所有用户
		RelatedID   uint   `json:"relatedId"`
		RelatedType string `json:"relatedType"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "请求参数错误"))
		return
	}

	// 开始事务
	tx := config.MysqlDataBase.Begin()

	// 如果未指定用户，则发送给所有用户
	if len(req.UserIDs) == 0 {
		var userIDs []uint
		err := tx.Model(&entity.User{}).Pluck("id", &userIDs).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "获取用户列表失败"))
			return
		}
		req.UserIDs = userIDs
	}

	// 批量创建通知
	for _, receiverID := range req.UserIDs {
		notification := entity.Notification{
			Type:        entity.SystemNotification,
			Title:       req.Title,
			Content:     req.Content,
			SenderID:    uint(userId.(int)),
			ReceiverID:  receiverID,
			RelatedID:   req.RelatedID,
			RelatedType: req.RelatedType,
		}

		if err := tx.Create(&notification).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "发送通知失败"))
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "发送通知失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse("发送系统通知成功"))
}
