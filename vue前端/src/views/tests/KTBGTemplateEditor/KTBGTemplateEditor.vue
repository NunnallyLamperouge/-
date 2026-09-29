<script setup>
import { ref, reactive, onMounted, watch } from 'vue';
import SectionEditor from './components/SectionEditor.vue';
import TemplatePreview from './components/TemplatePreview.vue';

// 当前活动的标签页
const activeTab = ref('basic'); // 'basic', 'sections', 'preview'

// 默认模板结构
const defaultTemplate = {
  id: 'ktbg',
  name: '开题报告模板',
  sections: [
    {
      id: 'cover',
      title: '封面信息',
      type: 'fields',
      fields: [
        {
          id: 'title',
          label: '论文题目',
          type: 'text',
          required: true,
          defaultValue: ''
        },
        {
          id: 'studentName',
          label: '学生姓名',
          type: 'text',
          required: true,
          defaultValue: ''
        },
        {
          id: 'studentId',
          label: '学号',
          type: 'text',
          required: true,
          defaultValue: ''
        },
        {
          id: 'major',
          label: '专业',
          type: 'text',
          required: true,
          defaultValue: ''
        },
        {
          id: 'supervisor',
          label: '指导教师',
          type: 'text',
          required: true,
          defaultValue: ''
        },
        {
          id: 'date',
          label: '日期',
          type: 'date',
          required: true,
          defaultValue: ''
        }
      ]
    },
    {
      id: 'background',
      title: '研究背景与意义',
      type: 'fields',
      fields: [
        {
          id: 'background',
          label: '研究背景与意义',
          type: 'textarea',
          required: true,
          minLength: 200,
          maxLength: 2000,
          defaultValue: ''
        }
      ]
    },
    {
      id: 'literature',
      title: '文献综述',
      type: 'fields',
      fields: [
        {
          id: 'literature',
          label: '文献综述',
          type: 'textarea',
          required: true,
          minLength: 300,
          maxLength: 3000,
          defaultValue: ''
        }
      ]
    },
    {
      id: 'methods',
      title: '研究方法',
      type: 'fields',
      fields: [
        {
          id: 'methods',
          label: '研究方法',
          type: 'textarea',
          required: true,
          minLength: 200,
          maxLength: 2000,
          defaultValue: ''
        }
      ]
    },
    {
      id: 'schedule',
      title: '工作计划',
      type: 'schedule',
      defaultItems: [
        {
          task: '文献调研',
          startDate: '',
          endDate: ''
        },
        {
          task: '方案设计',
          startDate: '',
          endDate: ''
        }
      ]
    },
    {
      id: 'references',
      title: '参考文献',
      type: 'references',
      minItems: 5
    }
  ]
};

// 当前模板状态
const currentTemplate = reactive({
  id: '',
  name: '',
  sections: []
});
const savedTemplates = ref([]);
const newTemplateName = ref('');
const currentTemplateIndex = ref(-1);

// 从localStorage加载已保存的模板
onMounted(() => {
  loadSavedTemplates();
});

// 监听模板变化，自动保存到localStorage
watch(currentTemplate, () => {
  if (currentTemplateIndex.value >= 0) {
    // 深拷贝当前模板，确保保存的是完整的数据
    savedTemplates.value[currentTemplateIndex.value] = JSON.parse(JSON.stringify(currentTemplate));
    saveTemplatesToLocalStorage();
  }
}, { deep: true });

// 加载已保存的模板
const loadSavedTemplates = () => {
  const savedData = localStorage.getItem('KTBGTemplates');
  if (savedData) {
    try {
      savedTemplates.value = JSON.parse(savedData);
      if (savedTemplates.value.length > 0) {
        // 加载第一个模板
        loadTemplate(0);
      } else {
        // 如果没有保存的模板，创建一个默认模板
        createDefaultTemplate();
      }
    } catch (e) {
      console.error('解析保存的模板失败:', e);
      createDefaultTemplate();
    }
  } else {
    // 如果没有保存的模板，创建一个默认模板
    createDefaultTemplate();
  }
};

// 创建默认模板
const createDefaultTemplate = () => {
  const template = JSON.parse(JSON.stringify(defaultTemplate));
  savedTemplates.value.push(template);
  currentTemplateIndex.value = 0;
  
  // 更新当前模板
  Object.assign(currentTemplate, template);
  
  saveTemplatesToLocalStorage();
};

