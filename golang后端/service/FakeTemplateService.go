package service

import (
	"example/config"
	"example/entity"
	"example/utility"
	"strconv"

	"github.com/gin-gonic/gin"
)

//  本人手搓已黑化 存个string意思意思得了

// 创建模板请求体
type CreateTemplateRequestBody struct {
	Title       string `json:"title" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
}

// 创建模板请求体
type ForkTemplateRequestBody struct {
	ForkCode    string `json:"forkCode" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// 更新模板请求体
type UpdateTemplateRequestBody struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

// FakeTemplateCreate 创建模板
func FakeTemplateCreate(c *gin.Context) {
	var body CreateTemplateRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "提交了错误的JSON"))
		return
	}

	// 获取当前用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}
	userIdInt := utility.GetUserIDFromContext(userId)
	code := utility.GenerateRandomString(10)
	template := entity.FakeTemplate{
		Title:       body.Title,
		Type:        body.Type,
		Description: body.Description,
		UserId:      uint(userIdInt),
		Content:     body.Content,
		Code:        code,
	}

	if err := config.MysqlDataBase.Create(&template).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "创建模板失败: "+err.Error()))
		return
	}
	history := entity.TemplateHistory{
		TemplateId: template.ID,
		Content:    body.Content,
	}
	if err := config.MysqlDataBase.Create(&history).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "创建模板历史失败: "+err.Error()))
		return
	}
	c.JSON(200, entity.SuccessResponse(template))
}

// FakeTemplateGetList 获取模板列表
func FakeTemplateGetList(c *gin.Context) {
	var templates []entity.FakeTemplate

	// 获取当前用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}

	if err := config.MysqlDataBase.Where("user_id = ?", userId).Find(&templates).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取模板列表失败: "+err.Error()))
		return
	}

	c.JSON(200, entity.SuccessResponse(templates))
}

// FakeTemplateGetById 根据ID获取模板
func FakeTemplateGetById(c *gin.Context) {
	id := c.Param("id")

	var template entity.FakeTemplate
	if err := config.MysqlDataBase.First(&template, id).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "模板不存在"))
		return
	}

	// 获取当前用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}
	userIdInt := utility.GetUserIDFromContext(userId)
	// 检查模板是否属于当前用户
	if template.UserId != uint(userIdInt) {
		var guide entity.Guide
		if err := config.MysqlDataBase.Where("student_user_id = ? AND teacher_user_id = ? AND status = 1", template.UserId, userIdInt).First(&guide).Error; err != nil {
			c.JSON(200, entity.ErrorResponse[string](403, "无权访问此模板"))
			return
		}
	}

	c.JSON(200, entity.SuccessResponse(template))
}

// FakeTemplateUpdate 更新模板
func FakeTemplateUpdate(c *gin.Context) {
	id := c.Param("id")
	templateId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "无效的模板ID"))
		return
	}

	var body UpdateTemplateRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "提交了错误的JSON"))
		return
	}

	// 获取当前用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}
	userIdInt := utility.GetUserIDFromContext(userId)

	// 查找模板
	var template entity.FakeTemplate
	if err := config.MysqlDataBase.First(&template, templateId).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "模板不存在"))
		return
	}

	// 检查模板是否属于当前用户
	if template.UserId != uint(userIdInt) {
		var guide entity.Guide
		if err := config.MysqlDataBase.Where("student_user_id = ? AND teacher_user_id = ? AND status = 1", template.UserId, userIdInt).First(&guide).Error; err != nil {
			c.JSON(200, entity.ErrorResponse[string](403, "无权修改此模板"))
			return
		}
	}

	// 更新模板
	updates := map[string]interface{}{}
	if body.Title != "" {
		updates["title"] = body.Title
	}
	if body.Type != "" {
		updates["type"] = body.Type
	}
	if body.Description != "" {
		updates["description"] = body.Description
	}
	if body.Content != "" {
		updates["content"] = body.Content
	}

	if err := config.MysqlDataBase.Model(&template).Updates(updates).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "更新模板失败: "+err.Error()))
		return
	}
	history := entity.TemplateHistory{
		TemplateId: template.ID,
		Content:    body.Content,
	}
	// 获取更新后的模板
	if err := config.MysqlDataBase.First(&template, templateId).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取更新后的模板失败"))
		return
	}

	if err := config.MysqlDataBase.Create(&history).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "创建模板历史失败: "+err.Error()))
		return
	}
	c.JSON(200, entity.SuccessResponse(template))
}

