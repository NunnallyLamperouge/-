import {defineStore} from 'pinia'

export const useThemeStore = defineStore('theme', {
    state: () => ({
        isDark: false,
    }),

    actions: {
        initTheme() {
            // 优先从本地存储获取主题设置
            const savedTheme = localStorage.getItem('theme')
            if (savedTheme) {
                this.isDark = savedTheme === 'dark'
            } else {
                // 如果没有保存的主题设置，则检查系统主题偏好
                this.isDark = window.matchMedia('(prefers-color-scheme: dark)').matches
            }
            // 初始化时应用主题
            if (this.isDark) {
                document.documentElement.classList.add('dark')
            }
        },

        toggleTheme() {
            this.isDark = !this.isDark
            // 更新 HTML 根元素的 class
            if (this.isDark) {
                document.documentElement.classList.add('dark')
            } else {
                document.documentElement.classList.remove('dark')
            }
            // 保存到本地存储
            localStorage.setItem('theme', this.isDark ? 'dark' : 'light')
        },

        applyTheme() {
            const root = document.documentElement
            if (this.isDark) {
                root.classList.add('dark')
            } else {
                root.classList.remove('dark')
            }
            // 添加过渡类
            root.classList.add('theme-transition')
            // 保存主题到本地存储
            localStorage.setItem('theme', this.isDark ? 'dark' : 'light')

            // 可选：在过渡完成后移除过渡类
            setTimeout(() => {
                root.classList.remove('theme-transition')
            }, 300) // 与 CSS 中的 duration 保持一致
        }
    },

    getters: {
        currentTheme: (state) => state.isDark ? 'dark' : 'light'
    }
})