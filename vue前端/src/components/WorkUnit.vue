<script setup>
import { defineProps } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    required: true,
    default: () => []
  },
  // 是否在表格中使用
  inTable: {
    type: Boolean,
    default: false
  }
})
</script>

<template>
  <div :class="['work-unit', { 'in-table': inTable }]">
    <div v-for="(item, index) in items" :key="index" class="work-item">
      <div class="flex justify-between items-center py-2">
        <span class="task-info">
          <span class="task-index">（{{ index + 1 }}）</span>
          <span class="task-name">{{ item.task }}</span>
        </span>
        <span class="date-range">{{ item.startDate }} — {{ item.endDate }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.work-unit {
  width: 100%;
}

.work-item {
  font-size: 14px;
  color: #333;
}

.in-table .work-item {
  border-bottom: 1px solid #eee;
}

.in-table .work-item:last-child {
  border-bottom: none;
}

.task-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.task-index {
  color: #666;
}

.task-name {
  font-weight: 500;
}

.date-range {
  color: #666;
  font-size: 13px;
}
</style> 