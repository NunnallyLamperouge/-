<template>
  <div class="progress-tracker relative rounded-lg overflow-hidden shadow-xl" :style="{ height: height }">
    <div class="absolute inset-0 bg-white dark:bg-gray-900">
      <div class="h-full flex flex-col">
        <!-- 标题栏 -->
        <div class="bg-gray-100 dark:bg-gray-800 px-4 py-3 border-b border-gray-200 dark:border-gray-700">
          <div class="flex justify-between items-center">
            <h3 class="text-sm font-medium text-gray-700 dark:text-gray-300">论文写作进度追踪</h3>
            <div class="flex items-center space-x-2">
              <span class="text-xs text-gray-500 dark:text-gray-400">预计完成日期: 2023-06-15</span>
              <div class="w-2 h-2 rounded-full bg-green-500"></div>
            </div>
          </div>
        </div>
        
        <!-- 内容区域 -->
        <div class="flex-1 p-4 overflow-auto">
          <!-- 进度概览 -->
          <div class="mb-6" v-motion :initial="{ opacity: 0, y: 20 }" :enter="{ opacity: 1, y: 0, transition: { duration: 500 } }">
            <div class="flex justify-between items-center mb-2">
              <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300">总体进度</h4>
              <span class="text-sm font-medium text-blue-600 dark:text-blue-400">68%</span>
            </div>
            <div class="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
              <div class="h-full bg-blue-500 rounded-full" style="width: 68%"></div>
            </div>
          </div>
          
          <!-- 热力图 -->
          <div class="mb-6" v-motion :initial="{ opacity: 0, y: 20 }" :enter="{ opacity: 1, y: 0, transition: { delay: 200, duration: 500 } }">
            <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">写作热力图</h4>
            <div class="heatmap-container">
              <div class="heatmap-labels">
                <div class="heatmap-month">一月</div>
                <div class="heatmap-month">二月</div>
                <div class="heatmap-month">三月</div>
                <div class="heatmap-month">四月</div>
                <div class="heatmap-month">五月</div>
                <div class="heatmap-month">六月</div>
              </div>
              <div class="heatmap-grid">
                <div v-for="i in 180" :key="i" class="heatmap-cell" :class="getHeatmapClass(i)"></div>
              </div>
              <div class="heatmap-legend">
                <div class="legend-item">
                  <div class="legend-color level-0"></div>
                  <span>无活动</span>
                </div>
                <div class="legend-item">
                  <div class="legend-color level-1"></div>
                  <span>低</span>
                </div>
                <div class="legend-item">
                  <div class="legend-color level-2"></div>
                  <span>中</span>
                </div>
                <div class="legend-item">
                  <div class="legend-color level-3"></div>
                  <span>高</span>
                </div>
                <div class="legend-item">
                  <div class="legend-color level-4"></div>
                  <span>极高</span>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 章节进度 -->
          <div class="mb-6" v-motion :initial="{ opacity: 0, y: 20 }" :enter="{ opacity: 1, y: 0, transition: { delay: 400, duration: 500 } }">
            <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">章节完成度</h4>
            <div class="space-y-3">
              <div class="chapter-progress">
                <div class="chapter-info">
                  <div class="chapter-name">第一章：绪论</div>
                  <div class="chapter-percent">100%</div>
                </div>
                <div class="progress-bar">
                  <div class="progress-fill" style="width: 100%"></div>
                </div>
              </div>
              
              <div class="chapter-progress">
                <div class="chapter-info">
                  <div class="chapter-name">第二章：相关工作</div>
                  <div class="chapter-percent">85%</div>
                </div>
                <div class="progress-bar">
                  <div class="progress-fill" style="width: 85%"></div>
                </div>
              </div>
              
              <div class="chapter-progress">
                <div class="chapter-info">
                  <div class="chapter-name">第三章：系统设计</div>
                  <div class="chapter-percent">70%</div>
                </div>
                <div class="progress-bar">
                  <div class="progress-fill" style="width: 70%"></div>
                </div>
              </div>
              
              <div class="chapter-progress">
                <div class="chapter-info">
                  <div class="chapter-name">第四章：实验与结果</div>
                  <div class="chapter-percent">45%</div>
                </div>
                <div class="progress-bar">
                  <div class="progress-fill" style="width: 45%"></div>
                </div>
              </div>
              
              <div class="chapter-progress">
                <div class="chapter-info">
                  <div class="chapter-name">第五章：结论与展望</div>
                  <div class="chapter-percent">20%</div>
                </div>
                <div class="progress-bar">
                  <div class="progress-fill" style="width: 20%"></div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 里程碑 -->
          <div v-motion :initial="{ opacity: 0, y: 20 }" :enter="{ opacity: 1, y: 0, transition: { delay: 600, duration: 500 } }">
            <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">里程碑</h4>
            <div class="milestones">
              <div class="milestone completed">
                <div class="milestone-dot"></div>
                <div class="milestone-line"></div>
                <div class="milestone-content">
                  <div class="milestone-date">2023-01-15</div>
                  <div class="milestone-title">开题报告提交</div>
                  <div class="milestone-status">已完成</div>
                </div>
              </div>
              
              <div class="milestone completed">
                <div class="milestone-dot"></div>
                <div class="milestone-line"></div>
                <div class="milestone-content">
                  <div class="milestone-date">2023-03-01</div>
                  <div class="milestone-title">中期检查</div>
                  <div class="milestone-status">已完成</div>
                </div>
              </div>
              
              <div class="milestone current">
                <div class="milestone-dot"></div>
                <div class="milestone-line"></div>
                <div class="milestone-content">
                  <div class="milestone-date">2023-05-20</div>
                  <div class="milestone-title">论文初稿完成</div>
                  <div class="milestone-status">进行中</div>
                </div>
              </div>
              
              <div class="milestone">
                <div class="milestone-dot"></div>
                <div class="milestone-line"></div>
                <div class="milestone-content">
                  <div class="milestone-date">2023-06-10</div>
                  <div class="milestone-title">论文定稿</div>
                  <div class="milestone-status">未开始</div>
                </div>
              </div>
              
              <div class="milestone">
                <div class="milestone-dot"></div>
                <div class="milestone-content">
                  <div class="milestone-date">2023-06-25</div>
                  <div class="milestone-title">答辩</div>
                  <div class="milestone-status">未开始</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  height: {
    type: String,
    default: '500px'
  }
})

