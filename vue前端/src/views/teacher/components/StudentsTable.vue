<template>
  <div class="students-cards-container">
    <a-empty v-if="data.length === 0 && !loading" description="暂无指导学生" />
    <a-spin :spinning="loading" tip="加载中...">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div 
          v-for="student in data" 
          :key="student.id" 
          class="student-card bg-white dark:bg-gray-800 rounded-xl shadow-md hover:shadow-lg transition-all duration-300 overflow-hidden"
          v-motion
          :initial="{ opacity: 0, y: 20 }"
          :enter="{ opacity: 1, y: 0, transition: { duration: 300 } }"
          @click="showStudentDetail(student)"
        >
          <!-- 卡片头部 -->
          <div class="card-header p-4 flex items-center space-x-3 border-b dark:border-gray-700">
            <a-avatar :src="student.avatar" v-if="student.avatar" :size="50">
              <template #icon><img src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png"/></template>
            </a-avatar>
            <a-avatar v-else class="bg-blue-500" :size="50">
              <template #icon><user-outlined /></template>
            </a-avatar>
            <div class="flex-1 min-w-0">
              <h3 class="text-lg font-bold text-gray-800 dark:text-white truncate">
                {{ student.username }}
              </h3>
              <p class="text-sm text-gray-500 dark:text-gray-400 truncate">
                {{ student.studentId || '暂无学号' }}
              </p>
            </div>
            
            <a-tag :color="getStatusColor(student.status)" class="whitespace-nowrap">
              {{ getStatusText(student.status) }}
            </a-tag>
          </div>
          
          <!-- 卡片内容 -->
          <div class="card-body p-4 space-y-3">
            <div class="grid grid-cols-2 gap-2">
              <div class="text-sm">
                <div class="text-gray-500 dark:text-gray-400">学院</div>
                <div class="font-medium">{{ student.faculty || '未设置' }}</div>
              </div>
              <div class="text-sm">
                <div class="text-gray-500 dark:text-gray-400">专业</div>
                <div class="font-medium">{{ student.major || '未设置' }}</div>
              </div>
              <div class="text-sm col-span-2">
                <div class="text-gray-500 dark:text-gray-400">开始指导时间</div>
                <div class="font-medium">{{ formatDate(student.createdAt) }}</div>
              </div>
            </div>
            
            <div class="progress-section pt-2">
              <div class="flex justify-between items-center mb-1">
                <span class="text-sm text-gray-500 dark:text-gray-400">完成进度</span>
                <span class="text-sm font-medium">{{ student.completionRate || 0 }}%</span>
              </div>
              <a-progress 
                :percent="student.completionRate || 0" 
                :status="getProgressStatus(student.completionRate || 0)"
                :stroke-width="8"
              />
            </div>
          </div>
          
          <!-- 卡片底部 -->
          <div class="card-footer p-4 border-t dark:border-gray-700 flex space-x-2 justify-between">
            <a-button 
              type="primary" 
              @click.stop="$emit('create-mission', student.id)"
              v-motion
              :initial="{ scale: 1 }"
              :hover="{ scale: 1.05 }"
            >
              <template #icon><file-add-outlined /></template>
              创建任务
            </a-button>
            
            <a-popconfirm
              title="确定解除与该学生的指导关系吗?"
              @confirm.stop="confirmRemove(student.id)"
              ok-text="确定"
              cancel-text="取消"
            >
              <a-button 
                type="danger"
                @click.stop
              >
                <template #icon><disconnect-outlined /></template>
                解除关系
              </a-button>
            </a-popconfirm>
          </div>
        </div>
      </div>
      
      <!-- 分页 -->
      <div class="mt-6 flex justify-center">
        <a-pagination
          v-model:current="currentPage"
          :total="data.length"
          :pageSize="pageSize"
          @change="handlePageChange"
          show-size-changer
          :pageSizeOptions="['6', '12', '24']"
          @showSizeChange="handleSizeChange"
        />
      </div>
    </a-spin>
    
    <!-- 学生详情抽屉 -->
    <a-drawer
      :title="`学生详情 - ${selectedStudent?.username || ''}`"
      :visible="drawerVisible"
      @close="drawerVisible = false"
      width="500"
      placement="right"
    >
      <template v-if="selectedStudent">
        <div class="flex items-center space-x-4 mb-6">
          <a-avatar :size="64" :src="selectedStudent.avatar" v-if="selectedStudent.avatar">
            <template #icon><img src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png"/></template>
          </a-avatar>
          <a-avatar :size="64" v-else class="bg-blue-500">
            <template #icon><user-outlined /></template>
          </a-avatar>
          
          <div>
            <h2 class="text-xl font-bold">{{ selectedStudent.username }}</h2>
            <p class="text-gray-500">{{ selectedStudent.studentId || '暂无学号' }}</p>
          </div>
        </div>
        
        <a-descriptions bordered :column="1">
          <a-descriptions-item label="邮箱">
            {{ selectedStudent.email || '未设置' }}
          </a-descriptions-item>
          <a-descriptions-item label="学院">
            {{ selectedStudent.faculty || '未设置' }}
          </a-descriptions-item>
          <a-descriptions-item label="专业">
            {{ selectedStudent.major || '未设置' }}
          </a-descriptions-item>
          <a-descriptions-item label="联系电话">
            {{ selectedStudent.phone || '未设置' }}
          </a-descriptions-item>
          <a-descriptions-item label="指导开始时间">
            {{ formatDate(selectedStudent.guidanceStartDate) }}
          </a-descriptions-item>
        </a-descriptions>
        
        <div class="mt-6">
          <h3 class="text-lg font-bold mb-3">任务完成情况</h3>
          <a-progress 
            :percent="selectedStudent.completionRate || 0" 
            :status="getProgressStatus(selectedStudent.completionRate || 0)"
          />
        </div>
      </template>
    </a-drawer>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue'
