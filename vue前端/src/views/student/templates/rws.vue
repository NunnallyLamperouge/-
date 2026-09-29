<script setup>
import { ref } from 'vue'
import { Document, Packer, Paragraph, Table, TableCell, TableRow, TextRun, HeadingLevel, AlignmentType, convertInchesToTwip, BorderStyle } from 'docx'

const formData = ref({
  // 基本信息
  title: '', // 中文题目
  titleEn: '', // 英文题目
  studentName: '', // 学生姓名
  studentId: '', // 学号
  major: '', // 专业班级
  college: '计算机与软件工程学院', // 学院
  supervisor: '', // 指导老师
  departmentHead: '', // 教研室负责人
  submitDate: '', // 下达日期
  
  // 研究内容
  researchBackground: '', // 研究背景及意义
  researchContent: '', // 研究内容
  
  // 工作进度
  schedule: [
    { task: '完成选题', startDate: '2021-11-15', endDate: '2021-11-26' },
    { task: '完成任务书下达', startDate: '2021-11-27', endDate: '2021-12-15' },
    { task: '完成开题报告', startDate: '2021-12-16', endDate: '2022-01-24' },
    { task: '完成中期检查', startDate: '2022-01-25', endDate: '2022-02-25' },
    { task: '完成论文初稿', startDate: '2022-02-26', endDate: '2022-03-31' },
    { task: '完成论文定稿', startDate: '2022-04-01', endDate: '2022-04-30' },
    { task: '完成相似比检测', startDate: '2022-04-15', endDate: '2022-05-15' },
    { task: '完成答辩', startDate: '2022-05-14', endDate: '2022-05-22' },
    { task: '材料归档', startDate: '2022-05-25', endDate: '2022-06-03' }
  ],
  
  // 参考文献
  references: []
})

// 添加参考文献
const addReference = () => {
  formData.value.references.push('')
}

