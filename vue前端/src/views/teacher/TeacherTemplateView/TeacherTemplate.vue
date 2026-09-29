<script setup>
import { ref, reactive, onMounted, watch, onUnmounted } from 'vue'
import { Document, Packer, ImageRun } from 'docx'
import { useMotion } from '@vueuse/motion'
import { message, notification } from 'ant-design-vue'
import { 
  DownloadOutlined, 
  LoadingOutlined,
  EditOutlined,
  DeleteOutlined,
  SaveOutlined,
  FileTextOutlined,
  CommentOutlined,
  HistoryOutlined,
  EyeOutlined,
  RollbackOutlined,
  LeftOutlined
} from '@ant-design/icons-vue'
import TemplateCanvas from './components/TemplateCanvas.vue'
import PropertyPanel from './components/PropertyPanel.vue'
import ExportDialog from './components/ExportDialog.vue'
const route = useRoute()
const templateId = route.params.id
// 导入工具函数
import { defaultTemplate } from './utils/templateUtils'
import { generateDoc } from './utils/documentExporter'
// 导入请求函数
import { get, post, postJSON } from '@/utility/request'
import { useRoute } from 'vue-router'

// 当前模板状态
const currentTemplate = ref(JSON.parse(JSON.stringify(defaultTemplate)))
const currentTemplateIndex = ref(-1)

// 历史记录状态
const historyIDs = ref([])
const showHistory = ref(false)
const leftPanelRef = ref(null)

// 编辑器状态
const selectedComponentId = ref(null)
const isExporting = ref(false) // 添加导出状态
const exportProgress = ref(0) // 添加导出进度
const exportStatus = ref('') // 添加导出状态信息
const hasUnsavedChanges = ref(false) // 是否有未保存的更改
const lastSavedTemplate = ref(null) // 最后一次保存的模板内容
const annotation = ref('') // 当前标注内容

// 创建DOM引用
const canvasAreaRef = ref(null)
const rightPanelRef = ref(null)

// 动画相关 - 确保在设置动画前已经创建了ref
const canvasArea = useMotion(canvasAreaRef, {
  initial: { opacity: 0, scale: 0.9 },
  enter: { opacity: 1, scale: 1, transition: { duration: 600, delay: 200 } }
})

const rightPanel = useMotion(rightPanelRef, {
  initial: { opacity: 0, x: 50 },
  enter: { opacity: 1, x: 0, transition: { duration: 500, delay: 100 } }
})

// 加载模板历史记录
const loadHistoryList = () => {
  get(`/api/template/history/${templateId}/list`, null,
    (message, data) => {
      historyIDs.value = data || []
    },
    (msg) => {
      console.error('获取历史记录失败:', msg)
      message.error('获取历史记录失败')
    }
  )
}

// 加载特定版本的模板历史
const loadTemplateHistory = (historyId) => {
  get(`/api/template/history/${historyId}/content`, null,
    (message, data) => {
      // 解析Content字段中存储的模板内容
      const templateContent = JSON.parse(data)
      // 合并模板信息
      const fullTemplate = {
        id: data.id,
        name: data.title,
        type: data.type,
        description: data.description,
        ...templateContent
      }
      
      // 更新当前模板
      currentTemplate.value = fullTemplate
      // 保存最后加载的模板状态，用于比较变更
      lastSavedTemplate.value = JSON.parse(JSON.stringify(fullTemplate))
      hasUnsavedChanges.value = false
    },
    (msg, data) => {
      console.error('获取模板详情失败:', msg)
      message.error('获取模板详情失败')
    },
    (msg, data) => {
      console.error('获取模板详情失败:', msg)
      message.error('获取模板详情失败')
    }
  )
}

