<template>
  <div class="er-diagram-preview relative rounded-lg overflow-hidden shadow-xl" :style="{ height: height }">
    <div class="absolute inset-0 bg-gray-100 dark:bg-gray-800">
      <svg class="w-full h-full" viewBox="0 0 800 600" xmlns="http://www.w3.org/2000/svg">
        <!-- 网格背景 -->
        <defs>
          <pattern id="grid" width="20" height="20" patternUnits="userSpaceOnUse">
            <path d="M 20 0 L 0 0 0 20" fill="none" stroke="rgba(100,100,100,0.1)" stroke-width="1"/>
          </pattern>
        </defs>
        <rect width="100%" height="100%" fill="url(#grid)" />
        
        <!-- 实体和关系 -->
        <g transform="translate(400, 300)" class="er-diagram-content">
          <!-- 用户表 -->
          <g transform="translate(-200, -150)">
            <rect width="160" height="120" rx="5" fill="#3b82f6" fill-opacity="0.2" stroke="#3b82f6" stroke-width="2" />
            <rect width="160" height="30" rx="5" fill="#3b82f6" stroke="#3b82f6" stroke-width="2" />
            <text x="80" y="20" text-anchor="middle" fill="#3b82f6" font-weight="bold">用户表</text>
            <text x="10" y="50" fill="#333333" class="dark:text-gray-200">+ 用户ID (PK)</text>
            <text x="10" y="70" fill="#333333" class="dark:text-gray-200">+ 用户名</text>
            <text x="10" y="90" fill="#333333" class="dark:text-gray-200">+ 邮箱</text>
            <text x="10" y="110" fill="#333333" class="dark:text-gray-200">+ 角色</text>
          </g>
          
          <!-- 论文表 -->
          <g transform="translate(100, -100)">
            <rect width="160" height="140" rx="5" fill="#8b5cf6" fill-opacity="0.2" stroke="#8b5cf6" stroke-width="2" />
            <rect width="160" height="30" rx="5" fill="#8b5cf6" stroke="#8b5cf6" stroke-width="2" />
            <text x="80" y="20" text-anchor="middle" fill="#8b5cf6" font-weight="bold">论文表</text>
            <text x="10" y="50" fill="#333333" class="dark:text-gray-200">+ 论文ID (PK)</text>
            <text x="10" y="70" fill="#333333" class="dark:text-gray-200">+ 标题</text>
            <text x="10" y="90" fill="#333333" class="dark:text-gray-200">+ 内容</text>
            <text x="10" y="110" fill="#333333" class="dark:text-gray-200">+ 创建时间</text>
            <text x="10" y="130" fill="#333333" class="dark:text-gray-200">+ 用户ID (FK)</text>
          </g>
          
          <!-- 模板表 -->
          <g transform="translate(-150, 80)">
            <rect width="160" height="120" rx="5" fill="#10b981" fill-opacity="0.2" stroke="#10b981" stroke-width="2" />
            <rect width="160" height="30" rx="5" fill="#10b981" stroke="#10b981" stroke-width="2" />
            <text x="80" y="20" text-anchor="middle" fill="#10b981" font-weight="bold">模板表</text>
            <text x="10" y="50" fill="#333333" class="dark:text-gray-200">+ 模板ID (PK)</text>
            <text x="10" y="70" fill="#333333" class="dark:text-gray-200">+ 名称</text>
            <text x="10" y="90" fill="#333333" class="dark:text-gray-200">+ 描述</text>
            <text x="10" y="110" fill="#333333" class="dark:text-gray-200">+ 配置JSON</text>
          </g>
          
          <!-- 连接线 -->
          <line x1="-120" y1="-90" x2="100" y2="-50" stroke="#3b82f6" stroke-width="2" stroke-dasharray="5,5" />
          <line x1="-70" y1="80" x2="100" y2="-30" stroke="#10b981" stroke-width="2" stroke-dasharray="5,5" />
          
          <!-- 关系标记 -->
          <circle cx="-120" cy="-90" r="5" fill="#3b82f6" />
          <circle cx="100" cy="-50" r="5" fill="#3b82f6" />
          <circle cx="-70" cy="80" r="5" fill="#10b981" />
          <circle cx="100" cy="-30" r="5" fill="#10b981" />
          
          <!-- 关系文本 -->
          <text x="-10" y="-80" text-anchor="middle" fill="#3b82f6" font-weight="bold">创建</text>
          <text x="20" y="30" text-anchor="middle" fill="#10b981" font-weight="bold">使用</text>
        </g>
      </svg>
    </div>
    
    <!-- 悬浮控制按钮 -->
    <div class="absolute top-4 right-4 flex space-x-2">
      <div class="w-8 h-8 rounded-full bg-white dark:bg-gray-700 shadow flex items-center justify-center cursor-pointer">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-600 dark:text-gray-300" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clip-rule="evenodd" />
        </svg>
      </div>
      <div class="w-8 h-8 rounded-full bg-white dark:bg-gray-700 shadow flex items-center justify-center cursor-pointer">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-600 dark:text-gray-300" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
        </svg>
      </div>
    </div>
    
    <!-- 工具栏 -->
    <div class="absolute bottom-0 left-0 right-0 bg-white dark:bg-gray-900 bg-opacity-80 dark:bg-opacity-80 backdrop-blur-sm p-2 flex justify-between items-center">
      <div class="flex space-x-2">
        <div class="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center cursor-pointer">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-white" viewBox="0 0 20 20" fill="currentColor">
            <path d="M5 3a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2V5a2 2 0 00-2-2H5zM5 11a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2v-2a2 2 0 00-2-2H5zM11 5a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V5zM11 13a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
        </div>
        <div class="w-8 h-8 rounded-full bg-purple-500 flex items-center justify-center cursor-pointer">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-white" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="w-8 h-8 rounded-full bg-green-500 flex items-center justify-center cursor-pointer">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-white" viewBox="0 0 20 20" fill="currentColor">
            <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
          </svg>
        </div>
      </div>
      <div class="text-xs text-gray-600 dark:text-gray-300">
        数据模型预览
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'

defineProps({
  height: {
    type: String,
    default: '400px'
  }
})

onMounted(() => {
  // 可以在这里添加交互动画逻辑
})
</script>

<style scoped>
.er-diagram-content {
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0% {
    transform: translateY(0px) scale(1);
  }
  50% {
    transform: translateY(-10px) scale(1.02);
  }
  100% {
    transform: translateY(0px) scale(1);
  }
}
</style> 