// FakeTemplateDelete 删除模板
func FakeTemplateDelete(c *gin.Context) {
	id := c.Param("id")
	templateId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "无效的模板ID"))
		return
	}

	// 获取当前用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}

	// 查找模板
	var template entity.FakeTemplate
	if err := config.MysqlDataBase.First(&template, templateId).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "模板不存在"))
		return
	}

	// 检查模板是否属于当前用户
	if template.UserId != userId.(uint) {
		c.JSON(200, entity.ErrorResponse[string](403, "无权删除此模板"))
		return
	}

	// 删除模板
	if err := config.MysqlDataBase.Delete(&template).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "删除模板失败: "+err.Error()))
		return
	}

	c.JSON(200, entity.SuccessResponse("模板删除成功"))
}

// FakeTemplateFork 创建模板
func FakeTemplateFork(c *gin.Context) {
	var body ForkTemplateRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "提交了错误的JSON"))
		return
	}

	// 获取当前用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "用户未登录"))
		return
	}
	userIdInt := utility.GetUserIDFromContext(userId)
	var templateOrigin entity.FakeTemplate
	if err := config.MysqlDataBase.Where("code = ?", body.ForkCode).First(&templateOrigin).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](404, "模板不存在"))
		return
	}
	template := entity.FakeTemplate{
		Title:       body.Title,
		Type:        templateOrigin.Type,
		Description: body.Description,
		UserId:      uint(userIdInt),
		Content:     templateOrigin.Content,
	}

	if err := config.MysqlDataBase.Create(&template).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "创建模板失败: "+err.Error()))
		return
	}
	history := entity.TemplateHistory{
		TemplateId: template.ID,
		Content:    template.Content,
	}
	if err := config.MysqlDataBase.Create(&history).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "创建模板历史失败: "+err.Error()))
		return
	}
	c.JSON(200, entity.SuccessResponse(template))
}

func FakeTemplateGetListStudent(c *gin.Context) {
	userId, _ := c.Get("userId")
	userIdInt := utility.GetUserIDFromContext(userId)
	var templates []entity.FakeTemplate

	// 查找该教师指导的所有学生（status为1表示指导关系已确认）
	var guides []entity.Guide
	if err := config.MysqlDataBase.Where("teacher_user_id = ? AND status = 1", userIdInt).Find(&guides).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取指导关系失败: "+err.Error()))
		return
	}

	// 如果没有指导的学生，返回空数组
	if len(guides) == 0 {
		c.JSON(200, entity.SuccessResponse(templates))
		return
	}

	// 提取所有学生的ID
	var studentIds []uint
	for _, guide := range guides {
		studentIds = append(studentIds, guide.StudentUserId)
	}

	// 查询这些学生的所有模板
	if err := config.MysqlDataBase.Where("user_id IN ?", studentIds).Preload("User").Find(&templates).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取学生模板列表失败: "+err.Error()))
		return
	}

	c.JSON(200, entity.SuccessResponse(templates))
}

func GetTemplateHistortyIDs(c *gin.Context) {
	id := c.Param("id")
	var templateHistory []entity.TemplateHistory
	if err := config.MysqlDataBase.Omit("content").Where("template_id = ?", id).Order("id desc").Find(&templateHistory).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取模板历史记录列表失败: "+err.Error()))
		return
	}
	c.JSON(200, entity.SuccessResponse(templateHistory))
}

func GetTemplateHistoryContent(c *gin.Context) {
	id := c.Param("id")
	var content string
	if err := config.MysqlDataBase.Table("template_histories").Select("content").Where("id = ?", id).Find(&content).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "获取模板历史内容失败: "+err.Error()))
		return
	}
	c.JSON(200, entity.SuccessResponse(content))
}
