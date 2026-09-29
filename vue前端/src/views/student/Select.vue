<template>
  <div class="teacher-selection-container px-6 py-8 max-w-7xl mx-auto bg-elevated dark:bg-gray-900 rounded-lg shadow-md">
    <div class="page-header flex justify-between items-center mb-6" 
      v-motion
      :initial="{ opacity: 0, y: -20 }"
      :enter="{ opacity: 1, y: 0 }">
      <h1 class="text-2xl font-bold text-gray-800 dark:text-white">选择指导老师</h1>
      <a-input
        v-model:value="searchQuery"
        placeholder="搜索老师姓名或专业"
        class="w-72"
        @input="filterTeachers"
        allow-clear
      >
        <template #prefix>
          <search-outlined />
        </template>
      </a-input>
    </div>
    
    <a-alert
      v-if="successMessage"
      type="success"
      class="mb-6"
      show-icon
      :message="'申请成功'"
      :description="successMessage"
      :closable="true"
      @close="successMessage = ''"
      v-motion
      :initial="{ opacity: 0, scale: 0.9 }"
      :enter="{ opacity: 1, scale: 1 }"
    />
    
    <a-alert
      v-if="errorMessage"
      type="error"
      class="mb-6"
      show-icon
      :message="'操作失败'"
      :description="errorMessage"
      :closable="true" 
      @close="errorMessage = ''"
      v-motion
      :initial="{ opacity: 0, scale: 0.9 }"
      :enter="{ opacity: 1, scale: 1 }"
    />
    
    <div class="teachers-grid" v-motion-slide-visible>
      <a-empty v-if="filteredTeachers.length === 0 && !loading" description="没有找到符合条件的老师" class="mt-12" />
      
      <a-spin :spinning="loading" tip="加载中..." class="mt-8">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <a-card 
            v-for="(teacher, index) in filteredTeachers" 
            :key="teacher.id" 
            class="card hover:-translate-y-2 transition-all duration-300 overflow-hidden"
            :bodyStyle="{ padding: '20px' }"
            v-motion
            :initial="{ opacity: 0, y: 50 }"
            :enter="{ opacity: 1, y: 0, transition: { delay: index * 100, duration: 500 } }"
          >
            <div class="flex items-start space-x-4">
              <div class="relative">
                <a-avatar :size="64" :src="teacher.avatar" v-if="teacher.avatar">
                  <template #icon><img src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png"/></template>
                </a-avatar>
                <a-avatar :size="64" v-else class="flex items-center justify-center bg-blue-500">
                  <span class="text-2xl text-white">{{ teacher.username ? teacher.username[0] : 'T' }}</span>
                </a-avatar>
                <div class="absolute -bottom-1 -right-1 w-5 h-5 bg-green-500 rounded-full border-2 border-white dark:border-gray-800"></div>
              </div>
              
              <div class="flex-1">
                <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-1">{{ teacher.username }}</h2>
                <p class="text-sm text-gray-600 dark:text-gray-300">{{ teacher.title || '讲师' }}</p>
                <p class="text-sm text-gray-600 dark:text-gray-300">{{ teacher.apartment || '暂无院系信息' }}</p>
                <div class="mt-3 py-2 px-3 bg-gray-100 dark:bg-gray-800 rounded-lg">
                  <p class="text-xs text-gray-500 dark:text-gray-400">研究方向</p>
                  <p class="text-sm font-medium">{{ teacher.profession || '暂无研究方向信息' }}</p>
                </div>

                <a-button 
                  type="primary" 
                  block
                  class="mt-4"
                  @click="applyForGuide(teacher)"
                >
                  <template #icon><user-add-outlined /></template>
                  申请指导
                </a-button>
              </div>
            </div>
          </a-card>
        </div>
      </a-spin>
    </div>
    
    <!-- 申请确认对话框 -->
    <a-modal
      :title="'申请指导确认'"
      v-model:visible="dialogVisible"
      :footer="null"
      class="custom-modal"
      :maskClosable="false"
      :width="420"
    >
      <div class="p-2" v-motion :initial="{ opacity: 0 }" :enter="{ opacity: 1 }">
        <p class="text-gray-700 dark:text-gray-200 mb-4">您确定要申请 <span class="font-bold text-blue-600 dark:text-blue-400">{{ selectedTeacher?.username }}</span> 作为您的指导老师吗？</p>
        
        <a-form :model="applyForm" layout="vertical" class="mt-6">
          <a-form-item label="申请理由" required>
            <a-textarea
              v-model:value="applyForm.reason"
              :rows="4"
              placeholder="请简要描述您的研究兴趣和申请理由"
              :auto-size="{ minRows: 4, maxRows: 6 }"
              show-count
              :maxlength="200"
            />
          </a-form-item>
        </a-form>

        <div class="flex justify-end mt-6 space-x-3">
          <a-button @click="dialogVisible = false">取消</a-button>
          <a-button 
            type="primary" 
            :loading="submitting"
            @click="submitApplication"
            v-motion
            :initial="{ scale: 1 }"
            :hover="{ scale: 1.05 }"
          >
            确认申请
          </a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { get, post, postJSON } from '../../utility/request'
