<template>
  <div class="min-h-screen bg-app theme-transition flex flex-col">
    <!-- 顶部导航栏 -->
    <StudentHeader 
      :projects="projects" 
      :current-project="currentProject" 
      @project-change="handleProjectChange"
      @logout="handleLogout"
      class="fixed top-0 left-0 right-0 z-20"
      v-motion
      :initial="{ opacity: 0, y: -20 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 600 } }"
    />

    <div class="flex flex-1 pt-16 h-screen">
      <!-- 侧边栏 -->
      <div
        ref="sidebarRef"
        class="fixed top-16 bottom-0 left-0 z-10 transition-all duration-500 ease-in-out text-theme-switch shadow-lg transform-gpu"
        :class="[
          sidebarExpanded ? 'sidebar-expanded' : 'sidebar-collapsed',
          'bg-white/90 dark:bg-[#1a1e2c]/95 backdrop-blur-md'
        ]"
        @mouseenter="expandSidebar"
        @mouseleave="collapseSidebar"
      >
        <div class="h-full flex flex-col p-3 gap-3 overflow-hidden">
          <!-- 顶部折叠/展开按钮 -->
          <div class="flex justify-end p-2">
            <a-tooltip placement="right" :title="sidebarExpanded ? '收起菜单' : '展开菜单'">
              <menu-unfold-outlined 
                v-if="!sidebarExpanded" 
                class="sidebar-toggle text-muted cursor-pointer text-lg hover:text-primary transition-all duration-300"
                @click="toggleSidebar"
              />
              <menu-fold-outlined 
                v-else 
                class="sidebar-toggle text-muted cursor-pointer text-lg hover:text-primary transition-all duration-300"
                @click="toggleSidebar"
              />
            </a-tooltip>
          </div>

          <!-- 导航菜单 -->
          <div class="flex-1 overflow-y-auto scrollbar-hide mt-2">
            <nav class="flex flex-col gap-2">
              <router-link
                v-for="(item, index) in menuItems"
                :key="item.path"
                :to="item.path"
                class="flex items-center p-2 rounded-lg cursor-pointer transition-all duration-300 relative overflow-hidden group"
                :class="[
                  $route.path.includes(item.path)
                    ? 'bg-gradient-to-r from-blue-50/80 to-transparent dark:from-indigo-900/20 dark:to-transparent text-blue-600 dark:text-indigo-300 font-medium'
                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800/30'
                ]"
                v-motion
                :initial="{ opacity: 0, x: -20 }"
                :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 100 * index } }"
              >
                <div 
                  v-if="$route.path.includes(item.path)" 
                  class="absolute left-0 top-0 h-full w-1 bg-gradient-to-b from-blue-500 to-blue-600 dark:from-indigo-500 dark:to-purple-600"
                ></div>
                <component 
                  :is="item.icon" 
                  class="flex-shrink-0 size-6 transition-all duration-300 z-10 ml-2"
                  :class="{
                    'scale-110': !sidebarExpanded,
                    'text-blue-600 dark:text-indigo-400': $route.path.includes(item.path),
                    'group-hover:rotate-6': !$route.path.includes(item.path)
                  }"
                />
                <div 
                  class="ml-3 overflow-hidden transition-all duration-500 ease-in-out z-10"
                  :class="[
                    sidebarExpanded 
                      ? (textVisible ? 'opacity-100 max-w-40 translate-x-0' : 'opacity-0 max-w-0 -translate-x-4') 
                      : 'opacity-0 max-w-0 -translate-x-4'
                  ]"
                >
                  <span class="whitespace-nowrap">{{ item.name }}</span>
                </div>
                <a-badge v-if="item.badge" :count="item.badge" :dot="item.dot" class="ml-auto" :class="{'opacity-0': !sidebarExpanded}"/>
                <div 
                  v-if="$route.path.includes(item.path)"
                  class="absolute right-0 bottom-0 w-2 h-2 rounded-full bg-blue-500 dark:bg-indigo-500 transform translate-x-1/2 translate-y-1/2"
                ></div>
              </router-link>
            </nav>
          </div>
        </div>
      </div>

      <!-- 主内容区 -->
      <main 
        ref="mainContentRef"
        class="flex-1 transition-all duration-300 ease-in-out will-change-transform transform-gpu overflow-auto ml-20"
        style="backface-visibility: hidden; perspective: 1000px; height: calc(100vh - 4rem);"
        v-motion
        :initial="{ opacity: 0 }"
        :enter="{ opacity: 1, transition: { duration: 800, delay: 200 } }"
      >
        <div class="py-6 min-h-full">
          <div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
            <router-view v-slot="{ Component }">
              <transition 
                name="fade" 
                mode="out-in"
                @before-enter="onBeforeEnter"
                @after-leave="onAfterLeave"
              >
                <component :is="Component" />
              </transition>
            </router-view>
          </div>
        </div>
      </main>
    </div>

    <!-- 页面底部 - 可选的 copyright 或其他信息 -->
    <div 
      class="text-center py-4 text-sm text-gray-500 dark:text-gray-400 bg-white/80 dark:bg-[#1a1e2c]/80 backdrop-blur-sm border-t border-gray-200 dark:border-gray-800"
      v-motion
      :initial="{ opacity: 0, y: 20 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 600, delay: 800 } }"
    >
      © {{ new Date().getFullYear() }} 研途智伴.
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useMotion } from '@vueuse/motion'
import { Badge, Tooltip, message } from 'ant-design-vue'
import StudentHeader from '@/components/StudentHeader.vue'
import {
  EditOutlined,
  TeamOutlined,
  CheckCircleOutlined,
  DashboardOutlined,
  ApartmentOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserSwitchOutlined,
  BellOutlined,
  CalendarOutlined,
  CopyOutlined
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 侧边栏状态管理
const sidebarExpanded = ref(false)
const textVisible = ref(true)
const sidebarCollapsed = ref(false)
const sidebarRef = ref(null)
const mainContentRef = ref(null)

// 应用动画效果到侧边栏
const sidebar = useMotion(sidebarRef, {
  initial: { opacity: 0.5, x: -20 },
  enter: { opacity: 1, x: 0, transition: { duration: 800 } },
  visible: { opacity: 1, x: 0 },
  hovered: { scale: 1.01, transition: { type: 'spring', stiffness: 500, damping: 20 } }
})

// 监听路由变化来触发某些动画效果
watch(() => route.path, (newPath) => {
  message.destroy()
  // 可以在这里添加路由变化时的过渡效果
})

// 监听侧边栏折叠状态变化
watch(sidebarCollapsed, (newVal) => {
  if (newVal) {
    collapseSidebar(false)
  } else {
    expandSidebar(false)
  }
})

// 页面加载完成时的初始化
onMounted(() => {
  // 可在此添加初始化逻辑
})

// 页面过渡前的回调函数
const onBeforeEnter = (el) => {
  // 过渡前的准备工作
}

// 页面过渡后的回调函数
const onAfterLeave = (el) => {
  // 过渡后的清理工作
}

// 菜单项 (添加了徽章示例)
const menuItems = [
  { name: '概览', path: '/student/dashboard', icon: DashboardOutlined },
  { name: '选择导师', path: '/student/select-teacher', icon: UserSwitchOutlined },
  { name: '写作', path: '/student/editor', icon: EditOutlined, badge: 1 },
  { name: '导师指导', path: '/student/chat', icon: TeamOutlined, dot: true },
  { name: '模拟答辩', path: '/student/format', icon: CheckCircleOutlined },
  { name: '任务管理', path: '/student/missions', icon: CalendarOutlined },
  { name: '数据建模', path: '/student/modeling', icon: ApartmentOutlined },
  { name: '知识库', path: '/student/documents', icon:CopyOutlined }
]

// 展开侧边栏
const expandSidebar = (updateCollapsed = true) => {
  sidebarExpanded.value = true
  // 延迟显示文本，创建平滑过渡效果
  setTimeout(() => {
    textVisible.value = true
  }, 200)
  
  if (updateCollapsed) {
    sidebarCollapsed.value = false
  }
}

// 收起侧边栏
const collapseSidebar = (updateCollapsed = true) => {
  textVisible.value = false
  // 延迟收起侧边栏，让文本先消失
  setTimeout(() => {
    sidebarExpanded.value = false
  }, 200)
  
  if (updateCollapsed) {
    sidebarCollapsed.value = true
  }
}

// 切换侧边栏状态
const toggleSidebar = () => {
  if (sidebarExpanded.value) {
    collapseSidebar()
  } else {
    expandSidebar()
  }
}

// 项目列表
const projects = ref([
  { id: 1, name: '毕业论文' },
  { id: 2, name: '开题报告' },
  { id: 3, name: '文献综述' }
])
const currentProject = ref(1)

const handleProjectChange = (value) => {
  console.log('切换项目:', value)
}

const handleLogout = async () => {
  try {
    await userStore.logout()
    message.success('退出成功')
    router.push('/auth/login')
  } catch (error) {
    console.error('退出失败:', error)
    message.error('退出失败: ' + error.message)
  }
}
</script>

<style scoped>
/* 侧边栏状态 */
.sidebar-expanded {
  width: 16rem; /* 256px */
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  z-index: 15; /* 确保展开时在主内容上方 */
}

.sidebar-collapsed {
  width: 5rem; /* 80px */
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);
}

