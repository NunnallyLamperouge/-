import { ref, computed, onMounted } from 'vue'
import { get, postJSON } from '@/utility/request'
import { message } from 'ant-design-vue'
import { useUserStore } from '@/stores/user'

/**
 * 群组聊天列表逻辑
 * @param {Object} props - 组件属性
 * @param {Function} emit - 事件发射器
 * @returns {Object} 包含所有群组列表所需的状态和方法
 */
export function useGroupChatList(props, emit) {
  const userStore = useUserStore()
  const groups = ref([])
  const loading = ref(false)
  const searchText = ref('')
  const createModalVisible = ref(false)
  const createGroupForm = ref({
    name: '',
    description: '',
    avatarUrl: ''
  })

  // 获取群组列表
  const fetchGroups = () => {
    loading.value = true
    get('/api/chat/groups', {}, (_, data) => {
      groups.value = data.groups
      loading.value = false
    }, () => {
      loading.value = false
      message.error('获取群组失败')
    })
  }

  // 过滤搜索结果
  const filteredGroups = computed(() => {
    if (!searchText.value) return groups.value
    
    const query = searchText.value.toLowerCase()
    return groups.value.filter(group => 
      group.Name.toLowerCase().includes(query) || 
      (group.Description && group.Description.toLowerCase().includes(query))
    )
  })

  // 选择群组
  const selectGroup = (group) => {
    emit('select-group', group)
  }

  // 打开创建群组模态框
  const openCreateModal = () => {
    createGroupForm.value = {
      name: '',
      description: '',
      avatarUrl: ''
    }
    createModalVisible.value = true
  }

  // 创建群组
  const handleCreateGroup = () => {
    if (!createGroupForm.value.name) {
      message.warning('请输入群组名称')
      return
    }
    
    // 提交创建请求
    postJSON('/api/chat/groups', {
      name: createGroupForm.value.name,
      description: createGroupForm.value.description,
      avatarUrl: createGroupForm.value.avatarUrl
    }, (_, data) => {
      message.success('群组创建成功')
      createModalVisible.value = false
      
      // 刷新群组列表或直接添加到列表
      if (data) {
        groups.value.unshift(data)
        emit('create-group', data)
      } else {
        fetchGroups()
      }
    })
  }

  // 组件挂载时获取群组列表
  onMounted(() => {
    fetchGroups()
  })

  return {
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
  }
} 