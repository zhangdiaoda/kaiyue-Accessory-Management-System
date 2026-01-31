import request from '@/utils/request'

/**
 * 补货/入库
 * @param {Object} data { part_id, quantity, remark, batch_no }
 */
export function restockPart(data) {
    return request({
        url: '/inbound',
        method: 'post',
        data
    })
}

/**
 * 获取入库流水
 * @param {Object} params { page, pageSize, part_name }
 */
export function getInboundList(params) {
    return request({
        url: '/inbound',
        method: 'get',
        params
    })
}
