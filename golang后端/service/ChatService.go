package service

import (
	"context"
	"encoding/json"
	"example/config"
	"example/entity"
	"example/utility"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该更严格
	},
}

// WebSocketMessage 定义WebSocket消息结构
type WebSocketMessage struct {
	Type      string      `json:"type"`      // 消息类型: "chat", "system", "join", "leave"
	Content   string      `json:"content"`   // 消息内容
	SenderID  uint        `json:"senderId"`  // 发送者ID
	GroupID   uint        `json:"groupId"`   // 群组ID
	MediaType string      `json:"mediaType"` // 媒体类型，如 "text", "image" 等
	MediaURL  string      `json:"mediaUrl"`  // 媒体URL
	Timestamp time.Time   `json:"timestamp"` // 发送时间
	Nickname  string      `json:"nickname"`  // 发送者昵称
	AvatarURL string      `json:"avatarUrl"` // 发送者头像
	Data      interface{} `json:"data"`      // 额外数据
}

// ClientManager 管理所有WebSocket客户端
type ClientManager struct {
	clients    map[*Client]bool          // 所有连接的客户端
	broadcast  chan []byte               // 广播消息的通道
	register   chan *Client              // 注册客户端的通道
	unregister chan *Client              // 注销客户端的通道
	groups     map[uint]map[*Client]bool // 群组到客户端的映射
	mutex      sync.Mutex                // 互斥锁，保护并发访问
}

// Client 表示一个WebSocket客户端
type Client struct {
	id           uint              // 用户ID
	socket       *websocket.Conn   // WebSocket连接
	send         chan []byte       // 发送消息的通道
	groups       map[uint]struct{} // 用户所在的群组
	isAuthorized bool              // 是否已认证
}

// 全局的客户端管理器
var Manager = ClientManager{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	groups:     make(map[uint]map[*Client]bool),
}

// Start 启动WebSocket管理器
func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.register:
			// 注册新客户端
			manager.mutex.Lock()
			manager.clients[client] = true
			manager.mutex.Unlock()

		case client := <-manager.unregister:
			// 客户端断开连接
			if _, ok := manager.clients[client]; ok {
				manager.mutex.Lock()
				delete(manager.clients, client)
				// 从所有群组中移除该客户端
				for groupID := range client.groups {
					if _, exists := manager.groups[groupID]; exists {
						delete(manager.groups[groupID], client)
						// 如果群组没有成员了，删除该群组
						if len(manager.groups[groupID]) == 0 {
							delete(manager.groups, groupID)
						}
					}
				}
				manager.mutex.Unlock()
				close(client.send)
			}
		}
	}
}

// 添加初始化消息结构体
type GroupChatInitMessage struct {
	Token string `json:"token"`
}

