import request from '@/utils/request'

// 获取所有权限（按分类分组）
export const getAllPermissions = () => {
    return request({
        url: '/permissions',
        method: 'get'
    })
}

// 获取用户权限
export const getUserPermissions = (userId) => {
    return request({
        url: `/permissions/users/${userId}`,
        method: 'get'
    })
}

// 设置用户权限
export const setUserPermissions = (userId, permissions) => {
    return request({
        url: `/permissions/users/${userId}`,
        method: 'put',
        data: { permissions }
    })
}

// 获取角色默认权限
export const getRolePermissions = (role) => {
    return request({
        url: `/permissions/roles/${role}`,
        method: 'get'
    })
}

// 设置角色默认权限
export const setRolePermissions = (role, permissions) => {
    return request({
        url: `/permissions/roles/${role}`,
        method: 'put',
        data: { permissions }
    })
}
