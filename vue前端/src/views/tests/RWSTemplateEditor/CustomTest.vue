<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { Document, Packer, Paragraph, Table, TableCell, TableRow, TextRun, 
  HeadingLevel, AlignmentType, convertInchesToTwip, BorderStyle } from 'docx'
import DocumentTitleEditor from './components/DocumentTitleEditor.vue'
import DocumentCoverEditor from './components/DocumentCoverEditor.vue'
import DocumentTableEditor from './components/DocumentTableEditor.vue'
import TemplatePreview from './components/TemplatePreview.vue'

// 模板的默认结构
const defaultTemplate = {
  name: '自定义模板',
  // 文档标题部分
  documentTitle: {
    text: '本科毕业设计（论文）任务书',
    fontName: '黑体',
    fontSize: 32,
    alignment: 'CENTER'
  },
  // 封面部分
  coverFields: [
    { id: 1, label: '题目（中文）', field: 'title', type: 'text', required: true, value: '' },
    { id: 2, label: '题目（英文）', field: 'titleEn', type: 'text', required: true, value: '' },
    { id: 3, label: '学生姓名', field: 'studentName', type: 'text', required: true, value: '' },
    { id: 4, label: '学号', field: 'studentId', type: 'text', required: true, value: '' },
    { id: 5, label: '专业班级', field: 'major', type: 'text', required: true, value: '' },
    { id: 6, label: '学院', field: 'college', type: 'text', required: true, value: '计算机与软件工程学院' },
    { id: 7, label: '指导老师', field: 'supervisor', type: 'text', required: true, value: '' },
    { id: 8, label: '教研室负责人', field: 'departmentHead', type: 'text', required: true, value: '' },
    { id: 9, label: '下达日期', field: 'submitDate', type: 'date', required: true, value: '' }
  ],
  // 表格内容部分
  tableContents: [
    {
      id: 1,
      title: '研究背景及意义',
      type: 'longtext',
      minLength: 200,
      value: '',
      required: true
    },
    {
      id: 2,
      title: '研究内容',
      type: 'longtext',
      minLength: 300,
      value: '',
      required: true
    },
    {
      id: 3,
      title: '工作进度',
      type: 'schedule',
      value: [
        { task: '完成选题', startDate: '', endDate: '' },
        { task: '完成任务书下达', startDate: '', endDate: '' },
        { task: '完成开题报告', startDate: '', endDate: '' },
        { task: '完成中期检查', startDate: '', endDate: '' },
        { task: '完成论文初稿', startDate: '', endDate: '' }
      ],
      required: true
    },
    {
      id: 4,
      title: '参考文献',
      type: 'references',
      value: ['', '', ''],
      minCount: 3,
      required: true
    }
  ]
}

// 当前模板状态
const currentTemplate = ref(JSON.parse(JSON.stringify(defaultTemplate)))
const savedTemplates = ref([])
const newTemplateName = ref('')
const currentTemplateIndex = ref(-1)
const activeTab = ref('title') // 'title', 'cover', 'content', 'preview'

// 从localStorage加载已保存的模板
onMounted(() => {
  loadSavedTemplates()
})

// 监听模板变化，自动保存到localStorage
watch(currentTemplate, (newVal) => {
  if (currentTemplateIndex.value >= 0) {
    savedTemplates.value[currentTemplateIndex.value] = JSON.parse(JSON.stringify(newVal))
    saveTemplatesToLocalStorage()
  }
}, { deep: true })

// 加载已保存的模板
const loadSavedTemplates = () => {
  const savedData = localStorage.getItem('custom-docx-templates')
  if (savedData) {
    savedTemplates.value = JSON.parse(savedData)
    if (savedTemplates.value.length > 0) {
      currentTemplate.value = JSON.parse(JSON.stringify(savedTemplates.value[0]))
      currentTemplateIndex.value = 0
    }
  }
}

// 保存模板到localStorage
const saveTemplatesToLocalStorage = () => {
  localStorage.setItem('custom-docx-templates', JSON.stringify(savedTemplates.value))
}

