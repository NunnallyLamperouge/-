<script setup>
import { ref } from 'vue'
import WorkUnit from '../../components/WorkUnit.vue'

const workItems = ref([
  { task: '完成选题', startDate: '2024-03-01', endDate: '2024-03-15' },
  { task: '完成开题报告', startDate: '2024-03-16', endDate: '2024-04-15' },
  { task: '完成中期检查', startDate: '2024-04-16', endDate: '2024-05-15' },
  { task: '完成论文初稿', startDate: '2024-05-16', endDate: '2024-06-15' }
])
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-6">工作单位组件示例</h1>
    
    <!-- 普通使用 -->
    <div class="mb-8">
      <h2 class="text-lg font-semibold mb-4">普通列表形式</h2>
      <div class="border rounded-lg p-4">
        <WorkUnit :items="workItems" />
      </div>
    </div>
    
    <!-- 在表格中使用 -->
    <div>
      <h2 class="text-lg font-semibold mb-4">表格中使用</h2>
      <table class="w-full border-collapse border">
        <thead>
          <tr>
            <th class="border p-2 bg-gray-50">工作进度</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td class="border p-4">
              <WorkUnit :items="workItems" :in-table="true" />
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 800px;
}
</style> 