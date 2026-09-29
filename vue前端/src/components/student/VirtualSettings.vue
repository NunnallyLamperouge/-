<template>
  <div class="container mx-auto p-4 md:p-6 lg:p-8 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 min-h-screen">
    <div 
      class="mb-8 border-l-4 border-blue-500 dark:border-indigo-400 pl-4"
      v-motion
      :initial="{ opacity: 0, x: -50 }"
      :enter="{ opacity: 1, x: 0, transition: { duration: 500 } }"
    >
      <h1 class="text-2xl md:text-3xl font-bold mb-2 text-gray-800 dark:text-white bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-500 dark:from-blue-400 dark:to-indigo-300">虚拟答辩</h1>
      <p class="text-gray-600 dark:text-gray-400">模拟真实答辩场景，提前做好准备</p>
    </div>

    <!-- 选择论文区域 -->
    <div 
      class="mb-10 bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6 backdrop-blur-sm bg-opacity-80 dark:bg-opacity-80 border border-gray-100 dark:border-gray-700"
      v-motion
      :initial="{ opacity: 0, y: 30 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 600 } }"
    >
      <h2 class="text-xl font-semibold mb-6 flex items-center text-gray-800 dark:text-white">
        <span class="inline-block w-10 h-10 rounded-full bg-gradient-to-r from-blue-400 to-indigo-500 flex items-center justify-center mr-3 shadow-md">
          <FileTextOutlined class="text-white" />
        </span>
        选择论文
      </h2>

      <div v-if="loading" class="flex justify-center py-6">
        <div class="loading-spinner">
          <LoadingOutlined class="text-4xl text-blue-500 dark:text-indigo-400" />
        </div>
      </div>

      <div v-else-if="thesisList.length === 0" class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-8 text-center">
        <a-empty description="暂无可用论文" />
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div 
          v-for="(thesis, index) in thesisList" 
          :key="thesis.id"
          class="thesis-card bg-white dark:bg-gray-800 rounded-xl overflow-hidden cursor-pointer group transition-all duration-300 transform hover:scale-105 hover:-rotate-1 shadow-md hover:shadow-2xl border border-gray-100 dark:border-gray-700 relative hover:z-10 h-full flex flex-col"
          :class="{'ring-2 ring-offset-2 dark:ring-offset-gray-800': selectedThesis === thesis.id}"
          :style="{'--ring-color': getThesisColor(thesis.type).ring}"
          v-motion
          :initial="{ opacity: 0, y: 50 }"
          :enter="{ opacity: 1, y: 0, transition: { delay: 100 * index, duration: 300 } }"
          @click="$emit('update:selectedThesis', thesis.id)"
        >
          <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r" :class="getThesisColor(thesis.type).gradient"></div>
          <div class="absolute inset-0 opacity-0 group-hover:opacity-100 bg-gradient-to-br from-white/5 to-white/20 dark:from-white/5 dark:to-white/10 transition-opacity duration-500 pointer-events-none"></div>
          
          <div class="p-4 flex flex-col flex-grow">
            <div class="flex items-start justify-between mb-3">
              <div class="w-10 h-10 rounded-full flex items-center justify-center shadow-lg" :class="getThesisColor(thesis.type).icon">
                <component :is="getThesisIcon(thesis.type)" class="text-white" />
              </div>
              <a-tag :color="getThesisColor(thesis.type).tag" class="rounded-full">{{ getThesisTypeName(thesis.type) }}</a-tag>
            </div>
            <h3 class="text-lg font-medium mb-2 text-gray-800 dark:text-white line-clamp-2" :title="thesis.title">{{ thesis.title }}</h3>
            <p v-if="thesis.description" class="text-sm text-gray-500 dark:text-gray-400 mt-auto mb-3 line-clamp-2">
              {{ thesis.description }}
            </p>
            <div v-if="thesis.author || thesis.date" class="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-auto mb-2">
              <span v-if="thesis.author" class="flex items-center">
                <UserOutlined class="mr-1" />
                {{ thesis.author }}
              </span>
              <span v-if="thesis.date" class="flex items-center">
                <CalendarOutlined class="mr-1" />
                {{ thesis.date }}
              </span>
            </div>
          </div>
          
          <div class="bg-gradient-to-r p-2 text-center relative overflow-hidden" :class="getThesisColor(thesis.type).button">
            <div class="absolute inset-0 w-full h-full bg-white/0 group-hover:bg-white/10 transition-all duration-300"></div>
            <span class="text-white font-medium flex items-center justify-center relative z-10">
              选择
              <CheckOutlined class="ml-1 text-xs opacity-0 group-hover:opacity-100 transition-all duration-300 group-hover:translate-x-1" />
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 设置区域 -->
    <div 
      class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6 backdrop-blur-sm bg-opacity-80 dark:bg-opacity-80 border border-gray-100 dark:border-gray-700"
      v-motion
      :initial="{ opacity: 0, y: 30 }"
      :enter="{ opacity: 1, y: 0, transition: { duration: 600, delay: 200 } }"
    >
      <h2 class="text-xl font-semibold mb-6 flex items-center text-gray-800 dark:text-white">
        <span class="inline-block w-10 h-10 rounded-full bg-gradient-to-r from-purple-400 to-pink-500 flex items-center justify-center mr-3 shadow-md">
          <SettingOutlined class="text-white" />
        </span>
        答辩设置
      </h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div 
          class="bg-white dark:bg-gray-800 rounded-xl overflow-hidden shadow-md border border-gray-100 dark:border-gray-700 p-5"
          v-motion
          :initial="{ opacity: 0, x: -20 }"
          :enter="{ opacity: 1, x: 0, transition: { delay: 200 } }"
        >
          <h3 class="text-lg font-medium mb-4 text-gray-800 dark:text-white flex items-center">
            <DashboardOutlined class="mr-2 text-blue-500" />
            问题难度
          </h3>
          <a-radio-group 
            :value="difficulty" 
            @change="(e) => $emit('update:difficulty', e.target.value)"
            class="w-full flex flex-col space-y-3"
          >
            <a-radio 
              v-for="level in difficultyLevels" 
              :key="level.value" 
              :value="level.value"
              class="p-2 w-full border border-gray-100 dark:border-gray-700 rounded-lg transition-all duration-300 hover:bg-gray-50 dark:hover:bg-gray-700/50"
              :class="{'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-700': difficulty === level.value}"
            >
              <div class="flex items-center">
                <div class="w-4 h-4 rounded-full mr-2" :class="getDifficultyColor(level.value)"></div>
                {{ level.label }}
              </div>
            </a-radio>
          </a-radio-group>
        </div>
        
        <div 
          class="bg-white dark:bg-gray-800 rounded-xl overflow-hidden shadow-md border border-gray-100 dark:border-gray-700 p-5"
          v-motion
          :initial="{ opacity: 0, x: 20 }"
          :enter="{ opacity: 1, x: 0, transition: { delay: 300 } }"
        >
          <h3 class="text-lg font-medium mb-4 text-gray-800 dark:text-white flex items-center">
            <OrderedListOutlined class="mr-2 text-purple-500" />
            问题数量
          </h3>
          <div class="flex flex-col items-center justify-center h-[calc(100%-3.5rem)]">
            <div class="text-5xl font-bold text-gray-800 dark:text-white mb-4">{{ questionCount }}</div>
            <div class="flex items-center gap-4">
              <a-button
                type="primary"
                shape="circle"
                size="large"
                @click="decreaseQuestionCount"
                :disabled="questionCount <= 3"
                class="shadow-md hover:shadow-lg transition-all duration-300"
              >
                <template #icon><MinusOutlined /></template>
              </a-button>
              <a-slider 
                :value="questionCount" 
                :min="3" 
                :max="10" 
                :step="1" 
                @change="(value) => $emit('update:questionCount', value)" 
                class="w-40 mx-4" 
              />
              <a-button
                type="primary"
                shape="circle"
                size="large"
                @click="increaseQuestionCount"
                :disabled="questionCount >= 10"
                class="shadow-md hover:shadow-lg transition-all duration-300"
              >
                <template #icon><PlusOutlined /></template>
              </a-button>
            </div>
          </div>
        </div>
      </div>
      
      <div 
        class="mt-8"
        v-motion
        :initial="{ opacity: 0, y: 20 }"
        :enter="{ opacity: 1, y: 0, transition: { delay: 400 } }"
      >
        <a-button
          type="primary"
          size="large"
          :disabled="!selectedThesis || loading"
          @click="$emit('startDefense')"
          class="w-full h-14 text-lg font-medium rounded-xl shadow-lg hover:shadow-xl transition-all duration-300 bg-gradient-to-r from-blue-500 to-indigo-600 border-0 hover:scale-[1.02]"
        >
          <template #icon><RocketOutlined /></template>
          开始答辩
        </a-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';
