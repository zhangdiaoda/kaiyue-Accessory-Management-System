import request from '@/utils/request'

// 获取操作日志列表
export const getOperationLogs = (params) => {
    return request({
        url: '/operation-logs',
        method: 'get',
        params
    })
}

// 获取日志详情
export const getLogDetail = (logId) => {
    return request({
        url: `/operation-logs/${logId}`,
        method: 'get'
    })
}

// 获取操作统计
export const getOperationStats = (params) => {
    return request({
        url: '/operation-logs/stats',
        method: 'get',
        params
    })
}

// 清空操作日志
export const clearLogs = () => {
    return request({
        url: '/operation-logs/clear',
        method: 'post'
    })
}

// 获取日志表大小
export const getLogSize = () => {
    return request({
        url: '/operation-logs/size',
        method: 'get'
    })
}

// 获取清理配置
export const getCleanupConfig = () => {
    return request({
        url: '/operation-logs/cleanup-config',
        method: 'get'
    })
}

// 更新清理配置
export const updateCleanupConfig = (data) => {
    return request({
        url: '/operation-logs/cleanup-config',
        method: 'post',
        data
    })
}
