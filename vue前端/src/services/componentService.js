import { get, post, postJSON } from '@/utility/request'

// 获取模板的所有组件
export const getTemplateComponents = (templateId, callback, errorCallback) => {
  get(`/api/template/${templateId}/components`, null, 
    (msg, data) => {
      if (callback) callback(msg, data)
    },
    (msg) => {
      if (errorCallback) errorCallback(msg)
    }
  )
}

// 获取单个组件详情
export const getComponent = (componentId, callback, errorCallback) => {
  get(`/api/component/${componentId}`, null, 
    (msg, data) => {
      if (callback) callback(msg, data)
    },
    (msg) => {
      if (errorCallback) errorCallback(msg)
    }
  )
}

// 创建新组件
export const createComponent = (templateId, component, callback, errorCallback) => {
  postJSON(`/api/template/${templateId}/component`, component, 
    (msg, data) => {
      if (callback) callback(msg, data)
    },
    (msg) => {
      if (errorCallback) errorCallback(msg)
    }
  )
}

// 更新组件
export const updateComponent = (componentId, componentData, callback, errorCallback) => {
  postJSON(`/api/component/${componentId}`, componentData, 
    (msg, data) => {
      if (callback) callback(msg, data)
    },
    (msg) => {
      if (errorCallback) errorCallback(msg)
    }
  )
}

// 删除组件
export const deleteComponent = (componentId, callback, errorCallback) => {
  get(`/api/component/${componentId}`, null, 
    {
      method: 'DELETE',
    },
    (msg, data) => {
      if (callback) callback(msg, data)
    },
    (msg) => {
      if (errorCallback) errorCallback(msg)
    }
  )
}

// 更新组件顺序
export const updateComponentOrder = (templateId, componentIds, callback, errorCallback) => {
  postJSON(`/api/template/${templateId}/component-order`, { component_ids: componentIds }, 
    (msg, data) => {
      if (callback) callback(msg, data)
    },
    (msg) => {
      if (errorCallback) errorCallback(msg)
    }
  )
}

// 获取组件评论
export const getComponentComments = (componentId, callback, errorCallback) => {
  get(`/api/component/${componentId}/comments`, null, 
    (msg, data) => {
      if (callback) callback(msg, data)
    },
    (msg) => {
      if (errorCallback) errorCallback(msg)
    }
  )
}

// 添加组件评论
export const addComponentComment = (componentId, comment, callback, errorCallback) => {
  postJSON(`/api/component/${componentId}/comment`, comment, 
    (msg, data) => {
      if (callback) callback(msg, data)
    },
    (msg) => {
      if (errorCallback) errorCallback(msg)
    }
  )
}

// 更新评论状态
export const updateCommentStatus = (commentId, status, callback, errorCallback) => {
  postJSON(`/api/comment/${commentId}/status`, { status }, 
    (msg, data) => {
      if (callback) callback(msg, data)
    },
    (msg) => {
      if (errorCallback) errorCallback(msg)
    }
  )
}

// 回复评论
export const replyToComment = (commentId, content, callback, errorCallback) => {
  postJSON(`/api/comment/${commentId}/reply`, { content }, 
    (msg, data) => {
      if (callback) callback(msg, data)
    },
    (msg) => {
      if (errorCallback) errorCallback(msg)
    }
  )
}

// 将前端组件数据转换为后端存储格式
export const convertComponentToBackend = (component) => {
  // 基础组件数据
  const baseData = {
    component_id: component.id,
    type: component.type,
    position: component.position,
    order: component.order || 0
  }
  
  // 根据组件类型添加特定数据
  const data = {}
  
  switch (component.type) {
    case 'title':
      Object.assign(data, {
        content: component.content,
        font_size: component.fontSize,
        font_family: component.fontFamily,
        alignment: component.alignment,
        bold: component.bold
      })
      break
    case 'text':
      Object.assign(data, {
        content: component.content,
        font_size: component.fontSize,
        font_family: component.fontFamily,
        alignment: component.alignment,
        line_height: component.lineHeight
      })
      break
    case 'table':
      Object.assign(data, {
        title: component.title,
        layout: component.layout,
        rows: component.rows,
        columns: component.columns,
        data: component.data,
        header_row: component.headerRow,
        borders: component.borders,
        title_alignment: component.titleAlignment,
        title_font_family: component.titleFontFamily,
        title_font_size: component.titleFontSize,
        column_widths: component.columnWidths,
        fields: component.fields,
        settings: component.settings
      })
      break
    case 'image':
      Object.assign(data, {
        image_url: component.imageUrl,
        image_alignment: component.imageAlignment,
        max_width: component.maxWidth,
        max_height: component.maxHeight,
        show_caption: component.showCaption,
        caption: component.caption,
        caption_alignment: component.captionAlignment,
        caption_font_size: component.captionFontSize,
        caption_font_family: component.captionFontFamily,
        caption_line_height: component.captionLineHeight
      })
      break
    case 'list':
      Object.assign(data, {
        list_type: component.listType,
        font_size: component.fontSize,
        font_family: component.fontFamily,
        line_height: component.lineHeight,
        alignment: component.alignment,
        items: component.items
      })
      break
    case 'reference':
      Object.assign(data, {
        title: component.title,
        title_font_family: component.titleFontFamily,
        title_font_size: component.titleFontSize,
        font_size: component.fontSize,
        font_family: component.fontFamily,
        line_height: component.lineHeight,
        alignment: component.alignment,
        items: component.items
      })
      break
    case 'toc':
      Object.assign(data, {
        title: component.title,
        title_font_family: component.titleFontFamily,
        title_font_size: component.titleFontSize,
        font_size: component.fontSize,
        font_family: component.fontFamily,
        line_height: component.lineHeight,
        alignment: component.alignment,
        items: component.items
      })
      break
    case 'schedule':
      Object.assign(data, {
        items: component.items,
        settings: component.settings
      })
      break
    case 'cover-table':
      Object.assign(data, {
        fields: component.fields,
        settings: component.settings
      })
      break
    default:
      // 对于未知类型，保留所有属性
      Object.assign(data, component)
  }
  
  return {
    ...baseData,
    data
  }
}