// HandleGroupChat 处理群组聊天的WebSocket连接
func HandleGroupChat(c *gin.Context) {
	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// fmt.Printf("WebSocket升级失败: %v\n", err)
		return
	}

	// 获取群组ID
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		// fmt.Printf("无效的群组ID: %v\n", err)
		conn.WriteJSON(entity.ErrorResponse[string](400, "无效的群组ID"))
		conn.Close()
		return
	}

	// fmt.Printf("新的WebSocket连接请求: 群组ID=%d\n", groupID)

	// 创建未认证的客户端
	client := &Client{
		socket:       conn,
		send:         make(chan []byte, 256), // 增加缓冲区大小
		groups:       make(map[uint]struct{}),
		isAuthorized: false,
	}

	// 等待认证消息
	_, message, err := conn.ReadMessage()
	if err != nil {
		// fmt.Printf("读取认证消息失败: %v\n", err)
		conn.Close()
		return
	}

	// fmt.Printf("收到认证消息: %s\n", string(message))

	var initMsg GroupChatInitMessage
	if err := json.Unmarshal(message, &initMsg); err != nil {
		// fmt.Printf("解析认证消息失败: %v\n", err)
		conn.WriteJSON(entity.ErrorResponse[string](500, "无法解析认证消息"))
		conn.Close()
		return
	}

	// 验证token并获取用户ID
	claims, err := utility.ParseToken(initMsg.Token)
	if err != nil {
		// fmt.Printf("Token验证失败: %v\n", err)
		conn.WriteJSON(entity.ErrorResponse[string](500, "token验证失败"))
		conn.Close()
		return
	}
	userID := claims.UserID

	// fmt.Printf("用户 %d 认证成功\n", userID)

	// 验证用户是否为群组成员
	var member entity.GroupMember
	err = config.MysqlDataBase.Where("group_id = ? AND user_id = ? AND status != 3", groupID, userID).First(&member).Error
	if err != nil {
		// fmt.Printf("群组成员验证失败: %v\n", err)
		conn.WriteJSON(entity.ErrorResponse[string](403, "您不是该群组成员或群组不存在"))
		conn.Close()
		return
	}

	// 设置客户端信息
	client.id = uint(userID)
	client.isAuthorized = true
	client.groups[uint(groupID)] = struct{}{}

	// 向管理器注册客户端
	Manager.register <- client

	// 向群组添加客户端
	Manager.mutex.Lock()
	if _, exists := Manager.groups[uint(groupID)]; !exists {
		Manager.groups[uint(groupID)] = make(map[*Client]bool)
	}
	Manager.groups[uint(groupID)][client] = true
	Manager.mutex.Unlock()

	// fmt.Printf("用户 %d 已加入群组 %d\n", userID, groupID)

	// 发送认证成功消息
	err = conn.WriteJSON(entity.SuccessResponse("认证成功"))
	if err != nil {
		// fmt.Printf("发送认证成功消息失败: %v\n", err)
		conn.Close()
		return
	}

	// 查询用户信息
	var user entity.User
	config.MysqlDataBase.Select("id, username, avatar").Where("id = ?", userID).First(&user)

	// 发送用户加入的系统消息
	joinMessage := WebSocketMessage{
		Type:      "join",
		SenderID:  uint(userID),
		GroupID:   uint(groupID),
		Timestamp: time.Now(),
		Nickname:  user.Username,
		AvatarURL: user.Avatar,
		Content:   fmt.Sprintf("%s 上线了", user.Username),
	}

	// 广播加入消息
	broadcastToGroup(joinMessage, uint(groupID))

	// 创建一个 WaitGroup 来等待 goroutines 完成
	var wg sync.WaitGroup
	wg.Add(2)

	// 启动读取 goroutine
	go func() {
		defer wg.Done()
		client.readPump(uint(groupID))
	}()

	// 启动写入 goroutine
	go func() {
		defer wg.Done()
		client.writePump()
	}()

	// 等待 goroutines 完成
	wg.Wait()

	// fmt.Printf("用户 %d 的WebSocket连接已关闭\n", userID)
}

