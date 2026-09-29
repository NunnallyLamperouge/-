<template>
  <div class="pending-students-container">
    <a-empty v-if="data.length === 0 && !loading" description="暂无待审核申请" />
    
    <a-spin :spinning="loading" tip="加载中...">
      <div class="mb-4 flex justify-end">
        <a-button 
          type="primary" 
          @click="$emit('refresh')"
        >
          <template #icon><reload-outlined /></template>
          刷新
        </a-button>
      </div>
      
      <div class="cards-container grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <a-card 
          v-for="record in data" 
          :key="record.id" 
          class="student-card rounded-lg shadow-md hover:shadow-lg transition-shadow"
          :class="{'border-yellow-400': record.status === 0, 
                  'border-green-400': record.status === 1, 
                  'border-red-400': record.status === 2}"
          :bordered="true"
        >
          <template #cover>
            <div class="p-4 dark:bg-gray-700 py-3 flex justify-center items-center">
              <a-avatar class="bg-blue-500" :size="64">
                <template #icon><user-outlined /></template>
              </a-avatar>
            </div>
          </template>
          
          <a-card-meta :title="record.Student.username">
            <template #description>
              <div class="student-info space-y-2 mt-2">
                <div class="flex items-center">
                  <span class="text-gray-500 inline-block w-16">学院：</span>
                  <span>{{ record.Student.apartment }}</span>
                </div>
                <div class="flex items-center">
                  <span class="text-gray-500 inline-block w-16">专业：</span>
                  <span>{{ record.Student.profession }}</span>
                </div>
                <div class="flex items-center">
                  <span class="text-gray-500 inline-block w-16">申请时间：</span>
                  <span>{{ formatDate(record.CreatedAt) }}</span>
                </div>
                <div class="flex items-center">
                  <span class="text-gray-500 inline-block w-16">状态：</span>
                  <a-tag 
                    :color="record.status === 0 ? 'warning' : record.status === 1 ? 'success' : 'error'"
                  >
                    {{ getStatusText(record.status) }}
                  </a-tag>
                </div>
              </div>
            </template>
          </a-card-meta>
          
          <div class="card-actions mt-4 flex items-center justify-between">
            <a-button type="link" @click="showReasonModal(record)">
              <template #icon><file-text-outlined /></template>
              查看申请理由
            </a-button>
            
            <div class="space-x-2" v-if="record.status === 0">
              <a-button 
                type="primary" 
                size="small"
                @click="updateStatus(record.id, 1)"
                v-motion
                :initial="{ scale: 1 }"
                :hover="{ scale: 1.05 }"
                :loading="loadingId === record.id && actionType === 'approve'"
              >
                <template #icon><check-outlined /></template>
                通过
              </a-button>
              
              <a-button 
                danger
                size="small"
                @click="updateStatus(record.id, 0)"
                v-motion
                :initial="{ scale: 1 }"
                :hover="{ scale: 1.05 }"
                :loading="loadingId === record.id && actionType === 'reject'"
              >
                <template #icon><close-outlined /></template>
                拒绝
              </a-button>
            </div>
          </div>
        </a-card>
      </div>
      
      <!-- 分页 -->
      <div class="mt-4 flex justify-center">
        <a-pagination
          v-model:current="currentPage"
          :total="data.length"
          :pageSize="pageSize"
          @change="handlePageChange"
          show-size-changer
          :pageSizeOptions="['5', '10', '20']"
          @showSizeChange="handleSizeChange"
        />
      </div>
    </a-spin>
    
    <!-- 申请理由模态框 -->
    <a-modal
      v-model:visible="reasonModalVisible"
      title="申请理由"
      :footer="null"
      class="custom-modal"
    >
      <div v-if="selectedApplication" class="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
        <p class="whitespace-pre-line text-gray-700 dark:text-gray-300">{{ selectedApplication.Message || '无申请理由' }}</p>
      </div>
      <div class="flex justify-end mt-4">
        <a-button type="primary" @click="reasonModalVisible = false">关闭</a-button>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits, computed } from 'vue'
import { useMotion } from '@vueuse/motion'
import { 
  UserOutlined,
  CheckOutlined,
  CloseOutlined,
  ReloadOutlined,
  FileTextOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { message } from 'ant-design-vue'
import { putJSON } from '@/utility/request'

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

const emit = defineEmits(['approve', 'reject', 'refresh'])

const reasonModalVisible = ref(false)
const selectedApplication = ref(null)
const currentPage = ref(1)
const pageSize = ref(9)
const loadingId = ref(null)
const actionType = ref('')

// 格式化日期
const formatDate = (date) => {
  if (!date) return '未知'
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 获取申请状态文本
const getStatusText = (status) => {
  const statusMap = {
    0: '待审核',
    1: '已通过',
    2: '已拒绝'
  }
  return statusMap[status] || '未知状态'
}

// 显示申请理由
const showReasonModal = (application) => {
  selectedApplication.value = application
  reasonModalVisible.value = true
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

// 更新指导申请状态
const updateStatus = (id, status) => {
  loadingId.value = id
  actionType.value = status === 1 ? 'approve' : 'reject'
  
  putJSON(
    `/api/guide/${id}/status`, 
    { status: status },
    (msg, data) => {
      // 成功回调
      message.success(status === 1 ? '已通过申请' : '已拒绝申请')
      emit('refresh')
      loadingId.value = null
      actionType.value = ''
    },
    (msg) => {
      // 失败回调
      message.error(msg || '操作失败')
      loadingId.value = null
      actionType.value = ''
    },
    (error) => {
      // 错误回调
      console.error('更新状态失败:', error)
      message.error('操作失败: 服务器错误')
      loadingId.value = null
      actionType.value = ''
    }
  )
}
</script>

<style scoped>
.custom-modal :deep(.ant-modal-content) {
  @apply bg-white dark:bg-gray-800 rounded-xl overflow-hidden border dark:border-gray-700;
}

.custom-modal :deep(.ant-modal-header) {
  @apply bg-white dark:bg-gray-800 border-b-0 pt-5;
}

.custom-modal :deep(.ant-modal-title) {
  @apply text-gray-900 dark:text-white text-lg font-bold;
}

.student-card {
  @apply transition-all duration-300;
}

.student-card:hover {
  @apply transform -translate-y-1;
}
</style> 