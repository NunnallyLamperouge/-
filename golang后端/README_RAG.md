# RAG 服务使用指南

本文档提供了 RAG (Retrieval-Augmented Generation) 服务的使用指南，该服务允许用户上传个人文档并基于这些文档进行知识检索和问答。

## 功能概述

1. **文档管理**：
   - 用户可以上传个人文档（支持PDF、Word、TXT等多种格式）
   - 查看已上传的文档列表
   - 删除不需要的文档

2. **知识库问答**：
   - 基于用户上传的文档进行问答
   - 支持普通问答和流式问答两种模式
   - 提供引用来源，便于追溯信息

## API 接口说明

### 1. 文档管理

#### 上传文档

```
POST /api/rag/documents
```

**请求参数**：
- `file`: 文档文件 (multipart/form-data)

**响应示例**：
```json
{
  "message": "文件上传成功",
  "data": {
    "fileName": "论文.pdf",
    "filePath": "user_documents/12345/20230523112233_abcdef12.pdf",
    "fileSize": 1024567,
    "contentType": "application/pdf",
    "uploadTime": "2023-05-23T11:22:33Z",
    "userId": 12345
  }
}
```

#### 获取文档列表

```
GET /api/rag/documents
```

**响应示例**：
```json
{
  "message": "获取文档列表成功",
  "data": [
    {
      "fileName": "论文.pdf",
      "filePath": "user_documents/12345/20230523112233_abcdef12.pdf",
      "fileSize": 1024567,
      "contentType": "application/pdf",
      "uploadTime": "2023-05-23T11:22:33Z",
      "userId": 12345
    }
  ]
}
```

#### 删除文档

```
DELETE /api/rag/documents/:fileName
```

**响应示例**：
```json
{
  "message": "文档删除成功"
}
```

### 2. 知识库问答

#### 普通问答

```
POST /api/rag/chat
```

**请求参数**：
```json
{
  "message": "请总结一下我上传的论文中的主要观点",
  "chatId": "chat_123456",
  "temperature": 0.7
}
```

**响应示例**：
```json
{
  "response": "根据您上传的论文，主要观点包括以下几点：1. ...",
  "chatId": "chat_123456",
  "sources": [
    {
      "text": "...",
      "source": "论文.pdf",
      "page": 12
    }
  ],
  "tokenUsage": {
    "prompt": 245,
    "completion": 128,
    "total": 373
  }
}
```

#### 流式问答

```
GET /api/rag/ws/chat
```

通过 WebSocket 连接进行流式问答，首先需要发送初始消息：

```json
{
  "message": "请总结一下我上传的论文中的主要观点",
  "chatId": "chat_123456",
  "temperature": 0.7
}
```

然后接收流式响应，每个响应包含：

```json
{
  "content": "根据您上传的论文，",
  "done": false
}
```

最后一条消息的 `done` 字段为 `true`。

## 支持的文件类型

- PDF (`.pdf`)
- Microsoft Word (`.doc`, `.docx`)
- 文本文件 (`.txt`, `.md`)
- Excel 表格 (`.xls`, `.xlsx`)
- PowerPoint (`.ppt`, `.pptx`)
- CSV 文件 (`.csv`)

## 技术实现

本服务基于以下技术实现：
1. Go 语言后端服务
2. AnythingLLM 作为知识检索和向量数据库管理
3. WebSocket 用于流式响应

## 使用限制

1. 单个文件大小限制：20MB
2. 用户文档总存储空间限制：根据用户权限级别
3. 请求频率限制：单用户每分钟10次查询

## 最佳实践

1. 上传高质量文档以获得更准确的检索结果
2. 提出具体明确的问题以获得更精确的回答
3. 使用流式响应模式可以获得更快的首字响应速度 