import { 
  PlusOutlined, 
  MinusOutlined, 
  FileTextOutlined, 
  SettingOutlined,
  CheckOutlined,
  LoadingOutlined,
  RocketOutlined,
  DashboardOutlined,
  OrderedListOutlined,
  UserOutlined,
  CalendarOutlined,
  BookOutlined,
  FileWordOutlined,
  CheckSquareOutlined
} from '@ant-design/icons-vue';

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  thesisList: {
    type: Array,
    default: () => []
  },
  selectedThesis: {
    type: String,
    default: ''
  },
  difficulty: {
    type: String,
    default: 'medium'
  },
  questionCount: {
    type: Number,
    default: 5
  }
});

const emit = defineEmits(['update:selectedThesis', 'update:difficulty', 'update:questionCount', 'startDefense']);

// 难度级别选项
const difficultyLevels = [
  { value: 'easy', label: '简单' },
  { value: 'medium', label: '中等' },
  { value: 'hard', label: '困难' }
];

// 获取难度对应的颜色
const getDifficultyColor = (difficultyValue) => {
  switch (difficultyValue) {
    case 'easy': return 'bg-green-500';
    case 'medium': return 'bg-yellow-500';
    case 'hard': return 'bg-red-500';
    default: return 'bg-blue-500';
  }
};

