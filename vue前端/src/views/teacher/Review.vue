<template>
    <div class="container mx-auto p-4 md:p-6 lg:p-8 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 min-h-screen">
      <div 
        class="mb-8 border-l-4 border-blue-500 dark:border-indigo-400 pl-4"
        v-motion
        :initial="{ opacity: 0, x: -50 }"
        :enter="{ opacity: 1, x: 0, transition: { duration: 500 } }"
      >
        <h1 class="text-2xl md:text-3xl font-bold mb-2 text-gray-800 dark:text-white bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-500 dark:from-blue-400 dark:to-indigo-300">评审中心</h1>
        <p class="text-gray-600 dark:text-gray-400">在这里创建和管理您的所有写作项目</p>
      </div>
  
  
      <!-- 我的项目 -->
      <div 
        class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6 backdrop-blur-sm bg-opacity-80 dark:bg-opacity-80 border border-gray-100 dark:border-gray-700"
        v-motion
        :initial="{ opacity: 0, y: 30 }"
        :enter="{ opacity: 1, y: 0, transition: { duration: 600, delay: 200 } }"
      >
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4 bg-gray-50 dark:bg-gray-700/50 p-4 rounded-lg backdrop-blur-sm">
          <h2 class="text-xl font-semibold flex items-center text-gray-800 dark:text-white">
            <span class="inline-block w-10 h-10 rounded-full bg-gradient-to-r from-purple-400 to-pink-500 flex items-center justify-center mr-3 shadow-md">
              <component :is="FolderOutlined" class="text-white" />
            </span>
            我的项目
          </h2>
          <div class="flex flex-wrap gap-2 w-full md:w-auto">
            <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索项目"
              class="w-full md:w-48"
              @search="handleSearch"
              v-motion
              :initial="{ opacity: 0, x: 20 }"
              :enter="{ opacity: 1, x: 0 }"
            />
            <a-select 
              v-model:value="filterType" 
              class="w-full md:w-32"
              :options="filterOptions"
              v-motion
              :initial="{ opacity: 0, x: 20 }"
              :enter="{ opacity: 1, x: 0, transition: { delay: 100 } }"
            />
            <a-select 
              v-model:value="sortBy" 
              class="w-full md:w-32"
              :options="sortOptions"
              v-motion
              :initial="{ opacity: 0, x: 20 }"
              :enter="{ opacity: 1, x: 0, transition: { delay: 200 } }"
            />
          </div>
        </div>
  
        <!-- 项目列表 -->
        <div v-if="!loading">
          <div v-if="filteredProjects.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <a-card 
              v-for="(project, index) in filteredProjects" 
              :key="project.id"
              class="w-full overflow-hidden rounded-xl border-0 project-card group"
              :bordered="false"
              v-motion
              :initial="{ opacity: 0, scale: 0.9 }"
              :enter="{ opacity: 1, scale: 1, transition: { delay: 50 * index, duration: 300 } }"
              :hovered="{ y: -5 }"
            >
              <template #title>
                <div class="flex justify-between items-start">
                  <h3 class="text-lg font-medium truncate text-gray-800 dark:text-white" :title="project.title">{{ project.title }}</h3>
                  <div class="flex items-center gap-2">
                    <a-tag 
                      :color="getStatusColor(project.status)"
                      class="rounded-full px-2"
                      v-if="project.status !== undefined"
                    >
                      {{ getStatusName(project.status) }}
                    </a-tag>
                    <a-tag 
                      :color="getTypeColor(project.type)"
                      class="rounded-full px-3 shadow-sm"
                    >
                      {{ getTypeName(project.type) }}
                    </a-tag>
                  </div>
                </div>
              </template>
              <template #cover>
                <div class="h-20 flex items-center justify-center relative overflow-hidden project-cover" :class="getCoverClass(project.type)">
                  <div class="absolute inset-0 bg-pattern opacity-10"></div>
                  <div class="absolute top-0 right-0 w-32 h-32 rounded-full bg-white/10 -mr-10 -mt-10 blur-2xl group-hover:scale-125 transition-transform duration-700"></div>
                  <div class="absolute bottom-0 left-0 w-24 h-24 rounded-full bg-white/10 -ml-10 -mb-10 blur-xl group-hover:scale-150 transition-transform duration-700"></div>
                  <div class="absolute inset-0 opacity-0 group-hover:opacity-100 bg-gradient-to-br from-white/0 to-white/20 transition-opacity duration-500 pointer-events-none"></div>
                  <component 
                    :is="getTypeIcon(project.type)" 
                    class="text-6xl text-white drop-shadow-lg transform transition-all duration-500 group-hover:scale-125 group-hover:rotate-12 relative z-10"
                  />
                </div>
              </template>
              <p class="text-sm text-gray-500 dark:text-gray-400 mt-2 line-clamp-2 h-10">{{ project.description }}</p>
              
              <!-- 作者信息 -->
              <div class="flex items-center mt-3 mb-2" v-if="project.user">
                <a-avatar :src="project.user?.avatar" class="mr-2" size="small">
                  {{ project.user?.username ? project.user.username.charAt(0).toUpperCase() : 'U' }}
                </a-avatar>
                <span class="text-xs text-gray-600 dark:text-gray-300">
                  {{ project.user?.username || '未知用户' }}
                  <span class="text-gray-400 dark:text-gray-500 ml-1">{{ project.user?.apartment || '' }}</span>
                </span>
              </div>

              <!-- 项目代码 -->
              <div v-if="project.code" class="mt-2 mb-2">
                <a-tag class="bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 border-0">
                  <span class="font-mono">{{ project.code }}</span>
                </a-tag>
              </div>
              
              <div class="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-4 mb-3">
                <span class="flex items-center">
                  <component :is="CalendarOutlined" class="mr-1" />
                  {{ formatDate(project.createdAt) }}
                </span>
                <span class="flex items-center">
                  <component :is="ClockCircleOutlined" class="mr-1" />
                  {{ formatDate(project.updatedAt) }}
                </span>
              </div>
              <a-divider class="my-2" />
              <div class="flex justify-between">
                <a-button 
                  type="primary"
                  @click="editProject(project)"
                  class="flex items-center shadow-md hover:shadow-lg transition-all duration-300 hover:scale-105"
                >
                  <template #icon><EditOutlined /></template>
                  评审
                </a-button>
                
              </div>
            </a-card>
          </div>
  
          <!-- 空状态 -->
          <a-empty 
            v-else 
            class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-8 text-center"
            description="暂无项目"
            v-motion
            :initial="{ opacity: 0 }"
            :enter="{ opacity: 1, transition: { duration: 500 } }"
          >
            <template #image>
              <div class="flex justify-center">
                <div class="w-32 h-32 rounded-full bg-gradient-to-br from-gray-200 to-gray-300 dark:from-gray-700 dark:to-gray-600 flex items-center justify-center shadow-inner">
                  <InboxOutlined class="text-7xl text-white opacity-70" />
                </div>
              </div>
            </template>
            <template #description>
              <p class="text-gray-500 dark:text-gray-400 mb-4">创建一个新项目开始您的写作之旅</p>
            </template>
            <a-button 
              type="primary"
              @click="showNewProjectDialog(projectTypes[0].id)"
              size="large"
              class="mt-2 shadow-lg hover:shadow-xl transition-shadow"
            >
              <template #icon><PlusOutlined /></template>
              创建项目
            </a-button>
          </a-empty>
        </div>
  
        <!-- 加载状态 -->
        <div v-else class="flex justify-center py-12">
            <div class="h-64 w-full flex items-center justify-center">
              <div class="loading-spinner">
                <LoadingOutlined class="text-5xl text-blue-500 dark:text-indigo-400" />
              </div>
            </div>
        </div>
      </div>
  
      <!-- 删除确认对话框 -->
      <a-modal
        v-model:visible="deleteModalVisible"
        title="确认删除"
        :footer="null"
        @cancel="deleteModalVisible = false"
        class="custom-modal"
      >
        <div class="text-center mb-6">
          <div class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-red-50 dark:bg-red-900/20 mb-4">
            <ExclamationCircleOutlined class="text-4xl text-red-500" />
          </div>
          <p class="text-lg font-medium">确定要删除项目 "{{ projectToDelete?.title }}" 吗？</p>
          <p class="text-gray-500 dark:text-gray-400 text-sm mt-2">此操作不可撤销。</p>
        </div>
        <div class="flex justify-center mt-6 gap-3">
          <a-button size="large" @click="deleteModalVisible = false" class="px-6">取消</a-button>
          <a-button type="primary" danger size="large" :loading="deleteLoading" @click="confirmDelete" class="px-6">
            删除
          </a-button>
        </div>
      </a-modal>
  
      <!-- 新建项目对话框 -->
      <a-modal
        v-model:visible="newProjectModalVisible"
        :title="`创建新${selectedProjectType ? getTypeName(selectedProjectType) : '项目'}`"
        :footer="null"
        @cancel="newProjectModalVisible = false"
        class="custom-modal"
        width="500px"
      >
        <div class="p-2">
          <div class="mb-6">
            <div class="flex justify-center mb-6">
              <div class="w-20 h-20 rounded-full flex items-center justify-center shadow-lg" :class="selectedProjectType ? getCoverClass(selectedProjectType) : 'bg-blue-500'">
                <component 
                  :is="selectedProjectType ? getTypeIcon(selectedProjectType) : FileTextOutlined" 
                  class="text-4xl text-white"
                />
              </div>
            </div>
            
            <a-form :model="newProjectForm" layout="vertical">
              <a-form-item label="项目标题" required>
                <a-input 
                  v-model:value="newProjectForm.title" 
                  placeholder="请输入项目标题"
                  :maxLength="50"
                  show-count
                />
              </a-form-item>
              
              <a-form-item label="项目类型">
                <a-select 
                  v-model:value="newProjectForm.type"
                  :options="projectTypes.map(type => ({ value: type.id, label: type.name }))"
                />
              </a-form-item>
              
              <a-form-item label="项目描述">
                <a-textarea 
                  v-model:value="newProjectForm.description" 
                  placeholder="请输入项目描述"
                  :rows="4"
                  :maxLength="200"
                  show-count
                />
              </a-form-item>
            </a-form>
          </div>
          
          <div class="flex justify-end gap-3">
            <a-button @click="newProjectModalVisible = false">取消</a-button>
            <a-button 
              type="primary" 
              :loading="createLoading"
              @click="confirmCreateProject"
              :disabled="!newProjectForm.title"
            >
              创建项目
            </a-button>
          </div>
        </div>
      </a-modal>
  
      <!-- 模板发放对话框 -->
      <a-modal
        v-model:visible="distributeModalVisible"
        title="发放模板"
        :footer="null"
        @cancel="distributeModalVisible = false"
        class="custom-modal"
        width="600px"
      >
        <div class="p-2">
          <div class="mb-6">
            <div class="flex flex-col items-center mb-6">
              <div class="w-20 h-20 rounded-full flex items-center justify-center shadow-lg mb-3" :class="currentProject ? getCoverClass(currentProject.type) : 'bg-blue-500'">
                <component 
                  :is="currentProject ? getTypeIcon(currentProject.type) : FileTextOutlined" 
                  class="text-4xl text-white"
                />
              </div>
              <h3 class="text-lg font-semibold">{{ currentProject?.title }}</h3>
              <p class="text-gray-500 dark:text-gray-400 text-sm mt-1">选择要发放的学生</p>
            </div>
            
            <div class="mb-4">
              <div class="flex justify-between items-center mb-3">
                <a-input-search
                  v-model:value="studentSearchKeyword"
                  placeholder="搜索学生"
                  class="w-64"
                  @search="searchStudents"
                />
                <div class="flex items-center gap-2">
                  <a-button size="small" @click="selectAllStudents">全选</a-button>
                  <a-button size="small" @click="unselectAllStudents">取消全选</a-button>
                </div>
              </div>
              
              <div class="max-h-60 overflow-y-auto border border-gray-200 dark:border-gray-700 rounded-lg bg-gray-50 dark:bg-gray-800/50 p-2">
                <a-spin :spinning="loadingStudents">
                  <div v-if="filteredStudents.length > 0">
                    <a-checkbox-group v-model:value="selectedStudentIds" class="w-full">
                      <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
                        <div 
                          v-for="student in filteredStudents" 
                          :key="student.id"
                          class="flex items-center p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                        >
                          <a-checkbox :value="student.id" class="w-full">
                            <div class="flex items-center">
                              <a-avatar :src="student.avatar" class="mr-2">
                                {{ student.name ? student.name.charAt(0) : 'U' }}
                              </a-avatar>
                              <div>
                                <div class="font-medium">{{ student.name }}</div>
                                <div class="text-xs text-gray-500 dark:text-gray-400">{{ student.id }}</div>
                              </div>
                            </div>
                          </a-checkbox>
                        </div>
                      </div>
                    </a-checkbox-group>
                  </div>
                  <a-empty 
                    v-else 
                    description="暂无学生" 
                    class="py-4"
                  />
                </a-spin>
              </div>
            </div>
            
            <div v-if="selectedStudentIds.length > 0" class="bg-blue-50 dark:bg-blue-900/20 p-3 rounded-lg mb-4">
              <div class="text-blue-600 dark:text-blue-400 text-sm mb-1">已选择 {{ selectedStudentIds.length }} 名学生</div>
              <div class="flex flex-wrap gap-1">
                <a-tag 
                  v-for="id in selectedStudentIds" 
                  :key="id"
                  closable
                  @close="removeSelectedStudent(id)"
                  class="mb-1"
                >
                  {{ getStudentNameById(id) }}
                </a-tag>
              </div>
            </div>
          </div>
          
          <div class="flex justify-end gap-3">
            <a-button @click="distributeModalVisible = false">取消</a-button>
            <a-button 
              type="primary" 
              :loading="distributingTemplate"
              @click="confirmDistribute"
              :disabled="selectedStudentIds.length === 0"
            >
              发放模板
            </a-button>
          </div>
        </div>
      </a-modal>
  
      <!-- 导入模板对话框 -->
      <a-modal
        v-model:visible="importModalVisible"
        title="导入模板"
        :footer="null"
        @cancel="importModalVisible = false"
        class="custom-modal"
        width="500px"
      >
        <div class="p-2">
          <div class="mb-6">
            <div class="flex justify-center mb-6">
              <div class="w-20 h-20 rounded-full bg-gradient-to-br from-purple-400 to-pink-600 flex items-center justify-center shadow-lg">
                <component 
                  :is="ImportOutlined" 
                  class="text-4xl text-white"
                />
              </div>
            </div>
            
            <a-form :model="importForm" layout="vertical">
              <a-form-item 
                label="模板代码" 
                required
                :validateStatus="importCodeError ? 'error' : ''"
                :help="importCodeError"
              >
                <a-input 
                  v-model:value="importForm.forkCode" 
                  placeholder="请输入模板代码"
                  :maxLength="100"
                />
              </a-form-item>
              
              <a-form-item label="项目标题" required>
                <a-input 
                  v-model:value="importForm.title" 
                  placeholder="请输入项目标题"
                  :maxLength="50"
                  show-count
                />
              </a-form-item>
              
              <a-form-item label="项目描述">
                <a-textarea 
                  v-model:value="importForm.description" 
                  placeholder="请输入项目描述"
                  :rows="4"
                  :maxLength="200"
                  show-count
                />
              </a-form-item>
            </a-form>
          </div>
          
          <div class="flex justify-end gap-3">
            <a-button @click="importModalVisible = false">取消</a-button>
            <a-button 
              type="primary" 
              :loading="importLoading"
              @click="confirmImport"
              :disabled="!importForm.forkCode || !importForm.title"
            >
              导入
            </a-button>
          </div>
        </div>
      </a-modal>
    </div>
  </template>
  
  <script setup>
  import { ref, computed, onMounted, reactive } from 'vue'
  import { useRouter } from 'vue-router'
  import { message, Modal } from 'ant-design-vue'
  import { 
    FileTextOutlined, 
    BookOutlined, 
    CheckSquareOutlined, 
    FileWordOutlined,
    DownloadOutlined,
    DeleteOutlined,
    InboxOutlined,
    LoadingOutlined,
    EditOutlined,
    PlusOutlined,
    ArrowRightOutlined,
    FolderOutlined,
    CalendarOutlined,
    ClockCircleOutlined,
    ExclamationCircleOutlined,
    SendOutlined,
    ImportOutlined
  } from '@ant-design/icons-vue'
  import { 
    getTemplateListStudent, 
    createTemplate, 
    deleteTemplate, 
    exportTemplate,
    sendTemplateToStudents,
    forkTemplate
  } from '../../api/template'
  import { getAllStudents } from '../../api/user'
  import { useMotion } from '@vueuse/motion'
  import { useUserStore } from '@/stores/user'
  const userStore = useUserStore()
  const router = useRouter()
  
  // 项目类型定义
  const projectTypes = [
    { 
      id: 'task', 
      name: '任务书', 
      icon: FileTextOutlined,
      description: '毕业设计任务书，明确研究内容和进度安排',
      route: '/student/editor/rws'
    },
    { 
      id: 'proposal', 
      name: '开题报告', 
      icon: BookOutlined,
      description: '阐述研究背景、目标、方法和预期成果',
      route: '/student/editor/proposal'
    },
    { 
      id: 'midterm', 
      name: '中期检查', 
      icon: CheckSquareOutlined,
      description: '汇报研究进展、遇到的问题和解决方案',
      route: '/student/editor/midterm'
    },
    { 
      id: 'thesis', 
      name: '毕业论文', 
      icon: FileWordOutlined,
      description: '完整的毕业论文，包含研究全过程和成果',
      route: '/student/editor/thesis'
    }
  ]
  
  // 项目状态定义
  const projectStatus = [
    { id: 0, name: '修改中', color: 'default' },
    { id: 1, name: '审核中', color: 'processing' },
    { id: 2, name: '已通过', color: 'success' },
    { id: 3, name: '待修改', color: 'warning' },
    { id: 4, name: '已拒绝', color: 'error' }
  ]
  
  // 状态变量
  const projects = ref([])
  const loading = ref(false)
  const filterType = ref('all')
  const sortBy = ref('newest')
  const searchKeyword = ref('')
  const deleteModalVisible = ref(false)
  const projectToDelete = ref(null)
  const deleteLoading = ref(false)
  
  // 新建项目相关状态
  const newProjectModalVisible = ref(false)
  const selectedProjectType = ref(null)
  const createLoading = ref(false)
  const newProjectForm = reactive({
    title: '',
    type: '',
    description: ''
  })
  
  // 下拉选项
  const filterOptions = [
    { value: 'all', label: '全部类型' },
    ...projectTypes.map(type => ({ value: type.id, label: type.name }))
  ]
  
  const sortOptions = [
    { value: 'newest', label: '最新创建' },
    { value: 'oldest', label: '最早创建' },
    { value: 'updated', label: '最近更新' },
    { value: 'name', label: '名称排序' }
  ]
  
  // 获取项目列表
  const fetchProjects = () => {
    loading.value = true
    getTemplateListStudent(
      (msg, data) => {
        projects.value = data.map(item => ({
          id: item.id,
          title: item.title || '未命名项目',
          type: item.type || 'task',
          description: item.description || '暂无描述',
          createdAt: new Date(item.createdAt || Date.now()),
          updatedAt: new Date(item.updatedAt || Date.now()),
          code: item.code,
          user: item.user || null,
          status: item.status || 0,
          titleId: item.titleId,
          coverId: item.coverId,
          contentId: item.contentId
        }))
        loading.value = false
      },
      (error) => {
        message.error('获取项目列表失败')
        loading.value = false
      }
    )
  }
  
  // 根据筛选和排序条件获取项目列表
  const filteredProjects = computed(() => {
    let result = [...projects.value]
    
    // 关键词搜索
    if (searchKeyword.value) {
      const keyword = searchKeyword.value.toLowerCase()
      result = result.filter(project => 
        project.title.toLowerCase().includes(keyword) || 
        project.description.toLowerCase().includes(keyword)
      )
    }
    
    // 类型筛选
    if (filterType.value !== 'all') {
      result = result.filter(project => project.type === filterType.value)
    }
    
    // 排序
    switch (sortBy.value) {
      case 'newest':
        result.sort((a, b) => b.createdAt - a.createdAt)
        break
      case 'oldest':
        result.sort((a, b) => a.createdAt - b.createdAt)
        break
      case 'updated':
        result.sort((a, b) => b.updatedAt - a.updatedAt)
        break
      case 'name':
        result.sort((a, b) => a.title.localeCompare(b.title))
        break
    }
    
    return result
  })
  
  // 获取项目类型名称
  const getTypeName = (typeId) => {
    const type = projectTypes.find(t => t.id === typeId)
    return type ? type.name : '未知类型'
  }
  
  // 获取项目类型图标
  const getTypeIcon = (typeId) => {
    switch (typeId) {
      case 'task': return FileTextOutlined
      case 'proposal': return BookOutlined
      case 'midterm': return CheckSquareOutlined
      case 'thesis': return FileWordOutlined
      default: return FileTextOutlined
    }
  }
  
  // 获取项目类型颜色
  const getTypeColor = (typeId) => {
    switch (typeId) {
      case 'task': return 'blue'
      case 'proposal': return 'green'
      case 'midterm': return 'gold'
      case 'thesis': return 'purple'
      default: return 'default'
    }
  }
  
  // 获取渐变类名
  const getGradientClass = (typeId) => {
    switch (typeId) {
      case 'task': return 'from-blue-400 to-blue-600'
      case 'proposal': return 'from-green-400 to-green-600'
      case 'midterm': return 'from-yellow-400 to-yellow-600'
      case 'thesis': return 'from-purple-400 to-purple-600'
      default: return 'from-blue-400 to-blue-600'
    }
  }
  
  // 获取渐变背景类名
  const getGradientBgClass = (typeId) => {
    switch (typeId) {
      case 'task': return 'from-blue-500 to-blue-700'
      case 'proposal': return 'from-green-500 to-green-700'
      case 'midterm': return 'from-yellow-500 to-yellow-700'
      case 'thesis': return 'from-purple-500 to-purple-700'
      default: return 'from-blue-500 to-blue-700'
    }
  }
  
  // 获取封面类名
  const getCoverClass = (typeId) => {
    switch (typeId) {
      case 'task': return 'bg-gradient-to-br from-blue-400 to-blue-600'
      case 'proposal': return 'bg-gradient-to-br from-green-400 to-green-600'
      case 'midterm': return 'bg-gradient-to-br from-yellow-400 to-yellow-600'
      case 'thesis': return 'bg-gradient-to-br from-purple-400 to-purple-600'
      default: return 'bg-gradient-to-br from-blue-400 to-blue-600'
    }
  }
  
  // 格式化日期
  const formatDate = (date) => {
    return new Intl.DateTimeFormat('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    }).format(date)
  }
  
  // 搜索处理
  const handleSearch = () => {
    // 搜索逻辑已在计算属性中实现
  }
  
  // 显示新建项目对话框
  const showNewProjectDialog = (typeId) => {
    selectedProjectType.value = typeId
    
    // 初始化表单数据
    const projectType = projectTypes.find(t => t.id === typeId)
    newProjectForm.type = typeId
    newProjectForm.title = `新${projectType.name}`
    newProjectForm.description = projectType.description
    
    // 显示对话框
    newProjectModalVisible.value = true
  }
  
  // 确认创建项目
  const confirmCreateProject = () => {
    if (!newProjectForm.title.trim()) {
      message.error('请输入项目标题')
      return
    }
    
    createLoading.value = true
    
    // 创建基础模板对象
    const newTemplate = {
      name: newProjectForm.title.trim(),
      type: newProjectForm.type,
      components: []
    }
    
    // 准备发送到后端的数据，格式与UniversalTemplateEditor一致
    const templateData = {
      title: newProjectForm.title.trim(),
      type: newProjectForm.type,
      description: newProjectForm.description.trim() || '',
      content: JSON.stringify(newTemplate)
    }
    
    createTemplate(
      templateData,
      (msg, data) => {
        message.success('创建项目成功')
        createLoading.value = false
        newProjectModalVisible.value = false
        fetchProjects()
        // 导航到相应的编辑器页面
        const projectType = projectTypes.find(t => t.id === newProjectForm.type)
        if (projectType) {
          
        }
      },
      (error) => {
        message.error('创建项目失败')
        createLoading.value = false
      }
    )
  }
  
  // 编辑项目
  const editProject = (project) => {
    router.push(`/view-template/${project.id}`)
  }
  
  // 下载项目
  const downloadProject = (project) => {
    loading.value = true
    exportTemplate(
      project.id,
      (msg, data) => {
        // 处理下载逻辑，通常后端会返回一个下载链接
        if (data && data.downloadUrl) {
          window.open(data.downloadUrl, '_blank')
        } else {
          // 如果没有直接的下载链接，可能需要处理二进制数据
          message.success('下载已开始')
        }
        loading.value = false
      },
      (error) => {
        message.error('下载失败')
        loading.value = false
      }
    )
  }
  
  // 显示删除确认
  const showDeleteConfirm = (project) => {
    projectToDelete.value = project
    deleteModalVisible.value = true
  }
  
  // 确认删除
  const confirmDelete = () => {
    if (!projectToDelete.value) return
    
    deleteLoading.value = true
    deleteTemplate(
      projectToDelete.value.id,
      (msg) => {
        message.success('删除成功')
        fetchProjects() // 重新获取项目列表
        deleteModalVisible.value = false
        deleteLoading.value = false
      },
      (error) => {
        message.error('删除失败')
        deleteLoading.value = false
      }
    )
  }
  
  // 模板发放相关状态
  const distributeModalVisible = ref(false)
  const currentProject = ref(null)
  const students = ref([])
  const loadingStudents = ref(false)
  const selectedStudentIds = ref([])
  const studentSearchKeyword = ref('')
  const distributingTemplate = ref(false)
  
  // 筛选后的学生列表
  const filteredStudents = computed(() => {
    if (!studentSearchKeyword.value) {
      return students.value
    }
    
    const keyword = studentSearchKeyword.value.toLowerCase()
    return students.value.filter(student => 
      student.name.toLowerCase().includes(keyword) || 
      student.id.toString().includes(keyword)
    )
  })
  
  // 显示模板发放对话框
  const showDistributeModal = (project) => {
    currentProject.value = project
    distributeModalVisible.value = true
    selectedStudentIds.value = []
    fetchStudents()
  }
  
  // 获取学生列表
  const fetchStudents = () => {
    loadingStudents.value = true
    getAllStudents(
      (msg, data) => {
        students.value = data.users.map(student => ({
          id: student.id,
          name: student.username,
          avatar: student.avatar || '',
        }))
        loadingStudents.value = false
      },
      (msg,data) => {
        message.error('获取学生列表失败')
        loadingStudents.value = false
      },
      (msg,data) => {
        message.error('获取学生列表失败')
        loadingStudents.value = false
      }
    )
  }
  
  // 搜索学生
  const searchStudents = () => {
    // 搜索逻辑已在计算属性中实现
  }
  
  // 全选学生
  const selectAllStudents = () => {
    selectedStudentIds.value = filteredStudents.value.map(student => student.id)
  }
  
  // 取消全选
  const unselectAllStudents = () => {
    selectedStudentIds.value = []
  }
  
  // 移除已选择的学生
  const removeSelectedStudent = (studentId) => {
    selectedStudentIds.value = selectedStudentIds.value.filter(id => id !== studentId)
  }
  
  // 根据学生ID获取学生姓名
  const getStudentNameById = (studentId) => {
    const student = students.value.find(s => s.id === studentId)
    return student ? student.name : `学生${studentId}`
  }
  
  // 确认发放模板
  const confirmDistribute = () => {
    if (!currentProject.value || selectedStudentIds.value.length === 0) {
      message.error('请选择至少一名学生')
      return
    }
    
    distributingTemplate.value = true
    sendTemplateToStudents(
      currentProject.value.code,
      selectedStudentIds.value,
      (msg) => {
        message.success('模板发放成功')
        distributingTemplate.value = false
        distributeModalVisible.value = false
      },
      (error) => {
        message.error('模板发放失败')
        distributingTemplate.value = false
      }
    )
  }
  
  // 导入模板相关状态
  const importModalVisible = ref(false)
  const importLoading = ref(false)
  const importCodeError = ref('')
  const importForm = reactive({
    forkCode: '',
    title: '',
    description: ''
  })
  
  // 显示导入模板对话框
  const showImportDialog = () => {
    importModalVisible.value = true
    importCodeError.value = ''
    // 清空表单
    importForm.forkCode = ''
    importForm.title = ''
    importForm.description = ''
  }
  
  // 确认导入模板
  const confirmImport = () => {
    // 表单验证
    if (!importForm.forkCode.trim()) {
      importCodeError.value = '请输入有效的模板代码'
      return
    }
    
    if (!importForm.title.trim()) {
      message.error('请输入项目标题')
      return
    }
    
    importLoading.value = true
    
    // 准备请求数据
    const requestData = {
      forkCode: importForm.forkCode.trim(),
      title: importForm.title.trim(),
      description: importForm.description.trim() || ''
    }
    
    forkTemplate(
      requestData,
      (msg, data) => {
        message.success('模板导入成功')
        importLoading.value = false
        importModalVisible.value = false
        fetchProjects() // 重新获取项目列表
      },
      (error) => {
        message.error('模板导入失败：' + (error || '未知错误'))
        importLoading.value = false
      }
    )
  }
  
  // 获取项目状态名称
  const getStatusName = (statusId) => {
    const status = projectStatus.find(s => s.id === statusId)
    return status ? status.name : '未知状态'
  }
  
  // 获取项目状态颜色
  const getStatusColor = (statusId) => {
    const status = projectStatus.find(s => s.id === statusId)
    return status ? status.color : 'default'
  }
  
  // 页面加载时获取项目列表
  onMounted(() => {
    fetchProjects()
  })
  </script>
  
  <style scoped>
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  
  /* 项目类型卡片样式 */
  .project-type-card {
    transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
    will-change: transform, box-shadow;
    transform-style: preserve-3d;
    backface-visibility: hidden;
  }
  
  .project-type-card:hover {
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04), 0 0 15px 0 rgba(66, 153, 225, 0.2);
  }
  
  .project-type-card:hover::after {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: 0.75rem;
    padding: 2px;
    background: linear-gradient(45deg, rgba(255,255,255,0.1), rgba(255,255,255,0.5), rgba(255,255,255,0.1));
    -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
    pointer-events: none;
    animation: border-glow 2s linear infinite;
  }
  
  @keyframes border-glow {
    0% {
      background-position: 0% 50%;
    }
    100% {
      background-position: 200% 50%;
    }
  }
  
  /* 项目卡片样式 */
  .project-card {
    transition: all 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
    box-shadow: 0 10px 20px rgba(0,0,0,0.05), 0 6px 6px rgba(0,0,0,0.07);
    will-change: transform, box-shadow;
    transform-style: preserve-3d;
    perspective: 1000px;
  }
  
  .project-card:hover {
    transform: translateY(-8px) scale(1.01) rotateX(2deg) rotateY(-2deg);
    box-shadow: 0 20px 30px rgba(0,0,0,0.1), 0 10px 10px rgba(0,0,0,0.05), 0 0 15px rgba(59, 130, 246, 0.2);
  }
  
  .project-card:hover::after {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: 0.75rem;
    padding: 2px;
    background: linear-gradient(45deg, rgba(59, 130, 246, 0.3), rgba(147, 197, 253, 0.3), rgba(59, 130, 246, 0.3));
    -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
    pointer-events: none;
    animation: border-glow 2s linear infinite;
  }
  
  /* 背景图案 */
  .bg-pattern {
    background-image: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23FFFFFF' fill-opacity='0.2'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  }
  
  /* 加载动画 */
  .loading-spinner {
    animation: spin 1.5s cubic-bezier(0.4, 0, 0.2, 1) infinite;
  }
  
  @keyframes spin {
    0% {
      transform: rotate(0deg) scale(1);
    }
    50% {
      transform: rotate(180deg) scale(1.2);
    }
    100% {
      transform: rotate(360deg) scale(1);
    }
  }
  
  /* 暗色模式适配 */
  :deep(.ant-empty-image) {
    filter: var(--tw-dark, none);
  }
  
  :deep(.custom-modal .ant-modal-content) {
    border-radius: 1rem;
    overflow: hidden;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  }
  
  :deep(.custom-modal .ant-modal-header) {
    border-bottom: none;
    padding-top: 1.5rem;
  }
  
  /* 项目封面效果 */
  .project-cover {
    position: relative;
    overflow: hidden;
  }
  
  .project-cover::before {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle, rgba(255,255,255,0.2) 0%, rgba(255,255,255,0) 60%);
    transform: rotate(30deg);
    animation: shimmer 8s infinite linear;
  }
  
  @keyframes shimmer {
    0% {
      transform: rotate(0deg) translate(-50%, -50%);
    }
    100% {
      transform: rotate(360deg) translate(-50%, -50%);
    }
  }
  
  @media (prefers-color-scheme: dark) {
    :deep(.ant-card-head) {
      background-color: #1f2937;
      border-bottom-color: #374151;
    }
    
    :deep(.ant-card-body) {
      background-color: #1f2937;
    }
    
    :deep(.ant-divider) {
      border-color: #374151;
    }
    
    :deep(.ant-modal-content),
    :deep(.ant-modal-header) {
      background-color: #1f2937;
      color: white;
    }
    
    :deep(.ant-modal-title) {
      color: white;
    }
  }
  </style>
  