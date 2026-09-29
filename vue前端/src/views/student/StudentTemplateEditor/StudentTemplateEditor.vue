<script setup>
import { ref, reactive, onMounted, watch, onUnmounted, computed } from "vue";
import { Document, Packer, ImageRun } from "docx";
import { useMotion } from "@vueuse/motion";
import { message, notification } from "ant-design-vue";
import {
  PlusOutlined,
  DownloadOutlined,
  LoadingOutlined,
  EditOutlined,
  DeleteOutlined,
  SaveOutlined,
  CommentOutlined,
  HistoryOutlined,
  EyeOutlined,
  RollbackOutlined,
} from "@ant-design/icons-vue";
import ComponentPalette from "./components/ComponentPalette.vue";
import TemplateCanvas from "./components/TemplateCanvas.vue";
import PropertyPanel from "./components/PropertyPanel.vue";
import NewTemplateDialog from "./components/NewTemplateDialog.vue";
import ExportDialog from "./components/ExportDialog.vue";
import TextEnhancePanel from "./components/TextEnhancePanel.vue";
const route = useRoute();
const templateId = route.params.id;
// 导入工具函数
import {
  defaultTemplate,
  createComponent,
  uploadImage,
} from "./utils/templateUtils";
import { generateDoc } from "./utils/documentExporter";
// 导入请求函数
import { get, post, postJSON } from "@/utility/request";
import { useRoute } from "vue-router";
import { parseDateTime } from "@/utility/common";
import router from "@/router";

// 当前模板状态
const currentTemplate = ref(JSON.parse(JSON.stringify(defaultTemplate)));
const savedTemplates = ref([]);
const currentTemplateIndex = ref(-1);

// 编辑器状态
const selectedComponentId = ref(null);
const isDragging = ref(false);
const showNewTemplateDialog = ref(false);
const isExporting = ref(false); // 添加导出状态
const exportProgress = ref(0); // 添加导出进度
const exportStatus = ref(""); // 添加导出状态信息
const hasUnsavedChanges = ref(false); // 是否有未保存的更改
const lastSavedTemplate = ref(null); // 最后一次保存的模板内容

// 创建DOM引用
const leftPanelRef = ref(null);
const canvasAreaRef = ref(null);
const rightPanelRef = ref(null);

const historyIDs = ref(null);
const fetchHistoryIds = () => {
  get(
    "/api/template/history/" + route.params.id + "/list",
    {},
    (message, data) => {
      historyIDs.value = data;
    },
    (message, data) => {},
    (message, data) => {}
  );
};
// 动画相关 - 确保在设置动画前已经创建了ref
const leftPanel = useMotion(leftPanelRef, {
  initial: { opacity: 0, x: -50 },
  enter: { opacity: 1, x: 0, transition: { duration: 500 } },
});

const canvasArea = useMotion(canvasAreaRef, {
  initial: { opacity: 0, scale: 0.9 },
  enter: { opacity: 1, scale: 1, transition: { duration: 600, delay: 200 } },
});

const rightPanel = useMotion(rightPanelRef, {
  initial: { opacity: 0, x: 50 },
  enter: { opacity: 1, x: 0, transition: { duration: 500, delay: 100 } },
});

// 文本选择相关状态
const textSelectionPanel = reactive({
  visible: false,
  x: 0,
  y: 0,
  selectedText: "",
  selectedComponent: null,
});

const activeToolTab = ref('annotations')

// 从后端加载已保存的模板
onMounted(() => {
  loadTemplates();
  fetchHistoryIds();
  // 设置自动保存定时器 (3分钟)
  const autoSaveInterval = setInterval(() => {
    if (hasUnsavedChanges.value && templateId) {
      saveTemplate();
    }
  }, 3 * 60 * 1000);

  // 组件卸载时清除定时器
  onUnmounted(() => {
    clearInterval(autoSaveInterval);
  });
  document.addEventListener("mouseup", handleTextSelection);
  document.addEventListener("click", handleGlobalClick);
});

// 监听模板变化
watch(
  currentTemplate,
  (newVal) => {
    // 检查是否有实际变化
    if (lastSavedTemplate.value) {
      // 比较当前模板和最后保存的模板是否不同
      const currentStr = JSON.stringify(newVal);
      const savedStr = JSON.stringify(lastSavedTemplate.value);
      hasUnsavedChanges.value = currentStr !== savedStr;
    } else {
      hasUnsavedChanges.value = true;
    }
  },
  { deep: true }
);

