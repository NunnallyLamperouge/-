<template>
  <div class="page-container">
    <!-- 侧边栏 -->
    <div class="sidebar">
      <h2 class="sidebar-title">历史答辩记录</h2>
      <div class="history-list">
        <div 
          v-for="record in historyRecords" 
          :key="record.id"
          class="history-item"
          :class="{ active: virtualId === record.id }"
          @click="viewHistoryResult(record.id)"
        >
          <div class="history-title">{{ record.fake_template.title }}的答辩记录</div>
          <div class="history-info">
            <span class="difficulty" :class="record.difficulty">{{ getDifficultyText(record.difficulty) }}</span>
            <span class="date">{{ formatDate(record.CreatedAt) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-content">
      <h1 class="title">虚拟答辩</h1>
      
      <VirtualSettings 
        v-if="currentView === 'settings'"
        :loading="loading"
        :thesisList="thesisList"
        v-model:selectedThesis="selectedThesis"
        v-model:difficulty="difficulty"
        v-model:questionCount="questionCount"
        @startDefense="startDefense"
      />
      
      <VirtualDefense
        v-if="currentView === 'defense'"
        :thesisId="selectedThesis"
        :difficulty="difficulty"
        :questionCount="questionCount"
        @back="currentView = 'settings'"
        @showResult="showResult"
      />

      <VirtualResult
        v-if="currentView === 'result'"
        :virtualId="virtualId"
        @back="currentView = 'settings'"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { get } from '@/utility/request';
import VirtualSettings from '@/components/student/VirtualSettings.vue';
import VirtualDefense from '@/components/student/VirtualDefense.vue';
import VirtualResult from '@/components/student/VirtualResult.vue';

const loading = ref(true);
const thesisList = ref([]);
const currentView = ref('settings');
const historyRecords = ref([]);

// 表单数据
const selectedThesis = ref('');
const difficulty = ref('medium');
const questionCount = ref(5);
const virtualId = ref('');

// 获取论文列表
const fetchThesisList = () => {
  loading.value = true;
  get('/api/template/list', null, 
    (message, data) => {
      thesisList.value = data.map(item => ({
        id: item.id,
        title: item.title,
        description: item.description,
        type: item.type,
        content: item.content
      }));
      loading.value = false;
    }, 
    (errorMsg) => {
      console.error('获取论文列表失败:', errorMsg);
      loading.value = false;
    },
    (error) => {
      console.error('获取论文列表出错:', error);
      loading.value = false;
    }
  );
};

// 获取历史答辩记录
const fetchHistoryRecords = () => {
  get('/api/virtual/list', null,
    (message, data) => {
      historyRecords.value = data;
    },
    (errorMsg) => {
      console.error('获取历史记录失败:', errorMsg);
    }
  );
};

// 开始答辩
const startDefense = () => {
  if (!selectedThesis.value) return;
  currentView.value = 'defense';
};

// 显示结果
const showResult = (id) => {
  virtualId.value = id;
  currentView.value = 'result';
};

// 查看历史记录结果
const viewHistoryResult = (id) => {
  virtualId.value = id;
  currentView.value = 'result';
};

// 格式化难度显示
const getDifficultyText = (difficulty) => {
  const map = {
    easy: '简单',
    medium: '中等',
    hard: '困难'
  };
  return map[difficulty] || difficulty;
};

// 格式化日期
const formatDate = (dateStr) => {
  const date = new Date(dateStr);
  return `${date.getMonth() + 1}月${date.getDate()}日`;
};

// 组件挂载时获取数据
onMounted(() => {
  fetchThesisList();
  fetchHistoryRecords();
});
</script>

<style scoped>
.page-container {
  display: flex;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.sidebar {
  width: 280px;
  background-color: #fff;
  border-right: 1px solid #e8e8e8;
  padding: 20px;
  overflow-y: auto;
}

.sidebar-title {
  font-size: 18px;
  color: #2c3e50;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e8e8e8;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.history-item {
  padding: 12px;
  background-color: #f9f9f9;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s;
}

.history-item:hover {
  background-color: #f0f0f0;
}

.history-item.active {
  background-color: #e6f7ff;
  border-left: 3px solid #1890ff;
}

.history-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 8px;
}

.history-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.difficulty {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.difficulty.easy {
  background-color: #f6ffed;
  color: #52c41a;
}

.difficulty.medium {
  background-color: #fff7e6;
  color: #fa8c16;
}

.difficulty.hard {
  background-color: #fff1f0;
  color: #f5222d;
}

.date {
  color: #999;
}

.main-content {
  flex: 1;
  padding: 20px;
}

.title {
  text-align: center;
  margin-bottom: 30px;
  color: #2c3e50;
}
</style>

