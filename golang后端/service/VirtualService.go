package service

import (
	"context"
	"encoding/json"
	"example/config"
	"example/entity"
	"example/utility"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var tts *utility.AzureTTS
var stt *utility.AzureSTT

// Initialize after app.yml has been read, before serving requests.
func InitSpeechServices() {
	tts = utility.NewAzureTTS(config.Config.AzureTTS.SubscriptionKey, config.Config.AzureTTS.Region)
	stt = utility.NewAzureSTT(config.Config.AzureSpeech.SubscriptionKey, config.Config.AzureSpeech.Region)
}

type GetVirtualQuestionRequest struct {
	TemplateId    int    `json:"template_id"`
	QuestionCount int    `json:"question_count"`
	Difficulty    string `json:"difficulty"`
}

func GetVirtualQuestion(c *gin.Context) {
	var request GetVirtualQuestionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(200, entity.ErrorResponse[string](401, "请求参数错误"))
		return
	}
	var template entity.FakeTemplate
	config.MysqlDataBase.Where("id = ?", request.TemplateId).First(&template)
	if template.ID == 0 {
		c.JSON(200, entity.ErrorResponse[string](401, "模板不存在"))
		return
	}
	content := utility.ExtractQuotedChinese(template.Content)

	prompt := "你是一个专业的论文答辩老师，现在需要根据论文内容生成答辩问题，要求如下："
	prompt += "1. 问题数量为" + strconv.Itoa(request.QuestionCount) + "个。"
	prompt += "2. 问题难度为" + request.Difficulty + "。"
	prompt += "3. 你最后需要返回一个json格式，格式为一个数组，数组中每个元素为一个字符串，字符串为问题内容。"
	question := "问题内容需要有关于" + content + "中的内容。"
	var message = []utility.Message{}
	chatRequest := utility.ChatRequest{
		Model:       config.Config.OpenAI.UseModelName,
		Messages:    message,
		Prompt:      prompt,
		Question:    question,
		Temperature: float32(config.Config.OpenAI.GlobalTemperature),
		MaxTokens:   8000,
	}

	// 最大重试次数
	maxRetries := 3
	var questions []string
	var answer string

	for retry := 0; retry < maxRetries; retry++ {
		resp, err := utility.ChatHandler(chatRequest)
		if err != nil {
			c.JSON(200, entity.ErrorResponse[string](500, "生成问题失败"))
			return
		}

		answer = cleanJSONResponse(resp.Choices[0].Message.Content)
		// 尝试解析JSON
		err = json.Unmarshal([]byte(answer), &questions)
		if err == nil && len(questions) > 0 {
			// 成功解析，跳出循环
			break
		}

		// 如果解析失败并且达到了最大重试次数，返回错误
		if retry == maxRetries-1 {
			c.JSON(200, entity.ErrorResponse[string](500, "无法解析AI生成的问题，请重试"))
			return
		}
	}

	c.JSON(200, entity.SuccessResponse(questions))
}
func cleanJSONResponse(response string) string {
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	return response
}

type GetVirtualQuestionAudioRequest struct {
	Question string `json:"question"`
}

func GetVirtualQuestionAudio(c *gin.Context) {
	var request GetVirtualQuestionAudioRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(200, entity.ErrorResponse[string](401, "请求参数错误"))
		return
	}
	audioFileName, err := tts.TextToSpeech(request.Question, "zh-CN", "zh-CN-XiaoYanNeural", "Female")
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "生成音频失败"+err.Error()))
		return
	}
	c.JSON(200, entity.SuccessResponse(audioFileName))
}

type UserAnswer struct {
	QuestionID      int    `json:"question_Id"`
	QuestionContent string `json:"questionContent"`
	Answer          string `json:"answer"`
	AudioBase64     string `json:"audioBase64"`
}

type SubmitVirtualDefenseRequest struct {
	TemplateId     uint         `json:"template_id"`
	QuestionCount  uint         `json:"question_count"`
	Difficulty     string       `json:"difficulty"`
	Questions      []string     `json:"questions"`
	AudioFileNames []string     `json:"audioFileNames"`
	Answers        []UserAnswer `json:"answers"`
	Seconds        int          `json:"seconds"`
}