// 手动保存模板
const saveTemplate = () => {
  if (!hasUnsavedChanges.value) {
    message.info("没有任何修改需要保存");
    return;
  }

  if (templateId) {
    updateTemplateToServer(templateId, currentTemplate.value);
    // 更新最后保存的模板状态
    lastSavedTemplate.value = JSON.parse(JSON.stringify(currentTemplate.value));
    hasUnsavedChanges.value = false;
  }
};

// 处理组件重新排序
const handleComponentReorder = (componentId, fromIndex, toIndex) => {
  // 获取要移动的组件
  const component = currentTemplate.value.components[fromIndex];

  // 从原位置删除
  currentTemplate.value.components.splice(fromIndex, 1);

  // 插入到新位置
  currentTemplate.value.components.splice(toIndex, 0, component);

  // 标记为有未保存的更改
  hasUnsavedChanges.value = true;
};

// 处理模板内容变更通知
const handleTemplateChanged = () => {
  hasUnsavedChanges.value = true;
};

// 从后端加载已保存的模板
const loadTemplates = () => {
  loadTemplate();
};

// 更新模板到服务器
const updateTemplateToServer = (templateId, templateData) => {
  // 准备发送到后端的数据
  const updateData = {
    title: templateData.name,
    type: templateData.type,
    description: templateData.description || "",
    content: JSON.stringify(templateData),
  };

  // 发送到后端更新
  postJSON(
    `/api/template/update/${templateId}`,
    updateData,
    (message, data) => {
      notification.success({
        message: "保存成功",
        description: "模板已成功更新",
        placement: "bottomRight",
        duration: 2,
      });
      fetchHistoryIds();
    },
    (errorMsg) => {
      notification.error({
        message: "保存失败",
        description: errorMsg || "更新模板时发生错误",
        placement: "bottomRight",
      });
    }
  );
};

// 加载模板
const loadTemplate = () => {
  get(
    `/api/template/detail/${templateId}`,
    null,
    (message, data) => {
      // 解析Content字段中存储的模板内容
      const templateContent = JSON.parse(data.content);
      // 合并模板信息
      const fullTemplate = {
        id: data.id,
        name: data.title,
        type: data.type,
        description: data.description,
        ...templateContent,
      };

      // 更新当前模板
      currentTemplate.value = fullTemplate;
      // 保存最后加载的模板状态，用于比较变更
      lastSavedTemplate.value = JSON.parse(JSON.stringify(fullTemplate));
      hasUnsavedChanges.value = false;
    },
    (msg, data) => {
      console.error("获取模板详情失败:", msg);
      message.error("获取模板详情失败");
    },
    (msg, data) => {
      console.error("获取模板详情失败:", msg);
      message.error("获取模板详情失败");
    }
  );
};

// 删除模板
const deleteTemplate = (index) => {
  const templateId = savedTemplates.value[index].id;

  // 发送删除请求到后端
  post(
    `/api/template/delete/${templateId}`,
    {},
    (message, data) => {
      // 从本地列表中移除
      savedTemplates.value.splice(index, 1);

      if (currentTemplateIndex.value === index) {
        if (savedTemplates.value.length > 0) {
          loadTemplate(0);
        } else {
          currentTemplateIndex.value = -1;
          currentTemplate.value = JSON.parse(JSON.stringify(defaultTemplate));
        }
      } else if (currentTemplateIndex.value > index) {
        currentTemplateIndex.value--;
      }

      notification.success({
        message: "删除成功",
        description: "模板已成功删除",
        placement: "bottomRight",
        duration: 2,
      });
    },
    (errorMsg) => {
      notification.error({
        message: "删除失败",
        description: errorMsg || "删除模板时发生错误",
        placement: "bottomRight",
      });
    }
  );
};

// 处理组件拖放
const handleComponentDrop = (componentType, position) => {
  const newComponent = createComponent(componentType, position);
  currentTemplate.value.components.push(newComponent);
  message.success("组件添加成功");

  // 标记为有未保存的更改
  hasUnsavedChanges.value = true;
};

// 更新组件属性
const updateComponentProperty = (componentId, property, value) => {
  const component = currentTemplate.value.components.find(
    (c) => c.id === componentId
  );
  if (component) {
    component[property] = value;
    // 标记为有未保存的更改
    hasUnsavedChanges.value = true;
  }
};

// 删除组件
const deleteComponent = (componentId) => {
  const index = currentTemplate.value.components.findIndex(
    (c) => c.id === componentId
  );
  if (index !== -1) {
    currentTemplate.value.components.splice(index, 1);
    if (selectedComponentId.value === componentId) {
      selectedComponentId.value = null;
    }
    message.success("组件已删除");

    // 标记为有未保存的更改
    hasUnsavedChanges.value = true;
  }
};