// 创建新模板
const createNewTemplate = () => {
  if (!newTemplateName.value.trim()) {
    alert('请输入模板名称')
    return
  }
  
  const newTemplate = JSON.parse(JSON.stringify(defaultTemplate))
  newTemplate.name = newTemplateName.value
  savedTemplates.value.push(newTemplate)
  currentTemplateIndex.value = savedTemplates.value.length - 1
  currentTemplate.value = JSON.parse(JSON.stringify(newTemplate))
  saveTemplatesToLocalStorage()
  newTemplateName.value = ''
}

// 加载已有模板
const loadTemplate = (index) => {
  currentTemplateIndex.value = index
  currentTemplate.value = JSON.parse(JSON.stringify(savedTemplates.value[index]))
}

// 删除模板
const deleteTemplate = (index) => {
  if (confirm('确定要删除此模板吗？')) {
    savedTemplates.value.splice(index, 1)
    saveTemplatesToLocalStorage()
    
    if (currentTemplateIndex.value === index) {
      if (savedTemplates.value.length > 0) {
        currentTemplateIndex.value = 0
        currentTemplate.value = JSON.parse(JSON.stringify(savedTemplates.value[0]))
      } else {
        currentTemplateIndex.value = -1
        currentTemplate.value = JSON.parse(JSON.stringify(defaultTemplate))
      }
    } else if (currentTemplateIndex.value > index) {
      currentTemplateIndex.value--
    }
  }
}

// 生成文档
const generateDoc = async () => {
  try {
    // 创建文档对象
    const doc = new Document({
      sections: [{
        properties: {
          page: {
            margin: {
              top: convertInchesToTwip(1),
              bottom: convertInchesToTwip(1),
              left: convertInchesToTwip(1.2),
              right: convertInchesToTwip(1.2),
            }
          }
        },
        children: generateDocumentContent()
      }]
    });

    // 生成并下载文档
    const blob = await Packer.toBlob(doc);
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `${currentTemplate.value.name}.docx`;
    link.click();
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error('文档生成错误:', error);
    alert('文档生成失败，请查看控制台获取详细错误信息');
  }
}