import { useMotion } from '@vueuse/motion'
import { message, Modal } from 'ant-design-vue'
import { 
  UserOutlined,
  EyeOutlined,
  DisconnectOutlined,
  FileAddOutlined,
  ReloadOutlined,
  FileTextOutlined
} from '@ant-design/icons-vue'
import { del, get } from '../../../utility/request'
import dayjs from 'dayjs'
import { useRouter } from 'vue-router'

defineProps({
  data: {
    type: Array,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['refresh', 'create-mission'])

const drawerVisible = ref(false)
const selectedStudent = ref(null)
const router = useRouter()
const currentPage = ref(1)
const pageSize = ref(6)

// 获取状态颜色
const getStatusColor = (status) => {
  const statusMap = {
    active: 'green',
    warning: 'orange',
    danger: 'red',
    inactive: 'gray'
  }
  return statusMap[status] || 'blue'
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    active: '正常',
    warning: '需关注',
    danger: '逾期',
    inactive: '暂停'
  }
  return statusMap[status] || '正常'
}

// 获取进度状态
const getProgressStatus = (percent) => {
  if (percent >= 100) return 'success'
  if (percent >= 80) return 'active'
  if (percent >= 0) return 'normal'
  return 'exception'
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return '未知'
  return dayjs(date).format('YYYY-MM-DD')
}

// 页码变更
const handlePageChange = (page) => {
  currentPage.value = page
}

// 每页数量变更
const handleSizeChange = (current, size) => {
  pageSize.value = size
  currentPage.value = 1
}

// 显示学生详情
const showStudentDetail = (student) => {
  selectedStudent.value = student
  drawerVisible.value = true
}

// 确认移除学生
const confirmRemove = (studentId) => {
  del(`/api/guide/${studentId}`, null,
    (msg) => {
      message.success('已解除指导关系')
      // 刷新列表
      this.$emit('refresh')
    },
    (msg) => {
      message.error(msg || '操作失败')
    }
  )
}

// 查看学生论文
const viewStudentThesis = (student) => {
  // 实际应用中需要获取学生的论文模板ID
  // 这里我们先查询学生的所有论文模板，然后选择最新的一个进行查看
  get(`/api/template/student/${student.id}`, null,
    (msg, data) => {
      if (data && data.length > 0) {
        // 假设我们选择最新的一个模板
        const latestTemplate = data[0]
        // 导航到论文查看页面
        router.push(`/teacher/view-template/${latestTemplate.id}`)
      } else {
        message.warning('该学生尚未创建任何论文')
      }
    },
    (errorMsg) => {
      message.error(errorMsg || '获取学生论文失败')
    }
  )
}
</script> 