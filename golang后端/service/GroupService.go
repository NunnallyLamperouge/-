package service

import (
	"example/config"
	"example/entity"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateChatGroup 创建聊天群组
func CreateChatGroup(c *gin.Context) {
	userId, _ := c.Get("userId")

	var group entity.ChatGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "请求参数错误"))
		return
	}

	// 设置创建者ID
	group.CreatorID = uint(userId.(int))
	group.MemberCount = 1

	// 开始事务
	tx := config.MysqlDataBase.Begin()

	// 创建群组
	if err := tx.Create(&group).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "创建群组失败"))
		return
	}

	// 添加创建者为群组成员且设为管理员
	member := entity.GroupMember{
		GroupID:  group.ID,
		UserID:   uint(userId.(int)),
		JoinTime: time.Now(),
		IsAdmin:  true,
		Status:   1, // 正常状态
	}

	if err := tx.Create(&member).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "创建群组失败"))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "创建群组失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(group))
}

// GetUserChatGroups 获取用户加入的群组列表
func GetUserChatGroups(c *gin.Context) {
	userId, _ := c.Get("userId")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.DefaultQuery("keyword", "")

	offset := (pageNum - 1) * pageSize
	var groups []entity.ChatGroup
	var total int64

	// 查询用户所在的群组ID
	var groupIds []uint
	err := config.MysqlDataBase.Model(&entity.GroupMember{}).
		Where("user_id = ? AND status != 3", userId). // 排除已退出的群组
		Pluck("group_id", &groupIds).Error

	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "获取群组列表失败"))
		return
	}

	if len(groupIds) == 0 {
		c.JSON(http.StatusOK, entity.SuccessResponse(gin.H{
			"total":    int64(0),
			"groups":   []entity.ChatGroup{},
			"pageNum":  pageNum,
			"pageSize": pageSize,
		}))
		return
	}

	// 查询群组信息
	db := config.MysqlDataBase.Model(&entity.ChatGroup{}).Where("id IN ?", groupIds)

	// 如果有关键字，添加模糊搜索
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 获取总数
	db.Count(&total)

	// 获取分页数据
	err = db.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&groups).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "获取群组列表失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(gin.H{
		"total":    total,
		"groups":   groups,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	}))
}

// GetChatGroupDetail 获取群组详情
func GetChatGroupDetail(c *gin.Context) {
	userId, _ := c.Get("userId")
	groupId, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// 检查用户是否是群成员
	var member entity.GroupMember
	err := config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND status != 3", groupId, userId).First(&member).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您不是该群成员或该群不存在"))
		return
	}

	var group entity.ChatGroup
	err = config.MysqlDataBase.First(&group, groupId).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "群组不存在"))
		return
	}

	// 获取成员信息，仅返回前10个
	var members []entity.GroupMember
	config.MysqlDataBase.Where("group_id = ? AND status != 3", groupId).
		Order("is_admin DESC, join_time ASC").
		Limit(10).
		Find(&members)

	c.JSON(http.StatusOK, entity.SuccessResponse(gin.H{
		"group":         group,
		"userRole":      member,
		"memberPreview": members,
		"totalMembers":  group.MemberCount,
		"isCreator":     group.CreatorID == uint(userId.(int)),
		"isAdmin":       member.IsAdmin,
		"currentStatus": member.Status,
	}))
}

