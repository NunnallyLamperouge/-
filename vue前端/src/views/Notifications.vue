<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { get, post, postJSON } from '@/utility/request'
import { message } from 'ant-design-vue'
import SpinLoaderLarge from '@/components/spinLoaderLarge.vue'

const userStore = useUserStore()
const notifications = ref([])
const loading = ref(true)
const pageNum = ref(1)
const pageSize = ref(10)
const total = ref(0)
const selectedTab = ref('all')
const notificationSettings = ref(null)
const settingsLoading = ref(false)
const showSettings = ref(false)
const showBackToTop = ref(false)

// 获取通知列表
const fetchNotifications = () => {
  loading.value = true
  get('/api/notifications', {
    pageNum: pageNum.value,
    pageSize: pageSize.value,
    type: selectedTab.value === 'all' ? '' : selectedTab.value
  }, (_, data) => {
    notifications.value = data.notifications
    total.value = data.total
    loading.value = false
  }, () => {
    loading.value = false
    message.error('获取通知失败')
  })
}

// 获取通知设置
const fetchNotificationSettings = () => {
  settingsLoading.value = true
  get('/api/notifications/settings', {}, (_, data) => {
    notificationSettings.value = data
    settingsLoading.value = false
  }, () => {
    settingsLoading.value = false
    message.error('获取通知设置失败')
  })
}

// 更新通知设置
const updateNotificationSettings = () => {
  post('/api/notifications/settings', notificationSettings.value, (_, data) => {
    message.success('设置已更新')
  }, () => {
    message.error('更新设置失败')
  })
}

// 标记通知为已读
const markAsRead = (id) => {
  post(`/api/notifications/${id}/read`, {}, (_, data) => {
    const notification = notifications.value.find(n => n.ID === id)
    if (notification) {
      notification.IsRead = true
      message.success('已标记为已读')
    }
  }, () => {
    message.error('标记已读失败')
  })
}

// 标记所有通知为已读
const markAllAsRead = () => {
  post('/api/notifications/read-all', {}, (_, data) => {
    notifications.value.forEach(notification => {
      notification.IsRead = true
    })
    message.success('已全部标记为已读')
  }, () => {
    message.error('标记全部已读失败')
  })
}

// 删除通知
const deleteNotification = (id) => {
  del(`/api/notifications/${id}`, {}, (_, data) => {
    notifications.value = notifications.value.filter(n => n.ID !== id)
    message.success('删除成功')
  }, () => {
    message.error('删除失败')
  })
}

// 删除所有通知
const deleteAllNotifications = () => {
  if (!confirm('确定要删除所有通知吗？')) return
  
  del('/api/notifications', {}, (_, data) => {
    notifications.value = []
    total.value = 0
    message.success('已清空所有通知')
  }, () => {
    message.error('清空通知失败')
  })
}

// 切换标签
const changeTab = (tab) => {
  selectedTab.value = tab
  pageNum.value = 1
  fetchNotifications()
}

// 切换页码
const changePage = (page) => {
  pageNum.value = page
  fetchNotifications()
}

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取通知类型名称
const getNotificationTypeName = (type) => {
  const typeMap = {
    1: '系统通知',
    2: '指导通知',
    3: '评论通知',
    4: '关注通知',
    5: '协作邀请',
    6: '内容更新'
  }
  return typeMap[type] || '未知类型'
}

// 获取通知图标
const getNotificationIcon = (type) => {
  const iconMap = {
    1: 'i-carbon-notification', // 系统通知
    2: 'i-carbon-thumbs-up', // 点赞通知
    3: 'i-carbon-chat', // 评论通知
    4: 'i-carbon-user-follow', // 关注通知
    5: 'i-carbon-collaborate', // 协作邀请
    6: 'i-carbon-document-update' // 内容更新
  }
  return iconMap[type] || 'i-carbon-notification'
}

// 计算未读通知数量
const unreadCount = computed(() => {
  return notifications.value.filter(n => !n.IsRead).length
})

// 滚动到顶部
const scrollToTop = () => {
  const container = document.querySelector('.notifications-scroll-container')
  if (container) {
    container.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
  }
}

// 处理滚动事件
const handleScroll = (e) => {
  if (e.target.scrollTop > 300) {
    showBackToTop.value = true
  } else {
    showBackToTop.value = false
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchNotifications()
  fetchNotificationSettings()
})
</script>

