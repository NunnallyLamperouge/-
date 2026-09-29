<template>
  <div class="max-w-4xl mx-auto p-6 bg-white rounded-lg shadow-lg">
    <div v-if="loading" class="flex flex-col items-center justify-center py-16">
      <a-spin size="large" />
      <p class="mt-4 text-gray-600">正在加载测评结果...</p>
    </div>
    
    <div v-else-if="error" class="flex flex-col items-center justify-center py-16">
      <a-result status="error" :title="'加载失败'" :sub-title="error">
        <template #extra>
          <a-button type="primary" @click="fetchResult">重试</a-button>
        </template>
      </a-result>
    </div>

    <div v-else class="result-content">
      <h2 
        class="text-2xl font-bold text-center text-gray-800 mb-8"
        v-motion
        :initial="{ opacity: 0, y: 50 }"
        :enter="{ opacity: 1, y: 0, transition: { delay: 100, duration: 500 } }"
      >
        {{ result.Virtual.user.username }}同学的答辩测评结果
      </h2>
      
      <div 
        class="bg-blue-50 rounded-xl p-6 shadow-sm mb-8"
        v-motion
        :initial="{ opacity: 0, scale: 0.9 }"
        :enter="{ opacity: 1, scale: 1, transition: { delay: 200, duration: 500 } }"
      >
        <div class="flex flex-col md:flex-row justify-between items-center">
          <div class="total-score text-center mb-6 md:mb-0">
            <h3 class="text-lg font-semibold text-gray-700">总分</h3>
            <div class="text-5xl font-bold text-blue-600 mt-2">{{ result.Virtual.score || 0 }}</div>
          </div>
          
          <div class="flex-1 md:ml-10 grid grid-cols-1 gap-4">
            <div class="stat-item bg-white p-4 rounded-lg shadow-sm flex justify-between items-center">
              <span class="text-gray-600">答题数量:</span>
              <span class="font-semibold text-gray-800">{{ result.VirtualQuestions.length }}</span>
            </div>
            <div class="stat-item bg-white p-4 rounded-lg shadow-sm flex justify-between items-center">
              <span class="text-gray-600">难度等级:</span>
              <span class="font-semibold text-gray-800">{{ difficultyText[result.Virtual.difficulty] || result.Virtual.difficulty }}</span>
            </div>
            <div class="stat-item bg-white p-4 rounded-lg shadow-sm flex justify-between items-center">
              <span class="text-gray-600">测评时间:</span>
              <span class="font-semibold text-gray-800">{{ formatDate(result.Virtual.CreatedAt) }}</span>
            </div>
          </div>
        </div>
      </div>

      <div
        class="bg-white border border-gray-200 rounded-xl p-6 shadow-sm mb-8"
        v-motion
        :initial="{ opacity: 0 }"
        :enter="{ opacity: 1, transition: { delay: 300, duration: 500 } }"
      >
        <h3 class="text-xl font-semibold text-gray-800 mb-4">测评详情</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="bg-gray-50 p-4 rounded-lg">
            <div class="text-sm text-gray-500">测评模板</div>
            <div class="font-medium">{{ result.Virtual.fake_template?.title || '未命名模板' }}</div>
          </div>
          <div class="bg-gray-50 p-4 rounded-lg">
            <div class="text-sm text-gray-500">学生姓名</div>
            <div class="font-medium">{{ result.Virtual.user?.username || '未知用户' }}</div>
          </div>
          <div class="bg-gray-50 p-4 rounded-lg">
            <div class="text-sm text-gray-500">专业班级</div>
            <div class="font-medium">{{ result.Virtual.user?.profession || '未设置' }}</div>
          </div>
        </div>
        <div class="mt-4 bg-gray-50 p-4 rounded-lg">
          <div class="text-sm text-gray-500">测评描述</div>
          <MdPreview
            style="background: transparent"
            editorId="analysis-preview"
            :modelValue="result.Virtual.description"
            previewTheme="github"
          />
        </div>
      </div>

      <div
        class="questions-review"
        v-motion
        :initial="{ opacity: 0 }"
        :enter="{ opacity: 1, transition: { delay: 400, duration: 500 } }"
      >
        <h3 class="text-xl font-semibold text-gray-800 mb-4">答题回顾</h3>
        <a-collapse>
          <a-collapse-panel 
            v-for="(question, index) in result.VirtualQuestions" 
            :key="index" 
            :header="`问题 ${question.question}`"
          >
            <div 
              class="question-content"
              v-motion
              :initial="{ opacity: 0, y: 20 }"
              :enter="{ opacity: 1, y: 0, transition: { delay: 100 * index, duration: 300 } }"
            >              
              <div class="bg-gray-50 p-4 rounded-lg mb-4">
                <div class="font-medium text-gray-700 mb-2">你的回答:</div>
                <p class="text-gray-800">{{ question.answer }}</p>
              </div>
              
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
                <div v-if="question.question_audio_path" class="audio-player">
                  <div class="text-sm text-gray-600 mb-2">问题音频:</div>
                  <audio 
                    controls 
                    class="w-full" 
                    :src="getAudioUrl(question.question_audio_path)"
                  ></audio>
                </div>
                
                <div v-if="question.answer_audio_path" class="audio-player">
                  <div class="text-sm text-gray-600 mb-2">回答音频:</div>
                  <audio 
                    controls 
                    class="w-full" 
                    :src="getAudioUrl(question.answer_audio_path)"
                  ></audio>
                </div>
              </div>
            </div>
          </a-collapse-panel>
        </a-collapse>
      </div>

      <div 
        class="mt-10 text-center"
        v-motion
        :initial="{ opacity: 0, y: 20 }"
        :enter="{ opacity: 1, y: 0, transition: { delay: 500, duration: 500 } }"
      >
        <a-button type="primary" size="large" @click="$emit('back')" class="mr-4">返回</a-button>
        <a-button type="primary" size="large" @click="exportToPdf" :loading="exporting">
          <template #icon><download-outlined /></template>
          导出PDF
        </a-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, h, reactive } from 'vue';