// 生成文档内容
const generateDocumentContent = () => {
  const docContent = []
  
  // 添加标题
  docContent.push(
    new Paragraph({
      children: [
        new TextRun({
          text: currentTemplate.value.documentTitle.text,
          font: {
            name: currentTemplate.value.documentTitle.fontName,
            size: currentTemplate.value.documentTitle.fontSize
          }
        })
      ],
      heading: HeadingLevel.TITLE,
      alignment: AlignmentType[currentTemplate.value.documentTitle.alignment]
    })
  )
  
  // 添加空行作为间隔
  docContent.push(
    new Paragraph({
      spacing: { after: 400 }
    })
  )
  
  // 生成封面表格
  if (currentTemplate.value.coverFields.length > 0) {
    const tableRows = currentTemplate.value.coverFields.map(field => {
      return new TableRow({
        children: [
          // 字段标签
          new TableCell({
            children: [
              new Paragraph({
                children: [
                  new TextRun({
                    text: field.label,
                    size: 16 * 2,
                    bold: true,
                    font: "宋体"
                  })
                ],
                alignment: AlignmentType.LEFT,
                spacing: { before: 200, after: 200 }
              })
            ],
            width: { size: 30, type: 'pct' },
            verticalAlign: AlignmentType.CENTER,
            borders: {
              top: { style: BorderStyle.NONE },
              bottom: { style: BorderStyle.NONE },
              left: { style: BorderStyle.NONE },
              right: { style: BorderStyle.NONE }
            }
          }),
          // 字段值
          new TableCell({
            children: [
              new Paragraph({
                children: [
                  new TextRun({
                    text: field.value || ' ',
                    size: 16 * 2,
                    bold: true,
                    font: "宋体"
                  })
                ],
                alignment: AlignmentType.CENTER,
                spacing: { before: 200, after: 200 }
              })
            ],
            borders: {
              top: { style: BorderStyle.NONE },
              bottom: { style: BorderStyle.SINGLE, size: 1 },
              left: { style: BorderStyle.NONE },
              right: { style: BorderStyle.NONE }
            }
          })
        ]
      })
    })
    
    docContent.push(
      new Table({
        width: {
          size: 100,
          type: 'pct'
        },
        borders: {
          top: { style: BorderStyle.NONE },
          bottom: { style: BorderStyle.NONE },
          left: { style: BorderStyle.NONE },
          right: { style: BorderStyle.NONE },
          insideHorizontal: { style: BorderStyle.NONE },
          insideVertical: { style: BorderStyle.NONE }
        },
        rows: tableRows
      })
    )
    
    // 添加空行作为间隔
    docContent.push(
      new Paragraph({
        spacing: { after: 400 }
      })
    )
  }
  
  // 生成主表格内容
  if (currentTemplate.value.tableContents.length > 0) {
    const mainTableRows = []
    
    // 遍历每个内容项，生成对应的表格行
    currentTemplate.value.tableContents.forEach(content => {
      // 创建表格行
      const row = new TableRow({
        children: [
          // 左侧标题单元格
          new TableCell({
            children: generateTitleCellContent(content.title),
            width: { size: 5, type: 'pct' },
            verticalAlign: AlignmentType.CENTER
          }),
          // 右侧内容单元格
          new TableCell({
            children: generateContentCellChildren(content)
          })
        ]
      })
      
      mainTableRows.push(row)
    })
    
    // 添加主表格
    docContent.push(
      new Table({
        width: {
          size: 100,
          type: 'pct'
        },
        borders: {
          top: { style: BorderStyle.SINGLE, size: 1 },
          bottom: { style: BorderStyle.SINGLE, size: 1 },
          left: { style: BorderStyle.SINGLE, size: 1 },
          right: { style: BorderStyle.SINGLE, size: 1 },
          insideHorizontal: { style: BorderStyle.SINGLE, size: 1 },
          insideVertical: { style: BorderStyle.SINGLE, size: 1 }
        },
        rows: mainTableRows
      })
    )
  }
  
  return docContent
}

// 生成标题单元格内容
const generateTitleCellContent = (title) => {
  // 将标题拆分为单个字符，每个字符一行
  const characters = title.split('')
  
  return characters.map(char => 
    new Paragraph({
      children: [
        new TextRun({
          text: char,
          font: {
            name: '宋体',
            size: 24
          },
          bold: true
        })
      ],
      alignment: AlignmentType.CENTER
    })
  )
}

// 根据内容类型生成不同的单元格内容
const generateContentCellChildren = (content) => {
  const children = []
  
  // 添加标题
  children.push(
    new Paragraph({
      children: [
        new TextRun({
          text: content.title,
          font: {
            name: '宋体',
            size: 24,
            bold: true
          }
        })
      ],
      spacing: { before: 200, after: 200 }
    })
  )
  
  // 根据不同类型生成内容
  switch (content.type) {
    case 'text':
    case 'longtext':
      children.push(
        new Paragraph({
          children: [
            new TextRun({
              text: content.value || '',
              font: {
                name: '宋体',
                size: 24
              }
            })
          ],
          indent: { firstLine: 480 }, // 首行缩进2字符
          spacing: { line: 360 } // 1.5倍行距
        })
      )
      break
      
    case 'schedule':
      // 工作进度表
      content.value.forEach((item, index) => {
        children.push(
          new Paragraph({
            children: [
              new TextRun({
                text: `（${index + 1}）${item.task || ''}    ${item.startDate || ''} — ${item.endDate || ''}`,
                font: {
                  name: '宋体',
                  size: 24
                }
              })
            ],
            spacing: { line: 480 } // 增加行距
          })
        )
      })
      break
      
    case 'references':
    case 'list':
      // 参考文献或普通列表
      content.value.forEach((item, index) => {
        children.push(
          new Paragraph({
            children: [
              new TextRun({
                text: content.type === 'references' ? `[${index + 1}] ` : `（${index + 1}）`,
                font: {
                  name: '宋体',
                  size: 24
                }
              }),
              new TextRun({
                text: item || '',
                font: {
                  name: '宋体',
                  size: 24
                }
              })
            ],
            indent: { firstLine: 480 }, // 首行缩进2字符
            spacing: { line: 360 } // 1.5倍行距
          })
        )
      })
      break
  }
  
  return children
}