// 删除参考文献
const removeReference = (index) => {
  formData.value.references.splice(index, 1)
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
              top: convertInchesToTwip(1), // 2.5cm
              bottom: convertInchesToTwip(1), // 2.5cm
              left: convertInchesToTwip(1.2), // 3cm
              right: convertInchesToTwip(1.2), // 3cm
            }
          }
        },
        children: [
          // 标题
          new Paragraph({
            children: [
              new TextRun({
                text: "本科毕业设计（论文）任务书",
                font: {
                  name: '黑体',
                  size: 32 // 二号
                }
              })
            ],
            heading: HeadingLevel.TITLE,
            alignment: AlignmentType.CENTER
          }),
          
          // 封面部分 - 添加空行作为间隔
          new Paragraph({
            spacing: { after: 400 }
          }),
          
          // 封面表格 - 优化版本
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
            rows: [
              // 题目（中文）
              new TableRow({
                children: [
                  new TableCell({
                    children: [new Paragraph({
                      children: [
                        new TextRun({
                          text: "题目（中文）",
                          size: 16 * 2,
                          bold: true,
                          font: "宋体"
                        })
                      ],
                      alignment: AlignmentType.LEFT,
                      spacing: { before: 200, after: 200 }
                    })],
                    width: { size: 30, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER,
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
                            text: formData.value.title,
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
              }),
              
              // 题目（英文）
              new TableRow({
                children: [
                  new TableCell({
                    children: [new Paragraph({
                      children: [
                        new TextRun({
                          text: "题目（英文）",
                          size: 16 * 2,
                          bold: true,
                          font: "宋体"
                        })
                      ],
                      alignment: AlignmentType.LEFT,
                      spacing: { before: 200, after: 200 }
                    })],
                    width: { size: 30, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER,
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
                            text: formData.value.titleEn,
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
              }),
              
              // 学生姓名
              new TableRow({
                children: [
                  new TableCell({
                    children: [new Paragraph({
                      children: [
                        new TextRun({
                          text: "学生姓名",
                          size: 16 * 2,
                          bold: true,
                          font: "宋体"
                        })
                      ],
                      alignment: AlignmentType.LEFT,
                      spacing: { before: 200, after: 200 }
                    })],
                    width: { size: 30, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER,
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
                            text: formData.value.studentName,
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
              }),
              
              // 学号
              new TableRow({
                children: [
                  new TableCell({
                    children: [new Paragraph({
                      children: [
                        new TextRun({
                          text: "学号",
                          size: 16 * 2,
                          bold: true,
                          font: "宋体"
                        })
                      ],
                      alignment: AlignmentType.LEFT,
                      spacing: { before: 200, after: 200 }
                    })],
                    width: { size: 30, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER,
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
                            text: formData.value.studentId,
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
              }),
              
              // 专业班级
              new TableRow({
                children: [
                  new TableCell({
                    children: [new Paragraph({
                      children: [
                        new TextRun({
                          text: "专业班级",
                          size: 16 * 2,
                          bold: true,
                          font: "宋体"
                        })
                      ],
                      alignment: AlignmentType.LEFT,
                      spacing: { before: 200, after: 200 }
                    })],
                    width: { size: 30, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER,
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
                            text: formData.value.major,
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
              }),
              
              // 学院
              new TableRow({
                children: [
                  new TableCell({
                    children: [new Paragraph({
                      children: [
                        new TextRun({
                          text: "学院",
                          size: 16 * 2,
                          bold: true,
                          font: "宋体"
                        })
                      ],
                      alignment: AlignmentType.LEFT,
                      spacing: { before: 200, after: 200 }
                    })],
                    width: { size: 30, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER,
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
                            text: formData.value.college,
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
              }),
              
              // 指导老师（手签）
              new TableRow({
                children: [
                  new TableCell({
                    children: [new Paragraph({
                      children: [
                        new TextRun({
                          text: "指导老师",
                          size: 16 * 2,
                          bold: true,
                          font: "宋体"
                        })
                      ],
                      alignment: AlignmentType.LEFT,
                      spacing: { before: 200, after: 200 }
                    })],
                    width: { size: 30, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER,
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
                            text: " ",
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
              }),
              
              // 教研室负责人（手签）
              new TableRow({
                children: [
                  new TableCell({
                    children: [new Paragraph({
                      children: [
                        new TextRun({
                          text: "教研室负责人",
                          size: 16 * 2,
                          bold: true,
                          font: "宋体"
                        })
                      ],
                      alignment: AlignmentType.LEFT,
                      spacing: { before: 200, after: 200 }
                    })],
                    width: { size: 30, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER,
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
                            text: " ",
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
              }),
              
              // 下达日期
              new TableRow({
                children: [
                  new TableCell({
                    children: [new Paragraph({
                      children: [
                        new TextRun({
                          text: "下达日期",
                          size: 16 * 2,
                          bold: true,
                          font: "宋体"
                        })
                      ],
                      alignment: AlignmentType.LEFT,
                      spacing: { before: 200, after: 200 }
                    })],
                    width: { size: 30, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER,
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
                            text: formData.value.submitDate ? formData.value.submitDate : "（此处填任务书的下达日期）",
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
            ]
          }),
          
          // 添加空行作为间隔
          new Paragraph({
            spacing: { after: 400 }
          }),
          
          // 主表格
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
            rows: [
              // 题目行
              new TableRow({
                children: [
                  new TableCell({
                    children: [new Paragraph({
                      children: [
                        new TextRun({
                          text: "题目",
                          font: {
                            name: '宋体',
                            size: 24 // 小四号
                          },
                          bold: true
                        })
                      ],
                      alignment: AlignmentType.CENTER
                    })],
                    width: { size: 5, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER
                  }),
                  new TableCell({
                    children: [new Paragraph({
                      children: [
                        new TextRun({
                          text: formData.value.title,
                          font: {
                            name: '宋体',
                            size: 24
                          }
                        })
                      ],
                      alignment: AlignmentType.CENTER
                    })]
                  })
                ]
              }),
              
              // 论文研究背景及意义、内容、成果形式
              new TableRow({
                children: [
                  new TableCell({
                    children: [
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "论",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "文",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "研",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "究",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "背",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "景",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      })
                    ],
                    width: { size: 5, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER
                  }),
                  new TableCell({
                    children: [
                      // 研究背景及意义
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "研究背景及意义",
                            font: {
                              name: '宋体',
                              size: 24,
                              bold: true
                            }
                          })
                        ],
                        spacing: { before: 200, after: 200 }
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: formData.value.researchBackground,
                            font: {
                              name: '宋体',
                              size: 24
                            }
                          })
                        ],
                        indent: { firstLine: 480 }, // 首行缩进2字符
                        spacing: { line: 360 } // 1.5倍行距
                      }),
                      
                      // 研究内容
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "研究内容",
                            font: {
                              name: '宋体',
                              size: 24,
                              bold: true
                            }
                          })
                        ],
                        spacing: { before: 200, after: 200 }
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: formData.value.researchContent,
                            font: {
                              name: '宋体',
                              size: 24
                            }
                          })
                        ],
                        indent: { firstLine: 480 },
                        spacing: { line: 360 }
                      }),
                      
                      // 成果形式
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "成果形式",
                            font: {
                              name: '宋体',
                              size: 24,
                              bold: true
                            }
                          })
                        ],
                        spacing: { before: 200, after: 200 }
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "毕业设计（论文）正文",
                            font: {
                              name: '宋体',
                              size: 24
                            }
                          })
                        ],
                        spacing: { line: 360 }
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "程序源代码和可执行程序（或视频）",
                            font: {
                              name: '宋体',
                              size: 24
                            }
                          })
                        ],
                        spacing: { line: 360 }
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "程序设计相关文档",
                            font: {
                              name: '宋体',
                              size: 24
                            }
                          })
                        ],
                        spacing: { line: 360 }
                      })
                    ]
                  })
                ]
              }),
              
              // 工作进度
              new TableRow({
                children: [
                  new TableCell({
                    children: [
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "工",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "作",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "进",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "度",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      })
                    ],
                    width: { size: 5, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER
                  }),
                  new TableCell({
                    children: formData.value.schedule.map((item, index) => 
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: `（${index + 1}）${item.task}    ${item.startDate} — ${item.endDate}`,
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
                ]
              }),
              
              // 参考文献
              new TableRow({
                children: [
                  new TableCell({
                    children: [
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "参",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "考",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "文",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      }),
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: "献",
                            font: {
                              name: '宋体',
                              size: 24
                            },
                            bold: true
                          })
                        ],
                        alignment: AlignmentType.CENTER
                      })
                    ],
                    width: { size: 5, type: 'pct' },
                    verticalAlign: AlignmentType.CENTER
                  }),
                  new TableCell({
                    children: formData.value.references.map((ref, index) => 
                      new Paragraph({
                        children: [
                          new TextRun({
                            text: `[${index + 1}] `,
                            font: {
                              name: '宋体',
                              size: 24
                            }
                          }),
                          new TextRun({
                            text: ref,
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
                ]
              })
            ]
          })
        ]
      }]
    });

    // 生成并下载文档
    const blob = await Packer.toBlob(doc);
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = '毕业设计任务书.docx';
    link.click();
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error('文档生成错误:', error);
    alert('文档生成失败，请查看控制台获取详细错误信息');
  }
}

// 字数验证
const validateLength = (text, minLength) => {
  return text.length >= minLength
}
</script>

<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-6">毕业设计任务书生成器</h1>
    
    <div class="space-y-4">
      <!-- 基本信息 -->
      <div class="grid grid-cols-2 gap-4">
        <div class="form-group">
          <label>题目</label>
          <input v-model="formData.title" class="form-input" required />
        </div>
        
        <div class="form-group">
          <label>英文题目</label>
          <input v-model="formData.titleEn" class="form-input" required />
        </div>
        
        <div class="form-group">
          <label>学生姓名</label>
          <input v-model="formData.studentName" class="form-input" required />
        </div>
        
        <div class="form-group">
          <label>学号</label>
          <input v-model="formData.studentId" class="form-input" required />
        </div>
        
        <div class="form-group">
          <label>专业班级</label>
          <input v-model="formData.major" class="form-input" required />
        </div>
        
        <div class="form-group">
          <label>学院</label>
          <input v-model="formData.college" class="form-input" disabled />
        </div>
        
        <div class="form-group">
          <label>指导老师（手签）</label>
          <input v-model="formData.supervisor" class="form-input" required />
        </div>
        
        <div class="form-group">
          <label>教研室负责人（手签）</label>
          <input v-model="formData.departmentHead" class="form-input" required />
        </div>
        
        <div class="form-group">
          <label>下达日期</label>
          <input type="date" v-model="formData.submitDate" class="form-input" required />
        </div>
      </div>

      <!-- 研究内容 -->
      <div class="form-group">
        <label>研究背景及意义（不少于200字）</label>
        <textarea 
          v-model="formData.researchBackground" 
          class="form-textarea"
          rows="6"
          required
          :class="{ 'border-red-500': !validateLength(formData.researchBackground, 200) }"
        ></textarea>
        <span class="text-sm text-gray-500">当前字数: {{ formData.researchBackground.length }}</span>
      </div>

      <div class="form-group">
        <label>研究内容（不少于300字）</label>
        <textarea 
          v-model="formData.researchContent" 
          class="form-textarea"
          rows="8"
          required
          :class="{ 'border-red-500': !validateLength(formData.researchContent, 300) }"
        ></textarea>
        <span class="text-sm text-gray-500">当前字数: {{ formData.researchContent.length }}</span>
      </div>

      <!-- 工作进度 -->
      <div class="form-group">
        <label>工作进度</label>
        <div class="border rounded p-4">
          <div v-for="(item, index) in formData.schedule" :key="index" class="grid grid-cols-3 gap-4 mb-2">
            <input v-model="item.task" class="form-input" :placeholder="`任务${index + 1}`" />
            <input type="date" v-model="item.startDate" class="form-input" />
            <input type="date" v-model="item.endDate" class="form-input" />
          </div>
        </div>
      </div>

      <!-- 参考文献 -->
      <div class="form-group">
        <label class="flex justify-between">
          <span>参考文献（至少3条）</span>
          <button 
            @click="addReference"
            class="px-2 py-1 bg-blue-500 text-white rounded"
          >
            添加文献
          </button>
        </label>
        
        <div v-for="(ref, index) in formData.references" :key="index" class="flex gap-2 mt-2">
          <input 
            v-model="formData.references[index]"
            class="form-input flex-1"
            :placeholder="`参考文献 ${index + 1}`"
            required
          />
          <button 
            @click="() => removeReference(index)"
            class="px-2 py-1 bg-red-500 text-white rounded"
          >
            删除
          </button>
        </div>
        <span class="text-sm text-gray-500" :class="{ 'text-red-500': formData.references.length < 3 }">
          当前文献数: {{ formData.references.length }}/3
        </span>
      </div>

      <!-- 生成按钮 -->
      <div class="text-center">
        <button
          @click="generateDoc"
          class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
          :disabled="
            !formData.title ||
            !formData.titleEn ||
            !formData.studentName ||
            !formData.studentId ||
            !formData.major ||
            !formData.supervisor ||
            !formData.departmentHead ||
            !formData.submitDate ||
            !validateLength(formData.researchBackground, 200) ||
            !validateLength(formData.researchContent, 300) ||
            formData.references.length < 3
          "
        >
          生成文档
        </button>
      </div>
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