// readPump 从WebSocket读取消息
func (c *Client) readPump(groupID uint) {
	defer func() {
		// fmt.Printf("用户 %d 的readPump正在退出\n", c.id)
		Manager.unregister <- c
		c.socket.Close()
	}()

	c.socket.SetReadLimit(512000)
	c.socket.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.socket.SetPongHandler(func(string) error {
		// fmt.Printf("用户 %d 收到pong消息\n", c.id)
		c.socket.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// fmt.Printf("用户 %d 的readPump开始运行\n", c.id)

	for {
		messageType, message, err := c.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("用户 %d WebSocket读取错误: %v\n", c.id, err)
			} else {
				fmt.Printf("用户 %d WebSocket连接关闭: %v\n", c.id, err)
			}
			break
		}

		fmt.Printf("用户 %d 收到消息类型: %d, 内容: %s\n", c.id, messageType, string(message))

		// 检查认证状态
		if !c.isAuthorized {
			// fmt.Printf("用户 %d 未认证\n", c.id)
			c.socket.WriteJSON(entity.ErrorResponse[string](401, "未认证的连接"))
			continue
		}

		// 解析消息
		var wsMessage WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			// fmt.Printf("用户 %d 消息解析失败: %v\n", c.id, err)
			continue
		}

		// fmt.Printf("用户 %d 消息解析成功: 类型=%s, 内容=%s\n", c.id, wsMessage.Type, wsMessage.Content)

		// 处理ping消息
		if wsMessage.Type == "ping" {
			// fmt.Printf("用户 %d 收到ping消息\n", c.id)
			// 回复pong消息
			pongMessage := WebSocketMessage{
				Type:      "pong",
				Timestamp: time.Now(),
			}

			if data, err := json.Marshal(pongMessage); err == nil {
				select {
				case c.send <- data:
					// fmt.Printf("用户 %d pong消息已发送\n", c.id)
				default:
					// fmt.Printf("用户 %d pong消息发送失败: 通道已满\n", c.id)
				}
			}

			// 更新读取截止时间
			c.socket.SetReadDeadline(time.Now().Add(60 * time.Second))
			continue
		}

		// 设置消息属性
		wsMessage.SenderID = c.id
		wsMessage.GroupID = groupID
		wsMessage.Timestamp = time.Now()

		// 查询发送者信息
		var user entity.User
		if err := config.MysqlDataBase.Select("username, avatar").Where("id = ?", c.id).First(&user).Error; err == nil {
			wsMessage.Nickname = user.Username
			wsMessage.AvatarURL = user.Avatar
			// fmt.Printf("用户 %d 发送者信息: 用户名=%s, 头像=%s\n", c.id, user.Username, user.Avatar)
		} else {
			// fmt.Printf("用户 %d 获取用户信息失败: %v\n", c.id, err)
		}

		// 保存消息到数据库
		saveMessage(wsMessage)

		// fmt.Printf("用户 %d 准备广播消息到群组 %d\n", c.id, groupID)
		// 广播消息到群组
		broadcastToGroup(wsMessage, groupID)
	}
}

// writePump 向WebSocket写入消息
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// 通道已关闭
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 添加队列中的消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// 发送ping保持连接
			c.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// broadcastToGroup 向群组广播消息
func broadcastToGroup(message WebSocketMessage, groupID uint) {
	data, err := json.Marshal(message)
	if err != nil {
		// fmt.Printf("消息序列化失败: %v\n", err)
		return
	}

	Manager.mutex.Lock()
	defer Manager.mutex.Unlock()

	if clients, exists := Manager.groups[groupID]; exists {
		for client := range clients {
			if client.isAuthorized {
				select {
				case client.send <- data:
					// fmt.Printf("消息已发送到客户端 %d\n", client.id)
				default:
					// fmt.Printf("客户端 %d 的消息队列已满，正在关闭连接\n", client.id)
					close(client.send)
					delete(Manager.groups[groupID], client)
					delete(Manager.clients, client)
				}
			}
		}
	} else {
		// fmt.Printf("群组 %d 不存在或没有活跃客户端\n", groupID)
	}
}

// saveMessage 将消息保存到数据库
func saveMessage(message WebSocketMessage) {
	// 创建消息记录
	dbMessage := entity.Message{
		Type:      entity.GroupMessage, // 群聊消息类型
		Content:   message.Content,
		SenderID:  message.SenderID,
		GroupID:   message.GroupID,
		IsRead:    false,
		MediaType: message.MediaType,
		MediaURL:  message.MediaURL,
	}

	// 保存到数据库
	config.MysqlDataBase.Create(&dbMessage)
}

// HandleStreamChat 处理流式聊天的WebSocket连接
func HandleStreamChat(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	// 读取初始请求
	_, message, err := ws.ReadMessage()
	if err != nil {
		return
	}

	var chatRequest utility.ChatRequest
	if err := json.Unmarshal(message, &chatRequest); err != nil {
		ws.WriteJSON(map[string]interface{}{
			"error": "请求JSON有误" + err.Error(),
		})
		return
	}

	// 创建上下文，支持取消
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 获取流式响应
	streamChan, err := utility.StreamChatCompletion(ctx, chatRequest)
	if err != nil {
		ws.WriteJSON(map[string]interface{}{
			"error": "流失请求失败" + err.Error(),
		})
		return
	}

	// 发送流式响应
	for response := range streamChan {
		if err := ws.WriteJSON(response); err != nil {
			return
		}
	}
}

// 初始化函数，启动WebSocket管理器
func init() {
	go Manager.Start()
}

