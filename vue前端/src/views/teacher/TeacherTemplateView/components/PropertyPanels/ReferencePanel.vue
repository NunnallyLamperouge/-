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

// 添加参考文献条目
const addReferenceItem = () => {
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
      content: '新参考文献'
    })
    updateProperty('items', [...component.items])
  }
}

// 删除参考文献条目
const removeReferenceItem = (itemId) => {
  const component = props.component
  if (component && component.items && Array.isArray(component.items)) {
    if (component.items.length > 1) {
      const index = component.items.findIndex(i => i.id === itemId)
      if (index !== -1) {
        component.items.splice(index, 1)
        updateProperty('items', [...component.items])
      }
    } else {
      alert('至少保留一条参考文献')
    }
  }
}
</script>

<template>
  <div>
    <div class="mb-4">
      <label class="block mb-2">标题</label>
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
    
    <!-- 内容字体和字号 -->
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
    
    <!-- 参考文献条目 -->
    <div class="mb-4">
      <label class="block mb-2">参考文献条目</label>
      <div v-for="(item, index) in component.items" :key="item.id" class="mb-2">
        <div class="flex items-center">
          <span class="mr-2 text-gray-500">[{{ index + 1 }}]</span>
          <input
            v-model="item.content"
            type="text"
            class="flex-1 px-3 py-2 border rounded"
            @input="updateItemContent(item.id, $event.target.value)"
          />
          <button
            class="ml-2 p-1 text-red-500 hover:bg-red-100 rounded"
            @click="removeReferenceItem(item.id)"
          >
            <span class="text-lg">×</span>
          </button>
        </div>
      </div>
      <button
        class="mt-2 px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
        @click="addReferenceItem"
      >
        添加条目
      </button>
    </div>
  </div>
</template> 