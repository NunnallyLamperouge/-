<script setup>
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'
import { useGroupChatMessages } from './groupChatMessagesLogic'

const props = defineProps({
  groupId: {
    type: Number,
    required: true
  }
})

const emits = defineEmits(['previewImage'])

// 导入并使用封装的群组聊天消息逻辑
const {
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
} = useGroupChatMessages(props, emits)
</script>

<template>
  <div class="chat-messages-container h-full flex flex-col animate__animated animate__fadeIn">
    <!-- 消息列表 -->
    <div class="messages-container flex-1 overflow-y-auto p-6 space-y-4 bg-white dark:bg-zinc-900 scrollbar-thin scrollbar-thumb-gray-300 dark:scrollbar-thumb-gray-600 scrollbar-track-transparent">
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin h-8 w-8 border-2 border-blue-500 dark:border-blue-400 rounded-full border-t-transparent"></div>
      </div>
      
      <div v-else-if="messages.length === 0" class="flex flex-col items-center justify-center py-12 animate__animated animate__fadeIn animate__delay-1s">
        <div class="text-6xl text-gray-300 dark:text-gray-600 mb-4 animate__animated animate__fadeIn animate__delay-2s">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-20 w-20" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
        </div>
        <div class="border border-gray-200 dark:border-gray-700 p-6 w-fit max-w-md rounded-xl bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/30 dark:to-indigo-900/30 animate__animated animate__fadeIn animate__delay-3s">
          <p class="text-lg text-gray-800 dark:text-gray-200 font-medium">开始聊天吧</p>
          <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">发送消息开始与群组成员交流</p>
        </div>
      </div>
      
      <template v-else>
        <!-- 系统消息（加入/离开提示） -->
        <div 
          v-for="message in messages" 
          :key="message.ID"
          class="animate__animated animate__fadeInUp animate__faster"
        >
          <!-- 系统消息（加入/离开提示） -->
          <div v-if="message.type === 'join' || message.type === 'leave'" 
               class="my-8 flex justify-center animate__animated animate__fadeIn">
            <div class="px-5 py-2 rounded-full bg-gradient-to-r from-gray-100 to-gray-50 dark:from-zinc-800/80 dark:to-zinc-800/40 text-xs text-gray-500 dark:text-gray-400 inline-flex items-center gap-2 shadow-sm border border-gray-200 dark:border-gray-700 backdrop-blur-sm">
              <div v-if="message.type === 'join'" class="w-4 h-4 rounded-full bg-blue-100 dark:bg-blue-900/40 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 text-blue-500 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                </svg>
              </div>
              <div v-else class="w-4 h-4 rounded-full bg-orange-100 dark:bg-orange-900/40 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 text-orange-500 dark:text-orange-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 12h-15" />
                </svg>
              </div>
              <span class="font-medium">{{ message.Content }}</span>
              <span class="text-[10px] text-gray-400 dark:text-gray-500 ml-1 opacity-70">{{ parseDateTime(message.CreatedAt) }}</span>
            </div>
          </div>
          
          <!-- 普通聊天消息 -->
          <div 
            v-else
            class="flex items-start gap-3 group mb-8 transition-all duration-300"
            :class="message.SenderID === userStore.user.id ? 'flex-row-reverse' : ''"
          >
            <!-- 发送者头像 -->
            <div class="w-10 h-10 rounded-full overflow-hidden bg-gradient-to-br from-blue-100 to-indigo-100 dark:from-blue-900/30 dark:to-indigo-900/30 flex-shrink-0 border border-gray-200 dark:border-gray-700 shadow-sm transition-all duration-300 group-hover:scale-110 group-hover:shadow-md group-hover:shadow-blue-500/10">
              <img 
                v-if="message.senderAvatar" 
                :src="BACKEND_DOMAIN + message.senderAvatar" 
                class="w-full h-full object-cover"
                alt="用户头像"
                onerror="this.style.display='none'"
              >
              <div v-else class="w-full h-full flex items-center justify-center text-blue-500 dark:text-blue-400 text-sm font-bold">
                {{ message.senderName?.charAt(0) || '?' }}
              </div>
              <div>{{ message.senderName }}</div>
            </div>
            
            <!-- 消息内容 -->
            <div 
              class="max-w-[70%] group-hover:shadow-md transition-all duration-300"
              :class="[
                message.SenderID === userStore.user.id 
                  ? 'bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-2xl rounded-tr-sm hover:-translate-y-0.5' 
                  : 'bg-gray-50 dark:bg-zinc-800/50 dark:text-gray-200 border border-gray-200 dark:border-gray-700 hover:-translate-y-0.5 rounded-2xl rounded-tl-sm'
              ]"
            >
              <!-- 发送者名称 -->
              <div 
                v-if="message.SenderID !== userStore.user.id"
                class="px-3.5 pt-2.5 text-xs font-medium flex items-center gap-1.5"
                :class="message.SenderID === userStore.user.id 
                  ? 'text-blue-200' 
                  : 'text-blue-500 dark:text-blue-400'"
              >
                {{ message.senderName }}
                <div class="w-1 h-1 rounded-full bg-current opacity-50"></div>
                <span class="text-[10px] text-gray-400 dark:text-gray-500 opacity-70">{{ parseDateTime(message.CreatedAt).split(' ')[1] }}</span>
              </div>
            
              <!-- 文本消息 -->
              <div class="p-3.5" :class="{'pt-2': message.SenderID !== userStore.user.id}">
                <p v-if="message.MediaType === 'text'" class="break-words text-sm leading-relaxed">
                  {{ message.Content }}
                </p>
              
                <!-- 图片消息 -->
                <div v-else-if="message.MediaType === 'image'" class="relative group/image overflow-hidden rounded-lg">
                  <img 
                    :src="message.MediaURL"
                    class="max-w-full rounded-lg shadow-sm group-hover/image:shadow-lg transition-all duration-500 cursor-pointer transform group-hover/image:scale-[1.02]"
                    alt="图片消息"
                    @click="$emit('previewImage', message.MediaURL)"
                  >
                  <div class="absolute inset-0 bg-black/0 group-hover/image:bg-black/10 rounded-lg transition-all duration-300 flex items-center justify-center opacity-0 group-hover/image:opacity-100">
                    <div class="bg-black/50 p-1.5 rounded-full transform scale-90 group-hover/image:scale-100 transition-transform">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                      </svg>
                    </div>
                  </div>
                </div>
              
                <!-- 音频消息 -->
                <div v-else-if="message.MediaType === 'audio'" class="bg-white/20 dark:bg-black/20 rounded-lg p-2 shadow-sm backdrop-blur-sm">
                  <div class="flex items-center gap-2 mb-1 text-xs opacity-70">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
                    </svg>
                    <span>语音消息</span>
                  </div>
                  <audio 
                    :src="message.MediaURL"
                    controls
                    class="max-w-full rounded-md w-full"
                  ></audio>
                </div>
              
                <!-- 视频消息 -->
                <div v-else-if="message.MediaType === 'video'" class="rounded-lg overflow-hidden shadow-sm">
                  <div class="bg-white/20 dark:bg-black/20 p-2 flex items-center gap-2 text-xs opacity-70 backdrop-blur-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                    <span>视频消息</span>
                  </div>
                  <video 
                    :src="message.MediaURL"
                    controls
                    class="max-w-full rounded-b-lg"
                  ></video>
                </div>
              </div>
              
              <!-- 发送时间 (仅自己的消息) -->
              <div 
                v-if="message.SenderID === userStore.user.id"
                class="text-[10px] px-3.5 pb-2.5 opacity-70 text-right"
                :class="message.SenderID === userStore.user.id 
                  ? 'text-blue-200' 
                  : 'text-gray-400 dark:text-gray-500'"
              >
                {{ parseDateTime(message.CreatedAt).split(' ')[1] }}
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>

    <!-- 输入区域 -->
    <div class="p-4 border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-zinc-900 shadow-lg shadow-black/5">
      <!-- 连接状态指示器 -->
      <div v-if="!isConnected" class="mb-2 flex items-center justify-center">
        <span class="text-xs px-3 py-1 rounded-full bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-400 flex items-center gap-1.5 shadow-sm animate-pulse">
          <span class="h-1.5 w-1.5 rounded-full bg-orange-500"></span>
          正在连接...
        </span>
      </div>

      <!-- 媒体预览 -->
      <div v-if="mediaPreview" class="mb-4 relative animate__animated animate__fadeIn">
        <div class="rounded-xl overflow-hidden border border-gray-200 dark:border-gray-700 shadow-sm hover:shadow-md transition-all duration-300">
          <img 
            v-if="mediaType === 'image'"
            :src="mediaPreview"
            class="max-h-40 w-full object-contain bg-gray-50 dark:bg-zinc-800/50"
            alt="媒体预览"
          >
          <div v-else class="p-4 bg-gray-50 dark:bg-zinc-800/50 flex items-center gap-3">
            <div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-full">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <span class="text-gray-700 dark:text-gray-300 font-medium">
              {{ mediaFile?.name }}
            </span>
          </div>
        </div>
        <button 
          @click="cancelMediaSelect"
          class="absolute -top-2 -right-2 w-7 h-7 bg-red-500 text-white rounded-full flex items-center justify-center hover:bg-red-600 transition-colors shadow-lg hover:shadow-red-500/30 transition-all duration-300 transform hover:scale-110"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      
      <div class="flex items-center gap-3">
        <!-- 左侧按钮区 -->
        <div class="flex items-center space-x-1">
          <!-- 媒体上传按钮 -->
          <div class="relative">
            <input 
              type="file" 
              accept="image/*,audio/*,video/*"
              class="hidden"
              @change="handleMediaSelect"
              ref="fileInput"
            >
            <button 
              @click="$refs.fileInput.click()"
              class="p-2.5 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 hover:shadow-sm transform hover:scale-105"
              title="上传媒体文件"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v6m3-3H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </button>
          </div>

          <!-- 表情按钮 -->
          <button 
            class="p-2.5 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 hover:shadow-sm transform hover:scale-105"
            title="选择表情"
            @click="isEmojiPickerVisible = !isEmojiPickerVisible"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </button>
        </div>
        
        <!-- 消息输入框 -->
        <input 
          v-model="messageInput"
          @keyup.enter="sendMessage"
          type="text"
          placeholder="输入消息..."
          class="flex-1 px-5 py-3 bg-gray-50 dark:bg-zinc-800/50 border border-gray-200 dark:border-gray-700 rounded-full focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 dark:text-gray-300 dark:placeholder-gray-500 transition-all shadow-sm hover:shadow-md text-sm"
        >
        
        <!-- 发送按钮 -->
        <button 
          @click="sendMessage"
          class="px-4 py-3 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-full hover:opacity-90 hover:shadow-md hover:shadow-blue-500/20 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 transform hover:scale-105"
          :disabled="!messageInput.trim() && !mediaFile"
        >
          <span class="hidden sm:inline font-medium">发送</span>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 自定义滚动条样式 */
