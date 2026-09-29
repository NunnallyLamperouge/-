<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { Document, Packer, ImageRun } from 'docx'
import ComponentPalette from './components/ComponentPalette.vue'
import TemplateCanvas from './components/TemplateCanvas.vue'
import PropertyPanel from './components/PropertyPanel.vue'
import NewTemplateDialog from './components/NewTemplateDialog.vue'
import ExportDialog from './components/ExportDialog.vue'
import TemplateList from './components/TemplateList.vue'

// 导入工具函数
import { defaultTemplate, createComponent, uploadImage } from './utils/templateUtils'
import { generateDoc } from './utils/documentExporter'
// 导入请求函数
import { get, post, postJSON } from '@/utility/request'

// 当前模板状态
const currentTemplate = ref(JSON.parse(JSON.stringify(defaultTemplate)))
const savedTemplates = ref([])
const currentTemplateIndex = ref(-1)

// 编辑器状态
const selectedComponentId = ref(null)
const isDragging = ref(false)
const showNewTemplateDialog = ref(false)
const isExporting = ref(false) // 添加导出状态
const exportProgress = ref(0) // 添加导出进度
const exportStatus = ref('') // 添加导出状态信息

// 从后端加载已保存的模板
onMounted(() => {
  loadTemplates()
})

// 监听模板变化，自动保存
watch(currentTemplate, (newVal) => {
  if (currentTemplateIndex.value >= 0) {
    const templateToUpdate = savedTemplates.value[currentTemplateIndex.value]
    if (templateToUpdate && templateToUpdate.id) {
      updateTemplateToServer(templateToUpdate.id, newVal)
    }
  }
}, { deep: true })

// 处理组件重新排序
const handleComponentReorder = (componentId, fromIndex, toIndex) => {
  // 获取要移动的组件
  const component = currentTemplate.value.components[fromIndex]
  
  // 从原位置删除
  currentTemplate.value.components.splice(fromIndex, 1)
  
  // 插入到新位置
  currentTemplate.value.components.splice(toIndex, 0, component)
}

// 从后端加载已保存的模板
const loadTemplates = () => {
  get('/api/template/list', null, 
    (message, data) => {
      // 将后端返回的模板数据转换为前端需要的格式
      savedTemplates.value = data.map(item => {
        const template = JSON.parse(item.content);
        return {
          ...template,
          id: item.id,
          name: item.title,
          type: item.type, 
          description: item.description
        };
      });
      
      if (savedTemplates.value.length > 0) {
        loadTemplate(0)
      }
    },
    (errorMsg) => {
      console.error('加载模板失败:', errorMsg)
    }
  )
}

// 创建新模板
const createNewTemplate = (formData) => {
  const newTemplate = JSON.parse(JSON.stringify(defaultTemplate))
  newTemplate.name = formData.name
  newTemplate.type = formData.type
  
  // 准备发送到后端的数据
  const templateData = {
    title: formData.name,
    type: formData.type,
    description: formData.description || '',
    content: JSON.stringify(newTemplate)
  }
  
  // 发送到后端创建
  postJSON('/api/template/create', templateData,
    (message, data) => {
      // 创建成功后，将返回的数据添加到模板列表
      const createdTemplate = {
        id: data.id,
        name: data.title,
        type: data.type,
        description: data.description,
        ...newTemplate
      }
      
      savedTemplates.value.push(createdTemplate)
      currentTemplateIndex.value = savedTemplates.value.length - 1
      currentTemplate.value = JSON.parse(JSON.stringify(createdTemplate))
      
      showNewTemplateDialog.value = false
    },
    (errorMsg) => {
      console.error('创建模板失败:', errorMsg)
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
      console.log('模板更新成功')
    },
    (errorMsg) => {
      console.error('更新模板失败:', errorMsg)
    }
  )
}

// 加载模板
const loadTemplate = (index) => {
  currentTemplateIndex.value = index
  currentTemplate.value = JSON.parse(JSON.stringify(savedTemplates.value[index]))
  
  // 如果需要从后端获取最新数据，可以添加以下代码
  const templateId = savedTemplates.value[index].id
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
      console.log(fullTemplate)
      // 更新保存的模板列表中的对应项
      savedTemplates.value[index] = JSON.parse(JSON.stringify(fullTemplate))
    },
    (errorMsg) => {
      console.error('获取模板详情失败:', errorMsg)
    }
  )
}