import { useMotion } from '@vueuse/motion'
import { message } from 'ant-design-vue'
import { 
  SearchOutlined, 
  UserAddOutlined,
} from '@ant-design/icons-vue'

const teachers = ref([])
const filteredTeachers = ref([])
const searchQuery = ref('')
const loading = ref(false)
const successMessage = ref('')
const errorMessage = ref('')
const dialogVisible = ref(false)
const selectedTeacher = ref(null)
const submitting = ref(false)
const applyForm = ref({
  reason: ''
})

const fetchTeachers = () => {
  loading.value = true
  get('/api/user/all-teacher', null, 
    (message, data) => {
      teachers.value = data.users
      filteredTeachers.value = [...teachers.value]
      loading.value = false
    }, 
    (msg, data) => {
      console.error('获取老师列表失败:', msg)
      errorMessage.value = '获取老师列表失败，请刷新页面重试'
      loading.value = false
    },
    (msg, data) => {
      // 可能的其他处理
    }
  )
}

const filterTeachers = () => {
  const query = searchQuery.value.toLowerCase().trim()
  if (!query) {
    filteredTeachers.value = [...teachers.value]
    return
  }
  
  filteredTeachers.value = teachers.value.filter(teacher => {
    return (
      (teacher.username && teacher.username.toLowerCase().includes(query)) ||
      (teacher.faculty && teacher.faculty.toLowerCase().includes(query)) ||
      (teacher.researchFields && teacher.researchFields.toLowerCase().includes(query))
    )
  })
}

const applyForGuide = (teacher) => {
  selectedTeacher.value = teacher
  applyForm.value.reason = ''
  dialogVisible.value = true
}

const submitApplication = () => {
  if (!applyForm.value.reason.trim()) {
    message.error('请填写申请理由')
    return
  }

  submitting.value = true
  postJSON('/api/guide/apply', {
    teacherId: selectedTeacher.value.id,
    message: applyForm.value.reason
  }, 
  (msg, data) => {
    successMessage.value = `已成功向${selectedTeacher.value.username}老师发送指导申请，请等待老师审核`
    dialogVisible.value = false
    submitting.value = false
  }, 
  (msg, data) => {
    console.error('申请失败:', msg)
    errorMessage.value = msg || '申请失败，请稍后重试'
    submitting.value = false
  },
  (msg, data) => {
    // 可能的其他处理
  })
}

onMounted(() => {
  fetchTeachers()
})
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

.custom-modal :deep(.ant-modal-close) {
  @apply text-gray-500 dark:text-gray-400;
}

.custom-modal :deep(.ant-form-item-label > label) {
  @apply text-gray-700 dark:text-gray-300;
}

/* 卡片悬停效果 */
.card {
  @apply shadow-md dark:shadow-[0_4px_20px_rgba(0,0,0,0.2)] border border-gray-200 dark:border-gray-700;
  transform-origin: center bottom;
}

.card:hover {
  @apply shadow-xl;
}

/* 加载动画 */
:deep(.ant-spin-dot-item) {
  @apply bg-blue-500 dark:bg-blue-400;
}
</style>