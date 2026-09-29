<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue';
import { generateDocx } from '../utils/docxGenerator';

const props = defineProps({
  template: {
    type: Object,
    required: true
  }
});

// 预览数据
const previewData = reactive({});

// 初始化预览数据
const initPreviewData = () => {
  // 清空现有数据
  Object.keys(previewData).forEach(key => {
    delete previewData[key];
  });
  
  // 重新填充数据
  props.template.sections.forEach(section => {
    if (section.type === 'schedule') {
      previewData.schedule = [...(section.defaultItems || [])];
    } else if (section.type === 'references') {
      previewData.references = [];
    } else if (section.fields) {
      section.fields.forEach(field => {
        previewData[field.id] = field.defaultValue || '';
      });
    }
  });
};

// 监听模板变化，重新初始化预览数据
watch(() => props.template, () => {
  initPreviewData();
}, { deep: true });

// 添加工作计划项
const addScheduleItem = () => {
  if (!previewData.schedule) {
    previewData.schedule = [];
  }
  previewData.schedule.push({
    task: '',
    startDate: '',
    endDate: ''
  });
};

// 删除工作计划项
const removeScheduleItem = (index) => {
  previewData.schedule.splice(index, 1);
};

// 添加参考文献
const addReference = () => {
  if (!previewData.references) {
    previewData.references = [];
  }
  previewData.references.push('');
};

// 删除参考文献
const removeReference = (index) => {
  previewData.references.splice(index, 1);
};

// 验证字段
const validateField = (field, value) => {
  if (field.required && (!value || value.length === 0)) {
    return '此字段为必填项';
  }
  
  if (field.type === 'textarea') {
    if (field.minLength && value.length < field.minLength) {
      return `最少需要 ${field.minLength} 个字符`;
    }
    if (field.maxLength && value.length > field.maxLength) {
      return `最多允许 ${field.maxLength} 个字符`;
    }
  }
  
  return '';
};

// 验证整个表单
const formErrors = computed(() => {
  const errors = {};
  
  props.template.sections.forEach(section => {
    if (section.fields) {
      section.fields.forEach(field => {
        const error = validateField(field, previewData[field.id]);
        if (error) {
          errors[field.id] = error;
        }
      });
    }
  });
  
  // 验证参考文献
  const referencesSection = props.template.sections.find(s => s.id === 'references');
  if (referencesSection && referencesSection.minItems) {
    if (!previewData.references || previewData.references.length < referencesSection.minItems) {
      errors.references = `至少需要 ${referencesSection.minItems} 条参考文献`;
    }
  }
  
  return errors;
});

// 表单是否有效
const isFormValid = computed(() => {
  return Object.keys(formErrors.value).length === 0;
});

// 导出为Word文档
const exportToWord = async () => {
  if (!isFormValid.value) {
    alert('表单中存在错误，请修正后再导出');
    return;
  }
  
  try {
    await generateDocx(props.template, previewData);
  } catch (error) {
    console.error('导出文档失败:', error);
    alert(`导出文档失败: ${error.message}`);
  }
};

// 初始化预览数据
onMounted(() => {
  initPreviewData();
});
</script>

