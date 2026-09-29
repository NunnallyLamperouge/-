<template>
  <div class="bg-app theme-transition">
    <!-- 页面标题 -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-app">数据建模助手</h1>
      <p class="mt-2 text-muted">
        导入SQL语句或文件，自动生成标准ER图和UML类图
      </p>
    </div>

    <div class="grid grid-cols-1 gap-8">
      <!-- 输入区域 -->
      <div class="card">
        <div class="p-6">
          <div>
            <h2 class="text-lg font-medium text-app mb-4">SQL输入</h2>
            <div class="space-y-4">
              <!-- SQL文件上传 -->
              <div class="border-2 border-dashed border-app rounded-lg p-6">
                <div class="text-center">
                  <upload-outlined class="text-4xl text-muted" />
                  <p class="mt-2 text-sm text-muted">
                    拖拽SQL文件到此处，或
                    <a-upload
                      :show-upload-list="false"
                      :before-upload="handleFileUpload"
                      accept=".sql"
                    >
                      <a class="text-blue-600 hover:text-blue-500">点击上传</a>
                    </a-upload>
                  </p>
                  <p class="mt-1 text-xs text-muted">支持 .sql 文件</p>
                </div>
              </div>

              <!-- SQL文本输入 -->
              <a-textarea
                v-model:value="sqlInput"
                :rows="10"
                placeholder="或直接输入SQL语句..."
                class="w-full"
              />

              <!-- 操作按钮 -->
              <div class="flex justify-end space-x-4">
                <a-button @click="clearInput" class="btn-secondary">清空</a-button>
                <a-button type="primary" :loading="generating" @click="generateDiagrams" class="btn-primary">
                  生成图表
                </a-button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 解析结果区域 -->
      <div class="card" v-if="parsedData.tables.length > 0">
        <div class="p-6">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-medium text-app mb-4">解析结果</h2>
            <a-button type="primary" @click="exportDatabaseDoc" class="btn-primary flex items-center">
              <file-text-outlined />
              <span class="ml-2">导出结构文档</span>
            </a-button>
          </div>
          <div class="space-y-4">
            <div v-for="table in parsedData.tables" :key="table.name" class="border border-app rounded-lg p-4">
              <div class="flex items-center justify-between">
                <h4 class="text-base font-medium text-app">
                  {{ table.name }}
                </h4>
                <a-tag color="blue">{{ table.type }}</a-tag>
              </div>
              <div class="mt-2 space-y-1">
                <div v-for="field in table.fields" :key="field.name" class="text-sm text-muted">
                  {{ field.name }} ({{ field.type }})
                  <a-tag v-if="field.isPrimary" color="green" size="small">主键</a-tag>
                  <a-tag v-if="field.isForeign" color="orange" size="small">外键</a-tag>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 图表预览区域 - 使用标签页 -->
      <div class="card">
        <div class="p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center">
              <h2 class="text-lg font-medium text-app">图表预览</h2>
              <a-tooltip title="编辑标题">
                <a-button type="text" @click="showTitleEdit = true" class="ml-2">
                  <template #icon><edit-outlined /></template>
                </a-button>
              </a-tooltip>
            </div>
            <div class="space-x-2">
              <a-tooltip title="放大">
                <a-button shape="circle" @click="handleZoomIn">
                  <template #icon><zoom-in-outlined /></template>
                </a-button>
              </a-tooltip>
              <a-tooltip title="缩小">
                <a-button shape="circle" @click="handleZoomOut">
                  <template #icon><zoom-out-outlined /></template>
                </a-button>
              </a-tooltip>
              <a-tooltip title="重置">
                <a-button shape="circle" @click="handleResetZoom">
                  <template #icon><undo-outlined /></template>
                </a-button>
              </a-tooltip>
              <a-tooltip title="导出图片">
                <a-button shape="circle" @click="handleExport">
                  <template #icon><download-outlined /></template>
                </a-button>
              </a-tooltip>
            </div>
          </div>

          <!-- 图表容器 -->
          <div v-if="parsedData.tables.length === 0" class="flex items-center justify-center h-[800px] border border-app rounded-lg bg-gray-50 dark:bg-gray-800">
            <div class="text-center text-muted">
              <database-outlined class="text-5xl mb-4" />
              <p>请输入SQL语句并点击"生成图表"按钮</p>
            </div>
          </div>
          
          <div v-else class="flex flex-col h-full">
            <!-- 图表类型切换 -->
            <a-tabs v-model:activeKey="activeTabKey" class="w-full mb-4" type="card">
              <a-tab-pane key="er" tab="ER图">
                <ERDiagram
                  ref="erDiagram"
                  :data="parsedData"
                  :theme="theme"
                  :chartTitle="chartTitle"
                  class="w-full h-[680px] border border-app border-2 rounded-lg mb-4 overflow-hidden"
                />
                
                <!-- ER图文本描述区域 -->
                <div class="border border-app rounded-lg p-4 bg-gray-50 dark:bg-gray-800 overflow-auto h-[150px]">
                  <h3 class="text-sm font-medium text-app mb-2">ER图文本描述</h3>
                  <div class="text-xs text-muted">
                    <div v-for="(table, index) in parsedData.tables" :key="index" class="mb-1">
                      <span class="font-medium">实体：{{ table.name }}</span>
                      <span class="ml-2">({{ table.fields.map(f => f.name).join(', ') }})</span>
                    </div>
                    <div v-for="(rel, index) in parsedData.relationships" :key="`rel-${index}`" class="mb-1">
                      <span class="font-medium">关系：</span>
                      <span>{{ rel.source }}.{{ rel.sourceField }} {{ getRelationshipSymbol(rel.type) }} {{ rel.target }}.{{ rel.targetField }}</span>
                    </div>
                  </div>
                </div>
              </a-tab-pane>
              
              <a-tab-pane key="uml" tab="UML类图">
                <UMLClassDiagram
                  ref="umlDiagram"
                  :data="parsedData"
                  :theme="theme"
                  :chartTitle="umlChartTitle"
                  class="w-full h-[680px] border border-app border-2 rounded-lg mb-4 overflow-hidden"
                />
                
                <!-- UML类图文本描述区域 -->
                <div class="border border-app rounded-lg p-4 bg-gray-50 dark:bg-gray-800 overflow-auto h-[150px]">
                  <h3 class="text-sm font-medium text-app mb-2">UML类图文本描述</h3>
                  <div class="text-xs text-muted">
                    <div v-for="(table, index) in parsedData.tables" :key="index" class="mb-1">
                      <span class="font-medium">类：{{ table.name }}</span>
                      <span class="ml-2">
                        <span class="text-pink-600 dark:text-pink-400">
                          <span v-if="table.fields.some(f => f.isPrimary)">
                            - 主键: {{ table.fields.filter(f => f.isPrimary).map(f => f.name + ": " + f.type).join(', ') }}
                          </span>
                        </span>
                        <span class="text-purple-600 dark:text-purple-400">
                          <span v-if="table.fields.some(f => f.isForeign && !f.isPrimary)">
                            - 外键: {{ table.fields.filter(f => f.isForeign && !f.isPrimary).map(f => f.name + ": " + f.type).join(', ') }}
                          </span>
                        </span>
                        <span>
                          <span v-if="table.fields.some(f => !f.isPrimary && !f.isForeign)">
                            - 属性: {{ table.fields.filter(f => !f.isPrimary && !f.isForeign).map(f => f.name + ": " + f.type).join(', ') }}
                          </span>
                        </span>
                      </span>
                    </div>
                    <div v-for="(rel, index) in parsedData.relationships" :key="`rel-${index}`" class="mb-1">
                      <span class="font-medium">关联：</span>
                      <span>{{ rel.source }} {{ getUMLRelationshipSymbol(rel.type) }} {{ rel.target }}</span>
                    </div>
                  </div>
                </div>
              </a-tab-pane>
            </a-tabs>
          </div>
        </div>
      </div>
    </div>

    <!-- 导出选项对话框 -->
    <a-modal
      v-model:visible="exportModalVisible"
      title="导出图表"
      @ok="handleExportConfirm"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-app">
            导出格式
          </label>
          <a-radio-group v-model:value="exportFormat" class="mt-2">
            <a-radio value="png">PNG</a-radio>
            <a-radio value="svg">SVG</a-radio>
            <a-radio value="pdf">PDF</a-radio>
          </a-radio-group>
        </div>
        <div>
          <label class="block text-sm font-medium text-app">
            图片尺寸
          </label>
          <div class="mt-2 flex items-center space-x-4">
            <a-input-number v-model:value="exportWidth" addon-after="px" />
            <span class="text-muted">×</span>
            <a-input-number v-model:value="exportHeight" addon-after="px" />
          </div>
        </div>
      </div>
    </a-modal>

    <!-- 标题编辑对话框 -->
    <a-modal
      v-model:visible="showTitleEdit"
      title="编辑图表标题"
      @ok="handleTitleEditConfirm"
      :okText="'确定'"
      :cancelText="'取消'"
      @open="handleTitleEditOpen"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-app mb-2">
            图表标题
          </label>
          <a-input 
            v-model:value="editingTitle" 
            placeholder="请输入图表标题" 
            :maxLength="30"
            ref="titleInput"
            @pressEnter="handleTitleEditConfirm"
          />
          <div class="mt-1 text-xs text-muted text-right">
            {{ editingTitle.length }}/30
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick, watch } from 'vue'
import { message } from 'ant-design-vue'
import {
  UploadOutlined,
  ZoomInOutlined,
  ZoomOutOutlined,
  UndoOutlined,
  DownloadOutlined,
  DatabaseOutlined,
  EditOutlined,
  FileTextOutlined
} from '@ant-design/icons-vue'
import ERDiagram from '@/components/ERDiagram.vue'
import UMLClassDiagram from '@/components/UMLClassDiagram.vue'
import { parseSQL, analyzeRelationshipTypes } from '@/utility/sqlParser'
import { Document, Packer, Paragraph, Table, TableRow, TableCell, BorderStyle, WidthType, TextRun, AlignmentType } from 'docx'

