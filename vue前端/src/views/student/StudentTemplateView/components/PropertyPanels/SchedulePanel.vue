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

// 更新工作项目
const updateItems = () => {
  updateProperty('items', [...props.component.items])
}

// 添加工作项目
const addItem = () => {
  if (!props.component.items) {
    props.component.items = []
  }
  
  props.component.items.push({
    task: '',
    startDate: '',
    endDate: ''
  })
  updateItems()
}

// 删除工作项目
const removeItem = (index) => {
  if (props.component.items.length > 1) {
    props.component.items.splice(index, 1)
    updateItems()
  } else {
    alert('至少保留一个工作项目')
  }
}
</script>

<template>
  <div>
    <div class="mb-4">
      <label class="block mb-2">工作项目</label>
      <div v-for="(item, index) in component.items" :key="index" class="mb-2 p-3 border rounded">
        <div class="grid grid-cols-1 gap-2">
          <input
            v-model="item.task"
            type="text"
            class="w-full px-3 py-2 border rounded"
            placeholder="任务名称"
            @input="updateItems"
          />
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-sm mb-1">开始日期</label>
              <input
                v-model="item.startDate"
                type="date"
                class="w-full px-3 py-2 border rounded"
                @input="updateItems"
              />
            </div>
            <div>
              <label class="block text-sm mb-1">结束日期</label>
              <input
                v-model="item.endDate"
                type="date"
                class="w-full px-3 py-2 border rounded"
                @input="updateItems"
              />
            </div>
          </div>
          <button
            @click="removeItem(index)"
            class="px-3 py-1 bg-red-500 text-white rounded"
            :disabled="component.items.length <= 1"
          >
            删除
          </button>
        </div>
      </div>
      <button
        @click="addItem"
        class="mt-2 px-3 py-2 bg-blue-500 text-white rounded w-full"
      >
        添加工作项
      </button>
    </div>

    <div class="grid grid-cols-2 gap-4 mb-4">
      <div>
        <label class="block mb-2">字体</label>
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
        <label class="block mb-2">字号</label>
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
    
    <div class="mb-4">
      <label class="block mb-2">行高</label>
      <input
        v-model.number="ensureComponentSettings(component).lineHeight"
        type="number"
        class="w-full px-3 py-2 border rounded"
        min="1"
        max="3"
        step="0.1"
        @input="updateProperty('settings', {...component.settings})"
      />
    </div>
    
    <div class="mb-4">
      <label class="block mb-2">标题</label>
      <input
        v-model="component.title"
        type="text"
        class="w-full px-3 py-2 border rounded"
        placeholder="进度表标题"
        @input="updateProperty('title', $event.target.value)"
      />
    </div>
    
    <div class="mb-4">
      <label class="block mb-2">显示方式</label>
      <div class="grid grid-cols-2 gap-4">
        <button
          :class="[
            'p-3 border rounded flex flex-col items-center',
            component.displayMode === 'table' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('displayMode', 'table')"
        >
          <div class="text-2xl mb-1">📋</div>
          <div class="text-sm">表格模式</div>
        </button>
        <button
          :class="[
            'p-3 border rounded flex flex-col items-center',
            component.displayMode === 'gantt' ? 'bg-blue-100 border-blue-500' : 'bg-white'
          ]"
          @click="updateProperty('displayMode', 'gantt')"
        >
          <div class="text-2xl mb-1">📊</div>
          <div class="text-sm">甘特图模式</div>
        </button>
      </div>
    </div>
  </div>
</template> 