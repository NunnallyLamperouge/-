<template>
  <nav class="navbar fixed w-full z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex">
          <!-- Logo -->
          <div class="flex-shrink-0 flex items-center">
            <img class="h-8 w-auto" src="@/assets/img/logo.webp" alt="研途智伴" />
            <span class="ml-2 text-xl font-bold text-app hidden md:block">研途智伴</span>
          </div>
          
          
        </div>

        <div class="flex items-center space-x-4">
          
          
          <!-- 主题切换 -->
          <ThemeSwitcher />
          
          <!-- 通知中心 -->
          <a-badge :count="notificationCount">
            <bell-outlined 
              class="text-muted text-xl cursor-pointer" 
              @click="router.push('/teacher/notifications')"
            />
          </a-badge>

          <!-- 用户菜单 -->
          <a-dropdown>
            <a-avatar class="cursor-pointer bg-blue-500">
              {{ userInitials }}
            </a-avatar>
            <template #overlay>
              <a-menu>
                <a-menu-item key="profile">
                  <user-outlined />
                  <span>个人信息</span>
                </a-menu-item>
                <a-menu-item key="settings">
                  <setting-outlined />
                  <span>偏好设置</span>
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="logout" @click="handleLogout">
                  <logout-outlined />
                  <span>退出登录</span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
    </div>

    <!-- 通知抽屉 -->
    <a-drawer
      v-model:visible="notificationDrawerVisible"
      title="通知中心"
      placement="right"
      width="380"
    >
      <div class="space-y-4">
        <div v-for="notice in notifications" :key="notice.id" class="notification-item">
          <div class="flex items-start p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg">
            <div class="flex-shrink-0 mt-1">
              <component 
                :is="notice.icon" 
                class="text-lg"
                :class="notice.iconColor"
              />
            </div>
            <div class="ml-3 flex-1">
              <p class="text-sm font-medium text-app">{{ notice.title }}</p>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {{ notice.content }}
              </p>
              <p class="mt-2 text-xs text-gray-400">{{ notice.time }}</p>
            </div>
          </div>
        </div>
      </div>
    </a-drawer>
  </nav>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import {
  UserOutlined,
  SettingOutlined,
  LogoutOutlined,
  BellOutlined,
  FileTextOutlined,
  TeamOutlined,
  CalendarOutlined
} from '@ant-design/icons-vue'
import ThemeSwitcher from '@/components/ThemeSwitcher.vue'

const props = defineProps({

})

const emit = defineEmits(['department-change', 'logout'])

const userStore = useUserStore()
const router = useRouter()

// 院系选择
const selectedDepartment = ref(props.currentDepartment)
const handleDepartmentChange = (value) => {
  emit('department-change', value)
}

// 评审进度
const reviewProgress = ref(75)

// 用户信息
const userInitials = computed(() => {
  return userStore.user?.username?.charAt(0).toUpperCase() || '教'
})

// 通知功能
const notificationCount = ref(3)
const notificationDrawerVisible = ref(false)
const notifications = ref([
  {
    id: 1,
    title: '论文评审提醒',
    content: '张三的论文《机器学习在图像识别中的应用》等待您的评审',
    time: '10分钟前',
    icon: FileTextOutlined,
    iconColor: 'text-blue-500'
  },
  {
    id: 2,
    title: '答辩安排确认',
    content: '请确认本周五下午的答辩时间安排',
    time: '30分钟前',
    icon: CalendarOutlined,
    iconColor: 'text-purple-500'
  },
  {
    id: 3,
    title: '指导学生消息',
    content: '李四提交了论文新版本，请查看修改内容',
    time: '2小时前',
    icon: TeamOutlined,
    iconColor: 'text-green-500'
  }
])

const showNotifications = () => {
  notificationDrawerVisible.value = true
}

const handleLogout = async () => {
  try {
    await userStore.logout()
    emit('logout')
    router.push('/auth/login')
  } catch (error) {
    console.error('退出失败:', error)
  }
}
</script>

<style scoped>
.navbar {
  @apply bg-white/90 dark:bg-[#1a1e2c]/95 backdrop-blur-md border-b border-gray-200 dark:border-gray-800;
}

.notification-item:hover {
  transform: translateY(-1px);
  transition: all 0.3s ease;
}

/* 暗色模式优化 */
:deep(.dark) .ant-select:not(.ant-select-customize-input) .ant-select-selector {
  @apply bg-gray-800/50 border-gray-700/20;
}

:deep(.dark) .ant-drawer-content {
  @apply bg-[#1a1e2c];
}

:deep(.dark) .ant-drawer-header {
  @apply bg-[#1a1e2c] border-gray-700/20;
}

:deep(.dark) .ant-drawer-title {
  @apply text-gray-200;
}

:deep(.dark) .ant-menu {
  @apply bg-[#1a1e2c] border-gray-700/20;
}

:deep(.dark) .ant-menu-item {
  @apply text-gray-200;
}

:deep(.dark) .ant-menu-item:hover {
  @apply bg-gray-700/20;
}

:deep(.dark) .ant-dropdown-menu {
  @apply bg-[#1a1e2c] border-gray-700/20;
}

:deep(.dark) .ant-dropdown-menu-item {
  @apply text-gray-200;
}

:deep(.dark) .ant-dropdown-menu-item:hover {
  @apply bg-gray-700/20;
}

:deep(.text-muted) {
  @apply text-gray-600 dark:text-gray-400;
}

:deep(.text-app) {
  @apply text-gray-900 dark:text-gray-100;
}
</style> 