<script setup>
import { ref } from 'vue'
import TextEditModal from '../TextEditModal.vue'

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

// 字体列表
const fontFamilies = [
  { value: '宋体', label: '宋体' },
  { value: '黑体', label: '黑体' },
  { value: '楷体', label: '楷体' },
  { value: '仿宋', label: '仿宋' }
]
</script>

<template>
  <div>
    <div class="mb-4">
      <label class="block mb-2">文本内容</label>
      <TextEditModal
        v-model="component.content"
        title="编辑文本内容"
        :rows="8"
        @update:modelValue="value => updateProperty('content', value)"
      />
    </div>
    
    <div class="grid grid-cols-2 gap-4 mb-4">
      <div>
        <label class="block mb-2">字体</label>
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
        <label class="block mb-2">字号</label>
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
    
    <div class="mb-4">
      <label class="block mb-2">行高</label>
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
  </div>
</template> 