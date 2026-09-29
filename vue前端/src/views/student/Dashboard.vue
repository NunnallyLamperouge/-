<template>
  <div class="py-6 bg-app theme-transition min-h-screen">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
      <!-- 快捷操作卡片 -->
      <div 
        v-motion
        :initial="{ opacity: 0, y: 50 }"
        :enter="{ opacity: 1, y: 0, transition: { duration: 500 } }"
        class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3"
      >
        <!-- 继续写作 -->
        <div 
          class="card hover:shadow-lg transform hover:-translate-y-1 transition-all duration-300"
          v-motion
          :initial="{ scale: 0.9, opacity: 0 }"
          :enter="{ scale: 1, opacity: 1, transition: { delay: 200 } }"
        >
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <edit-outlined class="h-6 w-6 text-blue-600" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-muted truncate">
                    继续写作
                  </dt>
                  <dd class="flex items-baseline">
                    <div class="text-2xl font-semibold text-app">
                      第三章
                    </div>
                    <div class="ml-2 flex items-baseline text-sm font-semibold text-green-600">
                      <arrow-up-outlined class="self-center flex-shrink-0 h-5 w-5" />
                      <span class="sr-only">增加</span>
                      2.5%
                    </div>
                  </dd>
                </dl>
              </div>
            </div>
          </div>
          <div class="bg-sunken px-5 py-3">
            <div class="text-sm">
              <router-link
                to="/student/editor"
                class="font-medium text-blue-600 hover:text-blue-500"
              >
                继续编辑
              </router-link>
            </div>
          </div>
        </div>

        <!-- 模拟答辩 -->
        <div 
          class="card hover:shadow-lg transform hover:-translate-y-1 transition-all duration-300"
          v-motion
          :initial="{ scale: 0.9, opacity: 0 }"
          :enter="{ scale: 1, opacity: 1, transition: { delay: 400 } }"
        >
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <team-outlined class="h-6 w-6 text-green-600" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-muted truncate">
                    模拟答辩
                  </dt>
                  <dd class="flex items-baseline">
                    <div class="text-2xl font-semibold text-app">
                      3 次
                    </div>
                    <div class="ml-2 flex items-baseline text-sm font-semibold text-green-600">
                      已完成
                    </div>
                  </dd>
                </dl>
              </div>
            </div>
          </div>
          <div class="bg-sunken px-5 py-3">
            <div class="text-sm">
              <router-link
                to="/student/defense"
                class="font-medium text-blue-600 hover:text-blue-500"
              >
                开始练习
              </router-link>
            </div>
          </div>
        </div>

        <!-- 格式检查 -->
        <div 
          class="card hover:shadow-lg transform hover:-translate-y-1 transition-all duration-300"
          v-motion
          :initial="{ scale: 0.9, opacity: 0 }"
          :enter="{ scale: 1, opacity: 1, transition: { delay: 600 } }"
        >
          <div class="p-5">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <check-circle-outlined class="h-6 w-6 text-yellow-600" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-muted truncate">
                    格式检查
                  </dt>
                  <dd class="flex items-baseline">
                    <div class="text-2xl font-semibold text-app">
                      5 项
                    </div>
                    <div class="ml-2 flex items-baseline text-sm font-semibold text-yellow-600">
                      待处理
                    </div>
                  </dd>
                </dl>
              </div>
            </div>
          </div>
          <div class="bg-sunken px-5 py-3">
            <div class="text-sm">
              <router-link
                to="/student/format"
                class="font-medium text-blue-600 hover:text-blue-500"
              >
                查看详情
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- 数据可视化 -->
      <div 
        class="mt-8 grid grid-cols-1 gap-6 lg:grid-cols-2"
        v-motion
        :initial="{ opacity: 0, x: -50 }"
        :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 800 } }"
      >
        <!-- 模板编辑次数热力图 -->
        <div class="card hover:shadow-xl transition-shadow duration-300">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex justify-between items-center mb-4">
              <h3 class="text-lg leading-6 font-medium text-app">
                模板编辑热力图
              </h3>
              <a-select 
                v-model:value="selectedTimeRange" 
                style="width: 120px" 
                @change="fetchDashboardData"
              >
                <a-select-option value="7">近7天</a-select-option>
                <a-select-option value="30">近30天</a-select-option>
                <a-select-option value="90">近90天</a-select-option>
                <a-select-option value="365">近1年</a-select-option>
              </a-select>
            </div>
            <div class="h-72" ref="heatmapRef">
              <div v-if="isLoading" class="flex h-full items-center justify-center">
                <a-spin />
              </div>
            </div>
          </div>
        </div>

        <!-- 任务完成率和模板类型分布 -->
        <div class="card hover:shadow-xl transition-shadow duration-300">
          <div class="px-4 py-5 sm:p-6">
            <h3 class="text-lg leading-6 font-medium text-app mb-4">
              模板类型分布
            </h3>
            <div class="h-72" ref="pieChartRef">
              <div v-if="isLoading" class="flex h-full items-center justify-center">
                <a-spin />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 写作进度 -->
      <div 
        class="mt-8"
        v-motion
        :initial="{ opacity: 0, x: 50 }"
        :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 1000 } }"
      >
        <div class="card hover:shadow-xl transition-shadow duration-300">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex justify-between items-center">
              <h3 class="text-lg leading-6 font-medium text-app">
                写作进度
              </h3>
              <div class="text-app font-medium">
                总体进度：{{ Math.round(taskCompletionRate) }}%
              </div>
            </div>
            <div class="mt-4">
              <div class="flex items-center justify-between mb-2">
                <div class="text-sm font-medium text-app">第一章 绪论</div>
                <div class="text-sm text-muted">100%</div>
              </div>
              <a-progress :percent="100" status="success" />

              <div class="flex items-center justify-between mb-2 mt-4">
                <div class="text-sm font-medium text-app">第二章 文献综述</div>
                <div class="text-sm text-muted">85%</div>
              </div>
              <a-progress :percent="85" />

              <div class="flex items-center justify-between mb-2 mt-4">
                <div class="text-sm font-medium text-app">第三章 研究方法</div>
                <div class="text-sm text-muted">45%</div>
              </div>
              <a-progress :percent="45" status="active" />
            </div>
          </div>
        </div>
      </div>

      <!-- 通知统计 -->
      <div 
        class="mt-8"
        v-motion
        :initial="{ opacity: 0, y: 50 }"
        :enter="{ opacity: 1, y: 0, transition: { duration: 500, delay: 1200 } }"
      >
        <div class="card hover:shadow-xl transition-shadow duration-300">
          <div class="px-4 py-5 sm:p-6">
            <h3 class="text-lg leading-6 font-medium text-app mb-4">
              通知统计
            </h3>
            <div class="h-64" ref="barChartRef">
              <div v-if="isLoading" class="flex h-full items-center justify-center">
                <a-spin />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 最近活动 -->
      <div 
        class="mt-8 mb-8"
        v-motion
        :initial="{ opacity: 0, y: 50 }"
        :enter="{ opacity: 1, y: 0, transition: { duration: 500, delay: 1400 } }"
      >
        <div class="card hover:shadow-xl transition-shadow duration-300">
          <div class="px-4 py-5 sm:p-6">
            <h3 class="text-lg leading-6 font-medium text-app mb-4">
              最近活动
            </h3>
            <a-timeline>
              <a-timeline-item>
                <template #dot>
                  <edit-outlined style="font-size: 16px;" />
                </template>
                <p class="text-app">更新了第三章研究方法</p>
                <p class="text-sm text-muted">2小时前</p>
              </a-timeline-item>
              <a-timeline-item>
                <template #dot>
                  <message-outlined style="font-size: 16px;" />
                </template>
                <p class="text-app">收到导师评审意见</p>
                <p class="text-sm text-muted">昨天 15:30</p>
              </a-timeline-item>
              <a-timeline-item>
                <template #dot>
                  <check-circle-outlined style="font-size: 16px;" />
                </template>
                <p class="text-app">完成格式自动修正</p>
                <p class="text-sm text-muted">昨天 10:20</p>
              </a-timeline-item>
            </a-timeline>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useMotion } from '@vueuse/motion'
