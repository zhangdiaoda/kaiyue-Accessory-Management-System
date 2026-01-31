import request from '@/utils/request'

/**
 * 获取备份配置
 */
export function getBackupConfig() {
    return request({
        url: '/system/backups/config',
        method: 'get'
    })
}

/**
 * 更新备份配置
 */
export function updateBackupConfig(data) {
    return request({
        url: '/system/backups/config',
        method: 'put',
        data
    })
}

/**
 * 立即执行备份
 */
export function runBackup() {
    return request({
        url: '/system/backups/run',
        method: 'post'
    })
}

/**
 * 获取备份列表
 */
export function getBackupList() {
    return request({
        url: '/system/backups',
        method: 'get'
    })
}

/**
 * 删除备份
 */
export function deleteBackup(name) {
    return request({
        url: '/system/backups/',
        method: 'delete',
        params: { name }
    })
}

/**
 * 恢复备份
 */
export function restoreBackup(name) {
    return request({
        url: '/system/backups/restore',
        method: 'post',
        data: { name }
    })
}

/**
 * 下载备份 (返回URL，通常直接window.open)
 */
export function getDownloadUrl(name) {
    const baseURL = import.meta.env.VITE_APP_BASE_API || '/api'
    const token = localStorage.getItem('token') // 简单处理，URL通常需要鉴权，这里假设Cookie或QueryToken
    // 由于是文件下载，后端 handler 需要做鉴权兼容，或者前端通过 blob 下载
    // 简单起见，这里假设后端 download 接口不做严格 Header Token 校验，或者使用 param 传递
    return `${baseURL}/system/backups/download?name=${name}&token=${token}`
}