func SubmitVirtualDefense(c *gin.Context) {
	userId, _ := c.Get("userId")
	userIdInt := utility.GetUserIDFromContext(userId)
	var request SubmitVirtualDefenseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(200, entity.ErrorResponse[string](401, "请求参数错误"))
		return
	}

	virtual := entity.Virtual{
		TemplateId: request.TemplateId,
		StudentId:  uint(userIdInt),
		Difficulty: request.Difficulty,
		Seconds:    request.Seconds,
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Create(&virtual).Error; err != nil {
		tx.Rollback()
		c.JSON(200, entity.ErrorResponse[string](500, "创建虚拟答辩失败"))
		return
	}

	for index, answer := range request.Answers {
		randomString := func() string {
			const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
			b := make([]byte, 6)
			for i := range b {
				b[i] = charset[rand.Intn(len(charset))]
			}
			return string(b)
		}
		timestamp := time.Now().Format("20060102")
		audioFileName := fmt.Sprintf("%s_%s.wav", timestamp, randomString())
		audioPath := "./audio"

		filePath, err := utility.SaveBase64ToFile(answer.AudioBase64, audioPath, audioFileName)
		if err != nil {
			tx.Rollback()
			c.JSON(200, entity.ErrorResponse[string](500, "保存音频文件失败:"+err.Error()))
			return
		}
		answerText, err := stt.SpeechToText(filePath, "zh-CN")
		if err != nil {
			tx.Rollback()
			c.JSON(200, entity.ErrorResponse[string](500, "语音识别失败:"+err.Error()))
			return
		}

		question := entity.VirtualQuestion{
			VirtualId:         virtual.ID,
			Question:          answer.QuestionContent,
			Answer:            answerText,
			QuestionAudioPath: request.AudioFileNames[index],
			AnswerAudioPath:   audioFileName,
		}

		// 将问题保存到数据库
		if err := tx.Create(&question).Error; err != nil {
			tx.Rollback()
			c.JSON(200, entity.ErrorResponse[string](500, "保存答辩问题失败"))
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(200, entity.ErrorResponse[string](500, "提交事务失败"))
		return
	}

	c.JSON(200, entity.SuccessResponse(virtual.ID))
}

type VirtualGetResult struct {
	Virtual          entity.Virtual
	VirtualQuestions []entity.VirtualQuestion
}

func GetVirtualDefense(c *gin.Context) {
	virtualId := c.Query("virtual_id")
	var virtual entity.Virtual
	config.MysqlDataBase.Where("id = ?", virtualId).Preload("FakeTemplate").Preload("User").First(&virtual)
	if virtual.ID == 0 {
		c.JSON(200, entity.ErrorResponse[string](401, "虚拟答辩不存在"))
		return
	}
	var questions []entity.VirtualQuestion
	config.MysqlDataBase.Where("virtual_id = ?", virtualId).Find(&questions)
	c.JSON(200, entity.SuccessResponse(VirtualGetResult{
		Virtual:          virtual,
		VirtualQuestions: questions,
	}))
}

func GetVirtualList(c *gin.Context) {
	userId, _ := c.Get("userId")
	userIdInt := utility.GetUserIDFromContext(userId)
	var virtuals []entity.Virtual
	config.MysqlDataBase.Where("student_id = ?", userIdInt).Preload("FakeTemplate").Find(&virtuals)
	c.JSON(200, entity.SuccessResponse(virtuals))
}

type StreamChatCompletionRequest struct {
	Content   VirtualGetResult `json:"content"`
	AuthToken string           `json:"auth_token"`
}

func StreamChatCompletion(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, entity.ErrorResponse[string](500, "WebSocket升级失败"))
		return
	}
	defer ws.Close()

	// 读取并解析初始消息
	_, message, err := ws.ReadMessage()
	if err != nil {
		return
	}
	var request StreamChatCompletionRequest
	if err := json.Unmarshal(message, &request); err != nil {
		ws.WriteJSON(entity.ErrorResponse[string](500, "无法解析请求参数"+err.Error()))
		return
	}

	// 验证token并获取userId
	claims, err := utility.ParseToken(request.AuthToken)
	if err != nil {
		ws.WriteJSON(entity.ErrorResponse[string](500, "token验证失败"))
		return
	}
	userId := claims.UserID
	if userId == 0 {
		ws.WriteJSON(entity.ErrorResponse[string](500, "token验证失败"))
		return
	}
	content := utility.ExtractQuotedChinese(request.Content.Virtual.FakeTemplate.Content)
	prompt := "你是一名毕业论文答辩教研室负责人，现在有学生给你他的答辩记录，你需要根据他的论文内容，对他的答辩表现进行全方面的评价和优化建议。论文标题:" + request.Content.Virtual.FakeTemplate.Title + "论文内容:" + content
	questionStr := "请根据学生的答辩记录，对他的答辩表现进行全方面的评价和优化建议。"
	dbjl := ""
	for _, question := range request.Content.VirtualQuestions {
		dbjl += "问题:" + question.Question + "\n"
		dbjl += "学生回答:" + question.Answer + "\n"
	}
	questionStr += dbjl
	// 调用流式聊天
	streamChan, err := utility.StreamChatCompletion(context.Background(), utility.ChatRequest{
		Model:       config.Config.OpenAI.ThinkModelName,
		Messages:    []utility.Message{},
		Prompt:      prompt,
		Question:    questionStr,
		Temperature: float32(config.Config.OpenAI.GlobalTemperature),
		MaxTokens:   8192,
	})

	if err != nil {
		ws.WriteJSON(entity.ErrorResponse[string](500, "启动流式生成失败"+err.Error()))
		return
	}

	// 读取流式响应并通过WebSocket发送
	for response := range streamChan {
		if err := ws.WriteJSON(response); err != nil {
			return
		}
		if response.Done {
			break
		}
	}

}

