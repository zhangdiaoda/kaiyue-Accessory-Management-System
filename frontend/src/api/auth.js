import request from '@/utils/request'

/**
 * 登录
 */
export function login(data) {
    return request({
        url: '/auth/login',
        method: 'post',
        data
    })
}

/**
 * 获取当前用户信息
 */
export function getUserInfo() {
    return request({
        url: '/auth/userinfo',
        method: 'get'
    })
}

/**
 * 获取所有系统用户列表
 */
export function getAllUsers() {
    return request({
        url: '/auth/users',
        method: 'get'
    })
}

/**
 * 创建用户
 */
export function createUser(data) {
    return request({
        url: '/auth/users',
        method: 'post',
        data
    })
}

/**
 * 更新用户
 */
export function updateUser(id, data) {
    return request({
        url: '/auth/users/' + id,
        method: 'put',
        data
    })
}

/**
 * 删除用户
 */
export function deleteUser(id) {
    return request({
        url: '/auth/users/' + id,
        method: 'delete'
    })
}

/**
 * 重置密码
 */
export function resetUserPassword(id, password) {
    return request({
        url: `/auth/users/${id}/reset-password`,
        method: 'put',
        data: { password }
    })
}

/**
 * 更新个人档案
 */
export function updateProfile(data) {
    return request({
        url: '/auth/profile',
        method: 'put',
        data
    })
}

/**
 * 退出登录
 */
export function logout() {
    return request({
        url: '/auth/logout',
        method: 'post'
    })
}
