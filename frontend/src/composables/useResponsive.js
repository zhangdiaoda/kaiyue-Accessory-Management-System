import { ref, onMounted, onUnmounted, computed } from 'vue'

/**
 * 响应式设备检测组合式函数
 * 用于判断当前设备类型和屏幕尺寸
 */
export function useResponsive() {
    // 强制布局模式: 'auto' | 'mobile' | 'desktop'
    const forceMode = ref(localStorage.getItem('layout-mode') || 'auto')

    // 设备类型标识（原始检测结果）
    const actualIsMobile = ref(false)
    const isTablet = ref(false)
    const isDesktop = ref(false)
    const isLargeScreen = ref(false)

    // 当前屏幕宽度
    const screenWidth = ref(0)

    /**
     * 更新设备类型
     */
    const updateDeviceType = () => {
        const width = window.innerWidth
        screenWidth.value = width

        actualIsMobile.value = width < 768
        isTablet.value = width >= 768 && width < 1024
        isDesktop.value = width >= 1024 && width < 1920
        isLargeScreen.value = width >= 1920
    }

    // 计算最终的 isMobile 状态（考虑强制模式）
    const isMobile = computed(() => {
        if (forceMode.value === 'mobile') return true
        if (forceMode.value === 'desktop') return false
        return actualIsMobile.value
    })

    /**
     * 切换布局模式
     * @param {string} mode - 'auto' | 'mobile' | 'desktop'
     */
    const toggleLayoutMode = (mode) => {
        forceMode.value = mode
        localStorage.setItem('layout-mode', mode)
    }

    /**
     * 获取弹窗宽度
     * @param {string} size - 弹窗大小 'small' | 'medium' | 'large'
     * @returns {string} 弹窗宽度
     */
    const getDialogWidth = (size = 'medium') => {
        const width = screenWidth.value

        if (isMobile.value) {
            return '95%'
        }

        if (isTablet.value) {
            return size === 'large' ? '85%' : '70%'
        }

        if (isLargeScreen.value) {
            const widths = {
                small: '600px',
                medium: '800px',
                large: '1200px'
            }
            return widths[size] || widths.medium
        }

        // 桌面端
        const widths = {
            small: '480px',
            medium: '600px',
            large: '900px'
        }
        return widths[size] || widths.medium
    }

    /**
     * 是否应该隐藏表格列
     * @param {string} priority - 列优先级 'low' | 'medium' | 'high'
     * @returns {boolean}
     */
    const shouldHideColumn = (priority = 'medium') => {
        if (isMobile.value) {
            return priority === 'low' || priority === 'medium'
        }
        if (isTablet.value) {
            return priority === 'low'
        }
        return false
    }

    // 挂载时注册事件
    onMounted(() => {
        updateDeviceType()
        window.addEventListener('resize', updateDeviceType)
    })

    // 卸载时移除事件
    onUnmounted(() => {
        window.removeEventListener('resize', updateDeviceType)
    })

    return {
        isMobile,
        isTablet,
        isDesktop,
        isLargeScreen,
        screenWidth,
        forceMode,
        toggleLayoutMode,
        getDialogWidth,
        shouldHideColumn
    }
}
