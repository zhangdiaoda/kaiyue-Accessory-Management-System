<template>
  <el-container class="layout-container">
    <!-- Scrolling Announcement Bar -->
    <div v-if="scrollingAnnouncement" class="announcement-bar">
      <div class="announcement-content">
        <el-icon><Bell /></el-icon>
        <div class="scroll-container">
          <div class="scroll-text">{{ scrollingAnnouncement.content }}</div>
        </div>
      </div>
    </div>

    <el-container class="main-layout" :class="{ 'has-announcement': scrollingAnnouncement }">
      <!-- 移动端抽屉侧边栏 -->
      <el-drawer
        v-if="isMobile"
        v-model="drawerVisible"
        direction="ltr"
        size="260px"
        :with-header="false"
        class="mobile-drawer"
      >
        <div class="drawer-content">
          <div class="logo">
            <span class="logo-icon">{{ brandConfig.brand_logo || '📦' }}</span>
            <h3 class="logo-text">{{ brandConfig.system_name || '仓储管理' }}</h3>
          </div>
          <el-menu :default-active="activeMenu" router class="apple-menu" @select="drawerVisible = false">
            <el-menu-item index="/dashboard">
              <el-icon><DataLine /></el-icon>
              <span>仪表盘</span>
            </el-menu-item>
            
            <el-menu-item v-if="!isMobile" @click.native="openDatavInNewWindow">
              <el-icon><Monitor /></el-icon>
              <span>可视化大屏</span>
            </el-menu-item>
            
            <el-sub-menu index="parts">
              <template #title>
                <el-icon><Box /></el-icon>
                <span>配件管理</span>
              </template>
              <el-menu-item index="/parts">配件总库</el-menu-item>
              <el-menu-item index="/parts/categories">分类管理</el-menu-item>
              <el-menu-item index="/parts/inbound">入库明细</el-menu-item>
              <el-menu-item index="/parts/warning">库存预警</el-menu-item>
            </el-sub-menu>
            
            <el-sub-menu index="borrow">
              <template #title>
                <el-icon><DocumentCopy /></el-icon>
                <span>领用管理</span>
              </template>
              <el-menu-item index="/borrow/create">登记领用</el-menu-item>
              <el-menu-item index="/borrow/history">领用记录</el-menu-item>
              <el-menu-item index="/borrow/old-inventory">旧件仓库</el-menu-item>
              <el-menu-item index="/borrow/scrap-inventory">废品仓库</el-menu-item>
            </el-sub-menu>
            
            <el-sub-menu index="reports">
              <template #title>
                <el-icon><DataAnalysis /></el-icon>
                <span>报表中心</span>
              </template>
              <el-menu-item index="/reports/query">维度查询</el-menu-item>
              <el-menu-item index="/reports/detailed">明细查询</el-menu-item>
              <el-menu-item index="/reports/employee-annual">员工年度报表</el-menu-item>
              <el-menu-item index="/reports/scrap-analytics">废品损毁分析</el-menu-item>
            </el-sub-menu>
            
            <el-sub-menu index="system" v-if="userStore.userInfo.role === 'SUPER_ADMIN'">
              <template #title>
                <el-icon><Setting /></el-icon>
                <span>系统管理</span>
              </template>
              <el-menu-item index="/system/users">用户管理</el-menu-item>
              <el-menu-item index="/system/employees">员工管理</el-menu-item>
              <el-menu-item index="/system/brand-config">品牌与公告</el-menu-item>
              <el-menu-item index="/system/backup">数据库备份</el-menu-item>
              <el-menu-item index="/system/notification-config">通知配置</el-menu-item>
              <el-menu-item index="/system/notification-center">通知中心</el-menu-item>
              <el-menu-item index="/system/config">高级参数</el-menu-item>
            </el-sub-menu>
          </el-menu>
        </div>
      </el-drawer>

      <!-- 桌面端固定侧边栏 -->
      <el-aside v-else width="200px">
        <div class="logo">
          <span class="logo-icon">{{ brandConfig.brand_logo || '📦' }}</span>
          <h3 class="logo-text">{{ brandConfig.system_name || '仓储管理' }}</h3>
        </div>
        <el-menu :default-active="activeMenu" router class="apple-menu">
          <el-menu-item index="/dashboard">
            <el-icon><DataLine /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>
          
          <el-menu-item v-if="!isMobile" @click.native="openDatavInNewWindow">
            <el-icon><Monitor /></el-icon>
            <span>可视化大屏</span>
          </el-menu-item>
          
          <el-sub-menu index="parts">
            <template #title>
              <el-icon><Box /></el-icon>
              <span>配件管理</span>
            </template>
            <el-menu-item index="/parts">配件总库</el-menu-item>
            <el-menu-item index="/parts/categories">分类管理</el-menu-item>
            <el-menu-item index="/parts/inbound">入库明细</el-menu-item>
            <el-menu-item index="/parts/warning">库存预警</el-menu-item>
          </el-sub-menu>

          
          <el-sub-menu index="borrow">
            <template #title>
              <el-icon><DocumentCopy /></el-icon>
              <span>领用管理</span>
            </template>
            <el-menu-item index="/borrow/create">登记领用</el-menu-item>
            <el-menu-item index="/borrow/history">领用记录</el-menu-item>
            <el-menu-item index="/borrow/old-inventory">旧件仓库</el-menu-item>
            <el-menu-item index="/borrow/scrap-inventory">废品仓库</el-menu-item>
          </el-sub-menu>

          
          <el-sub-menu index="reports">
            <template #title>
              <el-icon><DataAnalysis /></el-icon>
              <span>报表中心</span>
            </template>
            <el-menu-item index="/reports/query">维度查询</el-menu-item>
            <el-menu-item index="/reports/detailed">明细查询</el-menu-item>
            <el-menu-item index="/reports/employee-annual">员工年度报表</el-menu-item>
            <el-menu-item index="/reports/scrap-analytics">废品损毁分析</el-menu-item>
          </el-sub-menu>
          
          <el-sub-menu index="system" v-if="userStore.userInfo.role === 'SUPER_ADMIN'">
            <template #title>
              <el-icon><Setting /></el-icon>
              <span>系统管理</span>
            </template>
            <el-menu-item index="/system/users">用户管理</el-menu-item>
            <el-menu-item index="/system/employees">员工管理</el-menu-item>
            <el-menu-item index="/system/brand-config">品牌与公告</el-menu-item>
            <el-menu-item index="/system/backup">数据库备份</el-menu-item>
            <el-menu-item index="/system/notification-config">通知配置</el-menu-item>
            <el-menu-item index="/system/notification-center">通知中心</el-menu-item>
            <el-menu-item index="/system/config">高级参数</el-menu-item>
            <el-menu-item index="/system/operation-log">操作日志</el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-aside>
      
      <el-container>
        <el-header class="glass-effect">
          <div class="header-content">
            <!-- 移动端汉堡菜单按钮 -->
            <el-icon v-if="isMobile" class="hamburger-menu" @click="drawerVisible = true" :size="24">
              <MenuIcon />
            </el-icon>
            
            <h3 class="system-title">{{ isMobile ? (brandConfig.system_name || '仓储管理') : (brandConfig.company_name || brandConfig.system_name || '机械加工配件仓储管理系统') }}</h3>
            <div class="header-right">
              <!-- Message Bell -->
              <div class="message-bell" @click="router.push('/system/messages')">
                <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="bell-badge">
                  <el-icon :size="20"><ChatDotRound /></el-icon>
                </el-badge>
              </div>

              <el-divider direction="vertical" />

              <div class="user-info">
                <span class="username" v-if="!isMobile">{{ userStore.userInfo.realName }}</span>
                <el-dropdown @command="handleCommand">
                  <span class="el-dropdown-link">
                    <el-icon><User /></el-icon>
                  </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="profile">个人设置</el-dropdown-item>
                      <el-dropdown-item divided disabled style="opacity: 1; cursor: default; font-weight: 600; font-size: 12px; color: #86868b;">
                        布局模式
                      </el-dropdown-item>
                      <el-dropdown-item command="layout:auto">
                        <el-icon v-if="forceMode === 'auto'"><Select /></el-icon>
                        <span :style="{ marginLeft: forceMode === 'auto' ? '0' : '20px' }">自动适应</span>
                      </el-dropdown-item>
                      <el-dropdown-item command="layout:mobile">
                        <el-icon v-if="forceMode === 'mobile'"><Select /></el-icon>
                        <span :style="{ marginLeft: forceMode === 'mobile' ? '0' : '20px' }">移动端模式</span>
                      </el-dropdown-item>
                      <el-dropdown-item command="layout:desktop">
                        <el-icon v-if="forceMode === 'desktop'"><Select /></el-icon>
                        <span :style="{ marginLeft: forceMode === 'desktop' ? '0' : '20px' }">桌面端模式</span>
                      </el-dropdown-item>
                      <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
          </div>
        </el-header>
        
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>

    <!-- Popup Announcement -->
    <el-dialog
      v-model="popupVisible"
      :title="popupAnnouncement?.title"
      width="450px"
      class="apple-dialog popup-announcement"
      center
    >
      <div class="popup-content">{{ popupAnnouncement?.content }}</div>
      <template #footer>
        <el-button type="primary" round @click="popupVisible = false">我知道了</el-button>
      </template>
    </el-dialog>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Bell, User, DataLine, Box, DocumentCopy, DataAnalysis, Setting, ChatDotRound, Monitor, Menu as MenuIcon, Select } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import { useResponsive } from '@/composables/useResponsive'
