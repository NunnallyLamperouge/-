<template>
  <div class="mission-container">
    <a-row>
      <a-col :span="24">
        <a-card title="任务管理" :bordered="false">
          <template #extra>
            <a-button type="primary" @click="showAddModal">创建任务</a-button>
          </template>
          <a-table
            :dataSource="missions"
            :columns="columns"
            :loading="loading"
            :pagination="{ showSizeChanger: true, pageSizeOptions: ['10', '20', '50', '100'] }"
            rowKey="id"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="getStatusColor(record.status)">{{ getStatusText(record.status) }}</a-tag>
              </template>
              <template v-if="column.key === 'priority'">
                <a-tag :color="getPriorityColor(record.priority)">{{ getPriorityText(record.priority) }}</a-tag>
              </template>
              <template v-if="column.key === 'progress'">
                <a-progress :percent="record.progress" size="small" />
              </template>
              <template v-if="column.key === 'student'">
                <a-avatar size="small" style="margin-right: 8px" />
                {{ record.Student.username }}
              </template>
              <template v-if="column.key === 'endTime'">
                <span :style="{ color: isOverdue(record.endTime) ? '#ff4d4f' : '' }">
                  {{ formatDate(record.endTime) }}
                </span>
              </template>
              <template v-if="column.key === 'action'">
                <a-space>
                  <a-button type="link" @click="showDetailModal(record)">详情</a-button>
                  <a-button type="link" @click="router.push('/view-template/' + record.templateId)">审核</a-button>
                  <a-button type="link" @click="showFeedbackModal(record)">反馈</a-button>
                  <a-popconfirm
                    title="确定要删除此任务吗？"
                    ok-text="确定"
                    cancel-text="取消"
                    @confirm="deleteMission(record.id)"
                  >
                    <a-button type="link" danger>删除</a-button>
                  </a-popconfirm>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 添加任务模态框 -->
    <a-modal
      v-model:visible="addModalVisible"
      title="创建新任务"
      @ok="handleAddMission"
      :confirmLoading="submitLoading"
      width="700px"
    >
      <a-form :model="missionForm" :label-col="{ span: 4 }" :wrapper-col="{ span: 20 }">
        <a-form-item label="任务名称" name="name" :rules="[{ required: true, message: '请输入任务名称' }]">
          <a-input v-model:value="missionForm.name" placeholder="请输入任务名称" />
        </a-form-item>
        <a-form-item label="任务描述" name="description" :rules="[{ required: true, message: '请输入任务描述' }]">
          <a-textarea v-model:value="missionForm.description" :rows="4" placeholder="请输入任务描述" />
        </a-form-item>
        <a-form-item label="选择模板" name="templateId" :rules="[{ required: true, message: '请选择模板' }]">
          <a-select v-model:value="missionForm.templateId" placeholder="请选择模板">
            <a-select-option v-for="template in templates" :key="template.id" :value="template.id">
              {{ template.title }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="开始时间" name="startTime" :rules="[{ required: true, message: '请选择开始时间' }]">
          <a-date-picker
            v-model:value="missionForm.startTime"
            show-time
            format="YYYY-MM-DD HH:mm:ss"
            placeholder="请选择开始时间"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="截止时间" name="endTime" :rules="[{ required: true, message: '请选择截止时间' }]">
          <a-date-picker
            v-model:value="missionForm.endTime"
            show-time
            format="YYYY-MM-DD HH:mm:ss"
            placeholder="请选择截止时间"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="优先级" name="priority">
          <a-select v-model:value="missionForm.priority" placeholder="请选择优先级">
            <a-select-option :value="0">正常</a-select-option>
            <a-select-option :value="1">重要</a-select-option>
            <a-select-option :value="2">紧急</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="学生" name="studentIds" :rules="[{ required: true, message: '请选择学生' }]">
          <a-select
            v-model:value="missionForm.studentIds"
            mode="multiple"
            placeholder="请选择学生"
            :options="studentOptions"
            :fieldNames="{ label: 'label', value: 'value' }"
          ></a-select>
        </a-form-item>
        <a-form-item label="标签" name="tags">
          <a-select v-model:value="missionForm.tags" mode="tags" placeholder="输入标签后按回车">
          </a-select>
        </a-form-item>
        <a-form-item label="附件" name="attachments">
          <a-upload
            v-model:fileList="fileList"
            action="/api/upload"
            :headers="{ Authorization: getAuthToken() }"
            @change="handleUploadChange"
          >
            <a-button>
              <upload-outlined />
              上传附件
            </a-button>
          </a-upload>
        </a-form-item>
      </a-form>
    </a-modal>

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
        <a-descriptions-item label="学生">
          <a-avatar :src="currentMission.student?.avatar" size="small" style="margin-right: 8px" v-if="currentMission.student" />
          {{ currentMission.Student?.username }}
        </a-descriptions-item>
        <a-descriptions-item label="开始时间">{{ formatDate(currentMission.startTime) }}</a-descriptions-item>
        <a-descriptions-item label="截止时间">{{ formatDate(currentMission.endTime) }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentMission.status)">{{ getStatusText(currentMission.status) }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="优先级">
          <a-tag :color="getPriorityColor(currentMission.priority)">{{ getPriorityText(currentMission.priority) }}</a-tag>
        </a-descriptions-item>
        <!-- <a-descriptions-item label="任务类型">{{ getTypeText(currentMission.type) }}</a-descriptions-item> -->
        <a-descriptions-item label="进度">
          <a-progress :percent="currentMission.progress" />
        </a-descriptions-item>
        <a-descriptions-item label="标签" v-if="currentMission.tags">
          <template v-if="currentMission.tags">
            <a-tag v-for="(tag, index) in parseTags(currentMission.tags)" :key="index" color="blue">{{ tag }}</a-tag>
          </template>
        </a-descriptions-item>
        <a-descriptions-item label="附件" v-if="currentMission.attachments">
          <template v-if="currentMission.attachments">
            <a-space direction="vertical">
              <a-button
                v-for="(attachment, index) in parseAttachments(currentMission.attachments)"
                :key="index"
                type="link"
                @click="downloadAttachment(attachment)"
              >
                {{ attachment.name }}
              </a-button>
            </a-space>
          </template>
        </a-descriptions-item>
        <a-descriptions-item label="反馈" v-if="currentMission.feedback">
          {{ currentMission.feedback }}
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>

    <!-- 反馈模态框 -->
    <a-modal
      v-model:visible="feedbackModalVisible"
      title="提供反馈"
      @ok="handleSubmitFeedback"
      :confirmLoading="submitLoading"
    >
      <a-form :model="feedbackForm" :label-col="{ span: 4 }" :wrapper-col="{ span: 20 }">
        <a-form-item label="任务状态" name="status">
          <a-select v-model:value="feedbackForm.status">
            <a-select-option :value="MissionStatus.InProgress">进行中</a-select-option>
            <a-select-option :value="MissionStatus.Submitted">待审核</a-select-option>
            <a-select-option :value="MissionStatus.Rejected">打回修改</a-select-option>
            <a-select-option :value="MissionStatus.Completed">已完成</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item 
          v-if="feedbackForm.status === MissionStatus.Rejected" 
          label="打回原因" 
          name="rejectionReason"
          :rules="[{ required: true, message: '请输入打回原因' }]"
        >
          <a-textarea v-model:value="feedbackForm.rejectionReason" :rows="4" placeholder="请输入打回原因" />
        </a-form-item>
        <a-form-item label="反馈内容" name="feedback">
          <a-textarea v-model:value="feedbackForm.feedback" :rows="6" placeholder="请输入反馈内容" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue';
import { message } from 'ant-design-vue';
import { get, postJSON, del, putJSON } from '../../utility/request';
import dayjs from 'dayjs';
import { UploadOutlined } from '@ant-design/icons-vue';
import router from '@/router';

// 状态变量
const missions = ref([]);
const loading = ref(false);
const addModalVisible = ref(false);
const detailModalVisible = ref(false);
const feedbackModalVisible = ref(false);
const submitLoading = ref(false);
const currentMission = ref({});
const fileList = ref([]);
const studentOptions = ref([]);
const templates = ref([]);

// 表单数据
const missionForm = reactive({
  name: '',
  description: '',
  studentIds: [],
  startTime: '',
  endTime: '',
  priority: 0,
  type: 0,
  tags: [],
  attachments: '',
  templateId: undefined,
});

const feedbackForm = reactive({
  status: 0,
  feedback: '',
  rejectionReason: ''
});

// 表格列定义
const columns = [
  {
    title: '任务名称',
    dataIndex: 'name',
    key: 'name',
    ellipsis: true,
  },
  {
    title: '学生',
    dataIndex: 'student',
    key: 'student',
  },
  {
    title: '模板',
    dataIndex: ['template', 'title'],
    key: 'template',
  },
  {
    title: '优先级',
    dataIndex: 'priority',
    key: 'priority',
  },
  {
    title: '进度',
    dataIndex: 'progress',
    key: 'progress',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
  },
  {
    title: '截止时间',
    dataIndex: 'endTime',
    key: 'endTime',
  },
  {
    title: '操作',
    key: 'action',
  },
];

// 生命周期钩子
onMounted(() => {
  fetchMissions();
  fetchStudents();
  fetchMissionStats();
  fetchTemplates();
});

// 获取任务列表
const fetchMissions = () => {
  loading.value = true;
  get('/api/mission/teacher', null, 
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

// 获取任务统计信息
const fetchMissionStats = () => {
  get('/api/mission/statistics', null,
    (_, data) => {
      // 可以在这里处理统计数据
      console.log('任务统计:', data);
    },
    (msg) => {
      message.error(msg || '获取任务统计失败');
    }
  );
};

// 获取学生列表
const fetchStudents = () => {
  get('/api/guide/teacher/list', null, 
    (_, data) => {
      studentOptions.value = (data || []).map(student => ({
        label: student.Student.username,
        value: student.Student.id
      }));
    },
    (msg) => {
      message.error(msg || '获取学生列表失败');
    }
  );
};

// 获取模板列表
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

// 展示添加模态框
const showAddModal = () => {
  // 重置表单
  Object.assign(missionForm, {
    name: '',
    description: '',
    studentIds: [],
    startTime: '',
    endTime: '',
    priority: 0,
    type: 0,
    tags: [],
    attachments: '',
    templateId: undefined,
  });
  fileList.value = [];
  addModalVisible.value = true;
};

// 展示详情模态框
const showDetailModal = (record) => {
  currentMission.value = {...record};
  detailModalVisible.value = true;
};

// 展示反馈模态框
const showFeedbackModal = (record) => {
  currentMission.value = {...record};
  feedbackForm.status = record.status;
  feedbackForm.feedback = record.feedback || '';
  feedbackModalVisible.value = true;
};

// 处理添加任务
const handleAddMission = () => {
  if (!missionForm.name || !missionForm.description || !missionForm.studentIds.length || !missionForm.startTime || !missionForm.endTime) {
    message.warning('请填写必填字段');
    return;
  }
  
  submitLoading.value = true;
  
  const payload = {
    name: missionForm.name,
    description: missionForm.description,
    studentIds: missionForm.studentIds,
    startTime: dayjs(missionForm.startTime).format('YYYY-MM-DD HH:mm:ss'),
    endTime: dayjs(missionForm.endTime).format('YYYY-MM-DD HH:mm:ss'),
    priority: missionForm.priority,
    type: missionForm.type,
    tags: missionForm.tags.join(','),
    attachments: JSON.stringify(fileList.value.map(file => {
      return {
        uid: file.uid,
        name: file.name,
        url: file.response?.url || '',
        status: file.status
      };
    }).filter(file => file.status === 'done' && file.url)),
    templateId: missionForm.templateId
  };
  
  postJSON('/api/mission/create', payload, 
    (msg) => {
      message.success(msg || '任务创建成功');
      addModalVisible.value = false;
      submitLoading.value = false;
      fetchMissions();
    },
    (msg) => {
      message.error(msg || '任务创建失败');
      submitLoading.value = false;
    }
  );
};

// 处理提交反馈
const handleSubmitFeedback = () => {
  submitLoading.value = true;
  
  const payload = {
    status: feedbackForm.status,
    feedback: feedbackForm.feedback,
    rejectionReason: feedbackForm.status === MissionStatus.Rejected ? feedbackForm.rejectionReason : '',
    progress: currentMission.value.progress
  };
  
  putJSON(`/api/mission/${currentMission.value.id}/status`, payload, 
    (msg) => {
      message.success(msg || '反馈提交成功');
      feedbackModalVisible.value = false;
      submitLoading.value = false;
      fetchMissions();
    },
    (msg) => {
      message.error(msg || '反馈提交失败');
      submitLoading.value = false;
    }
  );
};

// 删除任务
const deleteMission = (id) => {
  del(`/api/mission/${id}`, null, 
    (msg) => {
      message.success(msg || '任务删除成功');
      fetchMissions();
    },
    (msg) => {
      message.error(msg || '任务删除失败');
    }
  );
};

// 处理上传变更
const handleUploadChange = (info) => {
  fileList.value = [...info.fileList];
};

// 下载附件
const downloadAttachment = (attachment) => {
  window.open(attachment.url, '_blank');
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

// 获取类型文本
const getTypeText = (type) => {
  const texts = ['普通任务', '论文任务', '实验任务', '调研任务'];
  return texts[type] || '未知类型';
};

// 格式化日期
const formatDate = (date) => {
  if (!date) return '';
  return dayjs(date).format('YYYY-MM-DD HH:mm');
};

// 检查是否过期
const isOverdue = (date) => {
  if (!date) return false;
  return dayjs(date).isBefore(dayjs());
};

// 解析标签
const parseTags = (tags) => {
  if (!tags) return [];
  return tags.split(',').filter(tag => tag);
};

// 解析附件
const parseAttachments = (attachments) => {
  if (!attachments) return [];
  try {
    return JSON.parse(attachments);
  } catch {
    return [];
  }
};

// 获取认证Token
const getAuthToken = () => {
  return localStorage.getItem('authToken') || '';
};

// 任务状态枚举
const MissionStatus = {
  Draft: 1,
  InProgress: 2,
  Submitted: 3,
  Rejected: 4,
  Completed: 5
};
</script>

<style scoped>
.mission-container {
  padding: 20px;
}

:deep(.ant-descriptions-item-label) {
  width: 120px;
}
</style>