<script setup>
import { ref, computed, watch } from 'vue'
import WorkUnit from '../../../../components/WorkUnit.vue'

const props = defineProps({
  template: {
    type: Object,
    required: true
  },
  isDragging: {
    type: Boolean,
    default: false
  },
  selectedComponentId: {
    type: [String, Number],
    default: null
  },
  // 添加属性表示是否有未保存的更改
  hasUnsavedChanges: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:selectedComponentId', 'component-drop', 'component-update', 'component-delete', 'component-reorder', 'template-changed'])

// 拖拽状态管理
const draggedComponentId = ref(null)
const dragOverIndex = ref(-1)
const isDraggingComponent = ref(false)

// 计算A4纸张样式
const pageStyle = computed(() => {
  // 添加默认值
  const defaultSettings = {
    pageSize: {
      width: 21,
      height: 29.7
    },
    margins: {
      top: 2.54,
      bottom: 2.54,
      left: 3.18,
      right: 3.18
    }
  }

  // 使用可选链操作符和空值合并操作符来安全地访问属性
  const settings = props.template?.documentSettings ?? defaultSettings
  const { width, height } = settings.pageSize
  const { top, bottom, left, right } = settings.margins
  
  return {
    width: `${width}cm`,
    minHeight: `${height}cm`,
    padding: `${top}cm ${right}cm ${bottom}cm ${left}cm`
  }
})

// 监听组件内部变化，通知父组件内容已更改
const emitTemplateChanged = () => {
  emit('template-changed')
}

// 处理组件拖放
const handleDrop = (e) => {
  e.preventDefault()
  const componentType = e.dataTransfer.getData('componentType')
  
  // 如果是从组件面板拖入的新组件
  if (componentType && !e.dataTransfer.getData('componentId')) {
    // 获取画布元素
    const canvas = e.currentTarget
    const rect = canvas.getBoundingClientRect()
    
    // 计算相对于画布的位置
    const position = {
      x: e.clientX - rect.left,
      y: e.clientY - rect.top
    }
    
    emit('component-drop', componentType, position)
    // 发出内容已更改事件
    emitTemplateChanged()
    return
  }
  
  // 如果是已有组件的重新排序
  const componentId = e.dataTransfer.getData('componentId')
  if (componentId && dragOverIndex.value !== -1) {
    // 找到拖拽组件的当前索引
    const currentIndex = props.template.components.findIndex(c => c.id === Number(componentId))
    if (currentIndex !== -1) {
      // 发出重新排序事件
      emit('component-reorder', Number(componentId), currentIndex, dragOverIndex.value)
      // 发出内容已更改事件
      emitTemplateChanged()
    }
  }
  
  // 重置拖拽状态
  draggedComponentId.value = null
  dragOverIndex.value = -1
  isDraggingComponent.value = false
}

const handleDragOver = (e) => {
  e.preventDefault()
  
  // 如果没有正在拖拽的组件，不处理
  if (!isDraggingComponent.value) return
  
  // 获取鼠标位置，计算应该插入的位置
  const canvas = e.currentTarget
  const rect = canvas.getBoundingClientRect()
  const mouseY = e.clientY - rect.top
  
  // 获取所有组件元素
  const componentElements = Array.from(canvas.querySelectorAll('.component-wrapper'))
  
  // 如果没有组件，设置为第一个位置
  if (componentElements.length === 0) {
    dragOverIndex.value = 0
    return
  }
  
  // 找到鼠标位置对应的组件索引
  for (let i = 0; i < componentElements.length; i++) {
    const element = componentElements[i]
    const elementRect = element.getBoundingClientRect()
    const elementTop = elementRect.top - rect.top
    const elementBottom = elementRect.bottom - rect.top
    const elementMiddle = (elementTop + elementBottom) / 2
    
    // 如果鼠标在组件上半部分，插入到该组件前面
    if (mouseY < elementMiddle) {
      dragOverIndex.value = i
      return
    }
  }
  
  // 如果鼠标在所有组件下方，插入到最后
  dragOverIndex.value = componentElements.length
}

// 处理组件拖拽开始
const handleDragStart = (e, componentId, index) => {
  // 设置拖拽数据
  e.dataTransfer.setData('componentId', componentId)
  e.dataTransfer.effectAllowed = 'move'
  
  // 更新拖拽状态
  draggedComponentId.value = componentId
  isDraggingComponent.value = true
  
  // 设置拖拽时的视觉效果
  setTimeout(() => {
    const element = document.querySelector(`.component-wrapper[data-id="${componentId}"]`)
    if (element) {
      element.classList.add('dragging')
    }
  }, 0)
}

// 处理组件拖拽结束
const handleDragEnd = () => {
  // 重置拖拽状态
  const element = document.querySelector('.component-wrapper.dragging')
  if (element) {
    element.classList.remove('dragging')
  }
  
  draggedComponentId.value = null
  dragOverIndex.value = -1
  isDraggingComponent.value = false
}

const handleDragEnter = (e, index) => {
  if (isDraggingComponent.value) {
    dragOverIndex.value = index
  }
}

// 选择组件
const selectComponent = (componentId) => {
  emit('update:selectedComponentId', componentId)
}

// 包装更新组件属性的事件，同时通知有内容更改
const handleComponentUpdate = (componentId, property, value) => {
  emit('component-update', componentId, property, value)
  emitTemplateChanged()
}

// 包装删除组件的事件，同时通知有内容更改
const handleComponentDelete = (componentId) => {
  emit('component-delete', componentId)
  emitTemplateChanged()
}

// 计算组件样式
const getComponentStyle = (component) => {
  if (!component) return {}
  
  const baseStyle = {
    textAlign: component.alignment || 'left',
    fontSize: `${component.fontSize || 16}px`,
    fontFamily: component.fontFamily || '宋体'
  }
  
  if (component.type === 'title') {
    return {
      ...baseStyle,
      fontWeight: component.bold ? 'bold' : 'normal'
    }
  } else if (component.type === 'text') {
    return {
      ...baseStyle,
      lineHeight: component.lineHeight || 1.5
    }
  } else if (component.type === 'toc') {
    return {
      ...baseStyle,
      lineHeight: component.lineHeight || 1.8
    }
  } else if (component.type === 'image') {
    return {
      ...baseStyle,
      textAlign: component.captionAlignment || 'center',
      lineHeight: component.captionLineHeight || 1.5
    }
  } else if (component.type === 'reference') {
    return {
      ...baseStyle,
      lineHeight: component.lineHeight || 1.5
    }
  } else if (component.type === 'list') {
    return {
      ...baseStyle,
      lineHeight: component.lineHeight || 1.5
    }
  }
  
  return baseStyle
}

// 计算组件包装器的类名
const getComponentWrapperClass = (component, index) => {
  const classes = ['component-wrapper']
  
  // 如果是被拖拽的组件
  if (component.id === Number(draggedComponentId.value)) {
    classes.push('dragging')
  }
  
  // 如果是拖拽目标位置的前一个组件
  if (dragOverIndex.value === index && Number(draggedComponentId.value) !== component.id) {
    classes.push('drag-over-top')
  }
  
  // 如果是拖拽目标位置的后一个组件
  if (dragOverIndex.value === index + 1 && Number(draggedComponentId.value) !== component.id) {
    classes.push('drag-over-bottom')
  }
  
  return classes.join(' ')
}

// 添加一个新函数：处理画布点击事件
const handleCanvasClick = (e) => {
  // 检查点击的是否是画布本身（不是组件）
  // 使用e.target和e.currentTarget比较
  if (e.target === e.currentTarget) {
    // 如果点击的是画布空白区域，取消组件选中状态
    emit('update:selectedComponentId', null)
  }
}
</script>

<template>
  <div 
    class="template-canvas bg-white shadow-lg mx-auto"
    :style="pageStyle"
    @drop="handleDrop"
    @dragover="handleDragOver"
    @click="handleCanvasClick"
  >
    <!-- 未保存内容提示 -->
    <div v-if="hasUnsavedChanges" class="unsaved-changes-indicator">
      <span class="dot"></span>
      <span class="text">内容已更改，尚未保存</span>
    </div>
    
    <template v-if="!template?.components?.length">
      <div class="empty-state text-center py-8 text-gray-400">
        <p class="text-lg">拖放组件到此处开始编辑</p>
        <p class="text-sm">支持标题、文本、表格等多种组件</p>
      </div>
    </template>
    
    <template v-else>
      <div 
        v-for="(component, index) in template.components"
        :key="component.id"
        :class="getComponentWrapperClass(component, index)"
        :data-id="component.id"
        draggable="true"
        @dragstart="handleDragStart($event, component.id, index)"
        @dragend="handleDragEnd"
        @dragenter="handleDragEnter($event, index)"
        @click.stop="selectComponent(component.id)"
      >
        <!-- 标题组件 -->
        <div
          v-if="component.type === 'title'"
          :class="['component', 'title-component', { selected: String(component.id) === String(selectedComponentId) }]"
          :style="getComponentStyle(component)"
          @click.stop="selectComponent(component.id)"
        >
          {{ component.content || '标题' }}
        </div>
        
        <!-- 文本组件 -->
        <div
          v-else-if="component.type === 'text'"
          :class="['component', 'text-component', { selected: String(component.id) === String(selectedComponentId) }]"
          :style="getComponentStyle(component)"
          @click.stop="selectComponent(component.id)"
        >
          {{ component.content || '点击编辑文本' }}
        </div>
        
        <!-- 普通表格组件 -->
        <div
          v-else-if="component.type === 'table'"
          :class="['component', 'table-component', { selected: String(component.id) === String(selectedComponentId) }]"
          @click.stop="selectComponent(component.id)"
        >
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
                  class="border p-2 font-bold"
                  :style="{ 
                    borderWidth: component.settings?.borderStyle?.left ? '1px' : '0',
                    fontSize: `${component.settings?.fontSize || 16}px`,
                    fontFamily: component.settings?.fontFamily || '宋体',
                    textAlign: component.settings?.cellStyle?.alignment || 'left',
                    width: component.columnWidths ? `${component.columnWidths[0] <= 10 ? 1 : (component.columnWidths[0] >= 90 ? 99 : component.columnWidths[0])}%` : '50%'
                  }"
                >
                  {{ field.label }}
                </td>
                <td
                  class="border p-2"
                  :style="{ 
                    borderWidth: component.settings?.borderStyle?.right ? '1px' : '0',
                    fontSize: `${component.settings?.fontSize || 16}px`,
                    fontFamily: component.settings?.fontFamily || '宋体',
                    textAlign: component.settings?.cellStyle?.alignment || 'left',
                    width: component.columnWidths ? `${component.columnWidths[0] <= 10 ? 99 : (component.columnWidths[0] >= 90 ? 1 : component.columnWidths[1])}%` : '50%'
                  }"
                >
                  <!-- 如果有内容组件，则渲染组件 -->
                  <template v-if="field.contentComponents && field.contentComponents.length > 0">
                    <div v-for="(contentComp, compIndex) in field.contentComponents" :key="contentComp.id" class="mb-2">
                      <!-- 文本组件 -->
                      <div v-if="contentComp.type === 'text'" class="content-component text-component">
                        <p 
                          :style="{
                            fontSize: `${contentComp.fontSize || 16}px`,
                            fontFamily: contentComp.fontFamily || '宋体',
                            textAlign: contentComp.alignment || 'left',
                            lineHeight: contentComp.lineHeight || 1.5
                          }"
                        >
                          {{ contentComp.content }}
                        </p>
                      </div>
                      
                      <!-- 大标题组件 -->
                      <div v-else-if="contentComp.type === 'title'" class="content-component title-component">
                        <h2
                          :style="{
                            fontSize: `${contentComp.fontSize || 24}px`,
                            fontFamily: contentComp.fontFamily || '黑体',
                            textAlign: contentComp.alignment || 'center',
                            fontWeight: contentComp.bold ? 'bold' : 'normal'
                          }"
                        >
                          {{ contentComp.content }}
                        </h2>
                      </div>
                      
                      <!-- 普通文本段落组件 -->
                      <div v-else-if="contentComp.type === 'paragraph'" class="content-component paragraph-component">
                        <p 
                          :style="{
                            fontSize: `${contentComp.fontSize || 16}px`,
                            fontFamily: contentComp.fontFamily || '宋体',
                            textAlign: contentComp.alignment || 'left',
                            lineHeight: contentComp.lineHeight || 1.5
                          }"
                        >
                          {{ contentComp.content }}
                        </p>
                      </div>
                      
                      <!-- 进度表组件 -->
                      <div v-else-if="contentComp.type === 'schedule'" class="content-component schedule-component">
                        <div 
                          v-if="contentComp.title" 
                          class="font-bold mb-2"
                          :style="{
                            fontSize: `${contentComp.settings?.fontSize || 16}px`,
                            fontFamily: contentComp.settings?.fontFamily || '宋体'
                          }"
                        >
                          {{ contentComp.title }}
                        </div>
                        <div :style="{
                          fontSize: `${contentComp.settings?.fontSize || 16}px`,
                          fontFamily: contentComp.settings?.fontFamily || '宋体',
                          lineHeight: contentComp.settings?.lineHeight || 1.5
                        }">
                          <WorkUnit :items="contentComp.items" :in-table="true" />
                        </div>
                      </div>
                      
                      <!-- 列表组件 -->
                      <div v-else-if="contentComp.type === 'list'" class="content-component list-component">
                        <ul v-if="contentComp.listType === 'bullet'" class="list-disc pl-5">
                          <li 
                            v-for="item in contentComp.items" 
                            :key="item.id"
                            :style="{
                              fontSize: `${contentComp.fontSize || 16}px`,
                              fontFamily: contentComp.fontFamily || '宋体',
                              textAlign: contentComp.alignment || 'left'
                            }"
                          >
                            {{ item.content }}
                          </li>
                        </ul>
                        <ol v-else class="list-decimal pl-5">
                          <li 
                            v-for="item in contentComp.items" 
                            :key="item.id"
                            :style="{
                              fontSize: `${contentComp.fontSize || 16}px`,
                              fontFamily: contentComp.fontFamily || '宋体',
                              textAlign: contentComp.alignment || 'left'
                            }"
                          >
                            {{ item.content }}
                          </li>
                        </ol>
                      </div>
                      
                      <!-- 参考文献组件 -->
                      <div v-else-if="contentComp.type === 'reference'" class="content-component reference-component">
                        <div 
                          v-for="(item, itemIndex) in contentComp.items" 
                          :key="item.id"
                          class="mb-2"
                          :style="{
                            fontSize: `${contentComp.fontSize || 16}px`,
                            fontFamily: contentComp.fontFamily || '宋体'
                          }"
                        >
                          <!-- APA格式 -->
                          <p v-if="contentComp.style === 'apa'" class="pl-4 text-indent-[-1rem]">
                            {{ item.author }} ({{ item.year }}). <span class="italic">{{ item.title }}</span>. {{ item.source }}.
                          </p>
                          <!-- MLA格式 -->
                          <p v-else-if="contentComp.style === 'mla'" class="pl-4 text-indent-[-1rem]">
                            {{ item.author }}. "{{ item.title }}." <span class="italic">{{ item.source }}</span>, {{ item.year }}.
                          </p>
                          <!-- Chicago格式 -->
                          <p v-else-if="contentComp.style === 'chicago'" class="pl-4 text-indent-[-1rem]">
                            {{ item.author }}. <span class="italic">{{ item.title }}</span>. {{ item.source }}, {{ item.year }}.
                          </p>
                          <!-- GB/T 7714格式 -->
                          <p v-else class="pl-4 text-indent-[-1rem]">
                            [{{ itemIndex + 1 }}] {{ item.author }}. {{ item.title }}[M]. {{ item.source }}, {{ item.year }}.
                          </p>
                        </div>
                      </div>
                      
                      <!-- 图片组件 -->
                      <div v-else-if="contentComp.type === 'image'" class="content-component image-component">
                        <div :style="{ textAlign: contentComp.alignment || 'center' }">
                          <img 
                            v-if="contentComp.url" 
                            :src="contentComp.url" 
                            :alt="contentComp.alt || ''" 
                            :style="{ maxWidth: `${contentComp.width || 80}%`, margin: '0 auto' }"
                          />
                          <div v-else class="bg-gray-200 text-gray-500 flex items-center justify-center" style="height: 100px;">
                            图片占位符
                          </div>
                          <p 
                            v-if="contentComp.caption" 
                            class="text-center text-sm text-gray-600 mt-1"
                            :style="{
                              fontSize: `${(contentComp.fontSize || 16) - 2}px`,
                              fontFamily: contentComp.fontFamily || '宋体'
                            }"
                          >
                            {{ contentComp.caption }}
                          </p>
                        </div>
                      </div>
                    </div>
                  </template>
                  <!-- 如果没有内容组件，则显示默认值 -->
                  <template v-else>
                    {{ field.value || '点击编辑' }}
                  </template>
                </td>
              </tr>
            </template>
            
            <!-- 数据行 - 标题内容布局 -->
            <template v-else>
              <tr v-for="field in component.fields" :key="field.id">
                <td
                  class="border p-2"
                  :style="{ 
                    borderWidth: component.settings?.borderStyle?.left ? '1px' : '0',
                    fontSize: `${component.settings?.fontSize || 16}px`,
                    fontFamily: component.settings?.fontFamily || '宋体'
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

        <!-- 封面表格组件 -->
        <div
          v-else-if="component.type === 'cover-table'"
          :class="['component', 'cover-table-component', { selected: component.id === selectedComponentId }]"
          @click.stop="selectComponent(component.id)"
        >
          <table class="w-full border-collapse">
            <tr v-for="field in component.fields" :key="field.id">
              <td
                class="p-2"
                :style="{ 
                  width: `${field.width <= 10 ? 1 : (field.width >= 90 ? 99 : field.width)}%`,
                  fontSize: `${component.settings?.fontSize || 16}px`,
                  fontFamily: component.settings?.fontFamily || '宋体',
                  fontWeight: component.settings?.labelStyle?.bold ? 'bold' : 'normal',
                  textAlign: component.settings?.labelStyle?.alignment || 'left'
                }"
              >
                {{ field.label }}
              </td>
              <td
                class="p-2 border-b"
                :style="{ 
                  width: `${field.width <= 10 ? 99 : (field.width >= 90 ? 1 : 100 - field.width)}%`,
                  fontSize: `${component.settings?.fontSize || 16}px`,
                  fontFamily: component.settings?.fontFamily || '宋体',
                  textAlign: component.settings?.valueStyle?.alignment || 'center'
                }"
              >
                {{ field.value || '点击编辑' }}
              </td>
            </tr>
          </table>
        </div>

        <!-- 进度表组件 -->
        <div
          v-else-if="component.type === 'schedule'"
          :class="['component', 'schedule-component', { selected: component.id === selectedComponentId }]"
          @click.stop="selectComponent(component.id)"
        >
          <WorkUnit :items="component.items" :in-table="true" />
        </div>

        <!-- 论文目录组件 -->
        <div
          v-else-if="component.type === 'toc'"
          :class="['component', 'toc-component', { selected: component.id === selectedComponentId }]"
          @click.stop="selectComponent(component.id)"
        >
          <div 
            v-if="component.title" 
            class="font-bold mb-4 text-center" 
            :style="{ 
              fontFamily: component.titleFontFamily || '黑体',
              fontSize: `${component.titleFontSize || 18}px`
            }"
          >
            {{ component.title || '目录' }}
          </div>
          <div class="toc-items">
            <div 
              v-for="(item, index) in component.items" 
              :key="index"
              class="toc-item"
              :style="getComponentStyle(component)"
            >
              <div class="toc-line">
                <span class="toc-number">{{ item.number }}</span>
                <span class="toc-title">{{ item.title }}</span>
                <span class="toc-dots"></span>
                <span class="toc-page">{{ item.page }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 图片组件 -->
        <div
          v-else-if="component.type === 'image'"
          :class="['component', 'image-component', { selected: component.id === selectedComponentId }]"
          @click.stop="selectComponent(component.id)"
        >
          <div class="image-container" :style="{ textAlign: component.imageAlignment || 'center' }">
            <img 
              v-if="component.imageUrl" 
              :src="component.imageUrl" 
              :alt="component.caption || '图片'" 
              :style="{ 
                maxWidth: `${component.maxWidth || 100}%`,
                maxHeight: `${component.maxHeight || 'auto'}`
              }"
              class="mx-auto"
            />
            <div v-else class="image-placeholder">
              <p>请添加图片URL</p>
            </div>
          </div>
          <div 
            v-if="component.showCaption" 
            class="image-caption mt-2"
            :style="getComponentStyle(component)"
          >
            {{ component.caption || '图片说明（例如：图1-1 系统概览）' }}
          </div>
        </div>
        
        <!-- 参考文献组件 -->
        <div
          v-else-if="component.type === 'reference'"
          :class="['component', 'reference-component', { selected: component.id === selectedComponentId }]"
          @click.stop="selectComponent(component.id)"
        >
          <div 
            v-if="component.title" 
            class="font-bold mb-4" 
            :style="{ 
              fontFamily: component.titleFontFamily || '黑体',
              fontSize: `${component.titleFontSize || 18}px`,
              textAlign: component.alignment || 'left'
            }"
          >
            {{ component.title || '参考文献' }}
          </div>
          <div class="reference-items">
            <div 
              v-if="component.items && Array.isArray(component.items)"
              v-for="(item, index) in component.items" 
              :key="item.id"
              class="reference-item mb-2"
              :style="getComponentStyle(component)"
            >
              [{{ index + 1 }}] {{ item.content }}
            </div>
            <div v-else class="reference-item mb-2 text-gray-500 italic" :style="getComponentStyle(component)">
              [暂无参考文献]
            </div>
          </div>
        </div>
        
        <!-- 列表组件 -->
        <div
          v-else-if="component.type === 'list'"
          :class="['component', 'list-component', { selected: component.id === selectedComponentId }]"
          @click.stop="selectComponent(component.id)"
        >
          <div class="list-items">
            <div 
              v-if="component.items && Array.isArray(component.items)"
              v-for="(item, index) in component.items" 
              :key="item.id"
              class="list-item mb-2"
              :style="getComponentStyle(component)"
            >
              <span class="list-marker">{{ component.listType === 'bullet' ? '•' : (index + 1) + '.' }}</span>
              <span class="list-content">{{ item.content }}</span>
            </div>
            <div v-else class="list-item mb-2 text-gray-500 italic" :style="getComponentStyle(component)">
              [空列表]
            </div>
          </div>
        </div>
        
        <!-- 拖拽指示器 -->
        <div class="drag-handle" title="拖动调整位置">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="3" y1="12" x2="21" y2="12"></line>
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <line x1="3" y1="18" x2="21" y2="18"></line>
          </svg>
        </div>
        
        <!-- 操作按钮 -->
        <div class="component-actions">
          <button 
            class="delete-btn" 
            title="删除组件"
            @click.stop="handleComponentDelete(component.id)"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
      </div>
      
      <!-- 最后一个位置的拖拽指示器 -->
      <div 
        v-if="isDraggingComponent && dragOverIndex.value === template.components.length"
        class="drop-indicator"
      ></div>
    </template>
  </div>
