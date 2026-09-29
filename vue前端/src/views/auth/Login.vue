<template>
  <div class="min-h-screen bg-gray-100 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <img class="mx-auto h-12 w-auto" src="@/assets/img/logo.webp" alt="研途智伴" />
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">登录到 研途智伴</h2>
      <p class="mt-2 text-center text-sm text-gray-600">
        还没有账号？
        <router-link to="/auth/register" class="font-medium text-blue-600 hover:text-blue-500">
          立即注册
        </router-link>
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
        <a-form
          :model="formState"
          name="login"
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

          <!-- 密码输入框 -->
          <a-form-item
            name="password"
            :rules="[{ required: true, message: '请输入密码！' }]"
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

          <!-- 记住我和忘记密码 -->
          <div class="flex items-center justify-between">
            <a-form-item name="remember" class="mb-0">
              <a-checkbox v-model:checked="formState.remember">
                记住我
              </a-checkbox>
            </a-form-item>

            <div class="text-sm">
              <a href="#" class="font-medium text-blue-600 hover:text-blue-500">
                忘记密码？
              </a>
            </div>
          </div>

          <!-- 登录按钮 -->
          <div>
            <a-button
              type="primary"
              html-type="submit"
              class="w-full"
              size="large"
              :loading="loading"
            >
              登录
            </a-button>
          </div>
        </a-form>

        <div class="mt-6">
          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-gray-300"></div>
            </div>
            <div class="relative flex justify-center text-sm">
              <span class="px-2 bg-white text-gray-500">或者使用以下方式登录</span>
            </div>
          </div>

          <div class="mt-6 grid grid-cols-3 gap-3">
            <div>
              <a href="#" class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.17 6.839 9.49.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.604-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.462-1.11-1.462-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.137 20.167 22 16.418 22 12c0-5.523-4.477-10-10-10z" />
                </svg>
              </a>
            </div>

            <div>
              <a href="#" class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm6.36 14.83c-1.43-1.74-4.9-2.33-6.36-2.33s-4.93.59-6.36 2.33C4.62 15.49 4 13.82 4 12c0-4.41 3.59-8 8-8s8 3.59 8 8c0 1.82-.62 3.49-1.64 4.83zM12 6c-1.94 0-3.5 1.56-3.5 3.5S10.06 13 12 13s3.5-1.56 3.5-3.5S13.94 6 12 6z" />
                </svg>
              </a>
            </div>

            <div>
              <a href="#" class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M23.643 4.937c-.835.37-1.732.62-2.675.733.962-.576 1.7-1.49 2.048-2.578-.9.534-1.897.922-2.958 1.13-.85-.904-2.06-1.47-3.4-1.47-2.572 0-4.658 2.086-4.658 4.66 0 .364.042.718.12 1.06-3.873-.195-7.304-2.05-9.602-4.868-.4.69-.63 1.49-.63 2.342 0 1.616.823 3.043 2.072 3.878-.764-.025-1.482-.234-2.11-.583v.06c0 2.257 1.605 4.14 3.737 4.568-.392.106-.803.162-1.227.162-.3 0-.593-.028-.877-.082.593 1.85 2.313 3.198 4.352 3.234-1.595 1.25-3.604 1.995-5.786 1.995-.376 0-.747-.022-1.112-.065 2.062 1.323 4.51 2.093 7.14 2.093 8.57 0 13.255-7.098 13.255-13.254 0-.2-.005-.402-.014-.602.91-.658 1.7-1.477 2.323-2.41z" />
                </svg>
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { postJSON } from '@/utility/request'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)

const formState = reactive({
  username: '',
  password: '',
  remember: true
})

const onFinish = async (values) => {
  loading.value = true
  try {
    // 使用 postJSON 发送登录请求
    postJSON('/api/auth/login', values, 
      (msg, data) => {
        // 登录成功，保存 token 和用户信息
        localStorage.setItem('authToken', data)
        userStore.login(values)
        message.success('登录成功！')
        router.push('/')
      },
      (msg) => {
        message.error(msg || '登录失败，请检查用户名和密码！')
      }
    )
  } catch (error) {
    message.error('登录请求发生错误！')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const onFinishFailed = (errorInfo) => {
  console.log('Failed:', errorInfo)
}
</script> 