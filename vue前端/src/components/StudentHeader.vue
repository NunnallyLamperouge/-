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
          <a-badge :count="3" @click="router.push('/student/notifications')">
            <bell-outlined class="text-muted text-xl cursor-pointer" />
          </a-badge>

          <!-- 用户菜单 -->
          <a-dropdown>
            <a-avatar class="cursor-pointer bg-blue-500">
              {{ userStore.user?.username?.charAt(0).toUpperCase() }}
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
  </nav>
</template>

<script setup>
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import {
  UserOutlined,
  SettingOutlined,
  LogoutOutlined,
  BellOutlined
} from '@ant-design/icons-vue'
import ThemeSwitcher from '@/components/ThemeSwitcher.vue'

const emit = defineEmits(['project-change', 'logout'])

const userStore = useUserStore()
const router = useRouter()

const handleProjectChange = (value) => {
  emit('project-change', value)
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