import './assets/main.css'
import { MotionPlugin } from '@vueuse/motion'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/reset.css';
import 'animate.css'
import axios from "axios";
import { useThemeStore } from './stores/theme'

axios.defaults.baseURL="http://localhost:8080"

const app = createApp(App)
const pinia = createPinia()
app.use(pinia).use(MotionPlugin)
const themeStore = useThemeStore()
themeStore.initTheme()
app.use(router)
app.use(ElementPlus).use(Antd)
app.mount('#app')
