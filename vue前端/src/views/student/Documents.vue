<template>
  <div
    class="min-h-screen bg-gradient-to-br rounded-3xl from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800"
    v-motion-fade
  >
    <div class="max-w-7xl mx-auto">
      <!-- 头部区域 -->
      <div class="flex justify-between items-center" v-motion-slide-left>
        <div class="flex items-center space-x-4">
          <a-button
            v-if="!workspaceDetails || workspaceDetails.id === 0"
            type="primary"
            @click="handleCreateWorkspace"
            class="flex items-center space-x-2 bg-gradient-to-r from-blue-500 to-indigo-500 hover:from-blue-600 hover:to-indigo-600"
          >
            <plus-outlined></plus-outlined>
            <span>创建工作台</span>
          </a-button>
        </div>
      </div>

      <!-- 标签页区域 -->
      <a-tabs
        v-model:activeKey="activeTab"
        v-motion-slide-right
        class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6"
      >
        <a-tab-pane key="workspace" tab="工作台文档">
          <template v-if="workspaceDetails && workspaceDetails.id !== 0">
            <!-- 工作台详情 -->
            <div class="mb-8" v-motion-fade>
              <a-descriptions
                title="工作台详情"
                bordered
                class="bg-white dark:bg-gray-800 rounded-lg overflow-hidden"
              >
                <a-descriptions-item
                  label="工作台名称"
                  class="text-gray-700 dark:text-gray-300"
                >
                  {{ workspaceDetails.name }}
                </a-descriptions-item>
                <a-descriptions-item
                  label="创建时间"
                  class="text-gray-700 dark:text-gray-300"
                >
                  {{ parseDateTime(workspaceDetails.createdAt) }}
                </a-descriptions-item>
                <a-descriptions-item
                  label="最后更新时间"
                  class="text-gray-700 dark:text-gray-300"
                >
                  {{ parseDateTime(workspaceDetails.lastUpdatedAt) }}
                </a-descriptions-item>
                <a-descriptions-item
                  label="OpenAI 温度"
                  class="text-gray-700 dark:text-gray-300"
                >
                  {{ workspaceDetails.openAiTemp }}
                </a-descriptions-item>
                <a-descriptions-item
                  label="历史记录数"
                  class="text-gray-700 dark:text-gray-300"
                >
                  {{ workspaceDetails.openAiHistory }}
                </a-descriptions-item>
                <a-descriptions-item
                  label="相似度阈值"
                  class="text-gray-700 dark:text-gray-300"
                >
                  {{ workspaceDetails.similarityThreshold }}
                </a-descriptions-item>
                <a-descriptions-item
                  label="Top N"
                  class="text-gray-700 dark:text-gray-300"
                >
                  {{ workspaceDetails.topN }}
                </a-descriptions-item>
                <a-descriptions-item
                  label="聊天模式"
                  class="text-gray-700 dark:text-gray-300"
                >
                  {{ workspaceDetails.chatMode }}
                </a-descriptions-item>
                <a-descriptions-item
                  label="向量搜索模式"
                  class="text-gray-700 dark:text-gray-300"
                >
                  {{ workspaceDetails.vectorSearchMode }}
                </a-descriptions-item>
              </a-descriptions>
            </div>
            <div class="flex flex-nowrap gap-2">
              <a-button
                type="link"
                @click="fetchWorkspaceDetails"
                class="text-blue-500 hover:text-blue-600"
              >
                刷新工作台
              </a-button>
              <div
              v-if="selectedWorkspaceItems.length > 0"
              class="mb-6"
              v-motion-slide-up
            >
              <a-button
                type="primary"
                danger
                @click="handleBatchDeleteFromWorkspace"
                class="bg-red-500 hover:bg-red-600"
              >
                批量移除 ({{ selectedWorkspaceItems.length }})
              </a-button>
            </div>
            </div>
          

            <!-- 文档网格 -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <a-card
                v-for="doc in workspaceDocuments"
                :key="doc.docId"
                class="document-card bg-white dark:bg-gray-800 rounded-xl shadow-md hover:shadow-xl transition-all duration-300"
                :class="{ 'border-2 border-blue-500': doc.pinned }"
                v-motion
                :initial="{ opacity: 0, y: 20 }"
                :enter="{ opacity: 1, y: 0, transition: { duration: 500 } }"
              >
                <template #extra>
                  <div class="flex items-center space-x-2">
                    <a-button
                      type="text"
                      @click="handleTogglePin(doc)"
                      class="flex items-center"
                      :class="{ 'text-blue-500': doc.pinned, 'text-gray-400': !doc.pinned }"
                    >
                      <pushpin-outlined :rotate="doc.pinned ? 0 : -45" />
                    </a-button>
                    <a-checkbox
                      v-model:checked="doc.selected"
                      @change="handleWorkspaceSelection(doc)"
                      class="text-blue-500"
                    />
                  </div>
                </template>
                <a-card-meta
                  :title="doc.filename"
                  class="text-sm font-semibold text-gray-800 dark:text-gray-200"
                >
                  <template #description>
                    <p class="text-gray-600 dark:text-gray-400">
                      创建时间: {{ parseDateTime(doc.createdAt) }}
                    </p>
                    <p class="text-gray-600 dark:text-gray-400">
                      更新时间: {{ parseDateTime(doc.lastUpdatedAt) }}
                    </p>
                  </template>
                </a-card-meta>
                <div class="mt-4 flex justify-end">
                  <a-button
                    type="link"
                    danger
                    @click="handleDeleteFromWorkspace(doc.filename)"
                    class="text-red-500 hover:text-red-600"
                  >
                    从工作台移除
                  </a-button>
                </div>
              </a-card>
            </div>
          </template>
          <template v-else>
            <div
              class="flex flex-col items-center justify-center py-12"
              v-motion-fade
            >
              <a-empty
                description="您还没有工作台，请先创建工作台"
                class="text-gray-500 dark:text-gray-400"
              >
                <a-button
                  type="primary"
                  @click="handleCreateWorkspace"
                  class="bg-gradient-to-r from-blue-500 to-indigo-500 hover:from-blue-600 hover:to-indigo-600"
                >
                  创建工作台
                </a-button>
              </a-empty>
            </div>
          </template>
        </a-tab-pane>

        <a-tab-pane key="archive" tab="云平台文档">
          <!-- 批量操作 -->
          <div class="flex flex-nowrap gap-2 mb-4">
            <a-upload
              :customRequest="handleUpload"
              :showUploadList="false"
              accept=".pdf,.doc,.docx,.txt"
            >
              <a-button
                type="primary"
                class="flex items-center space-x-2 bg-gradient-to-r from-green-500 to-emerald-500 hover:from-green-600 hover:to-emerald-600"
              >
                <upload-outlined></upload-outlined>
                <span>上传文档</span>
              </a-button>
            </a-upload>
            <div
              v-if="selectedArchiveItems.length > 0"
              class="mb-6"
              v-motion-slide-up
            >
              <a-button
                type="primary"
                @click="handleBatchAddToWorkspace"
                class="bg-gradient-to-r from-green-500 to-emerald-500 hover:from-green-600 hover:to-emerald-600"
              >
                添加到工作台 ({{ selectedArchiveItems.length }})
              </a-button>
            </div>
          </div>
          <!-- 文档网格 -->
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <a-card
              v-for="doc in archiveDocuments"
              :key="doc.id"
              class="document-card bg-white dark:bg-gray-800 rounded-xl shadow-md hover:shadow-xl transition-all duration-300"
              v-motion
              :initial="{ opacity: 0, y: 20 }"
              :enter="{ opacity: 1, y: 0, transition: { duration: 500 } }"
            >
              <template #extra>
                <a-checkbox
                  v-model:checked="doc.selected"
                  @change="handleArchiveSelection(doc)"
                  :disabled="isDocumentInWorkspace(doc.name)"
                  class="text-blue-500"
                />
              </template>
              <a-card-meta
                :title="doc.title"
                class="text-sm font-semibold text-gray-800 dark:text-gray-200"
              >
                <template #description>
                  <p class="text-gray-600 dark:text-gray-400">
                    作者: {{ doc.docAuthor || "未知作者" }}
                  </p>
                  <p class="text-gray-600 dark:text-gray-400">
                    来源: {{ doc.docSource || "未知来源" }}
                  </p>
                  <p class="text-gray-600 dark:text-gray-400">
                    发布时间: {{ doc.published }}
                  </p>
                  <p class="text-gray-600 dark:text-gray-400">
                    字数: {{ doc.wordCount }}
                  </p>
                  <p
                    v-if="isDocumentInWorkspace(doc.name)"
                    class="mt-2 text-green-500 flex items-center"
                  >
                    <check-circle-outlined class="mr-1" /> 已加入工作台
                  </p>
                </template>
              </a-card-meta>
              <div class="mt-4 flex justify-end">
                <a-button
                  type="link"
                  @click="handleAddToWorkspace(doc.name)"
                  :disabled="isDocumentInWorkspace(doc.name)"
                  class="text-blue-500 hover:text-blue-600 disabled:text-gray-400"
                >
                  添加到工作台
                </a-button>
              </div>
            </a-card>
          </div>
        </a-tab-pane>
      </a-tabs>
    </div>

    <spin-loader
      v-if="loading"
      class="fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-50"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { message } from "ant-design-vue";
