import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import Noir from '../presets/Noir'
// 使用浏览器 localStorage 存储主题状态

export const useThemeStore = defineStore('theme', () => {
    const isDark = ref(false)

    const STORAGE_KEY = 'theme-dark'

    function saveTheme(dark: boolean) {
        try {
            localStorage.setItem(STORAGE_KEY, dark ? 'true' : 'false')
        } catch (e) {
            console.warn('Failed to save theme to localStorage:', e)
        }
    }

    function loadSavedTheme(): boolean | null {
        try {
            const v = localStorage.getItem(STORAGE_KEY)
            if (v === 'true') return true
            if (v === 'false') return false
            return null
        } catch (e) {
            console.warn('Failed to load theme from localStorage:', e)
            return null
        }
    }

    // 当前主题的图标配置
    const currentTheme = computed(() => {
        return {
            icons: isDark.value ? Noir.semantic.icons.dark : Noir.semantic.icons.light
        }
    })

    // 切换主题
    async function toggleTheme() {
        const root = document.getElementsByTagName('html')[0]
        root.classList.toggle('p-dark')
        isDark.value = root.classList.contains('p-dark')
        // 保存到 localStorage
        saveTheme(isDark.value)
    }

    // 设置主题
    async function setTheme(dark: boolean) {
        const root = document.getElementsByTagName('html')[0]
        if (dark) {
            root.classList.add('p-dark')
        } else {
            root.classList.remove('p-dark')
        }
        isDark.value = dark
        // 保存到 localStorage
        saveTheme(dark)
    }

    // 初始化主题状态
    async function initTheme() {
        // 从 localStorage 读取主题状态
        const saved = loadSavedTheme()
        if (saved !== null) {
            await setTheme(saved)
        } else {
            const root = document.getElementsByTagName('html')[0]
            isDark.value = root.classList.contains('p-dark')
        }
    }

    return {
        isDark,
        currentTheme,
        toggleTheme,
        setTheme,
        initTheme
    }
})
