import { Document, Packer, Paragraph, TextRun, AlignmentType, convertInchesToTwip } from 'docx'
import { 
  fetchImageAndAddToDoc, 
  getImageDimensions, 
  calculateImageDimensions,
  generateTableComponent,
  generateTOCComponent,
  generateScheduleComponent
} from './documentGenerator'

// 生成文档内容
export const generateDocumentContent = async (template, updateProgress, updateStatus) => {
  const content = []
  const totalComponents = template.components.length
  let processedComponents = 0
  
  // 处理所有组件
  for (const component of template.components) {
    // 更新进度
    processedComponents++
    const progressPercent = 10 + Math.floor((processedComponents / totalComponents) * 70)
    updateProgress(progressPercent)
    
    switch (component.type) {
      case 'title':
        updateStatus('处理标题组件...')
        content.push(new Paragraph({
          children: [
            new TextRun({
              text: component.content,
              size: (component.fontSize || 32) * 2,
              font: component.fontFamily || '黑体',
              bold: component.bold !== false
            })
          ],
          alignment: AlignmentType[component.alignment?.toUpperCase() || 'CENTER']
        }))
        break
      case 'table':
      case 'cover-table':
        // 如果有表格标题，先添加标题段落
        if (component.title) {
          content.push(new Paragraph({
            children: [
              new TextRun({
                text: component.title,
                size: (component.titleFontSize || 18) * 2,
                font: component.titleFontFamily || '黑体',
                bold: true
              })
            ],
            alignment: AlignmentType[(component.titleAlignment || 'CENTER').toUpperCase()],
            spacing: { after: 200 }
          }))
        }
        // 添加表格
        content.push(await generateTableComponent(component))
        break
      case 'text':
        content.push(new Paragraph({
          children: [
            new TextRun({
              text: component.content,
              size: (component.fontSize || 16) * 2,
              font: component.fontFamily || '宋体'
            })
          ],
          alignment: AlignmentType[component.alignment?.toUpperCase() || 'LEFT'],
          spacing: { line: (component.lineHeight || 1.5) * 240 }
        }))
        break
      case 'toc':
        content.push(generateTOCComponent(component))
        break
      case 'schedule':
        content.push(generateScheduleComponent(component))
        break
      case 'list':
        updateStatus('处理列表组件...')
        // 添加列表项
        if (component.items && Array.isArray(component.items)) {
          component.items.forEach((item, index) => {
            const bulletOrNumber = component.listType === 'bullet' 
              ? '• ' 
              : `${index + 1}. `;
            
            content.push(new Paragraph({
              children: [
                new TextRun({
                  text: bulletOrNumber + item.content,
                  size: (component.fontSize || 16) * 2,
                  font: component.fontFamily || '宋体'
                })
              ],
              alignment: AlignmentType[(component.alignment || 'LEFT').toUpperCase()],
              spacing: { line: (component.lineHeight || 1.5) * 240, after: 120 },
              indent: { left: 720 } // 约0.5英寸的缩进
            }))
          })
        } else {
          // 如果没有items属性或不是数组，添加一个默认的空列表提示
          content.push(new Paragraph({
            children: [
              new TextRun({
                text: '[空列表]',
                size: (component.fontSize || 16) * 2,
                font: component.fontFamily || '宋体',
                italics: true,
                color: '888888'
              })
            ],
            alignment: AlignmentType[(component.alignment || 'LEFT').toUpperCase()],
            spacing: { line: (component.lineHeight || 1.5) * 240, after: 120 }
          }))
        }
        break
      case 'reference':
        updateStatus('处理参考文献组件...')
        // 添加参考文献标题
        if (component.title) {
          content.push(new Paragraph({
            children: [
              new TextRun({
                text: component.title,
                size: (component.titleFontSize || 18) * 2,
                font: component.titleFontFamily || '黑体',
                bold: true
              })
            ],
            alignment: AlignmentType[(component.alignment || 'LEFT').toUpperCase()],
            spacing: { after: 200 }
          }))
        }
        
        // 添加参考文献条目
        if (component.items && Array.isArray(component.items)) {
          component.items.forEach((item, index) => {
            content.push(new Paragraph({
              children: [
                new TextRun({
                  text: `[${index + 1}] ${item.content}`,
                  size: (component.fontSize || 16) * 2,
                  font: component.fontFamily || '宋体'
                })
              ],
              alignment: AlignmentType[(component.alignment || 'LEFT').toUpperCase()],
              spacing: { line: (component.lineHeight || 1.5) * 240, after: 120 }
            }))
          })
        } else {
          // 如果没有items属性或不是数组，添加一个默认的空参考文献提示
          content.push(new Paragraph({
            children: [
              new TextRun({
                text: '[暂无参考文献]',
                size: (component.fontSize || 16) * 2,
                font: component.fontFamily || '宋体',
                italics: true,
                color: '888888'
              })
            ],
            alignment: AlignmentType[(component.alignment || 'LEFT').toUpperCase()],
            spacing: { line: (component.lineHeight || 1.5) * 240, after: 120 }
          }))
        }
        break
      case 'image':
        try {
          if (component.imageUrl) {
            updateStatus(`加载图片: ${component.imageUrl.substring(0, 30)}...`)
            
            // 获取图片数据
            const imageBuffer = await fetchImageAndAddToDoc(component.imageUrl)
            
            // 获取图片尺寸
            let imageDimensions
            try {
              updateStatus('计算图片尺寸...')
              const dimensions = await getImageDimensions(component.imageUrl)
              // 根据用户设置的最大宽度计算图片尺寸
              imageDimensions = calculateImageDimensions(
                dimensions.width,
                dimensions.height,
                component.maxWidth || 100
              )
            } catch (error) {
              console.warn('无法获取图片尺寸，使用默认尺寸', error)
              imageDimensions = {
                width: 500,
                height: 350
              }
            }
            
            updateStatus('添加图片到文档...')
            
            // 创建图片段落
            const imageParagraph = new Paragraph({
              children: [
                new ImageRun({
                  data: imageBuffer,
                  transformation: {
                    width: imageDimensions.width,
                    height: imageDimensions.height,
                  }
                })
              ],
              alignment: AlignmentType[(component.imageAlignment || 'CENTER').toUpperCase()],
              spacing: { after: component.showCaption ? 120 : 240 }
            })
            
            content.push(imageParagraph)
            
            // 如果有图片说明，添加说明文字
            if (component.showCaption && component.caption) {
              updateStatus('添加图片说明...')
              content.push(new Paragraph({
                children: [
                  new TextRun({
                    text: component.caption,
                    size: (component.captionFontSize || 14) * 2,
                    font: component.captionFontFamily || '宋体'
                  })
                ],
                alignment: AlignmentType[(component.captionAlignment || 'CENTER').toUpperCase()],
                spacing: { line: (component.captionLineHeight || 1.5) * 240, after: 240 }
              }))
            }
          } else {
            // 如果没有图片URL，添加占位符
            content.push(new Paragraph({
              children: [
                new TextRun({
                  text: '[图片位置 - 未提供URL]',
                  size: 24,
                  italics: true
                })
              ],
              alignment: AlignmentType[(component.imageAlignment || 'CENTER').toUpperCase()],
              spacing: { after: 240 }
            }))
          }
        } catch (error) {
          // 处理图片加载失败的情况
          console.error('加载图片失败:', error)
          updateStatus(`图片加载失败: ${error.message}`)
          content.push(new Paragraph({
            children: [
              new TextRun({
                text: `[图片加载失败: ${error.message}]`,
                size: 24,
                italics: true,
                color: 'FF0000'
              })
            ],
            alignment: AlignmentType[(component.imageAlignment || 'CENTER').toUpperCase()],
            spacing: { after: 240 }
          }))
          
          // 如果有图片说明，仍然添加说明文字
          if (component.showCaption && component.caption) {
            content.push(new Paragraph({
              children: [
                new TextRun({
                  text: component.caption,
                  size: (component.captionFontSize || 14) * 2,
                  font: component.captionFontFamily || '宋体'
                })
              ],
              alignment: AlignmentType[(component.captionAlignment || 'CENTER').toUpperCase()],
              spacing: { line: (component.captionLineHeight || 1.5) * 240, after: 240 }
            }))
          }
        }
        break
      default:
        content.push(new Paragraph({ children: [] }))
    }
  }
  
  return content
}

