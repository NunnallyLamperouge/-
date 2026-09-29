<script setup>
import { computed, ref } from 'vue'
import { componentPanelMap } from './PropertyPanels'

const props = defineProps({
  component: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update', 'delete', 'template-changed', 'add-annotation'])

// 标注文本
const annotation = ref('')

// 删除组件
const deleteComponent = (componentId) => {
  if (confirm('确定要删除此组件吗？')) {
    emit('delete', componentId)
    emit('template-changed')
  }
}

// 处理组件更新并通知模板变更
const handleUpdate = (id, property, value) => {
  emit('update', id, property, value)
  emit('template-changed')
}

// 添加标注
const addAnnotation = (id, annotationText) => {
  if (annotationText.trim()) {
    emit('add-annotation', id, annotation.value)
    // 清空输入框
    annotation.value = ''
  }
}

// 删除标注
const deleteAnnotation = (componentId, annotationId) => {
  if (confirm('确定要删除这条标注吗？')) {
    // 找到该组件
    const component = props.component
    if (component && component.annotations) {
      // 从annotations数组中移除该标注
      const index = component.annotations.findIndex(a => a.id === annotationId)
      if (index !== -1) {
        component.annotations.splice(index, 1)
        // 通知模板变更
        emit('template-changed')
      }
    }
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
      <h3 class="text-lg font-bold">标注面板</h3>
      <button 
        @click="deleteComponent(component.id)" 
        class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
      >
        删除组件
      </button>
    </div>

    <!-- 添加标注区域 -->
    <div class="mb-4">
      <h4 class="text-md font-medium text-gray-700 mb-2">添加标注</h4>
      <a-textarea
        v-model:value="annotation"
        placeholder="输入对该组件的标注内容..."
        :rows="4"
        class="w-full mb-2"
      />
      <a-button
        type="primary"
        @click="addAnnotation(component.id, annotation)"
        class="w-full"
      >
        保存标注
      </a-button>
    </div>

    <!-- 现有标注列表 -->
    <div v-if="component.annotations && component.annotations.length > 0">
      <h4 class="text-md font-medium text-gray-700 mb-2">现有标注</h4>
      <a-list
        :data-source="component.annotations"
        class="annotation-list"
      >
        <template #renderItem="{ item }">
          <a-list-item>
            <div class="annotation-item p-2 bg-blue-50 rounded mb-2 w-full">
              <div class="flex justify-between">
                <p class="text-sm flex-grow">{{ item.content }}</p>
                <a-button 
                  type="text" 
                  danger 
                  size="small" 
                  @click="deleteAnnotation(component.id, item.id)"
                >
                  <template #icon><delete-outlined /></template>
                </a-button>
              </div>
              <div class="text-xs text-gray-500 mt-1">{{ item.date }}</div>
            </div>
          </a-list-item>
        </template>
      </a-list>
    </div>
    
    <!-- 组件信息展示 -->
    <div class="mt-4">
      <h4 class="text-md font-medium text-gray-700 mb-2">组件信息</h4>
      <div class="bg-gray-50 p-3 rounded">
        <p><span class="font-medium">类型:</span> {{ component.type }}</p>
        <p v-if="component.content"><span class="font-medium">内容:</span> {{ component.content }}</p>
      </div>
    </div>

    <!-- 未选择组件的提示 -->
    <div v-if="!component" class="p-4 bg-gray-100 text-gray-600 rounded text-center">
      请选择一个组件进行标注
    </div>
  </div>
</template>

<style scoped>
.property-panel {
  height: 100%;
  overflow-y: auto;
}

.annotation-item {
  border-left: 3px solid #1890ff;
}
</style> 