<script setup>
import { defineProps, defineEmits, ref, computed } from 'vue'
import axios from 'axios';

const props = defineProps({
  modelValue: {
    type: Object,
    required: true
  },
  isValid: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'generate'])

// 更新封面字段值
const updateCoverFieldValue = (field, value) => {
  const updatedFields = [...props.modelValue.coverFields]
  const fieldIndex = updatedFields.findIndex(f => f.id === field.id)
  
  if (fieldIndex !== -1) {
    updatedFields[fieldIndex] = { ...updatedFields[fieldIndex], value }
    
    const updatedTemplate = {
      ...props.modelValue,
      coverFields: updatedFields
    }
    
    emit('update:modelValue', updatedTemplate)
  }
}

// 更新表格内容项值
const updateContentValue = (content, value) => {
  const updatedContents = [...props.modelValue.tableContents]
  const contentIndex = updatedContents.findIndex(c => c.id === content.id)
  
  if (contentIndex !== -1) {
    updatedContents[contentIndex] = { ...updatedContents[contentIndex], value }
    
    const updatedTemplate = {
      ...props.modelValue,
      tableContents: updatedContents
    }
    
    emit('update:modelValue', updatedTemplate)
  }
}

// 添加一个工作进度项
const addScheduleItem = (content) => {
  const updatedContents = [...props.modelValue.tableContents]
  const contentIndex = updatedContents.findIndex(c => c.id === content.id)
  
  if (contentIndex !== -1) {
    const newValue = [...updatedContents[contentIndex].value, {
      task: '',
      startDate: '',
      endDate: ''
    }]
    
    updatedContents[contentIndex] = { 
      ...updatedContents[contentIndex], 
      value: newValue
    }
    
    const updatedTemplate = {
      ...props.modelValue,
      tableContents: updatedContents
    }
    
    emit('update:modelValue', updatedTemplate)
  }
}

// 删除一个工作进度项
const removeScheduleItem = (content, index) => {
  const updatedContents = [...props.modelValue.tableContents]
  const contentIndex = updatedContents.findIndex(c => c.id === content.id)
  
  if (contentIndex !== -1 && updatedContents[contentIndex].value.length > 1) {
    const newValue = [...updatedContents[contentIndex].value]
    newValue.splice(index, 1)
    
    updatedContents[contentIndex] = { 
      ...updatedContents[contentIndex], 
      value: newValue
    }
    
    const updatedTemplate = {
      ...props.modelValue,
      tableContents: updatedContents
    }
    
    emit('update:modelValue', updatedTemplate)
  }
}

// 更新工作进度项
const updateScheduleItem = (content, index, key, value) => {
  const updatedContents = [...props.modelValue.tableContents]
  const contentIndex = updatedContents.findIndex(c => c.id === content.id)
  
  if (contentIndex !== -1) {
    const newValue = [...updatedContents[contentIndex].value]
    newValue[index] = { ...newValue[index], [key]: value }
    
    updatedContents[contentIndex] = { 
      ...updatedContents[contentIndex], 
      value: newValue
    }
    
    const updatedTemplate = {
      ...props.modelValue,
      tableContents: updatedContents
    }
    
    emit('update:modelValue', updatedTemplate)
  }
}

// 添加一个列表项
const addListItem = (content) => {
  const updatedContents = [...props.modelValue.tableContents]
  const contentIndex = updatedContents.findIndex(c => c.id === content.id)
  
  if (contentIndex !== -1) {
    const newValue = [...updatedContents[contentIndex].value, '']
    
    updatedContents[contentIndex] = { 
      ...updatedContents[contentIndex], 
      value: newValue
    }
    
    const updatedTemplate = {
      ...props.modelValue,
      tableContents: updatedContents
    }
    
    emit('update:modelValue', updatedTemplate)
  }
}

// 删除一个列表项
const removeListItem = (content, index) => {
  const updatedContents = [...props.modelValue.tableContents]
  const contentIndex = updatedContents.findIndex(c => c.id === content.id)
  
  if (contentIndex !== -1 && updatedContents[contentIndex].value.length > 1) {
    const newValue = [...updatedContents[contentIndex].value]
    newValue.splice(index, 1)
    
    updatedContents[contentIndex] = { 
      ...updatedContents[contentIndex], 
      value: newValue
    }
    
    const updatedTemplate = {
      ...props.modelValue,
      tableContents: updatedContents
    }
    
    emit('update:modelValue', updatedTemplate)
  }
}

// 更新列表项
const updateListItem = (content, index, value) => {
  const updatedContents = [...props.modelValue.tableContents]
  const contentIndex = updatedContents.findIndex(c => c.id === content.id)
  
  if (contentIndex !== -1) {
    const newValue = [...updatedContents[contentIndex].value]
    newValue[index] = value
    
    updatedContents[contentIndex] = { 
      ...updatedContents[contentIndex], 
      value: newValue
    }
    
    const updatedTemplate = {
      ...props.modelValue,
      tableContents: updatedContents
    }
    
    emit('update:modelValue', updatedTemplate)
  }
}

// 生成文档
const generateDocument = () => {
  const formData = new FormData();
  formData.append('coverFields', JSON.stringify(props.modelValue.coverFields));
  formData.append('tableContents', JSON.stringify(props.modelValue.tableContents));
  console.log(formData)
  axios.post('http://localhost:8080/api/receive', formData, {
    headers: {
      'Content-Type': 'multipart/form-data' // 使用 multipart/form-data
    }
  })
      .then((response) => {
        console.log('上传成功:', response.data);
       // emit('generate')
      })
      .catch((error) => {
        console.error('上传失败:', error);
      });
}
</script>