// UpdateChatGroup 更新群组信息
func UpdateChatGroup(c *gin.Context) {
	userId, _ := c.Get("userId")
	groupId, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// 检查用户是否是管理员
	var member entity.GroupMember
	err := config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND is_admin = ?", groupId, userId, true).First(&member).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您不是该群管理员或该群不存在"))
		return
	}

	var updateData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		AvatarURL   string `json:"avatarUrl"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "请求参数错误"))
		return
	}

	var group entity.ChatGroup
	err = config.MysqlDataBase.First(&group, groupId).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "群组不存在"))
		return
	}

	// 更新群组信息
	if updateData.Name != "" {
		group.Name = updateData.Name
	}
	if updateData.Description != "" {
		group.Description = updateData.Description
	}
	if updateData.AvatarURL != "" {
		group.AvatarURL = updateData.AvatarURL
	}

	err = config.MysqlDataBase.Save(&group).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "更新群组信息失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse("更新群组信息成功"))
}

// DeleteChatGroup 删除/解散群组
func DeleteChatGroup(c *gin.Context) {
	userId, _ := c.Get("userId")
	groupId, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// 检查用户是否是群主
	var group entity.ChatGroup
	err := config.MysqlDataBase.Where("id = ? AND creator_id = ?", groupId, userId).First(&group).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您不是该群创建者或该群不存在"))
		return
	}

	// 开始事务
	tx := config.MysqlDataBase.Begin()

	// 删除所有群成员
	if err := tx.Where("group_id = ?", groupId).Delete(&entity.GroupMember{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "解散群组失败"))
		return
	}

	// 删除群组
	if err := tx.Delete(&group).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "解散群组失败"))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "解散群组失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse("解散群组成功"))
}

// AddGroupMember 添加群成员
func AddGroupMember(c *gin.Context) {
	userId, _ := c.Get("userId")
	groupId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		UserIDs []uint `json:"userIds" binding:"omitempty"`
		Email   string `json:"email" binding:"omitempty,email,max=50"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "请求参数错误"+err.Error()))
		return
	}

	// 检查是否同时提供了UserIDs和Email
	if len(req.UserIDs) == 0 && req.Email == "" {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "请提供要邀请的用户ID列表或邮箱"))
		return
	}

	// 检查用户是否是管理员
	var member entity.GroupMember
	err := config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND is_admin = ?", groupId, userId, true).First(&member).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您不是该群管理员或该群不存在"))
		return
	}

	// 检查群组是否存在
	var group entity.ChatGroup
	err = config.MysqlDataBase.First(&group, groupId).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "群组不存在"))
		return
	}

	// 开始事务
	tx := config.MysqlDataBase.Begin()

	// 添加新成员
	now := time.Now()
	addedCount := 0

	// 处理Email邀请
	if req.Email != "" {
		// 查找用户
		var user entity.User
		result := tx.Where("email = ?", req.Email).First(&user)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "未找到该邮箱对应的用户"))
			return
		}

		// 检查用户是否已在群组中
		var existingMember entity.GroupMember
		result = tx.Where("group_id = ? AND user_id = ?", groupId, user.ID).First(&existingMember)

		if result.Error == nil {
			// 用户已存在，如果是已退出状态则重新激活
			if existingMember.Status == 3 {
				existingMember.Status = 1
				existingMember.JoinTime = now
				if err := tx.Save(&existingMember).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "添加群成员失败"))
					return
				}
				addedCount++
			} else {
				tx.Rollback()
				c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "该用户已是群成员"))
				return
			}
		} else {
			// 创建新成员
			newMember := entity.GroupMember{
				GroupID:  uint(groupId),
				UserID:   user.ID,
				JoinTime: now,
				IsAdmin:  false,
				Status:   1,
			}

			if err := tx.Create(&newMember).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "添加群成员失败"))
				return
			}

			// 发送通知给新成员
			notification := entity.Notification{
				Type:        entity.SystemNotification,
				Title:       "群组邀请",
				Content:     "您已被邀请加入群组: " + group.Name,
				SenderID:    uint(userId.(int)),
				ReceiverID:  user.ID,
				RelatedID:   uint(groupId),
				RelatedType: "chat_group",
			}

			tx.Create(&notification)
			addedCount++
		}
	} else {
		// 处理UserIDs邀请
		for _, userID := range req.UserIDs {
			// 检查用户是否已是群成员
			var existingMember entity.GroupMember
			result := tx.Where("group_id = ? AND user_id = ?", groupId, userID).First(&existingMember)

			if result.Error == nil {
				// 用户已存在，如果是已退出状态则重新激活
				if existingMember.Status == 3 {
					existingMember.Status = 1
					existingMember.JoinTime = now
					if err := tx.Save(&existingMember).Error; err != nil {
						continue
					}
					addedCount++
				}
				continue
			}

			// 创建新成员
			newMember := entity.GroupMember{
				GroupID:  uint(groupId),
				UserID:   userID,
				JoinTime: now,
				IsAdmin:  false,
				Status:   1,
			}

			if err := tx.Create(&newMember).Error; err != nil {
				continue
			}

			// 发送通知给新成员
			notification := entity.Notification{
				Type:        entity.SystemNotification,
				Title:       "群组邀请",
				Content:     "您已被邀请加入群组: " + group.Name,
				SenderID:    uint(userId.(int)),
				ReceiverID:  userID,
				RelatedID:   uint(groupId),
				RelatedType: "chat_group",
			}

			tx.Create(&notification)
			addedCount++
		}
	}

	// 更新群组成员数量
	if addedCount > 0 {
		group.MemberCount += addedCount
		if err := tx.Save(&group).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "添加群成员失败"))
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "添加群成员失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(gin.H{
		"message": "添加群成员成功",
		"added":   addedCount,
	}))
}

