<script setup>
import { ref, computed } from 'vue';
import FieldEditor from './FieldEditor.vue';

const props = defineProps({
  section: {
    type: Object,
    required: true
  }
});

const emit = defineEmits(['update:section', 'delete']);

// 部分类型
const sectionTypes = [
  { value: 'fields', label: '普通字段' },
  { value: 'schedule', label: '工作计划' },
  { value: 'references', label: '参考文献' }
];

// 编辑中的部分
const editingSection = ref({ ...props.section });

// 监听编辑中的部分变化
const updateSection = () => {
  emit('update:section', { ...editingSection.value });
};

// 添加字段
const addField = () => {
  if (!editingSection.value.fields) {
    editingSection.value.fields = [];
  }
  
  const newField = {
    id: `field_${Date.now()}`,
    label: '新字段',
    type: 'text',
    required: false,
    defaultValue: ''
  };
  
  editingSection.value.fields.push(newField);
  updateSection();
};

// 更新字段
const updateField = (index, field) => {
  editingSection.value.fields[index] = { ...field };
  updateSection();
};

// 删除字段
const deleteField = (index) => {
  editingSection.value.fields.splice(index, 1);
  updateSection();
};

// 移动字段
const moveField = (index, direction) => {
  if (direction === 'up' && index > 0) {
    const temp = editingSection.value.fields[index];
    editingSection.value.fields[index] = editingSection.value.fields[index - 1];
    editingSection.value.fields[index - 1] = temp;
    updateSection();
  } else if (direction === 'down' && index < editingSection.value.fields.length - 1) {
    const temp = editingSection.value.fields[index];
    editingSection.value.fields[index] = editingSection.value.fields[index + 1];
    editingSection.value.fields[index + 1] = temp;
    updateSection();
  }
};

// 删除部分
const deleteSection = () => {
  emit('delete');
};

// 是否显示字段列表
const showFieldsList = computed(() => {
  return editingSection.value.type === 'fields' || !editingSection.value.type;
});

// 是否显示工作计划设置
const showScheduleSettings = computed(() => {
  return editingSection.value.type === 'schedule';
});

// 是否显示参考文献设置
const showReferencesSettings = computed(() => {
  return editingSection.value.type === 'references';
});

// 更新部分类型
const updateSectionType = (type) => {
  editingSection.value.type = type;
  
  // 根据类型初始化相关属性
  if (type === 'fields') {
    if (!editingSection.value.fields) {
      editingSection.value.fields = [];
    }
  } else if (type === 'schedule') {
    editingSection.value.defaultItems = editingSection.value.defaultItems || [];
  } else if (type === 'references') {
    editingSection.value.minItems = editingSection.value.minItems || 0;
  }
  
  updateSection();
};

// 添加默认工作计划项
const addDefaultScheduleItem = () => {
  if (!editingSection.value.defaultItems) {
    editingSection.value.defaultItems = [];
  }
  
  editingSection.value.defaultItems.push({
    task: '',
    startDate: '',
    endDate: ''
  });
  
  updateSection();
};

// 删除默认工作计划项
const removeDefaultScheduleItem = (index) => {
  editingSection.value.defaultItems.splice(index, 1);
  updateSection();
};
</script>

<template>
  <div class="section-editor border rounded-lg p-4 mb-4 bg-white shadow-sm">
    <div class="flex justify-between items-center mb-4">
      <h3 class="text-lg font-semibold">{{ editingSection.title || '未命名部分' }}</h3>
      <button 
        @click="deleteSection"
        class="px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600"
      >
        删除部分
      </button>
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
      <!-- 部分ID -->
      <div class="form-group">
        <label class="block mb-2 font-medium">部分ID</label>
        <input 
          v-model="editingSection.id"
          type="text"
          class="form-input w-full"
          @change="updateSection"
        />
      </div>
      
      <!-- 部分标题 -->
      <div class="form-group">
        <label class="block mb-2 font-medium">部分标题</label>
        <input 
          v-model="editingSection.title"
          type="text"
          class="form-input w-full"
          @change="updateSection"
        />
      </div>
      
      <!-- 部分类型 -->
      <div class="form-group col-span-1 md:col-span-2">
        <label class="block mb-2 font-medium">部分类型</label>
        <div class="flex space-x-4">
          <div 
            v-for="type in sectionTypes" 
            :key="type.value"
            class="flex items-center"
          >
            <input 
              type="radio" 
              :id="`type_${type.value}`"
              :value="type.value"
              v-model="editingSection.type"
              @change="updateSectionType(type.value)"
              class="mr-2"
            />
            <label :for="`type_${type.value}`">{{ type.label }}</label>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 字段列表 -->
    <div v-if="showFieldsList" class="mb-4">
      <div class="flex justify-between items-center mb-2">
        <h4 class="font-medium">字段列表</h4>
        <button 
          @click="addField"
          class="px-2 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          添加字段
        </button>
      </div>
      
      <div v-if="editingSection.fields && editingSection.fields.length > 0">
        <div 
          v-for="(field, index) in editingSection.fields" 
          :key="field.id"
          class="border rounded p-3 mb-3"
        >
          <div class="flex justify-between items-center mb-2">
            <h5 class="font-medium">{{ field.label || '未命名字段' }}</h5>
            <div class="flex space-x-2">
              <button 
                @click="moveField(index, 'up')"
                class="px-2 py-1 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
                :disabled="index === 0"
              >
                ↑
              </button>
              <button 
                @click="moveField(index, 'down')"
                class="px-2 py-1 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
                :disabled="index === editingSection.fields.length - 1"
              >
                ↓
              </button>
              <button 
                @click="deleteField(index)"
                class="px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600"
              >
                删除
              </button>
            </div>
          </div>
          
          <FieldEditor 
            :field="field"
            @update:field="updateField(index, $event)"
          />
        </div>
      </div>
      
      <div v-else class="text-center py-4 text-gray-500">
        暂无字段，点击"添加字段"按钮创建
      </div>
    </div>
    
    <!-- 工作计划设置 -->
    <div v-if="showScheduleSettings" class="mb-4">
      <h4 class="font-medium mb-2">工作计划设置</h4>
      
      <div class="mb-4">
        <label class="block mb-2">默认工作计划项</label>
        
        <div 
          v-for="(item, index) in editingSection.defaultItems" 
          :key="index"
          class="grid grid-cols-3 gap-2 mb-2"
        >
          <input 
            v-model="item.task" 
            class="form-input" 
            placeholder="任务描述"
            @change="updateSection"
          />
          <input 
            type="date" 
            v-model="item.startDate" 
            class="form-input"
            @change="updateSection"
          />
          <div class="flex">
            <input 
              type="date" 
              v-model="item.endDate" 
              class="form-input flex-1 mr-2"
              @change="updateSection"
            />
            <button 
              @click="removeDefaultScheduleItem(index)"
              class="px-2 py-1 bg-red-500 text-white rounded"
            >
              删除
            </button>
          </div>
        </div>
        
        <button 
          @click="addDefaultScheduleItem"
          class="px-4 py-2 bg-blue-500 text-white rounded w-full"
        >
          添加默认任务
        </button>
      </div>
    </div>
    
    <!-- 参考文献设置 -->
    <div v-if="showReferencesSettings" class="mb-4">
      <h4 class="font-medium mb-2">参考文献设置</h4>
      
      <div class="form-group">
        <label class="block mb-2">最少参考文献数量</label>
        <input 
          v-model.number="editingSection.minItems"
          type="number"
          min="0"
          class="form-input w-full"
          @change="updateSection"
        />
      </div>
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