// 格式化日期时间
const parseDateTime = (dateTimeStr) => {
  const date = new Date(dateTimeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 从后端加载已保存的模板
onMounted(() => {
  loadTemplate()
  loadHistoryList() // 加载历史记录
  // 设置自动保存定时器 (3分钟)
  const autoSaveInterval = setInterval(() => {
    if (hasUnsavedChanges.value && templateId) {
      saveTemplate()
    }
  }, 3 * 60 * 1000)
  
  // 组件卸载时清除定时器
  onUnmounted(() => {
    clearInterval(autoSaveInterval)
  })
})

// 监听模板变化
watch(currentTemplate, (newVal) => {
  // 检查是否有实际变化
  if (lastSavedTemplate.value) {
    // 比较当前模板和最后保存的模板是否不同
    const currentStr = JSON.stringify(newVal)
    const savedStr = JSON.stringify(lastSavedTemplate.value)
    hasUnsavedChanges.value = currentStr !== savedStr
  } else {
    hasUnsavedChanges.value = true
  }
}, { deep: true })

// 手动保存模板
const saveTemplate = () => {
  if (!hasUnsavedChanges.value) {
    message.info('没有任何修改需要保存')
    return
  }
  
  if (templateId) {
    updateTemplateToServer(templateId, currentTemplate.value)
    // 更新最后保存的模板状态
    lastSavedTemplate.value = JSON.parse(JSON.stringify(currentTemplate.value))
    hasUnsavedChanges.value = false
  }
}

// 处理模板内容变更通知
const handleTemplateChanged = () => {
  hasUnsavedChanges.value = true
}

// 加载模板
const loadTemplate = () => {
  get(`/api/template/detail/${templateId}`, null,
    (message, data) => {
      // 解析Content字段中存储的模板内容
      const templateContent = JSON.parse(data.content)
      // 合并模板信息
      const fullTemplate = {
        id: data.id,
        name: data.title,
        type: data.type,
        description: data.description,
        ...templateContent
      }
      
      // 更新当前模板
      currentTemplate.value = fullTemplate
      // 保存最后加载的模板状态，用于比较变更
      lastSavedTemplate.value = JSON.parse(JSON.stringify(fullTemplate))
      hasUnsavedChanges.value = false
    },
    (msg, data) => {
      console.error('获取模板详情失败:', msg)
      message.error('获取模板详情失败')
    },
    (msg, data) => {
      console.error('获取模板详情失败:', msg)
      message.error('获取模板详情失败')
    }
  )
}

// 更新模板到服务器
const updateTemplateToServer = (templateId, templateData) => {
  // 准备发送到后端的数据
  const updateData = {
    title: templateData.name,
    type: templateData.type,
    description: templateData.description || '',
    content: JSON.stringify(templateData)
  }
  
  // 发送到后端更新
  postJSON(`/api/template/update/${templateId}`, updateData,
    (message, data) => {
      notification.success({
        message: '保存成功',
        description: '模板标注已成功更新',
        placement: 'bottomRight',
        duration: 2
      })
    },
    (errorMsg) => {
      notification.error({
        message: '保存失败',
        description: errorMsg || '更新模板时发生错误',
        placement: 'bottomRight'
      })
    }
  )
}

// 处理组件属性变更
const updateComponentProperty = (componentId, property, value) => {
  const component = currentTemplate.value.components.find(c => c.id === componentId)
  if (component) {
    component[property] = value
    // 标记为有未保存的更改
    hasUnsavedChanges.value = true
  }
}

// 添加标注
const addAnnotation = (componentId, annotationText) => {
  if (!annotationText.trim()) {
    message.warning('标注内容不能为空')
    return
  }
  
  const component = currentTemplate.value.components.find(c => c.id === componentId)
  if (component) {
    // 如果组件没有annotations属性，创建一个空数组
    if (!component.annotations) {
      component.annotations = []
    }
    
    // 添加新标注
    component.annotations.push({
      id: Date.now(),
      content: annotationText,
      date: new Date().toLocaleString()
    })
    
    // 清空标注输入框
    annotation.value = ''
    
    // 标记为有未保存的更改
    hasUnsavedChanges.value = true
    
    message.success('标注已添加')
  }
}

// 删除组件
const deleteComponent = (componentId) => {
  const index = currentTemplate.value.components.findIndex(c => c.id === componentId)
  if (index !== -1) {
    currentTemplate.value.components.splice(index, 1)
    if (selectedComponentId.value === componentId) {
      selectedComponentId.value = null
    }
    message.success('组件已删除')
    
    // 标记为有未保存的更改
    hasUnsavedChanges.value = true
  }
}

// 导出文档
const handleExportDoc = async () => {
  try {
    isExporting.value = true
    
    await generateDoc(
      currentTemplate.value, 
      (progress) => { exportProgress.value = progress }, 
      (status) => { exportStatus.value = status }
    )
    
    // 延迟关闭导出状态
    setTimeout(() => {
      isExporting.value = false
      notification.success({
        message: '导出成功',
        description: '文档已成功生成',
        placement: 'bottomRight',
        duration: 3
      })
    }, 1500)
  } catch (error) {
    console.error('文档生成错误:', error)
    exportStatus.value = `导出失败: ${error.message}`
    setTimeout(() => {
      isExporting.value = false
      notification.error({
        message: '导出失败',
        description: error.message || '生成文档时发生错误',
        placement: 'bottomRight'
      })
    }, 3000)
  }
}

// 关闭导出对话框
const handleCloseExportDialog = () => {
  isExporting.value = false
}
</script>

<template>
  <div class="flex h-screen bg-gray-50 overflow-hidden">
    <!-- 左侧历史记录栏 -->
    <div 
      ref="leftPanelRef"
      class="w-72 bg-white shadow-xl p-5 overflow-y-auto border-r border-gray-200 transition-all duration-300 ease-in-out"
      v-show="showHistory"
    >
      <div class="text-lg font-semibold text-gray-700 mb-4 flex items-center justify-between">
        <div class="flex items-center">
          <HistoryOutlined class="mr-2 text-blue-500" />
          历史记录
        </div>
        <a-button type="text" @click="showHistory = false">
          <template #icon><LeftOutlined /></template>
        </a-button>
      </div>

      <div class="space-y-3">
        <div
          v-for="(history, index) in historyIDs"
          :key="history.id"
          class="component-item p-3 bg-white rounded-lg shadow hover:shadow-md cursor-pointer transition-all"
          :class="{ 'transform hover:-translate-y-1': true }"
        >
          <div class="flex items-center gap-3">
            <div class="flex-1 min-w-0" @click="loadTemplateHistory(history.id)"
              >
                <h3 class="font-medium text-gray-800 truncate">
                  第{{ historyIDs.length - index }}次修改
                </h3>
                <p class="text-xs text-gray-500 truncate">
                  修改于{{ parseDateTime(history.CreatedAt) }}
                </p>
              </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 中间画布区域 -->
    <div 
      ref="canvasAreaRef"
      class="flex-1 p-6 overflow-auto bg-gradient-to-br from-gray-50 to-gray-100 transition-all duration-300 ease-in-out"
      :class="{'pr-6': !selectedComponentId, 'pr-0': selectedComponentId}"
    >
      <div class="mb-4 flex justify-between items-center">
        <div class="flex items-center">
          <FileTextOutlined class="text-blue-500 mr-2 text-xl" />
          <h2 class="text-xl font-semibold text-gray-700">
            {{ currentTemplate.name || '未命名模板' }}
          </h2>
          <a-tag color="blue" class="ml-2">{{ currentTemplate.type || '标准模板' }}</a-tag>
        </div>
        
        <div class="flex space-x-2">
          <a-button 
            type="default"
            class="flex items-center"
            @click="showHistory = !showHistory"
          >
            <template #icon><HistoryOutlined /></template>
            历史记录
          </a-button>

          <a-button 
            type="primary"
            class="flex items-center"
            @click="saveTemplate"
            :disabled="!hasUnsavedChanges"
          >
            <template #icon><SaveOutlined /></template>
            保存标注
            <a-badge v-if="hasUnsavedChanges" status="processing" color="red" />
          </a-button>
          
          <a-button 
            type="primary"
            class="flex items-center"
            :loading="isExporting"
            @click="handleExportDoc"
            danger
          >
            <template #icon><DownloadOutlined v-if="!isExporting" /></template>
            {{ isExporting ? '导出中...' : '导出文档' }}
          </a-button>
        </div>
      </div>
      
      <div class="p-6 min-h-[calc(100vh-120px)]">
        <TemplateCanvas
          v-model:selectedComponentId="selectedComponentId"
          :template="currentTemplate"
          :hasUnsavedChanges="hasUnsavedChanges"
          @component-update="updateComponentProperty"
          @component-delete="deleteComponent"
          @template-changed="handleTemplateChanged"
        />
      </div>
    </div>
    
    <!-- 右侧标注面板 -->
    <div 
      ref="rightPanelRef"
      class="w-96 bg-white shadow-xl p-5 overflow-y-auto border-l border-gray-200 transition-all duration-300 ease-in-out"
      v-if="selectedComponentId"
    >
      <h2 class="text-lg font-semibold text-gray-700 mb-4 flex items-center">
        <CommentOutlined class="mr-2 text-blue-500" />
        标注面板
      </h2>
      
      <a-divider class="my-3" />
      
      <PropertyPanel
        :component="currentTemplate.components.find(c => c.id === selectedComponentId)"
        @update="updateComponentProperty"
        @delete="deleteComponent"
        @add-annotation="addAnnotation"
        @template-changed="handleTemplateChanged"
      />
    </div>
    
    <!-- 导出进度对话框 -->
    <ExportDialog
      :isExporting="isExporting"
      :exportProgress="exportProgress"
      :exportStatus="exportStatus"
      @close="handleCloseExportDialog"
    />
  </div>
</template>

<style scoped>
:deep(.ant-btn-primary) {
  background-color: #1890ff;
}

:deep(.ant-btn-primary.ant-btn-dangerous) {
  background-color: #ff4d4f;
}

:deep(.ant-empty-image) {
  height: 80px;
}

:deep(.ant-divider) {
  margin: 12px 0;
}

.component-item {
  border: 1px solid #f0f0f0;
}

.component-item:hover {
  border-color: #1890ff;
}
</style>

