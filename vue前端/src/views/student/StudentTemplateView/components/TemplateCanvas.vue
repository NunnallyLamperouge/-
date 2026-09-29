<script setup>
import { ref, computed, watch } from 'vue'
import WorkUnit from '../../../../components/WorkUnit.vue'

const props = defineProps({
  template: {
    type: Object,
    required: true
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

const emit = defineEmits(['update:selectedComponentId', 'component-update', 'component-delete', 'template-changed'])

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
const getComponentWrapperClass = (component) => {
  const classes = ['component-wrapper']
  
  // 如果组件有标注，添加一个特殊的类
  if (component.annotations && component.annotations.length > 0) {
    classes.push('has-annotations')
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

// 获取组件的标注数量
const getAnnotationCount = (component) => {
  return component.annotations?.length || 0
}
</script>

<template>
  <div 
    class="template-canvas bg-white shadow-lg mx-auto"
    :style="pageStyle"
    @click="handleCanvasClick"
  >
    <!-- 未保存内容提示 -->
    <div v-if="hasUnsavedChanges" class="unsaved-changes-indicator">
      <span class="dot"></span>
      <span class="text">内容已更改，尚未保存</span>
    </div>
    
    <template v-if="!template?.components?.length">
      <div class="empty-state text-center py-8 text-gray-400">
        <p class="text-lg">该模板暂无内容</p>
        <p class="text-sm">请在学生端创建模板内容</p>
      </div>
    </template>
    
    <template v-else>
      <div 
        v-for="(component, index) in template.components"
        :key="component.id"
        :class="getComponentWrapperClass(component)"
        :data-id="component.id"
        @click.stop="selectComponent(component.id)"
      >
        <!-- 标注标记 -->
        <div v-if="getAnnotationCount(component) > 0" class="annotation-badge">
          {{ getAnnotationCount(component) }}
        </div>
        
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
                    {{ field.value }}
                  </template>
                </td>
              </tr>
            </template>
          </table>
        </div>

        <!-- 封面表格组件 -->
        <div
          v-else-if="component.type === 'cover-table'"
          :class="['component', 'cover-table-component', { selected: String(component.id) === String(selectedComponentId) }]"
          @click.stop="selectComponent(component.id)"
        >
          <table class="w-full border-collapse">
            <tr v-for="(field, index) in component.fields" :key="field.id">
              <td
                class="border p-3"
                :style="{ 
                  width: `${field.width || 30}%`,
                  fontSize: `${component.settings?.fontSize || 16}px`,
                  fontFamily: component.settings?.fontFamily || '宋体',
                  textAlign: component.settings?.labelStyle?.alignment || 'center',
                  fontWeight: component.settings?.labelStyle?.bold ? 'bold' : 'normal'
                }"
              >
                {{ field.label }}
              </td>
              <td
                class="border p-3"
                :style="{ 
                  fontSize: `${component.settings?.fontSize || 16}px`,
                  fontFamily: component.settings?.fontFamily || '宋体',
                  textAlign: component.settings?.valueStyle?.alignment || 'center'
                }"
              >
                {{ field.value }}
              </td>
            </tr>
          </table>
        </div>

        <!-- 论文目录组件 -->
        <div
          v-else-if="component.type === 'toc'"
          :class="['component', 'toc-component', { selected: String(component.id) === String(selectedComponentId) }]"
          @click.stop="selectComponent(component.id)"
        >
          <div 
            v-if="component.title" 
            class="font-bold mb-4 pb-2 border-b" 
            :style="{ 
              textAlign: 'center',
              fontFamily: component.titleFontFamily || '黑体',
              fontSize: `${component.titleFontSize || 18}px`,
              borderBottomWidth: '1px'
            }"
          >
            {{ component.title }}
          </div>
          
          <div class="toc-items">
            <div 
              v-for="(item, index) in component.items" 
              :key="item.id"
              class="toc-item mb-2 flex items-baseline justify-between"
              :style="{ 
                paddingLeft: item.number.split('.').length > 1 ? `${(item.number.split('.').length - 1) * 1.5}em` : '0',
                fontSize: `${component.fontSize || 16}px`,
                fontFamily: component.fontFamily || '宋体',
                lineHeight: component.lineHeight || 1.8
              }"
            >
              <div class="flex-1">
                <span class="toc-number mr-2">{{ item.number }}</span>
                <span class="toc-title">{{ item.title }}</span>
              </div>
              <div class="toc-dots flex-1 border-b border-dotted border-gray-400 mx-2"></div>
              <div class="toc-page text-right">{{ item.page }}</div>
            </div>
          </div>
        </div>
        
        <!-- 参考文献组件 -->
        <div
          v-else-if="component.type === 'reference'"
          :class="['component', 'reference-component', { selected: String(component.id) === String(selectedComponentId) }]"
          @click.stop="selectComponent(component.id)"
        >
          <div 
            v-if="component.title" 
            class="font-bold mb-4" 
            :style="{ 
              textAlign: 'center',
              fontFamily: component.titleFontFamily || '黑体',
              fontSize: `${component.titleFontSize || 18}px`
            }"
          >
            {{ component.title }}
          </div>
          
          <div class="reference-items">
            <div 
              v-for="(item, index) in component.items" 
              :key="item.id"
              class="reference-item mb-2"
              :style="{ 
                fontSize: `${component.fontSize || 16}px`,
                fontFamily: component.fontFamily || '宋体',
                lineHeight: component.lineHeight || 1.5,
                textAlign: component.alignment || 'left'
              }"
            >
              [{{ index + 1 }}] {{ item.content }}
            </div>
          </div>
        </div>
        
        <!-- 列表组件 -->
        <div
          v-else-if="component.type === 'list'"
          :class="['component', 'list-component', { selected: String(component.id) === String(selectedComponentId) }]"
          @click.stop="selectComponent(component.id)"
        >
          <ul v-if="component.listType === 'bullet'" class="list-disc pl-5">
            <li 
              v-for="item in component.items" 
              :key="item.id"
              :style="{ 
                fontSize: `${component.fontSize || 16}px`,
                fontFamily: component.fontFamily || '宋体',
                lineHeight: component.lineHeight || 1.5,
                textAlign: component.alignment || 'left'
              }"
              class="mb-2"
            >
              {{ item.content }}
            </li>
          </ul>
          
          <ol v-else-if="component.listType === 'number'" class="list-decimal pl-5">
            <li 
              v-for="item in component.items" 
              :key="item.id"
              :style="{ 
                fontSize: `${component.fontSize || 16}px`,
                fontFamily: component.fontFamily || '宋体',
                lineHeight: component.lineHeight || 1.5,
                textAlign: component.alignment || 'left'
              }"
              class="mb-2"
            >
              {{ item.content }}
            </li>
          </ol>
        </div>
        
        <!-- 进度表组件 -->
        <div
          v-else-if="component.type === 'schedule'"
          :class="['component', 'schedule-component', { selected: String(component.id) === String(selectedComponentId) }]"
          @click.stop="selectComponent(component.id)"
        >
          <div v-if="component.title" class="font-bold mb-2">{{ component.title }}</div>
          <div
            :style="{ 
              fontSize: `${component.settings?.fontSize || 16}px`,
              fontFamily: component.settings?.fontFamily || '宋体',
              lineHeight: component.settings?.lineHeight || 1.5
            }"
          >
            <WorkUnit :items="component.items" />
          </div>
        </div>
        
        <!-- 图片组件 -->
        <div
          v-else-if="component.type === 'image'"
          :class="['component', 'image-component', { selected: String(component.id) === String(selectedComponentId) }]"
          @click.stop="selectComponent(component.id)"
          >
          <div 
            class="image-container mb-2"
            :style="{ 
              textAlign: component.imageAlignment || 'center',
              maxWidth: `${component.maxWidth || 100}%`,
              margin: '0 auto'
            }"
          >
            <img 
              v-if="component.imageUrl && !component.imagePreviewError"
              :src="component.imageUrl" 
              @error="component.imagePreviewError = true"
              :style="{ 
                maxWidth: '100%',
                maxHeight: component.maxHeight === 'auto' ? 'auto' : `${component.maxHeight}px`
              }" 
              class="mx-auto"
            />
            <div v-else class="image-placeholder bg-gray-100 p-10 text-center text-gray-400">
              {{ component.imagePreviewError ? '图片加载失败' : '暂无图片' }}
            </div>
          </div>
          <div 
            v-if="component.showCaption" 
            class="image-caption"
            :style="{ 
              textAlign: component.captionAlignment || 'center',
              fontSize: `${component.captionFontSize || 14}px`,
              fontFamily: component.captionFontFamily || '宋体',
              lineHeight: component.captionLineHeight || 1.5
            }"
          >
            {{ component.caption || '图片说明' }}
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.template-canvas {
  position: relative;
  border: 1px solid #e0e0e0;
  margin: 0 auto;
}

