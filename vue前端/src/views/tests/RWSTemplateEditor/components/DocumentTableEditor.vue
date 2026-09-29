<script setup>
import { defineProps, defineEmits, ref } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

// 内容类型选项
const contentTypes = [
  { value: 'text', label: '单行文本' },
  { value: 'longtext', label: '多行文本' },
  { value: 'schedule', label: '工作进度表' },
  { value: 'references', label: '参考文献列表' },
  { value: 'list', label: '普通列表' }
]

// 新内容项
const newContent = ref({
  title: '',
  type: 'longtext',
  minLength: 0,
  value: '',
  minCount: 1,
  required: true
})

// 是否处于编辑状态
const editingIndex = ref(-1)

// 生成唯一ID
const generateId = () => {
  return Date.now() + Math.floor(Math.random() * 1000)
}

// 添加内容项
const addContent = () => {
  if (!newContent.value.title) {
    alert('请填写内容标题')
    return
  }
  
  // 根据类型设置默认值
  let value = ''
  if (newContent.value.type === 'schedule') {
    value = [
      { task: '完成任务1', startDate: '', endDate: '' },
      { task: '完成任务2', startDate: '', endDate: '' }
    ]
  } else if (newContent.value.type === 'references' || newContent.value.type === 'list') {
    value = ['', '']
  }
  
  const contentToAdd = {
    ...newContent.value,
    id: generateId(),
    value: newContent.value.value || value
  }
  
  if (editingIndex.value >= 0) {
    // 更新现有内容
    const updatedContents = [...props.modelValue]
    
    // 保留原有数据结构
    if (['schedule', 'references', 'list'].includes(newContent.value.type) && 
        !Array.isArray(updatedContents[editingIndex.value].value)) {
      if (newContent.value.type === 'schedule') {
        contentToAdd.value = [
          { task: '完成任务1', startDate: '', endDate: '' },
          { task: '完成任务2', startDate: '', endDate: '' }
        ]
      } else {
        contentToAdd.value = ['', '']
      }
    } else if (!['schedule', 'references', 'list'].includes(newContent.value.type) && 
               Array.isArray(updatedContents[editingIndex.value].value)) {
      contentToAdd.value = ''
    } else {
      contentToAdd.value = updatedContents[editingIndex.value].value
    }
    
    updatedContents[editingIndex.value] = {
      ...contentToAdd,
      id: props.modelValue[editingIndex.value].id // 保留原始ID
    }
    
    emit('update:modelValue', updatedContents)
    editingIndex.value = -1
  } else {
    // 添加新内容
    emit('update:modelValue', [...props.modelValue, contentToAdd])
  }
  
  // 重置表单
  newContent.value = {
    title: '',
    type: 'longtext',
    minLength: 0,
    value: '',
    minCount: 1,
    required: true
  }
}

// 编辑内容
const editContent = (content, index) => {
  newContent.value = { ...content }
  editingIndex.value = index
}

// 删除内容
const deleteContent = (index) => {
  if (confirm('确定删除此内容项？')) {
    const updatedContents = [...props.modelValue]
    updatedContents.splice(index, 1)
    emit('update:modelValue', updatedContents)
  }
}

// 移动内容项
const moveContent = (index, direction) => {
  if ((direction === -1 && index === 0) || 
      (direction === 1 && index === props.modelValue.length - 1)) {
    return
  }
  
  const updatedContents = [...props.modelValue]
  const temp = updatedContents[index]
  updatedContents[index] = updatedContents[index + direction]
  updatedContents[index + direction] = temp
  
  emit('update:modelValue', updatedContents)
}

// 根据类型获取表单项说明
const getTypeDescription = (type) => {
  switch (type) {
    case 'longtext':
      return '可设置最小字数的文本区域，适合研究背景、内容描述等'
    case 'schedule':
      return '工作进度表，包含任务名称和时间段'
    case 'references':
      return '参考文献列表，可添加多条参考文献'
    case 'list':
      return '普通列表项，可添加多个项目'
    default:
      return '普通文本字段'
  }
}
</script>