// 校验表单是否有效
const validateForm = () => {
  // 检查封面字段
  for (const field of currentTemplate.value.coverFields) {
    if (field.required && !field.value) {
      return false
    }
  }
  
  // 检查表格内容
  for (const content of currentTemplate.value.tableContents) {
    if (content.required) {
      if (content.type === 'longtext' && (!content.value || content.value.length < content.minLength)) {
        return false
      } else if (content.type === 'references' && content.value.length < content.minCount) {
        return false
      } else if (content.type === 'schedule') {
        for (const item of content.value) {
          if (!item.task || !item.startDate || !item.endDate) {
            return false
          }
        }
      }
    }
  }
  
  return true
}
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-6">自定义文档模板生成器</h1>
    
    <!-- 模板选择和管理 -->
    <div class="bg-gray-100 p-4 rounded mb-4">
      <h2 class="text-xl mb-2">模板管理</h2>
      <div class="flex gap-4 mb-4">
        <div class="flex-1">
          <input v-model="newTemplateName" placeholder="新模板名称" class="w-full px-3 py-2 border rounded" />
        </div>
        <button @click="createNewTemplate" class="px-4 py-2 bg-blue-500 text-white rounded">创建模板</button>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="(template, index) in savedTemplates" :key="index" 
             class="border p-3 rounded cursor-pointer"
             :class="{'bg-blue-100': index === currentTemplateIndex}">
          <div class="flex justify-between items-center">
            <span>{{ template.name }}</span>
            <div>
              <button @click="loadTemplate(index)" class="px-2 py-1 bg-green-500 text-white rounded mr-2">加载</button>
              <button @click="deleteTemplate(index)" class="px-2 py-1 bg-red-500 text-white rounded">删除</button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 编辑部分的标签页 -->
    <div class="mb-4">
      <div class="flex border-b">
        <button @click="activeTab = 'title'" 
                class="px-4 py-2" 
                :class="{'border-b-2 border-blue-500': activeTab === 'title'}">
          文档标题
        </button>
        <button @click="activeTab = 'cover'" 
                class="px-4 py-2" 
                :class="{'border-b-2 border-blue-500': activeTab === 'cover'}">
          封面字段
        </button>
        <button @click="activeTab = 'content'" 
                class="px-4 py-2" 
                :class="{'border-b-2 border-blue-500': activeTab === 'content'}">
          表格内容
        </button>
        <button @click="activeTab = 'preview'" 
                class="px-4 py-2" 
                :class="{'border-b-2 border-blue-500': activeTab === 'preview'}">
          预览填写
        </button>
      </div>
    </div>
    
    <!-- 标题编辑区 -->
    <div v-if="activeTab === 'title'" class="mb-6">
      <DocumentTitleEditor v-model="currentTemplate.documentTitle" />
    </div>
    
    <!-- 封面字段编辑区 -->
    <div v-if="activeTab === 'cover'" class="mb-6">
      <DocumentCoverEditor v-model="currentTemplate.coverFields" />
    </div>
    
    <!-- 表格内容编辑区 -->
    <div v-if="activeTab === 'content'" class="mb-6">
      <DocumentTableEditor v-model="currentTemplate.tableContents" />
    </div>
    
    <!-- 预览和填写 -->
    <div v-if="activeTab === 'preview'" class="mb-6">
      <TemplatePreview 
        v-model="currentTemplate" 
        @generate="generateDoc" 
        :is-valid="validateForm()" 
      />
    </div>
  </div>
</template>

<style scoped>
.form-group {
  @apply mb-4;
}

.form-group label {
  @apply block mb-2 font-medium;
}

.form-input {
  @apply w-full px-3 py-2 border rounded;
}

.form-textarea {
  @apply w-full px-3 py-2 border rounded;
}

button:disabled {
  @apply opacity-50 cursor-not-allowed;
}
</style>
