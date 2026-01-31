<template>
  <div class="notification-settings">
    <el-card>
      <template #header>
        <div class="card-header">
          <span><el-icon><Bell /></el-icon> 通知偏好设置</span>
        </div>
      </template>

      <el-alert 
        title="功能说明" 
        type="info" 
        :closable="false"
        style="margin-bottom: 20px">
        <p>您可以为每个通知场景独立选择接收渠道。勾选的渠道将接收相应的通知。</p>
        <p><strong>钉钉</strong>: 发送至配置的钉钉群 | <strong>微信</strong>: 发送至绑定的微信 | <strong>站内信</strong>: 系统内部消息</p>
      </el-alert>

      <el-table :data="settings" border stripe v-loading="loading">
        <el-table-column prop="scene_name" label="通知场景" width="180" />
        <el-table-column prop="scene_description" label="说明" />
        <el-table-column label="接收渠道" width="350">
          <template #default="{ row }">
            <el-checkbox-group v-model="row.channels" @change="handleChannelChange(row)">
              <el-checkbox label="dingtalk">
                <el-icon color="#2ba1ff"><ChatDotSquare /></el-icon> 钉钉
              </el-checkbox>
              <el-checkbox label="wechat">
                <el-icon color="#07c160"><ChatLineSquare /></el-icon> 微信
              </el-checkbox>
              <el-checkbox label="internal">
                <el-icon color="#909399"><Message /></el-icon> 站内信
              </el-checkbox>
            </el-checkbox-group>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.channels.length === 0" type="info" size="small">未订阅</el-tag>
            <el-tag v-else type="success" size="small">已订阅</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Bell, ChatDotSquare, ChatLineSquare, Message } from '@element-plus/icons-vue'
import { getUserNotificationSettings, updateUserNotificationSetting } from '@/api/notification'

const loading = ref(false)

// 通知场景定义
const sceneDefinitions = [
  { 
    scene_type: 'stock_warning', 
    scene_name: '库存预警', 
    scene_description: '配件库存低于安全阈值时通知' 
  },
  { 
    scene_type: 'borrow_created', 
    scene_name: '领用通知', 
    scene_description: '有新的领用记录创建时通知' 
  },
  { 
    scene_type: 'return_reminder', 
    scene_name: '归还提醒', 
    scene_description: '未归还配件超过7天时提醒' 
  },
  { 
    scene_type: 'restock', 
    scene_name: '入库通知', 
    scene_description: '配件补货入库时通知' 
  },
  { 
    scene_type: 'daily_report', 
    scene_name: '每日报表', 
    scene_description: '每日自动生成并推送统计报表' 
  },
  { 
    scene_type: 'weekly_report', 
    scene_name: '周报推送', 
    scene_description: '每周汇总数据并推送' 
  },
  {
    scene_type: 'monthly_report',
    scene_name: '月报推送',
    scene_description: '每月数据统计与分析'
  },
  {
    scene_type: 'system_announcement',
    scene_name: '系统公告',
    scene_description: '系统发布重要公告时通知'
  }
]

const settings = ref([])

// 加载用户设置
const loadSettings = async () => {
  loading.value = true
  try {
    const res = await getUserNotificationSettings()
    if (res.code === 200) {
      // 创建场景设置映射
      const settingsMap = {}
      res.data.forEach(setting => {
        settingsMap[setting.scene_type] = setting
      })

      // 初始化所有场景的设置
      settings.value = sceneDefinitions.map(scene => {
        const userSetting = settingsMap[scene.scene_type]
        const channels = []
        
        if (userSetting) {
          if (userSetting.enable_dingtalk) channels.push('dingtalk')
          if (userSetting.enable_wechat) channels.push('wechat')
          if (userSetting.enable_internal) channels.push('internal')
        } else {
          // 默认启用站内信
          channels.push('internal')
        }

        return {
          ...scene,
          channels,
          id: userSetting?.id
        }
      })
    }
  } catch (error) {
    ElMessage.error('加载设置失败')
  } finally {
    loading.value = false
  }
}

// 处理渠道变更
const handleChannelChange = async (row) => {
  try {
    const res = await updateUserNotificationSetting({
      scene_type: row.scene_type,
      enable_dingtalk: row.channels.includes('dingtalk'),
      enable_wechat: row.channels.includes('wechat'),
      enable_internal: row.channels.includes('internal')
    })

    if (res.code === 200) {
      ElMessage.success('保存成功')
    } else {
      ElMessage.error(res.message || '保存失败')
      // 失败时重新加载
      loadSettings()
    }
  } catch (error) {
    ElMessage.error('保存失败')
    loadSettings()
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.notification-settings {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
}

:deep(.el-checkbox) {
  margin-right: 20px;
}

:deep(.el-checkbox__label) {
  display: flex;
  align-items: center;
  gap: 4px;
}

:deep(.el-alert p) {
  margin: 5px 0;
}
</style>
