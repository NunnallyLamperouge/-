<template>
  <div class="uml-diagram" ref="chartContainer">
    <div v-if="loading" class="loading-overlay">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>正在生成类图...</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue';
import * as echarts from 'echarts';
import { getBaseChartConfig, getChartOption } from './config/chartConfig';
import { processNodes, processEdges } from './utils/dataProcessor';

const props = defineProps({
  nodes: {
    type: Array,
    required: true
  },
  edges: {
    type: Array,
    required: true
  },
  chartTitle: {
    type: String,
    default: 'UML类图'
  },
  isDark: {
    type: Boolean,
    default: false
  }
});

const chartContainer = ref(null);
const chart = ref(null);
const loading = ref(false);

const initChart = () => {
  if (!chartContainer.value) return;
  
  chart.value = echarts.init(chartContainer.value, getBaseChartConfig(props.isDark));
  updateChart();
};

const updateChart = () => {
  if (!chart.value) return;
  
  const processedNodes = processNodes(props.nodes);
  const processedEdges = processEdges(props.edges);
  
  const option = getChartOption(props, processedNodes, processedEdges);
  chart.value.setOption(option);
};

const handleResize = () => {
  chart.value?.resize();
};

watch(() => [props.nodes, props.edges, props.isDark], () => {
  updateChart();
}, { deep: true });

onMounted(() => {
  initChart();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  chart.value?.dispose();
});
</script>

<style scoped>
.uml-diagram {
  width: 100%;
  height: 100%;
  min-height: 600px;
  position: relative;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.loading-overlay .el-icon {
  font-size: 24px;
  margin-bottom: 10px;
}
</style> 