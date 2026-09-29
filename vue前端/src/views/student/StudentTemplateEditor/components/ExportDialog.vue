<script setup>
import { ref, defineProps, defineEmits } from 'vue'

const props = defineProps({
  isExporting: {
    type: Boolean,
    default: false
  },
  exportProgress: {
    type: Number,
    default: 0
  },
  exportStatus: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close'])

const handleClose = () => {
  if (props.exportProgress === 100) {
    emit('close')
  }
}
</script>

<template>
  <div v-if="isExporting" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white p-6 rounded-lg w-96 shadow-xl">
      <h2 class="text-xl font-bold mb-4">导出文档</h2>
      <div class="mb-4">
        <div class="flex justify-between mb-1">
          <span>{{ exportStatus }}</span>
          <span>{{ exportProgress }}%</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2.5">
          <div 
            class="bg-blue-600 h-2.5 rounded-full transition-all duration-300" 
            :style="{ width: `${exportProgress}%` }"
          ></div>
        </div>
      </div>
      <p class="text-sm text-gray-500 mt-2">
        正在处理图片和文档内容，请稍候...
      </p>
      <div v-if="exportProgress === 100" class="mt-4 flex justify-end">
        <button 
          class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          @click="handleClose"
        >
          关闭
        </button>
      </div>
    </div>
  </div>
</template> 