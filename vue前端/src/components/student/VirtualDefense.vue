<template>
  <div class="virtual-defense">
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>正在生成答辩问题，请稍候...</p>
    </div>

    <div v-else-if="error" class="error-container">
      <h3>加载失败</h3>
      <p>{{ errorMessage }}</p>
      <button class="action-button retry-button" @click="fetchQuestions">重试</button>
      <button class="action-button back-button" @click="$emit('back')">返回设置</button>
    </div>

    <div v-else class="defense-content">
      <div class="defense-header">
        <h2>虚拟答辩</h2>
        <div class="defense-info">
          <span>难度：{{ getDifficultyLabel(difficulty) }}</span>
          <span>问题数量：{{ questions.length }}</span>
        </div>
      </div>

      <div class="questions-container">
        <div v-if="currentQuestion" class="question-item">
          <div class="question-header">
            <h3>问题 {{ currentQuestionIndex + 1 }} / {{ questions.length }}</h3>
            <div class="question-status" :class="getQuestionStatusClass(currentQuestion)">
              {{ getQuestionStatusText(currentQuestion) }}
            </div>
          </div>
          
          <div class="question-content">
            <div class="audio-text-container">
              <div class="audio-controls">
                <audio
                  ref="audioPlayer"
                  @ended="onAudioEnded"
                  @timeupdate="onTimeUpdate"
                  style="display: none;"
                ></audio>
                <div 
                  v-if="isLoadingAudio || audioLoadingStates[currentQuestion?.id]" 
                  class="loading-indicator"
                >
                  加载问题中...
                </div>
                <div v-else-if="isPlaying" class="playing-indicator">
                  正在播放问题...
                </div>
              </div>
              
              <div class="typewriter-text text-xl" ref="typewriterText">
                {{ displayedText }}
              </div>
            </div>
            
            <template v-if="!isSubmitted">
              <div v-if="hasPlayed" class="recording-section">
                <div class="recording-status" :class="{ 'active-recording': isRecording }">
                  <div v-if="isRecording" class="recording-indicator">
                    <span class="recording-dot"></span>
                    正在录制您的回答...（{{ Math.floor(recordingDuration) }}秒）
                  </div>
                  <div v-else-if="currentQuestion.userAnswer">
                    已完成录音
                  </div>
                  <div v-else>
                    准备开始录音...
                  </div>
                </div>
                <div v-if="currentQuestion.userAnswer" class="recognized-text">
                  <div class="recognized-label">录音状态：</div>
                  <div class="recognized-content">{{ currentQuestion.userAnswer }}</div>
                </div>
              </div>
              
              <div class="question-actions">
                <button 
                  v-if="hasNextQuestion" 
                  class="next-button"
                  @click="goToNextQuestion"
                  :disabled="!hasPlayed || !currentQuestion.userAnswer || isRecording"
                >
                  {{ isRecording ? '录音中...' : '完成作答，进入下一题' }}
                </button>
                <button 
                  v-if="!hasNextQuestion && hasPlayed"
                  class="submit-all-button"
                  @click="submitAllAnswers"
                  :disabled="!isAllQuestionsAnswered || isRecording"
                >
                  {{ isRecording ? '录音中...' : '提交所有答案' }}
                </button>
              </div>
            </template>
            
            <template v-else-if="currentQuestion.status === 'submitted'">
              <div class="user-answer">
                <div class="answer-label">您的答案：</div>
                <div class="answer-content">{{ currentQuestion.userAnswer }}</div>
              </div>
              <div class="evaluation-loading">评估中...</div>
            </template>
            
            <template v-else-if="currentQuestion.status === 'evaluated'">
              <div class="user-answer">
                <div class="answer-label">您的答案：</div>
                <div class="answer-content">{{ currentQuestion.userAnswer }}</div>
              </div>
              
              <div class="evaluation-result">
                <div class="evaluation-score">
                  得分: <span>{{ currentQuestion.score }} / 10</span>
                </div>
                <div class="evaluation-feedback">
                  <div class="feedback-label">评价：</div>
                  <div class="feedback-content">{{ currentQuestion.feedback }}</div>
                </div>
                <div class="evaluation-reference">
                  <div class="reference-label">参考答案：</div>
                  <div class="reference-content">{{ currentQuestion.referenceAnswer }}</div>
                </div>
              </div>

              <div class="question-actions">
                <button 
                  v-if="hasNextQuestion" 
                  class="next-button"
                  @click="goToNextQuestion"
                >
                  下一题
                </button>
              </div>
            </template>
          </div>
        </div>

        <div class="progress-indicator">
          <div 
            v-for="(question, index) in questions" 
            :key="index"
            class="progress-dot"
            :class="{
              'active': index === currentQuestionIndex,
              'completed': question.status === 'evaluated',
              'answered': question.userAnswer && question.status === 'pending'
            }"
            @click="goToQuestion(index)"
          ></div>
        </div>
      </div>

      <div class="defense-actions">
        <button class="action-button back-button" @click="$emit('back')">返回设置</button>
        <button 
          v-if="isAllQuestionsEvaluated" 
          class="action-button finish-button" 
          @click="finishDefense"
        >
          完成答辩
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { get, post , postJSON } from '@/utility/request';
import { BACKEND_DOMAIN } from '@/utility/config';
import { message } from 'ant-design-vue';

