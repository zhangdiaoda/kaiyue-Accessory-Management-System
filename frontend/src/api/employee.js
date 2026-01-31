import request from '@/utils/request'

/**
 * 获取员工列表
 */
export function getEmployeeList(params) {
    return request({
        url: '/employees',
        method: 'get',
        params
    })
}

/**
 * 获取所有在职员工
 */
export function getAllEmployees() {
    return request({
        url: '/employees/all',
        method: 'get'
    })
}

/**
 * 创建员工
 */
export function createEmployee(data) {
    return request({
        url: '/employees',
        method: 'post',
        data
    })
}

/**
 * 更新员工
 */
export function updateEmployee(id, data) {
    return request({
        url: `/employees/${id}`,
        method: 'put',
        data
    })
}

/**
 * 删除员工
 */
export function deleteEmployee(id) {
    return request({
        url: `/employees/${id}`,
        method: 'delete'
    })
}
/**
 * 下载导入模板
 */
export const downloadTemplateUrl = '/api/employees/template'

/**
 * 导出员工
 */
export const exportEmployeesUrl = '/api/employees/export'

/**
 * 导入员工
 */
export function importEmployees(data) {
    return request({
        url: '/employees/import',
        method: 'post',
        data,
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    })
}
