<script setup>
import { computed } from 'vue'
import { componentPanelMap } from './PropertyPanels'

const props = defineProps({
  component: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update', 'delete'])

// 删除组件
const deleteComponent = (componentId) => {
  if (confirm('确定要删除此组件吗？')) {
    emit('delete', componentId)
  }
}

// 获取当前组件对应的面板组件
const currentPanel = computed(() => {
  if (!props.component) return null
  return componentPanelMap[props.component.type] || null
})
</script>

<template>
  <div class="property-panel">
    <div class="flex justify-between items-center mb-4">
      <h3 class="text-lg font-bold">属性设置</h3>
      <button 
        @click="deleteComponent(component.id)" 
        class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
      >
        删除组件
      </button>
    </div>

    <!-- 动态加载对应的属性面板组件 -->
    <component 
      v-if="currentPanel && component" 
      :is="currentPanel" 
      :component="component"
      @update="(id, property, value) => $emit('update', id, property, value)"
    />
    
    <!-- 未找到对应面板组件的提示 -->
    <div v-else-if="component" class="p-4 bg-yellow-100 text-yellow-800 rounded">
      暂不支持编辑 {{ component.type }} 类型的组件属性
    </div>
    
    <!-- 未选择组件的提示 -->
    <div v-else class="p-4 bg-gray-100 text-gray-600 rounded text-center">
      请选择一个组件进行编辑
    </div>
  </div>
</template>

<style scoped>
.property-panel {
  height: 100%;
  overflow-y: auto;
}
</style> 