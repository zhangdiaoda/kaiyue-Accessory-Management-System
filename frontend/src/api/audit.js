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