const props = defineProps({
  thesisId: {
    type: String,
    required: true
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

const emit = defineEmits(['back', 'showResult']);

const loading = ref(true);
const error = ref(false);
const errorMessage = ref('');
const questions = ref([]);
const questionContents = ref([]); // 存储原始问题内容
const audioFileNames = ref([]); // 存储音频文件名

// 新增音频缓存相关的响应式变量
const audioLoadingStates = ref({}); // 存储每个问题的音频加载状态
const audioCachePromises = ref({}); // 存储音频加载的Promise

const currentQuestionIndex = ref(0);
const currentQuestion = computed(() => questions.value[currentQuestionIndex.value]);
const hasNextQuestion = computed(() => currentQuestionIndex.value < questions.value.length - 1);

const isSubmitted = ref(false);

const audioPlayer = ref(null);
const typewriterText = ref(null);
const isPlaying = ref(false);
const isLoadingAudio = ref(false);
const hasPlayed = ref(false);
const displayedText = ref('');
const audioCache = ref({});

// 添加录音相关状态
const isRecording = ref(false);
const mediaRecorder = ref(null);
const audioChunks = ref([]);
const userAudioBlobs = ref({});
const recordingTimer = ref(null);
const recordingDuration = ref(0);

// 预加载下一个问题的音频
const preloadNextQuestionAudio = async () => {
  const nextIndex = currentQuestionIndex.value + 1;
  if (nextIndex < questions.value.length) {
    const nextQuestion = questions.value[nextIndex];
    if (nextQuestion && !audioCache.value[nextQuestion.id]) {
      loadQuestionAudio(nextQuestion.id, nextQuestion.content, true);
    }
  }
};

// 加载问题音频的统一方法
const loadQuestionAudio = async (questionId, questionContent, isPreload = false) => {
  // 如果已经在加载中或已经缓存，则返回已存在的promise
  if (audioCachePromises.value[questionId]) {
    return audioCachePromises.value[questionId];
  }

  // 如果已经在缓存中，直接返回
  if (audioCache.value[questionId]) {
    return Promise.resolve(audioCache.value[questionId]);
  }

  // 设置加载状态
  if (!isPreload) {
    audioLoadingStates.value[questionId] = true;
  }

  // 创建新的加载promise
  const loadPromise = new Promise((resolve, reject) => {
    postJSON('/api/virtual/audio', {
      question: questionContent
    }, 
    (message, data) => {
      const audioUrl = `${BACKEND_DOMAIN}audio/${data}`;
      audioFileNames.value[questionId - 1] = data;
      audioCache.value[questionId] = audioUrl;
      audioLoadingStates.value[questionId] = false;
      resolve(audioUrl);
    },
    (errorMsg) => {
      audioLoadingStates.value[questionId] = false;
      reject(errorMsg);
    });
  });

  // 存储promise
  audioCachePromises.value[questionId] = loadPromise;

  return loadPromise;
};

// 获取问题列表
const fetchQuestions = async () => {
  loading.value = true;
  error.value = false;
  
  return new Promise((resolve, reject) => {
    postJSON('/api/virtual/question', {
      template_id: parseInt(props.thesisId),
      question_count: props.questionCount,
      difficulty: props.difficulty
    }, 
      (message, data) => {
        // 存储原始问题内容
        questionContents.value = [...data];
        
        // 后端返回的是字符串数组，需要转换为问题对象数组
        questions.value = data.map((content, index) => ({
          id: index + 1,
          content: content,
          status: 'pending',
          userAnswer: '',
          score: null,
          feedback: '',
          referenceAnswer: ''
        }));
        loading.value = false;
        resolve();
      }, 
      (errorMsg) => {
        error.value = true;
        errorMessage.value = errorMsg;
        loading.value = false;
        reject(errorMsg);
      },
      (error) => {
        error.value = true;
        errorMessage.value = '获取问题失败，请稍后重试';
        loading.value = false;
        reject(error);
      }
    );
  });
};

// 检查是否所有问题都已回答
const isAllQuestionsAnswered = computed(() => {
  return questions.value.every(q => q.userAnswer.trim() !== '');
});

// 提交所有答案
const submitAllAnswers = async () => {
  // 停止当前录音并等待完全停止
  await stopRecording();
  
  if (!isAllQuestionsAnswered.value) {
    alert('请回答所有问题后再提交');
    return;
  }

  isSubmitted.value = true;
  questions.value.forEach(question => {
    question.status = 'submitted';
  });

  // 准备所有数据（包含音频）
  await prepareSubmissionData();
};

// 准备提交数据，包括将音频转为base64
const prepareSubmissionData = async () => {
  try {
    console.log('准备提交数据 - 当前录音状态:', {
      totalQuestions: questions.value.length,
      recordedBlobs: Object.keys(userAudioBlobs.value).length,
      questionIds: questions.value.map(q => q.id),
      recordedIds: Object.keys(userAudioBlobs.value)
    });

    // 确保每个问题都有对应的录音
    const emptyAudioBlob = new Blob([], { type: 'audio/webm' });
    questions.value.forEach(question => {
      if (!userAudioBlobs.value[question.id]) {
        console.warn(`问题 ${question.id} 没有对应的录音，使用空录音代替`);
        userAudioBlobs.value[question.id] = emptyAudioBlob;
      }
    });

    // 转换音频为base64
    const audioPromises = Object.keys(userAudioBlobs.value).map(async (questionId) => {
      const blob = userAudioBlobs.value[questionId];
      return {
        questionId: parseInt(questionId),
        audioBase64: await blobToBase64(blob)
      };
    });
    
    const audioBase64Results = await Promise.all(audioPromises);
    
    // 构建答案数组，包含base64音频
    const answers = await Promise.all(questions.value.map(async q => {
      const audioData = audioBase64Results.find(a => a.questionId === q.id);
      if (!audioData) {
        console.warn(`问题 ${q.id} 未找到对应的音频数据，使用空音频`);
      }
      return {
        questionId: q.id,
        questionContent: q.content,
        answer: q.userAnswer,
        audioBase64: audioData ? audioData.audioBase64 : await blobToBase64(emptyAudioBlob)
      };
    }));

    console.log('准备提交的答案数量:', answers.length);
    console.log('每个答案是否包含音频:', answers.map(a => ({
      questionId: a.questionId,
      hasAudio: !!a.audioBase64
    })));
    
    // 发送所有数据到后端
    submitFinalData(answers);
  } catch (error) {
    console.error('准备提交数据失败:', error);
    alert('准备提交数据失败，请稍后重试');
    isSubmitted.value = false;
    questions.value.forEach(q => q.status = 'pending');
  }
};

// 将Blob转换为base64
const blobToBase64 = (blob) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => {
      // 获取base64字符串，去掉开头的数据类型标识（如"data:audio/wav;base64,"）
      const base64String = reader.result.split(',')[1];
      resolve(base64String);
    };
    reader.onerror = reject;
    reader.readAsDataURL(blob);
  });
};

