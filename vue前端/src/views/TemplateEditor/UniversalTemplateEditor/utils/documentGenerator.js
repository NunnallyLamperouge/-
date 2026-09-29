import { Document, Packer, Paragraph, Table, TableCell, TableRow, TextRun, 
  HeadingLevel, AlignmentType, convertInchesToTwip, BorderStyle, ImageRun } from 'docx'

// 获取图片并转换为二进制数据
export const fetchImageAndAddToDoc = async (imageUrl) => {
  try {
    // 检查URL是否有效
    if (!imageUrl || !imageUrl.trim()) {
      throw new Error('图片URL为空')
    }

    // 获取图片
    const response = await fetch(imageUrl)
    if (!response.ok) {
      throw new Error(`获取图片失败: ${response.statusText}`)
    }

    // 将图片转换为ArrayBuffer
    const imageBuffer = await response.arrayBuffer()
    
    return imageBuffer
  } catch (error) {
    console.error('获取图片失败:', error)
    throw error
  }
}

// 获取图片尺寸
export const getImageDimensions = async (imageUrl) => {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => {
      resolve({
        width: img.width,
        height: img.height,
        aspectRatio: img.width / img.height
      })
    }
    img.onerror = () => {
      reject(new Error('获取图片尺寸失败'))
    }
    img.src = imageUrl
  })
}

// 计算图片在文档中的尺寸
export const calculateImageDimensions = (originalWidth, originalHeight, maxWidthPercent) => {
  // 文档宽度（单位：像素）- 假设A4纸张宽度约为8.27英寸，1英寸=96像素
  const docWidth = 8.27 * 96
  
  // 根据最大宽度百分比计算图片宽度
  const maxWidth = docWidth * (maxWidthPercent / 100)
  
  // 如果原始宽度小于最大宽度，使用原始尺寸
  if (originalWidth <= maxWidth) {
    return {
      width: originalWidth,
      height: originalHeight
    }
  }
  
  // 否则，按比例缩放
  const aspectRatio = originalWidth / originalHeight
  const height = maxWidth / aspectRatio
  
  return {
    width: maxWidth,
    height: height
  }
}

