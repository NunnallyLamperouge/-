<script setup>
import { ref, computed } from 'vue'

// 定义可用的组件类型
const components = ref([
  {
    type: 'title',
    name: '大标题',
    icon: '📝',
    description: '用于添加文档标题',
    category: 'basic'
  },
  {
    type: 'text',
    name: '普通文本段落',
    icon: '📄',
    description: '用于添加文本内容',
    category: 'basic'
  },
  {
    type: 'table',
    name: '数据表格',
    icon: '📊',
    description: '用于添加表格数据',
    category: 'data'
  },
  {
    type: 'cover-table',
    name: '封面表格',
    icon: '📋',
    description: '用于添加封面样式的表格',
    category: 'data'
  },
  {
    type: 'image',
    name: '图片',
    icon: '🖼️',
    description: '添加图片及说明文字',
    category: 'media'
  },
  {
    type: 'toc',
    name: '论文目录',
    icon: '📑',
    description: '用于添加论文目录',
    category: 'structure'
  },
  {
    type: 'list',
    name: '列表',
    icon: '📋',
    description: '有序或无序列表',
    category: 'basic'
  },
  {
    type: 'reference',
    name: '参考文献',
    icon: '📚',
    description: '参考文献列表',
    category: 'structure'
  },
  {
    type: 'schedule',
    name: '进度表',
    icon: '📅',
    description: '工作进度安排表',
    category: 'data'
  }
])

// 组件分类
const categories = [
  { id: 'basic', name: '基础组件' },
  { id: 'structure', name: '结构组件' },
  { id: 'data', name: '数据组件' },
  { id: 'media', name: '媒体组件' }
]

// 当前选中的分类
const selectedCategory = ref('all')

// 根据分类筛选组件
const filteredComponents = computed(() => {
  if (selectedCategory.value === 'all') {
    return components.value
  }
  return components.value.filter(component => component.category === selectedCategory.value)
})

// 拖动处理
const handleDragStart = (e, component) => {
  e.dataTransfer.setData('componentType', component.type)
  emit('component-drag-start')
}

const handleDragEnd = () => {
  emit('component-drag-end')
}

// 定义事件
const emit = defineEmits(['component-drag-start', 'component-drag-end'])
</script>

<template>
  <div class="component-palette">
    <div class="category-tabs mb-4">
      <div class="flex space-x-1 overflow-x-auto pb-1">
        <button 
          class="category-tab px-3 py-1 rounded-md text-sm whitespace-nowrap"
          :class="selectedCategory === 'all' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
          @click="selectedCategory = 'all'"
        >
          全部
        </button>
        <button 
          v-for="category in categories" 
          :key="category.id"
          class="category-tab px-3 py-1 rounded-md text-sm whitespace-nowrap"
          :class="selectedCategory === category.id ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
          @click="selectedCategory = category.id"
        >
          {{ category.name }}
        </button>
      </div>
    </div>
    
    <div class="components-list space-y-2">
      <div
        v-for="component in filteredComponents"
        :key="component.type"
        class="component-item p-3 bg-white rounded-lg shadow hover:shadow-md cursor-move transition-all"
        :class="{'transform hover:-translate-y-1': true}"
        draggable="true"
        @dragstart="handleDragStart($event, component)"
        @dragend="handleDragEnd"
      >
        <div class="flex items-center gap-3">
          <div class="component-icon flex-shrink-0 w-10 h-10 flex items-center justify-center bg-blue-50 text-blue-500 rounded-lg">
            <span class="text-2xl">{{ component.icon }}</span>
          </div>
          <div class="flex-1 min-w-0">
            <h3 class="font-medium text-gray-800 truncate">{{ component.name }}</h3>
            <p class="text-xs text-gray-500 truncate">{{ component.description }}</p>
          </div>
        </div>
      </div>
    </div>
    
    <div class="mt-4 text-xs text-gray-500 text-center">
      <p>拖拽组件到画布中进行编辑</p>
    </div>
  </div>
</template>

<style scoped>
.component-item {
  touch-action: none;
  border: 1px solid #f0f0f0;
}

.component-item:hover {
  background-color: #f8fafc;
  border-color: #e2e8f0;
}

.component-icon {
  transition: all 0.2s;
}

.component-item:hover .component-icon {
  background-color: #dbeafe;
}

.category-tabs {
  position: sticky;
  top: 0;
  background-color: white;
  z-index: 10;
  padding-top: 0.5rem;
}

.components-list {
  max-height: calc(100vh - 250px);
  overflow-y: auto;
  padding-right: 0.25rem;
  margin-right: -0.25rem;
  scrollbar-width: thin;
}

.components-list::-webkit-scrollbar {
  width: 4px;
}

.components-list::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 10px;
}

.components-list::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 10px;
}

.components-list::-webkit-scrollbar-thumb:hover {
  background: #a1a1a1;
}
</style> 