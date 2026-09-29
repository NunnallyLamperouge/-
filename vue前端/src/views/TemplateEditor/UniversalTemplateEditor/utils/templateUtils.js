// 组件ID计数器
let componentIdCounter = 1;

// 默认模板结构
export const defaultTemplate = {
  id: '',
  name: '新模板',
  type: 'custom', // 'thesis', 'proposal', 'task', 'review'
  components: [], // 组件列表
  documentSettings: {
    pageSize: {
      width: 21, // cm
      height: 29.7 // cm
    },
    margins: {
      top: 2.54, // cm
      bottom: 2.54,
      left: 3.18,
      right: 3.18
    }
  }
}

// 从localStorage加载已保存的模板
export const loadSavedTemplates = () => {
  const savedData = localStorage.getItem('universal-docx-templates')
  if (savedData) {
    try {
      return JSON.parse(savedData)
    } catch (e) {
      console.error('解析保存的模板失败:', e)
      return []
    }
  }
  return []
}

// 保存模板到localStorage
export const saveTemplatesToLocalStorage = (templates) => {
  localStorage.setItem('universal-docx-templates', JSON.stringify(templates))
}

// 创建新组件
export const createComponent = (type, position) => {
  const baseComponent = {
    id: `${type}_${Date.now()}`,
    type,
    position: { ...position },
    style: getDefaultStyleForType(type)
  }
  
  // 根据类型添加特定属性
  switch (type) {
    case 'title':
      return {
        ...baseComponent,
        content: '标题',
        fontSize: 32,
        fontFamily: '黑体',
        alignment: 'center',
        bold: true
      }
    case 'table':
      return {
        ...baseComponent,
        title: '表格标题',
        type: 'table',
        layout: 'key-value', // 'key-value' 或 'title-content'
        rows: 3,
        columns: 2,
        data: Array(3).fill().map(() => Array(2).fill('')),
        headerRow: true,
        borders: true,
        titleAlignment: 'center', // 默认居中对齐
        titleFontFamily: '黑体', // 默认标题字体
        titleFontSize: 18, // 默认标题字号
        columnWidths: [50, 50], // 添加列宽设置，默认两列各占50%
        fields: [
          {
            id: 1,
            label: '字段1',
            field: 'field1',
            type: 'text',
            required: true,
            value: '',
            description: '',
            contentComponents: [] // 添加内容组件数组，用于存储内容列的子组件
          },
          {
            id: 2,
            label: '字段2',
            field: 'field2',
            type: 'text',
            required: true,
            value: '',
            description: '',
            contentComponents: [] // 添加内容组件数组，用于存储内容列的子组件
          }
        ],
        settings: {
          fontSize: 16,
          fontFamily: '宋体',
          headerStyle: {
            bold: true,
            alignment: 'center'
          },
          cellStyle: {
            alignment: 'left',
            padding: '8px'
          },
          borderStyle: {
            top: true,
            bottom: true,
            left: true,
            right: true,
            insideHorizontal: true,
            insideVertical: true
          }
        }
      }
    case 'cover-table':
      return {
        ...baseComponent,
        type: 'cover-table',
        fields: [
          {
            id: 1,
            label: '题目',
            value: '',
            width: 30 // 标签列宽度百分比
          }
        ],
        settings: {
          fontSize: 16,
          fontFamily: '宋体',
          labelStyle: {
            bold: true,
            alignment: 'left'
          },
          valueStyle: {
            alignment: 'center'
          }
        }
      }
    case 'text':
      return {
        ...baseComponent,
        content: '',
        fontSize: 16,
        fontFamily: '宋体',
        alignment: 'left',
        lineHeight: 1.5
      }
    case 'toc':
      return {
        ...baseComponent,
        type: 'toc',
        title: '目录',
        titleFontFamily: '黑体',
        titleFontSize: 18,
        fontSize: 16,
        fontFamily: '宋体',
        lineHeight: 1.8,
        alignment: 'left',
        items: [
          {
            id: 1,
            number: '1',
            title: '绪论',
            page: '1'
          },
          {
            id: 2,
            number: '1.1',
            title: '研究背景',
            page: '1'
          },
          {
            id: 3,
            number: '1.2',
            title: '研究意义',
            page: '2'
          },
          {
            id: 4,
            number: '2',
            title: '相关工作',
            page: '3'
          }
        ]
      }
    case 'schedule':
      return {
        ...baseComponent,
        type: 'schedule',
        items: [
          { task: '完成选题', startDate: '', endDate: '' },
          { task: '完成开题报告', startDate: '', endDate: '' },
          { task: '完成中期检查', startDate: '', endDate: '' }
        ],
        settings: {
          fontSize: 16,
          fontFamily: '宋体',
          lineHeight: 1.5
        }
      }
    case 'image':
      return {
        ...baseComponent,
        type: 'image',
        imageUrl: '',
        imagePreviewError: false,
        imageAlignment: 'center',
        maxWidth: 100, // 百分比
        maxHeight: 'auto',
        showCaption: true,
        caption: '图片说明（例如：图1-1 系统概览）',
        captionAlignment: 'center',
        captionFontSize: 14,
        captionFontFamily: '宋体',
        captionLineHeight: 1.5
      }
    case 'list':
      return {
        ...baseComponent,
        type: 'list',
        listType: 'bullet', // 'bullet' 或 'number'
        fontSize: 16,
        fontFamily: '宋体',
        lineHeight: 1.5,
        alignment: 'left',
        items: [
          { id: 1, content: '列表项 1' },
          { id: 2, content: '列表项 2' },
          { id: 3, content: '列表项 3' }
        ]
      }
    case 'reference':
      return {
        ...baseComponent,
        type: 'reference',
        title: '参考文献',
        titleFontFamily: '黑体',
        titleFontSize: 18,
        fontSize: 16,
        fontFamily: '宋体',
        lineHeight: 1.5,
        alignment: 'left',
        items: [
          { id: 1, content: '作者. 文献标题[J]. 期刊名称, 年份, 卷(期): 起止页码.' },
          { id: 2, content: '作者. 书名[M]. 出版地: 出版社, 出版年: 起止页码.' },
          { id: 3, content: '作者. 论文题目[D]. 保存地: 保存单位, 年份.' }
        ]
      }
    default:
      return baseComponent
  }
}

// 获取组件默认样式
export const getDefaultStyleForType = (type) => {
  const baseStyle = {
    width: '100%',
    marginTop: '1em',
    marginBottom: '1em'
  }
  
  switch (type) {
    case 'title':
      return {
        ...baseStyle,
        textAlign: 'center',
        fontWeight: 'bold'
      }
    case 'table':
      return {
        ...baseStyle,
        borderCollapse: 'collapse'
      }
    default:
      return baseStyle
  }
}

// 上传图片
export const uploadImage = async (file) => {
  try {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await fetch('http://localhost:8080/file/upload', {
      method: 'POST',
      body: formData
    })
    
    if (!response.ok) {
      throw new Error(`上传失败: ${response.statusText}`)
    }
    
    const data = await response.text()
    // 拼接完整的图片URL
    const imageUrl = `http://localhost:8080${data}`
    return imageUrl
  } catch (error) {
    console.error('图片上传失败:', error)
    throw error
  }
} 