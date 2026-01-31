import request from '@/utils/request'

// 获取所有通知配置
export const getNotificationConfigs = () => {
    return request({
        url: '/notifications/configs',
        method: 'get'
    })
}

// 获取指定类型的通知配置
export const getNotificationConfig = (type) => {
    return request({
        url: `/notifications/config/${type}`,
        method: 'get'
    })
}

// 更新通知配置
export const updateNotificationConfig = (data) => {
    return request({
        url: '/notifications/config',
        method: 'post',
        data
    })
}

// 测试钉钉通知
export const testDingTalk = (data) => {
    return request({
        url: '/notifications/test/dingtalk',
        method: 'post',
        data
    })
}

// 测试微信通知
export const testWechat = (data) => {
    return request({
        url: '/notifications/test/wechat',
        method: 'post',
        data
    })
}

// 获取通知日志
export const getNotificationLogs = (params) => {
    return request({
        url: '/notifications/logs',
        method: 'get',
        params
    })
}

// 获取统计信息
export const getNotificationStats = () => {
    return request({
        url: '/notifications/stats',
        method: 'get'
    })
}

// 手动发送通知
export const sendNotification = (data) => {
    return request({
        url: '/notifications/send',
        method: 'post',
        data
    })
}

// 获取用户通知设置
export const getUserNotificationSettings = () => {
    return request({
        url: '/notifications/user/settings',
        method: 'get'
    })
}

// 更新用户通知设置
export const updateUserNotificationSetting = (data) => {
    return request({
        url: '/notifications/user/setting',
        method: 'post',
        data
    })
}

// 获取微信绑定信息
export const getWechatBinding = () => {
    return request({
        url: '/notifications/wechat/binding',
        method: 'get'
    })
}

// 绑定微信用户
export const bindWechatUser = (data) => {
    return request({
        url: '/notifications/wechat/bind',
        method: 'post',
        data
    })
}

// 更新订阅场景
export const updateSubscribeScenes = (data) => {
    return request({
        url: '/notifications/wechat/subscribe',
        method: 'post',
        data
    })
}

// 立即运行每日报表
export const runDailyReport = () => {
    return request({
        url: '/notifications/run/daily-report',
        method: 'post'
    })
}

// 立即运行超期检查
export const runOverdueCheck = () => {
    return request({
        url: '/notifications/run/overdue-check',
        method: 'post'
    })
}

// 获取调度配置
export const getScheduleConfigs = () => {
    return request({
        url: '/notifications/schedules',
        method: 'get'
    })
}

// 更新调度配置
export const updateScheduleConfig = (data) => {
    return request({
        url: '/notifications/schedule',
        method: 'post',
        data
    })
}