// 提交最终数据到后端
const submitFinalData = (answers) => {
  // 构造最终提交的数据
  const finalData = {
    template_id: parseInt(props.thesisId),
    difficulty: props.difficulty,
    question_count: props.questionCount,
    questions: questionContents.value,
    audioFileNames: audioFileNames.value,
    answers: answers
  };
  
  // 发送数据到后端
  postJSON('/api/virtual/submit', finalData, 
    (msg, data) => {
      message.success('提交答辩数据成功');
      emit('showResult', data);
    },
    (msg,data) => {
      alert('提交答辩数据失败：' + errorMsg);
      isSubmitted.value = false;
      questions.value.forEach(q => q.status = 'pending');
    },
    (msg,dada) => {
      alert('提交答辩数据失败，请稍后重试');
      isSubmitted.value = false;
      questions.value.forEach(q => q.status = 'pending');
    }
  );
};

// 跳转到下一题
const goToNextQuestion = async () => {
  try {
    // 1. 如果正在录音，等待录音完成
    if (isRecording.value) {
      await stopRecording();
    }
    
    // 2. 确保当前问题已经回答
    if (!currentQuestion.value.userAnswer) {
      message.warning('请等待录音完成');
      return;
    }

    // 3. 切换到下一题
    if (hasNextQuestion.value) {
      // 重置状态
      resetAudioState();
      currentQuestionIndex.value++;
      
      // 加载并播放新问题
      await autoPlayCurrentQuestion();
    } else {
      // 如果是最后一题，提交所有答案
      await submitAllAnswers();
    }
  } catch (error) {
    console.error('切换问题失败:', error);
    message.error('切换问题时出错，请重试');
  }
};