// 状态变量
const sqlInput = ref('')
const generating = ref(false)
const parsedData = ref({ tables: [], relationships: [] })
const erDiagram = ref(null)
const umlDiagram = ref(null)
const exportModalVisible = ref(false)
const exportFormat = ref('png')
const exportWidth = ref(1920)
const exportHeight = ref(1080)

// 标签页控制
const activeTabKey = ref('er')

// 图表标题相关
const chartTitle = ref('ER图')
const umlChartTitle = ref('UML类图')
const editingTitle = ref('ER图')
const showTitleEdit = ref(false)

// 获取当前主题
const theme = computed(() => {
  return document.documentElement.classList.contains('dark') ? 'dark' : 'light'
})

// 文件上传处理
const handleFileUpload = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    sqlInput.value = e.target.result
  }
  reader.readAsText(file)
  return false
}

// 清空输入
const clearInput = () => {
  sqlInput.value = ''
  parsedData.value = { tables: [], relationships: [] }
  
  // 重置图表标题
  chartTitle.value = 'ER图'
  umlChartTitle.value = 'UML类图'
  editingTitle.value = 'ER图'
}

// 生成图表
const generateDiagrams = async () => {
  if (!sqlInput.value) {
    message.warning('请输入SQL语句或上传SQL文件')
    return
  }

  generating.value = true
  try {
    // 解析SQL语句
    const result = parseSQL(sqlInput.value)
    
    // 分析表之间的关系类型
    const enhancedResult = analyzeRelationshipTypes(result)
    
    // 更新解析结果
    parsedData.value = enhancedResult
    
    if (enhancedResult.tables.length === 0) {
      message.warning('未能解析出任何表结构，请检查SQL语句')
    } else {
      message.success(`成功解析 ${enhancedResult.tables.length} 张表和 ${enhancedResult.relationships.length} 个关系`)
      
      // 根据当前选中的标签页确保编辑标题框的值与当前标题一致
      if (activeTabKey.value === 'er') {
        editingTitle.value = chartTitle.value
      } else {
        editingTitle.value = umlChartTitle.value
      }
      
      // 确保图表组件在下一个渲染周期更新
      nextTick(() => {
        // 不管当前选中哪个标签页，都预先渲染两种图表
        // 先渲染当前显示的图表
        if (activeTabKey.value === 'er') {
          renderERDiagram();
          // 延迟渲染UML图表以确保DOM已完全准备好
          setTimeout(() => {
            renderUMLDiagram();
          }, 500);
        } else {
          renderUMLDiagram();
          // 延迟渲染ER图表以确保DOM已完全准备好
          setTimeout(() => {
            renderERDiagram();
          }, 500);
        }
      });
    }
  } catch (error) {
    console.error('SQL解析错误:', error)
    message.error('SQL解析失败: ' + (error.message || '未知错误'))
  } finally {
    generating.value = false
  }
}