// 导出文档
const handleExportDoc = async () => {
  try {
    isExporting.value = true;

    await generateDoc(
      currentTemplate.value,
      (progress) => {
        exportProgress.value = progress;
      },
      (status) => {
        exportStatus.value = status;
      }
    );

    // 延迟关闭导出状态
    setTimeout(() => {
      isExporting.value = false;
      notification.success({
        message: "导出成功",
        description: "文档已成功生成",
        placement: "bottomRight",
        duration: 3,
      });
    }, 1500);
  } catch (error) {
    console.error("文档生成错误:", error);
    exportStatus.value = `导出失败: ${error.message}`;
    setTimeout(() => {
      isExporting.value = false;
      notification.error({
        message: "导出失败",
        description: error.message || "生成文档时发生错误",
        placement: "bottomRight",
      });
    }, 3000);
  }
};

// 关闭导出对话框
const handleCloseExportDialog = () => {
  isExporting.value = false;
};

// 获取当前选中组件的标注
const getSelectedComponentAnnotations = computed(() => {
  if (!selectedComponentId.value) return [];
  const selectedComponent = currentTemplate.value.components.find(
    (c) => c.id === selectedComponentId.value
  );
  return selectedComponent?.annotations || [];
});

// 处理画布点击事件
const handleCanvasClick = (event) => {
  // 如果点击的是画布本身而不是组件，则取消选中
  if (event.target === canvasAreaRef.value) {
    selectedComponentId.value = null;
  }
};

// 处理文本选择事件
const handleTextSelection = (event) => {
  const selection = window.getSelection();
  const selectedText = selection.toString().trim();

  if (selectedText) {
    // 获取选区位置
    const range = selection.getRangeAt(0);
    const rect = range.getBoundingClientRect();

    // 更新面板状态
    textSelectionPanel.visible = true;
    textSelectionPanel.x = rect.left;
    textSelectionPanel.y = rect.bottom + 10;
    textSelectionPanel.selectedText = selectedText;

    // 查找选中文本所在的组件
    const componentElement = event.target.closest("[data-component-id]");
    if (componentElement) {
      textSelectionPanel.selectedComponent =
        currentTemplate.value.components.find(
          (c) => c.id === componentElement.dataset.componentId
        );
    }
  } else {
    textSelectionPanel.visible = false;
  }
};

// 处理文本更新
const handleTextUpdate = (enhancedText) => {
  if (textSelectionPanel.selectedComponent) {
    // 更新组件中的文本内容
    const component = textSelectionPanel.selectedComponent;
    if (component.type === "text" || component.type === "paragraph") {
      component.content = component.content.replace(
        textSelectionPanel.selectedText,
        enhancedText
      );
      hasUnsavedChanges.value = true;
    }
  }
};

// 关闭文本选择面板
const closeTextSelectionPanel = () => {
  textSelectionPanel.visible = false;
  textSelectionPanel.selectedText = "";
  textSelectionPanel.selectedComponent = null;
};

// 监听点击事件以关闭面板
const handleGlobalClick = (event) => {
  if (!event.target.closest(".text-enhance-panel")) {
    closeTextSelectionPanel();
  }
};

const goHistory = (id) => {
  window.open( window.location.origin + '/student/template/history/' + id)
}

// 加载模板
const loadTemplateHistory = (id) => {
  get(`/api/template/history/${id}/content`, null,
    (message, data) => {
      // 解析Content字段中存储的模板内容
      const templateContent = JSON.parse(data)
      // 合并模板信息
      const fullTemplate = {
        id: data.id,
        name: data.title,
        type: data.type,
        description: data.description,
        ...templateContent
      }
      
      // 更新当前模板
      currentTemplate.value = fullTemplate
      // 保存最后加载的模板状态，用于比较变更
      lastSavedTemplate.value = JSON.parse(JSON.stringify(fullTemplate))
      hasUnsavedChanges.value = false
    },
    (msg, data) => {
      console.error('获取模板详情失败:', msg)
      message.error('获取模板详情失败')
    },
    (msg, data) => {
      console.error('获取模板详情失败:', msg)
      message.error('获取模板详情失败')
    }
  )
}
</script>

