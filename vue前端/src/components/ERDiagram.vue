<template>
  <div class="er-diagram-container">
    <div ref="chartContainer" class="w-full h-full"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick, computed } from 'vue';
import * as echarts from 'echarts';
import { useThemeStore } from '@/stores/theme';

const props = defineProps({
  data: {
    type: Object,
    required: true,
    default: () => ({ tables: [], relationships: [] })
  },
  theme: {
    type: String,
    default: 'light'
  },
  chartTitle: {
    type: String,
    default: 'ER图'
  }
});

const themeStore = useThemeStore();
const isDark = computed(() => themeStore.isDark);

const chartContainer = ref(null);
const chart = ref(null);

// 初始化图表
const initChart = async () => {
  if (!chartContainer.value) return;

  // 确保DOM已经渲染完成并且有尺寸
  await nextTick();
  
  // 如果容器尺寸为0，使用setTimeout等待DOM完全渲染
  if (chartContainer.value.clientWidth === 0 || chartContainer.value.clientHeight === 0) {
    setTimeout(() => {
      initChart();
    }, 100);
    return;
  }

  if (chart.value) {
    chart.value.dispose();
  }

  // 创建图表实例
  chart.value = echarts.init(chartContainer.value, isDark.value ? 'dark' : null);
  
  // 设置图表配置
  updateChart();
  
  // 监听窗口大小变化
  window.addEventListener('resize', handleResize);
};