// 将后端组件数据转换为前端格式
export const convertComponentToFrontend = (backendComponent) => {
  if (!backendComponent) return null
  
  // 基础组件数据
  const baseComponent = {
    id: backendComponent.component_id,
    type: backendComponent.type,
    position: backendComponent.position,
    order: backendComponent.order || 0,
    style: getDefaultStyleForType(backendComponent.type)
  }
  
  // 根据组件类型添加特定属性
  const data = backendComponent.data || {}
  
  switch (backendComponent.type) {
    case 'title':
      return {
        ...baseComponent,
        content: data.content || '',
        fontSize: data.font_size || 32,
        fontFamily: data.font_family || '黑体',
        alignment: data.alignment || 'center',
        bold: data.bold !== undefined ? data.bold : true
      }
    case 'text':
      return {
        ...baseComponent,
        content: data.content || '',
        fontSize: data.font_size || 16,
        fontFamily: data.font_family || '宋体',
        alignment: data.alignment || 'left',
        lineHeight: data.line_height || 1.5
      }
    case 'table':
      return {
        ...baseComponent,
        title: data.title || '表格标题',
        layout: data.layout || 'key-value',
        rows: data.rows || 3,
        columns: data.columns || 2,
        data: data.data || Array(3).fill().map(() => Array(2).fill('')),
        headerRow: data.header_row !== undefined ? data.header_row : true,
        borders: data.borders !== undefined ? data.borders : true,
        titleAlignment: data.title_alignment || 'center',
        titleFontFamily: data.title_font_family || '黑体',
        titleFontSize: data.title_font_size || 18,
        columnWidths: data.column_widths || [50, 50],
        fields: data.fields || [],
        settings: data.settings || {}
      }
    case 'image':
      return {
        ...baseComponent,
        imageUrl: data.image_url || '',
        imagePreviewError: false,
        imageAlignment: data.image_alignment || 'center',
        maxWidth: data.max_width || 100,
        maxHeight: data.max_height || 'auto',
        showCaption: data.show_caption !== undefined ? data.show_caption : true,
        caption: data.caption || '图片说明',
        captionAlignment: data.caption_alignment || 'center',
        captionFontSize: data.caption_font_size || 14,
        captionFontFamily: data.caption_font_family || '宋体',
        captionLineHeight: data.caption_line_height || 1.5
      }
    case 'list':
      return {
        ...baseComponent,
        listType: data.list_type || 'bullet',
        fontSize: data.font_size || 16,
        fontFamily: data.font_family || '宋体',
        lineHeight: data.line_height || 1.5,
        alignment: data.alignment || 'left',
        items: data.items || [
          { id: 1, content: '列表项 1' },
          { id: 2, content: '列表项 2' },
          { id: 3, content: '列表项 3' }
        ]
      }
    case 'reference':
      return {
        ...baseComponent,
        title: data.title || '参考文献',
        titleFontFamily: data.title_font_family || '黑体',
        titleFontSize: data.title_font_size || 18,
        fontSize: data.font_size || 16,
        fontFamily: data.font_family || '宋体',
        lineHeight: data.line_height || 1.5,
        alignment: data.alignment || 'left',
        items: data.items || []
      }
    case 'toc':
      return {
        ...baseComponent,
        title: data.title || '目录',
        titleFontFamily: data.title_font_family || '黑体',
        titleFontSize: data.title_font_size || 18,
        fontSize: data.font_size || 16,
        fontFamily: data.font_family || '宋体',
        lineHeight: data.line_height || 1.8,
        alignment: data.alignment || 'left',
        items: data.items || []
      }
    case 'schedule':
      return {
        ...baseComponent,
        items: data.items || [],
        settings: data.settings || {
          fontSize: 16,
          fontFamily: '宋体',
          lineHeight: 1.5
        }
      }
    case 'cover-table':
      return {
        ...baseComponent,
        fields: data.fields || [],
        settings: data.settings || {}
      }
    default:
      // 对于未知类型，保留所有属性
      return {
        ...baseComponent,
        ...data
      }
  }
}

// 获取组件默认样式
const getDefaultStyleForType = (type) => {
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