<template>
  <div class="flex flex-col h-[calc(100vh-5rem)] animate__animated animate__fadeIn">
    <div class="flex-1 mx-4 mb-4 rounded-2xl shadow-lg overflow-hidden animate__animated animate__fadeIn animate__delay-2s">
      <div class="h-full flex flex-col bg-white dark:bg-zinc-900">
        <!-- 顶部标题栏 -->
        <div class="flex justify-between items-center px-6 py-4 border-b border-gray-200 dark:border-gray-700">
          <h1 class="text-2xl font-bold text-gray-800 dark:text-white">通知中心</h1>
          <div class="flex gap-2">
            <button 
              @click="showSettings = !showSettings" 
              class="px-3 py-1.5 rounded-lg bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-gray-200 transition-colors transform hover:scale-105"
            >
              <span class="i-carbon-settings mr-1"></span>
              设置
            </button>
            <button 
              @click="markAllAsRead" 
              class="px-3 py-1.5 rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 dark:bg-blue-900/30 dark:text-blue-300 dark:hover:bg-blue-800/50 transition-colors transform hover:scale-105"
            >
              <span class="i-carbon-checkmark mr-1"></span>
              全部已读
            </button>
            <button 
              @click="deleteAllNotifications" 
              class="px-3 py-1.5 rounded-lg bg-red-50 text-red-600 hover:bg-red-100 dark:bg-red-900/30 dark:text-red-300 dark:hover:bg-red-800/50 transition-colors transform hover:scale-105"
            >
              <span class="i-carbon-trash-can mr-1"></span>
              清空通知
            </button>
          </div>
        </div>

        <!-- 通知设置面板 -->
        <div v-if="showSettings" class="mx-6 my-4 p-4 bg-white dark:bg-[#18181b] rounded-lg shadow-sm animate__animated animate__fadeIn">
          <h2 class="text-lg font-semibold mb-4 dark:text-white">通知设置</h2>
          <div v-if="settingsLoading" class="flex justify-center py-4">
            <SpinLoaderLarge />
          </div>
          <div v-else-if="notificationSettings" class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="flex items-center justify-between p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <span class="dark:text-gray-300">系统通知</span>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.SystemNotification" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
              </label>
            </div>
            <div class="flex items-center justify-between p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <span class="dark:text-gray-300">点赞通知</span>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.LikeNotification" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
              </label>
            </div>
            <div class="flex items-center justify-between p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <span class="dark:text-gray-300">评论通知</span>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.CommentNotification" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
              </label>
            </div>
            <div class="flex items-center justify-between p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <span class="dark:text-gray-300">关注通知</span>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.FollowNotification" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
              </label>
            </div>
            <div class="flex items-center justify-between p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <span class="dark:text-gray-300">私信通知</span>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.MessageNotification" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
              </label>
            </div>
            <div class="flex items-center justify-between p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <span class="dark:text-gray-300">邮件通知</span>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.EmailNotification" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
              </label>
            </div>
            <div class="flex items-center justify-between p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <span class="dark:text-gray-300">推送通知</span>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.PushNotification" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
              </label>
            </div>
          </div>
          <div class="flex justify-end mt-4">
            <button 
              @click="updateNotificationSettings" 
              class="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 hover:shadow-md hover:shadow-blue-500/20 transition-all transform hover:scale-105"
            >
              保存设置
            </button>
          </div>
        </div>

        <!-- 通知类型标签 -->
        <div class="flex gap-2 mx-6 my-4 overflow-x-auto pb-2 animate__animated animate__fadeIn animate__delay-1s">
          <button 
            v-for="(name, type) in {all: '全部', '1': '系统通知', '2': '任务通知', '3': '指导通知', 
                                 '4': '关注通知', '5': '协作邀请', '6': '模版发放'}"
            :key="type"
            @click="changeTab(type)" 
            class="px-4 py-2 rounded-full whitespace-nowrap transition-colors transform hover:scale-105"
            :class="selectedTab === type ? 'bg-gradient-to-r from-blue-500 to-cyan-500 text-white shadow-md shadow-blue-500/20' : 'bg-gray-100 dark:bg-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'"
          >
            {{ name }}
          </button>
        </div>

        <!-- 通知列表 -->
        <div class="flex-1 overflow-y-auto px-6 notifications-scroll-container" @scroll="handleScroll">
          <div v-if="loading" class="flex justify-center py-8">
            <SpinLoaderLarge />
          </div>
          
          <div v-else-if="notifications.length === 0" class="flex flex-col items-center justify-center py-12 animate__animated animate__fadeIn">
            <div class="text-6xl text-gray-300 dark:text-gray-600 mb-4">
              <span class="i-carbon-notification"></span>
            </div>
            <p class="text-gray-500 dark:text-gray-400">暂无通知</p>
          </div>
          
          <div 
            v-else
            class="space-y-4 pb-6"
          >
            <div
              v-for="(notification, index) in notifications" 
              :key="notification.ID"
              class="p-4 bg-white dark:bg-zinc-800 rounded-lg shadow-sm transition-all hover:shadow-md relative animate__animated animate__fadeIn"
              :class="{
                'border-l-4 border-blue-500': !notification.IsRead,
                'animate__delay-1s': index < 3,
                'animate__delay-2s': index >= 3 && index < 6,
                'animate__delay-3s': index >= 6
              }"
            >
              <div class="flex items-start gap-4">
                <div :class="[getNotificationIcon(notification.Type), 'text-2xl text-blue-500 mt-1']"></div>
                <div class="flex-1">
                  <div class="flex justify-between items-start">
                    <h3 class="font-medium dark:text-white">{{ notification.Title || getNotificationTypeName(notification.Type) }}</h3>
                    <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatTime(notification.CreatedAt) }}</span>
                  </div>
                  <p class="text-gray-700 dark:text-gray-300 mt-1">{{ notification.Content }}</p>
                  <div class="flex justify-between items-center mt-2">
                    <span class="text-xs px-2 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400">
                      {{ getNotificationTypeName(notification.Type) }}
                    </span>
                    <div class="flex gap-2">
                      <button 
                        v-if="!notification.IsRead"
                        @click="markAsRead(notification.ID)" 
                        class="text-xs px-2 py-1 rounded bg-blue-50 text-blue-600 hover:bg-blue-100 dark:bg-blue-900/30 dark:text-blue-300 dark:hover:bg-blue-800/50 transition-colors transform hover:scale-105"
                      >
                        标为已读
                      </button>
                      <button 
                        @click="deleteNotification(notification.ID)" 
                        class="text-xs px-2 py-1 rounded bg-red-50 text-red-600 hover:bg-red-100 dark:bg-red-900/30 dark:text-red-300 dark:hover:bg-red-800/50 transition-colors transform hover:scale-105"
                      >
                        删除
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 回到顶部按钮 -->
          <div 
            v-show="showBackToTop" 
            class="fixed bottom-6 right-6 z-10 cursor-pointer bg-gradient-to-r from-blue-500 to-cyan-500 hover:opacity-90 text-white p-3 rounded-full shadow-lg hover:shadow-blue-500/20 transition-all duration-300 animate__animated animate__fadeIn transform hover:scale-110"
            @click="scrollToTop"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 10l7-7m0 0l7 7m-7-7v18" />
            </svg>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="total > pageSize" class="flex justify-center py-4 border-t border-gray-200 dark:border-gray-700">
          <div class="flex gap-2">
            <button 
              v-for="page in Math.ceil(total / pageSize)" 
              :key="page"
              @click="changePage(page)"
              class="w-10 h-10 flex items-center justify-center rounded-full transition-colors transform hover:scale-105"
              :class="pageNum === page 
                ? 'bg-gradient-to-r from-blue-500 to-cyan-500 text-white shadow-md shadow-blue-500/20' 
                : 'bg-gray-100 dark:bg-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'"
            >
              {{ page }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 动画相关样式 */
.animate__animated {
  animation-duration: 0.5s;
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

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate__fadeIn {
  animation-name: fadeIn;
}

/* 滚动条样式 */
.notifications-scroll-container::-webkit-scrollbar {
  width: 6px;
}

.notifications-scroll-container::-webkit-scrollbar-track {
  background: transparent;
}

.notifications-scroll-container::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 3px;
}

.notifications-scroll-container::-webkit-scrollbar-thumb:hover {
  background-color: rgba(107, 114, 128, 0.7);
}

/* 深色模式滚动条 */
.dark .notifications-scroll-container::-webkit-scrollbar-thumb {
  background-color: rgba(75, 85, 99, 0.5);
}

.dark .notifications-scroll-container::-webkit-scrollbar-thumb:hover {
  background-color: rgba(107, 114, 128, 0.7);
}
</style> 