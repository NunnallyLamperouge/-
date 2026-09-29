<script setup>
import { ref, reactive, defineProps, defineEmits } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'create'])

const newTemplateForm = reactive({
  name: '',
  type: 'custom'
})

const handleCreate = () => {
  if (!newTemplateForm.name.trim()) {
    alert('请输入模板名称')
    return
  }
  
  emit('create', { ...newTemplateForm })
  newTemplateForm.name = ''
  newTemplateForm.type = 'custom'
}

const handleCancel = () => {
  emit('close')
  newTemplateForm.name = ''
  newTemplateForm.type = 'custom'
}
</script>

<template>
  <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white p-6 rounded-lg w-96">
      <h2 class="text-xl font-bold mb-4">新建模板</h2>
      <div class="mb-4">
        <label class="block mb-2">模板名称</label>
        <input 
          v-model="newTemplateForm.name"
          type="text"
          class="w-full px-3 py-2 border rounded"
        />
      </div>
      <div class="mb-4">
        <label class="block mb-2">模板类型</label>
        <select 
          v-model="newTemplateForm.type"
          class="w-full px-3 py-2 border rounded"
        >
          <option value="custom">自定义模板</option>
          <option value="thesis">论文模板</option>
          <option value="proposal">开题报告</option>
          <option value="task">任务书</option>
          <option value="review">中期检查</option>
        </select>
      </div>
      <div class="flex justify-end gap-2">
        <button 
          class="px-4 py-2 bg-gray-300 rounded hover:bg-gray-400"
          @click="handleCancel"
        >
          取消
        </button>
        <button 
          class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          @click="handleCreate"
        >
          创建
        </button>
      </div>
    </div>
  </div>
</template> 