// RemoveGroupMember 移除群成员
func RemoveGroupMember(c *gin.Context) {
	userId, _ := c.Get("userId")
	groupId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	targetUserId, _ := strconv.ParseUint(c.Param("userId"), 10, 64)

	// 不能移除自己，应该使用退出群组功能
	if uint(userId.(int)) == uint(targetUserId) {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "不能移除自己，请使用退出群组功能"))
		return
	}

	// 检查用户是否是管理员
	var currentMember entity.GroupMember
	err := config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND is_admin = ?", groupId, userId, true).First(&currentMember).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您不是该群管理员或该群不存在"))
		return
	}

	// 检查目标用户是否在群中
	var targetMember entity.GroupMember
	err = config.MysqlDataBase.Where("group_id = ? AND user_id = ?", groupId, targetUserId).First(&targetMember).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "该用户不在群组中"))
		return
	}

	// 群主不能被管理员移除
	var group entity.ChatGroup
	config.MysqlDataBase.First(&group, groupId)
	if group.CreatorID == uint(targetUserId) && uint(userId.(int)) != group.CreatorID {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您无权移除群主"))
		return
	}

	// 开始事务
	tx := config.MysqlDataBase.Begin()

	// 标记成员为退出状态
	targetMember.Status = 3 // 退出状态
	if err := tx.Save(&targetMember).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "移除成员失败"))
		return
	}

	// 更新群组成员数量
	group.MemberCount--
	if err := tx.Save(&group).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "移除成员失败"))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "移除成员失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse("移除成员成功"))
}

// LeaveGroup 退出群组
func LeaveGroup(c *gin.Context) {
	userId, _ := c.Get("userId")
	groupId, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// 检查群组是否存在
	var group entity.ChatGroup
	err := config.MysqlDataBase.First(&group, groupId).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "群组不存在"))
		return
	}

	// 群主不能退出，只能解散群组
	if group.CreatorID == uint(userId.(int)) {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "群主不能退出群组，请使用解散群组功能"))
		return
	}

	// 检查用户是否在群中
	var member entity.GroupMember
	err = config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND status != 3", groupId, userId).First(&member).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "您不在该群组中"))
		return
	}

	// 开始事务
	tx := config.MysqlDataBase.Begin()

	// 标记成员为退出状态
	member.Status = 3 // 退出状态
	if err := tx.Save(&member).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "退出群组失败"))
		return
	}

	// 更新群组成员数量
	group.MemberCount--
	if err := tx.Save(&group).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "退出群组失败"))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "退出群组失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse("退出群组成功"))
}

