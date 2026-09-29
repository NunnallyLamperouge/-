package service

import (
	"example/config"
	"example/entity"
	"example/utility"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUserWorkspace 创建用户的AnythingLLM工作区
func CreateUserWorkspace(c *gin.Context) {
	// 1. 获取用户ID
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "未认证用户"))
		return
	}
	userId := uint(utility.GetUserIDFromContext(userIdInterface))

	// 2. 检查用户是否已有工作区
	var user entity.User
	if err := config.MysqlDataBase.Table("users").Where("id = ?", userId).Find(&user).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, fmt.Sprintf("获取用户信息失败: %v", err)))
		return
	}

	if user.WorkspaceId != "" {
		c.JSON(200, entity.ErrorResponse[string](400, "用户已有工作区，不能重复创建"))
		return
	}

	// 3. 创建工作区
	workspaceId := fmt.Sprintf("user_%d", userId)
	description := fmt.Sprintf("用户%d的个人知识库", userId)

	// 4. 调用AnythingLLM API创建工作区
	slug, err := utility.CreateWorkspace(workspaceId, description)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, fmt.Sprintf("创建工作区失败: %v", err)))
		return
	}

	// 5. 更新用户的工作区ID
	if err := config.MysqlDataBase.Table("users").Where("id = ?", userId).Update("workspace_id", slug).Error; err != nil {
		// 如果更新失败，尝试删除已创建的工作区（这里可以添加删除工作区的逻辑）
		c.JSON(200, entity.ErrorResponse[string](500, fmt.Sprintf("更新用户工作区信息失败: %v", err)))
		return
	}

	// 6. 返回成功响应
	c.JSON(200, entity.SuccessResponseWithMessage("工作区创建成功", gin.H{
		"workspaceId": slug,
	}))
}

func GetUserWorkspace(c *gin.Context) {
	userId, _ := c.Get("userId")
	userIdInt := utility.GetUserIDFromContext(userId)
	workspaceId := "user_" + strconv.Itoa(userIdInt)
	result, err := utility.QueryWorkspace(workspaceId)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "未创建工作台"))
		return
	}
	c.JSON(200, entity.SuccessResponse(result))
}

func UploadRAGDocument(c *gin.Context) {
	// 1. 获取用户ID
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "未认证用户"))
		return
	}
	userId := uint(utility.GetUserIDFromContext(userIdInterface))

	// 2. 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, fmt.Sprintf("获取文件失败: %v", err)))
		return
	}

	// 3. 调用服务上传文件
	docInfo, err := utility.UploadUserDocument(userId, file)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, fmt.Sprintf("上传文件失败: %v", err)))
		return
	}

	// 4. 返回成功响应
	c.JSON(200, entity.SuccessResponseWithMessage("文件上传成功", docInfo))
}

// GetUserRAGDocuments 获取用户文档列表
func GetUserRAGDocuments(c *gin.Context) {
	// 1. 获取用户ID
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "未认证用户"))
		return
	}
	userId := uint(utility.GetUserIDFromContext(userIdInterface))

	// 2. 获取用户文档列表
	documents, err := utility.GetUserDocuments(userId)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, fmt.Sprintf("获取文档列表失败: %v", err)))
		return
	}

	// 3. 返回文档列表
	c.JSON(200, entity.SuccessResponseWithMessage("获取文档列表成功", documents))
}

// DeleteRAGDocument 删除用户的RAG文档
func DeleteRAGDocument(c *gin.Context) {
	// 1. 获取用户ID
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "未认证用户"))
		return
	}
	userId := uint(utility.GetUserIDFromContext(userIdInterface))

	// 2. 获取要删除的文件名列表
	var request struct {
		FileNames []string `json:"fileNames" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, fmt.Sprintf("请求参数错误: %v", err)))
		return
	}

	// 3. 获取用户工作区ID
	workspaceId, err := utility.GetUserWorkspaceId(userId)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, fmt.Sprintf("获取用户工作区失败: %v", err)))
		return
	}

	// 4. 调用服务删除文档
	if err := utility.DeleteDocumentsFromWorkspace(workspaceId, request.FileNames); err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, fmt.Sprintf("删除文档失败: %v", err)))
		return
	}

	// 5. 返回成功响应
	c.JSON(200, entity.SuccessResponseWithMessage[string]("文档删除成功", ""))
}

// AddRAGDocument 添加文档到用户的RAG工作区
func AddRAGDocument(c *gin.Context) {
	// 1. 获取用户ID
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(200, entity.ErrorResponse[string](401, "未认证用户"))
		return
	}
	userId := uint(utility.GetUserIDFromContext(userIdInterface))

	// 2. 获取要添加的文件名列表
	var request struct {
		FileNames []string `json:"fileNames" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, fmt.Sprintf("请求参数错误: %v", err)))
		return
	}

	// 3. 获取用户工作区ID
	workspaceId, err := utility.GetUserWorkspaceId(userId)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, fmt.Sprintf("获取用户工作区失败: %v", err)))
		return
	}

	// 4. 调用服务添加文档
	if err := utility.AddDocumentsToWorkspace(workspaceId, request.FileNames); err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, fmt.Sprintf("添加文档失败: %v", err)))
		return
	}

	// 5. 返回成功响应
	c.JSON(200, entity.SuccessResponseWithMessage[string]("文档添加成功", ""))
}