// 更新图表数据
const updateChart = () => {
  if (!chart.value || !props.data) return;
  
  const { tables, relationships } = props.data;
  
  if (tables.length === 0) return;
  
  // 准备节点数据
  const nodes = [];
  
  // 添加实体节点（使用椭圆形状）
  tables.forEach((table) => {
    nodes.push({
      id: table.name,
      name: table.name,
      value: table.fields.length,
      symbolSize: [120, 80], // 矩形形状
      symbol: 'rect', // 使用矩形表示实体
      category: 0, // 实体类别
      itemStyle: {
        color: '#ffffff', // 白色背景
        borderColor: '#000000', // 黑色边框
        borderWidth: 2 // 加粗边框以符合ER规范
      },
      label: {
        show: true,
        position: 'inside',
        formatter: table.name,
        fontSize: 14, // 加大字体
        fontWeight: 'bold',
        color: '#000000', // 黑色文字
        textBorderColor: 'rgba(255,255,255,0.8)',
        textBorderWidth: 2
      },
      tooltip: {
        formatter: formatTableTooltip(table)
      }
    });
    
    // 为每个主键添加一个属性节点（使用椭圆形状，但更小）
    table.fields.filter(field => field.isPrimary).forEach(field => {
      const fieldNodeId = `${table.name}_${field.name}`;
      nodes.push({
        id: fieldNodeId,
        name: field.name,
        symbolSize: [80, 40], // 椭圆形状
        symbol: 'circle', // 属性使用椭圆形
        category: 1, // 属性类别
        itemStyle: {
          color: '#ffffff', // 白色背景
          borderColor: '#000000', // 黑色边框
          borderWidth: 1.5 // 略微加粗边框
        },
        label: {
          show: true,
          position: 'inside',
          formatter: field.name,
          fontSize: 12, // 调整字体大小
          color: '#000000', // 黑色文字
          textBorderColor: 'rgba(255,255,255,0.8)',
          textBorderWidth: 1
        },
        tooltip: {
          formatter: `<div><b>${field.name}</b></div><div>类型: ${field.type}</div><div>主键</div>`
        }
      });
    });
  });
  
  // 检查关系中的源和目标是否存在于表中
  const validRelationships = relationships.filter(rel => {
    const sourceExists = tables.some(table => table.name === rel.source);
    const targetExists = tables.some(table => table.name === rel.target);
    return sourceExists && targetExists;
  });
  
  // 为每个关系添加一个关系节点（使用菱形形状）
  validRelationships.forEach((rel, index) => {
    const relationNodeId = `relation_${index}`;
    nodes.push({
      id: relationNodeId,
      name: formatRelationshipLabel(rel),
      symbolSize: [60, 60],
      symbol: 'diamond', // 菱形形状表示关系
      category: 2, // 关系类别
      itemStyle: {
        color: '#ffffff', // 白色背景
        borderColor: '#000000', // 黑色边框
        borderWidth: 1.5 // 略微加粗边框
      },
      label: {
        show: true,
        position: 'inside',
        formatter: formatRelationshipLabel(rel),
        fontSize: 12, // 调整字体大小
        color: '#000000', // 黑色文字
        textBorderColor: 'rgba(255,255,255,0.8)',
        textBorderWidth: 1
      },
      tooltip: {
        formatter: formatRelationshipTooltip(rel)
      }
    });
  });
  
  // 准备边数据
  const edges = [];
  
  // 添加实体到属性的连接
  tables.forEach(table => {
    table.fields.filter(field => field.isPrimary).forEach(field => {
      const fieldNodeId = `${table.name}_${field.name}`;
      edges.push({
        source: table.name,
        target: fieldNodeId,
        lineStyle: {
          width: 1,
          type: 'solid',
          color: '#000000',
          curveness: 0 // 改为直线
        },
        tooltip: {
          formatter: `${table.name} 包含属性 ${field.name}`
        }
      });
    });
  });
  
  // 添加实体到关系的连接
  validRelationships.forEach((rel, index) => {
    const relationNodeId = `relation_${index}`;
    
    // 源实体到关系的连接
    edges.push({
      source: rel.source,
      target: relationNodeId,
      lineStyle: {
        width: 2,
        type: 'solid',
        color: '#000000',
        curveness: 0 // 改为直线
      },
      label: {
        show: true,
        formatter: rel.type.startsWith('one') ? '1' : 'N',
        fontSize: 12,
        fontWeight: 'bold'
      }
    });
    
    // 关系到目标实体的连接
    edges.push({
      source: relationNodeId,
      target: rel.target,
      lineStyle: {
        width: 2,
        type: 'solid',
        color: '#000000',
        curveness: 0 // 改为直线
      },
      label: {
        show: true,
        formatter: rel.type.endsWith('many') ? 'N' : '1',
        fontSize: 12,
        fontWeight: 'bold'
      }
    });
  });
  
  // 设置图表选项
  const option = {
    title: {
      text: props.chartTitle,
      subtext: `共 ${tables.length} 张表, ${validRelationships.length} 个关系`,
      left: 'center',
      textStyle: {
        color: isDark.value ? '#e5e7eb' : '#1f2937',
        textBorderColor: isDark.value ? 'rgba(0,0,0,0.3)' : 'rgba(255,255,255,0.3)',
        textBorderWidth: 1
      },
      subtextStyle: {
        color: isDark.value ? '#9ca3af' : '#6b7280'
      }
    },
    tooltip: {
      trigger: 'item',
      backgroundColor: isDark.value ? 'rgba(30,30,30,0.9)' : 'rgba(255,255,255,0.9)',
      borderColor: isDark.value ? '#4b5563' : '#e5e7eb',
      borderWidth: 1,
      textStyle: {
        color: isDark.value ? '#e5e7eb' : '#1f2937'
      },
      confine: true,
      enterable: true,
      extraCssText: 'max-width: 300px; white-space: normal; box-shadow: 0 4px 12px rgba(0,0,0,0.15);'
    },
    legend: {
      data: ['实体', '属性', '关系'],
      orient: 'vertical',
      right: 10,
      top: 20,
      textStyle: {
        color: isDark.value ? '#e5e7eb' : '#1f2937'
      },
      itemStyle: {
        borderColor: isDark.value ? '#4b5563' : '#e5e7eb'
      }
    },
    animationDurationUpdate: 1500,
    animationEasingUpdate: 'quinticInOut',
    series: [
      {
        type: 'graph',
        layout: 'force',
        data: nodes,
        links: edges,
        categories: [
          { name: '实体' },
          { name: '属性' },
          { name: '关系' }
        ],
        roam: true,
        draggable: true,
        focusNodeAdjacency: true,
        force: {
          repulsion: 2000, // 增加排斥力
          edgeLength: [150, 250], // 增加边长
          gravity: 0.05, // 减小引力
          layoutAnimation: true
        },
        edgeSymbol: ['none', 'none'],
        lineStyle: {
          color: '#000000',
          curveness: 0, // 改为直线
          width: 1.5,
          shadowColor: 'rgba(0,0,0,0.1)',
          shadowBlur: 2
        },
        emphasis: {
          focus: 'adjacency',
          lineStyle: {
            width: 4,
            color: isDark.value ? '#818cf8' : '#6366f1'
          }
        },
        // 确保图表能够适应容器大小
        left: '5%',
        right: '5%',
        top: '15%',
        bottom: '5%'
      }
    ]
  };
  
  try {
    chart.value.setOption(option, true);
    
    // 自动调整布局
    setTimeout(() => {
      if (chart.value) {
        // 先缩小再放大，触发布局重新计算
        chart.value.dispatchAction({
          type: 'graphRoam',
          zoom: 0.8
        });
        
        setTimeout(() => {
          chart.value.dispatchAction({
            type: 'graphRoam',
            zoom: 1.2
          });
          
          // 最后重置到正常大小
          setTimeout(() => {
            chart.value.dispatchAction({
              type: 'graphRoam',
              zoom: 1
            });
          }, 100);
        }, 100);
      }
    }, 200);
  } catch (error) {
    console.error('设置图表选项时出错:', error);
  }
};

