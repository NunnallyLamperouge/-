# 研途智伴源码

这是“研途智伴”源码。
后端为 Go + Gin + GORM，前端为 Vue 3 + Vite。
源码已移除原服务器配置、用户文档和真实凭据；`app.example.yml` 仅含占位值。

## 运行前提

- Go 1.23.4 工具链（`golang后端/go.mod` 声明）
- Node.js/npm
- MySQL 8、Redis
- 可选：AnythingLLM、兼容 OpenAI Chat Completions 的模型接口、Azure Speech TTS/STT

## 后端

```bash
cd golang后端
cp app.example.yml app.yml
# 编辑 app.yml，填入本地数据库、JWT、模型、知识库和语音服务配置
go mod download
go run .
```

后端默认监听 `:8080`。启动时会连接 MySQL 和 Redis，并自动创建 GORM 模型表；也可以在专用空数据库先执行上级目录的 `schema.sql`（含 DROP TABLE，勿对已有业务库执行）。注册需要有效邮件验证码；源码包不附带可直接登录的历史账号。

## 前端

```bash
cd vue前端
npm ci
npm run dev
```

前端使用 `src/utility/config.js` 和 `src/main.js` 中的本地 API 地址。部署到其他地址时两处同步调整，并配置后端 CORS。

## 主要接口

- `/api/auth/*`：注册、登录和验证码
- `/api/template/*`：模板与历史版本
- `/api/mission/*`：教师任务与学生提交
- `/api/guide/*`：师生关系
- `/api/rag/*`：知识库文档、工作区和 AI 润色
- `/api/virtual/*`：AI 答辩问题、TTS、录音提交与结果
- `/api/chat/*`、`/api/ws/*`：师生群组和实时消息

外部模型、AnythingLLM、邮箱和 Azure Speech 未配置时，相应功能不能返回真实结果。