// 渲染ER图
const renderERDiagram = () => {
  if (erDiagram.value && parsedData.value.tables.length > 0) {
    // 给ER图表容器一些时间来调整大小
    setTimeout(() => {
      // 强制重新初始化ER图表
      erDiagram.value.updateChart();
    }, 300);
  }
}

// 渲染UML类图
const renderUMLDiagram = () => {
  if (umlDiagram.value && parsedData.value.tables.length > 0) {
    // 给UML类图表容器一些时间来调整大小
    setTimeout(() => {
      // 强制重新初始化UML类图表
      umlDiagram.value.updateChart();
    }, 300);
  }
}

// 缩放控制
const handleZoomIn = () => {
  if (activeTabKey.value === 'er' && erDiagram.value) {
    erDiagram.value.zoomIn()
  } else if (activeTabKey.value === 'uml' && umlDiagram.value) {
    umlDiagram.value.zoomIn()
  }
}

const handleZoomOut = () => {
  if (activeTabKey.value === 'er' && erDiagram.value) {
    erDiagram.value.zoomOut()
  } else if (activeTabKey.value === 'uml' && umlDiagram.value) {
    umlDiagram.value.zoomOut()
  }
}

const handleResetZoom = () => {
  if (activeTabKey.value === 'er' && erDiagram.value) {
    erDiagram.value.resetZoom()
  } else if (activeTabKey.value === 'uml' && umlDiagram.value) {
    umlDiagram.value.resetZoom()
  }
}

