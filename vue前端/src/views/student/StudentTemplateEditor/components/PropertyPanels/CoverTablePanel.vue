<script setup>
import { ref } from 'vue'
import { fontFamilies, ensureComponentSettings } from './common/CommonData'

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

// 更新表格字段
const updateFields = () => {
  if (!props.component || !props.component.fields) return
  updateProperty('fields', [...props.component.fields])
}

// 添加表格字段
const addField = () => {
  if (!props.component.fields) {
    props.component.fields = []
  }
  
  const newField = {
    id: Date.now(),
    label: '新字段',
    value: '',
    width: 30
  }
  
  props.component.fields.push(newField)
  updateFields()
}

// 删除表格字段
const removeField = (fieldId) => {
  const index = props.component.fields.findIndex(f => f.id === fieldId)
  if (index !== -1) {
    props.component.fields.splice(index, 1)
    updateFields()
  }
}

// 添加拖拽排序相关函数
const draggedField = ref(null)
const dragOverIndex = ref(null)

// 开始拖拽字段
const startDragField = (field) => {
  draggedField.value = field
}

// 拖拽结束
const endDragField = () => {
  draggedField.value = null
  dragOverIndex.value = null
}

// 拖拽经过某个位置
const dragOverField = (index) => {
  dragOverIndex.value = index
}

// 放置字段
const dropField = (index) => {
  if (!draggedField.value) return
  
  const fromIndex = props.component.fields.findIndex(f => f.id === draggedField.value.id)
  if (fromIndex === -1) return
  
  // 移除原位置的字段
  const [movedField] = props.component.fields.splice(fromIndex, 1)
  
  // 计算目标位置（考虑到移除后索引的变化）
  let targetIndex = index
  if (fromIndex < index) {
    targetIndex--
  }
  
  // 插入到新位置
  props.component.fields.splice(targetIndex, 0, movedField)
  
  // 更新字段
  updateFields()
  
  // 重置拖拽状态
  endDragField()
}
</script>

<template>
  <div>
    <!-- 字段设置 -->
    <div class="mb-4">
      <div class="flex justify-between items-center mb-2">
        <label class="block">表格字段</label>
        <button
          @click="addField"
          class="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 flex items-center gap-1"
        >
          <span>+</span> 添加字段
        </button>
      </div>
      
      <div v-for="field in component.fields" :key="field.id" class="mb-4 p-4 border rounded"
        :class="{ 'border-blue-300 bg-blue-50': draggedField && draggedField.id === field.id }"
        draggable="true"
        @dragstart="startDragField(field)"
        @dragend="endDragField"
        @dragover.prevent="dragOverField(index)"
        @drop.prevent="dropField(index)"
      >
        <div class="flex gap-2 mb-2">
          <div class="mr-2 cursor-move text-gray-400 hover:text-gray-600">⋮⋮</div>
          <input
            v-model="field.label"
            type="text"
            class="flex-1 px-3 py-2 border rounded"
            placeholder="字段标签"
            @input="updateFields"
          />
          <input
            v-model.number="field.width"
            type="number"
            class="w-20 px-3 py-2 border rounded"
            min="10"
            max="50"
            placeholder="宽度%"
            @input="updateFields"
          />
          <button
            @click="removeField(field.id)"
            class="px-3 py-1 bg-red-500 text-white rounded"
          >
            删除
          </button>
        </div>
        
        <div>
          <label class="block mb-1 text-sm">值</label>
          <input
            v-model="field.value"
            type="text"
            class="w-full px-3 py-2 border rounded"
            placeholder="字段值"
            @input="updateFields"
          />
        </div>
      </div>
      
      <!-- 拖拽提示 -->
      <div class="mt-2 text-xs text-gray-500">
        提示：可拖拽字段调整顺序
      </div>
    </div>

    <!-- 样式设置 -->
    <div class="mb-4">
      <label class="block mb-2">样式设置</label>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block mb-1">字体</label>
          <select
            v-model="ensureComponentSettings(component).fontFamily"
            class="w-full px-3 py-2 border rounded"
            @change="updateProperty('settings', {...component.settings})"
          >
            <option v-for="font in fontFamilies" :key="font.value" :value="font.value">
              {{ font.label }}
            </option>
          </select>
        </div>
        <div>
          <label class="block mb-1">字号</label>
          <input
            v-model.number="ensureComponentSettings(component).fontSize"
            type="number"
            class="w-full px-3 py-2 border rounded"
            min="12"
            max="72"
            @input="updateProperty('settings', {...component.settings})"
          />
        </div>
      </div>
    </div>
    
    <!-- 表格边框设置 -->
    <div class="mb-4">
      <label class="block mb-2">边框样式</label>
      <div class="grid grid-cols-2 gap-4">
        <button
          :class="[
            'p-3 border rounded flex flex-col items-center',
            component.borderStyle === 'none' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('borderStyle', 'none')"
        >
          <div class="text-2xl mb-1">◻️</div>
          <div class="text-sm">无边框</div>
        </button>
        <button
          :class="[
            'p-3 border rounded flex flex-col items-center',
            component.borderStyle === 'outer' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('borderStyle', 'outer')"
        >
          <div class="text-2xl mb-1">🔲</div>
          <div class="text-sm">外边框</div>
        </button>
        <button
          :class="[
            'p-3 border rounded flex flex-col items-center',
            component.borderStyle === 'full' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('borderStyle', 'full')"
        >
          <div class="text-2xl mb-1">📊</div>
          <div class="text-sm">全边框</div>
        </button>
        <button
          :class="[
            'p-3 border rounded flex flex-col items-center',
            component.borderStyle === 'horizontal' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('borderStyle', 'horizontal')"
        >
          <div class="text-2xl mb-1">☰</div>
          <div class="text-sm">横线</div>
        </button>
      </div>
    </div>
    
    <!-- 表格对齐方式 -->
    <div class="mb-4">
      <label class="block mb-2">表格对齐方式</label>
      <div class="grid grid-cols-3 gap-2">
        <button
          :class="[
            'p-2 border rounded flex items-center justify-center',
            component.alignment === 'left' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('alignment', 'left')"
        >
          左对齐
        </button>
        <button
          :class="[
            'p-2 border rounded flex items-center justify-center',
            component.alignment === 'center' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('alignment', 'center')"
        >
          居中
        </button>
        <button
          :class="[
            'p-2 border rounded flex items-center justify-center',
            component.alignment === 'right' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('alignment', 'right')"
        >
          右对齐
        </button>
      </div>
    </div>
    
    <!-- 表格背景色 -->
    <div class="mb-4">
      <label class="block mb-2">表格背景色</label>
      <div class="flex items-center gap-2">
        <input
          v-model="component.backgroundColor"
          type="color"
          class="w-12 h-10 border rounded"
          @input="updateProperty('backgroundColor', $event.target.value)"
        />
        <input
          v-model="component.backgroundColor"
          type="text"
          class="flex-1 px-3 py-2 border rounded"
          placeholder="#FFFFFF"
          @input="updateProperty('backgroundColor', $event.target.value)"
        />
      </div>
    </div>
  </div>
</template> 