// 模拟热力图数据
const getHeatmapClass = (index) => {
  // 这里可以根据实际数据生成热力图
  const random = Math.random();
  if (random < 0.5) return 'level-0';
  if (random < 0.7) return 'level-1';
  if (random < 0.85) return 'level-2';
  if (random < 0.95) return 'level-3';
  return 'level-4';
}
</script>

<style scoped lang="postcss">
/* 热力图样式 */
.heatmap-container {
  @apply bg-white dark:bg-gray-800 rounded-lg p-3 border border-gray-200 dark:border-gray-700;
}

.heatmap-labels {
  @apply flex text-xs text-gray-500 dark:text-gray-400 mb-2;
}

.heatmap-month {
  @apply flex-1 text-center;
}

.heatmap-grid {
  @apply grid gap-1;
  grid-template-columns: repeat(30, minmax(0, 1fr));
}

.heatmap-cell {
  @apply w-2 h-2 rounded-sm;
}

.level-0 {
  @apply bg-gray-200 dark:bg-gray-700;
}

.level-1 {
  @apply bg-green-200 dark:bg-green-900;
}

.level-2 {
  @apply bg-green-300 dark:bg-green-700;
}

.level-3 {
  @apply bg-green-500 dark:bg-green-600;
}

.level-4 {
  @apply bg-green-700 dark:bg-green-500;
}

.heatmap-legend {
  @apply flex justify-end mt-2 space-x-3 text-xs text-gray-500 dark:text-gray-400;
}

.legend-item {
  @apply flex items-center;
}

.legend-color {
  @apply w-2 h-2 rounded-sm mr-1;
}

/* 章节进度样式 */
.chapter-progress {
  @apply bg-white dark:bg-gray-800 rounded-lg p-3 border border-gray-200 dark:border-gray-700;
}

.chapter-info {
  @apply flex justify-between items-center mb-2;
}

.chapter-name {
  @apply text-xs font-medium text-gray-700 dark:text-gray-300;
}

.chapter-percent {
  @apply text-xs font-medium text-blue-600 dark:text-blue-400;
}

.progress-bar {
  @apply h-1.5 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden;
}

.progress-fill {
  @apply h-full bg-blue-500 rounded-full;
}

/* 里程碑样式 */
.milestones {
  @apply relative pl-6;
}

.milestone {
  @apply relative mb-6 last:mb-0;
}

.milestone-dot {
  @apply absolute left-0 w-4 h-4 rounded-full bg-gray-200 dark:bg-gray-700 border-2 border-white dark:border-gray-900 -translate-x-1/2;
}

.milestone.completed .milestone-dot {
  @apply bg-green-500;
}

.milestone.current .milestone-dot {
  @apply bg-blue-500;
}

.milestone-line {
  @apply absolute left-0 top-4 bottom-0 w-0.5 bg-gray-200 dark:bg-gray-700 -translate-x-1/2;
}

.milestone:last-child .milestone-line {
  @apply hidden;
}

.milestone-content {
  @apply ml-4;
}

.milestone-date {
  @apply text-xs text-gray-500 dark:text-gray-400;
}

.milestone-title {
  @apply text-sm font-medium text-gray-700 dark:text-gray-300;
}

.milestone-status {
  @apply text-xs;
}

.milestone.completed .milestone-status {
  @apply text-green-600 dark:text-green-400;
}

.milestone.current .milestone-status {
  @apply text-blue-600 dark:text-blue-400;
}

.milestone:not(.completed):not(.current) .milestone-status {
  @apply text-gray-500 dark:text-gray-400;
}
</style> 