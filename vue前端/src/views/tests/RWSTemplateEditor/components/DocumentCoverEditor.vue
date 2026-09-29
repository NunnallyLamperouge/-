<script setup>
import { defineProps, defineEmits, ref } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

// 字段类型选项
const fieldTypes = [
  { value: 'text', label: '文本' },
  { value: 'date', label: '日期' },
  { value: 'number', label: '数字' },
  { value: 'select', label: '下拉选择' }
]

// 新字段
const newField = ref({
  label: '',
  field: '',
  type: 'text',
  required: true,
  value: '',
  options: [] // 用于select类型
})

// 新选项（用于select类型）
const newOption = ref('')

// 是否处于编辑状态
const editingIndex = ref(-1)

// 生成唯一ID
const generateId = () => {
  return Date.now() + Math.floor(Math.random() * 1000)
}

// 添加字段
const addField = () => {
  if (!newField.value.label || !newField.value.field) {
    alert('请填写字段名称和字段标识')
    return
  }
  
  // 检查字段标识是否重复
  const isDuplicate = props.modelValue.some(field => field.field === newField.value.field)
  if (isDuplicate && editingIndex.value === -1) {
    alert('字段标识已存在，请使用其他标识')
    return
  }
  
  const fieldToAdd = {
    ...newField.value,
    id: generateId()
  }
  
  if (editingIndex.value >= 0) {
    // 更新现有字段
    const updatedFields = [...props.modelValue]
    updatedFields[editingIndex.value] = {
      ...updatedFields[editingIndex.value],
      ...fieldToAdd,
      id: props.modelValue[editingIndex.value].id // 保留原始ID
    }
    emit('update:modelValue', updatedFields)
    editingIndex.value = -1
  } else {
    // 添加新字段
    emit('update:modelValue', [...props.modelValue, fieldToAdd])
  }
  
  // 重置表单
  newField.value = {
    label: '',
    field: '',
    type: 'text',
    required: true,
    value: '',
    options: []
  }
  newOption.value = ''
}

// 编辑字段
const editField = (field, index) => {
  newField.value = { ...field }
  editingIndex.value = index
  
  // 确保options是数组
  if (!newField.value.options) {
    newField.value.options = []
  }
}

// 删除字段
const deleteField = (index) => {
  if (confirm('确定删除此字段？')) {
    const updatedFields = [...props.modelValue]
    updatedFields.splice(index, 1)
    emit('update:modelValue', updatedFields)
  }
}

// 添加选项（用于select类型）
const addOption = () => {
  if (!newOption.value.trim()) return
  
  newField.value.options.push(newOption.value)
  newOption.value = ''
}

// 删除选项
const removeOption = (optionIndex) => {
  newField.value.options.splice(optionIndex, 1)
}

// 移动字段
const moveField = (index, direction) => {
  if ((direction === -1 && index === 0) || 
      (direction === 1 && index === props.modelValue.length - 1)) {
    return
  }
  
  const updatedFields = [...props.modelValue]
  const temp = updatedFields[index]
  updatedFields[index] = updatedFields[index + direction]
  updatedFields[index + direction] = temp
  
  emit('update:modelValue', updatedFields)
}
</script>