// 跳转到指定题目
const goToQuestion = async (index) => {
  // 停止当前问题的录音并等待完全停止
  await stopRecording();
  
  // 在提交前，只允许查看已回答的题目或下一题
  if (!isSubmitted.value) {
    const canGoToQuestion = index <= currentQuestionIndex.value + 1 && 
                           (index === 0 || questions.value[index - 1].userAnswer.trim() !== '');
    if (canGoToQuestion) {
      currentQuestionIndex.value = index;
      // 重置音频相关状态
      resetAudioState();
      // 自动开始播放新问题
      autoPlayCurrentQuestion();
    }
  } else {
    // 提交后允许自由查看所有题目
    currentQuestionIndex.value = index;
    // 重置音频相关状态
    resetAudioState();
    // 自动开始播放新问题
    autoPlayCurrentQuestion();
  }
};

// 重置音频状态
const resetAudioState = () => {
  displayedText.value = '';
  hasPlayed.value = false;
  isPlaying.value = false;
  isRecording.value = false;
  recordingDuration.value = 0;
  
  // 停止并清理音频播放器
  if (audioPlayer.value) {
    audioPlayer.value.pause();
    audioPlayer.value.currentTime = 0;
    audioPlayer.value.src = '';
  }
  
  // 停止并清理录音
  if (mediaRecorder.value?.stream) {
    mediaRecorder.value.stream.getTracks().forEach(track => track.stop());
  }
  if (recordingTimer.value) {
    clearInterval(recordingTimer.value);
    recordingTimer.value = null;
  }
};

// 完成答辩
const finishDefense = () => {
  const totalScore = questions.value.reduce((sum, q) => sum + (q.score || 0), 0);
  const averageScore = (totalScore / questions.value.length).toFixed(1);
  
  alert(`答辩已完成！您的平均得分为: ${averageScore}/10`);
  emit('back');
};

// 获取难度标签
const getDifficultyLabel = (value) => {
  const difficultyMap = {
    'easy': '简单',
    'medium': '中等',
    'hard': '困难'
  };
  return difficultyMap[value] || value;
};

// 获取问题状态类名
const getQuestionStatusClass = (question) => {
  return {
    'status-pending': question.status === 'pending',
    'status-submitted': question.status === 'submitted',
    'status-evaluated': question.status === 'evaluated'
  };
};

// 获取问题状态文本
const getQuestionStatusText = (question) => {
  const statusMap = {
    'pending': '待回答',
    'submitted': '评估中',
    'evaluated': '已完成'
  };
  return statusMap[question.status];
};

// 是否所有问题都已评估
const isAllQuestionsEvaluated = computed(() => {
  return questions.value.length > 0 && 
         questions.value.every(q => q.status === 'evaluated');
});

// 获取音频URL并播放
const playAudio = async () => {
  if (!currentQuestion.value || hasPlayed.value) return; // 如果已经播放过，不再允许播放
  
  const questionId = currentQuestion.value.id;
  
  // 检查缓存中是否已存在该问题的音频
  if (audioCache.value[questionId]) {
    audioPlayer.value.src = audioCache.value[questionId];
    startPlayback();
    return;
  }
  
  isLoadingAudio.value = true;
  try {
    await loadQuestionAudio(questionId, currentQuestion.value.content);
    audioPlayer.value.src = audioCache.value[questionId];
    isLoadingAudio.value = false;
    startPlayback();
  } catch (error) {
    alert('获取音频失败，请稍后重试');
    isLoadingAudio.value = false;
  }
};

// 开始播放
const startPlayback = () => {
  displayedText.value = '';
  isPlaying.value = true;
  audioPlayer.value.play();
};

// 音频播放结束后自动开始录音
const onAudioEnded = () => {
  isPlaying.value = false;
  hasPlayed.value = true;
  displayedText.value = currentQuestion.value.content;
  
  // 自动开始录音
  startRecording();
};

