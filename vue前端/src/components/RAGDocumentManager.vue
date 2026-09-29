<template>
  <div class="rag-document-manager">
    <div class="upload-section">
      <a-upload
        :customRequest="handleUpload"
        :showUploadList="false"
        accept=".pdf,.doc,.docx,.txt"
      >
        <a-button type="primary">
          <upload-outlined></upload-outlined>
          上传文档
        </a-button>
      </a-upload>
    </div>

    <a-table
      :columns="columns"
      :dataSource="documents"
      :loading="loading"
      class="document-table"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="handleDelete(record.fileName)">
              删除
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { UploadOutlined } from '@ant-design/icons-vue';
import { post, get, del } from '../utility/request';

const documents = ref([]);
const loading = ref(false);

const columns = [
  {
    title: '文件名',
    dataIndex: 'fileName',
    key: 'fileName',
  },
  {
    title: '上传时间',
    dataIndex: 'uploadTime',
    key: 'uploadTime',
  },
  {
    title: '操作',
    key: 'action',
  },
];

const fetchDocuments = () => {
  loading.value = true;
  get('/api/rag/documents', null, 
    (msg, data) => {
      documents.value = data;
      loading.value = false;
    },
    (msg) => {
      message.error(msg);
      loading.value = false;
    }
  );
};

const handleUpload = async ({ file }) => {
  const formData = new FormData();
  formData.append('file', file);
  
  post('/api/rag/documents', formData,
    (msg) => {
      message.success(msg);
      fetchDocuments();
    },
    (msg) => {
      message.error(msg);
    }
  );
};

const handleDelete = (fileName) => {
  del(`/api/rag/documents/${fileName}`, null,
    (msg) => {
      message.success(msg);
      fetchDocuments();
    },
    (msg) => {
      message.error(msg);
    }
  );
};

onMounted(() => {
  fetchDocuments();
});
</script>

<style scoped>
.rag-document-manager {
  padding: 20px;
}

.upload-section {
  margin-bottom: 20px;
}

.document-table {
  margin-top: 20px;
}
</style> 