type RAGParagraphAdvancedRequestBody struct {
	Source     string `json:"source"`
	ModifyDes  string `json:"modify_description"`
	FieldType  string `json:"field_type"`
	Style      string `json:"style"`
	TemplateId int    `json:"templateId"`
}

func RAGParagraphAdvancedRequest(c *gin.Context) {
	userId, _ := c.Get("userId")
	userIdInt := utility.GetUserIDFromContext(userId)
	var request RAGParagraphAdvancedRequestBody
	var fakeTemplate entity.FakeTemplate
	wkscSlug := fmt.Sprintf("user_%d", userIdInt)
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, fmt.Sprintf("错误的表单:%s请求id:%d", err.Error(), request.TemplateId)))
		return
	}
	if err := config.MysqlDataBase.Table("fake_templates").Where("id = ?", request.TemplateId).First(&fakeTemplate).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, fmt.Sprintf("无法找到对应模版:%s请求id:%d", err.Error(), request.TemplateId)))
		return
	}
	if fakeTemplate.Slug == "" {
		slug := fmt.Sprintf("template_%d", fakeTemplate.ID)
		result, err := utility.CreateNewThread(wkscSlug, slug)
		if err != nil {
			c.JSON(200, entity.ErrorResponse[string](500, "创建工作台thread时发生错误:"+err.Error()))
			return
		}
		fakeTemplate.Slug = result
		if err := config.MysqlDataBase.Save(fakeTemplate).Error; err != nil {
			c.JSON(200, entity.ErrorResponse[string](500, "保存thread信息时发生错误:"+err.Error()))
			return
		}
	}
	questionStr := fmt.Sprintf("这是需要优化的%s部分：%s", request.FieldType, request.Source)

	// 根据字段类型添加不同的提示词
	switch request.FieldType {
	case "title":
		questionStr = fmt.Sprintf("你需要帮助优化论文标题，使其更加专业、准确和吸引人。%s", questionStr)
	case "acknowledgment":
		questionStr = fmt.Sprintf("你需要帮助优化论文致谢部分，使其更加真诚、得体和专业。%s", questionStr)
	case "abstract":
		questionStr = fmt.Sprintf("你需要帮助优化论文摘要，确保其简洁、完整地概括论文的主要内容。%s", questionStr)
	case "content":
		questionStr = fmt.Sprintf("你需要帮助优化论文正文部分，提升其学术性和专业性。%s", questionStr)
	}

	// 添加风格要求
	if request.Style != "" {
		questionStr = fmt.Sprintf("%s\n请按照以下风格要求进行优化：%s", questionStr, request.Style)
	}

	// 添加具体优化要求
	if request.ModifyDes != "" {
		questionStr = fmt.Sprintf("%s\n具体优化要求：%s", questionStr, request.ModifyDes)
	}

	// 添加参考原文
	if request.TemplateId != 0 {
		var source string
		if err := config.MysqlDataBase.Table("fake_templates").Select("content").Where("id = ?", request.TemplateId).Find(&source).Error; err != nil {
			log.Println(err.Error())
		} else {
			questionStr = fmt.Sprintf("%s\n以下是原文内容，仅供参考以保持文章的连贯性和避免产生不实信息：%s", questionStr, source)
		}
	}

	response, err := utility.SendChatMessage(wkscSlug, fakeTemplate.Slug, questionStr)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "内容优化失败"))
		return
	}
	c.JSON(200, entity.SuccessResponse(response))
}

type SwitchFilePinStatusRequestBody struct {
	FileName string `json:"filename"`
	Status   bool   `json:"pin_status"`
}

func SwitchFilePinStatus(c *gin.Context) {
	userId, _ := c.Get("userId")
	userIdInt := utility.GetUserIDFromContext(userId)
	workSpaceSlug := fmt.Sprintf("user_%d", userIdInt)
	var request SwitchFilePinStatusRequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(200, entity.ErrorResponse[string](400, "请求表单有误"))
		return
	}
	_, err := utility.UpdateDocumentPinStatus(request.FileName, workSpaceSlug, request.Status)
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, fmt.Sprintf("发生错误，请稍后重试:%s", err.Error())))
		return
	}
	c.JSON(200, entity.SuccessResponse("更新成功！"))
}

type GetAdveriseUniversalRAGRequest struct {
	Source     string `json:"source"`             // 需要修改的原文内容
	ModifyDes  string `json:"modify_description"` // 修改要求描述
	FieldType  string `json:"field_type"`         // 论文字段类型（如：title, acknowledgment, content等）
	Style      string `json:"style"`              // 可选的风格要求
	TemplateId int    `json:"templateId"`
}