// 获取表的主色调
const getPrimaryColor = (table) => {
  // 根据表的类型或特性设置不同的颜色
  if (table.fields.some(f => f.isPrimary)) {
    return '#5470c6'; // 有主键的表
  } else if (table.foreignKeys && table.foreignKeys.length > 0) {
    return '#91cc75'; // 有外键的表
  } else if (table.fields.length > 10) {
    return '#fac858'; // 字段较多的表
  } else {
    return '#ee6666'; // 其他表
  }
};

// 格式化表的提示信息
const formatTableTooltip = (table) => {
  const borderColor = isDark.value ? '#4b5563' : '#e5e7eb';
  const headerBgColor = isDark.value ? '#374151' : '#f3f4f6';
  const rowBorderColor = isDark.value ? '#4b5563' : '#f3f4f6';
  const textColor = isDark.value ? '#e5e7eb' : '#1f2937';
  const headerTextColor = isDark.value ? '#f9fafb' : '#111827';
  
  let tooltip = `<div style="font-weight:bold;margin-bottom:8px;font-size:14px;color:${headerTextColor};">${table.name}</div>`;
  tooltip += `<table style="width:100%;border-collapse:collapse;border-radius:4px;overflow:hidden;">`;
  tooltip += `<tr style="background-color:${headerBgColor};">
    <th style="text-align:left;border-bottom:1px solid ${borderColor};padding:4px 8px;color:${headerTextColor};">字段</th>
    <th style="text-align:left;border-bottom:1px solid ${borderColor};padding:4px 8px;color:${headerTextColor};">类型</th>
    <th style="text-align:left;border-bottom:1px solid ${borderColor};padding:4px 8px;color:${headerTextColor};">属性</th>
  </tr>`;
  
  table.fields.forEach(field => {
    const attributes = [];
    if (field.isPrimary) attributes.push('<span style="color:#10b981;">PK</span>');
    if (field.isForeign) attributes.push('<span style="color:#f59e0b;">FK</span>');
    if (field.isNotNull) attributes.push('NOT NULL');
    if (field.isAutoIncrement) attributes.push('AUTO');
    
    tooltip += `<tr>
      <td style="text-align:left;padding:4px 8px;border-bottom:1px solid ${rowBorderColor};color:${textColor};">${field.name}</td>
      <td style="text-align:left;padding:4px 8px;border-bottom:1px solid ${rowBorderColor};color:${textColor};">${field.type}</td>
      <td style="text-align:left;padding:4px 8px;border-bottom:1px solid ${rowBorderColor};color:${textColor};">${attributes.join(', ')}</td>
    </tr>`;
  });
  
  tooltip += '</table>';
  
  // 添加外键信息
  if (table.foreignKeys && table.foreignKeys.length > 0) {
    tooltip += `<div style="margin-top:12px;font-weight:bold;color:${headerTextColor};">外键关系:</div>`;
    tooltip += `<ul style="margin:5px 0;padding-left:20px;color:${textColor};">`;
    table.foreignKeys.forEach(fk => {
      tooltip += `<li>${fk.field} → <span style="color:#f59e0b;">${fk.referencedTable}.${fk.referencedField}</span></li>`;
    });
    tooltip += '</ul>';
  }
  
  return tooltip;
};