import {
  EditOutlined,
  TeamOutlined,
  CheckCircleOutlined,
  MessageOutlined,
  ArrowUpOutlined
} from '@ant-design/icons-vue'
import { get } from '../../utility/request'
import * as echarts from 'echarts'

// 引用和状态
const heatmapRef = ref(null)
const pieChartRef = ref(null)
const barChartRef = ref(null)
const isLoading = ref(true)
const selectedTimeRange = ref('30')
const taskCompletionRate = ref(0)

// echarts 实例
let heatmapChart = null
let pieChart = null
let barChart = null

// 仪表盘数据
const dashboardData = ref({
  templateEditHeatmap: [],
  taskCompletionRate: 0,
  notificationStatistics: [],
  templateTypeDistribution: []
})

// 获取仪表盘数据
const fetchDashboardData = () => {
  isLoading.value = true
  
  get('/api/dashboard/student', { timeRange: selectedTimeRange.value }, 
    (message, data) => {
      dashboardData.value = data
      taskCompletionRate.value = data.taskCompletionRate
      isLoading.value = false
      
      // 更新图表
      renderHeatmapChart()
      renderPieChart()
      renderBarChart()
    },
    (errorMsg) => {
      console.error('获取仪表盘数据失败:', errorMsg)
      isLoading.value = false
    }
  )
}

