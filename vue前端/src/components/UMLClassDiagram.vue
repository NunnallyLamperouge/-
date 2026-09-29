<template>
  <div class="uml-diagram-container">
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
    default: 'UML类图'
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

  // 创建图表实例，根据当前主题设置
  chart.value = echarts.init(chartContainer.value, isDark.value ? 'dark' : null, {
    renderer: 'canvas',  // 使用canvas渲染器以获得更好的性能
    devicePixelRatio: window.devicePixelRatio  // 适配高分辨率屏幕
  });
  
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
  
  // 准备边数据
  const edges = [];
  
  // 计算节点的初始位置
  const nodePositions = calculateInitialPositions(tables);
  
  // 添加类节点
  tables.forEach((table, index) => {
    // 分类字段
    const primaryFields = table.fields.filter(field => field.isPrimary);
    const foreignFields = table.fields.filter(field => field.isForeign && !field.isPrimary);
    const normalFields = table.fields.filter(field => !field.isPrimary && !field.isForeign);
    
    // 计算节点高度 (标题 + 属性 + 方法)
    // 根据字段数量动态计算节点高度，并设置最小高度
    const totalFields = table.fields.length;
    const nodeHeight = Math.max(100, 70 + totalFields * 22);
    
    // 根据表名和字段名长度调整宽度
    const maxFieldNameLength = table.fields.reduce((max, field) => 
      Math.max(max, field.name.length + field.type.length), 0);
    const nodeWidth = Math.max(200, Math.min(350, 120 + maxFieldNameLength * 8));
    
    // 添加UML类节点
    nodes.push({
      id: table.name,
      name: table.name,
      value: [nodeWidth, nodeHeight],
      x: nodePositions[index].x,
      y: nodePositions[index].y,
      symbolSize: [nodeWidth, nodeHeight],
      symbol: 'rect',
      itemStyle: {
        color: '#ffffff',
        borderColor: '#000000',
        borderWidth: 2
      },
      label: {
        show: true,
        position: 'inside',
        formatter: (params) => {
          // 类名
          let html = `{title|${table.name}}\n`;
          
          // 属性
          html += '{hr|}\n';
          
          // 主键
          primaryFields.forEach(field => {
            html += `{primary|+ ${field.name}: ${field.type}}\n`;
          });
          
          // 外键
          foreignFields.forEach(field => {
            html += `{foreign|+ ${field.name}: ${field.type}}\n`;
          });
          
          // 普通字段
          normalFields.forEach(field => {
            html += `{normal|+ ${field.name}: ${field.type}}\n`;
          });
          
          return html;
        },
        rich: {
          title: {
            fontSize: 14,
            fontWeight: 'bold',
            color: '#000000',
            padding: [8, 10, 5, 10],
            align: 'center'
          },
          hr: {
            width: '100%',
            borderWidth: 1,
            height: 0,
            borderColor: '#000000',
            padding: [0, 10]
          },
          primary: {
            fontSize: 12,
            color: '#be185d',
            align: 'left',
            padding: [3, 10]
          },
          foreign: {
            fontSize: 12,
            color: '#7c3aed',
            align: 'left',
            padding: [3, 10]
          },
          normal: {
            fontSize: 12,
            color: '#000000',
            align: 'left',
            padding: [3, 10]
          }
        }
      },
      tooltip: {
        formatter: (params) => {
          let html = `<div style="font-weight:bold">${table.name}</div>`;
          html += '<div style="margin-top:5px">属性:</div>';
          html += '<ul style="margin:5px 0;padding-left:20px">';
          
          table.fields.forEach(field => {
            const isPrimary = field.isPrimary ? '(主键)' : '';
            const isForeign = field.isForeign ? '(外键)' : '';
            html += `<li>${field.name}: ${field.type} ${isPrimary} ${isForeign}</li>`;
          });
          
          html += '</ul>';
          return html;
        }
      }
    });
  });
  
  // 添加关系边
  relationships.forEach((rel, index) => {
    // 检查源表和目标表是否存在
    const sourceExists = tables.some(table => table.name === rel.source);
    const targetExists = tables.some(table => table.name === rel.target);
    
    if (sourceExists && targetExists) {
      // 根据关系类型设置不同的线型和箭头
      let lineStyle = { type: 'solid', width: 1 };
      let symbol = ['none', 'arrow'];
      
      if (rel.type === 'one-to-one') {
        lineStyle = { type: 'solid', width: 1 };
        symbol = ['none', 'arrow'];
      } else if (rel.type === 'one-to-many') {
        lineStyle = { type: 'solid', width: 1 };
        symbol = ['none', 'arrow'];
      } else if (rel.type === 'many-to-many') {
        lineStyle = { type: 'dashed', width: 1 };
        symbol = ['arrow', 'arrow'];
      }
      
      // 添加关系边，为相邻的边设置不同的曲率，避免重叠
      edges.push({
        source: rel.source,
        target: rel.target,
        lineStyle: {
          ...lineStyle,
          curveness: index * 0.05  // 为每条边设置递增的曲率
        },
        symbol: symbol,
        label: {
          show: true,
          formatter: formatRelationship(rel),
          fontSize: 12,
          color: '#000000',
          backgroundColor: '#ffffff',
          padding: [4, 8]
        }
      });
    }
  });
  
  // 设置图表选项
  const option = {
    title: {
      text: props.chartTitle,
      subtext: `共 ${tables.length} 个类, ${relationships.length} 个关系`,
      left: 'center',
      top: 20,
      textStyle: {
        color: '#000000'
      },
      subtextStyle: {
        color: '#666666'
      }
    },
    tooltip: {
      trigger: 'item',
      backgroundColor: '#ffffff',
      borderColor: '#cccccc',
      borderWidth: 1,
      textStyle: {
        color: '#000000'
      },
      confine: true
    },
    toolbox: {
      feature: {
        saveAsImage: { show: true, title: '保存为图片' },
        restore: { show: true, title: '还原' }
      }
    },
    legend: {
      data: [
        { name: '类', icon: 'rect' }
      ],
      bottom: 10,
      textStyle: {
        color: '#000000'
      }
    },
    animationDuration: 1000,
    animationEasingUpdate: 'quinticInOut',
    series: [
      {
        name: 'UML类图',
        type: 'graph',
        layout: 'none',  // 使用预先计算的位置
        data: nodes,
        links: edges,
        roam: true,      // 允许缩放和平移
        focusNodeAdjacency: true,
        edgeSymbol: ['none', 'arrow'],
        edgeSymbolSize: 10,
        // 根据节点数量动态调整初始缩放比例
        zoom: Math.max(0.8, 1 - tables.length * 0.02),
        // 居中显示
        center: [0, 0],
        // 确保图表占满容器
        left: '5%',
        right: '5%',
        top: '10%',
        bottom: '10%',
        itemStyle: {
          borderColor: '#000000',
          borderWidth: 2
        },
        lineStyle: {
          color: '#666666',
          curveness: 0.1,  // 默认曲率
          width: 1.5
        },
        emphasis: {
          focus: 'adjacency',
          lineStyle: {
            width: 3
          },
          itemStyle: {
            borderWidth: 3
          }
        },
        label: {
          show: true
        }
      }
    ]
  };
  
  // 设置图表
  try {
    chart.value.setOption(option);
    chart.value.resize();
    
    // 确保图表完全显示所有节点
    setTimeout(() => {
      if (chart.value) {
        chart.value.dispatchAction({
          type: 'graphRoam',
          zoom: option.series[0].zoom,
          originX: 0,
          originY: 0
        });
      }
    }, 100);
  } catch (error) {
    console.error('设置图表选项时出错:', error);
  }
};