type GetAdveriseRequest struct {
	// FakeTemplate entity.FakeTemplate `json:"fakeTemplate"`
	Source string `json:"source"`
}

func GetAdverise(c *gin.Context) {
	var request GetAdveriseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(200, entity.ErrorResponse[string](401, "请求参数错误"))
		return
	}
	//prompt := "你是一名毕业论文段落优化师，现在有学生给你他的论文段落部分，你需要根据他的论文内容，对这部分内容进行润色，最后你只需要返回一个润色后的段落即可。论文标题:" + request.FakeTemplate.Title + "论文整体内容:" + request.FakeTemplate.Content
	prompt := "你是一名毕业论文段落优化师，现在有学生给你他的论文段落部分，你需要根据他的论文内容，对这部分内容进行润色，最后你只需要返回一个润色后的段落即可。"

	questionStr := "这是他论文的段落部分:" + request.Source

	// 最大重试次数
	var scoreStr string
	score, err := utility.ChatHandler(utility.ChatRequest{
		Model:    config.Config.OpenAI.UseModelName,
		Messages: []utility.Message{},
		Prompt:   prompt,
		Question: questionStr,
	})
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "段落润色失败"))
		return
	}

	// 获取生成的分数
	scoreStr = score.Choices[0].Message.Content
	c.JSON(200, entity.SuccessResponse(scoreStr))
}

type GetAdveriseUniversalRequest struct {
	Source     string `json:"source"`             // 需要修改的原文内容
	ModifyDes  string `json:"modify_description"` // 修改要求描述
	FieldType  string `json:"field_type"`         // 论文字段类型（如：title, acknowledgment, content等）
	Style      string `json:"style,omitempty"`    // 可选的风格要求
	TemplateId int    `json:"templateId,omitempty"`
}

func GetAdveriseUniversal(c *gin.Context) {
	var request GetAdveriseUniversalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(200, entity.ErrorResponse[string](401, "请求参数错误"))
		return
	}

	// 根据字段类型构建不同的提示词
	basePrompt := "你是一名专业的学术论文优化专家。"
	fieldPrompt := ""
	switch request.FieldType {
	case "title":
		fieldPrompt = "你需要帮助优化论文标题，使其更加专业、准确和吸引人。"
	case "acknowledgment":
		fieldPrompt = "你需要帮助优化论文致谢部分，使其更加真诚、得体和专业。"
	case "abstract":
		fieldPrompt = "你需要帮助优化论文摘要，确保其简洁、完整地概括论文的主要内容。"
	case "content":
		fieldPrompt = "你需要帮助优化论文正文部分，提升其学术性和专业性。"
	default:
		fieldPrompt = "你需要帮助优化论文内容，提升其整体质量。"
	}

	stylePrompt := ""
	if request.Style != "" {
		stylePrompt = fmt.Sprintf("请按照以下风格要求进行优化：%s", request.Style)
	}

	modifyPrompt := ""
	if request.ModifyDes != "" {
		modifyPrompt = fmt.Sprintf("具体优化要求：%s", request.ModifyDes)
	}

	sourcePrompt := ""
	if request.TemplateId != 0 {
		var source string
		if err := config.MysqlDataBase.Table("fake_templates").Select("content").Where("id = ?", request.TemplateId).Find(&source).Error; err != nil {
			log.Println(err.Error())
		} else {
			sourcePrompt = fmt.Sprintf("以下是原文内容，仅供参考以保持文章的连贯性和避免产生不实信息：%s", source)
		}
	}

	prompt := strings.Join([]string{basePrompt, fieldPrompt, stylePrompt, modifyPrompt, sourcePrompt}, " ")
	questionStr := fmt.Sprintf("这是需要优化的%s部分：%s", request.FieldType, request.Source)

	score, err := utility.ChatHandler(utility.ChatRequest{
		Model:    config.Config.OpenAI.UseModelName,
		Messages: []utility.Message{},
		Prompt:   prompt,
		Question: questionStr,
	})
	if err != nil {
		c.JSON(200, entity.ErrorResponse[string](500, "内容优化失败"))
		return
	}

	resStr := score.Choices[0].Message.Content
	c.JSON(200, entity.SuccessResponse(resStr))
}
