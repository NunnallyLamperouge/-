<script setup>
import { ref, computed } from 'vue'
import { message, Dropdown, Menu, Tooltip, Input, Modal } from 'ant-design-vue'

const props = defineProps({
  members: {
    type: Array,
    default: () => []
  },
  isAdmin: {
    type: Boolean,
    default: false
  },
  isCreator: {
    type: Boolean,
    default: false
  },
  creatorId: {
    type: Number,
    default: 0
  },
  currentUserId: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['mute-member', 'remove-member', 'toggle-admin', 'update-nickname'])

// 搜索框
const searchText = ref('')
// 排序方式
const sortBy = ref('role') // 'role' | 'name' | 'joinTime'
// 编辑昵称相关
const editingMember = ref(null)
const newNickname = ref('')
const nicknameModalVisible = ref(false)

// 根据搜索和排序过滤成员
const filteredMembers = computed(() => {
  let result = [...props.members]
  
  // 搜索过滤
  if (searchText.value) {
    const query = searchText.value.toLowerCase()
    result = result.filter(member => 
      (member.Nickname && member.Nickname.toLowerCase().includes(query)) || 
      (member.Username && member.Username.toLowerCase().includes(query))
    )
  }
  
  // 排序
  switch (sortBy.value) {
    case 'role':
      // 创建者 > 管理员 > 普通成员
      result.sort((a, b) => {
        if (a.UserID === props.creatorId) return -1
        if (b.UserID === props.creatorId) return 1
        if (a.IsAdmin && !b.IsAdmin) return -1
        if (!a.IsAdmin && b.IsAdmin) return 1
        return 0
      })
      break
    case 'name':
      // 按昵称/用户名排序
      result.sort((a, b) => {
        const nameA = (a.Nickname || a.Username || '').toLowerCase()
        const nameB = (b.Nickname || b.Username || '').toLowerCase()
        return nameA.localeCompare(nameB)
      })
      break
    case 'joinTime':
      // 按加入时间排序
      result.sort((a, b) => new Date(b.JoinTime) - new Date(a.JoinTime))
      break
  }
  
  return result
})

// 获取成员状态文本
const getMemberStatus = (member) => {
  if (member.UserID === props.creatorId) return '创建者'
  if (member.IsAdmin) return '管理员'
  if (member.Status === 2) return '已禁言'
  return '成员'
}

// 获取成员状态样式
const getMemberStatusClass = (member) => {
  if (member.UserID === props.creatorId) return 'text-red-500'
  if (member.IsAdmin) return 'text-blue-500'
  if (member.Status === 2) return 'text-orange-500'
  return 'text-gray-500'
}

// 打开昵称编辑模态框
const openNicknameModal = (member) => {
  editingMember.value = member
  newNickname.value = member.Nickname || ''
  nicknameModalVisible.value = true
}

// 提交昵称修改
const submitNicknameChange = () => {
  if (!editingMember.value) return
  
  emit('update-nickname', editingMember.value, newNickname.value)
  nicknameModalVisible.value = false
}

// 禁言/解禁成员
const handleMute = (member) => {
  const isMuted = member.Status !== 2
  emit('mute-member', member, isMuted)
}

// 移除成员
const handleRemove = (member) => {
  emit('remove-member', member)
}

// 设置/取消管理员
const handleToggleAdmin = (member) => {
  emit('toggle-admin', member, !member.IsAdmin)
}

// 格式化时间
const formatJoinTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<template>
  <div class="member-list">
    <!-- 搜索和排序 -->
    <div class="p-4 border-b theme-border">
      <Input 
        v-model:value="searchText" 
        placeholder="搜索成员" 
        class="mb-2"
        allowClear
      >
        <template #prefix>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </template>
      </Input>
      
      <div class="flex justify-between text-xs text-gray-500 dark:text-gray-400">
        <span>共 {{ members.length }} 名成员</span>
        <div class="flex items-center gap-2">
          <span>排序：</span>
          <button 
            @click="sortBy = 'role'"
            class="px-1.5 py-0.5 rounded transition-colors"
            :class="sortBy === 'role' ? 'bg-blue-100 text-blue-600 dark:bg-blue-900/50 dark:text-blue-400' : 'hover:bg-gray-100 dark:hover:bg-gray-800'"
          >角色</button>
          <button 
            @click="sortBy = 'name'"
            class="px-1.5 py-0.5 rounded transition-colors"
            :class="sortBy === 'name' ? 'bg-blue-100 text-blue-600 dark:bg-blue-900/50 dark:text-blue-400' : 'hover:bg-gray-100 dark:hover:bg-gray-800'"
          >昵称</button>
          <button 
            @click="sortBy = 'joinTime'"
            class="px-1.5 py-0.5 rounded transition-colors"
            :class="sortBy === 'joinTime' ? 'bg-blue-100 text-blue-600 dark:bg-blue-900/50 dark:text-blue-400' : 'hover:bg-gray-100 dark:hover:bg-gray-800'"
          >加入时间</button>
        </div>
      </div>
    </div>
    
    <!-- 成员列表 -->
    <div class="py-2">
      <div 
        v-for="member in filteredMembers" 
        :key="member.UserID"
        class="flex items-center justify-between p-3 hover:bg-gray-50 dark:hover:bg-zinc-800 group transition-colors"
      >
        <div class="flex items-center gap-3">
          <!-- 成员头像 -->
          <div class="w-10 h-10 rounded-full overflow-hidden bg-gradient-to-br from-blue-100 to-indigo-100 dark:from-blue-900/30 dark:to-indigo-900/30 flex-shrink-0 border theme-border">
            <img 
              v-if="member.Avatar" 
              :src="member.Avatar" 
              class="w-full h-full object-cover"
              alt="用户头像"
            >
            <div v-else class="w-full h-full flex items-center justify-center text-blue-500 dark:text-blue-400 text-sm font-bold">
              {{ (member.Nickname || member.Username || 'U')?.charAt(0) }}
            </div>
          </div>
          
          <!-- 成员信息 -->
          <div>
            <div class="flex items-center gap-2">
              <span class="font-medium text-gray-800 dark:text-gray-200">
                {{ member.Nickname || member.Username }}
              </span>
              <Tooltip v-if="member.UserID === currentUserId" title="这是你">
                <span class="bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400 text-xs px-1.5 py-0.5 rounded">我</span>
              </Tooltip>
              <span 
                class="text-xs px-1.5 py-0.5 rounded-full"
                :class="getMemberStatusClass(member)"
              >
                {{ getMemberStatus(member) }}
              </span>
            </div>
            <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
              加入时间：{{ formatJoinTime(member.JoinTime) }}
            </div>
          </div>
        </div>
        
        <!-- 操作菜单 - 仅管理员可操作其他成员 -->
        <div v-if="(isAdmin || isCreator || member.UserID === currentUserId) && member.UserID !== creatorId" class="opacity-0 group-hover:opacity-100 transition-opacity">
          <Dropdown>
            <button class="p-1.5 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 rounded-full">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
              </svg>
            </button>
            <template #overlay>
              <Menu>
                <!-- 修改群内昵称 - 自己或管理员可操作 -->
                <Menu.Item v-if="member.UserID === currentUserId || isAdmin || isCreator" @click="openNicknameModal(member)">
                  <div class="flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                    </svg>
                    <span>修改昵称</span>
                  </div>
                </Menu.Item>
                
                <!-- 禁言/解禁 - 仅管理员可操作 -->
                <Menu.Item v-if="(isAdmin || isCreator) && member.UserID !== currentUserId && !(member.IsAdmin && !isCreator)" @click="handleMute(member)">
                  <div class="flex items-center gap-2" :class="{ 'text-orange-500': member.Status !== 2 }">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path v-if="member.Status === 2" stroke-linecap="round" stroke-linejoin="round" d="M8.5 8.5l7 7M8.5 15.5l7-7" />
                      <path v-else stroke-linecap="round" stroke-linejoin="round" d="M12 18.75a6 6 0 006-6v-1.5m-6 7.5a6 6 0 01-6-6v-1.5m6 7.5v3.75m-3.75 0h7.5M12 15.75a3 3 0 01-3-3V4.5a3 3 0 116 0v8.25a3 3 0 01-3 3z" />
                    </svg>
                    <span>{{ member.Status === 2 ? '解除禁言' : '禁言' }}</span>
                  </div>
                </Menu.Item>
                
                <!-- 设置/取消管理员 - 仅创建者可操作 -->
                <Menu.Item v-if="isCreator && member.UserID !== currentUserId" @click="handleToggleAdmin(member)">
                  <div class="flex items-center gap-2" :class="{ 'text-blue-500': !member.IsAdmin }">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                    </svg>
                    <span>{{ member.IsAdmin ? '取消管理员' : '设为管理员' }}</span>
                  </div>
                </Menu.Item>
                
                <!-- 移除成员 - 仅管理员可操作 -->
                <Menu.Item v-if="(isAdmin || isCreator) && member.UserID !== currentUserId && !(member.IsAdmin && !isCreator)" @click="handleRemove(member)">
                  <div class="flex items-center gap-2 text-red-500">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M13 7a4 4 0 11-8 0 4 4 0 018 0zM9 14a6 6 0 00-6 6v1h12v-1a6 6 0 00-6-6zM21 12h-6" />
                    </svg>
                    <span>移出群组</span>
                  </div>
                </Menu.Item>
              </Menu>
            </template>
          </Dropdown>
        </div>
      </div>
      
      <!-- 空列表提示 -->
      <div v-if="filteredMembers.length === 0" class="py-8 text-center text-gray-500 dark:text-gray-400">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-300 dark:text-gray-600 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        <p v-if="searchText">没有找到匹配的成员</p>
        <p v-else>群组暂无成员</p>
      </div>
    </div>
    
    <!-- 修改昵称模态框 -->
    <Modal
      v-model:visible="nicknameModalVisible"
      title="修改群内昵称"
      @ok="submitNicknameChange"
      okText="保存"
      cancelText="取消"
    >
      <div class="space-y-4">
        <p class="text-gray-500 dark:text-gray-400">为 {{ editingMember?.Nickname || editingMember?.Username }} 设置新的群内昵称</p>
        <Input v-model:value="newNickname" placeholder="输入新昵称" maxlength="20" showCount />
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.theme-border {
  border-color: var(--border-color);
}
</style> 