.component-wrapper {
  position: relative;
  margin-bottom: 16px;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.2s;
  cursor: pointer;
}

.component-wrapper:hover {
  background-color: rgba(224, 242, 254, 0.3);
}

.component {
  padding: 8px 4px;
}

.component-wrapper.has-annotations {
  border-left: 3px solid #1890ff;
  padding-left: 8px;
}

.component.selected {
  background-color: rgba(24, 144, 255, 0.1);
  outline: 2px solid rgba(24, 144, 255, 0.5);
}

.unsaved-changes-indicator {
  position: absolute;
  top: 10px;
  right: 10px;
  display: flex;
  align-items: center;
  background-color: rgba(255, 77, 79, 0.1);
  padding: 4px 8px;
  border-radius: 4px;
  z-index: 10;
}

.unsaved-changes-indicator .dot {
  width: 8px;
  height: 8px;
  background-color: #ff4d4f;
  border-radius: 50%;
  margin-right: 6px;
  animation: pulse 1.5s infinite;
}

.unsaved-changes-indicator .text {
  font-size: 12px;
  color: #ff4d4f;
}

.annotation-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 20px;
  height: 20px;
  background-color: #1890ff;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  z-index: 10;
}

@keyframes pulse {
  0% {
    transform: scale(0.8);
    opacity: 0.8;
  }
  50% {
    transform: scale(1.2);
    opacity: 1;
  }
  100% {
    transform: scale(0.8);
    opacity: 0.8;
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
}

.toc-dots {
  flex: 1;
  height: 0.5em;
  margin: 0 0.5em;
}
</style> 