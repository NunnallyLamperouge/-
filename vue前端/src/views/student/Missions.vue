<template>
  <div class="mission-container">
    <a-row :gutter="[16, 16]">
      <a-col :span="24">
        <a-card title="我的任务" :bordered="false">
          <template #extra>
            <a-space>
              <a-select v-model:value="filterStatus" style="width: 120px" placeholder="任务状态">
                <a-select-option :value="-1">全部状态</a-select-option>
                <a-select-option :value="MissionStatus.Draft">草稿</a-select-option>
                <a-select-option :value="MissionStatus.InProgress">进行中</a-select-option>
                <a-select-option :value="MissionStatus.Submitted">待审核</a-select-option>
                <a-select-option :value="MissionStatus.Rejected">已打回</a-select-option>
                <a-select-option :value="MissionStatus.Completed">已完成</a-select-option>
              </a-select>
              <a-select v-model:value="filterPriority" style="width: 120px" placeholder="优先级">
                <a-select-option :value="-1">全部优先级</a-select-option>
                <a-select-option :value="0">正常</a-select-option>
                <a-select-option :value="1">重要</a-select-option>
                <a-select-option :value="2">紧急</a-select-option>
              </a-select>
            </a-space>
          </template>
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :sm="12" :md="8" :lg="6" v-for="mission in filteredMissions" :key="mission.id">
              <a-card :class="['mission-card', getStatusClass(mission.status)]" hoverable>
                <template #cover>
                  <div class="mission-header" :class="getPriorityClass(mission.priority)">
                    <h3>{{ mission.name }}</h3>
                    <a-tag>{{ mission.template?.title }}</a-tag>
                  </div>
                </template>
                <a-card-meta>
                  <template #description>
                    <div class="mission-content">
                      <p class="description">{{ mission.description }}</p>
                      <div class="mission-info">
                        <p><clock-circle-outlined /> 截止时间：{{ formatDate(mission.endTime) }}</p>
                        <p><user-outlined /> 导师：{{ mission.Teacher?.username }}</p>
                      </div>
                      <div class="mission-footer">
                        <a-space>
                          <a-tag :color="getStatusColor(mission.status)">{{ getStatusText(mission.status) }}</a-tag>
                          <a-button type="link" @click="showDetailModal(mission)">详情</a-button>
                          <a-button 
                            type="link" 
                            @click="showSubmitModal(mission)" 
                            :disabled="mission.status === MissionStatus.Submitted || mission.status === MissionStatus.Completed"
                          >
                            提交
                          </a-button>
                        </a-space>
                      </div>
                    </div>
                  </template>
                </a-card-meta>
              </a-card>
            </a-col>
          </a-row>
        </a-card>
      </a-col>
    </a-row>

    <!-- 任务详情模态框 -->
    <a-modal
      v-model:visible="detailModalVisible"
      title="任务详情"
      :footer="null"
      width="700px"
    >
      <a-descriptions :column="1" bordered>
        <a-descriptions-item label="任务名称">{{ currentMission.name }}</a-descriptions-item>
        <a-descriptions-item label="任务描述">{{ currentMission.description }}</a-descriptions-item>
        <a-descriptions-item label="导师">
          <a-avatar :src="currentMission.Teacher?.avatar" size="small" style="margin-right: 8px" />
          {{ currentMission.Teacher?.username }}
        </a-descriptions-item>
        <a-descriptions-item label="模板">{{ currentMission.template?.title }}</a-descriptions-item>
        <a-descriptions-item label="开始时间">{{ formatDate(currentMission.startTime) }}</a-descriptions-item>
        <a-descriptions-item label="截止时间">{{ formatDate(currentMission.endTime) }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentMission.status)">{{ getStatusText(currentMission.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="优先级">
          <a-tag :color="getPriorityColor(currentMission.priority)">{{ getPriorityText(currentMission.priority) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="提交内容" v-if="currentMission.content">
          {{ currentMission.content }}
        </a-descriptions-item>
        <a-descriptions-item label="导师反馈" v-if="currentMission.feedback">
          {{ currentMission.feedback }}
        </a-descriptions-item>
        <a-descriptions-item label="打回原因" v-if="currentMission.rejectionReason">
          <a-alert type="warning" :message="currentMission.rejectionReason" />
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 提交任务模态框 -->
    <a-modal
      v-model:visible="submitModalVisible"
      title="提交任务"
      @ok="handleSubmitMission"
      :confirmLoading="submitLoading"
    >
      <a-form :model="submitForm" :label-col="{ span: 4 }" :wrapper-col="{ span: 20 }">
        <a-form-item label="提交内容" name="content" :rules="[{ required: true, message: '请输入提交内容' }]">
          <a-textarea v-model:value="submitForm.content" :rows="10" placeholder="请详细描述您的任务完成情况" />
        </a-form-item>
        <a-form-item label="选择模板" name="templateId">
          <a-select v-model:value="submitForm.templateId" placeholder="选择模板">
            <a-select-option v-for="template in templates" :key="template.id" :value="template.id">{{ template.title }}</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { ClockCircleOutlined, UserOutlined } from '@ant-design/icons-vue';
import dayjs from 'dayjs';
import { get, postJSON } from '../../utility/request';

// 任务状态枚举
const MissionStatus = {
  Draft: 1,
  InProgress: 2,
  Submitted: 3,
  Rejected: 4,
  Completed: 5
};

// 状态变量
const missions = ref([]);
const loading = ref(false);
const detailModalVisible = ref(false);
const submitModalVisible = ref(false);
const submitLoading = ref(false);
const currentMission = ref({});
const filterStatus = ref(-1);
const filterPriority = ref(-1);
const templates = ref([]);

// 提交表单数据
const submitForm = ref({
  content: '',
  templateId: undefined
});

// 计算过滤后的任务列表
const filteredMissions = computed(() => {
  return missions.value.filter(mission => {
    const statusMatch = filterStatus.value === -1 || mission.status === filterStatus.value;
    const priorityMatch = filterPriority.value === -1 || mission.priority === filterPriority.value;
    return statusMatch && priorityMatch;
  });
});

// 生命周期钩子
onMounted(() => {
  fetchMissions();
  fetchTemplates();
});

// 获取任务列表
const fetchMissions = () => {
  loading.value = true;
  get('/api/mission/student', null,
    (_, data) => {
      missions.value = data || [];
      loading.value = false;
    },
    (msg) => {
      message.error(msg || '获取任务列表失败');
      loading.value = false;
    }
  );
};

// 获取学生的模板列表
const fetchTemplates = () => {
  get('/api/template/list', null,
    (_, data) => {
      templates.value = data || [];
    },
    (msg) => {
      message.error(msg || '获取模板列表失败');
    }
  );
};

// 展示详情模态框
const showDetailModal = (mission) => {
  currentMission.value = {...mission};
  detailModalVisible.value = true;
};

// 展示提交模态框
const showSubmitModal = (mission) => {
  currentMission.value = {...mission};
  submitForm.value = {
    content: mission.content || '',
    templateId: undefined
  };
  submitModalVisible.value = true;
};

// 处理任务提交
const handleSubmitMission = () => {
  if (!submitForm.value.content) {
    message.warning('请填写提交内容');
    return;
  }
  if (!submitForm.value.templateId) {
    message.warning('请选择提交的模板');
    return;
  }

  submitLoading.value = true;
  postJSON(`/api/mission/${currentMission.value.id}/submit`, submitForm.value,
    (msg) => {
      message.success(msg || '提交任务成功');
      submitModalVisible.value = false;
      submitLoading.value = false;
      fetchMissions();
    },
    (msg) => {
      message.error(msg || '提交任务失败');
      submitLoading.value = false;
    }
  );
};

// 获取状态颜色
const getStatusColor = (status) => {
  const colors = {
    [MissionStatus.Draft]: 'default',
    [MissionStatus.InProgress]: 'processing',
    [MissionStatus.Submitted]: 'warning',
    [MissionStatus.Rejected]: 'error',
    [MissionStatus.Completed]: 'success'
  };
  return colors[status] || 'default';
};

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    [MissionStatus.Draft]: '草稿',
    [MissionStatus.InProgress]: '进行中',
    [MissionStatus.Submitted]: '待审核',
    [MissionStatus.Rejected]: '已打回',
    [MissionStatus.Completed]: '已完成'
  };
  return texts[status] || '未知状态';
};