// 导出处理
const handleExport = () => {
  if (!parsedData.value.tables.length) {
    message.warning('请先生成图表')
    return
  }
  exportModalVisible.value = true
}

const handleExportConfirm = () => {
  let diagramRef = null;
  let diagramTitle = '';
  
  if (activeTabKey.value === 'er') {
    diagramRef = erDiagram.value;
    diagramTitle = chartTitle.value;
  } else if (activeTabKey.value === 'uml') {
    diagramRef = umlDiagram.value;
    diagramTitle = umlChartTitle.value;
  }
  
  if (!diagramRef) return;

  const url = diagramRef.exportChart(exportFormat.value, 2)

  // 使用自定义标题作为文件名，去除特殊字符
  const safeTitle = diagramTitle.replace(/[^\w\u4e00-\u9fa5]/g, '_')
  const fileName = safeTitle || (activeTabKey.value === 'er' ? 'er-diagram' : 'uml-class-diagram')

  const link = document.createElement('a')
  link.download = `${fileName}.${exportFormat.value}`
  link.href = url
  link.click()

  exportModalVisible.value = false
  message.success('导出成功')
}

// 示例SQL
const loadExampleSQL = () => {
  sqlInput.value = `
CREATE TABLE users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(50) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  title VARCHAR(200) NOT NULL,
  content TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE comments (
  id INT PRIMARY KEY AUTO_INCREMENT,
  post_id INT NOT NULL,
  user_id INT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE tags (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE post_tags (
  post_id INT NOT NULL,
  tag_id INT NOT NULL,
  PRIMARY KEY (post_id, tag_id),
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);
  `.trim()
}

// 生命周期钩子
onMounted(() => {
  // 加载示例SQL
  loadExampleSQL()
})

// 获取ER关系类型的符号
const getRelationshipSymbol = (type) => {
  switch (type) {
    case 'one-to-one':
      return '1:1 (一对一)';
    case 'one-to-many':
      return '1:N (一对多)';
    case 'many-to-many':
      return 'N:M (多对多)';
    default:
      return '关联';
  }
}

