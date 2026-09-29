<template>
  <div class="student-dashboard px-6 py-8 mx-auto bg-gray-50 dark:bg-gray-900 min-h-screen">
    <div class="dashboard-header mb-8"
      v-motion
      :initial="{ opacity: 0, y: -20 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 500 } }">
      <h1 class="text-3xl font-bold text-gray-800 dark:text-white mb-2">学生指导管理仪表盘</h1>
      <p class="text-gray-600 dark:text-gray-400">查看和管理您的指导学生及申请</p>
    </div>
    
    <!-- 概览卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <!-- 当前指导学生数量卡片 -->
      <div class="stat-card bg-white dark:bg-gray-800 rounded-xl shadow-md p-6"
        v-motion
        :initial="{ opacity: 0, x: -20 }"
        :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 100 } }">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 dark:text-gray-400 text-sm">当前指导学生</p>
            <h2 class="text-3xl font-bold text-blue-600 dark:text-blue-400 mt-2">
              {{ allStudents.length }}
            </h2>
          </div>
          <div class="bg-blue-100 dark:bg-blue-900/30 p-3 rounded-lg">
            <UserOutlined class="text-2xl text-blue-600 dark:text-blue-400" />
          </div>
        </div>
      </div>
      
      <!-- 待审核申请卡片 -->
      <div class="stat-card bg-white dark:bg-gray-800 rounded-xl shadow-md p-6"
        v-motion
        :initial="{ opacity: 0, x: -20 }"
        :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 200 } }">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 dark:text-gray-400 text-sm">待审核申请</p>
            <h2 class="text-3xl font-bold text-orange-500 dark:text-orange-400 mt-2">
              {{ pendingCount }}
            </h2>
          </div>
          <div class="bg-orange-100 dark:bg-orange-900/30 p-3 rounded-lg">
            <ClockCircleOutlined class="text-2xl text-orange-500 dark:text-orange-400" />
          </div>
        </div>
      </div>
      
      <!-- 最近创建的任务 -->
      <div class="stat-card bg-white dark:bg-gray-800 rounded-xl shadow-md p-6"
        v-motion
        :initial="{ opacity: 0, x: -20 }"
        :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 300 } }">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 dark:text-gray-400 text-sm">最近创建的任务</p>
            <h2 class="text-3xl font-bold text-green-600 dark:text-green-400 mt-2">
              {{ recentMissions }}
            </h2>
          </div>
          <div class="bg-green-100 dark:bg-green-900/30 p-3 rounded-lg">
            <CheckCircleOutlined class="text-2xl text-green-600 dark:text-green-400" />
          </div>
        </div>
      </div>
      
      <!-- 快速操作卡片 -->
      <div class="stat-card bg-white dark:bg-gray-800 rounded-xl shadow-md p-6"
        v-motion
        :initial="{ opacity: 0, x: -20 }"
        :enter="{ opacity: 1, x: 0, transition: { duration: 500, delay: 400 } }">
        <div class="flex flex-col h-full">
          <p class="text-gray-500 dark:text-gray-400 text-sm mb-4">快速操作</p>
          <div class="flex gap-2 mt-auto">
            <a-button type="primary" @click="showCreateMissionModal()" class="flex-1">
              创建任务
            </a-button>
            <a-button @click="activeTab = 'pending'" :disabled="!pendingCount" class="flex-1">
              处理申请
            </a-button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 搜索和过滤栏 -->
    <div class="flex flex-col md:flex-row justify-between items-center bg-white dark:bg-gray-800 rounded-xl shadow-md p-4 mb-6"
      v-motion
      :initial="{ opacity: 0, y: 20 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 500 } }">
      <div class="flex items-center space-x-4 w-full md:w-auto mb-4 md:mb-0">
        <a-input-search
          v-model:value="searchQuery"
          placeholder="搜索学生姓名或学号"
          class="w-full md:w-72"
          @search="handleSearch"
          allow-clear
        />
      </div>
      
      <div class="flex space-x-2 w-full md:w-auto">
        <a-button 
          :type="activeTab === 'current' ? 'primary' : 'default'"
          @click="activeTab = 'current'"
          class="flex-1 md:flex-none"
        >
          当前指导学生
        </a-button>
        <a-button 
          :type="activeTab === 'pending' ? 'primary' : 'default'"
          @click="activeTab = 'pending'"
          class="flex-1 md:flex-none"
        >
          待审核申请
          <a-badge v-if="pendingCount" :count="pendingCount" />
        </a-button>
      </div>
    </div>
    
    <!-- 内容区域 -->
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-md p-6 mb-6"
      v-motion
      :initial="{ opacity: 0, y: 20 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 500, delay: 200 } }">
      
      <div v-if="activeTab === 'current'">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold text-gray-800 dark:text-white">指导学生列表</h2>
          <a-button type="primary" @click="fetchGuideList">
            <template #icon><ReloadOutlined /></template>
            刷新
          </a-button>
        </div>
        <students-table 
          :data="filteredStudents" 
          :loading="loading" 
          @refresh="fetchGuideList"
          @create-mission="showCreateMissionModal"
        />
      </div>
      
      <div v-else-if="activeTab === 'pending'">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold text-gray-800 dark:text-white">待审核申请列表</h2>
          <a-button type="primary" @click="fetchPendingList">
            <template #icon><ReloadOutlined /></template>
            刷新
          </a-button>
        </div>
        <pending-students-table 
          :data="pendingStudents" 
          :loading="pendingLoading"
          @approve="handleApprove"
          @reject="handleReject"
          @refresh="fetchPendingList"
        />
      </div>
    </div>
    
    <!-- 创建任务对话框 -->
    <a-modal
      v-model:visible="missionModalVisible"
      title="创建指导任务"
      :footer="null"
      class="custom-modal"
      :width="600"
    >
      <div class="p-2" v-motion :initial="{ opacity: 0 }" :enter="{ opacity: 1 }">
        <a-form
          :model="missionForm"
          layout="vertical"
          :rules="rules"
          ref="formRef"
        >
          <a-form-item name="title" label="任务标题" required>
            <a-input v-model:value="missionForm.title" placeholder="请输入任务标题" />
          </a-form-item>
          
          <a-form-item name="description" label="任务描述" required>
            <a-textarea 
              v-model:value="missionForm.description" 
              placeholder="请详细描述任务内容和要求"
              :rows="4"
              :auto-size="{ minRows: 4, maxRows: 8 }"
            />
          </a-form-item>
          
          <a-form-item name="deadline" label="截止日期" required>
            <a-date-picker 
              v-model:value="missionForm.deadline"
              style="width: 100%"
              :disabledDate="disabledDate"
            />
          </a-form-item>
          
          <a-form-item name="studentId" label="指派学生" required>
            <a-select
              v-model:value="missionForm.studentId"
              placeholder="选择指派的学生"
              :options="studentOptions"
            />
          </a-form-item>
          
          <a-form-item name="attachments" label="附件资料">
            <a-upload
              v-model:fileList="fileList"
              :action="uploadUrl"
              :headers="uploadHeaders"
              :before-upload="beforeUpload"
              @change="handleUploadChange"
            >
              <a-button type="dashed" style="width: 100%">
                <template #icon><upload-outlined /></template>
                添加附件资料
              </a-button>
            </a-upload>
          </a-form-item>
        </a-form>
        
        <div class="flex justify-end mt-6 space-x-3">
          <a-button @click="missionModalVisible = false">取消</a-button>
          <a-button 
            type="primary" 
            :loading="submitting"
            @click="createMission"
            v-motion
            :initial="{ scale: 1 }"
            :hover="{ scale: 1.05 }"
          >
            创建任务
          </a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useMotion } from '@vueuse/motion'
