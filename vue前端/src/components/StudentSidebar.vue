<template>
  <div 
    class="sidebar fixed h-[calc(100vh-4rem)] overflow-y-auto left-0 top-16 bg-app"
    :class="{ 'collapsed': isCollapsed, 'expanded': !isCollapsed }"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <div class="flex justify-end p-2 transition-all duration-300">
      <menu-unfold-outlined 
        v-if="isCollapsed" 
        class="sidebar-toggle text-muted cursor-pointer text-lg hover:text-primary transition-all duration-300"
        @click="toggleSidebar"
      />
      <menu-fold-outlined 
        v-else 
        class="sidebar-toggle text-muted cursor-pointer text-lg hover:text-primary transition-all duration-300"
        @click="toggleSidebar"
      />
    </div>
    <nav class="mt-2 px-2 space-y-1">
      <router-link
        v-for="item in menuItems"
        :key="item.path"
        :to="item.path"
        class="menu-item group flex items-center px-2 py-2 text-base font-medium rounded-md transition-all duration-300"
        :class="[
          $route.path.includes(item.path)
            ? 'bg-sunken text-blue-600 dark:text-blue-400'
            : 'text-muted hover:bg-elevated'
        ]"
      >
        <component 
          :is="item.icon" 
          class="menu-icon flex-shrink-0 transition-transform duration-300" 
          :class="[isCollapsed ? 'mx-auto' : 'mr-3']"
        />
        <span 
          class="menu-text whitespace-nowrap overflow-hidden transition-all duration-300"
          :class="{ 'opacity-0': isCollapsed, 'opacity-100': !isCollapsed }"
        >{{ item.name }}</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import {
  EditOutlined,
  TeamOutlined,
  CheckCircleOutlined,
  DashboardOutlined,
  ApartmentOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined
} from '@ant-design/icons-vue'

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false
  },
  autoCollapse: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:collapsed'])

const isCollapsed = ref(props.collapsed)
const route = useRoute()
const hoverTimer = ref(null)

// 监听props变化
watch(() => props.collapsed, (newVal) => {
  isCollapsed.value = newVal
})

// 菜单项
const menuItems = [
  { name: '概览', path: '/student/dashboard', icon: DashboardOutlined },
  { name: '写作', path: '/student/editor', icon: EditOutlined },
  { name: '模拟答辩', path: '/student/defense', icon: TeamOutlined },
  { name: '格式检查', path: '/student/format', icon: CheckCircleOutlined },
  { name: '数据建模', path: '/student/modeling', icon: ApartmentOutlined }
]

const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
  emit('update:collapsed', isCollapsed.value)
}

const handleMouseEnter = () => {
  if (props.autoCollapse && isCollapsed.value) {
    // 清除之前的定时器
    if (hoverTimer.value) {
      clearTimeout(hoverTimer.value)
      hoverTimer.value = null
    }
    
    // 设置新的定时器，延迟展开
    hoverTimer.value = setTimeout(() => {
      isCollapsed.value = false
      emit('update:collapsed', false)
    }, 200)
  }
}

const handleMouseLeave = () => {
  if (props.autoCollapse && !isCollapsed.value) {
    // 清除之前的定时器
    if (hoverTimer.value) {
      clearTimeout(hoverTimer.value)
      hoverTimer.value = null
    }
    
    // 设置新的定时器，延迟折叠
    hoverTimer.value = setTimeout(() => {
      isCollapsed.value = true
      emit('update:collapsed', true)
    }, 300)
  }
}

// 组件卸载时清除定时器
onMounted(() => {
  return () => {
    if (hoverTimer.value) {
      clearTimeout(hoverTimer.value)
    }
  }
})
</script>

<style scoped>
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 3px;
}

.dark ::-webkit-scrollbar-thumb {
  background-color: rgba(75, 85, 99, 0.5);
}

/* 侧边栏基础样式 */
.sidebar {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  z-index: 100;
  will-change: width;
  border-right: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 1s ease;
}

.dark .sidebar {
  border-right: 1px solid rgba(255, 255, 255, 0.05);
}

/* 侧边栏状态 */
.sidebar.collapsed {
  width: 4rem;
}

.sidebar.expanded {
  width: 16rem;
}

/* 菜单图标样式 */
.menu-icon {
  height: 1.5rem;
  width: 1.5rem;
  transition: all 1s ease;
}

/* 菜单文本样式 */
.menu-text {
  max-width: 200px;
  transform-origin: left;
  transition: all 1s ease;
}

.collapsed .menu-text {
  max-width: 0;
  opacity: 0;
  transform: translateX(-10px);
}

/* 菜单项悬停效果 */
.menu-item {
  transition: all 1s ease;
}

.menu-item:hover .menu-icon {
  transform: scale(1.05);
}

.sidebar-toggle {
  transition: all 0.3s ease;
}

.sidebar-toggle:hover {
  transform: scale(1.1);
}

.text-primary {
  color: #1890ff;
}
</style> 