// 获取论文类型图标
const getThesisIcon = (typeId) => {
  switch (typeId) {
    case 'task': return FileTextOutlined;
    case 'proposal': return BookOutlined;
    case 'midterm': return CheckSquareOutlined;
    case 'thesis': return FileWordOutlined;
    default: return FileTextOutlined;
  }
};

// 获取论文类型名称
const getThesisTypeName = (typeId) => {
  switch (typeId) {
    case 'task': return '任务书';
    case 'proposal': return '开题报告';
    case 'midterm': return '中期检查';
    case 'thesis': return '毕业论文';
    default: return '论文';
  }
};

// 获取论文类型对应的颜色方案
const getThesisColor = (typeId) => {
  switch (typeId) {
    case 'task':
      return {
        ring: 'ring-blue-500',
        gradient: 'from-blue-400 to-blue-600',
        icon: 'bg-gradient-to-br from-blue-400 to-blue-600',
        tag: 'blue',
        button: 'from-blue-500 to-blue-700'
      };
    case 'proposal':
      return {
        ring: 'ring-green-500',
        gradient: 'from-green-400 to-green-600',
        icon: 'bg-gradient-to-br from-green-400 to-green-600',
        tag: 'green',
        button: 'from-green-500 to-green-700'
      };
    case 'midterm':
      return {
        ring: 'ring-yellow-500',
        gradient: 'from-yellow-400 to-yellow-600',
        icon: 'bg-gradient-to-br from-yellow-400 to-yellow-600',
        tag: 'gold',
        button: 'from-yellow-500 to-yellow-700'
      };
    case 'thesis':
      return {
        ring: 'ring-purple-500',
        gradient: 'from-purple-400 to-purple-600',
        icon: 'bg-gradient-to-br from-purple-400 to-purple-600',
        tag: 'purple',
        button: 'from-purple-500 to-purple-700'
      };
    default:
      return {
        ring: 'ring-blue-500',
        gradient: 'from-blue-400 to-indigo-600',
        icon: 'bg-gradient-to-br from-blue-400 to-indigo-600',
        tag: 'blue',
        button: 'from-blue-500 to-indigo-600'
      };
  }
};

// 问题数量控制
const decreaseQuestionCount = () => {
  if (props.questionCount > 3) {
    emit('update:questionCount', props.questionCount - 1);
  }
};

const increaseQuestionCount = () => {
  if (props.questionCount < 10) {
    emit('update:questionCount', props.questionCount + 1);
  }
};
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 加载动画 */
.loading-spinner {
  animation: spin 1.5s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg) scale(1);
  }
  50% {
    transform: rotate(180deg) scale(1.2);
  }
  100% {
    transform: rotate(360deg) scale(1);
  }
}

/* 论文卡片样式 */
.thesis-card {
  transition: all 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  box-shadow: 0 10px 20px rgba(0,0,0,0.05), 0 6px 6px rgba(0,0,0,0.07);
  will-change: transform, box-shadow;
  transform-style: preserve-3d;
  perspective: 1000px;
}

.thesis-card:hover {
  transform: translateY(-8px) scale(1.01) rotateX(2deg) rotateY(-2deg);
  box-shadow: 0 20px 30px rgba(0,0,0,0.1), 0 10px 10px rgba(0,0,0,0.05), 0 0 15px rgba(59, 130, 246, 0.2);
}

.thesis-card:hover::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 0.75rem;
  padding: 2px;
  background: linear-gradient(45deg, rgba(59, 130, 246, 0.3), rgba(147, 197, 253, 0.3), rgba(59, 130, 246, 0.3));
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  pointer-events: none;
  animation: border-glow 2s linear infinite;
}

@keyframes border-glow {
  0% {
    background-position: 0% 50%;
  }
  100% {
    background-position: 200% 50%;
  }
}
</style>