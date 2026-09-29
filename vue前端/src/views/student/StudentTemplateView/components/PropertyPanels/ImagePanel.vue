<script setup>
import { ref } from 'vue'

const props = defineProps({
  component: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update'])

// 上传状态
const isUploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref('')

// 更新组件属性的辅助函数
const updateProperty = (property, value) => {
  emit('update', props.component.id, property, value)
}

// 处理图片上传
const handleImageUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  try {
    isUploading.value = true
    uploadError.value = ''
    uploadProgress.value = 10
    
    const formData = new FormData()
    formData.append('file', file)
    
    // 创建 XMLHttpRequest 以便跟踪上传进度
    const xhr = new XMLHttpRequest()
    
    // 设置进度监听
    xhr.upload.addEventListener('progress', (e) => {
      if (e.lengthComputable) {
        uploadProgress.value = Math.round((e.loaded / e.total) * 90) + 10
      }
    })
    
    // 创建 Promise 包装 XHR
    const uploadPromise = new Promise((resolve, reject) => {
      xhr.open('POST', 'http://localhost:8080/file/upload')
      
      xhr.onload = () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          resolve(xhr.responseText)
        } else {
          reject(new Error(`上传失败: ${xhr.statusText}`))
        }
      }
      
      xhr.onerror = () => reject(new Error('网络错误'))
      xhr.send(formData)
    })
    
    // 等待上传完成
    const response = await uploadPromise
    
    // 拼接完整的图片URL
    const imageUrl = `http://localhost:8080${JSON.parse(response).url}`
    
    // 更新组件属性
    updateProperty('imageUrl', imageUrl)
    
    // 重置图片预览错误状态
    updateProperty('imagePreviewError', false)
    
    uploadProgress.value = 100
    setTimeout(() => {
      isUploading.value = false
      uploadProgress.value = 0
    }, 1000)
    
    // 清空文件输入，允许再次选择相同文件
    event.target.value = ''
    
  } catch (error) {
    console.error('图片上传失败:', error)
    uploadError.value = error.message || '上传失败'
    uploadProgress.value = 0
    setTimeout(() => {
      isUploading.value = false
    }, 3000)
    
    // 清空文件输入
    event.target.value = ''
  }
}

// 字体列表
const fontFamilies = [
  { value: '宋体', label: '宋体' },
  { value: '黑体', label: '黑体' },
  { value: '楷体', label: '楷体' },
  { value: '仿宋', label: '仿宋' },
  { value: 'Arial', label: 'Arial' },
  { value: 'Times New Roman', label: 'Times New Roman' }
]

// 对齐方式
const alignments = [
  { value: 'left', label: '左对齐', icon: '⬅️' },
  { value: 'center', label: '居中', icon: '⬆️' },
  { value: 'right', label: '右对齐', icon: '➡️' },
  { value: 'justify', label: '两端对齐', icon: '↔️' }
]
</script>

<template>
  <div>
    <div class="mb-4">
      <label class="block mb-2">图片URL</label>
      <div class="flex gap-2">
        <input
          v-model="component.imageUrl"
          type="text"
          class="flex-1 px-3 py-2 border rounded"
          placeholder="输入图片URL"
          @input="updateProperty('imageUrl', $event.target.value)"
        />
        <label class="px-3 py-2 bg-blue-500 text-white rounded cursor-pointer hover:bg-blue-600">
          上传
          <input
            type="file"
            accept="image/*"
            class="hidden"
            @change="handleImageUpload"
          />
        </label>
      </div>
      
      <!-- 上传进度条 -->
      <div v-if="isUploading" class="mt-2">
        <div class="flex justify-between text-sm mb-1">
          <span>上传中...</span>
          <span>{{ uploadProgress }}%</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div 
            class="bg-blue-600 h-2 rounded-full transition-all duration-300" 
            :style="{ width: `${uploadProgress}%` }"
          ></div>
        </div>
        <p v-if="uploadError" class="text-red-500 text-sm mt-1">{{ uploadError }}</p>
      </div>
      
      <!-- 图片预览 -->
      <div v-if="component.imageUrl" class="mt-2">
        <img 
          :src="component.imageUrl" 
          alt="图片预览" 
          class="max-w-full h-auto border rounded"
          @error="updateProperty('imagePreviewError', true)"
        />
        <p v-if="component.imagePreviewError" class="text-red-500 text-sm mt-1">
          图片加载失败，请检查URL是否正确
        </p>
      </div>
    </div>

    <div class="mb-4">
      <label class="block mb-2">图片对齐方式</label>
      <div class="flex gap-2">
        <button
          v-for="align in alignments"
          :key="align.value"
          :class="[
            'px-3 py-2 border rounded flex-1',
            component.imageAlignment === align.value ? 'bg-blue-500 text-white' : 'bg-white'
          ]"
          @click="updateProperty('imageAlignment', align.value)"
        >
          {{ align.icon }}
        </button>
      </div>
    </div>

    <div class="mb-4">
      <label class="block mb-2">最大宽度 (%)</label>
      <input
        v-model.number="component.maxWidth"
        type="number"
        min="10"
        max="100"
        class="w-full px-3 py-2 border rounded"
        @input="updateProperty('maxWidth', parseInt($event.target.value))"
      />
    </div>

    <div class="mb-4">
      <label class="flex items-center">
        <input
          v-model="component.showCaption"
          type="checkbox"
          class="mr-2"
          @change="updateProperty('showCaption', component.showCaption)"
        />
        显示图片说明
      </label>
    </div>

    <template v-if="component.showCaption">
      <div class="mb-4">
        <label class="block mb-2">图片说明文字</label>
        <input
          v-model="component.caption"
          type="text"
          class="w-full px-3 py-2 border rounded"
          placeholder="例如：图1-1 系统概览"
          @input="updateProperty('caption', $event.target.value)"
        />
      </div>

      <div class="mb-4">
        <label class="block mb-2">说明文字对齐方式</label>
        <div class="flex gap-2">
          <button
            v-for="align in alignments"
            :key="align.value"
            :class="[
              'px-3 py-2 border rounded flex-1',
              component.captionAlignment === align.value ? 'bg-blue-500 text-white' : 'bg-white'
            ]"
            @click="updateProperty('captionAlignment', align.value)"
          >
            {{ align.icon }}
          </button>
        </div>
      </div>

      <div class="mb-4">
        <label class="block mb-2">说明文字字体</label>
        <select
          v-model="component.captionFontFamily"
          class="w-full px-3 py-2 border rounded"
          @change="updateProperty('captionFontFamily', $event.target.value)"
        >
          <option v-for="font in fontFamilies" :key="font.value" :value="font.value">
            {{ font.label }}
          </option>
        </select>
      </div>

      <div class="mb-4">
        <label class="block mb-2">说明文字字号</label>
        <input
          v-model.number="component.captionFontSize"
          type="number"
          min="8"
          max="36"
          class="w-full px-3 py-2 border rounded"
          @input="updateProperty('captionFontSize', parseInt($event.target.value))"
        />
      </div>

      <div class="mb-4">
        <label class="block mb-2">说明文字行高</label>
        <input
          v-model.number="component.captionLineHeight"
          type="number"
          min="1"
          max="3"
          step="0.1"
          class="w-full px-3 py-2 border rounded"
          @input="updateProperty('captionLineHeight', parseFloat($event.target.value))"
        />
      </div>
    </template>
  </div>
</template> 