</template>

<style scoped>
.template-canvas {
  background-color: white;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  margin: 0 auto;
  min-height: 100%;
  height: auto;
  overflow: visible;
  position: relative;
}

.component-wrapper {
  position: relative;
  margin: 0.5rem 0;
  transition: transform 0.2s, opacity 0.2s;
}

.component {
  position: relative;
  padding: 0.5rem;
  border: 2px solid transparent;
  border-radius: 4px;
  transition: all 0.2s;
  cursor: pointer;
}

.component:hover {
  border-color: #e2e8f0;
}

.component.selected {
  border-color: #3b82f6;
}

/* 拖拽样式 */
.component-wrapper.dragging {
  opacity: 0.5;
  transform: scale(0.98);
}

.component-wrapper.drag-over-top {
  border-top: 2px dashed #3b82f6;
  margin-top: 1rem;
  padding-top: 0.5rem;
}

.component-wrapper.drag-over-bottom {
  border-bottom: 2px dashed #3b82f6;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
}

.drop-indicator {
  height: 2px;
  background-color: #3b82f6;
  margin: 0.5rem 0;
}

.drag-handle {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  cursor: move;
  opacity: 0;
  transition: opacity 0.2s;
  background-color: rgba(255, 255, 255, 0.8);
  border-radius: 4px;
  padding: 2px;
  z-index: 10;
  color: #718096;
}