<template>
  <div class="flex h-screen bg-gray-50 overflow-hidden">
    <!-- 左侧组件面板 -->
    <div
      ref="leftPanelRef"
      class="w-64 bg-white shadow-xl p-5 overflow-y-auto flex flex-col border-r border-gray-200 transition-all duration-300 ease-in-out"
      v-show="!selectedComponentId"
    >
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-800 flex items-center">
          <EditOutlined class="mr-2 text-blue-500" />
          模板编辑器
        </h1>
        <p class="text-sm text-gray-500 mt-1">拖拽组件到右侧画布构建模板</p>
      </div>

      <div class="mt-4 pt-4 border-t border-gray-200">
        <div class="flex flex-col gap-3">
          <a-button
            type="primary"
            class="flex items-center justify-center h-10 text-base"
            @click="saveTemplate"
            :disabled="!hasUnsavedChanges"
          >
            <template #icon><SaveOutlined /></template>
            保存模板
            <a-badge v-if="hasUnsavedChanges" status="processing" color="red" />
          </a-button>

          <a-button
            type="primary"
            class="flex items-center justify-center h-10 text-base"
            :loading="isExporting"
            @click="handleExportDoc"
            danger
          >
            <template #icon><DownloadOutlined v-if="!isExporting" /></template>
            {{ isExporting ? "导出中..." : "导出文档" }}
          </a-button>
        </div>
      </div>

      <div class="flex-1 mt-6">
        <h2 class="text-lg font-semibold text-gray-700 mb-3">组件库</h2>
        <ComponentPalette
          @component-drag-start="isDragging = true"
          @component-drag-end="isDragging = false"
        />
      </div>
    </div>

    <!-- 中间画布区域 -->
    <div
      ref="canvasAreaRef"
      class="flex-1 p-6 overflow-auto bg-gradient-to-br from-gray-50 to-gray-100 transition-all duration-300 ease-in-out relative"
      :class="{ 'pr-6': !selectedComponentId, 'pr-0': selectedComponentId }"
      @click="handleCanvasClick"
    >
      <div class="mb-4 flex justify-between items-center">
        <h2 class="text-xl font-semibold text-gray-700">
          {{ currentTemplate.name || "未命名模板" }}
        </h2>
        <a-tag color="blue">{{ currentTemplate.type || "标准模板" }}</a-tag>
      </div>

      <div class="flex gap-6">
        <!-- 主渲染区域 -->
        <div class="flex-1 p-6 min-h-[calc(100vh-120px)]">
          <TemplateCanvas
            v-model:selectedComponentId="selectedComponentId"
            :template="currentTemplate"
            :isDragging="isDragging"
            :hasUnsavedChanges="hasUnsavedChanges"
            @component-drop="handleComponentDrop"
            @component-update="updateComponentProperty"
            @component-delete="deleteComponent"
            @component-reorder="handleComponentReorder"
            @template-changed="handleTemplateChanged"
          />
        </div>

        <!-- 标注显示区域 -->
        <div
          class="w-64 bg-white rounded-lg shadow-md p-4 max-h-screen overflow-auto sticky top-6 m-2"
        >
          <a-tabs v-model:activeKey="activeToolTab">
            <a-tab-pane key="annotations" tab="教师标注">
              <div class="text-lg font-semibold text-gray-700 mb-4 flex items-center">
                <CommentOutlined class="mr-2 text-blue-500" />
                教师标注
              </div>

              <div class="space-y-6">
                <template v-if="currentTemplate.components.length">
                  <div
                    v-for="component in currentTemplate.components"
                    :key="component.id"
                    class="component-annotations"
                    :class="{
                      'selected-component': component.id === selectedComponentId,
                    }"
                  >
                    <div class="component-info mb-2 flex items-center justify-between">
                      <span class="text-sm font-medium text-gray-600">
                        {{
                          component.type === "text"
                            ? "文本组件"
                            : component.type === "title"
                            ? "标题组件"
                            : component.type === "paragraph"
                            ? "段落组件"
                            : component.type === "table"
                            ? "表格组件"
                            : component.type === "image"
                            ? "图片组件"
                            : component.type === "list"
                            ? "列表组件"
                            : component.type === "cover-table"
                            ? "封面表格"
                            : component.type === "reference"
                            ? "参考文献组件"
                            : "未知组件"
                        }}
                      </span>
                      <span
                        v-if="component.annotations && component.annotations.length"
                        class="text-xs bg-blue-100 text-blue-600 px-2 py-1 rounded-full"
                      >
                        {{ component.annotations.length }}条标注
                      </span>
                    </div>

                    <template
                      v-if="component.annotations && component.annotations.length"
                    >
                      <div
                        v-for="(annotation, index) in component.annotations"
                        :key="index"
                        class="annotation-item bg-blue-50 p-3 rounded-md mb-2"
                      >
                        <div class="text-xs text-gray-500 mb-1">
                          {{
                            annotation.date
                              ? new Date(annotation.date).toLocaleString()
                              : "未知时间"
                          }}
                        </div>
                        <div class="text-sm text-gray-700">
                          {{ annotation.content }}
                        </div>
                      </div>
                    </template>
                    <div v-else class="text-sm text-gray-500 text-center py-2">
                      暂无标注
                    </div>

                    <a-divider class="my-3" />
                  </div>
                </template>
                <div v-else class="text-gray-500 text-center py-4">
                  模板中暂无组件
                </div>
              </div>
            </a-tab-pane>

            <a-tab-pane key="history" tab="修改历史">
              <div class="text-lg font-semibold text-gray-700 mb-4 flex items-center">
                <HistoryOutlined class="mr-2 text-blue-500" />
                修改历史
              </div>

              <div class="space-y-3">
                <div
                  v-for="(history, index) in historyIDs"
                  :key="history.id"
                  class="component-item p-3 bg-white rounded-lg shadow hover:shadow-md cursor-pointer transition-all"
                  :class="{ 'transform hover:-translate-y-1': true }"
                >
                  <div class="flex items-center gap-3">
                    <a-popover title="操作" placement="left">
                      <template #content>
                        <div class="space-y-2">
                          <a-button type="link" block @click="goHistory(history.id)">
                            <EyeOutlined />
                            预览
                          </a-button>
                          <a-button type="link" block @click="loadTemplateHistory(history.id)">
                            <RollbackOutlined />
                            还原
                          </a-button>
                        </div>
                      </template>
                      <div class="flex-1 min-w-0">
                        <h3 class="font-medium text-gray-800 truncate">
                          第{{ historyIDs.length - index }}次修改
                        </h3>
                        <p class="text-xs text-gray-500 truncate">
                          修改于{{ parseDateTime(history.CreatedAt) }}
                        </p>
                      </div>
                    </a-popover>
                  </div>
                </div>
              </div>
            </a-tab-pane>
          </a-tabs>
        </div>
      </div>
    </div>

    <!-- 右侧属性面板 -->
    <div
      ref="rightPanelRef"
      class="w-96 bg-white shadow-xl p-5 overflow-y-auto border-l border-gray-200 transition-all duration-300 ease-in-out"
      v-if="selectedComponentId"
    >
      <h2 class="text-lg font-semibold text-gray-700 mb-4 flex items-center">
        <SaveOutlined class="mr-2 text-blue-500" />
        属性设置
      </h2>

      <a-divider class="my-3" />

      <PropertyPanel
        :component="
          currentTemplate.components.find((c) => c.id === selectedComponentId)
        "
        @update="updateComponentProperty"
        @delete="deleteComponent"
        @template-changed="handleTemplateChanged"
      />
    </div>

    <!-- 导出进度对话框 -->
    <ExportDialog
      :isExporting="isExporting"
      :exportProgress="exportProgress"
      :exportStatus="exportStatus"
      @close="handleCloseExportDialog"
    />

    <!-- 文本选择增强面板 -->
    <TextEnhancePanel
      v-model:visible="textSelectionPanel.visible"
      :x="textSelectionPanel.x"
      :y="textSelectionPanel.y"
      :selected-text="textSelectionPanel.selectedText"
      @update="handleTextUpdate"
      @close="closeTextSelectionPanel"
      class="text-enhance-panel"
    />
  </div>
</template>

<style scoped>
.component-dragging {
  opacity: 0.5;
  cursor: move;
}

:deep(.ant-btn-primary) {
  background-color: #1890ff;
}

:deep(.ant-btn-primary.ant-btn-dangerous) {
  background-color: #ff4d4f;
}

:deep(.ant-empty-image) {
  height: 80px;
}

:deep(.ant-divider) {
  margin: 12px 0;
}

.component-annotations {
  padding: 8px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.selected-component {
  background-color: rgba(24, 144, 255, 0.05);
  border: 1px solid rgba(24, 144, 255, 0.3);
}

.annotation-item {
  position: relative;
  border-left: 3px solid #1890ff;
}

.annotation-item::before {
  content: "";
  position: absolute;
  left: -6px;
  top: 50%;
  transform: translateY(-50%);
  width: 8px;
  height: 8px;
  background-color: #1890ff;
  border-radius: 50%;
}

:deep(.ant-tabs-nav) {
  margin-bottom: 16px;
}

:deep(.ant-tabs-tab) {
  padding: 8px 16px;
}

:deep(.ant-tabs-tab-active) {
  font-weight: 500;
}

:deep(.ant-popover-buttons) {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
