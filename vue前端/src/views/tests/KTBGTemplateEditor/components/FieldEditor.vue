<script setup>
import { ref, watch } from 'vue';

const props = defineProps({
  field: {
    type: Object,
    required: true
  }
});

// 字段类型选项
const fieldTypes = [
  { value: 'text', label: '文本' },
  { value: 'textarea', label: '多行文本' },
  { value: 'number', label: '数字' },
  { value: 'date', label: '日期' },
  { value: 'radio', label: '单选框' },
  { value: 'checkbox', label: '复选框' },
  { value: 'select', label: '下拉选择' }
];

// 选项列表（用于单选、复选、下拉类型）
const options = ref(props.field.options || []);

// 添加选项
const addOption = () => {
  options.value.push('新选项');
  updateOptions();
};

// 删除选项
const removeOption = (index) => {
  options.value.splice(index, 1);
  updateOptions();
};

// 更新字段的选项
const updateOptions = () => {
  props.field.options = [...options.value];
};

// 监听选项变化
watch(options, () => {
  updateOptions();
}, { deep: true });

// 检查字段类型是否需要选项
const needsOptions = (type) => {
  return ['radio', 'checkbox', 'select'].includes(type);
};
</script>

<template>
  <div class="field-editor">
    <div class="mb-4">
      <label class="block mb-1 font-medium">字段标签</label>
      <input 
        v-model="field.label" 
        class="form-input w-full" 
        :disabled="field.disabled"
      />
    </div>
    
    <div class="mb-4">
      <label class="block mb-1 font-medium">字段类型</label>
      <select 
        v-model="field.type" 
        class="form-select w-full"
        :disabled="field.disabled"
      >
        <option 
          v-for="type in fieldTypes" 
          :key="type.value" 
          :value="type.value"
        >
          {{ type.label }}
        </option>
      </select>
    </div>
    
    <div class="mb-4">
      <label class="block mb-1 font-medium">是否必填</label>
      <div class="flex items-center">
        <input 
          type="checkbox" 
          v-model="field.required" 
          class="mr-2"
          :disabled="field.disabled"
        />
        <span>必填字段</span>
      </div>
    </div>
    
    <template v-if="field.type === 'textarea'">
      <div class="mb-4">
        <label class="block mb-1 font-medium">最小字数</label>
        <input 
          type="number" 
          v-model.number="field.minLength" 
          class="form-input w-full"
          min="0"
        />
      </div>
      
      <div class="mb-4">
        <label class="block mb-1 font-medium">最大字数</label>
        <input 
          type="number" 
          v-model.number="field.maxLength" 
          class="form-input w-full"
          min="0"
        />
      </div>
    </template>
    
    <template v-if="needsOptions(field.type)">
      <div class="mb-4">
        <label class="block mb-1 font-medium">选项列表</label>
        <div 
          v-for="(option, index) in options" 
          :key="index"
          class="flex mb-2"
        >
          <input 
            v-model="options[index]" 
            class="form-input flex-1 mr-2"
          />
          <button 
            @click="removeOption(index)"
            class="px-2 py-1 bg-red-500 text-white rounded"
          >
            删除
          </button>
        </div>
        <button 
          @click="addOption"
          class="px-4 py-2 bg-blue-500 text-white rounded w-full"
        >
          添加选项
        </button>
      </div>
    </template>
    
    <div class="mb-4">
      <label class="block mb-1 font-medium">默认值</label>
      <input 
        v-model="field.defaultValue" 
        class="form-input w-full"
        :disabled="field.disabled"
      />
    </div>
  </div>
</template>

<style scoped>
.form-input, .form-select {
  @apply px-3 py-2 border rounded;
}
</style> 