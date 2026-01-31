import request from '@/utils/request'

/**
 * 发送报表到钉钉
 */
export function sendToDingTalk(data) {
    return request({
        url: '/dingtalk/send',
        method: 'post',
        data
    })
}

/**
 * 测试Webhook
 */
export function testWebhook(data) {
    return request({
        url: '/dingtalk/test',
        method: 'post',
        data
    })
}