<template>
  <div class="border p-4 rounded">
    <h2 class="text-lg font-bold mb-4">表格内容设置</h2>
    
    <!-- 已有内容列表 -->
    <div class="mb-6">
      <h3 class="text-md font-semibold mb-2">当前内容项</h3>
      
      <div v-if="modelValue.length === 0" class="p-4 bg-gray-100 rounded text-center text-gray-500">
        暂无内容项，请添加新内容
      </div>
      
      <div 
        v-for="(content, index) in modelValue" 
        :key="content.id"
        class="p-3 border rounded mb-2 bg-white"
      >
        <div class="flex justify-between items-center">
          <div>
            <span class="font-medium">{{ content.title }}</span>
            <span class="ml-2 text-xs px-2 py-0.5 rounded bg-blue-100 text-blue-800">
              {{ contentTypes.find(t => t.value === content.type)?.label || content.type }}
            </span>
            <span class="ml-2 text-xs px-2 py-0.5 rounded" 
                  :class="content.required ? 'bg-red-100 text-red-800' : 'bg-gray-100'">
              {{ content.required ? '必填' : '选填' }}
            </span>
            
            <div class="text-sm text-gray-600 mt-1">
              <template v-if="content.type === 'longtext' && content.minLength > 0">
                最小字数: {{ content.minLength }} 字
              </template>
              <template v-else-if="['references', 'list'].includes(content.type) && content.minCount > 0">
                最少项目: {{ content.minCount }} 项
              </template>
            </div>
          </div>
          
          <div class="flex gap-2">
            <button 
              @click="moveContent(index, -1)" 
              class="p-1 text-gray-600 hover:text-blue-600"
              :disabled="index === 0"
              :class="{'opacity-30': index === 0}"
            >
              ↑
            </button>
            <button 
              @click="moveContent(index, 1)" 
              class="p-1 text-gray-600 hover:text-blue-600"
              :disabled="index === modelValue.length - 1"
              :class="{'opacity-30': index === modelValue.length - 1}"
            >
              ↓
            </button>
            <button @click="editContent(content, index)" class="p-1 text-blue-600 hover:text-blue-800">
              编辑
            </button>
            <button @click="deleteContent(index)" class="p-1 text-red-600 hover:text-red-800">
              删除
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 添加/编辑内容表单 -->
    <div class="border p-4 rounded bg-gray-50">
      <h3 class="text-md font-semibold mb-2">
        {{ editingIndex >= 0 ? '编辑内容项' : '添加新内容项' }}
      </h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="form-group md:col-span-2">
          <label>内容标题</label>
          <input 
            v-model="newContent.title" 
            class="form-input"
            placeholder="例如：研究背景及意义"
          />
        </div>
        
        <div class="form-group">
          <label>内容类型</label>
          <select v-model="newContent.type" class="form-input">
            <option v-for="type in contentTypes" :key="type.value" :value="type.value">
              {{ type.label }}
            </option>
          </select>
          <div class="text-xs text-gray-500 mt-1">
            {{ getTypeDescription(newContent.type) }}
          </div>
        </div>
        
        <div class="form-group">
          <label>是否必填</label>
          <div class="mt-2">
            <label class="inline-flex items-center">
              <input type="checkbox" v-model="newContent.required" class="form-checkbox" />
              <span class="ml-2">必填内容</span>
            </label>
          </div>
        </div>
        
        <!-- 根据类型显示不同的选项 -->
        <template v-if="newContent.type === 'longtext'">
          <div class="form-group">
            <label>最小字数限制</label>
            <input 
              v-model.number="newContent.minLength" 
              type="number" 
              min="0" 
              class="form-input"
            />
          </div>
        </template>
        
        <template v-if="['references', 'list'].includes(newContent.type)">
          <div class="form-group">
            <label>最少项目数</label>
            <input 
              v-model.number="newContent.minCount" 
              type="number" 
              min="0" 
              class="form-input"
            />
          </div>
        </template>
      </div>
      
      <div class="flex justify-end mt-4">
        <button 
          v-if="editingIndex >= 0"
          @click="editingIndex = -1; newContent = {title: '', type: 'longtext', minLength: 0, value: '', minCount: 1, required: true}"
          class="px-4 py-2 bg-gray-500 text-white rounded mr-2"
        >
          取消
        </button>
        <button 
          @click="addContent"
          class="px-4 py-2 bg-blue-500 text-white rounded"
        >
          {{ editingIndex >= 0 ? '保存修改' : '添加内容项' }}
        </button>
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

.form-checkbox {
  @apply rounded text-blue-600;
}
</style> 