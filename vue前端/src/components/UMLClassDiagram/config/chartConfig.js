export const getBaseChartConfig = (isDark) => ({
  renderer: 'canvas',
  devicePixelRatio: window.devicePixelRatio,
  backgroundColor: isDark ? '#1a1a1a' : '#ffffff'
});

export const getChartOption = (props, nodes, edges) => ({
  title: {
    text: props.chartTitle,
    subtext: `共 ${nodes.length} 个类, ${edges.length} 个关系`,
    left: 'center',
    top: 20,
    textStyle: {
      color: props.isDark ? '#ffffff' : '#000000'
    },
    subtextStyle: {
      color: props.isDark ? '#999999' : '#666666'
    }
  },
  tooltip: {
    trigger: 'item',
    backgroundColor: props.isDark ? '#2a2a2a' : '#ffffff',
    borderColor: props.isDark ? '#444444' : '#cccccc',
    borderWidth: 1,
    textStyle: {
      color: props.isDark ? '#ffffff' : '#000000'
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
    data: [{ name: '类', icon: 'rect' }],
    bottom: 10,
    textStyle: {
      color: props.isDark ? '#ffffff' : '#000000'
    }
  },
  animationDuration: 1000,
  animationEasingUpdate: 'quinticInOut',
  series: [
    {
      name: 'UML类图',
      type: 'graph',
      layout: 'none',
      data: nodes,
      links: edges,
      roam: true,
      focusNodeAdjacency: true,
      edgeSymbol: ['none', 'arrow'],
      edgeSymbolSize: 10,
      zoom: Math.max(0.8, 1 - nodes.length * 0.02),
      center: [0, 0],
      left: '5%',
      right: '5%',
      top: '10%',
      bottom: '10%',
      itemStyle: {
        borderColor: props.isDark ? '#ffffff' : '#000000',
        borderWidth: 2
      },
      lineStyle: {
        color: props.isDark ? '#999999' : '#666666',
        curveness: 0.1,
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
}); 