.scrollbar-thin::-webkit-scrollbar {
  width: 6px;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background: transparent;
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 3px;
}

.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background-color: rgba(156, 163, 175, 0.7);
}

/* 动画相关样式 */
.animate__animated {
  animation-duration: 0.4s;
}

.animate__faster {
  animation-duration: 0.25s;
}

.animate__delay-1s {
  animation-delay: 0.1s;
}

.animate__delay-2s {
  animation-delay: 0.2s;
}

.animate__delay-3s {
  animation-delay: 0.3s;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.animate__fadeInUp {
  animation-name: fadeInUp;
}

.animate__fadeIn {
  animation-name: fadeIn;
}

/* 音频播放器样式 */
audio::-webkit-media-controls-timeline {
  @apply bg-gray-200 dark:bg-zinc-700 rounded-full overflow-hidden;
}

audio::-webkit-media-controls-current-time-display,
audio::-webkit-media-controls-time-remaining-display {
  @apply text-gray-700 dark:text-gray-300;
}

audio::-webkit-media-controls-volume-slider {
  @apply bg-gray-200 dark:bg-zinc-700 rounded-full;
}

/* Firefox 样式 */
audio::-moz-range-track {
  @apply bg-gray-200 dark:bg-zinc-700 rounded-full;
}

audio::-moz-range-thumb {
  @apply bg-blue-500;
}
</style> 