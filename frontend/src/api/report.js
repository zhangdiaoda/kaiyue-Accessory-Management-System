import request from '@/utils/request'

/**
 * 按配件统计
 */
export function getPartReport(params) {
    return request({
        url: '/reports/by-part',
        method: 'get',
        params
    })
}

/**
 * 按员工统计
 */
export function getEmployeeReport(params) {
    return request({
        url: '/reports/by-employee',
        method: 'get',
        params
    })
}

/**
 * 按部门统计
 */
export function getDepartmentReport(params) {
    return request({
        url: '/reports/by-department',
        method: 'get',
        params
    })
}

/**
 * 获取详细报表（每人每月每产品）
 */
export function getDetailedReport(params) {
    return request({
        url: '/reports/detailed',
        method: 'get',
        params
    })
}

/**
 * 推送统计报表到通知渠道 (钉钉/微信等)
 */
export function pushReport(params) {
    return request({
        url: '/reports/push',
        method: 'post',
        params
    })
}
