<template>
  <div class="document-editor">
    <!-- 顶部工具栏 -->
    <div class="bg-white border-b shadow-sm py-2 px-4 flex items-center justify-between sticky top-0 z-10">
      <div class="flex items-center space-x-4">
        <h1 class="text-xl font-semibold text-gray-800">文档编辑器</h1>
        <div class="h-6 w-px bg-gray-300"></div>
        <button @click="saveDocument" class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition flex items-center">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
          </svg>
          保存文档
        </button>
        <button @click="exportDocument" class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition flex items-center">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          导出Word
        </button>
      </div>
      <div class="flex items-center space-x-4">
        <span class="text-sm text-gray-500">上次保存: {{ lastSaved }}</span>
        <button @click="toggleOutline" class="px-3 py-1 text-gray-600 hover:bg-gray-100 rounded-md transition flex items-center">
          <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
          大纲
        </button>
      </div>
    </div>

    <!-- 主要编辑区域 -->
    <div class="flex h-[calc(100vh-64px)]">
      <!-- 左侧大纲面板 -->
      <div v-show="showOutline" class="w-64 border-r bg-gray-50 p-4 overflow-y-auto">
        <h2 class="text-lg font-semibold mb-4">文档大纲</h2>
        <div class="space-y-2">
          <div v-for="(heading, index) in outline" :key="index" 
               class="cursor-pointer hover:bg-gray-200 p-2 rounded-md"
               :class="{'pl-4': heading.level === 2, 'pl-6': heading.level === 3}"
               @click="scrollToHeading(heading)">
            {{ heading.text }}
          </div>
        </div>
      </div>

      <!-- 编辑器区域 -->
      <div class="flex-1 bg-white">
        <Editor
          api-key="REPLACE_WITH_EDITOR_KEY"
          v-model="content"
          :init="{
            height: '100%',
            menubar: true,
            plugins: [
              'advlist', 'autolink', 'lists', 'link', 'image', 'charmap', 'preview',
              'anchor', 'searchreplace', 'visualblocks', 'code', 'fullscreen',
              'insertdatetime', 'media', 'table', 'help', 'wordcount'
            ],
            toolbar: 'undo redo | formatselect | ' +
              'bold italic backcolor | alignleft aligncenter ' +
              'alignright alignjustify | bullist numlist outdent indent | ' +
              'removeformat | help',
            content_style: 'body { font-family:Helvetica,Arial,sans-serif; font-size:14px }',
            language: 'zh_CN',
          }"
          @change="handleEditorChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Editor from '@tinymce/tinymce-vue'
import { message } from 'ant-design-vue'

// 状态变量
const content = ref('')
const showOutline = ref(false)
const outline = ref([])
const lastSaved = ref('未保存')

// 切换大纲显示
const toggleOutline = () => {
  showOutline.value = !showOutline.value
}

// 保存文档
const saveDocument = async () => {
  try {
    // 这里添加保存到后端的逻辑
    lastSaved.value = new Date().toLocaleTimeString()
    message.success('文档保存成功！')
  } catch (error) {
    message.error('保存失败，请重试！')
  }
}

// 导出为Word文档
const exportDocument = () => {
  try {
    // 这里添加导出Word的逻辑
    const blob = new Blob([content.value], { type: 'application/msword' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = '文档.doc'
    a.click()
    window.URL.revokeObjectURL(url)
    message.success('文档导出成功！')
  } catch (error) {
    message.error('导出失败，请重试！')
  }
}

// 处理编辑器内容变化
const handleEditorChange = (e) => {
  content.value = e.target.getContent()
  updateOutline()
}

// 更新大纲
const updateOutline = () => {
  const parser = new DOMParser()
  const doc = parser.parseFromString(content.value, 'text/html')
  const headings = doc.querySelectorAll('h1, h2, h3')
  outline.value = Array.from(headings).map(heading => ({
    text: heading.textContent,
    level: parseInt(heading.tagName.charAt(1)),
    id: heading.id
  }))
}

// 滚动到指定标题
const scrollToHeading = (heading) => {
  const element = document.getElementById(heading.id)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
  }
}

onMounted(() => {
  // 初始化编辑器后的操作
})
</script>

<style scoped>
.document-editor {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

:deep(.tox-tinymce) {
  border: none !important;
}

:deep(.tox-statusbar) {
  display: none !important;
}
</style> 