// 格式化关系文本
const formatRelationship = (rel) => {
  switch(rel.type) {
    case 'one-to-one':
      return '1..1';
    case 'one-to-many':
      return '1..*';
    case 'many-to-many':
      return '*..*';
    default:
      return '';
  }
};

// 计算初始节点位置
const calculateInitialPositions = (tables) => {
  const positions = [];
  const cols = Math.ceil(Math.sqrt(tables.length));
  const rows = Math.ceil(tables.length / cols);
  
  // 增加节点间的水平和垂直间距，扩大整体布局
  const horizontalSpacing = 500;
  const verticalSpacing = 350;
  
  // 创建更大的画布，扩大图表范围
  const centerX = (cols - 1) * horizontalSpacing / 2;
  const centerY = (rows - 1) * verticalSpacing / 2;
  
  for (let i = 0; i < tables.length; i++) {
    const row = Math.floor(i / cols);
    const col = i % cols;
    
    positions.push({
      x: col * horizontalSpacing - centerX,
      y: row * verticalSpacing - centerY
    });
  }
  
  return positions;
};

// 处理窗口大小变化
const handleResize = () => {
  if (chart.value) {
    chart.value.resize();
  }
};

// 缩放控制函数
const zoomIn = () => {
  if (chart.value) {
    const option = chart.value.getOption();
    const newOption = {
      series: [
        {
          zoom: option.series[0].zoom * 1.2
        }
      ]
    };
    chart.value.setOption(newOption);
  }
};

const zoomOut = () => {
  if (chart.value) {
    const option = chart.value.getOption();
    const newOption = {
      series: [
        {
          zoom: option.series[0].zoom / 1.2
        }
      ]
    };
    chart.value.setOption(newOption);
  }
};

const resetZoom = () => {
  if (chart.value) {
    const newOption = {
      series: [
        {
          zoom: 1
        }
      ]
    };
    chart.value.setOption(newOption);
  }
};

// 导出图表为图片
const exportChart = (type = 'png', pixelRatio = 1) => {
  if (!chart.value) return '';
  
  return chart.value.getDataURL({
    type,
    pixelRatio,
    backgroundColor: '#ffffff',
    excludeComponents: ['toolbox']
  });
};

// 生命周期钩子
onMounted(() => {
  initChart();
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  if (chart.value) {
    chart.value.dispose();
    chart.value = null;
  }
});

// 监听数据变化
watch(() => props.data, (newVal) => {
  nextTick(() => {
    updateChart();
  });
}, { deep: true });

// 监听主题变化
watch(() => props.theme, () => {
  nextTick(() => {
    // 重新初始化图表以应用新主题
    initChart();
  });
});

// 监听标题变化
watch(() => props.chartTitle, () => {
  if (chart.value) {
    nextTick(() => {
      updateChart();
    });
  }
});

// 向父组件暴露方法
defineExpose({
  zoomIn,
  zoomOut,
  resetZoom,
  exportChart,
  updateChart
});
</script>

<style scoped>
.uml-diagram-container {
  width: 100%;
  height: 100%;
  min-height: 700px;
  position: relative;
}
</style> 