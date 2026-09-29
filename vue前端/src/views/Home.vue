<template>
  <div class="min-h-screen bg-app theme-transition">
    <!-- 导航栏 -->
    <nav class="navbar fixed w-full z-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <img class="h-8 w-auto" src="@/assets/img/logo.webp" alt="研途智伴" />
            <span class="ml-2 text-xl font-bold text-app">研途智伴</span>
          </div>
          <div class="flex items-center space-x-4">
            <ThemeSwitcher />
            <router-link
              v-if="!userStore.isLogin"
              to="/auth/login"
              class="btn-primary"
            >
              登录
            </router-link>
            <a-dropdown v-else>
              <a-button class="flex items-center space-x-2">
                <user-outlined />
                <span>{{ userStore.user?.username }}</span>
              </a-button>
              <template #overlay>
                <a-menu>
                  <a-menu-item key="profile">
                    <user-outlined />
                    <span>个人中心</span>
                  </a-menu-item>
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

    <!-- Hero Section with Wave Background -->
    <WaveBackground height="600px" class="pt-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        <div class="pt-16 pb-20 text-center">
          <h1 class="text-4xl font-extrabold tracking-tight sm:text-5xl lg:text-6xl">
            <span class="block text-white">智能论文写作助手</span>
            <span class="block text-indigo-200 mt-2">让您的学术之路更轻松</span>
          </h1>
          <p class="mt-6 max-w-lg mx-auto text-center text-xl text-indigo-200 sm:max-w-3xl">
            提供全方位的论文写作支持，从格式规范到答辩准备，助您轻松完成学术任务
          </p>
          <div class="mt-10 max-w-sm mx-auto sm:max-w-none sm:flex sm:justify-center">
            <div class="space-y-4 sm:space-y-0 sm:space-x-4">
              <router-link
                to="/student/dashboard"
                class="btn-primary-light"
              >
                开始写作
              </router-link>
              <router-link
                to="/teacher/dashboard"
                class="btn-secondary-light"
              >
              我是老师
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </WaveBackground>

    <!-- 功能特点 -->
    <div id="features" class="py-16 bg-app overflow-hidden lg:py-24">
      <div class="relative max-w-xl mx-auto px-4 sm:px-6 lg:px-8 lg:max-w-7xl">
        <div class="relative">
          <h2 class="text-center text-3xl leading-8 font-extrabold tracking-tight text-app sm:text-4xl">
            为什么选择 研途智伴
          </h2>
          <p class="mt-4 max-w-3xl mx-auto text-center text-xl text-muted">
            全方位的论文写作解决方案，让您的论文写作更轻松、更高效
          </p>
        </div>

        <div class="mt-12 grid gap-8 lg:grid-cols-2">
          <!-- 语义增强预览 -->
          <FeatureCard
            title="智能语义增强"
            description="基于深度学习的语义分析，提供专业的学术表达建议和语法优化"
            color="blue"
            :delay="200"
          >
            <template #icon>
              <BulbOutlined class="h-6 w-6" />
            </template>
            <SemanticEnhancerPreview height="300px" class="mt-4" />
          </FeatureCard>

          <!-- 进度追踪 -->
          <FeatureCard
            title="智能进度管理"
            description="可视化的写作进度追踪，帮助您合理规划时间，按期完成论文"
            color="green"
            :delay="400"
          >
            <template #icon>
              <chart-outlined class="h-6 w-6" />
            </template>
            <ProgressTracker height="300px" class="mt-4" />
          </FeatureCard>

          <!-- 智能排版 -->
          <FeatureCard
            title="智能排版转换"
            description="一键转换不同院校和期刊的论文格式要求，省时省力"
            color="purple"
            :delay="600"
          >
            <template #icon>
              <layout-outlined class="h-6 w-6" />
            </template>
            <SmartFormattingPreview height="300px" class="mt-4" />
          </FeatureCard>

          <!-- 虚拟答辩 -->
          <FeatureCard
            title="虚拟答辩模拟"
            description="基于AI的答辩模拟系统，提前适应答辩环境，提高答辩成功率"
            color="indigo"
            :delay="800"
          >
            <template #icon>
              <team-outlined class="h-6 w-6" />
            </template>
            <VirtualDefenseSimulator height="300px" class="mt-4" />
          </FeatureCard>
        </div>
      </div>
    </div>

    <!-- 数据模型展示 -->
    <div class="bg-sunken py-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="lg:text-center">
          <h2 class="text-3xl font-extrabold text-app sm:text-4xl">
            专业的数据支持
          </h2>
          <p class="mt-4 max-w-2xl text-xl text-muted lg:mx-auto">
            完善的数据模型设计，确保您的论文数据安全可靠
          </p>
        </div>
        <div class="mt-10">
          <ERDiagramPreview height="400px" />
        </div>
      </div>
    </div>

    <!-- 统计数据 -->
    <div class="bg-app pt-12 sm:pt-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="max-w-4xl mx-auto text-center">
          <h2 class="text-3xl font-extrabold text-app sm:text-4xl">
            值得信赖的论文助手
          </h2>
          <p class="mt-3 text-xl text-muted sm:mt-4">
            已帮助数千名学生完成他们的学术论文
          </p>
        </div>
      </div>
      <div class="mt-10 pb-12 bg-app sm:pb-16">
        <div class="relative">
          <div class="absolute inset-0 h-1/2 bg-sunken"></div>
          <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div class="max-w-4xl mx-auto">
              <dl class="rounded-lg bg-app shadow-lg sm:grid sm:grid-cols-3">
                <div class="flex flex-col border-b border-app p-6 text-center sm:border-0 sm:border-r">
                  <dt class="order-2 mt-2 text-lg leading-6 font-medium text-muted">
                    活跃用户
                  </dt>
                  <dd class="order-1 text-5xl font-extrabold text-blue-600">
                    100,000+
                  </dd>
                </div>
                <div class="flex flex-col border-t border-b border-app p-6 text-center sm:border-0 sm:border-l sm:border-r">
                  <dt class="order-2 mt-2 text-lg leading-6 font-medium text-muted">
                    完成论文
                  </dt>
                  <dd class="order-1 text-5xl font-extrabold text-blue-600">
                    50,000+
                  </dd>
                </div>
                <div class="flex flex-col border-t border-app p-6 text-center sm:border-0 sm:border-l">
                  <dt class="order-2 mt-2 text-lg leading-6 font-medium text-muted">
                    用户满意度
                  </dt>
                  <dd class="order-1 text-5xl font-extrabold text-blue-600">
                    98%
                  </dd>
                </div>
              </dl>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部CTA -->
    <div class="bg-blue-600 dark:bg-blue-700">
      <div class="max-w-2xl mx-auto text-center py-16 px-4 sm:py-20 sm:px-6 lg:px-8">
        <h2 class="text-3xl font-extrabold text-white sm:text-4xl">
          <span class="block">准备开始您的论文写作之旅？</span>
        </h2>
        <p class="mt-4 text-lg leading-6 text-blue-200">
          立即注册，享受智能写作助手的全方位支持。
        </p>
        <router-link
          to="/auth/register"
          class="mt-8 w-full inline-flex items-center justify-center px-5 py-3 border border-transparent text-base font-medium rounded-md text-blue-600 bg-white hover:bg-blue-50 sm:w-auto"
        >
          免费注册
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useUserStore } from '@/stores/user'
import { message } from 'ant-design-vue'
import { 
  UserOutlined, 
  LogoutOutlined,
  BulbOutlined,
  TeamOutlined,
  LayoutOutlined
} from '@ant-design/icons-vue'
import ThemeSwitcher from '@/components/ThemeSwitcher.vue'
import WaveBackground from '@/components/WaveBackground.vue'
import FeatureCard from '@/components/FeatureCard.vue'
import SemanticEnhancerPreview from '@/components/SemanticEnhancerPreview.vue'
import ProgressTracker from '@/components/ProgressTracker.vue'
import SmartFormattingPreview from '@/components/SmartFormattingPreview.vue'
import VirtualDefenseSimulator from '@/components/VirtualDefenseSimulator.vue'
import ERDiagramPreview from '@/components/ERDiagramPreview.vue'

const userStore = useUserStore()

const handleLogout = async () => {
  try {
    await userStore.logout()
    message.success('退出登录成功')
  } catch (error) {
    message.error('退出登录失败')
  }
}
</script>

<style scoped>
.btn-primary-light {
  @apply inline-flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-blue-700 bg-white hover:bg-blue-50 md:py-4 md:text-lg md:px-10;
}

.btn-secondary-light {
  @apply inline-flex items-center justify-center px-8 py-3 border border-white text-base font-medium rounded-md text-white hover:bg-white hover:bg-opacity-10 md:py-4 md:text-lg md:px-10;
}

/* 优化暗色模式下的渐变效果 */
:deep(.dark) .bg-gradient-to-r {
  opacity: 0.9;
}
</style> 