// 渲染模板编辑热力图
const renderHeatmapChart = () => {
  if (!heatmapChart || !dashboardData.value.templateEditHeatmap) return
  
  // 准备数据
  const heatmapData = dashboardData.value.templateEditHeatmap.map(item => {
    return [item.date, item.count]
  })
  
  // 计算日期范围
  const dateSet = new Set(dashboardData.value.templateEditHeatmap.map(item => item.date))
  const dateList = Array.from(dateSet).sort()
  
  // 获取最大值，用于设置颜色梯度
  const maxCount = Math.max(...dashboardData.value.templateEditHeatmap.map(item => item.count), 1)
  
  // 设置GitHub风格的颜色
  const getGithubColor = (value) => {
    // 无编辑时为浅灰色
    if (value === 0) return '#ebedf0'
    
    // GitHub风格的绿色梯度
    const ratio = value / maxCount
    if (ratio <= 0.25) return '#9be9a8'
    if (ratio <= 0.5) return '#40c463'
    if (ratio <= 0.75) return '#30a14e'
    return '#216e39'
  }
  
  // 设置配置项
  const option = {
    tooltip: {
      position: 'top',
      formatter: function (params) {
        const count = params.data[1]
        const date = params.data[0]
        if (count === 0) {
          return `${date}: 无编辑记录`
        }
        return `${date}: ${count} 次编辑`
      },
      backgroundColor: 'rgba(50, 50, 50, 0.9)',
      borderColor: '#666',
      textStyle: {
        color: '#fff'
      }
    },
    visualMap: {
      min: 0,
      max: maxCount,
      calculable: false,
      orient: 'horizontal',
      left: 'center',
      bottom: '0%',
      show: false, // 隐藏图例
      inRange: {
        color: ['#ebedf0', '#9be9a8', '#40c463', '#30a14e', '#216e39'] // GitHub风格的颜色
      }
    },
    calendar: {
      top: 30,
      left: 40, 
      right: 20,
      cellSize: ['auto', 16],
      range: dateList.length > 0 ? [dateList[0], dateList[dateList.length - 1]] : '',
      itemStyle: {
        borderWidth: 1,
        borderColor: '#fff'
      },
      splitLine: {
        show: true,
        lineStyle: {
          color: '#f0f0f0',
          width: 2,
          type: 'solid'
        }
      },
      yearLabel: { show: false },
      dayLabel: {
        firstDay: 1,
        nameMap: ['日', '一', '二', '三', '四', '五', '六'],
        color: 'var(--color-text-muted)'
      },
      monthLabel: {
        nameMap: 'cn',
        color: 'var(--color-text-muted)',
        fontSize: 12
      }
    },
    series: {
      type: 'heatmap',
      coordinateSystem: 'calendar',
      data: heatmapData,
      label: {
        show: false
      },
      itemStyle: {
        color: function (params) {
          return getGithubColor(params.data[1])
        },
        borderColor: '#fff',
        borderWidth: 1
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 4,
          shadowColor: 'rgba(0, 0, 0, 0.2)'
        }
      }
    }
  }
  
  heatmapChart.setOption(option)
}

// 渲染模板类型饼图
const renderPieChart = () => {
  if (!pieChart || !dashboardData.value.templateTypeDistribution) return
  
  // 准备数据
  const pieData = dashboardData.value.templateTypeDistribution.map(item => {
    return { value: item.count, name: item.type }
  })
  
  // 设置配置项
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      textStyle: {
        color: 'var(--color-text)'
      }
    },
    series: [
      {
        name: '模板类型',
        type: 'pie',
        radius: '65%',
        center: ['60%', '50%'],
        data: pieData,
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        label: {
          color: 'var(--color-text)'
        }
      }
    ]
  }
  
  pieChart.setOption(option)
}

