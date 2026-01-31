import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/Login.vue'),
        meta: { title: '登录' }
    },
    {
        path: '/datav',
        name: 'DataV',
        component: () => import('@/views/datav/Index.vue'),
        meta: { title: '可视化监控大屏', requiresAuth: true }
    },
    {
        path: '/',
        name: 'Layout',
        component: () => import('@/views/Layout.vue'),
        redirect: '/dashboard',
        meta: { requiresAuth: true },
        children: [
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('@/views/dashboard/Index.vue'),
                meta: { title: '仪表盘', icon: 'DataLine' }
            },
            {
                path: 'parts',
                name: 'Parts',
                meta: { title: '配件管理', icon: 'Box' },
                redirect: '/parts/list',
                children: [
                    {
                        path: 'list',
                        name: 'PartList',
                        component: () => import('@/views/parts/PartList.vue'),
                        meta: { title: '配件列表' }
                    },
                    {
                        path: 'categories',
                        name: 'PartCategory',
                        component: () => import('@/views/parts/PartCategory.vue'),
                        meta: { title: '分类管理' }
                    },
                    {
                        path: 'warning',
                        name: 'StockWarning',
                        component: () => import('@/views/parts/StockWarning.vue'),
                        meta: { title: '库存预警' }
                    },
                    {
                        path: 'inbound',
                        name: 'InboundHistory',
                        component: () => import('@/views/parts/InboundHistory.vue'),
                        meta: { title: '入库明细' }
                    }
                ]
            },
            {
                path: 'borrow',
                name: 'Borrow',
                meta: { title: '领用管理', icon: 'DocumentCopy' },
                children: [
                    {
                        path: 'create',
                        name: 'BorrowCreate',
                        component: () => import('@/views/borrow/BorrowCreate.vue'),
                        meta: { title: '登记领用' }
                    },
                    {
                        path: 'return',
                        name: 'ReturnCreate',
                        component: () => import('@/views/borrow/ReturnCreate.vue'),
                        meta: { title: '登记归还' }
                    },
                    {
                        path: 'history',
                        name: 'BorrowHistory',
                        component: () => import('@/views/borrow/BorrowHistory.vue'),
                        meta: { title: '领用记录' }
                    },
                    {
                        path: 'old-inventory',
                        name: 'OldInventory',
                        component: () => import('@/views/borrow/OldInventory.vue'),
                        meta: { title: '旧件仓库' }
                    },
                    {
                        path: 'scrap-inventory',
                        name: 'ScrapInventory',
                        component: () => import('@/views/borrow/ScrapInventory.vue'),
                        meta: { title: '废品仓库' }
                    }
                ]
            },

            {
                path: 'reports',
                name: 'Reports',
                meta: { title: '报表中心', icon: 'DataAnalysis' },
                children: [
                    {
                        path: 'query',
                        name: 'ReportQuery',
                        component: () => import('@/views/reports/ReportQuery.vue'),
                        meta: { title: '报表查询' }
                    },
                    {
                        path: 'detailed',
                        name: 'DetailedReport',
                        component: () => import('@/views/reports/DetailedReport.vue'),
                        meta: { title: '明细查询' }
                    },
                    {
                        path: 'employee-annual',
                        name: 'EmployeeAnnualReport',
                        component: () => import('@/views/reports/EmployeeAnnualReport.vue'),
                        meta: { title: '员工年度报表' }
                    },
                    {
                        path: 'scrap-analytics',
                        name: 'ScrapAnalytics',
                        component: () => import('@/views/reports/ScrapAnalytics.vue'),
                        meta: { title: '废品损毁分析' }
                    },
                    {
                        path: 'analytics',
                        name: 'DataAnalytics',
                        component: () => import('@/views/analytics/DataAnalytics.vue'),
                        meta: { title: '数据分析' }
                    }
                ]
            },
            {
                path: 'system',
                name: 'System',
                meta: { title: '系统管理', icon: 'Setting', role: 'SUPER_ADMIN' },
                children: [
                    {
                        path: 'users',
                        name: 'UserManage',
                        component: () => import('@/views/system/UserManage.vue'),
                        meta: { title: '用户管理' }
                    },
                    {
                        path: 'employees',
                        name: 'EmployeeManage',
                        component: () => import('@/views/system/EmployeeManage.vue'),
                        meta: { title: '员工管理' }
                    },
                    {
                        path: 'config',
                        name: 'ConfigManage',
                        component: () => import('@/views/system/ConfigManage.vue'),
                        meta: { title: '系统配置' }
                    },
                    {
                        path: 'profile',
                        name: 'Profile',
                        component: () => import('@/views/system/Profile.vue'),
                        meta: { title: '个人信息' }
                    },
                    {
                        path: 'brand-config',
                        name: 'BrandConfig',
                        component: () => import('@/views/system/BrandConfig.vue'),
                        meta: { title: '品牌与公告设置' }
                    },
                    {
                        path: 'messages',
                        name: 'MessageCenter',
                        component: () => import('@/views/system/MessageCenter.vue'),
                        meta: { title: '消息中心' }
                    },
                    {
                        path: 'backup',
                        name: 'BackupManage',
                        component: () => import('@/views/system/BackupManage.vue'),
                        meta: { title: '数据库备份' }
                    },
                    {
                        path: 'notification-config',
                        name: 'NotificationConfig',
                        component: () => import('@/views/system/NotificationConfig.vue'),
                        meta: { title: '通知配置' }
                    },
                    {
                        path: 'notification-center',
                        name: 'NotificationCenter',
                        component: () => import('@/views/system/NotificationCenter.vue'),
                        meta: { title: '通知中心' }
                    },
                    {
                        path: 'notification-settings',
                        name: 'NotificationSettings',
                        component: () => import('@/views/system/NotificationSettings.vue'),
                        meta: { title: '我的通知偏好' }
                    },
                    {
                        path: 'schedule-config',
                        name: 'ScheduleConfig',
                        component: () => import('@/views/system/ScheduleConfig.vue'),
                        meta: { title: '定时任务配置' }
                    },
                    {
                        path: 'operation-log',
                        name: 'OperationLog',
                        component: () => import('@/views/system/OperationLog.vue'),
                        meta: { title: '操作日志' }
                    }
                ]
            }
        ]
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
    const userStore = useUserStore()

    // 设置页面标题
    document.title = to.meta.title ? `${to.meta.title} - 仓储管理系统` : '仓储管理系统'

    // 判断是否需要登录
    if (to.meta.requiresAuth) {
        if (userStore.token) {
            // 权限检查
            if (to.meta.role && to.meta.role !== userStore.userInfo.role) {
                next('/dashboard')
            } else {
                next()
            }
        } else {
            next('/login')
        }
    } else {
        next()
    }
})

export default router