// 获取UML关系类型的符号
const getUMLRelationshipSymbol = (type) => {
  switch (type) {
    case 'one-to-one':
      return '←→ (关联)';
    case 'one-to-many':
      return '←→ (聚合)';
    case 'many-to-many':
      return '←→ (多重关联)';
    default:
      return '←→';
  }
}

// 处理标题编辑确认
const handleTitleEditConfirm = () => {
  if (editingTitle.value.trim()) {
    // 根据当前选中的标签页更新对应的标题
    if (activeTabKey.value === 'er') {
      chartTitle.value = editingTitle.value.trim()
      // 如果ER图表已经生成，更新图表
      if (erDiagram.value && parsedData.value.tables.length > 0) {
        nextTick(() => {
          erDiagram.value.updateChart()
        })
      }
    } else if (activeTabKey.value === 'uml') {
      umlChartTitle.value = editingTitle.value.trim()
      // 如果UML图表已经生成，更新图表
      if (umlDiagram.value && parsedData.value.tables.length > 0) {
        nextTick(() => {
          umlDiagram.value.updateChart()
        })
      }
    }
    
    showTitleEdit.value = false
  } else {
    message.warning('标题不能为空')
  }
}

// 标题输入框引用
const titleInput = ref(null)

// 处理标题编辑对话框打开
const handleTitleEditOpen = () => {
  // 根据当前选中的标签页设置编辑标题框的值
  if (activeTabKey.value === 'er') {
    editingTitle.value = chartTitle.value
  } else if (activeTabKey.value === 'uml') {
    editingTitle.value = umlChartTitle.value
  }
  
  // 下一个渲染周期后聚焦输入框
  nextTick(() => {
    if (titleInput.value) {
      titleInput.value.focus()
    }
  })
}

// 导出数据库结构文档
const exportDatabaseDoc = async () => {
  if (!parsedData.value.tables.length) {
    message.warning('请先生成图表')
    return
  }

  try {
    // 创建表头行
    const headerRow = new TableRow({
      tableHeader: true,
      children: [
        new TableCell({
          children: [new Paragraph({
            children: [new TextRun({ text: '字段名', bold: true })],
            alignment: AlignmentType.CENTER
          })],
          width: { size: 25, type: WidthType.PERCENTAGE }
        }),
        new TableCell({
          children: [new Paragraph({
            children: [new TextRun({ text: '类型', bold: true })],
            alignment: AlignmentType.CENTER
          })],
          width: { size: 25, type: WidthType.PERCENTAGE }
        }),
        new TableCell({
          children: [new Paragraph({
            children: [new TextRun({ text: '可为空', bold: true })],
            alignment: AlignmentType.CENTER
          })],
          width: { size: 20, type: WidthType.PERCENTAGE }
        }),
        new TableCell({
          children: [new Paragraph({
            children: [new TextRun({ text: '描述', bold: true })],
            alignment: AlignmentType.CENTER
          })],
          width: { size: 30, type: WidthType.PERCENTAGE }
        })
      ]
    })

    // 创建文档内容
    const docContent = []
    
    // 为每个表创建一个表结构
    for (const table of parsedData.value.tables) {
      // 添加表名作为标题
      docContent.push(
        new Paragraph({
          text: `表名：${table.name}`,
          heading: 1,
          spacing: { before: 300, after: 200 }
        })
      )
      
      // 创建字段行
      const rows = [headerRow]
      
      for (const field of table.fields) {
        rows.push(
          new TableRow({
            children: [
              new TableCell({
                children: [
                  new Paragraph({
                    text: field.name,
                    alignment: AlignmentType.LEFT
                  })
                ]
              }),
              new TableCell({
                children: [
                  new Paragraph({
                    text: field.type,
                    alignment: AlignmentType.LEFT
                  })
                ]
              }),
              new TableCell({
                children: [
                  new Paragraph({
                    text: field.isNotNull ? '否' : '是',
                    alignment: AlignmentType.CENTER
                  })
                ]
              }),
              new TableCell({
                children: [
                  new Paragraph({
                    text: '',
                    alignment: AlignmentType.LEFT
                  })
                ]
              })
            ]
          })
        )
      }
      
      // 创建表格并添加到文档
      docContent.push(
        new Table({
          width: { size: 100, type: WidthType.PERCENTAGE },
          borders: {
            top: { style: BorderStyle.SINGLE, size: 1 },
            bottom: { style: BorderStyle.SINGLE, size: 1 },
            left: { style: BorderStyle.SINGLE, size: 1 },
            right: { style: BorderStyle.SINGLE, size: 1 },
            insideHorizontal: { style: BorderStyle.SINGLE, size: 1 },
            insideVertical: { style: BorderStyle.SINGLE, size: 1 }
          },
          rows
        })
      )
      
      // 添加空行作为间隔
      docContent.push(
        new Paragraph({
          spacing: { after: 200 }
        })
      )
    }
    
    // 创建文档
    const doc = new Document({
      sections: [{
        properties: {},
        children: docContent
      }]
    })
    
    // 生成并下载文档
    const blob = await Packer.toBlob(doc)
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = '数据库结构说明文档.docx'
    link.click()
    window.URL.revokeObjectURL(url)
    
    message.success('文档导出成功')
  } catch (error) {
    console.error('文档生成错误:', error)
    message.error('文档导出失败: ' + (error.message || '未知错误'))
  }
}