// 生成文档
export const generateDoc = async (template, updateProgress, updateStatus) => {
  try {
    updateProgress(0)
    updateStatus('准备导出文档...')

    // 添加默认文档设置
    const defaultSettings = {
      pageSize: {
        width: 21,
        height: 29.7
      },
      margins: {
        top: 2.54,
        bottom: 2.54,
        left: 3.18,
        right: 3.18
      }
    }

    // 使用默认值作为后备
    const documentSettings = template?.documentSettings ?? defaultSettings
    
    updateStatus('生成文档内容...')
    updateProgress(10)
    
    const documentContent = await generateDocumentContent(template, updateProgress, updateStatus)
    
    updateStatus('创建文档...')
    updateProgress(80)
    
    const doc = new Document({
      sections: [{
        properties: {
          page: {
            margin: {
              top: convertInchesToTwip((documentSettings.margins.top || 2.54) / 2.54),
              bottom: convertInchesToTwip((documentSettings.margins.bottom || 2.54) / 2.54),
              left: convertInchesToTwip((documentSettings.margins.left || 3.18) / 2.54),
              right: convertInchesToTwip((documentSettings.margins.right || 3.18) / 2.54)
            }
          }
        },
        children: documentContent
      }]
    })

    updateStatus('打包文档...')
    updateProgress(90)
    
    const blob = await Packer.toBlob(doc)
    
    updateStatus('下载文档...')
    updateProgress(95)
    
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${template?.name || '新文档'}.docx`
    link.click()
    window.URL.revokeObjectURL(url)
    
    updateProgress(100)
    updateStatus('导出完成！')
    
    return true
  } catch (error) {
    console.error('文档生成错误:', error)
    updateStatus(`导出失败: ${error.message}`)
    throw error
  }
} 