// 音频播放进度更新
const onTimeUpdate = () => {
  if (!audioPlayer.value || !currentQuestion.value) return;
  
  const progress = audioPlayer.value.currentTime / audioPlayer.value.duration;
  const textLength = currentQuestion.value.content.length;
  const charactersToShow = Math.floor(progress * textLength);
  displayedText.value = currentQuestion.value.content.slice(0, charactersToShow);
};

// 开始录制回答
const startRecording = async () => {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    mediaRecorder.value = new MediaRecorder(stream, {
      mimeType: 'audio/webm;codecs=opus'
    });
    
    // 重置录音状态
    audioChunks.value = [];
    recordingDuration.value = 0;
    isRecording.value = true;
    
    // 收集音频数据
    mediaRecorder.value.ondataavailable = (event) => {
      if (event.data.size > 0) {
        audioChunks.value.push(event.data);
      }
    };
    
    // 设置录音结束的处理
    mediaRecorder.value.onstop = () => {
      // 保存录音数据
      if (audioChunks.value.length > 0) {
        const audioBlob = new Blob(audioChunks.value, { type: 'audio/webm' });
        if (currentQuestion.value) {
          userAudioBlobs.value[currentQuestion.value.id] = audioBlob;
          currentQuestion.value.userAnswer = `已完成录音 - 时长：${Math.floor(recordingDuration.value)}秒`;
        }
      }
      
      // 清理资源
      if (recordingTimer.value) {
        clearInterval(recordingTimer.value);
        recordingTimer.value = null;
      }
      if (mediaRecorder.value?.stream) {
        mediaRecorder.value.stream.getTracks().forEach(track => track.stop());
      }
      
      isRecording.value = false;
    };
    
    // 开始录音
    mediaRecorder.value.start(100);
    
    // 开始计时，30秒后自动停止
    recordingTimer.value = setInterval(() => {
      recordingDuration.value += 0.1;
      if (recordingDuration.value >= 30) {
        stopRecording();
      }
    }, 100);
    
  } catch (error) {
    console.error('启动录音失败:', error);
    message.error('无法启动录音，请检查麦克风权限');
    if (currentQuestion.value) {
      currentQuestion.value.userAnswer = '录音功能不可用';
    }
    isRecording.value = false;
  }
};

// 停止录制回答
const stopRecording = () => {
  return new Promise((resolve) => {
    if (!mediaRecorder.value || !isRecording.value) {
      resolve();
      return;
    }

    try {
      // 添加一次性事件监听器来处理录音结束
      mediaRecorder.value.addEventListener('stop', () => {
        resolve();
      }, { once: true });
      
      // 停止录音
      mediaRecorder.value.stop();
      
    } catch (error) {
      console.error('停止录音时出错:', error);
      isRecording.value = false;
      resolve();
    }
  });
};

// 自动播放当前问题
const autoPlayCurrentQuestion = async () => {
  if (!currentQuestion.value) return;
  
  try {
    // 1. 重置状态
    hasPlayed.value = false;
    displayedText.value = '';
    
    // 2. 开始加载并播放当前问题的音频
    console.log('开始加载并播放问题:', currentQuestion.value.id);
    await playAudio();
    
    // 3. 预加载下一个问题的音频（如果有）
    if (hasNextQuestion.value) {
      await preloadNextQuestionAudio();
    }
  } catch (error) {
    console.error('自动播放失败:', error);
    message.error('播放问题音频失败，请刷新页面重试');
  }
};

// 组件挂载时获取问题
onMounted(async () => {
  try {
    await fetchQuestions();
    console.log('问题列表获取成功，开始自动播放');
    // 确保问题列表已加载完成后再开始播放
    if (questions.value.length > 0) {
      await autoPlayCurrentQuestion();
    }
  } catch (error) {
    console.error('初始化失败:', error);
  }
});
</script>