// 生成表格组件
export const generateTableComponent = async (component) => {
  const rows = []
  
  if (component.type === 'cover-table') {
    // 封面表格布局
    component.fields.forEach(field => {
      rows.push(
        new TableRow({
          children: [
            new TableCell({
              children: [
                new Paragraph({
                  children: [
                    new TextRun({
                      text: field.label,
                      size: (component.settings?.fontSize || 16) * 2,
                      font: component.settings?.fontFamily || '宋体',
                      bold: component.settings?.labelStyle?.bold !== false
                    })
                  ],
                  alignment: AlignmentType[
                    (component.settings?.labelStyle?.alignment || 'left').toUpperCase()
                  ]
                })
              ],
              width: { 
                size: field.width <= 10 ? 1 : (field.width >= 90 ? 99 : field.width), 
                type: 'pct' 
              },
              borders: {
                top: { style: BorderStyle.NONE },
                bottom: { style: BorderStyle.NONE },
                left: { style: BorderStyle.NONE },
                right: { style: BorderStyle.NONE }
              }
            }),
            new TableCell({
              children: [
                new Paragraph({
                  children: [
                    new TextRun({
                      text: field.value || '',
                      size: (component.settings?.fontSize || 16) * 2,
                      font: component.settings?.fontFamily || '宋体'
                    })
                  ],
                  alignment: AlignmentType[
                    (component.settings?.valueStyle?.alignment || 'center').toUpperCase()
                  ],
                  spacing: { before: 200, after: 200 }
                })
              ],
              width: { 
                size: field.width <= 10 ? 99 : (field.width >= 90 ? 1 : 100 - field.width), 
                type: 'pct' 
              },
              borders: {
                top: { style: BorderStyle.NONE },
                bottom: { style: BorderStyle.SINGLE, size: 1 },
                left: { style: BorderStyle.NONE },
                right: { style: BorderStyle.NONE }
              }
            })
          ]
        })
      )
    })

    return new Table({
      width: { size: 100, type: 'pct' },
      borders: {
        top: { style: BorderStyle.NONE },
        bottom: { style: BorderStyle.NONE },
        left: { style: BorderStyle.NONE },
        right: { style: BorderStyle.NONE },
        insideHorizontal: { style: BorderStyle.NONE },
        insideVertical: { style: BorderStyle.NONE }
      },
      rows
    })
  }
  
  // 确保列宽数组存在且长度正确
  const columnWidths = component.columnWidths || []
  if (component.layout === 'key-value' && columnWidths.length !== 2) {
    // 如果是键值对布局但列宽数组不是2个元素，则重置为默认值
    component.columnWidths = [50, 50]
  }
  
  // 键值对布局或标题内容布局
  if (component.layout === 'key-value') {
    // 键值对布局 - 不添加表头行，直接添加数据行
    for (const field of component.fields) {
      const cellContent = await generateCellContent(field, component)
      
      rows.push(
        new TableRow({
          children: [
            new TableCell({
              children: [
                new Paragraph({
                  children: [
                    new TextRun({
                      text: field.label,
                      size: (component.settings?.fontSize || 16) * 2,
                      font: component.settings?.fontFamily || '宋体',
                      bold: true
                    })
                  ],
                  alignment: AlignmentType[
                    (component.settings?.cellStyle?.alignment || 'left').toUpperCase()
                  ]
                })
              ],
              width: { 
                size: component.columnWidths[0] <= 10 ? 1 : (component.columnWidths[0] >= 90 ? 99 : component.columnWidths[0]), 
                type: 'pct' 
              }
            }),
            new TableCell({
              children: cellContent,
              width: { 
                size: component.columnWidths[0] <= 10 ? 99 : (component.columnWidths[0] >= 90 ? 1 : component.columnWidths[1]), 
                type: 'pct' 
              }
            })
          ]
        })
      )
    }
  } else {
    // 标题内容布局 - 每个字段占一行
    for (const field of component.fields) {
      rows.push(
        new TableRow({
          children: [
            new TableCell({
              children: [
                new Paragraph({
                  children: [
                    new TextRun({
                      text: field.label,
                      size: (component.settings?.fontSize || 16) * 2,
                      font: component.settings?.fontFamily || '宋体',
                      bold: true
                    })
                  ]
                }),
                new Paragraph({
                  children: [
                    new TextRun({
                      text: field.value || '',
                      size: (component.settings?.fontSize || 16) * 2,
                      font: component.settings?.fontFamily || '宋体'
                    })
                  ]
                }),
                ...(field.description ? [
                  new Paragraph({
                    children: [
                      new TextRun({
                        text: field.description,
                        size: ((component.settings?.fontSize || 16) - 2) * 2,
                        font: component.settings?.fontFamily || '宋体',
                        color: '666666'
                      })
                    ]
                  })
                ] : [])
              ]
            })
          ]
        })
      )
    }
  }
  
  // 创建表格
  return new Table({
    width: { size: 100, type: 'pct' },
    borders: {
      top: { style: component.settings?.borderStyle?.top ? BorderStyle.SINGLE : BorderStyle.NONE, size: 1 },
      bottom: { style: component.settings?.borderStyle?.bottom ? BorderStyle.SINGLE : BorderStyle.NONE, size: 1 },
      left: { style: component.settings?.borderStyle?.left ? BorderStyle.SINGLE : BorderStyle.NONE, size: 1 },
      right: { style: component.settings?.borderStyle?.right ? BorderStyle.SINGLE : BorderStyle.NONE, size: 1 },
      insideHorizontal: { style: component.settings?.borderStyle?.insideHorizontal ? BorderStyle.SINGLE : BorderStyle.NONE, size: 1 },
      insideVertical: { style: component.settings?.borderStyle?.insideVertical ? BorderStyle.SINGLE : BorderStyle.NONE, size: 1 }
    },
    rows
  })
}