// 保存模板到localStorage
const saveTemplatesToLocalStorage = () => {
  localStorage.setItem('KTBGTemplates', JSON.stringify(savedTemplates.value));
};

// 创建新模板
const createNewTemplate = () => {
  if (!newTemplateName.value.trim()) {
    alert('请输入模板名称');
    return;
  }
  
  const newTemplate = JSON.parse(JSON.stringify(defaultTemplate));
  newTemplate.name = newTemplateName.value;
  newTemplate.id = `ktbg_${Date.now()}`;
  savedTemplates.value.push(newTemplate);
  currentTemplateIndex.value = savedTemplates.value.length - 1;
  
  // 更新当前模板
  Object.assign(currentTemplate, newTemplate);
  
  saveTemplatesToLocalStorage();
  newTemplateName.value = '';
};

// 加载已有模板
const loadTemplate = (index) => {
  currentTemplateIndex.value = index;
  const template = savedTemplates.value[index];
  
  // 更新当前模板 - 使用深拷贝确保数据独立
  Object.assign(currentTemplate, JSON.parse(JSON.stringify(template)));
};

// 删除模板
const deleteTemplate = (index) => {
  if (confirm('确定要删除此模板吗？')) {
    savedTemplates.value.splice(index, 1);
    saveTemplatesToLocalStorage();
    
    if (currentTemplateIndex.value === index) {
      if (savedTemplates.value.length > 0) {
        currentTemplateIndex.value = 0;
        loadTemplate(0);
      } else {
        currentTemplateIndex.value = -1;
        // 重置为默认模板但不保存
        Object.assign(currentTemplate, JSON.parse(JSON.stringify(defaultTemplate)));
      }
    } else if (currentTemplateIndex.value > index) {
      currentTemplateIndex.value--;
    }
  }
};

// 添加新部分
const addSection = () => {
  const newSection = {
    id: `section_${Date.now()}`,
    title: '新部分',
    type: 'fields',
    fields: []
  };
  
  currentTemplate.sections.push(newSection);
};

// 更新部分
const updateSection = (index, section) => {
  currentTemplate.sections[index] = { ...section };
};

// 删除部分
const deleteSection = (index) => {
  if (confirm('确定要删除此部分吗？')) {
    currentTemplate.sections.splice(index, 1);
  }
};

// 移动部分
const moveSection = (index, direction) => {
  if (direction === 'up' && index > 0) {
    const temp = currentTemplate.sections[index];
    currentTemplate.sections[index] = currentTemplate.sections[index - 1];
    currentTemplate.sections[index - 1] = temp;
  } else if (direction === 'down' && index < currentTemplate.sections.length - 1) {
    const temp = currentTemplate.sections[index];
    currentTemplate.sections[index] = currentTemplate.sections[index + 1];
    currentTemplate.sections[index + 1] = temp;
  }
};

// 重置模板
const resetTemplate = () => {
  if (confirm('确定要重置模板吗？这将丢失所有当前的更改。')) {
    const template = JSON.parse(JSON.stringify(defaultTemplate));
    // 保留原模板的ID和名称
    template.id = currentTemplate.id;
    template.name = currentTemplate.name;
    
    // 更新当前模板
    Object.assign(currentTemplate, template);
    
    if (currentTemplateIndex.value >= 0) {
      savedTemplates.value[currentTemplateIndex.value] = JSON.parse(JSON.stringify(currentTemplate));
      saveTemplatesToLocalStorage();
    }
  }
};