import {
  UploadOutlined,
  PlusOutlined,
  CheckCircleOutlined,
  PushpinOutlined,
} from "@ant-design/icons-vue";
import { post, get, del, postJSON } from "@/utility/request";
import { parseDateTime } from "@/utility/common";
import SpinLoader from "@/components/spinLoader.vue";
import { useMotion } from "@vueuse/motion";
import { useUserStore } from "@/stores/user";

const userStore = useUserStore();
useMotion();

const activeTab = ref("workspace");
const archiveDocuments = ref([]);
const workspaceDocuments = ref([]);
const loading = ref(false);
const workspaceDetails = ref(null);
const selectedWorkspaceItems = ref([]);
const selectedArchiveItems = ref([]);

const formatFileSize = (bytes) => {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
};

const handleWorkspaceSelection = (doc) => {
  if (doc.selected) {
    selectedWorkspaceItems.value.push(doc.filename);
  } else {
    const index = selectedWorkspaceItems.value.indexOf(doc.filename);
    if (index > -1) {
      selectedWorkspaceItems.value.splice(index, 1);
    }
  }
};

const isDocumentInWorkspace = (docName) => {
  return workspaceDocuments.value.some((doc) => doc.filename === docName);
};

const handleArchiveSelection = (doc) => {
  if (isDocumentInWorkspace(doc.name)) {
    return;
  }

  if (doc.selected) {
    selectedArchiveItems.value.push(doc.name);
  } else {
    const index = selectedArchiveItems.value.indexOf(doc.name);
    if (index > -1) {
      selectedArchiveItems.value.splice(index, 1);
    }
  }
};