<template>
  <div class="border p-4 rounded">
    <h2 class="text-lg font-bold mb-4">预览与填写</h2>
    
    <!-- 封面字段预览 -->
    <div class="mb-6">
      <h3 class="text-md font-semibold border-b pb-2 mb-4">封面信息</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div 
          v-for="field in modelValue.coverFields" 
          :key="field.id"
          class="form-group"
        >
          <label>
            {{ field.label }}
            <span v-if="field.required" class="text-red-600">*</span>
          </label>
          
          <!-- 不同类型的输入控件 -->
          <input 
            v-if="field.type === 'text' || field.type === 'number'"
            :value="field.value"
            :type="field.type"
            @input="updateCoverFieldValue(field, $event.target.value)"
            class="form-input"
            :placeholder="`请输入${field.label}`"
            :required="field.required"
          />
          
          <input 
            v-else-if="field.type === 'date'"
            :value="field.value"
            type="date"
            @input="updateCoverFieldValue(field, $event.target.value)"
            class="form-input"
            :required="field.required"
          />
          
          <select 
            v-else-if="field.type === 'select'"
            :value="field.value"
            @change="updateCoverFieldValue(field, $event.target.value)"
            class="form-input"
            :required="field.required"
          >
            <option value="">请选择{{ field.label }}</option>
            <option v-for="(option, index) in field.options" :key="index" :value="option">
              {{ option }}
            </option>
          </select>
        </div>
      </div>
    </div>
    
    <!-- 表格内容预览 -->
    <div class="mb-6">
      <h3 class="text-md font-semibold border-b pb-2 mb-4">表格内容</h3>
      
      <div 
        v-for="content in modelValue.tableContents" 
        :key="content.id"
        class="mb-6"
      >
        <div class="form-group">
          <label>
            {{ content.title }}
            <span v-if="content.required" class="text-red-600">*</span>
            
            <span v-if="content.type === 'longtext' && content.minLength > 0" class="text-sm text-gray-500 ml-2">
              (最少 {{ content.minLength }} 字)
            </span>
            <span v-else-if="['references', 'list'].includes(content.type) && content.minCount > 0" class="text-sm text-gray-500 ml-2">
              (最少 {{ content.minCount }} 项)
            </span>
          </label>
          
          <!-- 不同类型的输入控件 -->
          <input 
            v-if="content.type === 'text'"
            :value="content.value"
            type="text"
            @input="updateContentValue(content, $event.target.value)"
            class="form-input"
            :placeholder="`请输入${content.title}`"
          />
          
          <textarea 
            v-else-if="content.type === 'longtext'"
            :value="content.value"
            @input="updateContentValue(content, $event.target.value)"
            class="form-textarea"
            rows="6"
            :placeholder="`请输入${content.title}`"
          ></textarea>
          
          <!-- 工作进度表 -->
          <div v-else-if="content.type === 'schedule'">
            <div 
              v-for="(item, index) in content.value" 
              :key="index"
              class="grid grid-cols-1 md:grid-cols-3 gap-2 mb-2"
            >
              <input 
                :value="item.task"
                @input="updateScheduleItem(content, index, 'task', $event.target.value)"
                class="form-input"
                placeholder="任务名称"
              />
              <input 
                :value="item.startDate"
                @input="updateScheduleItem(content, index, 'startDate', $event.target.value)"
                type="date"
                class="form-input"
              />
              <div class="flex items-center">
                <input 
                  :value="item.endDate"
                  @input="updateScheduleItem(content, index, 'endDate', $event.target.value)"
                  type="date"
                  class="form-input"
                />
                <button 
                  @click="removeScheduleItem(content, index)"
                  class="ml-2 p-2 text-red-600"
                  :disabled="content.value.length <= 1"
                >
                  &times;
                </button>
              </div>
            </div>
            
            <button 
              @click="addScheduleItem(content)"
              class="mt-2 px-4 py-2 bg-blue-500 text-white rounded"
            >
              添加进度项
            </button>
          </div>
          
          <!-- 列表项 (参考文献或普通列表) -->
          <div v-else-if="['references', 'list'].includes(content.type)">
            <div 
              v-for="(item, index) in content.value" 
              :key="index"
              class="flex items-center mb-2"
            >
              <div class="mr-2 text-gray-500">{{ index + 1 }}.</div>
              <input 
                :value="item"
                @input="updateListItem(content, index, $event.target.value)"
                class="form-input flex-1"
                :placeholder="content.type === 'references' ? '参考文献' : '列表项'"
              />
              <button 
                @click="removeListItem(content, index)"
                class="ml-2 p-2 text-red-600"
                :disabled="content.value.length <= 1"
              >
                &times;
              </button>
            </div>
            
            <button 
              @click="addListItem(content)"
              class="mt-2 px-4 py-2 bg-blue-500 text-white rounded"
            >
              添加{{ content.type === 'references' ? '参考文献' : '列表项' }}
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 生成按钮 -->
    <div class="flex justify-center">
      <button
        @click="generateDocument"
        class="px-6 py-3 bg-green-500 text-white rounded-lg font-bold text-lg hover:bg-green-600 transition"
        :disabled="!isValid"
        :class="{'opacity-50 cursor-not-allowed': !isValid}"
      >
        生成文档
      </button>
    </div>
    
    <div v-if="!isValid" class="mt-4 text-center text-red-500">
      请完成所有必填字段后再生成文档
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

.form-textarea {
  @apply w-full px-3 py-2 border rounded;
}
</style> 