import { logout } from '@/api/auth'
import { getSystemConfig, getAnnouncements } from '@/api/system'
import { getUnreadCount } from '@/api/message'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const { isMobile, isTablet, forceMode, toggleLayoutMode, getDialogWidth } = useResponsive()

const activeMenu = computed(() => route.path)
const brandConfig = ref({})
const announcements = ref([])
const unreadCount = ref(0)
const popupVisible = ref(false)
const drawerVisible = ref(false)

const scrollingAnnouncement = computed(() => 
  announcements.value.find(a => a.type === 'SCROLL' && a.status === 1)
)
const popupAnnouncement = computed(() => 
  announcements.value.find(a => a.type === 'POPUP' && a.status === 1)
)

const loadInitialData = async () => {
  try {
    const [configRes, annoRes, unreadRes] = await Promise.all([
      getSystemConfig(),
      getAnnouncements({ status: 1 }),
      getUnreadCount()
    ])
    brandConfig.value = configRes.data
    announcements.value = annoRes.data || []
    unreadCount.value = unreadRes.data || 0

    // Check if we should show popup
    if (popupAnnouncement.value) {
      const lastSeenId = localStorage.getItem('last_popup_id')
      if (lastSeenId !== String(popupAnnouncement.value.id)) {
        popupVisible.value = true
        localStorage.setItem('last_popup_id', popupAnnouncement.value.id)
      }
    }
  } catch (error) {
    console.error('Layout data load failed', error)
  }
}