// 生成目录组件
export const generateTOCComponent = (component) => {
  const children = []
  
  // 添加标题
  if (component.title) {
    children.push(
      new Paragraph({
        text: component.title,
        heading: HeadingLevel.HEADING_1,
        alignment: AlignmentType.CENTER,
        spacing: { after: 200 },
        font: component.titleFontFamily || '黑体',
        size: (component.titleFontSize || 18) * 2
      })
    )
  }
  
  // 添加目录项
  component.items.forEach(item => {
    children.push(
      new Paragraph({
        children: [
          new TextRun({
            text: `${item.number} ${item.title}`,
            size: (component.fontSize || 16) * 2,
            font: component.fontFamily || '宋体'
          }),
          new TextRun({
            text: "\t",
          }),
          new TextRun({
            text: item.page,
            size: (component.fontSize || 16) * 2,
            font: component.fontFamily || '宋体'
          })
        ],
        tabStops: [
          {
            type: AlignmentType.RIGHT,
            position: convertInchesToTwip(5.5),
            leader: "dot"
          }
        ],
        spacing: { line: (component.lineHeight || 1.8) * 240 }
      })
    )
  })
  
  return new Table({
    width: { size: 100, type: 'pct' },
    borders: {
      top: { style: BorderStyle.NONE },
      bottom: { style: BorderStyle.NONE },
      left: { style: BorderStyle.NONE },
      right: { style: BorderStyle.NONE },
      insideHorizontal: { style: BorderStyle.NONE },
      insideVertical: { style: BorderStyle.NONE }
    },
    rows: [
      new TableRow({
        children: [
          new TableCell({
            children,
            borders: {
              top: { style: BorderStyle.NONE },
              bottom: { style: BorderStyle.NONE },
              left: { style: BorderStyle.NONE },
              right: { style: BorderStyle.NONE }
            }
          })
        ]
      })
    ]
  })
}

// 生成进度表组件
export const generateScheduleComponent = (component) => {
  const children = []
  
  // 添加每个工作项
  component.items.forEach((item, index) => {
    children.push(
      new Paragraph({
        children: [
          new TextRun({
            text: `（${index + 1}）${item.task}`,
            size: (component.settings?.fontSize || 16) * 2,
            font: component.settings?.fontFamily || '宋体'
          }),
          new TextRun({
            text: "\t",
          }),
          new TextRun({
            text: `${item.startDate} — ${item.endDate}`,
            size: (component.settings?.fontSize || 16) * 2,
            font: component.settings?.fontFamily || '宋体'
          })
        ],
        tabStops: [
          {
            type: AlignmentType.RIGHT,
            position: convertInchesToTwip(5.5)
          }
        ],
        spacing: { line: (component.settings?.lineHeight || 1.5) * 240 }
      })
    )
  })
  
  return new Table({
    width: { size: 100, type: 'pct' },
    borders: {
      top: { style: BorderStyle.NONE },
      bottom: { style: BorderStyle.NONE },
      left: { style: BorderStyle.NONE },
      right: { style: BorderStyle.NONE },
      insideHorizontal: { style: BorderStyle.NONE },
      insideVertical: { style: BorderStyle.NONE }
    },
    rows: [
      new TableRow({
        children: [
          new TableCell({
            children,
            borders: {
              top: { style: BorderStyle.NONE },
              bottom: { style: BorderStyle.NONE },
              left: { style: BorderStyle.NONE },
              right: { style: BorderStyle.NONE }
            }
          })
        ]
      })
    ]
  })
}

