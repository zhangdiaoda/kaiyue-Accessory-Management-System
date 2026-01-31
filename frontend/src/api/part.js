import request from '@/utils/request'

/**
 * 获取配件列表
 */
export function getPartList(params) {
    return request({
        url: '/parts',
        method: 'get',
        params
    })
}

/**
 * 创建配件
 */
export function createPart(data) {
    return request({
        url: '/parts',
        method: 'post',
        data
    })
}

/**
 * 更新配件
 */
export function updatePart(id, data) {
    return request({
        url: `/parts/${id}`,
        method: 'put',
        data
    })
}

/**
 * 删除配件
 */
export function deletePart(id) {
    return request({
        url: `/parts/${id}`,
        method: 'delete'
    })
}

/**
 * 获取低库存配件
 */
export function getLowStockParts() {
    return request({
        url: '/parts/low-stock',
        method: 'get'
    })
}

/**
 * 获取模板下载链接
 */
export function getTemplateUrl() {
    const baseURL = import.meta.env.VITE_APP_BASE_API || '/api'
    const token = localStorage.getItem('token')
    return `${baseURL}/parts/template?token=${token}`
}

/**
 * 获取导出链接
 */
export function getExportUrl() {
    const baseURL = import.meta.env.VITE_APP_BASE_API || '/api'
    const token = localStorage.getItem('token')
    return `${baseURL}/parts/export?token=${token}`
}

/**
 * 导入配件
 */
export function importParts(data) {
    return request({
        url: '/parts/import',
        method: 'post',
        data,
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    })
}