const handleBatchDeleteFromWorkspace = () => {
  if (selectedWorkspaceItems.value.length === 0) {
    message.warning("请选择要移除的文档");
    return;
  }
  const fileNamesWithPrefix = selectedWorkspaceItems.value.map(
    (filename) => userStore.user.workspaceId + "/" + filename
  );
  postJSON(
    "/api/rag/documents/delete",
    {
      fileNames: fileNamesWithPrefix,
    },
    (msg) => {
      message.success(msg);
      selectedWorkspaceItems.value = [];
      fetchWorkspaceDetails();
    },
    (msg) => {
      message.error(msg);
    }
  );
};

const handleBatchAddToWorkspace = () => {
  if (selectedArchiveItems.value.length === 0) {
    message.warning("请选择要添加的文档");
    return;
  }

  const fileNamesWithPrefix = selectedArchiveItems.value.map(
    (name) => userStore.user.workspaceId + "/" + name
  );

  postJSON(
    "/api/rag/documents/add",
    {
      fileNames: fileNamesWithPrefix,
    },
    (msg) => {
      message.success(msg);
      selectedArchiveItems.value = [];
      fetchWorkspaceDetails();
    },
    (msg) => {
      message.error(msg);
    }
  );
};

const handleAddToWorkspace = (docId) => {
  docId = userStore.user.workspaceId + "/" + docId;
  postJSON(
    "/api/rag/documents/add",
    {
      fileNames: [`${docId}`],
    },
    (msg) => {
      message.success(msg);
      fetchWorkspaceDetails();
    },
    (msg) => {
      message.error(msg);
    }
  );
};

const handleUpload = async ({ file }) => {
  const formData = new FormData();
  formData.append("file", file);

  post(
    "/api/rag/documents",
    formData,
    (msg) => {
      message.success(msg);
      fetchArchiveDocuments();
    },
    (msg) => {
      message.error(msg);
    }
  );
};

