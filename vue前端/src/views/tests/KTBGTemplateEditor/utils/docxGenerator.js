import { Document, Packer, Paragraph, TextRun, Table, TableRow, TableCell, BorderStyle, WidthType, AlignmentType, HeadingLevel, TabStopPosition, TabStopType } from 'docx';
import { saveAs } from 'file-saver';

/**
 * 将模板数据转换为Word文档
 * @param {Object} template 模板数据
 * @param {Object} formData 表单数据
 * @returns {Promise} 返回一个Promise，成功时返回生成的文档
 */
export async function generateDocx(template, formData) {
  try {
    // 创建文档
    const doc = new Document({
      sections: [
        {
          properties: {},
          children: generateDocumentContent(template, formData)
        }
      ]
    });

    // 生成并保存文档
    const buffer = await Packer.toBlob(doc);
    saveAs(buffer, `${template.name || '开题报告'}.docx`);
    return doc;
  } catch (error) {
    console.error('生成文档时出错:', error);
    throw error;
  }
}

/**
 * 生成文档内容
 * @param {Object} template 模板数据
 * @param {Object} formData 表单数据
 * @returns {Array} 文档内容元素数组
 */
function generateDocumentContent(template, formData) {
  const content = [];
  
  // 添加标题
  content.push(
    new Paragraph({
      text: template.name || '开题报告',
      heading: HeadingLevel.HEADING_1,
      alignment: AlignmentType.CENTER,
      spacing: { before: 200, after: 200 }
    })
  );
  
  // 处理每个部分
  template.sections.forEach(section => {
    // 添加部分标题
    content.push(
      new Paragraph({
        text: section.title,
        heading: HeadingLevel.HEADING_2,
        spacing: { before: 200, after: 100 }
      })
    );
    
    // 根据部分类型生成内容
    if (section.type === 'fields') {
      content.push(...generateFieldsContent(section, formData));
    } else if (section.type === 'schedule') {
      content.push(...generateScheduleContent(section, formData));
    } else if (section.type === 'references') {
      content.push(...generateReferencesContent(section, formData));
    }
  });
  
  return content;
}

/**
 * 生成普通字段内容
 * @param {Object} section 部分数据
 * @param {Object} formData 表单数据
 * @returns {Array} 内容元素数组
 */
function generateFieldsContent(section, formData) {
  const content = [];
  
  // 如果是封面信息，使用表格布局
  if (section.id === 'cover') {
    const table = new Table({
      width: { size: 100, type: WidthType.PERCENTAGE },
      borders: {
        top: { style: BorderStyle.NONE },
        bottom: { style: BorderStyle.NONE },
        left: { style: BorderStyle.NONE },
        right: { style: BorderStyle.NONE },
        insideHorizontal: { style: BorderStyle.NONE },
        insideVertical: { style: BorderStyle.NONE }
      },
      rows: []
    });
    
    // 添加每个字段
    section.fields.forEach(field => {
      const value = formData[field.id] || '';
      
      table.root.push(
        new TableRow({
          children: [
            new TableCell({
              width: { size: 30, type: WidthType.PERCENTAGE },
              children: [new Paragraph(field.label)]
            }),
            new TableCell({
              width: { size: 70, type: WidthType.PERCENTAGE },
              children: [new Paragraph(value)]
            })
          ]
        })
      );
    });
    
    content.push(table);
  } else {
    // 其他部分使用段落布局
    section.fields.forEach(field => {
      const value = formData[field.id] || '';
      
      if (field.type === 'textarea') {
        // 多行文本按段落分割
        const paragraphs = value.split('\n');
        paragraphs.forEach(paragraph => {
          content.push(
            new Paragraph({
              text: paragraph,
              spacing: { line: 360 } // 1.5倍行距
            })
          );
        });
      } else {
        content.push(
          new Paragraph({
            children: [
              new TextRun({ text: `${field.label}: `, bold: true }),
              new TextRun(value)
            ],
            spacing: { line: 360 } // 1.5倍行距
          })
        );
      }
    });
  }
  
  return content;
}

/**
 * 生成工作计划内容
 * @param {Object} section 部分数据
 * @param {Object} formData 表单数据
 * @returns {Array} 内容元素数组
 */
function generateScheduleContent(section, formData) {
  const content = [];
  const scheduleItems = formData.schedule || [];
  
  // 添加工作计划说明
  content.push(
    new Paragraph({
      text: '工作计划安排如下：',
      spacing: { line: 360 } // 1.5倍行距
    })
  );
  
  // 添加每个工作计划项
  scheduleItems.forEach((item, index) => {
    const dateRange = item.startDate && item.endDate 
      ? `${item.startDate} 至 ${item.endDate}`
      : '';
    
    content.push(
      new Paragraph({
        children: [
          new TextRun(`${index + 1}. ${item.task}`),
          new TextRun({ text: '\t' }),
          new TextRun(dateRange)
        ],
        spacing: { line: 360 }, // 1.5倍行距
        tabStops: [
          {
            type: TabStopType.RIGHT,
            position: TabStopPosition.MAX
          }
        ]
      })
    );
  });
  
  return content;
}

/**
 * 生成参考文献内容
 * @param {Object} section 部分数据
 * @param {Object} formData 表单数据
 * @returns {Array} 内容元素数组
 */
function generateReferencesContent(section, formData) {
  const content = [];
  const references = formData.references || [];
  
  // 添加每条参考文献
  references.forEach((reference, index) => {
    content.push(
      new Paragraph({
        text: `[${index + 1}] ${reference}`,
        spacing: { line: 360 } // 1.5倍行距
      })
    );
  });
  
  return content;
}

/**
 * 将英寸转换为twip（文档中的度量单位）
 * @param {Number} inches 英寸数
 * @returns {Number} twip数
 */
export function inchesToTwip(inches) {
  return Math.round(inches * 1440);
}

/**
 * 将厘米转换为twip
 * @param {Number} cm 厘米数
 * @returns {Number} twip数
 */
export function cmToTwip(cm) {
  return Math.round(cm * 566.9);
}

/**
 * 将毫米转换为twip
 * @param {Number} mm 毫米数
 * @returns {Number} twip数
 */
export function mmToTwip(mm) {
  return Math.round(mm * 56.69);
} 