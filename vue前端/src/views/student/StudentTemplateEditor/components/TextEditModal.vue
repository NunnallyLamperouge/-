<script setup>
import { ref, watch } from 'vue'
import { EditOutlined, UndoOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { postJSON } from '@/utility/request'
import { useRoute } from 'vue-router'
const route = useRoute()
const currentTemplateId = route.params.id
const props = defineProps({
  modelValue: {
    type: String,
    required: true
  },
  title: {
    type: String,
    default: '编辑文本'
  },
  rows: {
    type: Number,
    default: 20
  }
})

const emit = defineEmits(['update:modelValue'])

const visible = ref(false)
const tempValue = ref('')
const loading = ref(false)
const modifyDescription = ref('')
const contentHistory = ref([])
const canUndo = ref(false)

// 字段类型选项
const fieldTypes = [
  { label: '标题', value: 'title' },
  { label: '致谢', value: 'acknowledgment' },
  { label: '摘要', value: 'abstract' },
  { label: '正文内容', value: 'content' },
]
const selectedFieldType = ref('content')

// 风格选项
const styleTypes = [
  { label: '学术严谨', value: '学术严谨，使用专业术语' },
  { label: '简明扼要', value: '简明扼要，突出重点' },
  { label: '详细具体', value: '详细具体，论述充分' },
  { label: '创新思维', value: '体现创新思维和独特见解' },
]
const selectedStyle = ref('')

// 打开模态框时，初始化临时值和历史记录
const showModal = () => {
  tempValue.value = props.modelValue
  contentHistory.value = [props.modelValue]
  canUndo.value = false
  visible.value = true
}

// 优化文本内容
const optimizeContent = async () => {
  if (!tempValue.value.trim()) {
    message.warning('请输入需要优化的内容')
    return
  }

  loading.value = true
  try {
    postJSON('/api/rag/chatComplention', 
      {
        source: tempValue.value,
        field_type: selectedFieldType.value,
        modify_description: modifyDescription.value,
        style: selectedStyle.value,
        templateId: Number(currentTemplateId)
      },
      (msg, data) => {
        // 保存当前内容到历史记录
        contentHistory.value.push(tempValue.value)
        // 处理新的返回数据结构
        tempValue.value = data.textResponse || data
        canUndo.value = true
        message.success('内容优化成功')
        loading.value = false
      },
      (errorMessage) => {
        message.error(errorMessage || '优化失败')
        loading.value = false
      },
      (error) => {
        message.error('请求失败：' + error.message)
        loading.value = false
      }
    )
  } catch (error) {
    message.error('请求失败：' + error.message)
    loading.value = false
  }
}

// 回滚到上一个版本
const undoChanges = () => {
  if (contentHistory.value.length > 0) {
    const previousContent = contentHistory.value.pop()
    tempValue.value = previousContent
    canUndo.value = contentHistory.value.length > 0
    message.success('已回退到上一个版本')
  }
}

// 确认修改
const handleOk = () => {
  emit('update:modelValue', tempValue.value)
  visible.value = false
  // 重置表单和历史记录
  modifyDescription.value = ''
  selectedStyle.value = ''
  selectedFieldType.value = 'content'
  contentHistory.value = []
  canUndo.value = false
}

// 取消修改
const handleCancel = () => {
  visible.value = false
  // 重置表单和历史记录
  modifyDescription.value = ''
  selectedStyle.value = ''
  selectedFieldType.value = 'content'
  contentHistory.value = []
  canUndo.value = false
}
</script>

<template>
  <div class="flex items-center">
    <div class="flex-1 truncate bg-gray-50 hover:bg-gray-100 transition-colors duration-200 px-4 py-3 rounded-lg cursor-pointer" @click="showModal">
      {{ modelValue || '暂无内容' }}
    </div>
    <a-button type="link" @click="showModal" class="ml-2">
      <template #icon><EditOutlined /></template>
    </a-button>
    
    <a-modal
      :title="title"
      v-model:visible="visible"
      @ok="handleOk"
      @cancel="handleCancel"
      width="800px"
      :okButtonProps="{ loading }"
      class="text-edit-modal"
      :bodyStyle="{ padding: '24px' }"
    >
      <div class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- 字段类型选择 -->
          <div class="field-type-section">
            <div class="text-base font-medium mb-3 text-gray-700">文本类型</div>
            <a-radio-group v-model:value="selectedFieldType" class="w-full flex flex-wrap gap-2">
              <a-radio-button 
                v-for="type in fieldTypes" 
                :key="type.value" 
                :value="type.value"
                class="flex-grow text-center"
              >
                {{ type.label }}
              </a-radio-button>
            </a-radio-group>
          </div>

          <!-- 风格选择 -->
          <div class="style-section">
            <div class="text-base font-medium mb-3 text-gray-700">优化风格</div>
            <a-radio-group v-model:value="selectedStyle" class="w-full flex flex-wrap gap-2">
              <a-radio-button 
                v-for="style in styleTypes" 
                :key="style.value" 
                :value="style.value"
                class="flex-grow text-center"
              >
                {{ style.label }}
              </a-radio-button>
            </a-radio-group>
          </div>
        </div>

        <!-- 修改要求输入 -->
        <div class="modify-description-section">
          <div class="text-base font-medium mb-3 text-gray-700">具体修改要求</div>
          <a-input
            v-model:value="modifyDescription"
            placeholder="请输入具体的修改要求（选填）"
            allow-clear
            class="hover:border-blue-400 focus:border-blue-500"
          />
        </div>

        <!-- 文本内容编辑区 -->
        <div class="content-section">
          <div class="text-base font-medium mb-3 text-gray-700 flex justify-between items-center">
            <span>编辑文本内容</span>
            <div class="flex items-center gap-3">
              <a-tooltip v-if="canUndo" title="回退到上一版本">
                <a-button 
                  type="default"
                  @click="undoChanges"
                  class="undo-btn"
                >
                  <template #icon><UndoOutlined /></template>
                  撤销修改
                </a-button>
              </a-tooltip>
              <a-button 
                type="primary" 
                :loading="loading"
                @click="optimizeContent"
                class="optimize-btn"
              >
                一键优化内容
              </a-button>
            </div>
          </div>
          <a-textarea
            v-model:value="tempValue"
            :rows="rows"
            :auto-size="{ minRows: rows, maxRows: 40 }"
            class="w-full hover:border-blue-400 focus:border-blue-500"
          />
        </div>
      </div>
    </a-modal>
  </div>
</template>

<style scoped>
:deep(.ant-modal-body) {
  max-height: 80vh;
  overflow-y: auto;
}

.space-y-6 > * + * {
  margin-top: 1.5rem;
}

:deep(.ant-radio-button-wrapper) {
  transition: all 0.3s;
}

:deep(.ant-radio-button-wrapper-checked) {
  background-color: #1890ff !important;
  color: white !important;
  border-color: #1890ff !important;
}

:deep(.ant-radio-button-wrapper:hover) {
  color: #1890ff;
}

.optimize-btn {
  background: linear-gradient(to right, #1890ff, #36cfc9);
  border: none;
  transition: all 0.3s;
}

.optimize-btn:hover {
  opacity: 0.9;
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(24, 144, 255, 0.2);
}

:deep(.ant-input:hover),
:deep(.ant-input:focus) {
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.text-edit-modal :deep(.ant-modal-content) {
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.undo-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  transition: all 0.3s;
  border-color: #d9d9d9;
}

.undo-btn:hover {
  color: #1890ff;
  border-color: #1890ff;
  background: rgba(24, 144, 255, 0.1);
}

.undo-btn :deep(.anticon) {
  font-size: 14px;
}
</style> 