import { get, postJSON } from '@/utility/request';
import { message } from 'ant-design-vue';
import { useMotion } from '@vueuse/motion';
import { BACKEND_DOMAIN } from '@/utility/config';
import { MdPreview } from "md-editor-v3"
import "md-editor-v3/lib/preview.css"
import { DownloadOutlined } from '@ant-design/icons-vue';
import html2canvas from 'html2canvas';
import jsPDF from 'jspdf';

const props = defineProps({
  virtualId: {
    type: String,
    required: true
  }
});

defineEmits(['back']);

const loading = ref(true);
const error = ref('');
const result = ref({
  Virtual: {
    id: 0,
    student_id: 0,
    difficulty: '',
    score: 0,
    seconds: 0,
    description: '',
    CreatedAt: '',
    user: {
      username: '',
      profession: ''
    },
    fake_template: {
      title: ''
    }
  },
  VirtualQuestions: []
});

const exporting = ref(false);

// 难度级别的文本映射
const difficultyText = {
  'easy': '简单',
  'medium': '中等',
  'hard': '困难'
};

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '未知';
  const date = new Date(dateString);
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`;
};

// 格式化秒数为时分秒
const formatSeconds = (seconds) => {
  if (!seconds) return '0秒';
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return mins > 0 ? `${mins}分${secs}秒` : `${secs}秒`;
};

// 获取音频URL
const getAudioUrl = (path) => {
  return BACKEND_DOMAIN + `audio/${path}`;
};

const fetchResult = async () => {
  loading.value = true;
  error.value = '';
  
  get('/api/virtual/get', { virtual_id: props.virtualId },
    (msg, data) => {
      result.value = data;
      loading.value = false;
      if (result.value.Virtual.score === 0 || result.value.Virtual.description === "") {
        generateDescription();
      }
    },
    (errorMsg) => {
      error.value = errorMsg;
      loading.value = false;
      message.error('获取测评结果失败');
    },
    (err) => {
      error.value = '系统错误，请稍后重试';
      loading.value = false;
      message.error('系统错误');
    }
  );
};

onMounted(() => {
  fetchResult();
});

const options = reactive({

})

const generateDescription = () => {
  // 创建 WebSocket 连接
  const token = localStorage.getItem("authToken")
  const baseUrl = BACKEND_DOMAIN.replace(/^http/, 'ws').replace(/\/$/, '')
  const ws = new WebSocket(
    `${baseUrl}/ws/projectSuggest`
  )

  ws.onopen = () => {
    ws.send(
      JSON.stringify({
        content: result.value,
        auth_token: token
      })
    )
  }

  ws.onmessage = (event) => {
    const response = JSON.parse(event.data)
    if (response.code === 500) {
      message.error(response.message)
      analysisOptions.isGenerating = false
      ws.close()
      return
    }

    if (response.done) {
      ws.close()
      return
    }

    result.value.Virtual.description += response.content
  }

  ws.onerror = (error) => {
    console.error("WebSocket error:", error)
    message.error("连接发生错误，请重试")
  }

  ws.onclose = () => {
    getScore()
  }
};

const getScore = () => {
  postJSON('/api/virtual/score', {
    result: result.value
  },
  (msg, data) => {
    fetchResult()
  },
  (errorMsg) => {
    message.error(errorMsg)
  }
  )
}

const exportToPdf = async () => {
  exporting.value = true;
  try {
    // 获取要导出的内容元素
    const content = document.querySelector('.result-content');
    if (!content) {
      throw new Error('无法找到要导出的内容');
    }

    // 创建PDF
    const pdf = new jsPDF('p', 'mm', 'a4');
    const contentWidth = content.offsetWidth;
    const contentHeight = content.offsetHeight;
    
    // 使用html2canvas将内容转换为canvas
    const canvas = await html2canvas(content, {
      scale: 2, // 提高清晰度
      useCORS: true, // 允许加载跨域图片
      logging: false,
      backgroundColor: '#ffffff'
    });
    
    // 计算PDF页面宽度
    const pdfWidth = pdf.internal.pageSize.getWidth();
    const pdfHeight = (contentHeight * pdfWidth) / contentWidth;
    
    // 将canvas内容添加到PDF
    const imgData = canvas.toDataURL('image/jpeg', 1.0);
    pdf.addImage(imgData, 'JPEG', 0, 0, pdfWidth, pdfHeight);
    
    // 如果内容太长，分页处理
    if (pdfHeight > pdf.internal.pageSize.getHeight()) {
      let heightLeft = pdfHeight;
      let position = 0;
      const pageHeight = pdf.internal.pageSize.getHeight();
      
      while (heightLeft > 0) {
        position -= pageHeight;
        pdf.addPage();
        pdf.addImage(imgData, 'JPEG', 0, position, pdfWidth, pdfHeight);
        heightLeft -= pageHeight;
      }
    }
    
    // 下载PDF
    const filename = `${result.value.Virtual.user?.username || '未知用户'}_答辩测评结果_${new Date().toISOString().slice(0, 10)}.pdf`;
    pdf.save(filename);
  } catch (err) {
    console.error('导出PDF失败:', err);
    alert('导出PDF失败，请重试');
  } finally {
    exporting.value = false;
  }
};
</script>

<style scoped>
/* 自定义样式可以在这里添加，但主要使用Tailwind的类 */
audio::-webkit-media-controls-panel {
  background-color: #f3f4f6;
}

audio {
  border-radius: 8px;
}

.result-content {
  max-width: 100%;
}

.audio-player audio {
  width: 100%;
  height: 40px;
}

/* 适配PDF导出 */
@media print {
  .result-content {
    padding: 0;
    margin: 0;
  }
}
</style> 