.component-wrapper:hover .drag-handle {
  opacity: 1;
}

.empty-state {
  border: 2px dashed #e2e8f0;
  border-radius: 8px;
  margin: 1rem;
}

.cover-table-component td {
  border: none;
}

.cover-table-component td.border-b {
  border-bottom: 1px solid #000;
}

/* 目录组件样式 */
.toc-component {
  width: 100%;
}

.toc-items {
  width: 100%;
}

.toc-item {
  margin-bottom: 0.5rem;
}

.toc-line {
  display: flex;
  width: 100%;
  position: relative;
}

.toc-number {
  margin-right: 0.5rem;
  white-space: nowrap;
}

.toc-title {
  white-space: nowrap;
  margin-right: 0.5rem;
}

.toc-dots {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.toc-dots::after {
  content: ".................................................................................................";
  position: absolute;
  width: 100%;
  overflow: hidden;
  white-space: nowrap;
  letter-spacing: 2px;
}

.toc-page {
  margin-left: 0.5rem;
  white-space: nowrap;
}

/* 图片组件样式 */
.image-component {
  width: 100%;
}

.image-container {
  width: 100%;
}

.image-placeholder {
  width: 100%;
  height: 150px;
  border: 2px dashed #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #a0aec0;
  font-size: 14px;
  border-radius: 4px;
}

.image-caption {
  font-size: 14px;
}

/* 参考文献组件样式 */
.reference-component {
  width: 100%;
}

.reference-items {
  width: 100%;
}

.reference-item {
  text-indent: -1.5em;
  padding-left: 1.5em;
}

/* 列表组件样式 */
.list-component {
  width: 100%;
}

.list-items {
  width: 100%;
}

.list-item {
  padding-left: 1.5em;
  position: relative;
}

.list-marker {
  position: absolute;
  left: 0;
  font-weight: bold;
}

.list-content {
  margin-left: 0.5em;
}

.unsaved-changes-indicator {
  position: absolute;
  top: 5px;
  right: 10px;
  padding: 4px 8px;
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 4px;
  font-size: 12px;
  color: #f56c6c;
  display: flex;
  align-items: center;
  z-index: 100;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  border: 1px solid #f0f0f0;
}

.unsaved-changes-indicator .dot {
  width: 8px;
  height: 8px;
  background-color: #f56c6c;
  border-radius: 50%;
  margin-right: 6px;
  animation: pulse 1.5s infinite;
}

.unsaved-changes-indicator .text {
  font-weight: 500;
}

@keyframes pulse {
  0% {
    opacity: 0.5;
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0.5;
  }
}

.component-actions {
  position: absolute;
  top: 0.5rem;
  right: 2rem;
  display: flex;
  gap: 5px;
  opacity: 0;
  transition: opacity 0.2s;
  z-index: 10;
}

.component-wrapper:hover .component-actions {
  opacity: 1;
}

.delete-btn {
  background-color: rgba(255, 255, 255, 0.8);
  color: #f56c6c;
  border: 1px solid #f0f0f0;
  border-radius: 4px;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.delete-btn:hover {
  background-color: #f56c6c;
  color: white;
}
</style> 