// 渲染通知统计柱状图
const renderBarChart = () => {
  if (!barChart || !dashboardData.value.notificationStatistics) return
  
  // 准备数据
  const categories = dashboardData.value.notificationStatistics.map(item => item.type)
  const values = dashboardData.value.notificationStatistics.map(item => item.count)
  
  // 设置配置项
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: categories,
      axisLabel: {
        color: 'var(--color-text)',
        rotate: 45,
        interval: 0
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: 'var(--color-text)'
      }
    },
    series: [
      {
        name: '数量',
        type: 'bar',
        data: values,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#6366f1' },
            { offset: 1, color: '#8b5cf6' }
          ])
        },
        emphasis: {
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#8b5cf6' },
              { offset: 1, color: '#6366f1' }
            ])
          }
        }
      }
    ]
  }
  
  barChart.setOption(option)
}

// 优化图表主题
const getChartTheme = () => {
  return {
    backgroundColor: 'transparent',
    textStyle: {
      fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial',
      color: 'var(--color-text)'
    },
    grid: {
      borderColor: 'rgba(128, 128, 128, 0.1)'
    },
    legend: {
      textStyle: {
        color: 'var(--color-text)'
      }
    }
  }
}

// 初始化图表
const initCharts = () => {
  if (heatmapRef.value) {
    heatmapChart = echarts.init(heatmapRef.value, getChartTheme())
    window.addEventListener('resize', () => heatmapChart.resize())
  }
  
  if (pieChartRef.value) {
    pieChart = echarts.init(pieChartRef.value, getChartTheme())
    window.addEventListener('resize', () => pieChart.resize())
  }
  
  if (barChartRef.value) {
    barChart = echarts.init(barChartRef.value, getChartTheme())
    window.addEventListener('resize', () => barChart.resize())
  }
  
  fetchDashboardData()
}

// 组件挂载
onMounted(() => {
  initCharts()
})

// 组件销毁
onUnmounted(() => {
  if (heatmapChart) {
    heatmapChart.dispose()
    window.removeEventListener('resize', () => heatmapChart.resize())
  }
  
  if (pieChart) {
    pieChart.dispose()
    window.removeEventListener('resize', () => pieChart.resize())
  }
  
  if (barChart) {
    barChart.dispose()
    window.removeEventListener('resize', () => barChart.resize())
  }
})
</script>

<style scoped>
.card {
  @apply bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden border border-gray-200 dark:border-gray-700;
  backdrop-filter: blur(10px);
}

.text-app {
  @apply text-gray-900 dark:text-gray-100;
}

.text-muted {
  @apply text-gray-500 dark:text-gray-400;
}

.bg-sunken {
  @apply bg-gray-50 dark:bg-gray-900/50;
}

/* 优化暗色模式下的卡片效果 */
:deep(.dark) .card {
  @apply shadow-xl;
  background: rgba(17, 24, 39, 0.8);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

/* 优化暗色模式下的图标颜色 */
:deep(.dark) .text-blue-600 {
  @apply text-indigo-400;
}

:deep(.dark) .text-green-600 {
  @apply text-emerald-400;
}

:deep(.dark) .text-yellow-600 {
  @apply text-amber-400;
}

/* 优化暗色模式下的链接颜色 */
:deep(.dark) a.text-blue-600 {
  @apply text-indigo-400 hover:text-indigo-300;
}

/* 优化暗色模式下的进度条 */
:deep(.dark) .ant-progress-bg {
  @apply bg-gradient-to-r from-indigo-500 to-purple-500;
}

:deep(.dark) .ant-progress-success-bg {
  @apply bg-gradient-to-r from-emerald-500 to-teal-500;
}

/* 优化暗色模式下的时间线 */
:deep(.dark) .ant-timeline-item-head {
  @apply bg-gray-800;
}

:deep(.dark) .ant-timeline-item-tail {
  @apply border-l-gray-700;
}

/* 添加过渡效果 */
.theme-transition {
  @apply transition-colors duration-300;
}

/* 优化选择器样式 */
:deep(.ant-select) {
  @apply rounded-md;
}

:deep(.ant-select-selector) {
  @apply border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100;
}

/* 优化加载动画 */
:deep(.ant-spin-dot-item) {
  @apply bg-indigo-500 dark:bg-indigo-400;
}
</style> 