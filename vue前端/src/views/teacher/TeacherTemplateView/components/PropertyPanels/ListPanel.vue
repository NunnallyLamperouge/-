<script setup>
import { ref } from 'vue'

const props = defineProps({
  component: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update'])

// 更新组件属性的辅助函数
const updateProperty = (property, value) => {
  emit('update', props.component.id, property, value)
}

// 更新参考文献条目内容
const updateItemContent = (itemId, content) => {
  const component = props.component
  if (component && component.items && Array.isArray(component.items)) {
    const item = component.items.find(i => i.id === itemId)
    if (item) {
      item.content = content
      updateProperty('items', [...component.items])
    }
  }
}

// 添加列表项
const addListItem = () => {
  const component = props.component
  if (component) {
    // 确保items属性存在且是数组
    if (!component.items || !Array.isArray(component.items)) {
      component.items = []
    }
    
    const newId = component.items.length > 0 
      ? Math.max(0, ...component.items.map(i => i.id)) + 1 
      : 1
    
    component.items.push({
      id: newId,
      content: '新列表项'
    })
    updateProperty('items', [...component.items])
  }
}

// 删除列表项
const removeListItem = (itemId) => {
  const component = props.component
  if (component && component.items && Array.isArray(component.items)) {
    if (component.items.length > 1) {
      const index = component.items.findIndex(i => i.id === itemId)
      if (index !== -1) {
        component.items.splice(index, 1)
        updateProperty('items', [...component.items])
      }
    } else {
      alert('至少保留一个列表项')
    }
  }
}

// 批量添加列表项
const addMultipleListItems = (itemsText) => {
  const component = props.component
  if (component) {
    // 确保items属性存在且是数组
    if (!component.items || !Array.isArray(component.items)) {
      component.items = []
    }
    
    // 分割文本为多行
    const lines = itemsText.split('\n').filter(line => line.trim() !== '')
    
    // 获取当前最大ID
    const maxId = component.items.length > 0 
      ? Math.max(0, ...component.items.map(i => i.id))
      : 0
    
    // 添加新项目
    lines.forEach((line, index) => {
      component.items.push({
        id: maxId + index + 1,
        content: line.trim()
      })
    })
    
    updateProperty('items', [...component.items])
  }
}

// 移动列表项
const moveListItem = (fromIndex, toIndex) => {
  const component = props.component
  if (component && component.items && Array.isArray(component.items)) {
    // 确保索引有效
    if (fromIndex >= 0 && fromIndex < component.items.length &&
        toIndex >= 0 && toIndex < component.items.length) {
      // 移动项目
      const item = component.items.splice(fromIndex, 1)[0]
      component.items.splice(toIndex, 0, item)
      updateProperty('items', [...component.items])
    }
  }
}

// 应用列表预设模板
const applyListTemplate = (templateName) => {
  const component = props.component
  if (component) {
    let templateItems = []
    
    switch (templateName) {
      case 'weekdays':
        templateItems = [
          { id: 1, content: '星期一' },
          { id: 2, content: '星期二' },
          { id: 3, content: '星期三' },
          { id: 4, content: '星期四' },
          { id: 5, content: '星期五' },
          { id: 6, content: '星期六' },
          { id: 7, content: '星期日' }
        ]
        break
      case 'months':
        templateItems = [
          { id: 1, content: '一月' },
          { id: 2, content: '二月' },
          { id: 3, content: '三月' },
          { id: 4, content: '四月' },
          { id: 5, content: '五月' },
          { id: 6, content: '六月' },
          { id: 7, content: '七月' },
          { id: 8, content: '八月' },
          { id: 9, content: '九月' },
          { id: 10, content: '十月' },
          { id: 11, content: '十一月' },
          { id: 12, content: '十二月' }
        ]
        break
      case 'steps':
        templateItems = [
          { id: 1, content: '第一步：准备材料' },
          { id: 2, content: '第二步：开始操作' },
          { id: 3, content: '第三步：完成工作' }
        ]
        break
      case 'clear':
        templateItems = [
          { id: 1, content: '新列表项' }
        ]
        break
    }
    
    component.items = templateItems
    updateProperty('items', [...component.items])
  }
}

// 处理列表项键盘事件
const handleListItemKeydown = (event, itemId, index) => {
  const component = props.component
  
  // 回车键：在当前项后添加新项
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    
    const newId = Math.max(0, ...component.items.map(i => i.id)) + 1
    const newItem = { id: newId, content: '' }
    
    component.items.splice(index + 1, 0, newItem)
    updateProperty('items', [...component.items])
    
    // 焦点移到新项
    setTimeout(() => {
      const inputs = document.querySelectorAll('.list-item-input')
      if (inputs[index + 1]) {
        inputs[index + 1].focus()
      }
    }, 10)
  }
  
  // Tab键：增加缩进（这里只是示例，实际实现可能需要更复杂的逻辑）
  else if (event.key === 'Tab' && !event.shiftKey) {
    event.preventDefault()
    const item = component.items.find(i => i.id === itemId)
    if (item) {
      if (!item.content.startsWith('    ')) {
        item.content = '    ' + item.content
        updateProperty('items', [...component.items])
      }
    }
  }
  
  // Shift+Tab：减少缩进
  else if (event.key === 'Tab' && event.shiftKey) {
    event.preventDefault()
    const item = component.items.find(i => i.id === itemId)
    if (item && item.content.startsWith('    ')) {
      item.content = item.content.substring(4)
      updateProperty('items', [...component.items])
    }
  }
  
  // Backspace键：如果内容为空，则删除当前项
  else if (event.key === 'Backspace') {
    const item = component.items.find(i => i.id === itemId)
    if (item && item.content === '' && component.items.length > 1) {
      event.preventDefault()
      removeListItem(itemId)
      
      // 焦点移到前一项
      setTimeout(() => {
        const inputs = document.querySelectorAll('.list-item-input')
        if (inputs[index - 1]) {
          inputs[index - 1].focus()
        }
      }, 10)
    }
  }
}

// 开始拖拽列表项
let draggedItemIndex = -1

const startDrag = (index) => {
  draggedItemIndex = index
}

const onDragOver = (index, event) => {
  event.preventDefault()
}

const onDrop = (index) => {
  if (draggedItemIndex !== -1 && draggedItemIndex !== index) {
    moveListItem(draggedItemIndex, index)
  }
  draggedItemIndex = -1
}

// 批量添加对话框
const showBatchAddDialog = ref(false)
const batchItemsText = ref('')

const openBatchAddDialog = () => {
  batchItemsText.value = ''
  showBatchAddDialog.value = true
}

const closeBatchAddDialog = () => {
  showBatchAddDialog.value = false
}

const confirmBatchAdd = () => {
  if (batchItemsText.value.trim()) {
    addMultipleListItems(batchItemsText.value)
    closeBatchAddDialog()
  }
}

// 字体列表
const fontFamilies = [
  { value: '宋体', label: '宋体' },
  { value: '黑体', label: '黑体' },
  { value: '楷体', label: '楷体' },
  { value: '仿宋', label: '仿宋' }
]
</script>

<template>
  <div>
    <div class="mb-4">
      <label class="block mb-2">列表类型</label>
      <div class="flex gap-2">
        <button
          class="flex-1 px-3 py-2 border rounded flex items-center justify-center"
          :class="{ 'bg-blue-500 text-white': component.listType === 'bullet', 'hover:bg-gray-100': component.listType !== 'bullet' }"
          @click="updateProperty('listType', 'bullet')"
        >
          <span class="mr-2 text-xl">•</span> 无序列表
        </button>
        <button
          class="flex-1 px-3 py-2 border rounded flex items-center justify-center"
          :class="{ 'bg-blue-500 text-white': component.listType === 'number', 'hover:bg-gray-100': component.listType !== 'number' }"
          @click="updateProperty('listType', 'number')"
        >
          <span class="mr-2 font-bold">1.</span> 有序列表
        </button>
      </div>
    </div>
    
    <!-- 预设模板 -->
    <div class="mb-4">
      <label class="block mb-2">预设模板</label>
      <div class="grid grid-cols-2 gap-2">
        <button
          class="px-3 py-2 border rounded text-sm hover:bg-gray-100"
          @click="applyListTemplate('weekdays')"
        >
          星期
        </button>
        <button
          class="px-3 py-2 border rounded text-sm hover:bg-gray-100"
          @click="applyListTemplate('months')"
        >
          月份
        </button>
        <button
          class="px-3 py-2 border rounded text-sm hover:bg-gray-100"
          @click="applyListTemplate('steps')"
        >
          步骤
        </button>
        <button
          class="px-3 py-2 border rounded text-sm hover:bg-gray-100"
          @click="applyListTemplate('clear')"
        >
          清空
        </button>
      </div>
    </div>
    
    <!-- 内容字体和字号 -->
    <div class="grid grid-cols-2 gap-4 mb-4">
      <div>
        <label class="block mb-2">字体</label>
        <select
          v-model="component.fontFamily"
          class="w-full px-3 py-2 border rounded"
          @change="updateProperty('fontFamily', $event.target.value)"
        >
          <option v-for="font in fontFamilies" :key="font.value" :value="font.value">
            {{ font.label }}
          </option>
        </select>
      </div>
      <div>
        <label class="block mb-2">字号</label>
        <input
          v-model.number="component.fontSize"
          type="number"
          class="w-full px-3 py-2 border rounded"
          min="12"
          max="72"
          @input="updateProperty('fontSize', $event.target.value)"
        />
      </div>
    </div>
    
    <div class="mb-4">
      <label class="block mb-2">行距</label>
      <input
        v-model.number="component.lineHeight"
        type="number"
        class="w-full px-3 py-2 border rounded"
        min="1"
        max="3"
        step="0.1"
        @input="updateProperty('lineHeight', $event.target.value)"
      />
    </div>
    
    <div class="mb-4">
      <label class="block mb-2">对齐方式</label>
      <select
        v-model="component.alignment"
        class="w-full px-3 py-2 border rounded"
        @change="updateProperty('alignment', $event.target.value)"
      >
        <option value="left">左对齐</option>
        <option value="center">居中</option>
        <option value="right">右对齐</option>
      </select>
    </div>
    
    <!-- 列表项 -->
    <div class="mb-4">
      <div class="flex justify-between items-center mb-2">
        <label class="block font-medium">列表项</label>
        <div class="flex gap-2">
          <button
            class="px-2 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 text-sm flex items-center"
            @click="addListItem"
            title="添加单个列表项"
          >
            <span class="text-lg mr-1">+</span> 添加
          </button>
          <button
            class="px-2 py-1 bg-green-500 text-white rounded hover:bg-green-600 text-sm flex items-center"
            @click="openBatchAddDialog"
            title="批量添加多个列表项"
          >
            <span class="text-lg mr-1">+</span> 批量添加
          </button>
        </div>
      </div>
      
      <!-- 列表项编辑区域 -->
      <div class="border rounded p-2 bg-gray-50 max-h-60 overflow-y-auto">
        <p v-if="!component.items || !component.items.length" class="text-gray-500 italic text-center py-2">
          暂无列表项，请添加
        </p>
        
        <div 
          v-for="(item, index) in component.items" 
          :key="item.id" 
          class="mb-2 bg-white border rounded p-1 shadow-sm hover:shadow-md transition-shadow duration-200"
          draggable="true"
          @dragstart="startDrag(index)"
          @dragover="onDragOver(index, $event)"
          @drop="onDrop(index)"
        >
          <div class="flex items-center">
            <div class="mr-2 text-gray-500 w-6 flex-shrink-0 flex justify-center cursor-move">
              <span v-if="component.listType === 'bullet'">•</span>
              <span v-else>{{ index + 1 }}.</span>
            </div>
            <input
              v-model="item.content"
              type="text"
              class="flex-1 px-3 py-2 border rounded list-item-input"
              placeholder="输入列表项内容"
              @input="updateItemContent(item.id, $event.target.value)"
              @keydown="handleListItemKeydown($event, item.id, index)"
            />
            <button
              class="ml-2 p-1 text-red-500 hover:bg-red-100 rounded flex-shrink-0"
              @click="removeListItem(item.id)"
              title="删除此列表项"
            >
              <span class="text-lg">×</span>
            </button>
          </div>
        </div>
      </div>
      
      <div class="mt-2 text-xs text-gray-500">
        <p>提示：</p>
        <ul class="list-disc pl-5">
          <li>拖拽列表项可以调整顺序</li>
          <li>按回车键在当前项后添加新项</li>
          <li>按Tab键增加缩进，Shift+Tab减少缩进</li>
          <li>当内容为空时按Backspace键删除项</li>
        </ul>
      </div>
    </div>
    
    <!-- 批量添加对话框 -->
    <div v-if="showBatchAddDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h3 class="text-lg font-bold mb-4">批量添加列表项</h3>
        <p class="mb-2 text-sm text-gray-600">每行将作为一个列表项添加</p>
        <textarea
          v-model="batchItemsText"
          class="w-full px-3 py-2 border rounded h-40"
          placeholder="输入多行文本，每行将成为一个列表项"
        ></textarea>
        <div class="flex justify-end gap-2 mt-4">
          <button
            class="px-4 py-2 border rounded hover:bg-gray-100"
            @click="closeBatchAddDialog"
          >
            取消
          </button>
          <button
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            @click="confirmBatchAdd"
          >
            添加
          </button>
        </div>
      </div>
    </div>
  </div>
</template> 