<template>
  <div class="border p-4 rounded">
    <h2 class="text-lg font-bold mb-4">封面字段设置</h2>
    
    <!-- 已有字段列表 -->
    <div class="mb-6">
      <h3 class="text-md font-semibold mb-2">当前字段列表</h3>
      
      <div v-if="modelValue.length === 0" class="p-4 bg-gray-100 rounded text-center text-gray-500">
        暂无字段，请添加新字段
      </div>
      
      <div 
        v-for="(field, index) in modelValue" 
        :key="field.id"
        class="p-3 border rounded mb-2 bg-white"
      >
        <div class="flex justify-between items-center">
          <div>
            <span class="font-medium">{{ field.label }}</span>
            <span class="ml-2 text-gray-500 text-sm">({{ field.field }})</span>
            <span class="ml-2 text-xs px-2 py-0.5 rounded" 
                  :class="field.required ? 'bg-red-100 text-red-800' : 'bg-gray-100'">
              {{ field.required ? '必填' : '选填' }}
            </span>
            <span class="ml-2 text-xs px-2 py-0.5 rounded bg-blue-100 text-blue-800">
              {{ fieldTypes.find(t => t.value === field.type)?.label || field.type }}
            </span>
          </div>
          
          <div class="flex gap-2">
            <button 
              @click="moveField(index, -1)" 
              class="p-1 text-gray-600 hover:text-blue-600"
              :disabled="index === 0"
              :class="{'opacity-30': index === 0}"
            >
              ↑
            </button>
            <button 
              @click="moveField(index, 1)" 
              class="p-1 text-gray-600 hover:text-blue-600"
              :disabled="index === modelValue.length - 1"
              :class="{'opacity-30': index === modelValue.length - 1}"
            >
              ↓
            </button>
            <button @click="editField(field, index)" class="p-1 text-blue-600 hover:text-blue-800">
              编辑
            </button>
            <button @click="deleteField(index)" class="p-1 text-red-600 hover:text-red-800">
              删除
            </button>
          </div>
        </div>
        
        <!-- 如果是select类型，显示选项 -->
        <div v-if="field.type === 'select' && field.options && field.options.length > 0" class="mt-2 pl-4 text-sm text-gray-600">
          <span>选项：{{ field.options.join(', ') }}</span>
        </div>
      </div>
    </div>
    
    <!-- 添加/编辑字段表单 -->
    <div class="border p-4 rounded bg-gray-50">
      <h3 class="text-md font-semibold mb-2">
        {{ editingIndex >= 0 ? '编辑字段' : '添加新字段' }}
      </h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="form-group">
          <label>字段名称</label>
          <input 
            v-model="newField.label" 
            class="form-input"
            placeholder="例如：学生姓名"
          />
        </div>
        
        <div class="form-group">
          <label>字段标识</label>
          <input 
            v-model="newField.field" 
            class="form-input"
            placeholder="例如：studentName"
            :disabled="editingIndex >= 0"
          />
          <div v-if="editingIndex >= 0" class="text-xs text-gray-500 mt-1">
            编辑模式下不能修改字段标识
          </div>
        </div>
        
        <div class="form-group">
          <label>字段类型</label>
          <select v-model="newField.type" class="form-input">
            <option v-for="type in fieldTypes" :key="type.value" :value="type.value">
              {{ type.label }}
            </option>
          </select>
        </div>
        
        <div class="form-group">
          <label>是否必填</label>
          <div class="mt-2">
            <label class="inline-flex items-center">
              <input type="checkbox" v-model="newField.required" class="form-checkbox" />
              <span class="ml-2">必填字段</span>
            </label>
          </div>
        </div>
        
        <div v-if="newField.type === 'select'" class="form-group md:col-span-2">
          <label>选项列表</label>
          <div class="flex gap-2 mb-2">
            <input v-model="newOption" class="form-input flex-1" placeholder="添加新选项" />
            <button @click="addOption" class="px-4 py-2 bg-blue-500 text-white rounded">
              添加选项
            </button>
          </div>
          
          <div class="flex flex-wrap gap-2 mt-2">
            <div 
              v-for="(option, optionIndex) in newField.options" 
              :key="optionIndex"
              class="px-3 py-1 bg-blue-100 rounded flex items-center"
            >
              {{ option }}
              <button @click="removeOption(optionIndex)" class="ml-2 text-red-600">
                &times;
              </button>
            </div>
          </div>
        </div>
        
        <div v-if="newField.type === 'text'" class="form-group">
          <label>默认值 (可选)</label>
          <input v-model="newField.value" class="form-input" placeholder="默认值" />
        </div>
      </div>
      
      <div class="flex justify-end mt-4">
        <button 
          v-if="editingIndex >= 0"
          @click="editingIndex = -1; newField = {label: '', field: '', type: 'text', required: true, value: '', options: []}"
          class="px-4 py-2 bg-gray-500 text-white rounded mr-2"
        >
          取消
        </button>
        <button 
          @click="addField"
          class="px-4 py-2 bg-blue-500 text-white rounded"
        >
          {{ editingIndex >= 0 ? '保存修改' : '添加字段' }}
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