// GetGroupMembers 获取群组成员列表
func GetGroupMembers(c *gin.Context) {
	userId, _ := c.Get("userId")
	groupId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.DefaultQuery("keyword", "")

	// 检查用户是否在群中
	var member entity.GroupMember
	err := config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND status != 3", groupId, userId).First(&member).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您不是该群成员或该群不存在"))
		return
	}

	offset := (pageNum - 1) * pageSize
	var members []struct {
		entity.GroupMember
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}
	var total int64

	// 构建查询
	query := config.MysqlDataBase.Table("group_members").
		Select("group_members.*, users.username, users.avatar").
		Joins("LEFT JOIN users ON group_members.user_id = users.id").
		Where("group_members.group_id = ? AND group_members.status != 3", groupId)

	// 如果有关键字，添加模糊搜索
	if keyword != "" {
		query = query.Where("users.username LIKE ? OR group_members.nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 获取分页数据
	err = query.Order("group_members.is_admin DESC, group_members.join_time ASC").
		Offset(offset).Limit(pageSize).Find(&members).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "获取群成员列表失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(gin.H{
		"total":    total,
		"members":  members,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	}))
}

// UpdateGroupMemberInfo 更新群成员信息
func UpdateGroupMemberInfo(c *gin.Context) {
	userId, _ := c.Get("userId")
	groupId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	targetUserId, _ := strconv.ParseUint(c.Param("userId"), 10, 64)

	var req struct {
		Nickname string `json:"nickname"`
		IsAdmin  *bool  `json:"isAdmin,omitempty"`
		Status   *int   `json:"status,omitempty"` // 1:正常 2:禁言
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "请求参数错误"))
		return
	}

	// 检查当前用户权限
	var currentMember entity.GroupMember
	err := config.MysqlDataBase.Where("group_id = ? AND user_id = ?", groupId, userId).First(&currentMember).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您不是该群成员或该群不存在"))
		return
	}

	// 检查目标用户是否在群中
	var targetMember entity.GroupMember
	err = config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND status != 3", groupId, targetUserId).First(&targetMember).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "该用户不在群组中"))
		return
	}

	// 获取群组信息
	var group entity.ChatGroup
	config.MysqlDataBase.First(&group, groupId)

	// 权限检查：只有自己可以修改自己的昵称
	if req.Nickname != "" && uint(userId.(int)) != uint(targetUserId) {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您只能修改自己的群昵称"))
		return
	}

	// 权限检查：管理员相关操作只有群主可以执行
	if req.IsAdmin != nil && group.CreatorID != uint(userId.(int)) {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "只有群主可以设置管理员"))
		return
	}

	// 权限检查：禁言等操作需要管理员权限
	if req.Status != nil && !currentMember.IsAdmin {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您不是管理员，无权执行此操作"))
		return
	}

	// 权限检查：不能对群主执行限制操作
	if (req.Status != nil || req.IsAdmin != nil) && group.CreatorID == uint(targetUserId) && uint(userId.(int)) != group.CreatorID {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "无法对群主执行此操作"))
		return
	}

	// 更新成员信息
	if req.Nickname != "" {
		targetMember.Nickname = req.Nickname
	}

	if req.IsAdmin != nil {
		targetMember.IsAdmin = *req.IsAdmin
	}

	if req.Status != nil {
		if *req.Status == 1 || *req.Status == 2 {
			targetMember.Status = *req.Status
		}
	}

	err = config.MysqlDataBase.Save(&targetMember).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "更新成员信息失败"))
		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse("更新成员信息成功"))
}