import { message, notification } from 'ant-design-vue'
import { 
  SearchOutlined, 
  UserOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  UploadOutlined,
  ReloadOutlined,
  ClockCircleOutlined,
  PlusOutlined
} from '@ant-design/icons-vue'
import { get, post, postJSON, put, del } from '../../utility/request'
import StudentsTable from './components/StudentsTable.vue'
import PendingStudentsTable from './components/PendingStudentsTable.vue'
import dayjs from 'dayjs'

// 状态变量
const activeTab = ref('current')
const searchQuery = ref('')
const loading = ref(false)
const pendingLoading = ref(false)
const missionModalVisible = ref(false)
const submitting = ref(false)
const fileList = ref([])
const formRef = ref(null)

// 数据
const allStudents = ref([])
const pendingStudents = ref([])
const pendingCount = computed(() => pendingStudents.value.length)
const recentMissions = ref(0) // 添加最近任务计数

// 筛选后的学生列表
const filteredStudents = computed(() => {
  if (!searchQuery.value) return allStudents.value
  
  const query = searchQuery.value.toLowerCase().trim()
  return allStudents.value.filter(student => {
    return (
      (student.username && student.username.toLowerCase().includes(query)) ||
      (student.studentId && student.studentId.toLowerCase().includes(query))
    )
  })
})

// 学生选项
const studentOptions = computed(() => {
  return allStudents.value.map(student => ({
    value: student.id,
    label: `${student.username} (${student.studentId || '无学号'})`
  }))
})

// 任务表单
const missionForm = reactive({
  title: '',
  description: '',
  deadline: null,
  studentId: undefined,
  attachments: []
})