const handleDeleteFromWorkspace = (docId) => {
  docId = userStore.user.workspaceId + "/" + docId;
  postJSON(
    "/api/rag/documents/delete",
    {
      fileNames: [docId],
    },
    (msg) => {
      message.success(msg);
      fetchWorkspaceDetails();
    },
    (msg) => {
      message.error(msg);
    }
  );
};

const handleCreateWorkspace = () => {
  post(
    "/api/rag/workspace",
    {},
    (msg) => {
      message.success(msg);
      fetchWorkspaceDetails();
    },
    (msg) => {
      message.error(msg);
    }
  );
};

const fetchArchiveDocuments = () => {
  loading.value = true;
  get(
    "/api/rag/documents",
    null,
    (msg, data) => {
      archiveDocuments.value = (data.documents || []).map((doc) => ({
        ...doc,
        selected: false,
      }));
      loading.value = false;
    },
    (msg) => {
      message.error(msg);
      loading.value = false;
    }
  );
};

const fetchWorkspaceDetails = () => {
  get(
    "/api/rag/workspace",
    null,
    (msg, data) => {
      workspaceDetails.value = data;
      workspaceDocuments.value = (data.documents || []).map((doc) => ({
        ...doc,
        selected: false,
      }));
    },
    (msg) => {
      message.error(msg);
    }
  );
};

const handleTogglePin = async (doc) => {
  postJSON(`/api/rag/documents/status`, {
      filename: userStore.user.workspaceId + '/' + doc.filename,
      pin_status: !doc.pinned
    },(msg,data) => {
      doc.pinned = !doc.pinned;
      message.success(doc.pinned ? '文档已固定' : '文档已取消固定');
    },(msg,data) => {
      message.warning('操作失败：' + msg);
    },(msg,data) => {
      message.error('操作失败：' + msg);
    });
};

onMounted(() => {
  fetchArchiveDocuments();
  fetchWorkspaceDetails();
});
</script>

<style scoped>
.document-card {
  @apply transition-all duration-300 hover:transform hover:scale-105;
}

:deep(.ant-card-head) {
  @apply border-b border-gray-200 dark:border-gray-700;
}

:deep(.ant-card-body) {
  @apply p-6;
}

:deep(.ant-descriptions-item-label) {
  @apply bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-300;
}

:deep(.ant-descriptions-item-content) {
  @apply bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300;
}

:deep(.ant-tabs-nav) {
  @apply mb-6;
}

:deep(.ant-tabs-tab) {
  @apply text-gray-600 dark:text-gray-400;
}

:deep(.ant-tabs-tab-active) {
  @apply text-blue-500;
}

:deep(.ant-tabs-ink-bar) {
  @apply bg-gradient-to-r from-blue-500 to-indigo-500;
}

:deep(.ant-checkbox-checked .ant-checkbox-inner) {
  @apply bg-blue-500 border-blue-500;
}

:deep(.ant-btn-primary) {
  @apply border-0;
}

:deep(.ant-btn-link) {
  @apply border-0;
}

:deep(.ant-empty-image) {
  @apply filter dark:brightness-75;
}

:deep(.ant-card-meta-title) {
  @apply text-lg font-semibold text-gray-800 dark:text-gray-200;
}

:deep(.ant-card-meta-description) {
  @apply text-gray-600 dark:text-gray-400;
}

:deep(.ant-upload) {
  @apply w-auto;
}

:deep(.ant-upload-list) {
  @apply hidden;
}

/* 暗色模式适配 */
@media (prefers-color-scheme: dark) {
  :deep(.ant-card) {
    @apply bg-gray-800 border-gray-700;
  }

  :deep(.ant-descriptions-item-label) {
    @apply bg-gray-700 text-gray-300;
  }

  :deep(.ant-descriptions-item-content) {
    @apply bg-gray-800 text-gray-300;
  }

  :deep(.ant-tabs-tab) {
    @apply text-gray-400;
  }

  :deep(.ant-tabs-tab-active) {
    @apply text-blue-400;
  }

  :deep(.ant-tabs-ink-bar) {
    @apply bg-gradient-to-r from-blue-400 to-indigo-400;
  }
}
</style>
