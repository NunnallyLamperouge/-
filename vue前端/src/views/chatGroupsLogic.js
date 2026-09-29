import { ref, computed, onMounted } from 'vue'
import { get, post, postJSON } from '@/utility/request'
import { message, Modal } from 'ant-design-vue'
import { useUserStore } from '@/stores/user'
import { BACKEND_DOMAIN } from '@/utility/config'

/**
 * 群组聊天页面逻辑
 * @returns {Object} 包含所有群组聊天页面所需的状态和方法
 */
export function useChatGroups() {
  const userStore = useUserStore()
  const activeGroup = ref(null)
  const members = ref([])
  const showMembers = ref(false)
  const isCreator = computed(() => activeGroup.value?.CreatorID === userStore.user?.id)
  const isAdmin = ref(false)
  const settingsVisible = ref(false)
  const inviteModalVisible = ref(false)
  const inviteEmail = ref('')
  const inviteActiveTab = ref('userList')
  const joinModalVisible = ref(false)
  const joinGroupCode = ref('')
  const groupSettingsForm = ref({
    name: '',
    description: '',
    avatarUrl: ''
  })
  const allUsers = ref([])
  const selectedUserIds = ref([])
  const userTableLoading = ref(false)
  const showBackToTop = ref(false)

  // 回到顶部
  const scrollToTop = (elementSelector) => {
    const element = document.querySelector(elementSelector)
    if (element) {
      element.scrollTop = 0
      showBackToTop.value = false
    }
  }

  // 监听滚动
  const handleScroll = (event) => {
    showBackToTop.value = event.target.scrollTop > 200
  }

  // 获取群组成员
  const fetchGroupMembers = () => {
    if (!activeGroup.value) return
    
    get(`/api/chat/groups/${activeGroup.value.ID}/members`, {}, (_, data) => {
      members.value = data.members
      // 检查当前用户是否为管理员
      const currentUserMember = members.value.find(m => m.UserID === userStore.user.id)
      isAdmin.value = currentUserMember?.IsAdmin || isCreator.value
    }, (msg) => {
      message.error('获取成员失败' + msg)
    })
  }

  // 选择群组时，获取群组详情和成员
  const handleSelectGroup = (group) => {
    // 获取群组详情
    get(`/api/chat/groups/${group.ID}`, {}, (_, data) => {
      activeGroup.value = data.group
      fetchGroupMembers()
    })
  }

  // 创建群组
  const handleCreateGroup = (group) => {
    postJSON('/api/chat/groups', group, (_, data) => {
      message.success('群组创建成功')
      activeGroup.value = data
      fetchGroupMembers()
    })
  }

  // 打开群组设置
  const openGroupSettings = () => {
    if (!activeGroup.value) return
    
    groupSettingsForm.value = {
      name: activeGroup.value.Name,
      description: activeGroup.value.Description || '',
      avatarUrl: activeGroup.value.AvatarURL || ''
    }
    
    settingsVisible.value = true
  }

  // 更新群组信息
  const updateGroupSettings = () => {
    if (!activeGroup.value) return
    
    put(`/api/chat/groups/${activeGroup.value.ID}`, groupSettingsForm.value, () => {
      message.success('群组信息更新成功')
      // 更新本地群组数据
      activeGroup.value = {
        ...activeGroup.value,
        Name: groupSettingsForm.value.name,
        Description: groupSettingsForm.value.description,
        AvatarURL: groupSettingsForm.value.avatarUrl
      }
      settingsVisible.value = false
    })
  }

  // 打开邀请成员模态框
  const openInviteModal = () => {
    inviteEmail.value = ''
    selectedUserIds.value = []
    inviteModalVisible.value = true
    fetchAllUsers()
  }

  // 获取所有用户列表
  const fetchAllUsers = () => {
    userTableLoading.value = true
    get('/api/user/all', {}, (_, data) => {
      allUsers.value = data.users.filter(user => {
        // 排除当前群组成员
        const memberIds = members.value.map(m => m.UserID)
        return !memberIds.includes(user.id)
      })
      userTableLoading.value = false
    }, (msg) => {
      message.error('获取用户列表失败' + msg)
      userTableLoading.value = false
    })
  }

  // 邀请成员加入群组
  const inviteMember = () => {
    if (!activeGroup.value) return
    
    // 通过邮箱邀请的情况
    if (inviteActiveTab.value === 'email') {
      if (!inviteEmail.value) {
        message.warning('请输入邮箱地址')
        return
      }
      
      post(`/api/chat/groups/${activeGroup.value.ID}/members`, {
        email: inviteEmail.value
      }, () => {
        message.success('邀请已发送')
        inviteModalVisible.value = false
        fetchGroupMembers()
      })
      return
    }
    
    // 通过选择用户列表邀请的情况
    if (selectedUserIds.value.length === 0) {
      message.warning('请选择要邀请的用户')
      return
    }
    
    postJSON(`/api/chat/groups/${activeGroup.value.ID}/members`, {
      userIds: selectedUserIds.value
    }, () => {
      message.success('邀请成功')
      inviteModalVisible.value = false
      fetchGroupMembers()
    }, (msg) => {
      message.error('邀请失败: ' + msg)
    })
  }

  // 打开加入群组模态框
  const openJoinModal = () => {
    joinGroupCode.value = ''
    joinModalVisible.value = true
  }

  // 加入群组
  const joinGroup = () => {
    if (!joinGroupCode.value) return
    
    // 发送加入请求
    post(`/api/chat/groups/join`, {
      code: joinGroupCode.value
    }, (_, data) => {
      message.success('成功加入群组')
      joinModalVisible.value = false
      // 设置当前活动群组
      activeGroup.value = data
      fetchGroupMembers()
    })
  }

  // 离开群组
  const leaveGroup = () => {
    if (!activeGroup.value) return
    
    Modal.confirm({
      title: '确认离开群组',
      content: '确定要离开该群组吗？',
      okText: '确认',
      cancelText: '取消',
      onOk: () => {
        post(`/api/chat/groups/${activeGroup.value.ID}/leave`, {}, () => {
          message.success('已离开群组')
          activeGroup.value = null
        })
      }
    })
  }

  // 删除群组
  const deleteGroup = () => {
    if (!activeGroup.value || !isCreator.value) return
    
    Modal.confirm({
      title: '确认删除群组',
      content: '此操作不可恢复，确定要删除该群组吗？',
      okText: '删除',
      okType: 'danger',
      cancelText: '取消',
      onOk: () => {
        del(`/api/chat/groups/${activeGroup.value.ID}`, {}, () => {
          message.success('群组已删除')
          activeGroup.value = null
        })
      }
    })
  }

  // 禁言/解禁成员
  const toggleMuteMember = (member, isMuted) => {
    if (!activeGroup.value || !isAdmin.value) return
    
    post(`/api/chat/groups/${activeGroup.value.ID}/members/${member.UserID}/mute`, {
      status: isMuted ? 2 : 1 // 状态：1-正常，2-禁言
    }, () => {
      message.success(isMuted ? '已禁言该成员' : '已解除禁言')
      fetchGroupMembers()
    })
  }

  // 移除成员
  const removeMember = (member) => {
    if (!activeGroup.value || !isAdmin.value) return
    
    Modal.confirm({
      title: '确认移除成员',
      content: `确定要将 ${member.Nickname || member.Username} 移出群组吗？`,
      okText: '移除',
      okType: 'danger',
      cancelText: '取消',
      onOk: () => {
        del(`/api/chat/groups/${activeGroup.value.ID}/members/${member.UserID}`, {}, () => {
          message.success('成员已移除')
          fetchGroupMembers()
        })
      }
    })
  }

  // 设置/取消管理员
  const toggleAdminStatus = (member, isAdminStatus) => {
    if (!activeGroup.value || !isCreator.value) return
    
    post(`/api/chat/groups/${activeGroup.value.ID}/members/${member.UserID}/admin`, {
      isAdmin: isAdminStatus
    }, () => {
      message.success(isAdminStatus ? '已设置为管理员' : '已取消管理员权限')
      fetchGroupMembers()
    })
  }

  // 更新成员昵称
  const updateMemberNickname = (member, nickname) => {
    if (!activeGroup.value) return
    
    put(`/api/chat/groups/${activeGroup.value.ID}/members/${member.UserID}`, {
      nickname: nickname
    }, () => {
      message.success('昵称已更新')
      fetchGroupMembers()
    })
  }

  // 用户列表表格列定义
  const userColumns = [
    {
      title: '选择',
      dataIndex: 'select',
      width: 60,
      slots: { 
        customRender: 'select' 
      }
    },
    {
      title: '用户名',
      dataIndex: 'username',
    },
    {
      title: '头像',
      dataIndex: 'avatar',
      width: 80,
      slots: { 
        customRender: 'avatar' 
      }
    }
  ]

  return {
    // 状态
    userStore,
    activeGroup,
    members,
    showMembers,
    isCreator,
    isAdmin,
    settingsVisible,
    inviteModalVisible,
    inviteEmail,
    inviteActiveTab,
    joinModalVisible,
    joinGroupCode,
    groupSettingsForm,
    allUsers,
    selectedUserIds,
    userTableLoading,
    showBackToTop,
    userColumns,
    BACKEND_DOMAIN,

    // 方法
    scrollToTop,
    handleScroll,
    fetchGroupMembers,
    handleSelectGroup,
    handleCreateGroup,
    openGroupSettings,
    updateGroupSettings,
    openInviteModal,
    fetchAllUsers,
    inviteMember,
    openJoinModal,
    joinGroup,
    leaveGroup,
    deleteGroup,
    toggleMuteMember,
    removeMember,
    toggleAdminStatus,
    updateMemberNickname
  }
} 