type GetScoreRequest struct {
	Result VirtualGetResult `json:"result"`
}

func GetScore(c *gin.Context) {
	var request GetScoreRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(200, entity.ErrorResponse[string](401, "请求参数错误"))
		return
	}
	prompt := "你是一名毕业论文答辩教研室负责人，现在有学生给你他的答辩记录，你需要根据他的论文内容，对他的答辩表现进行打分，满分100分，最后你只需要返回一个整型数字即可。论文标题:" + request.Result.Virtual.FakeTemplate.Title + "论文内容:" + request.Result.Virtual.FakeTemplate.Content + "这是你之前对他的评价:" + request.Result.Virtual.Description
	questionStr := "请根据学生的答辩记录，对他的答辩表现进行打分，满分100分，最后你只需要返回一个整型数字即可。"
	dbjl := ""
	for _, question := range request.Result.VirtualQuestions {
		dbjl += "问题:" + question.Question + "\n"
		dbjl += "学生回答:" + question.Answer + "\n"
	}
	questionStr += dbjl

	// 最大重试次数
	maxRetries := 3
	var scoreStr string
	var scoreInt int

	for retry := 0; retry < maxRetries; retry++ {
		score, err := utility.ChatHandler(utility.ChatRequest{
			Model:    config.Config.OpenAI.ThinkModelName,
			Messages: []utility.Message{},
			Prompt:   prompt,
			Question: questionStr,
		})
		if err != nil {
			c.JSON(200, entity.ErrorResponse[string](500, "生成分数失败"))
			return
		}

		// 获取生成的分数
		scoreStr = strings.TrimSpace(score.Choices[0].Message.Content)

		// 验证是否为整型数字
		var parseErr error
		scoreInt, parseErr = strconv.Atoi(scoreStr)
		if parseErr == nil {
			// 成功获取整型数字，跳出循环
			break
		}

		// 如果达到最大重试次数仍未成功，返回错误
		if retry == maxRetries-1 {
			c.JSON(200, entity.ErrorResponse[string](500, "无法获取有效的整型分数"))
			return
		}

		// 重试前稍微等待一下
		time.Sleep(500 * time.Millisecond)
	}

	request.Result.Virtual.Score = scoreInt
	if err := config.MysqlDataBase.Save(&request.Result.Virtual).Error; err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "保存分数失败"))
		return
	}
	c.JSON(200, entity.SuccessResponse(request.Result.Virtual))
}
