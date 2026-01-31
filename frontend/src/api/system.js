import request from '@/utils/request'

/**
 * 获取系统配置
 */
export function getSystemConfig() {
    return request({
        url: '/system/config',
        method: 'get'
    })
}

/**
 * 获取公开品牌配置
 */
export function getBrandingConfig() {
    return request({
        url: '/branding',
        method: 'get'
    })
}


/**
 * 更新系统配置
 */

export function updateSystemConfig(data) {
    return request({
        url: '/system/config',
        method: 'put',
        data
    })
}

/**
 * 获取公告列表
 */
export function getAnnouncements(params) {
    return request({
        url: '/announcements',
        method: 'get',
        params
    })
}

/**
 * 发布公告
 */
export function createAnnouncement(data) {
    return request({
        url: '/announcements',
        method: 'post',
        data
    })
}

/**
 * 更新公告
 */
export function updateAnnouncement(id, data) {
    return request({
        url: '/announcements/' + id,
        method: 'put',
        data
    })
}

/**
 * 删除公告
 */
export function deleteAnnouncement(id) {
    return request({
        url: '/announcements/' + id,
        method: 'delete'
    })
}
