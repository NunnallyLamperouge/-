<template>
  <div class="min-h-screen bg-gray-100 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <img class="mx-auto h-12 w-auto" src="@/assets/img/logo.webp" alt="研途智伴" />
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">注册 研途智伴</h2>
      <p class="mt-2 text-center text-sm text-gray-600">
        已有账号？
        <router-link to="/auth/login" class="font-medium text-blue-600 hover:text-blue-500">
          立即登录
        </router-link>
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
        <a-form
          :model="formState"
          name="register"
          @finish="onFinish"
          @finishFailed="onFinishFailed"
          class="space-y-6"
        >
          <!-- 用户名输入框 -->
          <a-form-item
            name="username"
            :rules="[{ required: true, message: '请输入用户名！' }]"
          >
            <a-input
              v-model:value="formState.username"
              size="large"
              placeholder="用户名"
            >
              <template #prefix>
                <user-outlined />
              </template>
            </a-input>
          </a-form-item>

          <!-- 邮箱输入框 -->
          <a-form-item
            name="email"
            :rules="[
              { required: true, message: '请输入邮箱地址！' },
              { type: 'email', message: '请输入有效的邮箱地址！' }
            ]"
          >
            <a-input
              v-model:value="formState.email"
              size="large"
              placeholder="邮箱地址"
            >
              <template #prefix>
                <mail-outlined />
              </template>
            </a-input>
          </a-form-item>

          <!-- 验证码输入框 -->
          <a-form-item
            name="code"
            :rules="[{ required: true, message: '请输入验证码！' }]"
          >
            <div class="flex space-x-2">
              <a-input
                v-model:value="formState.code"
                size="large"
                placeholder="验证码"
                class="flex-1"
              >
                <template #prefix>
                  <safety-outlined />
                </template>
              </a-input>
              <a-button 
                type="primary" 
                size="large" 
                :disabled="codeSending || cooldown > 0" 
                @click="sendVerifyCode"
              >
                {{ cooldown > 0 ? `${cooldown}秒后重试` : '获取验证码' }}
              </a-button>
            </div>
          </a-form-item>

          <!-- 密码输入框 -->
          <a-form-item
            name="password"
            :rules="[
              { required: true, message: '请输入密码！' },
              { min: 8, message: '密码长度至少为8个字符！' }
            ]"
          >
            <a-input-password
              v-model:value="formState.password"
              size="large"
              placeholder="密码"
            >
              <template #prefix>
                <lock-outlined />
              </template>
            </a-input-password>
          </a-form-item>

          <!-- 确认密码输入框 -->
          <a-form-item
            name="confirmPassword"
            :rules="[
              { required: true, message: '请确认密码！' },
              { validator: validateConfirmPassword }
            ]"
          >
            <a-input-password
              v-model:value="formState.confirmPassword"
              size="large"
              placeholder="确认密码"
            >
              <template #prefix>
                <lock-outlined />
              </template>
            </a-input-password>
          </a-form-item>

          <!-- 用户类型选择 -->
          <a-form-item
            name="userType"
            :rules="[{ required: true, message: '请选择用户类型！' }]"
          >
            <a-select
              v-model:value="formState.userType"
              size="large"
              placeholder="选择用户类型"
            >
              <a-select-option value="student">学生</a-select-option>
              <a-select-option value="teacher">指导老师</a-select-option>
            </a-select>
          </a-form-item>

          <!-- 注册按钮 -->
          <div>
            <a-button
              type="primary"
              html-type="submit"
              class="w-full"
              size="large"
              :loading="loading"
            >
              注册
            </a-button>
          </div>
        </a-form>
      </div>
    </div>

    <!-- 用户协议弹窗 -->
    <a-modal
      v-model:visible="agreementVisible"
      title="用户协议"
      width="700px"
      @ok="agreementVisible = false"
    >
      <div class="prose max-w-none">
        <h3>研途智伴 用户协议</h3>
        <p>欢迎使用 研途智伴！本协议是您与 研途智伴 平台之间的法律协议...</p>
        <!-- 这里添加更多协议内容 -->
      </div>
    </a-modal>

    <!-- 隐私政策弹窗 -->
    <a-modal
      v-model:visible="privacyVisible"
      title="隐私政策"
      width="700px"
      @ok="privacyVisible = false"
    >
      <div class="prose max-w-none">
        <h3>研途智伴 隐私政策</h3>
        <p>我们非常重视您的个人隐私和数据安全...</p>
        <!-- 这里添加更多隐私政策内容 -->
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { MailOutlined, LockOutlined, UserOutlined, SafetyOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { postJSON } from '@/utility/request'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const codeSending = ref(false)
const cooldown = ref(0)
const agreementVisible = ref(false)
const privacyVisible = ref(false)

const formState = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  userType: undefined,
  code: ''
})

const validateConfirmPassword = async (rule, value) => {
  if (value !== formState.password) {
    throw new Error('两次输入的密码不一致！')
  }
}

const showAgreement = () => {
  agreementVisible.value = true
}

const showPrivacy = () => {
  privacyVisible.value = true
}

// 发送验证码
const sendVerifyCode = async () => {
  // 验证邮箱格式
  if (!formState.email) {
    message.warning('请先输入邮箱地址！')
    return
  }
  
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(formState.email)) {
    message.warning('请输入有效的邮箱地址！')
    return
  }
  
  codeSending.value = true
  
  try {
    postJSON('/api/auth/sendCode', { email: formState.email }, 
      (msg) => {
        message.success(msg)
        // 设置冷却时间
        cooldown.value = 60
        const timer = setInterval(() => {
          cooldown.value--
          if (cooldown.value <= 0) {
            clearInterval(timer)
          }
        }, 1000)
      },
      (msg) => {
        message.error(msg || '发送验证码失败')
      }
    )
  } finally {
    codeSending.value = false
  }
}

const onFinish = async (values) => {
  loading.value = true
  try {
    // 构建注册请求数据
    const registerData = {
      username: formState.username,
      password: formState.password,
      email: formState.email,
      code: formState.code,
      type: formState.userType === 'student' ? 0 : 1 // 学生为0，老师为1
    }
    
    // 发送注册请求
    postJSON('/api/auth/register', registerData, 
      (msg) => {
        message.success(msg || '注册成功！')
        // 注册成功后跳转到登录页
        router.push('/auth/login')
      },
      (msg) => {
        message.error(msg || '注册失败，请重试！')
      }
    )
  } finally {
    loading.value = false
  }
}

const onFinishFailed = (errorInfo) => {
  console.log('Failed:', errorInfo)
}
</script> 