// 删除模板
const deleteTemplate = (index) => {
  const templateId = savedTemplates.value[index].id
  
  // 发送删除请求到后端
  post(`/api/template/delete/${templateId}`, {},
    (message, data) => {
      // 从本地列表中移除
      savedTemplates.value.splice(index, 1)
      
      if (currentTemplateIndex.value === index) {
        if (savedTemplates.value.length > 0) {
          loadTemplate(0)
        } else {
          currentTemplateIndex.value = -1
          currentTemplate.value = JSON.parse(JSON.stringify(defaultTemplate))
        }
      } else if (currentTemplateIndex.value > index) {
        currentTemplateIndex.value--
      }
    },
    (errorMsg) => {
      console.error('删除模板失败:', errorMsg)
    }
  )
}

// 处理组件拖放
const handleComponentDrop = (componentType, position) => {
  const newComponent = createComponent(componentType, position)
  currentTemplate.value.components.push(newComponent)
}

// 更新组件属性
const updateComponentProperty = (componentId, property, value) => {
  const component = currentTemplate.value.components.find(c => c.id === componentId)
  if (component) {
    component[property] = value
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
    }, 1500)
  } catch (error) {
    console.error('文档生成错误:', error)
    exportStatus.value = `导出失败: ${error.message}`
    setTimeout(() => {
      isExporting.value = false
    }, 3000)
    alert('文档生成失败，请检查模板内容')
  }
}

// 关闭导出对话框
const handleCloseExportDialog = () => {
  isExporting.value = false
}
</script>

<template>
  <div class="flex h-screen bg-gray-100">
    <!-- 左侧组件面板 -->
    <div class="w-72 bg-white shadow-lg p-4 overflow-y-auto flex flex-col border-r border-gray-200">
      <div class="mb-4">
        <h1 class="text-xl font-bold text-gray-800">模板编辑器</h1>
        <p class="text-sm text-gray-500">拖拽组件到右侧画布</p>
      </div>
      <div class="mt-4 pt-4 border-t border-gray-200">
        <div class="flex flex-col gap-2">
          <button 
            class="w-full px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 flex items-center justify-center gap-2"
            @click="showNewTemplateDialog = true"
          >
            <span class="text-lg">+</span> 新建模板
          </button>
          <button 
            class="w-full px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 flex items-center justify-center gap-2"
            @click="handleExportDoc"
            :disabled="isExporting"
          >
            <span v-if="!isExporting" class="text-lg">↓</span>
            <span v-else class="animate-spin">⟳</span>
            {{ isExporting ? '导出中...' : '导出文档' }}
          </button>
        </div>
        
        <!-- 模板列表 -->
        <TemplateList 
          :templates="savedTemplates"
          :currentIndex="currentTemplateIndex"
          @select="loadTemplate"
          @delete="deleteTemplate"
        />
      </div>
      <div class="flex-1">
        <ComponentPalette 
          @component-drag-start="isDragging = true"
          @component-drag-end="isDragging = false"
        />
      </div>
    </div>
    
    <!-- 中间画布区域 -->
    <div class="flex-1 p-4 overflow-auto">
      <TemplateCanvas
        v-model:selectedComponentId="selectedComponentId"
        :template="currentTemplate"
        :isDragging="isDragging"
        @component-drop="handleComponentDrop"
        @component-update="updateComponentProperty"
        @component-delete="deleteComponent"
        @component-reorder="handleComponentReorder"
      />
    </div>
    
    <!-- 右侧属性面板 -->
    <div class="w-80 bg-white shadow-lg p-4 overflow-y-auto">
      <PropertyPanel
        v-if="selectedComponentId"
        :component="currentTemplate.components.find(c => c.id === selectedComponentId)"
        @update="updateComponentProperty"
        @delete="deleteComponent"
      />
      <div v-else class="text-center py-8 text-gray-400">
        <p class="text-lg">请选择一个组件</p>
        <p class="text-sm">选择组件后可在此处编辑属性</p>
      </div>
    </div>
    
    <!-- 新建模板对话框 -->
    <NewTemplateDialog
      :show="showNewTemplateDialog"
      @close="showNewTemplateDialog = false"
      @create="createNewTemplate"
    />
    
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
.component-dragging {
  opacity: 0.5;
  cursor: move;
}
</style>

