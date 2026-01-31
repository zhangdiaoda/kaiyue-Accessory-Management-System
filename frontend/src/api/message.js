import request from '@/utils/request'

/**
 * 获取我的站内信
 */
export function getMyMessages() {
    return request({
        url: '/messages',
        method: 'get'
    })
}

/**
 * 发送站内信
 */
export function sendMessage(data) {
    return request({
        url: '/messages',
        method: 'post',
        data
    })
}

/**
 * 标记消息为已读
 */
export function markMessageRead(id) {
    return request({
        url: `/messages/${id}/read`,
        method: 'put'
    })
}

/**
 * 获取未读消息数
 */
export function getUnreadCount() {
    return request({
        url: '/messages/unread-count',
        method: 'get'
    })
}
