<script setup>
import { Modal, Input } from 'ant-design-vue'
import { useGroupChatList } from './groupChatListLogic'

const props = defineProps({
  activeGroupId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['select-group', 'create-group'])

// 导入并使用封装的群组列表逻辑
const {
  // 状态
  userStore,
  groups,
  loading,
  searchText,
  createModalVisible,
  createGroupForm,
  filteredGroups,

  // 方法
  fetchGroups,
  selectGroup,
  openCreateModal,
  handleCreateGroup
} = useGroupChatList(props, emit)
</script>

<template>
  <div class="group-chat-list h-full flex flex-col bg-white dark:bg-zinc-900">
    <!-- 顶部搜索和创建 -->
    <div class="p-4 border-b border-gray-200 dark:border-gray-700">
      <div class="flex items-center gap-2 mb-4">
        <h2 class="text-lg font-bold text-gray-800 dark:text-gray-200">我的群组</h2>
        <button 
          @click="openCreateModal"
          class="ml-auto p-1.5 text-gray-500 hover:text-blue-500 dark:text-gray-400 dark:hover:text-blue-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transform hover:scale-105"
          title="创建群组"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
        </button>
      </div>
      
      <div class="relative">
        <Input 
          v-model:value="searchText" 
          placeholder="搜索群组" 
          allowClear
          class="hover:border-blue-500 focus:border-blue-500  dark:border-gray-700 dark:text-gray-300 dark:placeholder-gray-500"
        >
          <template #prefix>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-400 dark:text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </template>
        </Input>
      </div>
    </div>
    
    <!-- 群组列表 -->
    <div class="flex-1 overflow-y-auto py-2 scrollbar-thin scrollbar-thumb-gray-300 dark:scrollbar-thumb-gray-600 scrollbar-track-transparent">
      <!-- 加载中 -->
      <div v-if="loading" class="py-8 text-center">
        <div class="animate-spin h-8 w-8 border-2 border-blue-500 dark:border-blue-400 rounded-full border-t-transparent mx-auto mb-3"></div>
        <p class="text-gray-500 dark:text-gray-400">加载中...</p>
      </div>
      
      <!-- 群组列表 -->
      <div 
        v-else-if="filteredGroups.length > 0"
        class="space-y-1 px-2"
      >
        <div 
          v-for="group in filteredGroups" 
          :key="group.ID"
          class="flex items-center p-3 rounded-xl cursor-pointer transition-all duration-300 group"
          :class="[
            group.ID === activeGroupId 
              ? 'bg-gradient-to-r from-blue-500/10 to-cyan-500/10 dark:from-blue-500/20 dark:to-cyan-500/20 text-blue-600 dark:text-blue-400' 
              : 'hover:bg-gray-100 dark:hover:bg-zinc-800/70'
          ]"
          @click="selectGroup(group)"
        >
          <!-- 群组头像 -->
          <div class="w-10 h-10 rounded-full overflow-hidden bg-gradient-to-br from-blue-100 to-indigo-100 dark:from-blue-900/30 dark:to-indigo-900/30 flex-shrink-0 border border-gray-200 dark:border-gray-700 shadow-sm group-hover:shadow-md group-hover:scale-105 transition-all duration-300">
            <img 
              v-if="group.AvatarURL" 
              :src="group.AvatarURL" 
              class="w-full h-full object-cover"
              alt="群组头像"
            >
            <div v-else class="w-full h-full flex items-center justify-center text-blue-500 dark:text-blue-400 text-sm font-bold">
              {{ group.Name.charAt(0) }}
            </div>
          </div>
          
          <!-- 群组信息 -->
          <div class="ml-3 overflow-hidden">
            <h3 class="font-medium text-gray-800 dark:text-gray-200 truncate group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
              {{ group.Name }}
            </h3>
            <p class="text-xs text-gray-500 dark:text-gray-400 truncate flex items-center gap-1.5">
              <span class="flex items-center gap-1">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
                {{ group.MemberCount }}人
              </span>
              <span v-if="group.CreatorID === userStore.user?.id" class="flex items-center gap-1 text-red-500">
                <span class="w-1 h-1 rounded-full bg-current"></span>
                创建者
              </span>
            </p>
          </div>
        </div>
      </div>
      
      <!-- 空列表提示 -->
      <div v-else-if="!loading" class="py-12 text-center">
        <div class="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/30 dark:to-indigo-900/30 rounded-2xl p-6 mx-4 border border-gray-200 dark:border-gray-700">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-300 dark:text-gray-600 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          <p class="text-gray-500 dark:text-gray-400 mb-4">{{ searchText ? '没有找到匹配的群组' : '暂无群组' }}</p>
          <button 
            @click="openCreateModal"
            class="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg hover:opacity-90 hover:shadow-md hover:shadow-blue-500/20 transition-all flex items-center gap-2 mx-auto transform hover:scale-105"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
            创建新群组
          </button>
        </div>
      </div>
    </div>
    
    <!-- 创建群组模态框 -->
    <Modal
      v-model:visible="createModalVisible"
      title="创建新群组"
      @ok="handleCreateGroup"
      okText="创建"
      cancelText="取消"
      class="custom-modal"
      :bodyStyle="{
        backgroundColor: 'var(--bg-color)',
        borderRadius: '0.75rem',
      }"
    >
      <div class="space-y-4 p-2">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组名称<span class="text-red-500">*</span></label>
          <Input 
            v-model:value="createGroupForm.name" 
            placeholder="输入群组名称" 
            class="hover:border-blue-500 focus:border-blue-500 dark:bg-zinc-800/50 dark:border-gray-700 dark:text-gray-300"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组描述</label>
          <Input.TextArea 
            v-model:value="createGroupForm.description" 
            placeholder="输入群组描述" 
            :rows="3" 
            class="hover:border-blue-500 focus:border-blue-500 dark:bg-zinc-800/50 dark:border-gray-700 dark:text-gray-300"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">群组头像链接</label>
          <Input 
            v-model:value="createGroupForm.avatarUrl" 
            placeholder="输入头像URL" 
            class="hover:border-blue-500 focus:border-blue-500 dark:bg-zinc-800/50 dark:border-gray-700 dark:text-gray-300"
          />
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
/* 自定义滚动条样式 */
.scrollbar-thin::-webkit-scrollbar {
  width: 6px;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background: transparent;
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 3px;
}

.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background-color: rgba(156, 163, 175, 0.7);
}

/* 自定义模态框样式 */
:deep(.custom-modal .ant-modal-content) {
  @apply rounded-2xl overflow-hidden border border-gray-200 dark:border-gray-700;
}

:deep(.custom-modal .ant-modal-header) {
  @apply bg-white dark:bg-zinc-900 border-b border-gray-200 dark:border-gray-700;
}

:deep(.custom-modal .ant-modal-footer) {
  @apply bg-white dark:bg-zinc-900 border-t border-gray-200 dark:border-gray-700;
}

:deep(.custom-modal .ant-btn-primary) {
  @apply bg-gradient-to-r from-blue-500 to-cyan-500 border-0;
}

:deep(.custom-modal .ant-btn-primary:hover) {
  @apply opacity-90 shadow-md shadow-blue-500/20;
}

/* 输入框激活状态 */
:deep(.ant-input-affix-wrapper:focus),
:deep(.ant-input-affix-wrapper-focused) {
  @apply border-blue-500 shadow-none;
}

:deep(.ant-input-affix-wrapper:hover) {
  @apply border-blue-500;
}
</style> 