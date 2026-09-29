# AnythingLLM 适配包使用指南

本适配包提供了对 AnythingLLM API 的易用封装，支持普通聊天和流式聊天响应。

## 配置方式

### 1. 通过配置文件（推荐）

在 `app.yml` 文件中添加以下配置：

```yaml
anythingLLM:
  baseURL: "http://localhost:3001"  # AnythingLLM 的基础URL
  apiKey: "your-api-key"            # AnythingLLM 的API密钥（如果需要）
```

### 2. 通过全局变量

在应用启动时，调用初始化函数：

```go
import "example/utility"

func init() {
    // 初始化 AnythingLLM 配置
    utility.InitAnythingLLM("http://localhost:3001", "your-api-key")
}
```

## 基本使用

### 普通聊天请求

```go
import (
    "example/utility"
    "fmt"
)

func main() {
    // 聊天参数
    workspaceId := "your-workspace-id" // 工作区ID，在AnythingLLM中创建的工作区
    message := "你好，请帮我总结一下这个文档的内容"
    chatId := ""                     // 如果是新对话，可以为空；如果是已有对话，则提供对话ID
    model := ""                      // 可选，如果为空则使用AnythingLLM默认模型
    temperature := float32(0.7)      // 可选，生成文本的随机性，值越大生成结果越多样

    // 发送消息并获取响应
    response, err := utility.AnythingLLMChatHandler(workspaceId, message, chatId, model, temperature)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    // 打印响应内容
    fmt.Printf("响应: %s\n", response.Response)
    fmt.Printf("会话ID: %s\n", response.ChatId)

    // 如果有令牌用量信息，则打印
    if response.TokenUsage.Total > 0 {
        fmt.Printf("提示令牌: %d\n", response.TokenUsage.Prompt)
        fmt.Printf("补全令牌: %d\n", response.TokenUsage.Completion)
        fmt.Printf("总令牌: %d\n", response.TokenUsage.Total)
    }

    // 打印引用源（如果有）
    if len(response.Sources) > 0 {
        fmt.Println("\n引用源:")
        for i, source := range response.Sources {
            fmt.Printf("  %d. %v\n", i+1, source)
        }
    }
}
```

### 流式聊天请求

```go
import (
    "context"
    "example/utility"
    "fmt"
    "time"
)

func main() {
    // 创建上下文，以便可以取消请求
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // 聊天参数
    workspaceId := "your-workspace-id" // 工作区ID
    message := "你好，请帮我总结一下这个文档的内容"
    chatId := ""                     // 如果是新对话，可以为空
    model := ""                      // 可选，如果为空则使用默认模型
    temperature := float32(0.7)      // 可选，生成文本的随机性

    // 开始流式聊天
    stream, err := utility.AnythingLLMStreamChatHandler(ctx, workspaceId, message, chatId, model, temperature)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    // 收集完整的响应文本
    var fullResponse string

    // 处理流式响应
    for response := range stream {
        if response.Done {
            break
        }
        // 打印或处理流式响应片段
        fmt.Print(response.Content)
        fullResponse += response.Content
    }

    fmt.Println("\n流式聊天完成")
    fmt.Printf("完整响应: %s\n", fullResponse)
}
```

### 使用客户端直接调用

如果需要更多控制，可以直接使用 AnythingLLMClient：

```go
import (
    "example/utility"
    "fmt"
)

func main() {
    // 创建客户端
    client := utility.NewAnythingLLMClient()
    
    // 聊天参数
    workspaceId := "your-workspace-id"
    message := "你好，请帮我总结一下这个文档的内容"
    chatId := ""
    model := ""
    temperature := float32(0.7)
    
    // 发送消息
    response, err := client.SendMessage(workspaceId, message, chatId, model, temperature)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    
    // 处理响应
    fmt.Printf("响应: %s\n", response.Response)
}
```

## API参数说明

| 参数 | 类型 | 说明 |
|------|------|------|
| workspaceId | string | 必填，AnythingLLM中创建的工作区ID |
| message | string | 必填，要发送的消息内容 |
| chatId | string | 可选，会话ID，用于继续现有会话 |
| model | string | 可选，要使用的模型名称 |
| temperature | float32 | 可选，生成文本的随机性，0.0-1.0 |

## 响应结构

### 普通聊天响应

```go
type AnythingLLMChatResponse struct {
    Response   string `json:"response"`   // 响应文本
    ChatId     string `json:"chatId"`     // 会话ID
    Sources    []any  `json:"sources"`    // 引用的知识库源
    TokenUsage struct {
        Prompt     int `json:"prompt"`     // 提示令牌数
        Completion int `json:"completion"` // 补全令牌数
        Total      int `json:"total"`      // 总令牌数
    } `json:"tokenUsage"`
}
```

### 流式聊天响应

```go
type AnythingLLMStreamResponse struct {
    Content string `json:"content"` // 响应文本片段
    Done    bool   `json:"done"`    // 是否完成
}
```

## 注意事项

1. 使用前请确保 AnythingLLM 服务已正确部署和运行
2. 检查 API 密钥是否正确设置（如果 AnythingLLM 启用了密钥认证）
3. workspaceId 必须是 AnythingLLM 中已创建的工作区
4. 对于流式响应，请确保正确处理上下文取消，避免资源泄漏 