// 获取状态类名
const getStatusClass = (status) => {
  const classes = {
    [MissionStatus.Draft]: 'status-draft',
    [MissionStatus.InProgress]: 'status-in-progress',
    [MissionStatus.Submitted]: 'status-pending',
    [MissionStatus.Rejected]: 'status-rejected',
    [MissionStatus.Completed]: 'status-completed'
  };
  return classes[status] || '';
};

// 获取优先级颜色
const getPriorityColor = (priority) => {
  const colors = ['default', 'orange', 'red'];
  return colors[priority] || 'default';
};

// 获取优先级文本
const getPriorityText = (priority) => {
  const texts = ['正常', '重要', '紧急'];
  return texts[priority] || '未知优先级';
};

// 获取优先级类名
const getPriorityClass = (priority) => {
  const classes = ['priority-normal', 'priority-important', 'priority-urgent'];
  return classes[priority] || '';
};

// 格式化日期
const formatDate = (date) => {
  if (!date) return '';
  return dayjs(date).format('YYYY-MM-DD HH:mm');
};
</script>

<style scoped>
.mission-container {
  padding: 20px;
}

.mission-card {
  height: 100%;
  transition: all 0.3s;
}

.mission-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.mission-header {
  padding: 16px;
  border-radius: 8px 8px 0 0;
}

.mission-header h3 {
  margin: 0 0 8px 0;
  color: #fff;
}

.mission-content {
  padding: 8px 0;
}

.mission-info {
  margin: 12px 0;
  color: #666;
}

.mission-footer {
  margin-top: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.description {
  margin: 0;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  color: #666;
}

/* 优先级样式 */
.priority-normal {
  background: #1890ff;
}

.priority-important {
  background: #fa8c16;
}

.priority-urgent {
  background: #f5222d;
}

/* 状态样式 */
.status-draft {
  border-left: 3px solid #d9d9d9;
}

.status-in-progress {
  border-left: 3px solid #1890ff;
}

.status-pending {
  border-left: 3px solid #faad14;
}

.status-rejected {
  border-left: 3px solid #ff4d4f;
}

.status-completed {
  border-left: 3px solid #52c41a;
}

:deep(.ant-descriptions-item-label) {
  width: 120px;
}
</style>
