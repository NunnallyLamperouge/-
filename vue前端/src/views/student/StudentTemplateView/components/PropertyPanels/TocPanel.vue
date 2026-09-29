<script setup>
import { ref } from 'vue'
import { fontFamilies } from './common/CommonData'

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

// 添加目录项
const addTocItem = () => {
  if (!props.component.items) {
    props.component.items = []
  }
  
  const lastItem = props.component.items[props.component.items.length - 1]
  let newNumber = '1'
  
  if (lastItem) {
    // 尝试解析最后一个项目的编号，并增加1
    const parts = lastItem.number.split('.')
    if (parts.length > 1) {
      // 如果是子章节，增加最后一个数字
      const lastPart = parseInt(parts[parts.length - 1])
      parts[parts.length - 1] = (lastPart + 1).toString()
      newNumber = parts.join('.')
    } else {
      // 如果是主章节，增加章节号
      newNumber = (parseInt(lastItem.number) + 1).toString()
    }
  }
  
  props.component.items.push({
    id: Date.now(),
    number: newNumber,
    title: '新章节',
    page: (parseInt(lastItem?.page || 0) + 1).toString()
  })
  
  updateProperty('items', [...props.component.items])
}

// 删除目录项
const removeTocItem = (index) => {
  if (props.component.items.length > 1) {
    props.component.items.splice(index, 1)
    updateProperty('items', [...props.component.items])
  } else {
    alert('至少保留一个目录项')
  }
}

// 添加子章节
const addSubchapter = (index) => {
  const parentItem = props.component.items[index]
  if (!parentItem) return
  
  // 创建子章节编号
  const newNumber = `${parentItem.number}.1`
  
  // 在父章节后插入子章节
  props.component.items.splice(index + 1, 0, {
    id: Date.now(),
    number: newNumber,
    title: '新子章节',
    page: (parseInt(parentItem.page) + 1).toString()
  })
  
  updateProperty('items', [...props.component.items])
}

// 移动目录项
let draggedItemIndex = -1

const startDrag = (index) => {
  draggedItemIndex = index
}

const onDragOver = (index, event) => {
  event.preventDefault()
}

const onDrop = (index) => {
  if (draggedItemIndex !== -1 && draggedItemIndex !== index) {
    // 移动项目
    const item = props.component.items.splice(draggedItemIndex, 1)[0]
    props.component.items.splice(index, 0, item)
    updateProperty('items', [...props.component.items])
  }
  draggedItemIndex = -1
}
</script>

<template>
  <div>
    <div class="mb-4">
      <label class="block mb-2">目录标题</label>
      <input
        v-model="component.title"
        type="text"
        class="w-full px-3 py-2 border rounded"
        @input="updateProperty('title', $event.target.value)"
      />
    </div>
    
    <!-- 标题字体和字号 -->
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
        <input
          v-model.number="component.titleFontSize"
          type="number"
          class="w-full px-3 py-2 border rounded"
          min="12"
          max="72"
          @input="updateProperty('titleFontSize', $event.target.value)"
        />
      </div>
    </div>
    
    <!-- 目录内容字体和字号 -->
    <div class="grid grid-cols-2 gap-4 mb-4">
      <div>
        <label class="block mb-2">内容字体</label>
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
        <label class="block mb-2">内容字号</label>
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
    
    <!-- 目录项设置 -->
    <div class="mb-4">
      <div class="flex justify-between items-center mb-2">
        <label class="block mb-2">目录项</label>
        <button
          @click="addTocItem"
          class="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          添加章节
        </button>
      </div>
      
      <div class="space-y-2">
        <div 
          v-for="(item, index) in component.items" 
          :key="item.id" 
          class="p-3 border rounded bg-white"
          draggable="true"
          @dragstart="startDrag(index)"
          @dragend="draggedItemIndex = -1"
          @dragover.prevent="onDragOver(index, $event)"
          @drop.prevent="onDrop(index)"
        >
          <div class="flex items-center mb-2">
            <div class="mr-2 cursor-move text-gray-400 hover:text-gray-600">⋮⋮</div>
            <div class="flex-1 grid grid-cols-3 gap-2">
              <input
                v-model="item.number"
                type="text"
                class="px-3 py-2 border rounded"
                placeholder="编号"
                @input="updateProperty('items', [...component.items])"
              />
              <input
                v-model="item.title"
                type="text"
                class="col-span-2 px-3 py-2 border rounded"
                placeholder="标题"
                @input="updateProperty('items', [...component.items])"
              />
            </div>
          </div>
          
          <div class="flex items-center pl-6">
            <input
              v-model="item.page"
              type="text"
              class="w-20 px-3 py-2 border rounded"
              placeholder="页码"
              @input="updateProperty('items', [...component.items])"
            />
            <div class="flex-1"></div>
            <div class="flex gap-2">
              <button
                @click="addSubchapter(index)"
                class="px-2 py-1 bg-blue-100 text-blue-700 rounded hover:bg-blue-200"
                title="添加子章节"
              >
                添加子章节
              </button>
              <button
                @click="removeTocItem(index)"
                class="px-2 py-1 bg-red-100 text-red-700 rounded hover:bg-red-200"
                title="删除此章节"
              >
                删除
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 拖拽提示 -->
      <div class="mt-2 text-xs text-gray-500">
        提示：可拖拽目录项调整顺序
      </div>
    </div>
    
    <!-- 目录样式设置 -->
    <div class="mb-4">
      <label class="block mb-2">目录样式</label>
      <div class="grid grid-cols-2 gap-4">
        <button
          :class="[
            'p-3 border rounded flex flex-col items-center',
            component.style === 'dots' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('style', 'dots')"
        >
          <div class="text-2xl mb-1">⋯⋯</div>
          <div class="text-sm">点线型</div>
        </button>
        <button
          :class="[
            'p-3 border rounded flex flex-col items-center',
            component.style === 'plain' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('style', 'plain')"
        >
          <div class="text-2xl mb-1">—</div>
          <div class="text-sm">简洁型</div>
        </button>
      </div>
    </div>
  </div>
</template> 