// 手动保存当前模板
const saveCurrentTemplate = () => {
  if (currentTemplateIndex.value >= 0) {
    savedTemplates.value[currentTemplateIndex.value] = JSON.parse(JSON.stringify(currentTemplate));
    saveTemplatesToLocalStorage();
    alert('模板已保存');
  } else {
    alert('请先创建或加载一个模板');
  }
};
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-6">开题报告模板编辑器</h1>
    
    <!-- 模板选择和管理 -->
    <div class="bg-gray-100 p-4 rounded mb-4">
      <h2 class="text-xl mb-2">模板管理</h2>
      <div class="flex gap-4 mb-4">
        <div class="flex-1">
          <input 
            v-model="newTemplateName" 
            placeholder="新模板名称" 
            class="w-full px-3 py-2 border rounded" 
          />
        </div>
        <button 
          @click="createNewTemplate" 
          class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          创建模板
        </button>
        <button 
          @click="saveCurrentTemplate" 
          class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
          :disabled="currentTemplateIndex < 0"
        >
          保存当前模板
        </button>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div 
          v-for="(template, index) in savedTemplates" 
          :key="index" 
          class="border p-3 rounded cursor-pointer"
          :class="{'bg-blue-100': index === currentTemplateIndex}"
        >
          <div class="flex justify-between items-center">
            <span>{{ template.name }}</span>
            <div>
              <button 
                @click="loadTemplate(index)" 
                class="px-2 py-1 bg-green-500 text-white rounded mr-2 hover:bg-green-600"
              >
                加载
              </button>
              <button 
                @click="deleteTemplate(index)" 
                class="px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600"
              >
                删除
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 编辑部分的标签页 -->
    <div class="mb-4">
      <div class="flex border-b">
        <button 
          @click="activeTab = 'basic'" 
          class="px-4 py-2" 
          :class="{'border-b-2 border-blue-500 text-blue-500': activeTab === 'basic'}"
        >
          基本信息
        </button>
        <button 
          @click="activeTab = 'sections'" 
          class="px-4 py-2" 
          :class="{'border-b-2 border-blue-500 text-blue-500': activeTab === 'sections'}"
        >
          模板部分
        </button>
        <button 
          @click="activeTab = 'preview'" 
          class="px-4 py-2" 
          :class="{'border-b-2 border-blue-500 text-blue-500': activeTab === 'preview'}"
        >
          预览填写
        </button>
      </div>
    </div>
    
    <!-- 基本信息编辑区 -->
    <div v-if="activeTab === 'basic'" class="mb-6">
      <div class="bg-white p-4 rounded shadow-sm">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- 模板ID -->
          <div class="form-group">
            <label class="block mb-2 font-medium">模板ID</label>
            <input 
              v-model="currentTemplate.id"
              type="text"
              class="form-input w-full"
              disabled
            />
            <p class="text-gray-500 text-sm mt-1">模板ID自动生成，不可修改</p>
          </div>
          
          <!-- 模板名称 -->
          <div class="form-group">
            <label class="block mb-2 font-medium">模板名称</label>
            <input 
              v-model="currentTemplate.name"
              type="text"
              class="form-input w-full"
            />
          </div>
        </div>
        
        <div class="mt-4">
          <button 
            @click="resetTemplate"
            class="px-4 py-2 bg-yellow-500 text-white rounded hover:bg-yellow-600"
          >
            重置为默认模板
          </button>
        </div>
      </div>
    </div>
    
    <!-- 模板部分编辑区 -->
    <div v-if="activeTab === 'sections'" class="mb-6">
      <div class="mb-4 flex justify-between items-center">
        <h2 class="text-xl font-semibold">模板部分</h2>
        <button 
          @click="addSection"
          class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          添加部分
        </button>
      </div>
      
      <div v-if="currentTemplate.sections && currentTemplate.sections.length > 0">
        <div 
          v-for="(section, index) in currentTemplate.sections" 
          :key="section.id"
          class="mb-6 relative"
        >
          <div class="absolute right-0 -top-4 flex space-x-2">
            <button 
              @click="moveSection(index, 'up')"
              class="px-2 py-1 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
              :disabled="index === 0"
            >
              上移
            </button>
            <button 
              @click="moveSection(index, 'down')"
              class="px-2 py-1 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
              :disabled="index === currentTemplate.sections.length - 1"
            >
              下移
            </button>
          </div>
          
          <SectionEditor 
            :section="section"
            @update:section="updateSection(index, $event)"
            @delete="deleteSection(index)"
          />
        </div>
      </div>
      
      <div v-else class="text-center py-8 text-gray-500 border rounded">
        暂无模板部分，点击"添加部分"按钮创建
      </div>
    </div>
    
    <!-- 预览内容 -->
    <div v-if="activeTab === 'preview'" class="mb-6">
      <TemplatePreview :template="currentTemplate" />
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