// 生成单元格内容
const generateCellContent = async (field, component) => {
  // 如果字段有内容组件，则生成组件内容
  if (field.contentComponents && field.contentComponents.length > 0) {
    const paragraphs = []
    
    // 处理每个内容组件
    for (const contentComp of field.contentComponents) {
      switch (contentComp.type) {
        case 'text':
          // 文本组件
          paragraphs.push(
            new Paragraph({
              children: [
                new TextRun({
                  text: contentComp.content || '',
                  size: (contentComp.fontSize || 16) * 2,
                  font: contentComp.fontFamily || '宋体'
                })
              ],
              alignment: AlignmentType[(contentComp.alignment || 'left').toUpperCase()],
              spacing: { line: (contentComp.lineHeight || 1.5) * 240 }
            })
          )
          break
          
        case 'list':
          // 列表组件
          if (contentComp.items && Array.isArray(contentComp.items)) {
            contentComp.items.forEach((item, index) => {
              const bulletOrNumber = contentComp.listType === 'bullet' 
                ? '• ' 
                : `${index + 1}. `
              
              paragraphs.push(
                new Paragraph({
                  children: [
                    new TextRun({
                      text: bulletOrNumber + item.content,
                      size: (contentComp.fontSize || 16) * 2,
                      font: contentComp.fontFamily || '宋体'
                    })
                  ],
                  alignment: AlignmentType[(contentComp.alignment || 'left').toUpperCase()],
                  indent: { left: 720 } // 约0.5英寸的缩进
                })
              )
            })
          }
          break
          
        case 'reference':
          // 参考文献组件
          if (contentComp.items && Array.isArray(contentComp.items)) {
            contentComp.items.forEach((item, index) => {
              let referenceText = ''
              
              // 根据不同的引用样式格式化文本
              switch (contentComp.style) {
                case 'apa':
                  referenceText = `${item.author} (${item.year}). ${item.title}. ${item.source}.`
                  break
                case 'mla':
                  referenceText = `${item.author}. "${item.title}." ${item.source}, ${item.year}.`
                  break
                case 'chicago':
                  referenceText = `${item.author}. ${item.title}. ${item.source}, ${item.year}.`
                  break
                case 'gb':
                default:
                  referenceText = `[${index + 1}] ${item.author}. ${item.title}[M]. ${item.source}, ${item.year}.`
                  break
              }
              
              paragraphs.push(
                new Paragraph({
                  children: [
                    new TextRun({
                      text: referenceText,
                      size: (contentComp.fontSize || 16) * 2,
                      font: contentComp.fontFamily || '宋体'
                    })
                  ],
                  alignment: AlignmentType.LEFT,
                  indent: { left: 720, hanging: 360 } // 悬挂缩进
                })
              )
            })
          }
          break
          
        case 'image':
          // 图片组件
          if (contentComp.url) {
            try {
              // 添加图片
              const imageBuffer = await fetchImageAndAddToDoc(contentComp.url)
              const dimensions = await getImageDimensions(contentComp.url)
              const scaledDimensions = calculateImageDimensions(
                dimensions.width, 
                dimensions.height, 
                contentComp.width || 80
              )
              
              paragraphs.push(
                new Paragraph({
                  children: [
                    new ImageRun({
                      data: imageBuffer,
                      transformation: {
                        width: scaledDimensions.width,
                        height: scaledDimensions.height
                      }
                    })
                  ],
                  alignment: AlignmentType[(contentComp.alignment || 'center').toUpperCase()]
                })
              )
              
              // 如果有图片说明，添加说明文字
              if (contentComp.caption) {
                paragraphs.push(
                  new Paragraph({
                    children: [
                      new TextRun({
                        text: contentComp.caption,
                        size: ((contentComp.fontSize || 16) - 2) * 2,
                        font: contentComp.fontFamily || '宋体'
                      })
                    ],
                    alignment: AlignmentType.CENTER
                  })
                )
              }
            } catch (error) {
              console.error('处理图片失败:', error)
              // 添加错误提示
              paragraphs.push(
                new Paragraph({
                  children: [
                    new TextRun({
                      text: '[图片加载失败]',
                      size: (contentComp.fontSize || 16) * 2,
                      font: contentComp.fontFamily || '宋体',
                      color: 'FF0000'
                    })
                  ],
                  alignment: AlignmentType.CENTER
                })
              )
            }
          }
          break
          
        case 'title':
          // 大标题组件
          paragraphs.push(
            new Paragraph({
              children: [
                new TextRun({
                  text: contentComp.content || '大标题',
                  size: (contentComp.fontSize || 24) * 2,
                  font: contentComp.fontFamily || '黑体',
                  bold: contentComp.bold !== false
                })
              ],
              alignment: AlignmentType[(contentComp.alignment || 'center').toUpperCase()],
              spacing: { before: 240, after: 120 }
            })
          )
          break
          
        case 'paragraph':
          // 普通文本段落组件
          paragraphs.push(
            new Paragraph({
              children: [
                new TextRun({
                  text: contentComp.content || '',
                  size: (contentComp.fontSize || 16) * 2,
                  font: contentComp.fontFamily || '宋体'
                })
              ],
              alignment: AlignmentType[(contentComp.alignment || 'left').toUpperCase()],
              spacing: { line: (contentComp.lineHeight || 1.5) * 240 }
            })
          )
          break
          
        case 'schedule':
          // 进度表组件
          // 添加标题（如果有）
          if (contentComp.title) {
            paragraphs.push(
              new Paragraph({
                children: [
                  new TextRun({
                    text: contentComp.title,
                    size: (contentComp.settings?.fontSize || 16) * 2,
                    font: contentComp.settings?.fontFamily || '宋体',
                    bold: true
                  })
                ],
                alignment: AlignmentType.CENTER,
                spacing: { before: 240, after: 120 }
              })
            )
          }
          
          // 添加每个工作项
          if (contentComp.items && Array.isArray(contentComp.items)) {
            contentComp.items.forEach((item, index) => {
              paragraphs.push(
                new Paragraph({
                  children: [
                    new TextRun({
                      text: `（${index + 1}）${item.task}`,
                      size: (contentComp.settings?.fontSize || 16) * 2,
                      font: contentComp.settings?.fontFamily || '宋体'
                    }),
                    new TextRun({
                      text: "\t",
                    }),
                    new TextRun({
                      text: `${item.startDate || ''} — ${item.endDate || ''}`,
                      size: (contentComp.settings?.fontSize || 16) * 2,
                      font: contentComp.settings?.fontFamily || '宋体'
                    })
                  ],
                  tabStops: [
                    {
                      type: AlignmentType.RIGHT,
                      position: convertInchesToTwip(3) // 表格内的进度表宽度较小
                    }
                  ],
                  spacing: { line: (contentComp.settings?.lineHeight || 1.5) * 240 }
                })
              )
            })
          }
          break
      }
    }
    
    // 如果没有生成任何段落，添加一个空段落
    if (paragraphs.length === 0) {
      paragraphs.push(
        new Paragraph({
          children: [
            new TextRun({
              text: field.value || '',
              size: (component.settings?.fontSize || 16) * 2,
              font: component.settings?.fontFamily || '宋体'
            })
          ],
          alignment: AlignmentType[
            (component.settings?.cellStyle?.alignment || 'left').toUpperCase()
          ]
        })
      )
    }
    
    return paragraphs
  }
  
  // 如果没有内容组件，返回默认段落
  return [
    new Paragraph({
      children: [
        new TextRun({
          text: field.value || '',
          size: (component.settings?.fontSize || 16) * 2,
          font: component.settings?.fontFamily || '宋体'
        })
      ],
      alignment: AlignmentType[
        (component.settings?.cellStyle?.alignment || 'left').toUpperCase()
      ]
    })
  ]
} 