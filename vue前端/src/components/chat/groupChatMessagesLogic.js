import { ref, onMounted, watch, onBeforeUnmount } from 'vue'
import { useUserStore } from '@/stores/user'
import { get, post, postJSON } from '@/utility/request'
import { message } from 'ant-design-vue'
import { BACKEND_DOMAIN } from '@/utility/config'
import { parseDateTime } from '@/utility/common'

/**
 * 群组聊天消息逻辑
 * @param {Object} props - 组件属性
 * @param {Function} emit - 事件发射器
 * @returns {Object} 包含所有群组聊天消息所需的状态和方法
 */
export function useGroupChatMessages(props, emit) {
  const userStore = useUserStore()
  const messages = ref([])
  const loading = ref(false)
  const messageInput = ref('')
  const mediaType = ref('text')
  const mediaFile = ref(null)
  const mediaPreview = ref('')
  const isEmojiPickerVisible = ref(false)
  const ws = ref(null) // WebSocket连接
  const retryCount = ref(0)
  const isClosing = ref(false)
  const isConnected = ref(false)
  const heartbeatInterval = ref(null) // 心跳定时器

  // 处理WebSocket消息
  const handleWebSocketMessage = (data) => {
    // 处理不同类型的消息
    switch (data.type) {
      case 'chat':
        // 添加新消息到消息列表，需要转换字段名
        messages.value.push({
          ...data,
          SenderID: data.senderId,
          SenderName: data.nickname,
          SenderAvatar: data.avatarUrl,
          Content: data.content,
          MediaType: data.mediaType,
          MediaURL: data.mediaUrl,
          CreatedAt: data.timestamp,
          type: data.type
        })
        // 滚动到底部
        setTimeout(scrollToBottom, 100)
        break
        
      case 'join':
      case 'leave':
        // 处理系统消息（加入/离开）
        messages.value.push({
          ...data,
          MediaType: 'text',
          SenderID: data.senderId || 0,
          Content: data.content,
          CreatedAt: data.timestamp,
          type: data.type // 保留原始类型以便于模板中区分
        })
        setTimeout(scrollToBottom, 100)
        break
        
      default:
        console.warn('未知的消息类型:', data.type)
    }
  }

  // 发送心跳消息，保持连接活跃
  const startHeartbeat = () => {
    // 清除现有的心跳定时器
    stopHeartbeat()
    
    // 每30秒发送一次心跳
    heartbeatInterval.value = setInterval(() => {
      if (ws.value && ws.value.readyState === WebSocket.OPEN) {
        try {
          // 发送一个简单的ping消息
          ws.value.send(JSON.stringify({ type: 'ping' }))
          console.log('发送心跳消息')
        } catch (error) {
          console.error('发送心跳消息失败:', error)
          connectWebSocket() // 尝试重连
        }
      } else if (ws.value && ws.value.readyState !== WebSocket.CONNECTING) {
        connectWebSocket() // 如果连接已断开且不在连接中，尝试重连
      }
    }, 30000)
  }

  // 停止心跳
  const stopHeartbeat = () => {
    if (heartbeatInterval.value) {
      clearInterval(heartbeatInterval.value)
      heartbeatInterval.value = null
    }
  }

  // 建立WebSocket连接
  const connectWebSocket = () => {
    if (!props.groupId) return
    
    // 如果已经连接或正在连接，不要重复连接
    if (ws.value && (ws.value.readyState === WebSocket.OPEN || ws.value.readyState === WebSocket.CONNECTING)) {
      console.log('WebSocket已连接或正在连接中')
      return
    }
    
    try {
      // 将 http:// 替换为 ws://，并移除末尾的斜杠
      const baseUrl = BACKEND_DOMAIN.replace(/^http/, 'ws').replace(/\/$/, '')
      const wsUrl = `${baseUrl}/ws/groupChat/${props.groupId}`
      
      console.log('正在连接WebSocket:', wsUrl)
      
      // 获取认证token
      const token = localStorage.getItem('authToken')
      if (!token) {
        message.error('未登录或登录已过期')
        return
      }

      // 如果已经有连接，先关闭
      if (ws.value) {
        stopHeartbeat()
        isClosing.value = true
        ws.value.close()
        isClosing.value = false
      }

      ws.value = new WebSocket(wsUrl)
      
      // 设置连接超时
      const connectionTimeout = setTimeout(() => {
        if (ws.value && ws.value.readyState === WebSocket.CONNECTING) {
          ws.value.close()
          message.error('WebSocket连接超时，请检查网络状态')
        }
      }, 5000)
      
      ws.value.onopen = () => {
        clearTimeout(connectionTimeout)
        console.log('WebSocket连接已建立，发送认证消息')
        
        // 发送认证消息
        const authMessage = {
          token: token
        }
        ws.value.send(JSON.stringify(authMessage))
      }
      
      ws.value.onclose = (event) => {
        clearTimeout(connectionTimeout)
        console.log('WebSocket连接已关闭:', event)
        isConnected.value = false
        stopHeartbeat()
        
        // 如果不是主动关闭，尝试重连
        if (!isClosing.value) {
          const delay = Math.min(1000 * Math.pow(2, retryCount.value), 30000)
          retryCount.value++
          
          setTimeout(() => {
            if (!isClosing.value) {
              console.log(`尝试第 ${retryCount.value} 次重连...`)
              connectWebSocket()
            }
          }, delay)
        }
      }
      
      ws.value.onerror = (error) => {
        clearTimeout(connectionTimeout)
        console.error('WebSocket错误:', error)
        isConnected.value = false
        message.error('WebSocket连接错误，正在尝试重连...')
      }
      
      ws.value.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          
          // 处理认证成功消息
          if (data.code === 200 && (data.message === '认证成功' || data.data === '认证成功' || data.data === '连接成功')) {
            console.log('WebSocket认证成功')
            isConnected.value = true
            retryCount.value = 0
            startHeartbeat() // 启动心跳
            return
          }
          
          // 处理pong响应，不做特殊处理
          if (data.type === 'pong') {
            console.log('收到pong响应')
            return
          }
          
          // 处理错误消息
          if (data.code !== undefined && data.code !== 200) {
            message.error(data.message || '连接错误')
            if (data.code === 401) {
              // 认证失败，重新连接
              ws.value.close()
            }
            return
          }
          
          // 处理正常的聊天消息
          if (data.type) {
            handleWebSocketMessage(data)
          }
        } catch (error) {
          console.error('解析WebSocket消息失败:', error)
        }
      }
    } catch (error) {
      console.error('创建WebSocket连接失败:', error)
      message.error('创建WebSocket连接失败，请稍后重试')
    }
  }

  // 关闭WebSocket连接
  const closeWebSocket = () => {
    stopHeartbeat()
    if (ws.value) {
      isClosing.value = true
      if (ws.value.readyState === WebSocket.OPEN || ws.value.readyState === WebSocket.CONNECTING) {
        ws.value.close()
      }
      ws.value = null
      isConnected.value = false
    }
  }

  // 获取历史消息列表
  const fetchMessages = () => {
    if (!props.groupId) return
    
    loading.value = true
    get(`/api/chat/groups/${props.groupId}/messages`, {}, (_, data) => {
      // 确保每条消息都有type属性，用于区分系统消息和普通消息
      messages.value = data.map(msg => {
        // 假设消息的内容中包含"加入群聊"或"离开群聊"的是系统消息
        if (msg.Content && (msg.Content.includes('加入群聊') || msg.Content.includes('离开群聊'))) {
          return { ...msg, type: msg.Content.includes('加入群聊') ? 'join' : 'leave' }
        }
        return { ...msg, type: 'chat' }
      })
      loading.value = false
      scrollToBottom()
    }, (msg) => {
      loading.value = false
      message.error('获取消息失败' + msg)
    })
  }

  // 发送消息 - 使用WebSocket
  const sendMessage = async () => {
    if (!messageInput.value.trim() && !mediaFile.value) return
    
    // 如果连接正在建立中，等待一段时间
    if (ws.value && ws.value.readyState === WebSocket.CONNECTING) {
      try {
        await new Promise((resolve, reject) => {
          const timeout = setTimeout(() => {
            reject(new Error('连接超时'))
          }, 5000)

          const checkConnection = () => {
            if (isConnected.value) {
              clearTimeout(timeout)
              resolve()
            } else if (ws.value.readyState === WebSocket.CLOSED) {
              clearTimeout(timeout)
              reject(new Error('连接已关闭'))
            } else {
              setTimeout(checkConnection, 100)
            }
          }
          checkConnection()
        })
      } catch (error) {
        message.error('聊天连接失败，正在重新连接...')
        connectWebSocket()
        return
      }
    }

    // 检查连接状态
    if (!isConnected.value || !ws.value || ws.value.readyState !== WebSocket.OPEN) {
      message.error('聊天连接已断开，正在重新连接...')
      connectWebSocket()
      return
    }
    
    // 如果有媒体文件，先上传
    if (mediaFile.value) {
      uploadMedia(() => {
        sendTextMessage()
      })
      return
    }
    
    sendTextMessage()
  }

  // 发送文本消息
  const sendTextMessage = () => {
    if (!isConnected.value) {
      message.error('等待连接建立...')
      return
    }

    const messageData = {
      type: 'chat',
      content: messageInput.value,
      mediaType: mediaType.value,
      mediaURL: mediaPreview.value
    }
    
    try {
      ws.value.send(JSON.stringify(messageData))
      messageInput.value = ''
      mediaFile.value = null
      mediaPreview.value = ''
      mediaType.value = 'text'
    } catch (error) {
      console.error('发送消息失败:', error)
      message.error('发送消息失败，正在重新连接...')
      connectWebSocket()
    }
  }

  // 上传媒体文件
  const uploadMedia = (callback) => {
    if (!mediaFile.value) {
      callback()
      return
    }
    
    const formData = new FormData()
    formData.append('media', mediaFile.value)
    formData.append('type', mediaType.value)
    
    post('/api/upload/chat-media', formData, (_, data) => {
      mediaPreview.value = data.url
      callback()
    }, () => {
      message.error('媒体上传失败')
    })
  }

  // 处理媒体文件选择
  const handleMediaSelect = (event) => {
    const file = event.target.files[0]
    if (!file) return
    
    mediaFile.value = file
    mediaType.value = file.type.startsWith('image/') ? 'image' : 
                      file.type.startsWith('video/') ? 'video' : 
                      file.type.startsWith('audio/') ? 'audio' : 'file'
    
    if (mediaType.value === 'image') {
      const reader = new FileReader()
      reader.onload = (e) => {
        mediaPreview.value = e.target.result
      }
      reader.readAsDataURL(file)
    }
  }

  // 取消媒体文件选择
  const cancelMediaSelect = () => {
    mediaFile.value = null
    mediaPreview.value = ''
    mediaType.value = 'text'
  }

  // 滚动到底部
  const scrollToBottom = () => {
    const container = document.querySelector('.messages-container')
    if (container) {
      container.scrollTop = container.scrollHeight
    }
  }

  // 监听groupId变化
  watch(() => props.groupId, (newVal, oldVal) => {
    if (newVal !== oldVal) {
      closeWebSocket() // 关闭旧连接
      messages.value = [] // 清空消息列表
      fetchMessages() // 获取新群组的历史消息
      connectWebSocket() // 建立新连接
    }
  })

  // 组件挂载时初始化
  onMounted(() => {
    fetchMessages() // 首先获取历史消息
    connectWebSocket() // 然后建立WebSocket连接
  })

  // 组件卸载时清理
  onBeforeUnmount(() => {
    stopHeartbeat()
    if (ws.value) {
      isClosing.value = true
      try {
        ws.value.close(1000, 'Component unmounted')
      } catch (error) {
        console.error('关闭WebSocket连接失败:', error)
      }
      ws.value = null
    }
  })

  return {
    // 状态
    userStore,
    messages,
    loading,
    messageInput,
    mediaType,
    mediaFile,
    mediaPreview,
    isEmojiPickerVisible,
    isConnected,
    BACKEND_DOMAIN,

    // 方法
    handleWebSocketMessage,
    sendMessage,
    handleMediaSelect,
    cancelMediaSelect,
    scrollToBottom,
    parseDateTime
  }
} 