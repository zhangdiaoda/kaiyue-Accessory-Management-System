import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import router from '@/router'

// 创建axios实例
const request = axios.create({
    baseURL: '/api',
    timeout: 30000
})

// 请求拦截器
request.interceptors.request.use(
    config => {
        const userStore = useUserStore()
        if (userStore.token) {
            config.headers['Authorization'] = `Bearer ${userStore.token}`
        }
        return config
    },
    error => {
        console.error('请求错误:', error)
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    response => {
        const res = response.data

        // 如果返回的状态码不是200，则判断为错误
        if (res.code !== 200) {
            ElMessage.error(res.message || '请求失败')

            // 401: 未登录或token过期
            if (res.code === 401) {
                const userStore = useUserStore()
                userStore.logout()
                ElMessage.warning('登录已失效,请重新登录')
                // 使用replace避免用户点击浏览器后退按钮回到错误页面
                router.replace('/login')
            }

            return Promise.reject(new Error(res.message || '请求失败'))
        } else {
            return res
        }
    },
    error => {
        console.error('响应错误:', error)

        if (error.response) {
            switch (error.response.status) {
                case 401:
                    const userStore = useUserStore()
                    userStore.logout()
                    ElMessage.warning('登录已失效,请重新登录')
                    // 使用replace避免用户点击浏览器后退按钮回到错误页面
                    router.replace('/login')
                    break
                case 403:
                    ElMessage.error('拒绝访问')
                    break
                case 404:
                    ElMessage.error('请求的资源不存在')
                    break
                case 500:
                    ElMessage.error('服务器内部错误')
                    break
                default:
                    ElMessage.error(error.response.data.message || '请求失败')
            }
        } else {
            ElMessage.error('网络错误，请检查网络连接')
        }

        return Promise.reject(error)
    }
)

export default request
