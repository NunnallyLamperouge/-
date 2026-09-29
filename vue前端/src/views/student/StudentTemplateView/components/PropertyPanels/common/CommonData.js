// 字体列表
export const fontFamilies = [
  { value: '宋体', label: '宋体' },
  { value: '黑体', label: '黑体' },
  { value: '楷体', label: '楷体' },
  { value: '仿宋', label: '仿宋' },
  { value: 'Arial', label: 'Arial' },
  { value: 'Times New Roman', label: 'Times New Roman' }
]

// 对齐方式
export const alignments = [
  { value: 'left', label: '左对齐', icon: '⬅️' },
  { value: 'center', label: '居中', icon: '⬆️' },
  { value: 'right', label: '右对齐', icon: '➡️' },
  { value: 'justify', label: '两端对齐', icon: '↔️' }
]

// 确保组件设置存在
export const ensureComponentSettings = (component) => {
  if (!component.settings) {
    component.settings = {}
  }
  
  if (!component.settings.borderStyle) {
    component.settings.borderStyle = {
      top: false,
      bottom: false,
      left: false,
      right: false,
      insideHorizontal: false,
      insideVertical: false
    }
  }
  
  if (!component.settings.cellStyle) {
    component.settings.cellStyle = {
      alignment: 'left',
      padding: '8px'
    }
  }
  
  if (!component.settings.headerStyle) {
    component.settings.headerStyle = {
      bold: true,
      alignment: 'center'
    }
  }
  
  if (!component.settings.fontSize) {
    component.settings.fontSize = 16
  }
  
  if (!component.settings.fontFamily) {
    component.settings.fontFamily = '宋体'
  }
  
  if (!component.settings.lineHeight) {
    component.settings.lineHeight = 1.5
  }
  
  if (!component.settings.labelStyle) {
    component.settings.labelStyle = {
      bold: true,
      alignment: 'left'
    }
  }
  
  if (!component.settings.valueStyle) {
    component.settings.valueStyle = {
      alignment: 'center'
    }
  }
  
  return component.settings
}

// 表格边框预设
export const borderPresets = [
  { name: '无边框', icon: '◻️', preset: { top: false, bottom: false, left: false, right: false, insideHorizontal: false, insideVertical: false } },
  { name: '外边框', icon: '🔲', preset: { top: true, bottom: true, left: true, right: true, insideHorizontal: false, insideVertical: false } },
  { name: '全边框', icon: '📊', preset: { top: true, bottom: true, left: true, right: true, insideHorizontal: true, insideVertical: true } },
  { name: '横线', icon: '☰', preset: { top: true, bottom: true, left: false, right: false, insideHorizontal: true, insideVertical: false } },
  { name: '竖线', icon: '▤', preset: { top: false, bottom: false, left: true, right: true, insideHorizontal: false, insideVertical: true } }
] 