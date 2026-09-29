package utility

import (
	"bytes"
	"encoding/json"
	"example/config"
	"example/entity"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

/**
 * AnythingLLM 服务适配器
 * 实现对 AnythingLLM API 的调用
 */

// ==================== AnythingLLM 配置 ====================

// AnythingLLMConfig 定义 AnythingLLM 的配置结构
type AnythingLLMConfig struct {
	BaseURL string // AnythingLLM API 的基础URL
	APIKey  string // AnythingLLM API的密钥
}

// 全局变量
var (
	AnythingLLMBaseURL string // AnythingLLM API的基础URL
	AnythingLLMAPIKey  string // AnythingLLM API的密钥
)

// InitAnythingLLM 初始化AnythingLLM配置
func InitAnythingLLM(baseURL, apiKey string) {
	AnythingLLMBaseURL = baseURL
	AnythingLLMAPIKey = apiKey
}

// AnythingLLMWorkspaceDocument AnythingLLM文档上传响应
type AnythingLLMWorkspaceDocument struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	DocumentId  string `json:"documentId,omitempty"`
	DocumentURL string `json:"documentURL,omitempty"`
}

// AnythingLLMThread 定义 AnythingLLM 工作区中的线程信息
type AnythingLLMThread struct {
	UserID *string `json:"user_id"`
	Slug   string  `json:"slug"`
	Name   string  `json:"name"`
}