<template>
  <div class="template-preview bg-white p-4 rounded shadow-sm">
    <h2 class="text-xl font-bold mb-4">{{ template.name }} 预览</h2>
    
    <div class="mb-6">
      <div 
        v-for="section in template.sections" 
        :key="section.id"
        class="mb-8"
      >
        <h3 class="text-lg font-semibold mb-4">{{ section.title }}</h3>
        
        <!-- 普通字段 -->
        <template v-if="section.fields">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div 
              v-for="field in section.fields" 
              :key="field.id"
              class="form-group"
              :class="{ 'col-span-2': field.type === 'textarea' }"
            >
              <label class="block mb-2 font-medium">
                {{ field.label }}
                <span v-if="field.required" class="text-red-500">*</span>
              </label>
              
              <!-- 文本输入 -->
              <input 
                v-if="field.type === 'text' || field.type === 'number' || field.type === 'date'"
                v-model="previewData[field.id]"
                :type="field.type"
                class="form-input w-full"
                :class="{ 'border-red-500': formErrors[field.id] }"
                :disabled="field.disabled"
              />
              
              <!-- 多行文本 -->
              <template v-else-if="field.type === 'textarea'">
                <textarea 
                  v-model="previewData[field.id]"
                  class="form-textarea w-full"
                  rows="6"
                  :class="{ 'border-red-500': formErrors[field.id] }"
                ></textarea>
                <div class="flex justify-between text-sm text-gray-500 mt-1">
                  <span>当前字数: {{ previewData[field.id]?.length || 0 }}</span>
                  <span v-if="field.minLength || field.maxLength">
                    {{ field.minLength || 0 }} - {{ field.maxLength || '∞' }} 字
                  </span>
                </div>
              </template>
              
              <!-- 单选框 -->
              <div v-else-if="field.type === 'radio'" class="space-y-2">
                <div 
                  v-for="option in field.options" 
                  :key="option"
                  class="flex items-center"
                >
                  <input 
                    type="radio" 
                    :id="`${field.id}_${option}`"
                    :value="option"
                    v-model="previewData[field.id]"
                    class="mr-2"
                  />
                  <label :for="`${field.id}_${option}`">{{ option }}</label>
                </div>
              </div>
              
              <!-- 复选框 -->
              <div v-else-if="field.type === 'checkbox'" class="space-y-2">
                <div 
                  v-for="option in field.options" 
                  :key="option"
                  class="flex items-center"
                >
                  <input 
                    type="checkbox" 
                    :id="`${field.id}_${option}`"
                    :value="option"
                    v-model="previewData[field.id]"
                    class="mr-2"
                  />
                  <label :for="`${field.id}_${option}`">{{ option }}</label>
                </div>
              </div>
              
              <!-- 下拉选择 -->
              <select 
                v-else-if="field.type === 'select'"
                v-model="previewData[field.id]"
                class="form-select w-full"
                :class="{ 'border-red-500': formErrors[field.id] }"
              >
                <option value="">请选择</option>
                <option 
                  v-for="option in field.options" 
                  :key="option"
                  :value="option"
                >
                  {{ option }}
                </option>
              </select>
              
              <!-- 错误提示 -->
              <div 
                v-if="formErrors[field.id]" 
                class="text-red-500 text-sm mt-1"
              >
                {{ formErrors[field.id] }}
              </div>
            </div>
          </div>
        </template>
        
        <!-- 工作计划 -->
        <template v-else-if="section.type === 'schedule'">
          <div class="border rounded p-4">
            <div 
              v-for="(item, index) in previewData.schedule" 
              :key="index"
              class="grid grid-cols-3 gap-4 mb-2"
            >
              <input 
                v-model="item.task" 
                class="form-input" 
                placeholder="任务描述"
              />
              <input 
                type="date" 
                v-model="item.startDate" 
                class="form-input"
              />
              <div class="flex">
                <input 
                  type="date" 
                  v-model="item.endDate" 
                  class="form-input flex-1 mr-2"
                />
                <button 
                  @click="removeScheduleItem(index)"
                  class="px-2 py-1 bg-red-500 text-white rounded"
                >
                  删除
                </button>
              </div>
            </div>
            <button 
              @click="addScheduleItem"
              class="mt-2 px-4 py-2 bg-blue-500 text-white rounded w-full"
            >
              添加任务
            </button>
          </div>
        </template>
        
        <!-- 参考文献 -->
        <template v-else-if="section.type === 'references'">
          <div class="mb-4">
            <div 
              v-for="(ref, index) in previewData.references" 
              :key="index"
              class="flex mb-2"
            >
              <input 
                v-model="previewData.references[index]"
                class="form-input flex-1 mr-2"
                placeholder="参考文献"
              />
              <button 
                @click="removeReference(index)"
                class="px-2 py-1 bg-red-500 text-white rounded"
              >
                删除
              </button>
            </div>
            <button 
              @click="addReference"
              class="px-4 py-2 bg-blue-500 text-white rounded w-full"
            >
              添加参考文献
            </button>
            <div 
              v-if="formErrors.references" 
              class="text-red-500 text-sm mt-1"
            >
              {{ formErrors.references }}
            </div>
          </div>
        </template>
      </div>
    </div>
    
    <div class="text-center">
      <button 
        @click="exportToWord"
        class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
        :disabled="!isFormValid"
      >
        导出为Word文档
      </button>
    </div>
  </div>
</template>

<style scoped>
.form-input, .form-textarea, .form-select {
  @apply px-3 py-2 border rounded;
}

button:disabled {
  @apply opacity-50 cursor-not-allowed;
}
</style> 