/* 暗色模式下的侧边栏增强 */
.dark .sidebar-expanded {
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
  border-right: 1px solid rgba(255, 255, 255, 0.05);
}

.dark .sidebar-collapsed {
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  border-right: 1px solid rgba(255, 255, 255, 0.03);
}

/* 侧边栏动画 */
@keyframes expandSidebar {
  from { width: 5rem; }
  to { width: 16rem; }
}

@keyframes collapseSidebar {
  from { width: 16rem; }
  to { width: 5rem; }
}

.sidebar-expanded {
  animation: expandSidebar 0.5s ease-in-out;
}

.sidebar-collapsed {
  animation: collapseSidebar 0.5s ease-in-out;
}

/* 主内容区滚动条样式 */
main::-webkit-scrollbar {
  width: 6px;
}

main::-webkit-scrollbar-track {
  background: transparent;
}

main::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 3px;
}

.dark main::-webkit-scrollbar-thumb {
  background-color: rgba(75, 85, 99, 0.5);
}

/* 隐藏滚动条 */
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

/* 悬停效果 */
.sidebar-toggle:hover {
  transform: scale(1.1);
}

.text-primary {
  color: #1890ff;
}

/* 页面过渡效果 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 菜单项选中效果 */
.router-link-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 100%;
  background: linear-gradient(90deg, rgba(59, 130, 246, 0.08) 0%, rgba(59, 130, 246, 0) 100%);
  z-index: 0;
}

.dark .router-link-active::before {
  background: linear-gradient(90deg, rgba(99, 102, 241, 0.1) 0%, rgba(139, 92, 246, 0.05) 50%, rgba(0, 0, 0, 0) 100%);
}

/* 暗色模式下的文本颜色 */
.dark .text-primary {
  color: #818cf8;
}

.dark .text-muted {
  color: #9ca3af;
}

/* 菜单项悬停效果 */
.router-link-exact-active .sidebar-icon,
.router-link-active .sidebar-icon {
  transform: scale(1.2);
}

/* 添加菜单项悬停效果 */
.nav-item {
  position: relative;
  overflow: hidden;
}

.nav-item::after {
  content: '';
  position: absolute;
  left: 0;
  bottom: 0;
  height: 2px;
  width: 0;
  background: linear-gradient(90deg, #3b82f6, #818cf8);
  transition: width 0.3s ease;
}

.nav-item:hover::after {
  width: 100%;
}

/* Ant Design组件样式调整 */
:deep(.ant-badge-count) {
  box-shadow: 0 0 0 1px #fff;
}

:deep(.ant-tooltip) {
  max-width: 300px;
}
</style> 