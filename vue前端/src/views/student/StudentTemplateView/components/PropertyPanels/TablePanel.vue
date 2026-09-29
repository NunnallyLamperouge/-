<script setup>
import { ref } from 'vue'

const props = defineProps({
  component: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update'])

// 预览表格
const showTablePreview = ref(false)

// 当前编辑的字段索引
const currentEditingFieldIndex = ref(-1)

// 可用的子组件类型
const availableContentComponents = [
  { type: 'text', name: '文本段落', icon: '📄' },
  { type: 'list', name: '列表', icon: '📋' },
  { type: 'reference', name: '参考文献', icon: '📚' },
  { type: 'image', name: '图片', icon: '🖼️' },
  { type: 'title', name: '大标题', icon: '📌' },
  { type: 'paragraph', name: '普通文本段落', icon: '📝' },
  { type: 'schedule', name: '进度表', icon: '📅' }
]

// 切换表格预览
const toggleTablePreview = () => {
  showTablePreview.value = !showTablePreview.value
}

// 表格边框预设
const borderPresets = [
  { name: '无边框', icon: '◻️', preset: { top: false, bottom: false, left: false, right: false, insideHorizontal: false, insideVertical: false } },
  { name: '外边框', icon: '🔲', preset: { top: true, bottom: true, left: true, right: true, insideHorizontal: false, insideVertical: false } },
  { name: '全边框', icon: '📊', preset: { top: true, bottom: true, left: true, right: true, insideHorizontal: true, insideVertical: true } },
  { name: '横线', icon: '☰', preset: { top: true, bottom: true, left: false, right: false, insideHorizontal: true, insideVertical: false } },
  { name: '竖线', icon: '▤', preset: { top: false, bottom: false, left: true, right: true, insideHorizontal: false, insideVertical: true } }
]

// 应用边框预设
const applyBorderPreset = (preset) => {
  if (!props.component.settings) {
    props.component.settings = {}
  }
  
  if (!props.component.settings.borderStyle) {
    props.component.settings.borderStyle = {}
  }
  
  props.component.settings.borderStyle = { ...preset }
  updateProperty('settings', { ...props.component.settings })
}

// 确保组件设置存在
const ensureComponentSettings = (component) => {
  if (!component.settings) {
    component.settings = {}
  }
  
  if (!component.settings.borderStyle) {
    component.settings.borderStyle = {
      top: false,
      bottom: false,
      left: false,
      right: false,
      insideHorizontal: false,
      insideVertical: false
    }
  }
  
  if (!component.settings.cellStyle) {
    component.settings.cellStyle = {
      alignment: 'left',
      padding: '8px'
    }
  }
  
  if (!component.settings.headerStyle) {
    component.settings.headerStyle = {
      bold: true,
      alignment: 'center'
    }
  }
  
  if (!component.settings.fontSize) {
    component.settings.fontSize = 16
  }
  
  if (!component.settings.fontFamily) {
    component.settings.fontFamily = '宋体'
  }
  
  if (!component.settings.lineHeight) {
    component.settings.lineHeight = 1.5
  }
  
  if (!component.settings.labelStyle) {
    component.settings.labelStyle = {
      bold: true,
      alignment: 'left'
    }
  }
  
  if (!component.settings.valueStyle) {
    component.settings.valueStyle = {
      alignment: 'center'
    }
  }
  
  return component.settings
}

// 更新表格字段
const updateFields = () => {
  if (!props.component || !props.component.fields) return
  
  // 更新数据结构以匹配字段
  const newData = Array(props.component.data?.length || 0).fill().map(() => 
    Array(props.component.fields.length).fill('')
  )
  props.component.data = newData
  props.component.columns = props.component.fields.length
  updateProperty('fields', [...props.component.fields])
}

// 添加表格字段
const addField = () => {
  const newField = {
    id: Date.now(),
    label: '新字段',
    value: '',
    field: `field_${Date.now()}`,
    type: 'text',
    required: true,
    description: ''
  }
  
  props.component.fields.push(newField)
  updateFields()
}

// 删除表格字段
const removeField = (fieldId) => {
  const index = props.component.fields.findIndex(f => f.id === fieldId)
  if (index !== -1) {
    props.component.fields.splice(index, 1)
    updateFields()
  }
}

// 更新列宽
const updateColumnWidths = () => {
  // 确保第一列宽度在有效范围内
  const firstColumnWidth = Math.min(Math.max(props.component.columnWidths[0], 5), 95)
  // 第二列宽度 = 100% - 第一列宽度
  const secondColumnWidth = 100 - firstColumnWidth
  
  // 创建新的数组以确保触发响应式更新
  const newColumnWidths = [firstColumnWidth, secondColumnWidth]
  
  // 更新列宽数组
  updateProperty('columnWidths', newColumnWidths)
}

// 添加拖拽排序相关函数
const draggedField = ref(null)
const dragOverIndex = ref(null)

// 开始拖拽字段
const startDragField = (field) => {
  draggedField.value = field
}

// 拖拽结束
const endDragField = () => {
  draggedField.value = null
  dragOverIndex.value = null
}

// 拖拽经过某个位置
const dragOverField = (index) => {
  dragOverIndex.value = index
}

// 放置字段
const dropField = (index) => {
  if (!draggedField.value) return
  
  const fromIndex = props.component.fields.findIndex(f => f.id === draggedField.value.id)
  if (fromIndex === -1) return
  
  // 移除原位置的字段
  const [movedField] = props.component.fields.splice(fromIndex, 1)
  
  // 计算目标位置（考虑到移除后索引的变化）
  let targetIndex = index
  if (fromIndex < index) {
    targetIndex--
  }
  
  // 插入到新位置
  props.component.fields.splice(targetIndex, 0, movedField)
  
  // 更新字段
  updateFields()
  
  // 重置拖拽状态
  endDragField()
}

// 更新组件属性的辅助函数
const updateProperty = (property, value) => {
  emit('update', props.component.id, property, value)
}

// 字体列表
const fontFamilies = [
  { value: '宋体', label: '宋体' },
  { value: '黑体', label: '黑体' },
  { value: '楷体', label: '楷体' },
  { value: '仿宋', label: '仿宋' },
  { value: 'Arial', label: 'Arial' },
  { value: 'Times New Roman', label: 'Times New Roman' }
]

// 对齐方式
const alignments = [
  { value: 'left', label: '左对齐', icon: '⬅️' },
  { value: 'center', label: '居中', icon: '⬆️' },
  { value: 'right', label: '右对齐', icon: '➡️' },
  { value: 'justify', label: '两端对齐', icon: '↔️' }
]

// 添加子组件到字段内容
const addContentComponent = (fieldIndex, componentType) => {
  // 确保字段存在contentComponents数组
  if (!props.component.fields[fieldIndex].contentComponents) {
    props.component.fields[fieldIndex].contentComponents = []
  }
  
  // 创建新的子组件
  const newComponent = {
    id: `content_component_${Date.now()}`,
    type: componentType,
    // 根据组件类型设置默认属性
    ...getDefaultPropertiesForComponentType(componentType)
  }
  
  // 添加到字段的contentComponents数组
  props.component.fields[fieldIndex].contentComponents.push(newComponent)
  
  // 更新字段
  updateFields()
}

// 根据组件类型获取默认属性
const getDefaultPropertiesForComponentType = (type) => {
  switch (type) {
    case 'text':
      return {
        content: '文本内容',
        fontSize: 16,
        fontFamily: '宋体',
        alignment: 'left',
        lineHeight: 1.5
      }
    case 'list':
      return {
        items: [
          { id: 1, content: '列表项1' },
          { id: 2, content: '列表项2' }
        ],
        listType: 'bullet',
        fontSize: 16,
        fontFamily: '宋体',
        alignment: 'left'
      }
    case 'reference':
      return {
        items: [
          { id: 1, author: '作者', title: '标题', year: '年份', source: '来源' }
        ],
        style: 'apa',
        fontSize: 16,
        fontFamily: '宋体'
      }
    case 'image':
      return {
        url: '',
        alt: '图片描述',
        width: 80, // 百分比
        alignment: 'center',
        caption: '图片说明'
      }
    case 'title':
      return {
        content: '大标题',
        fontSize: 24,
        fontFamily: '黑体',
        alignment: 'center',
        bold: true
      }
    case 'paragraph':
      return {
        content: '普通文本段落内容',
        fontSize: 16,
        fontFamily: '宋体',
        alignment: 'left',
        lineHeight: 1.5
      }
    case 'schedule':
      return {
        items: [
          { task: '完成选题', startDate: '', endDate: '' },
          { task: '完成开题报告', startDate: '', endDate: '' },
          { task: '完成中期检查', startDate: '', endDate: '' }
        ],
        displayMode: 'table',
        title: '进度安排',
        settings: {
          fontSize: 16,
          fontFamily: '宋体',
          lineHeight: 1.5
        }
      }
    default:
      return {}
  }
}

// 删除字段中的子组件
const removeContentComponent = (fieldIndex, componentIndex) => {
  if (props.component.fields[fieldIndex].contentComponents) {
    props.component.fields[fieldIndex].contentComponents.splice(componentIndex, 1)
    updateFields()
  }
}

// 编辑子组件属性
const updateContentComponentProperty = (fieldIndex, componentIndex, property, value) => {
  if (props.component.fields[fieldIndex].contentComponents?.[componentIndex]) {
    props.component.fields[fieldIndex].contentComponents[componentIndex][property] = value
    updateFields()
  }
}

// 更新进度表组件的工作项目
const updateItems = (fieldIndex, componentIndex) => {
  if (props.component.fields[fieldIndex].contentComponents?.[componentIndex]) {
    const items = [...props.component.fields[fieldIndex].contentComponents[componentIndex].items]
    updateContentComponentProperty(fieldIndex, componentIndex, 'items', items)
  }
}

// 添加进度表工作项目
const addScheduleItem = (fieldIndex, componentIndex) => {
  if (props.component.fields[fieldIndex].contentComponents?.[componentIndex]) {
    const component = props.component.fields[fieldIndex].contentComponents[componentIndex]
    if (!component.items) {
      component.items = []
    }
    
    component.items.push({
      task: '新工作项',
      startDate: '',
      endDate: ''
    })
    
    updateItems(fieldIndex, componentIndex)
  }
}

// 删除进度表工作项目
const removeScheduleItem = (fieldIndex, componentIndex, itemIndex) => {
  if (props.component.fields[fieldIndex].contentComponents?.[componentIndex]) {
    const component = props.component.fields[fieldIndex].contentComponents[componentIndex]
    if (component.items && component.items.length > 1) {
      component.items.splice(itemIndex, 1)
      updateItems(fieldIndex, componentIndex)
    }
  }
}

// 设置当前编辑的字段
const setCurrentEditingField = (index) => {
  currentEditingFieldIndex.value = index
}

// 更新进度表组件的设置
const updateScheduleSettings = (fieldIndex, componentIndex, property, value) => {
  if (props.component.fields[fieldIndex].contentComponents?.[componentIndex]) {
    const component = props.component.fields[fieldIndex].contentComponents[componentIndex]
    if (!component.settings) {
      component.settings = {}
    }
    
    component.settings[property] = value
    updateContentComponentProperty(fieldIndex, componentIndex, 'settings', { ...component.settings })
  }
}
</script>

<template>
  <!-- 表格预览切换按钮 -->
  <div class="mb-4 flex justify-end">
    <button 
      @click="toggleTablePreview" 
      class="px-3 py-2 bg-blue-100 text-blue-700 rounded hover:bg-blue-200 flex items-center gap-1"
    >
      <span v-if="!showTablePreview">👁️ 预览表格</span>
      <span v-else>✏️ 返回编辑</span>
    </button>
  </div>

  <!-- 表格预览模式 -->
  <div v-if="showTablePreview" class="mb-4 p-4 border rounded bg-white">
    <div 
      v-if="component.title" 
      class="font-bold mb-2" 
      :style="{ 
        textAlign: component.titleAlignment || 'center',
        fontFamily: component.titleFontFamily || '黑体',
        fontSize: `${component.titleFontSize || 18}px`
      }"
    >
      {{ component.title }}
    </div>
    <table class="w-full border-collapse">
      <!-- 数据行 - 键值对布局 -->
      <template v-if="component.layout === 'key-value'">
        <tr v-for="(field, index) in component.fields" :key="field.id">
          <td
            class="p-2"
            :style="{ 
              borderWidth: ensureComponentSettings(component).borderStyle?.left ? '1px' : '0',
              borderStyle: 'solid',
              borderColor: '#000',
              fontSize: `${ensureComponentSettings(component).fontSize || 16}px`,
              fontFamily: ensureComponentSettings(component).fontFamily || '宋体',
              textAlign: ensureComponentSettings(component).cellStyle?.alignment || 'left',
              fontWeight: 'bold'
            }"
          >
            {{ field.label }}
          </td>
          <td
            class="p-2"
            :style="{ 
              borderWidth: ensureComponentSettings(component).borderStyle?.right ? '1px' : '0',
              borderStyle: 'solid',
              borderColor: '#000',
              fontSize: `${ensureComponentSettings(component).fontSize || 16}px`,
              fontFamily: ensureComponentSettings(component).fontFamily || '宋体',
              textAlign: ensureComponentSettings(component).cellStyle?.alignment || 'left'
            }"
          >
            {{ field.value || '点击编辑' }}
          </td>
        </tr>
      </template>
      
      <!-- 数据行 - 标题内容布局 -->
      <template v-else>
        <tr v-for="field in component.fields" :key="field.id">
          <td
            class="p-2"
            :style="{ 
              borderWidth: ensureComponentSettings(component).borderStyle?.left ? '1px' : '0',
              borderStyle: 'solid',
              borderColor: '#000',
              fontSize: `${ensureComponentSettings(component).fontSize || 16}px`,
              fontFamily: ensureComponentSettings(component).fontFamily || '宋体'
            }"
          >
            <div class="font-bold">{{ field.label }}</div>
            <div>{{ field.value || '点击编辑' }}</div>
            <div v-if="field.description" class="text-gray-500 text-sm">{{ field.description }}</div>
          </td>
        </tr>
      </template>
    </table>
  </div>

  <!-- 表格编辑模式 -->
  <div v-else>
    <div class="mb-4">
      <label class="block mb-2">表格标题</label>
      <input
        v-model="component.title"
        type="text"
        class="w-full px-3 py-2 border rounded"
        @input="updateProperty('title', $event.target.value)"
      />
    </div>

    <!-- 表格标题对齐方式 -->
    <div class="mb-4">
      <label class="block mb-2">标题对齐方式</label>
      <div class="flex gap-2">
        <button
          v-for="align in alignments"
          :key="align.value"
          :class="[
            'px-3 py-2 border rounded flex-1',
            component.titleAlignment === align.value ? 'bg-blue-500 text-white' : 'bg-white'
          ]"
          @click="updateProperty('titleAlignment', align.value)"
        >
          {{ align.icon }}
        </button>
      </div>
    </div>

    <!-- 表格标题字体和字号 -->
    <div class="grid grid-cols-2 gap-4 mb-4">
      <div>
        <label class="block mb-2">标题字体</label>
        <select
          v-model="component.titleFontFamily"
          class="w-full px-3 py-2 border rounded"
          @change="updateProperty('titleFontFamily', $event.target.value)"
        >
          <option v-for="font in fontFamilies" :key="font.value" :value="font.value">
            {{ font.label }}
          </option>
        </select>
      </div>
      <div>
        <label class="block mb-2">标题字号</label>
        <div class="flex items-center">
          <input
            v-model.number="component.titleFontSize"
            type="range"
            class="flex-1 mr-2"
            min="12"
            max="36"
            step="1"
            @input="updateProperty('titleFontSize', $event.target.value)"
          />
          <span class="w-8 text-center">{{ component.titleFontSize || 18 }}</span>
        </div>
      </div>
    </div>

    <!-- 表格布局设置 -->
    <div class="mb-4">
      <label class="block mb-2">表格布局</label>
      <div class="grid grid-cols-2 gap-4">
        <button
          :class="[
            'p-3 border rounded flex flex-col items-center',
            component.layout === 'key-value' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('layout', 'key-value')"
        >
          <div class="text-2xl mb-1">📋</div>
          <div class="text-sm">键值对布局</div>
        </button>
        <button
          :class="[
            'p-3 border rounded flex flex-col items-center',
            component.layout === 'title-content' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('layout', 'title-content')"
        >
          <div class="text-2xl mb-1">📝</div>
          <div class="text-sm">标题内容布局</div>
        </button>
      </div>
    </div>

    <!-- 列宽设置 (仅键值对布局) -->
    <div class="mb-4" v-if="component.layout === 'key-value'">
      <label class="block mb-2">列宽设置</label>
      <div class="space-y-4">
        <!-- 确保columnWidths存在 -->
        <div v-if="!component.columnWidths" class="hidden">
          {{ updateProperty('columnWidths', [50, 50]) }}
        </div>
        
        <!-- 列宽预览 -->
        <div class="flex mb-2 h-8 bg-gray-100 rounded overflow-hidden">
          <div 
            class="bg-blue-200 flex items-center justify-center text-xs font-medium"
            :style="{ width: `${component.columnWidths?.[0] || 50}%` }"
          >
            标签列: {{ component.columnWidths?.[0] || 50 }}%
          </div>
          <div 
            class="bg-green-200 flex items-center justify-center text-xs font-medium"
            :style="{ width: `${component.columnWidths?.[1] || 50}%` }"
          >
            内容列: {{ component.columnWidths?.[1] || 50 }}%
          </div>
        </div>
        
        <!-- 第一列宽度滑块 -->
        <div>
          <div class="flex justify-between mb-1">
            <span class="text-sm">拖动调整列宽比例</span>
          </div>
          <input
            v-model.number="component.columnWidths[0]"
            type="range"
            class="w-full"
            min="5"
            max="95"
            step="5"
            @input="updateColumnWidths"
          />
          <div class="flex justify-between text-xs text-gray-500">
            <span>标签列窄 (5%)</span>
            <span>平均 (50%)</span>
            <span>内容列窄 (95%)</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 表格字段设置 -->
    <div class="mb-4">
      <div class="flex justify-between items-center mb-2">
        <label class="block">表格字段</label>
        <button
          @click="addField"
          class="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 flex items-center gap-1"
        >
          <span>+</span> 添加字段
        </button>
      </div>
      
      <div class="space-y-2">
        <div 
          v-for="(field, index) in component.fields" 
          :key="field.id" 
          class="p-3 border rounded bg-white"
          :class="{ 'border-blue-300 bg-blue-50': draggedField && draggedField.id === field.id }"
          draggable="true"
          @dragstart="startDragField(field)"
          @dragend="endDragField"
          @dragover.prevent="dragOverField(index)"
          @drop.prevent="dropField(index)"
        >
          <div class="flex items-center mb-2">
            <div class="mr-2 cursor-move text-gray-400 hover:text-gray-600">⋮⋮</div>
            <input
              v-model="field.label"
              type="text"
              class="flex-1 px-3 py-2 border rounded"
              placeholder="字段标签"
              @input="updateFields"
            />
            <select
              v-model="field.type"
              class="ml-2 px-3 py-2 border rounded"
              @change="updateFields"
            >
              <option value="text">文本</option>
              <option value="number">数字</option>
              <option value="date">日期</option>
            </select>
            <button
              @click="removeField(field.id)"
              class="ml-2 px-2 py-2 text-red-500 hover:bg-red-100 rounded"
            >
              🗑️
            </button>
          </div>
          
          <!-- 字段值和描述 -->
          <div class="space-y-2 pl-6">
            <div>
              <label class="block mb-1 text-sm">默认值</label>
              <input
                v-model="field.value"
                type="text"
                class="w-full px-3 py-2 border rounded"
                placeholder="默认值"
                @input="updateFields"
              />
            </div>
            <div>
              <label class="block mb-1 text-sm">描述</label>
              <input
                v-model="field.description"
                type="text"
                class="w-full px-3 py-2 border rounded"
                placeholder="字段描述"
                @input="updateFields"
              />
            </div>
            
            <!-- 内容组件管理 -->
            <div class="mt-3 border-t pt-3">
              <div class="flex justify-between items-center mb-2">
                <label class="block text-sm font-medium">内容组件</label>
                <button
                  @click="setCurrentEditingField(index)"
                  class="px-2 py-1 bg-blue-100 text-blue-700 rounded hover:bg-blue-200 text-sm"
                >
                  {{ currentEditingFieldIndex === index ? '完成编辑' : '编辑组件' }}
                </button>
              </div>
              
              <!-- 子组件编辑区域 -->
              <div v-if="currentEditingFieldIndex === index" class="bg-gray-50 p-3 rounded border">
                <!-- 添加组件按钮 -->
                <div class="mb-3">
                  <label class="block mb-1 text-sm">添加内容组件</label>
                  <div class="flex flex-wrap gap-2">
                    <button
                      v-for="compType in availableContentComponents"
                      :key="compType.type"
                      @click="addContentComponent(index, compType.type)"
                      class="px-2 py-1 bg-white border rounded hover:bg-gray-100 flex items-center gap-1"
                    >
                      <span>{{ compType.icon }}</span>
                      <span>{{ compType.name }}</span>
                    </button>
                  </div>
                </div>
                
                <!-- 已添加的组件列表 -->
                <div v-if="field.contentComponents && field.contentComponents.length > 0">
                  <label class="block mb-1 text-sm">已添加组件</label>
                  <div class="space-y-2">
                    <div 
                      v-for="(component, compIndex) in field.contentComponents" 
                      :key="component.id"
                      class="p-2 bg-white border rounded"
                    >
                      <div class="flex justify-between items-center mb-2">
                        <div class="flex items-center gap-1">
                          <span v-if="component.type === 'text'">📄</span>
                          <span v-else-if="component.type === 'list'">📋</span>
                          <span v-else-if="component.type === 'reference'">📚</span>
                          <span v-else-if="component.type === 'image'">🖼️</span>
                          <span class="font-medium">{{ component.type === 'text' ? '文本段落' : 
                                                      component.type === 'list' ? '列表' : 
                                                      component.type === 'reference' ? '参考文献' : 
                                                      component.type === 'image' ? '图片' : 
                                                      component.type }}</span>
                        </div>
                        <button
                          @click="removeContentComponent(index, compIndex)"
                          class="px-2 py-1 text-red-500 hover:bg-red-100 rounded"
                        >
                          🗑️
                        </button>
                      </div>
                      
                      <!-- 文本组件编辑 -->
                      <div v-if="component.type === 'text'" class="pl-4 space-y-2">
                        <div>
                          <label class="block mb-1 text-xs">内容</label>
                          <textarea
                            v-model="component.content"
                            class="w-full px-2 py-1 border rounded text-sm"
                            rows="3"
                            @input="updateContentComponentProperty(index, compIndex, 'content', component.content)"
                          ></textarea>
                        </div>
                        <div class="flex gap-2">
                          <div class="flex-1">
                            <label class="block mb-1 text-xs">字体</label>
                            <select
                              v-model="component.fontFamily"
                              class="w-full px-2 py-1 border rounded text-sm"
                              @change="updateContentComponentProperty(index, compIndex, 'fontFamily', component.fontFamily)"
                            >
                              <option v-for="font in fontFamilies" :key="font.value" :value="font.value">
                                {{ font.label }}
                              </option>
                            </select>
                          </div>
                          <div class="w-24">
                            <label class="block mb-1 text-xs">字号</label>
                            <input
                              v-model.number="component.fontSize"
                              type="number"
                              min="8"
                              max="72"
                              class="w-full px-2 py-1 border rounded text-sm"
                              @input="updateContentComponentProperty(index, compIndex, 'fontSize', component.fontSize)"
                            />
                          </div>
                        </div>
                      </div>
                      
                      <!-- 列表组件编辑 -->
                      <div v-else-if="component.type === 'list'" class="pl-4 space-y-2">
                        <div class="flex gap-2 mb-2">
                          <div class="flex-1">
                            <label class="block mb-1 text-xs">列表类型</label>
                            <select
                              v-model="component.listType"
                              class="w-full px-2 py-1 border rounded text-sm"
                              @change="updateContentComponentProperty(index, compIndex, 'listType', component.listType)"
                            >
                              <option value="bullet">无序列表</option>
                              <option value="number">有序列表</option>
                            </select>
                          </div>
                          <div class="w-24">
                            <label class="block mb-1 text-xs">字号</label>
                            <input
                              v-model.number="component.fontSize"
                              type="number"
                              min="8"
                              max="72"
                              class="w-full px-2 py-1 border rounded text-sm"
                              @input="updateContentComponentProperty(index, compIndex, 'fontSize', component.fontSize)"
                            />
                          </div>
                        </div>
                        
                        <div>
                          <label class="block mb-1 text-xs">列表项</label>
                          <div v-for="(item, itemIndex) in component.items" :key="item.id" class="flex gap-1 mb-1">
                            <input
                              v-model="item.content"
                              type="text"
                              class="flex-1 px-2 py-1 border rounded text-sm"
                              placeholder="列表项内容"
                              @input="updateFields"
                            />
                            <button
                              @click="component.items.splice(itemIndex, 1); updateFields()"
                              class="px-2 text-red-500 hover:bg-red-100 rounded text-sm"
                              :disabled="component.items.length <= 1"
                            >
                              -
                            </button>
                          </div>
                          <button
                            @click="component.items.push({id: Date.now(), content: '新列表项'}); updateFields()"
                            class="mt-1 px-2 py-1 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 text-sm w-full"
                          >
                            + 添加列表项
                          </button>
                        </div>
                      </div>
                      
                      <!-- 参考文献组件编辑 -->
                      <div v-else-if="component.type === 'reference'" class="pl-4 space-y-2">
                        <div class="flex gap-2 mb-2">
                          <div class="flex-1">
                            <label class="block mb-1 text-xs">引用样式</label>
                            <select
                              v-model="component.style"
                              class="w-full px-2 py-1 border rounded text-sm"
                              @change="updateContentComponentProperty(index, compIndex, 'style', component.style)"
                            >
                              <option value="apa">APA</option>
                              <option value="mla">MLA</option>
                              <option value="chicago">Chicago</option>
                              <option value="gb">GB/T 7714</option>
                            </select>
                          </div>
                          <div class="w-24">
                            <label class="block mb-1 text-xs">字号</label>
                            <input
                              v-model.number="component.fontSize"
                              type="number"
                              min="8"
                              max="72"
                              class="w-full px-2 py-1 border rounded text-sm"
                              @input="updateContentComponentProperty(index, compIndex, 'fontSize', component.fontSize)"
                            />
                          </div>
                        </div>
                        
                        <div class="mb-2">
                          <label class="block mb-1 text-xs">字体</label>
                          <select
                            v-model="component.fontFamily"
                            class="w-full px-2 py-1 border rounded text-sm"
                            @change="updateContentComponentProperty(index, compIndex, 'fontFamily', component.fontFamily)"
                          >
                            <option v-for="font in fontFamilies" :key="font.value" :value="font.value">
                              {{ font.label }}
                            </option>
                          </select>
                        </div>
                        
                        <div>
                          <label class="block mb-1 text-xs">参考文献</label>
                          <div v-for="(item, itemIndex) in component.items" :key="item.id" class="border p-2 rounded mb-2">
                            <div class="grid grid-cols-2 gap-1 mb-1">
                              <div>
                                <label class="block mb-1 text-xs">作者</label>
                                <input
                                  v-model="item.author"
                                  type="text"
                                  class="w-full px-2 py-1 border rounded text-sm"
                                  placeholder="作者"
                                  @input="updateFields"
                                />
                              </div>
                              <div>
                                <label class="block mb-1 text-xs">年份</label>
                                <input
                                  v-model="item.year"
                                  type="text"
                                  class="w-full px-2 py-1 border rounded text-sm"
                                  placeholder="年份"
                                  @input="updateFields"
                                />
                              </div>
                            </div>
                            <div class="mb-1">
                              <label class="block mb-1 text-xs">标题</label>
                              <input
                                v-model="item.title"
                                type="text"
                                class="w-full px-2 py-1 border rounded text-sm"
                                placeholder="标题"
                                @input="updateFields"
                              />
                            </div>
                            <div class="mb-1">
                              <label class="block mb-1 text-xs">来源</label>
                              <input
                                v-model="item.source"
                                type="text"
                                class="w-full px-2 py-1 border rounded text-sm"
                                placeholder="来源"
                                @input="updateFields"
                              />
                            </div>
                            <div class="flex justify-end">
                              <button
                                @click="component.items.splice(itemIndex, 1); updateFields()"
                                class="px-2 py-1 text-red-500 hover:bg-red-100 rounded text-sm"
                                :disabled="component.items.length <= 1"
                              >
                                删除
                              </button>
                            </div>
                          </div>
                          <button
                            @click="component.items.push({id: Date.now(), author: '', title: '', year: '', source: ''}); updateFields()"
                            class="mt-1 px-2 py-1 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 text-sm w-full"
                          >
                            + 添加参考文献
                          </button>
                        </div>
                      </div>
                      
                      <!-- 图片组件编辑 -->
                      <div v-else-if="component.type === 'image'" class="p-3 border rounded mb-2">
                        <h4 class="font-bold mb-2">图片组件</h4>
                        <div class="mb-2">
                          <label class="block mb-1">图片URL</label>
                          <input 
                            v-model="component.url" 
                            type="text" 
                            class="w-full px-2 py-1 border rounded"
                            @input="updateContentComponentProperty(index, compIndex, 'url', $event.target.value)"
                          />
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">图片说明</label>
                          <input 
                            v-model="component.caption" 
                            type="text" 
                            class="w-full px-2 py-1 border rounded"
                            @input="updateContentComponentProperty(index, compIndex, 'caption', $event.target.value)"
                          />
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">对齐方式</label>
                          <select 
                            v-model="component.alignment" 
                            class="w-full px-2 py-1 border rounded"
                            @change="updateContentComponentProperty(index, compIndex, 'alignment', $event.target.value)"
                          >
                            <option v-for="align in alignments" :key="align.value" :value="align.value">
                              {{ align.label }}
                            </option>
                          </select>
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">宽度 (%)</label>
                          <input 
                            v-model.number="component.width" 
                            type="number" 
                            min="10" 
                            max="100" 
                            class="w-full px-2 py-1 border rounded"
                            @input="updateContentComponentProperty(index, compIndex, 'width', $event.target.value)"
                          />
                        </div>
                      </div>
                      
                      <!-- 大标题组件编辑 -->
                      <div v-else-if="component.type === 'title'" class="p-3 border rounded mb-2">
                        <h4 class="font-bold mb-2">大标题组件</h4>
                        <div class="mb-2">
                          <label class="block mb-1">标题内容</label>
                          <input 
                            v-model="component.content" 
                            type="text" 
                            class="w-full px-2 py-1 border rounded"
                            @input="updateContentComponentProperty(index, compIndex, 'content', $event.target.value)"
                          />
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">字体</label>
                          <select 
                            v-model="component.fontFamily" 
                            class="w-full px-2 py-1 border rounded"
                            @change="updateContentComponentProperty(index, compIndex, 'fontFamily', $event.target.value)"
                          >
                            <option v-for="font in fontFamilies" :key="font.value" :value="font.value">
                              {{ font.label }}
                            </option>
                          </select>
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">字号</label>
                          <input 
                            v-model.number="component.fontSize" 
                            type="number" 
                            min="12" 
                            max="72" 
                            class="w-full px-2 py-1 border rounded"
                            @input="updateContentComponentProperty(index, compIndex, 'fontSize', $event.target.value)"
                          />
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">对齐方式</label>
                          <select 
                            v-model="component.alignment" 
                            class="w-full px-2 py-1 border rounded"
                            @change="updateContentComponentProperty(index, compIndex, 'alignment', $event.target.value)"
                          >
                            <option v-for="align in alignments" :key="align.value" :value="align.value">
                              {{ align.label }}
                            </option>
                          </select>
                        </div>
                        <div class="mb-2">
                          <label class="flex items-center">
                            <input 
                              v-model="component.bold" 
                              type="checkbox" 
                              class="mr-2"
                              @change="updateContentComponentProperty(index, compIndex, 'bold', $event.target.checked)"
                            />
                            加粗
                          </label>
                        </div>
                      </div>
                      
                      <!-- 普通文本段落组件编辑 -->
                      <div v-else-if="component.type === 'paragraph'" class="p-3 border rounded mb-2">
                        <h4 class="font-bold mb-2">文本段落组件</h4>
                        <div class="mb-2">
                          <label class="block mb-1">段落内容</label>
                          <textarea 
                            v-model="component.content" 
                            rows="3"
                            class="w-full px-2 py-1 border rounded"
                            @input="updateContentComponentProperty(index, compIndex, 'content', $event.target.value)"
                          ></textarea>
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">字体</label>
                          <select 
                            v-model="component.fontFamily" 
                            class="w-full px-2 py-1 border rounded"
                            @change="updateContentComponentProperty(index, compIndex, 'fontFamily', $event.target.value)"
                          >
                            <option v-for="font in fontFamilies" :key="font.value" :value="font.value">
                              {{ font.label }}
                            </option>
                          </select>
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">字号</label>
                          <input 
                            v-model.number="component.fontSize" 
                            type="number" 
                            min="12" 
                            max="72" 
                            class="w-full px-2 py-1 border rounded"
                            @input="updateContentComponentProperty(index, compIndex, 'fontSize', $event.target.value)"
                          />
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">对齐方式</label>
                          <select 
                            v-model="component.alignment" 
                            class="w-full px-2 py-1 border rounded"
                            @change="updateContentComponentProperty(index, compIndex, 'alignment', $event.target.value)"
                          >
                            <option v-for="align in alignments" :key="align.value" :value="align.value">
                              {{ align.label }}
                            </option>
                          </select>
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">行高</label>
                          <input 
                            v-model.number="component.lineHeight" 
                            type="number" 
                            min="1" 
                            max="3" 
                            step="0.1"
                            class="w-full px-2 py-1 border rounded"
                            @input="updateContentComponentProperty(index, compIndex, 'lineHeight', $event.target.value)"
                          />
                        </div>
                      </div>
                      
                      <!-- 进度表组件编辑 -->
                      <div v-else-if="component.type === 'schedule'" class="p-3 border rounded mb-2">
                        <h4 class="font-bold mb-2">进度表组件</h4>
                        <div class="mb-2">
                          <label class="block mb-1">标题</label>
                          <input 
                            v-model="component.title" 
                            type="text" 
                            class="w-full px-2 py-1 border rounded"
                            @input="updateContentComponentProperty(index, compIndex, 'title', $event.target.value)"
                          />
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">工作项目</label>
                          <div v-for="(item, itemIndex) in component.items" :key="itemIndex" class="mb-2 p-2 border rounded">
                            <div class="mb-1">
                              <label class="block text-xs mb-1">任务名称</label>
                              <input 
                                v-model="item.task" 
                                type="text" 
                                class="w-full px-2 py-1 border rounded"
                                @input="updateItems(index, compIndex)"
                              />
                            </div>
                            <div class="grid grid-cols-2 gap-2 mb-1">
                              <div>
                                <label class="block text-xs mb-1">开始日期</label>
                                <input 
                                  v-model="item.startDate" 
                                  type="date" 
                                  class="w-full px-2 py-1 border rounded"
                                  @input="updateItems(index, compIndex)"
                                />
                              </div>
                              <div>
                                <label class="block text-xs mb-1">结束日期</label>
                                <input 
                                  v-model="item.endDate" 
                                  type="date" 
                                  class="w-full px-2 py-1 border rounded"
                                  @input="updateItems(index, compIndex)"
                                />
                              </div>
                            </div>
                            <button 
                              @click="removeScheduleItem(index, compIndex, itemIndex)" 
                              class="px-2 py-1 bg-red-500 text-white rounded text-xs"
                              :disabled="component.items.length <= 1"
                            >
                              删除
                            </button>
                          </div>
                          <button 
                            @click="addScheduleItem(index, compIndex)" 
                            class="w-full px-2 py-1 bg-blue-500 text-white rounded text-sm"
                          >
                            添加工作项
                          </button>
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">字体</label>
                          <select 
                            :value="component.settings?.fontFamily || '宋体'" 
                            class="w-full px-2 py-1 border rounded"
                            @change="updateScheduleSettings(index, compIndex, 'fontFamily', $event.target.value)"
                          >
                            <option v-for="font in fontFamilies" :key="font.value" :value="font.value">
                              {{ font.label }}
                            </option>
                          </select>
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">字号</label>
                          <input 
                            :value="component.settings?.fontSize || 16" 
                            type="number" 
                            min="12" 
                            max="72" 
                            class="w-full px-2 py-1 border rounded"
                            @input="updateScheduleSettings(index, compIndex, 'fontSize', parseInt($event.target.value))"
                          />
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">行高</label>
                          <input 
                            :value="component.settings?.lineHeight || 1.5" 
                            type="number" 
                            min="1" 
                            max="3" 
                            step="0.1"
                            class="w-full px-2 py-1 border rounded"
                            @input="updateScheduleSettings(index, compIndex, 'lineHeight', parseFloat($event.target.value))"
                          />
                        </div>
                        <div class="mb-2">
                          <label class="block mb-1">显示方式</label>
                          <select 
                            v-model="component.displayMode" 
                            class="w-full px-2 py-1 border rounded"
                            @change="updateContentComponentProperty(index, compIndex, 'displayMode', $event.target.value)"
                          >
                            <option value="table">表格模式</option>
                            <option value="gantt">甘特图模式</option>
                          </select>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                
                <div v-else class="text-center py-2 text-gray-500 text-sm">
                  暂无内容组件，请点击上方按钮添加
                </div>
              </div>
              
              <!-- 组件数量提示 -->
              <div v-else class="text-sm text-gray-500">
                {{ field.contentComponents?.length ? `已添加 ${field.contentComponents.length} 个组件` : '暂无内容组件' }}
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 拖拽提示 -->
      <div class="mt-2 text-xs text-gray-500">
        提示：可拖拽字段调整顺序
      </div>
    </div>

    <!-- 表格样式设置 -->
    <div class="mb-4">
      <label class="block mb-2">表格样式</label>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block mb-1">字体</label>
          <select
            v-model="ensureComponentSettings(component).fontFamily"
            class="w-full px-3 py-2 border rounded"
            @change="updateProperty('settings', {...component.settings})"
          >
            <option v-for="font in fontFamilies" :key="font.value" :value="font.value">
              {{ font.label }}
            </option>
          </select>
        </div>
        <div>
          <label class="block mb-1">字号</label>
          <input
            v-model.number="ensureComponentSettings(component).fontSize"
            type="number"
            class="w-full px-3 py-2 border rounded"
            min="12"
            max="72"
            @input="updateProperty('settings', {...component.settings})"
          />
        </div>
      </div>
    </div>

    <!-- 表格边框设置 -->
    <div class="mb-4">
      <label class="block mb-2">边框样式</label>
      
      <!-- 边框预设 -->
      <div class="flex flex-wrap gap-2 mb-3">
        <button
          v-for="preset in borderPresets"
          :key="preset.name"
          class="px-3 py-2 border rounded hover:bg-gray-100 flex flex-col items-center"
          @click="applyBorderPreset(preset.preset)"
        >
          <span class="text-xl">{{ preset.icon }}</span>
          <span class="text-xs mt-1">{{ preset.name }}</span>
        </button>
      </div>
      
      <!-- 自定义边框设置 -->
      <div class="grid grid-cols-3 gap-2 p-3 border rounded bg-gray-50">
        <label class="flex items-center">
          <input
            v-model="ensureComponentSettings(component).borderStyle.top"
            type="checkbox"
            class="mr-2"
            @change="updateProperty('settings', {...component.settings})"
          />
          上边框
        </label>
        <label class="flex items-center">
          <input
            v-model="ensureComponentSettings(component).borderStyle.bottom"
            type="checkbox"
            class="mr-2"
            @change="updateProperty('settings', {...component.settings})"
          />
          下边框
        </label>
        <label class="flex items-center">
          <input
            v-model="ensureComponentSettings(component).borderStyle.left"
            type="checkbox"
            class="mr-2"
            @change="updateProperty('settings', {...component.settings})"
          />
          左边框
        </label>
        <label class="flex items-center">
          <input
            v-model="ensureComponentSettings(component).borderStyle.right"
            type="checkbox"
            class="mr-2"
            @change="updateProperty('settings', {...component.settings})"
          />
          右边框
        </label>
        <label class="flex items-center">
          <input
            v-model="ensureComponentSettings(component).borderStyle.insideHorizontal"
            type="checkbox"
            class="mr-2"
            @change="updateProperty('settings', {...component.settings})"
          />
          内部横线
        </label>
        <label class="flex items-center">
          <input
            v-model="ensureComponentSettings(component).borderStyle.insideVertical"
            type="checkbox"
            class="mr-2"
            @change="updateProperty('settings', {...component.settings})"
          />
          内部竖线
        </label>
      </div>
    </div>
  </div>
</template> 