// 格式化关系的标签
const formatRelationshipLabel = (rel) => {
  switch (rel.type) {
    case 'one-to-one':
      return '1:1';
    case 'one-to-many':
      return '1:N';
    case 'many-to-many':
      return 'N:M';
    default:
      return '';
  }
};

// 获取关系的颜色
const getRelationshipColor = (type) => {
  switch (type) {
    case 'one-to-one':
      return '#5470c6';
    case 'one-to-many':
      return '#91cc75';
    case 'many-to-many':
      return '#fac858';
    default:
      return '#ee6666';
  }
};

// 格式化关系的提示信息
const formatRelationshipTooltip = (rel) => {
  let tooltip = `<div style="font-weight:bold;margin-bottom:5px;">关系: ${rel.source} → ${rel.target}</div>`;
  tooltip += `<div>类型: ${getRelationshipTypeName(rel.type)}</div>`;
  tooltip += `<div>字段: ${rel.sourceField} → ${rel.targetField}</div>`;
  
  if (rel.throughTable) {
    tooltip += `<div>通过中间表: ${rel.throughTable}</div>`;
  }
  
  return tooltip;
};

// 获取关系类型的中文名称
const getRelationshipTypeName = (type) => {
  switch (type) {
    case 'one-to-one':
      return '一对一';
    case 'one-to-many':
      return '一对多';
    case 'many-to-many':
      return '多对多';
    default:
      return '未知';
  }
};

// 处理窗口大小变化
const handleResize = () => {
  if (chart.value) {
    chart.value.resize();
  }
};

// 缩放控制
const zoomIn = () => {
  if (chart.value) {
    const option = chart.value.getOption();
    const zoom = option.series[0].zoom || 1;
    option.series[0].zoom = zoom * 1.2;
    chart.value.setOption(option);
  }
};

const zoomOut = () => {
  if (chart.value) {
    const option = chart.value.getOption();
    const zoom = option.series[0].zoom || 1;
    option.series[0].zoom = zoom * 0.8;
    chart.value.setOption(option);
  }
};

const resetZoom = () => {
  if (chart.value) {
    try {
      // 重新初始化图表
      updateChart();
    } catch (error) {
      console.error('重置缩放时出错:', error);
    }
  }
};

// 导出图表为图片
const exportChart = (type = 'png', pixelRatio = 2) => {
  if (!chart.value) return '';
  
  return chart.value.getDataURL({
    type,
    pixelRatio,
    backgroundColor: 'transparent'
  });
};

// 暴露方法给父组件
defineExpose({
  zoomIn,
  zoomOut,
  resetZoom,
  exportChart,
  updateChart
});

// 监听数据变化
watch(() => props.data, (newData) => {
  if (newData.tables.length > 0) {
    if (chart.value) {
      updateChart();
    } else {
      initChart();
    }
  }
}, { deep: true });

// 监听主题变化
watch(() => props.theme, () => {
  if (chart.value) {
    chart.value.dispose();
    initChart();
  }
});

// 监听主题变化
watch([() => props.theme, () => isDark.value], () => {
  if (chart.value) {
    chart.value.dispose();
    initChart();
  }
});

// 监听标题变化
watch(() => props.chartTitle, () => {
  if (chart.value) {
    updateChart();
  }
});

// 生命周期钩子
onMounted(() => {
  // 使用nextTick确保DOM已经渲染
  nextTick(() => {
    initChart();
  });
});

onUnmounted(() => {
  if (chart.value) {
    chart.value.dispose();
    chart.value = null;
  }
  window.removeEventListener('resize', handleResize);
});
</script>

<style scoped>
.er-diagram-container {
  width: 100%;
  height: 100%;
  min-height: 700px;
  position: relative;
}

:deep(.echarts-tooltip) {
  max-width: 300px;
  overflow-y: auto;
  max-height: 300px;
}
</style> 