// 监听标签页切换
watch(() => activeTabKey.value, (newTabKey) => {
  nextTick(() => {
    if (newTabKey === 'er') {
      renderERDiagram();
    } else if (newTabKey === 'uml') {
      renderUMLDiagram();
    }
  });
});
</script>

<style scoped>
/* 暗色模式下的卡片增强效果 */
:deep(.dark) .card {
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.05);
}

/* 暗色模式下的上传区域 */
:deep(.dark) .border-dashed {
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.02);
}

:deep(.dark) .border-dashed:hover {
  border-color: rgba(99, 102, 241, 0.5);
  background: rgba(99, 102, 241, 0.05);
}

/* 暗色模式下的文本区域 */
:deep(.dark) .ant-input {
  background-color: #232838;
  border-color: rgba(255, 255, 255, 0.1);
  color: #f3f4f6;
}

:deep(.dark) .ant-input:focus,
:deep(.dark) .ant-input:hover {
  border-color: #6366f1;
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2);
}

/* 暗色模式下的按钮 */
:deep(.dark) .ant-btn-primary {
  background: linear-gradient(90deg, #6366f1, #8b5cf6);
  border: none;
}

:deep(.dark) .ant-btn-primary:hover {
  background: linear-gradient(90deg, #818cf8, #a78bfa);
}

/* 暗色模式下的图表区域 */
:deep(.dark) .diagram-container {
  background-color: #1a1e2c;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

/* 标签页样式 */
:deep(.ant-tabs-tab) {
  font-weight: 500;
}

:deep(.ant-tabs-card > .ant-tabs-nav .ant-tabs-tab) {
  border-radius: 8px 8px 0 0;
  margin-right: 4px;
  transition: all 0.3s;
}

:deep(.dark) .ant-tabs-card > .ant-tabs-nav .ant-tabs-tab-active {
  background-color: #374151;
  border-color: #4B5563;
}

:deep(.ant-tabs-card > .ant-tabs-nav .ant-tabs-tab-active) {
  background-color: #F3F4F6;
  border-color: #D1D5DB;
}

:deep(.dark) .ant-tabs-tab {
  color: #d1d5db;
}

:deep(.dark) .ant-tabs-tab.ant-tabs-tab-active .ant-tabs-tab-btn {
  color: #818cf8;
}

:deep(.dark) .ant-tabs-ink-bar {
  background-color: #818cf8;
}

/* ER图和UML图容器样式 */
.card {
  height: 100%;
}

.card .p-6 {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.card .p-6 > div:last-child {
  flex-grow: 1;
}

/* 暗色模式下的链接颜色 */
:deep(.dark) a.text-blue-600 {
  color: #818cf8;
}

:deep(.dark) a.text-blue-600:hover {
  color: #a5b4fc;
  text-shadow: 0 0 8px rgba(99, 102, 241, 0.4);
}
</style> 