// AnythingLLMDocument 定义 AnythingLLM 文档结构
type AnythingLLMDocument struct {
	ID            int    `json:"id"`
	DocID         string `json:"docId"`
	Filename      string `json:"filename"`
	DocPath       string `json:"docpath"`
	WorkspaceID   int    `json:"workspaceId"`
	Metadata      string `json:"metadata"`
	Pinned        bool   `json:"pinned"`
	Watched       bool   `json:"watched"`
	CreatedAt     string `json:"createdAt"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
}

// AnythingLLMWorkspace 定义 AnythingLLM 工作区信息
type AnythingLLMWorkspace struct {
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	Slug                 string                `json:"slug"`
	VectorTag            *string               `json:"vectorTag"`
	CreatedAt            string                `json:"createdAt"`
	OpenAiTemp           float64               `json:"openAiTemp"`
	OpenAiHistory        int                   `json:"openAiHistory"`
	LastUpdatedAt        string                `json:"lastUpdatedAt"`
	OpenAiPrompt         string                `json:"openAiPrompt"`
	SimilarityThreshold  float64               `json:"similarityThreshold"`
	ChatProvider         *string               `json:"chatProvider"`
	ChatModel            *string               `json:"chatModel"`
	TopN                 int                   `json:"topN"`
	ChatMode             string                `json:"chatMode"`
	PfpFilename          *string               `json:"pfpFilename"`
	AgentProvider        *string               `json:"agentProvider"`
	AgentModel           *string               `json:"agentModel"`
	QueryRefusalResponse string                `json:"queryRefusalResponse"`
	VectorSearchMode     string                `json:"vectorSearchMode"`
	Threads              []AnythingLLMThread   `json:"threads"`
	Documents            []AnythingLLMDocument `json:"documents"`
}

// UserDocumentInfo 存储上传的用户文档信息
type UserDocumentInfo struct {
	FileName    string    `json:"fileName"`
	FilePath    string    `json:"filePath"`
	FileSize    int64     `json:"fileSize"`
	ContentType string    `json:"contentType"`
	UploadTime  time.Time `json:"uploadTime"`
	UserId      uint      `json:"userId"`
}
type GetWorkspaceResponseBody struct {
	Workspace []AnythingLLMWorkspace `json:"workspace"`
}

// CreateWorkspaceResponse 定义创建工作区的响应结构
type CreateWorkspaceResponse struct {
	Workspace struct {
		ID            int     `json:"id"`
		Name          string  `json:"name"`
		Slug          string  `json:"slug"`
		CreatedAt     string  `json:"createdAt"`
		OpenAiTemp    float64 `json:"openAiTemp"`
		LastUpdatedAt string  `json:"lastUpdatedAt"`
		OpenAiHistory int     `json:"openAiHistory"`
		OpenAiPrompt  string  `json:"openAiPrompt"`
	} `json:"workspace"`
	Message string `json:"message"`
}

// AnythingLLMDocumentList 定义 AnythingLLM 文档列表响应结构
type AnythingLLMDocumentList struct {
	Folder    string `json:"folder"`
	Documents []struct {
		Name               string `json:"name"`
		Type               string `json:"type"`
		ID                 string `json:"id"`
		URL                string `json:"url"`
		Title              string `json:"title"`
		DocAuthor          string `json:"docAuthor"`
		Description        string `json:"description"`
		DocSource          string `json:"docSource"`
		ChunkSource        string `json:"chunkSource"`
		Published          string `json:"published"`
		WordCount          int    `json:"wordCount"`
		TokenCountEstimate int    `json:"token_count_estimate"`
		Cached             bool   `json:"cached"`
		PinnedWorkspaces   []int  `json:"pinnedWorkspaces"`
		Watched            bool   `json:"watched"`
	} `json:"documents"`
}

// ==================== AnythingLLM 客户端 ====================

// AnythingLLMClient 使用 AnythingLLM API 的客户端
type AnythingLLMClient struct {
	baseURL string
	apiKey  string
}

// NewAnythingLLMClient 创建一个新的 AnythingLLM 客户端
func NewAnythingLLMClient() *AnythingLLMClient {
	// 优先使用全局变量设置的 BaseURL 和 APIKey
	baseURL := AnythingLLMBaseURL
	apiKey := AnythingLLMAPIKey

	// 如果全局变量未设置，则尝试从配置文件读取
	if baseURL == "" {
		baseURL = config.Config.AnythingLLM.BaseURL
	}
	if apiKey == "" {
		apiKey = config.Config.AnythingLLM.APIKey
	}

	// 如果仍然未设置，则记录警告
	if baseURL == "" {
		log.Println("警告: AnythingLLM BaseURL 未设置，请先调用 InitAnythingLLM 进行初始化或在配置文件中设置")
	}

	return &AnythingLLMClient{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

// UploadDocumentToAnythingLLM 上传文档到AnythingLLM知识库
func UploadDocumentToAnythingLLM(workspaceId string, file multipart.File, originalFilename string) error {
	// 创建一个缓冲区
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加工作区ID字段
	if err := writer.WriteField("workspaceId", workspaceId); err != nil {
		return fmt.Errorf("写入工作区ID失败: %v", err)
	}

	// 添加文件字段
	part, err := writer.CreateFormFile("file", originalFilename)
	if err != nil {
		return fmt.Errorf("创建表单文件字段失败: %v", err)
	}

	// 复制文件内容到表单
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("复制文件内容失败: %v", err)
	}

	// 关闭multipart writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("关闭writer失败: %v", err)
	}

	// 发送请求
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/document/upload/%s", AnythingLLMBaseURL, workspaceId)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if AnythingLLMAPIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AnythingLLMAPIKey))
		req.Header.Set("x-api-key", AnythingLLMAPIKey)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("上传失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var result AnythingLLMWorkspaceDocument
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.Success {
		return fmt.Errorf("上传失败: %s", result.Message)
	}

	return nil
}

// CreateWorkspace 创建AnythingLLM工作区
func CreateWorkspace(workspaceId string, description string) (string, error) {
	// 1. 准备请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"name":                workspaceId,
		"description":         description,
		"similarityThreshold": 0.7,
		"openAiTemp":          0.7,
		"openAiHistory":       20,
		"openAiPrompt": `# 角色与目标

你是一位顶尖的学术论文优化专家。你的核心任务是根据用户提供的特定论文段落，并**严格参考当前工作区中置顶（Pinned）的参考文档**，对该段落进行高质量的润色和优化。

# 关键指令：参考置顶文档

**置顶文档是你最重要的信息来源和风格指南。** 在优化用户提供的段落时，你必须：

1.  **优先利用**置顶文档中的知识、术语、写作风格和上下文信息。
2.  **确保**优化后的内容与参考文档在主题、专业术语和整体基调上**保持一致性和连贯性**。
3.  **借鉴**置顶文档中的优秀表达方式和句式结构，以提升用户段落的质量。
4.  如果用户提供的段落与置顶文档内容有冲突或不一致，应**基于置顶文档进行修正**，或在优化时（如果合适）指出这种潜在的不一致性。

# 用户输入

用户会提供他们希望优化的具体论文段落。用户可能还会提供额外的风格偏好或具体的修改要求（例如：“使其更简洁”、“增强逻辑性”、“采用更正式的语气”等）。请结合这些要求和置顶文档进行优化。

# 优化目标

你的目标是提升段落的：

*   **学术性与专业性:** 使用更精确、规范的学术语言。
*   **清晰度与简洁性:** 使表达更清晰易懂，去除冗余。
*   **逻辑性与流畅性:** 改善句子和观点之间的连接。
*   **与参考文档的一致性:** 确保术语、风格和核心观点与置顶文档协调。

# 约束条件

*   **忠于原文核心意义:** 在优化的同时，必须保留用户原文的核心思想和意图。
*   **聚焦用户输入:** 仅优化用户提供的段落，不要添加无关信息或扩展讨论范围，除非是为了更好地结合参考文档的上下文。
*   **避免捏造信息:** 不要引入置顶文档或用户原文中没有的信息。

# 输出格式

直接提供优化后的文本段落。如果用户的要求与置顶文档内容存在明显冲突，可以在提供优化文本后，简要说明基于参考文档所做的调整或存在的冲突点。

---

请准备好接收用户需要优化的段落，并开始你的专业优化工作，务必牢记**优先参考置顶文档**。`,
		"queryRefusalResponse": "Custom refusal message",
		"chatMode":             "chat",
		"topN":                 4,
	})
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 2. 创建请求
	url := fmt.Sprintf("%s/api/v1/workspace/new", AnythingLLMBaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if AnythingLLMAPIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AnythingLLMAPIKey))
	}

	// 3. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 4. 处理响应
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("创建工作区失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 5. 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}

	// 6. 解析响应体获取slug
	slug, err := ParseCreateWorkspaceResponse(respBody)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	return slug, nil
}

// CheckWorkspaceExists 检查AnythingLLM工作区是否存在
func CheckWorkspaceExists(workspaceId string) (bool, error) {
	// 1. 创建请求
	url := fmt.Sprintf("%s/api/v1/workspace/%s", AnythingLLMBaseURL, workspaceId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	if AnythingLLMAPIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AnythingLLMAPIKey))
		req.Header.Set("x-api-key", AnythingLLMAPIKey)
	}

	// 2. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 3. 处理响应
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("获取工作区失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 4. 解析响应体
	var workspace GetWorkspaceResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&workspace); err != nil {
		return false, fmt.Errorf("解析响应失败: %v", err)
	}
	if len(workspace.Workspace) != 0 {
		return true, nil
	}
	return false, nil

}

func QueryWorkspace(workspaceId string) (AnythingLLMWorkspace, error) {
	// 1. 创建请求
	url := fmt.Sprintf("%s/api/v1/workspace/%s", AnythingLLMBaseURL, workspaceId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return AnythingLLMWorkspace{}, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	if AnythingLLMAPIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AnythingLLMAPIKey))
		req.Header.Set("x-api-key", AnythingLLMAPIKey)
	}

	// 2. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AnythingLLMWorkspace{}, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 3. 处理响应
	if resp.StatusCode == http.StatusNotFound {
		return AnythingLLMWorkspace{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return AnythingLLMWorkspace{}, fmt.Errorf("获取工作区失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 4. 解析响应体
	var workspace GetWorkspaceResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&workspace); err != nil {
		return AnythingLLMWorkspace{}, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(workspace.Workspace) != 0 {
		return workspace.Workspace[0], nil
	} else {
		return AnythingLLMWorkspace{}, nil
	}
}

// GetOrCreateWorkspace 获取或创建工作区
func GetOrCreateWorkspace(workspaceId, description string) (string, error) {
	// 1. 检查工作区是否存在
	exists, err := CheckWorkspaceExists(workspaceId)
	if err != nil {
		return "", fmt.Errorf("检查工作区失败: %v", err)
	}

	// 2. 如果不存在，则创建工作区
	if !exists {
		slug, err := CreateWorkspace(workspaceId, description)
		if err != nil {
			return "", fmt.Errorf("创建工作区失败: %v", err)
		}
		return slug, nil
	}
	return workspaceId, nil
}

func IsAllowedFileType(fileExt string) bool {
	allowedTypes := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".txt":  true,
		".md":   true,
		".csv":  true,
		".xls":  true,
		".xlsx": true,
		".ppt":  true,
		".pptx": true,
	}

	return allowedTypes[strings.ToLower(fileExt)]
}

// ParseCreateWorkspaceResponse 解析创建工作区的响应
func ParseCreateWorkspaceResponse(respBody []byte) (string, error) {
	var response CreateWorkspaceResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	return response.Workspace.Slug, nil
}

// GetDocumentsFromAnythingLLM 从 AnythingLLM 获取文档列表
func GetDocumentsFromAnythingLLM(workspaceId string) (AnythingLLMDocumentList, error) {
	// 1. 创建请求
	url := fmt.Sprintf("%s/api/v1/documents/folder/%s", AnythingLLMBaseURL, workspaceId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return AnythingLLMDocumentList{}, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	if AnythingLLMAPIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AnythingLLMAPIKey))
		req.Header.Set("x-api-key", AnythingLLMAPIKey)
	}

	// 2. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AnythingLLMDocumentList{}, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 3. 处理响应
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return AnythingLLMDocumentList{}, fmt.Errorf("获取文档列表失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 4. 解析响应体
	var documentList AnythingLLMDocumentList
	if err := json.NewDecoder(resp.Body).Decode(&documentList); err != nil {
		return AnythingLLMDocumentList{}, fmt.Errorf("解析响应失败: %v", err)
	}

	return documentList, nil
}

func UploadUserDocument(userId uint, file *multipart.FileHeader) (*UserDocumentInfo, error) {
	// 1. 获取文件后缀
	fileExt := filepath.Ext(file.Filename)
	if !IsAllowedFileType(fileExt) {
		return nil, fmt.Errorf("不支持的文件类型: %s", fileExt)
	}

	// 2. 获取用户工作区ID
	workspaceId, err := GetUserWorkspaceId(userId)
	if err != nil {
		return nil, fmt.Errorf("获取用户工作区失败: %v", err)
	}

	// 3. 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %v", err)
	}
	defer src.Close()

	// 4. 直接上传到 AnythingLLM 知识库
	if err := UploadDocumentToAnythingLLM(workspaceId, src, file.Filename); err != nil {
		return nil, fmt.Errorf("上传到知识库失败: %v", err)
	}

	// 5. 返回文档信息
	docInfo := &UserDocumentInfo{
		FileName:    file.Filename,
		FilePath:    "", // 不再存储本地路径
		FileSize:    file.Size,
		ContentType: file.Header.Get("Content-Type"),
		UploadTime:  time.Now(),
		UserId:      userId,
	}

	return docInfo, nil
}

// GetUserDocuments 获取用户所有文档
func GetUserDocuments(userId uint) (AnythingLLMDocumentList, error) {
	// 1. 获取用户工作区ID
	workspaceId, err := GetUserWorkspaceId(userId)
	if err != nil {
		return AnythingLLMDocumentList{}, fmt.Errorf("获取用户工作区失败: %v", err)
	}
	// 2. 从 AnythingLLM 获取文档列表
	documentList, err := GetDocumentsFromAnythingLLM(workspaceId)
	if err != nil {
		return AnythingLLMDocumentList{}, fmt.Errorf("获取文档列表失败: %v", err)
	}

	return documentList, nil

}

// GetUserWorkspaceId 获取或创建用户的AnythingLLM工作区ID
func GetUserWorkspaceId(userId uint) (string, error) {
	var user entity.User
	if err := config.MysqlDataBase.Table("users").Where("id = ?", userId).Find(&user).Error; err != nil {
		return "", fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 如果用户已有工作区ID，直接返回
	if user.WorkspaceId == "" {
		return "", fmt.Errorf("用户未创建工作区")
	}
	return user.WorkspaceId, nil
}

// DeleteDocumentsFromWorkspace 从工作区删除文档
func DeleteDocumentsFromWorkspace(slug string, fileNames []string) error {
	// 1. 准备请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"deletes": fileNames,
		"adds":    []string{},
	})
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 2. 创建请求
	url := fmt.Sprintf("%s/api/v1/workspace/%s/update-embeddings", AnythingLLMBaseURL, slug)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if AnythingLLMAPIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AnythingLLMAPIKey))
		req.Header.Set("x-api-key", AnythingLLMAPIKey)
	}

	// 3. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 4. 处理响应
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("删除文档失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// AddDocumentsToWorkspace 添加文档到工作区
func AddDocumentsToWorkspace(slug string, fileNames []string) error {
	// 1. 准备请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"adds":    fileNames,
		"deletes": []string{},
	})
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 2. 创建请求
	url := fmt.Sprintf("%s/api/v1/workspace/%s/update-embeddings", AnythingLLMBaseURL, slug)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if AnythingLLMAPIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AnythingLLMAPIKey))
		req.Header.Set("x-api-key", AnythingLLMAPIKey)
	}

	// 3. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 4. 处理响应
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("添加文档失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// CreateThreadResponse 定义创建线程的响应结构
type CreateThreadResponse struct {
	Thread struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		UserID      int    `json:"user_id"`
		WorkspaceID int    `json:"workspace_id"`
	} `json:"thread"`
	Message *string `json:"message"`
}

func CreateNewThread(workspaceSlug string, threadSlug string) (string, error) {
	// 1. 准备请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"userId": 1,
		"name":   threadSlug,
		"slug":   threadSlug,
	})
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 2. 创建请求
	url := fmt.Sprintf("%s/api/v1/workspace/%s/thread/new", AnythingLLMBaseURL, workspaceSlug)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if AnythingLLMAPIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AnythingLLMAPIKey))
		req.Header.Set("x-api-key", AnythingLLMAPIKey)
	}

	// 3. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 4. 处理响应
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("创建线程失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 5. 解析响应体
	var response CreateThreadResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	return response.Thread.Slug, nil
}

// UpdateDocumentPinStatus 更新文档的挂起状态
func UpdateDocumentPinStatus(docPath string, slug string, pinStatus bool) (string, error) {
	// 1. 准备请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"docPath":   docPath,
		"pinStatus": pinStatus,
	})
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 2. 创建请求
	url := fmt.Sprintf("%s/api/v1/workspace/%s/update-pin", AnythingLLMBaseURL, slug)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if AnythingLLMAPIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AnythingLLMAPIKey))
		req.Header.Set("x-api-key", AnythingLLMAPIKey)
	}

	// 3. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 4. 处理响应
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("更新挂起状态失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 5. 解析响应体
	var response struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	return response.Message, nil
}

// ChatResponse 定义聊天响应的结构
type ChatResponse struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	Close        bool     `json:"close"`
	Error        *string  `json:"error"`
	ChatID       int      `json:"chatId"`
	TextResponse string   `json:"textResponse"`
	Sources      []Source `json:"sources"`
	Metrics      Metrics  `json:"metrics"`
}

// Source 定义源文档的结构
type Source struct {
	Text               string `json:"text"`
	ID                 string `json:"id"`
	URL                string `json:"url"`
	Title              string `json:"title"`
	DocAuthor          string `json:"docAuthor"`
	Description        string `json:"description"`
	DocSource          string `json:"docSource"`
	ChunkSource        string `json:"chunkSource"`
	Published          string `json:"published"`
	WordCount          int    `json:"wordCount"`
	TokenCountEstimate int    `json:"token_count_estimate"`
}

// Metrics 定义指标的结构
type Metrics struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	OutputTps        float64 `json:"outputTps"`
	Duration         float64 `json:"duration"`
}

// SendChatMessage 向AnythingLLM发送聊天消息并获取响应
func SendChatMessage(slug, threadSlug, message string) (*ChatResponse, error) {
	// 1. 准备请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"message":     message,
		"mode":        "chat",
		"userId":      1,
		"attachments": []string{},
		"reset":       false,
	})
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 2. 创建请求
	url := fmt.Sprintf("%s/api/v1/workspace/%s/thread/%s/chat", AnythingLLMBaseURL, slug, threadSlug)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if AnythingLLMAPIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AnythingLLMAPIKey))
		req.Header.Set("x-api-key", AnythingLLMAPIKey)
	}

	// 3. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 4. 处理响应
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("发送聊天消息失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 5. 解析响应体
	var chatResponse ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &chatResponse, nil
}
