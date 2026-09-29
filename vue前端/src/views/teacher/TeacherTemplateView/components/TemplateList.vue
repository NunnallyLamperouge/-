<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  templates: {
    type: Array,
    default: () => []
  },
  currentIndex: {
    type: Number,
    default: -1
  },
  isLoading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['select', 'delete'])

const handleSelect = (index) => {
  emit('select', index)
}

const handleDelete = (index, event) => {
  event.stopPropagation()
  if (confirm('确定要删除此模板吗？')) {
    emit('delete', index)
  }
}
</script>

<template>
  <div class="mt-4">
    <h3 class="text-sm font-medium text-gray-700 mb-2">已保存的模板</h3>
    
    <!-- 加载状态 -->
    <div v-if="isLoading" class="flex justify-center items-center py-4 border rounded">
      <div class="animate-spin mr-2">⟳</div>
      <span class="text-sm text-gray-500">加载中...</span>
    </div>
    
    <!-- 空状态 -->
    <div v-else-if="templates.length === 0" class="py-4 text-center border rounded">
      <p class="text-sm text-gray-500">暂无保存的模板</p>
    </div>
    
    <!-- 模板列表 -->
    <div v-else class="max-h-60 overflow-y-auto border rounded">
      <div 
        v-for="(template, index) in templates" 
        :key="template.id"
        class="px-3 py-2 flex justify-between items-center cursor-pointer hover:bg-gray-100"
        :class="{ 'bg-blue-50': index === currentIndex }"
        @click="handleSelect(index)"
      >
        <div class="flex items-center">
          <span class="text-sm">{{ template.name }}</span>
          <span class="ml-2 text-xs text-gray-500">{{ template.type }}</span>
        </div>
        <button 
          class="text-red-500 hover:text-red-700"
          @click="(e) => handleDelete(index, e)"
        >
          删除
        </button>
      </div>
    </div>
  </div>
</template> 