<style scoped>
.virtual-defense {
  width: 100%;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-container {
  text-align: center;
  padding: 40px;
  background-color: #fff3f3;
  border-radius: 8px;
}

.defense-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.defense-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.defense-info {
  display: flex;
  gap: 20px;
  color: #666;
}

.questions-container {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.question-item {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.question-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.question-status {
  padding: 5px 10px;
  border-radius: 15px;
  font-size: 12px;
}

.status-pending {
  background-color: #ffeeba;
  color: #856404;
}

.status-submitted {
  background-color: #b8daff;
  color: #004085;
}

.status-evaluated {
  background-color: #c3e6cb;
  color: #155724;
}

.question-content {
  line-height: 1.6;
}

.audio-text-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.audio-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.audio-controls audio {
  display: none;
}

.play-button {
  background-color: #28a745;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: background-color 0.3s ease;
}

.play-button:hover {
  background-color: #218838;
}

.play-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.playing-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #28a745;
  font-weight: bold;
}

.playing-indicator::after {
  content: '';
  display: inline-block;
  width: 12px;
  height: 12px;
  background-color: #28a745;
  border-radius: 50%;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0% {
    transform: scale(0.95);
    opacity: 0.5;
  }
  50% {
    transform: scale(1.1);
    opacity: 1;
  }
  100% {
    transform: scale(0.95);
    opacity: 0.5;
  }
}

/* 录音相关样式 */
.recording-section {
  margin: 15px 0;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.recording-status {
  display: flex;
  align-items: center;
  padding: 10px;
  border-radius: 6px;
  margin-bottom: 10px;
}

.active-recording {
  background-color: #ffeeee;
  border: 1px solid #ffcccc;
}

.recording-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #d9534f;
  font-weight: bold;
}

.recording-dot {
  display: inline-block;
  width: 12px;
  height: 12px;
  background-color: #d9534f;
  border-radius: 50%;
  animation: record-pulse 1s infinite;
}

@keyframes record-pulse {
  0% {
    transform: scale(0.95);
    opacity: 0.8;
  }
  50% {
    transform: scale(1.2);
    opacity: 1;
  }
  100% {
    transform: scale(0.95);
    opacity: 0.8;
  }
}

.recognized-text {
  background-color: #f0f8ff;
  border-radius: 6px;
  padding: 10px 15px;
  margin-top: 10px;
}

.recognized-label {
  font-weight: bold;
  margin-bottom: 5px;
  color: #0056b3;
}

.recognized-content {
  white-space: pre-line;
  line-height: 1.5;
}

.answer-input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin: 15px 0;
  font-family: inherit;
  resize: vertical;
}

.question-actions {
  display: flex;
  justify-content: flex-end;
}

.submit-button {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}

.submit-button:hover {
  background-color: #0069d9;
}

.user-answer {
  margin: 15px 0;
  padding: 15px;
  background-color: #f0f8ff;
  border-radius: 4px;
}

.answer-label, .feedback-label, .reference-label {
  font-weight: bold;
  margin-bottom: 5px;
}

.evaluation-loading {
  text-align: center;
  padding: 20px;
  color: #666;
  font-style: italic;
}

.evaluation-result {
  margin-top: 20px;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
}

.evaluation-score {
  font-size: 18px;
  margin-bottom: 15px;
}

.evaluation-feedback, .evaluation-reference {
  margin-top: 15px;
}

.feedback-content, .reference-content {
  line-height: 1.6;
  margin-top: 5px;
  white-space: pre-line;
}

.defense-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.action-button {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
}

.back-button {
  background-color: #f8f9fa;
  color: #333;
  border: 1px solid #ddd;
}

.retry-button {
  background-color: #17a2b8;
  color: white;
}

.finish-button {
  background-color: #28a745;
  color: white;
}

.action-button:hover {
  opacity: 0.9;
}

.progress-indicator {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-top: 20px;
}

.progress-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background-color: #ddd;
  cursor: pointer;
  transition: all 0.3s ease;
}

.progress-dot.active {
  background-color: #007bff;
  transform: scale(1.2);
}

.progress-dot.completed {
  background-color: #28a745;
}

.progress-dot.answered {
  background-color: #ffc107;
}

.next-button {
  background-color: #28a745;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 15px;
}

.next-button:hover {
  background-color: #218838;
}

.next-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.submit-all-button {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
}

.submit-all-button:hover {
  background-color: #0056b3;
}

.submit-all-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.typewriter-text {
  font-family: monospace;
  white-space: pre-wrap;
  line-height: 1.6;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
  min-height: 100px;
  width: 100%;
}

.stop-recording-button {
  background-color: #dc3545;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  margin-top: 15px;
}

.stop-recording-button:hover {
  background-color: #c82333;
}

.loading-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #666;
  font-style: italic;
}

.loading-indicator::after {
  content: '';
  display: inline-block;
  width: 12px;
  height: 12px;
  border: 2px solid #666;
  border-top-color: transparent;
  border-radius: 50%;
  animation: loading-spin 1s linear infinite;
}

@keyframes loading-spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style> 