// 表单验证规则
const rules = {
  title: [{ required: true, message: '请输入任务标题', trigger: 'blur' }],
  description: [{ required: true, message: '请输入任务描述', trigger: 'blur' }],
  deadline: [{ required: true, message: '请选择截止日期', type: 'object', trigger: 'change' }],
  studentId: [{ required: true, message: '请选择指派的学生', trigger: 'change' }]
}

// 上传配置
const uploadUrl = '/api/mission/upload'
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token') || ''}`
}

// 禁用过去的日期
const disabledDate = (current) => {
  return current && current < dayjs().startOf('day')
}

// 处理上传前验证
const beforeUpload = (file) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    message.error('文件必须小于10MB!')
  }
  return isLt10M
}

// 处理上传状态变化
const handleUploadChange = (info) => {
  if (info.file.status === 'done') {
    message.success(`${info.file.name} 上传成功`)
    // 将文件信息添加到表单中
    missionForm.attachments.push({
      name: info.file.name,
      url: info.file.response.url
    })
  } else if (info.file.status === 'error') {
    message.error(`${info.file.name} 上传失败`)
  }
}

// 获取当前指导的学生列表
const fetchGuideList = () => {
  loading.value = true
  get('/api/guide/teacher/list', null, 
    (msg, data) => {
      // 处理返回的数据格式
      allStudents.value = data.map(guide => {
        return {
          id: guide.id,
          username: guide.Student.username,
          studentId: guide.Student.id,
          createdAt: guide.CreatedAt,
          faculty: guide.Student.apartment,
          major: guide.Student.profession,
          status: guide.status
        }
      })
      loading.value = false
      // 模拟最近任务数据
      fetchRecentMissionsCount()
    },
    (msg) => {
      message.error(msg || '获取学生列表失败')
      loading.value = false
    }
  )
}

// 获取待审核的申请列表
const fetchPendingList = () => {
  pendingLoading.value = true
  get('/api/guide/teacher/pending', null, 
    (msg, data) => {
      // 处理返回的数据格式
      pendingStudents.value = data
      pendingLoading.value = false
    },
    (msg) => {
      message.error(msg || '获取待审核申请失败')
      pendingLoading.value = false
    }
  )
}

// 模拟获取最近任务数量
const fetchRecentMissionsCount = () => {
  // 这里应该发送请求获取实际的任务数量，这里使用模拟数据
  recentMissions.value = Math.floor(Math.random() * 10)
}

// 处理搜索
const handleSearch = () => {
  // 直接使用计算属性筛选
}

// 显示创建任务对话框
const showCreateMissionModal = (studentId) => {
  resetMissionForm()
  if (studentId) {
    missionForm.studentId = studentId
  }
  missionModalVisible.value = true
}

// 重置任务表单
const resetMissionForm = () => {
  Object.assign(missionForm, {
    title: '',
    description: '',
    deadline: null,
    studentId: undefined,
    attachments: []
  })
  fileList.value = []
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

// 创建任务
const createMission = () => {
  formRef.value.validate().then(() => {
    submitting.value = true
    
    const missionData = {
      title: missionForm.title,
      description: missionForm.description,
      deadline: missionForm.deadline.format('YYYY-MM-DD'),
      studentId: missionForm.studentId,
      attachments: missionForm.attachments
    }
    
    postJSON('/api/mission/create', missionData, 
      (msg) => {
        message.success('任务创建成功')
        missionModalVisible.value = false
        submitting.value = false
        resetMissionForm()
      },
      (msg) => {
        message.error(msg || '创建任务失败')
        submitting.value = false
      }
    )
  }).catch(errorInfo => {
    console.log('验证失败:', errorInfo)
  })
}

// 审核通过申请
const handleApprove = (applicationId) => {
  put(`/api/guide/${applicationId}/status`, { status: 'approved' }, 
    (msg) => {
      message.success('已批准申请')
      fetchPendingList()
      fetchGuideList()
    },
    (msg) => {
      message.error(msg || '操作失败')
    }
  )
}

// 拒绝申请
const handleReject = (applicationId) => {
  put(`/api/guide/${applicationId}/status`, { status: 'rejected' }, 
    (msg) => {
      message.success('已拒绝申请')
      fetchPendingList()
    },
    (msg) => {
      message.error(msg || '操作失败')
    }
  )
}

onMounted(() => {
  fetchGuideList()
  fetchPendingList()
})
</script>

<style scoped>
.student-dashboard {
  @apply transition-all duration-300;
}

.stat-card {
  @apply transition-all duration-300 hover:shadow-lg;
}

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

/* 标签页样式 */
:deep(.ant-tabs-tab) {
  @apply text-gray-600 dark:text-gray-300;
}

:deep(.ant-tabs-tab.ant-tabs-tab-active .ant-tabs-tab-btn) {
  @apply text-blue-600 dark:text-blue-400;
}

:deep(.ant-tabs-ink-bar) {
  @apply bg-blue-500 dark:bg-blue-400;
}
</style>