// MuteGroupMember 禁言/解除禁言群成员
func MuteGroupMember(c *gin.Context) {
	userId, _ := c.Get("userId")
	groupId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	targetUserId, _ := strconv.ParseUint(c.Param("userId"), 10, 64)
	action := c.Query("action") // "mute" 或 "unmute"

	if action != "mute" && action != "unmute" {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "无效的操作类型"))
		return
	}

	// 检查当前用户是否是管理员
	var currentMember entity.GroupMember
	err := config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND is_admin = ?", groupId, userId, true).First(&currentMember).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您不是该群管理员或该群不存在"))
		return
	}

	// 检查目标用户是否在群中
	var targetMember entity.GroupMember
	err = config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND status != 3", groupId, targetUserId).First(&targetMember).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "该用户不在群组中"))
		return
	}

	// 获取群组信息
	var group entity.ChatGroup
	config.MysqlDataBase.First(&group, groupId)

	// 不能禁言群主
	if group.CreatorID == uint(targetUserId) {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "无法对群主执行此操作"))
		return
	}

	// 普通管理员不能禁言其他管理员
	if targetMember.IsAdmin && uint(userId.(int)) != group.CreatorID {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "无法对其他管理员执行此操作"))
		return
	}

	// 更新成员状态
	if action == "mute" {
		targetMember.Status = 2 // 禁言状态
	} else {
		targetMember.Status = 1 // 正常状态
	}

	err = config.MysqlDataBase.Save(&targetMember).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "操作失败"))
		return
	}

	message := "解除禁言成功"
	if action == "mute" {
		message = "禁言成功"
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(message))
}

// SetGroupAdmin 设置/取消群管理员
func SetGroupAdmin(c *gin.Context) {
	userId, _ := c.Get("userId")
	groupId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	targetUserId, _ := strconv.ParseUint(c.Param("userId"), 10, 64)
	action := c.Query("action") // "set" 或 "cancel"

	if action != "set" && action != "cancel" {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "无效的操作类型"))
		return
	}

	// 检查当前用户是否是群主
	var group entity.ChatGroup
	err := config.MysqlDataBase.Where("id = ? AND creator_id = ?", groupId, userId).First(&group).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您不是该群群主或该群不存在"))
		return
	}

	// 不能对自己执行此操作
	if uint(userId.(int)) == uint(targetUserId) {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](400, "不能对自己执行此操作"))
		return
	}

	// 检查目标用户是否在群中
	var targetMember entity.GroupMember
	err = config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND status != 3", groupId, targetUserId).First(&targetMember).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](404, "该用户不在群组中"))
		return
	}

	// 更新管理员状态
	isAdmin := action == "set"

	// 如果状态已经是目标状态
	if targetMember.IsAdmin == isAdmin {
		message := "该用户已经是管理员"
		if !isAdmin {
			message = "该用户已经不是管理员"
		}
		c.JSON(http.StatusOK, entity.SuccessResponse(message))
		return
	}

	targetMember.IsAdmin = isAdmin
	err = config.MysqlDataBase.Save(&targetMember).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "操作失败"))
		return
	}

	message := "设置管理员成功"
	if action == "cancel" {
		message = "取消管理员成功"
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(message))
}

// GetGroupMessages 获取群组消息
func GetGroupMessages(c *gin.Context) {
	userID, _ := c.Get("userId")
	groupID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	// 验证用户是否为群组成员
	var member entity.GroupMember
	err := config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND status != 3", groupID, userID).First(&member).Error
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](403, "您不是该群成员或该群不存在"))
		return
	}

	// 查询消息
	var messages []struct {
		entity.Message
		SenderName   string `json:"senderName"`
		SenderAvatar string `json:"senderAvatar"`
	}

	err = config.MysqlDataBase.Table("messages").
		Select("messages.*, users.username as sender_name, users.avatar as sender_avatar").
		Joins("LEFT JOIN users ON messages.sender_id = users.id").
		Where("messages.group_id = ? AND messages.deleted_at IS NULL", groupID).
		Order("messages.created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&messages).Error

	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "获取消息失败"))
		return
	}

	// 将结果反转，使最新消息在底部
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(messages))
}
