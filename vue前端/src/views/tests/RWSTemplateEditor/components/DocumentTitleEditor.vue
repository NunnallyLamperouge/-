<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

const updateTitle = (key, value) => {
  const updatedTitle = { ...props.modelValue, [key]: value }
  emit('update:modelValue', updatedTitle)
}

// 字体选项
const fonts = [
  { value: '黑体', label: '黑体' },
  { value: '宋体', label: '宋体' },
  { value: '楷体', label: '楷体' },
  { value: '微软雅黑', label: '微软雅黑' },
  { value: 'Arial', label: 'Arial' },
  { value: 'Times New Roman', label: 'Times New Roman' }
]

// 字号选项
const fontSizes = [
  { value: 16, label: '八号 (4.2mm)' },
  { value: 24, label: '小四 (5.25mm)' },
  { value: 28, label: '四号 (6.3mm)' },
  { value: 32, label: '三号 (7.35mm)' },
  { value: 36, label: '小二 (8.4mm)' },
  { value: 44, label: '二号 (10.5mm)' },
  { value: 54, label: '小一 (12.6mm)' },
  { value: 63, label: '一号 (14.7mm)' }
]

// 对齐方式选项
const alignments = [
  { value: 'LEFT', label: '左对齐' },
  { value: 'CENTER', label: '居中' },
  { value: 'RIGHT', label: '右对齐' },
  { value: 'JUSTIFIED', label: '两端对齐' }
]
</script>

<template>
  <div class="border p-4 rounded">
    <h2 class="text-lg font-bold mb-4">文档标题设置</h2>
    
    <div class="form-group">
      <label>标题文本</label>
      <input 
        :value="modelValue.text" 
        @input="updateTitle('text', $event.target.value)"
        class="form-input"
        placeholder="输入文档标题"
      />
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="form-group">
        <label>字体</label>
        <select 
          :value="modelValue.fontName" 
          @change="updateTitle('fontName', $event.target.value)"
          class="form-input"
        >
          <option v-for="font in fonts" :key="font.value" :value="font.value">
            {{ font.label }}
          </option>
        </select>
      </div>
      
      <div class="form-group">
        <label>字号</label>
        <select 
          :value="modelValue.fontSize" 
          @change="updateTitle('fontSize', parseInt($event.target.value))"
          class="form-input"
        >
          <option v-for="size in fontSizes" :key="size.value" :value="size.value">
            {{ size.label }}
          </option>
        </select>
      </div>
      
      <div class="form-group">
        <label>对齐方式</label>
        <select 
          :value="modelValue.alignment" 
          @change="updateTitle('alignment', $event.target.value)"
          class="form-input"
        >
          <option v-for="align in alignments" :key="align.value" :value="align.value">
            {{ align.label }}
          </option>
        </select>
      </div>
    </div>
    
    <!-- 预览 -->
    <div class="mt-4 p-4 border rounded bg-gray-50">
      <h3 class="text-sm text-gray-500 mb-2">预览效果：</h3>
      <div 
        class="p-4 border-dashed border-2 border-gray-300"
        :style="{
          fontFamily: modelValue.fontName,
          fontSize: `${modelValue.fontSize / 4}px`,
          textAlign: modelValue.alignment.toLowerCase()
        }"
      >
        {{ modelValue.text || '请输入标题文本' }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.form-group {
  @apply mb-4;
}

.form-group label {
  @apply block mb-2 font-medium;
}

.form-input {
  @apply w-full px-3 py-2 border rounded;
}
</style> 