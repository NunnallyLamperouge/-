<template>
  <div
    v-show="visible"
    class="fixed bg-white rounded-lg shadow-lg p-3 z-50 border border-gray-200"
    :style="{ top: `${y}px`, left: `${x}px` }"
  >
    <div class="text-sm text-gray-600 mb-2">已选择文本：</div>
    <div class="text-sm text-gray-800 mb-3 max-w-md overflow-hidden text-ellipsis">
      {{ selectedText }}
    </div>
    <a-button
      type="primary"
      size="small"
      :loading="loading"
      @click="handleEnhance"
    >
      智能润色
    </a-button>
    <a-button
      type="info"
      size="small"
      :loading="loading"
      @click="close()"
    >
      取消
    </a-button>
    <div v-if="options.result">
    <div class="text-sm text-gray-600 mb-2">润色结果：</div>
    <div class="text-sm text-gray-800 mb-3 max-w-md overflow-hidden text-ellipsis">
      {{ options.result }}
    </div>
    <a-button
      type="primary"
      size="small"
      :loading="loading"
      @click="copyResult"
    >
      复制结果
    </a-button>
  </div>
  </div>
  
</template>

<script setup>
import { defineProps, defineEmits , ref, reactive } from 'vue'
import { postJSON } from '@/utility/request'
import { message } from 'ant-design-vue'
import { copyText } from 'vue3-clipboard'

const options = reactive({
    result: '',
})
const copyResult = () => {
    copyText(options.result)
    message.success('复制成功')
    close()
}

const close = () => {
    emit('close')
    options.result = ''
}
const props = defineProps({
  visible: Boolean,
  x: Number,
  y: Number,
  selectedText: String
})

const emit = defineEmits(['update', 'close'])

const loading = ref(false)

const handleEnhance = async () => {
  if (!props.selectedText) return
  
  loading.value = true
  postJSON(
    '/api/template/adverise',
    { source: props.selectedText },
    (msg, data) => {
        options.result = data
      message.success('文本润色成功')
    },
    (error) => {
      message.error('文本润色失败：' + error)
    }
  )
  loading.value = false
}
</script>

<style scoped>
.max-w-md {
  max-width: 28rem;
}
</style> 