// Polling for new messages
let timer = null
onMounted(() => {
  loadInitialData()
  timer = setInterval(async () => {
    const res = await getUnreadCount()
    unreadCount.value = res.data || 0
  }, 30000)
})

onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})

// 在新标签页打开可视化大屏
const openDatavInNewWindow = () => {
  window.open('/datav', '_blank')
}

const handleCommand = async (command) => {
  if (command === 'logout') {
    try {
      await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      
      await logout()
      userStore.logout()
      ElMessage.success('退出登录成功')
      router.push('/login')
    } catch (error) {
    }
  } else if (command === 'profile') {
    router.push('/system/profile')
  } else if (command.startsWith('layout:')) {
    const mode = command.split(':')[1]
    toggleLayoutMode(mode)
    ElMessage.success(`已切换到${mode === 'auto' ? '自动适应' : mode === 'mobile' ? '移动端' : '桌面端'}模式`)
    // 刷新页面以应用新布局
    setTimeout(() => window.location.reload(), 300)
  }
}
</script>

<style scoped>
.layout-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.announcement-bar {
  height: 32px;
  background: rgba(244, 63, 94, 0.08);
  color: #f43f5e;
  font-size: 13px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  border-bottom: 1px solid rgba(244, 63, 94, 0.1);
  z-index: 2000;
}

.announcement-content {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}

.scroll-container {
  flex: 1;
  overflow: hidden;
  white-space: nowrap;
}

.scroll-text {
  display: inline-block;
  padding-left: 100%;
  animation: scroll-left 20s linear infinite;
  font-weight: 500;
}

@keyframes scroll-left {
  0% { transform: translateX(0); }
  100% { transform: translateX(-100%); }
}


.main-layout {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.el-aside {
  background-color: white;
  height: 100%;
  border-right: 1px solid rgba(0,0,0,0.05);
}

.logo {
  height: 70px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  gap: 10px;
}

.logo-icon {
  font-size: 24px;
}

.logo-text {
  font-size: 16px;
  font-weight: 800;
  color: #1d1d1f;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.el-header {
  height: 64px;
  background-color: rgba(255, 255, 255, 0.82);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(0,0,0,0.05);
}

.header-content {
  height: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.system-title {
  font-size: 17px;
  font-weight: 600;
  color: #1d1d1f;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.message-bell {
  cursor: pointer;
  color: #86868b;
  transition: color 0.2s;
  display: flex;
  padding: 5px;
}

.message-bell:hover {
  color: var(--apple-accent);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username {
  font-size: 14px;
  font-weight: 500;
  color: #4b4b4d;
}

.el-dropdown-link {
  cursor: pointer;
  font-size: 20px;
  color: #86868b;
}

.popup-content {
  padding: 10px 0;
  line-height: 1.6;
  color: #1d1d1f;
  font-size: 15px;
}

.el-main {
  background-color: #f5f5f7;
  padding: 24px;
  overflow-y: auto;
}

/* ========== 移动端样式 ========== */
.hamburger-menu {
  cursor: pointer;
  color: #1d1d1f;
  margin-right: 12px;
  transition: transform 0.2s;
}

.hamburger-menu:active {
  transform: scale(0.9);
}

.mobile-drawer :deep(.el-drawer__body) {
  padding: 0;
}

.drawer-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.drawer-content .logo {
  border-bottom: 1px solid rgba(0,0,0,0.05);
}

.drawer-content .apple-menu {
  flex: 1;
  border-right: none;
}

/* 移动端响应式 */
@media (max-width: 768px) {
  .announcement-bar {
    font-size: 12px;
    padding: 0 12px;
    height: 28px;
  }
  
  .el-header {
    height: 56px;
  }
  
  .system-title {
    font-size: 15px;
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .header-right {
    gap: 8px;
  }
  
  .message-bell {
    padding: 4px;
  }
  
  .el-main {
    padding: 16px;
  }
  
  .popup-announcement :deep(.el-dialog) {
    width: 90%;
  }
}

/* 平板优化 */
@media (min-width: 768px) and (max-width: 1024px) {
  .el-aside {
    width: 180px !important;
  }
  
  .logo-text {
    font-size: 14px;
  }
  
  .system-title {
    font-size: 16px;
  }
}
</style>
