import request from '@/utils/request'

/**
 * 获取领用记录列表
 */
export function getBorrowRecordList(params) {
    return request({
        url: '/borrows',
        method: 'get',
        params
    })
}

/**
 * 创建领用记录
 */
export function createBorrowRecord(data) {
    return request({
        url: '/borrows',
        method: 'post',
        data
    })
}

/**
 * 归还登记
 */
export function returnBorrowRecord(id, data) {
    return request({
        url: `/borrows/${id}/return`,
        method: 'post',
        data
    })
}

/**
 * 检查未归还记录
 */
export function checkUnreturned(params) {
    return request({
        url: '/borrows/check-unreturned',
        method: 'get',
        params
    })
}
/**
 * 获取旧件库列表 (带溯源)
 */
export function getOldInventory() {
    return request({
        url: '/borrows/old-inventory',
        method: 'get'
    })
}

/**
 * 获取废品库列表 (带溯源)
 */
export function getScrapInventory() {
    return request({
        url: '/borrows/scrap-inventory',
        method: 'get'
    })
}
