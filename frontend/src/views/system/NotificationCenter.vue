<template>
  <div class="notification-center">
    <el-card class="header-card">
      <template #header>
        <div class="card-header">
          <span><el-icon><MessageBox /></el-icon> 通知中心</span>
          <div class="actions">
            <el-button type="success" size="small" @click="handleRunDailyReport">
              <el-icon><DataAnalysis /></el-icon> 立即生成日报
            </el-button>
            <el-button type="warning" size="small" @click="handleRunOverdueCheck">
              <el-icon><WarnTriangleFilled /></el-icon> 立即检查超期
            </el-button>
            <el-button type="primary" size="small" @click="showSendDialog = true">
              <el-icon><Position /></el-icon> 手动发送通知
            </el-button>
          </div>
        </div>
      </template>

      <!-- 统计卡片 -->
      <el-row :gutter="20" style="margin-bottom: 20px">
        <el-col :span="6">
          <el-statistic title="总发送数" :value="stats.total">
            <template #prefix>
              <el-icon><DataLine /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="6">
          <el-statistic title="成功" :value="stats.success">
            <template #prefix>
              <el-icon color="green"><CircleCheck /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="6">
          <el-statistic title="失败" :value="stats.failed">
            <template #prefix>
              <el-icon color="red"><CircleClose /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="6">
          <el-statistic title="队列中" :value="stats.queue_size">
            <template #prefix>
              <el-icon color="orange"><Clock /></el-icon>
            </template>
          </el-statistic>
        </el-col>
      </el-row>

      <!-- 筛选 -->
      <el-form :inline="true" :model="queryForm">
        <el-form-item label="通知渠道">
          <el-select v-model="queryForm.provider_type" placeholder="全部" clearable style="width: 150px">
            <el-option label="钉钉" value="dingtalk" />
            <el-option label="微信" value="wechat" />
            <el-option label="站内信" value="internal" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="待发送" value="pending" />
            <el-option label="成功" value="success" />
            <el-option label="失败" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadLogs">
            <el-icon><Search /></el-icon> 查询
          </el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 通知日志列表 -->
      <el-table :data="logs" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="渠道" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.provider_type === 'dingtalk'" type="primary" size="small">钉钉</el-tag>
            <el-tag v-else-if="row.provider_type === 'wechat'" type="success" size="small">微信</el-tag>
            <el-tag v-else type="info" size="small">站内信</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scene_type" label="场景" width="140">
          <template #default="{ row }">
            {{ getSceneName(row.scene_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" width="200" show-overflow-tooltip />
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'success'" type="success" size="small">成功</el-tag>
            <el-tag v-else-if="row.status === 'failed'" type="danger" size="small">失败</el-tag>
            <el-tag v-else type="warning" size="small">待发送</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="重试" width="80">
          <template #default="{ row }">
            {{ row.retry_count }}次
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button 
              v-if="row.status === 'failed'" 
              type="text" 
              size="small" 
              @click="viewError(row)"
            >
              查看错误
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="queryForm.page"
        v-model:page-size="queryForm.page_size"
        :total="total"
        :page-sizes="[20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadLogs"
        @size-change="loadLogs"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <!-- 手动发送对话框 -->
    <el-dialog v-model="showSendDialog" title="手动发送通知" width="600px">
      <el-form :model="sendForm" label-width="100px">
        <el-form-item label="通知标题" required>
          <el-input v-model="sendForm.title" placeholder="通知标题" />
        </el-form-item>
        <el-form-item label="通知内容" required>
          <el-input 
            v-model="sendForm.content" 
            type="textarea" 
            :rows="5"
            placeholder="通知内容"
          />
        </el-form-item>
        <el-form-item label="发送渠道" required>
          <el-checkbox-group v-model="sendForm.providers">
            <el-checkbox label="dingtalk">钉钉</el-checkbox>
            <el-checkbox label="wechat">微信</el-checkbox>
            <el-checkbox label="internal">站内信</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showSendDialog = false">取消</el-button>
        <el-button type="primary" @click="sendNotification" :loading="sending">发送</el-button>
      </template>
    </el-dialog>

    <!-- 错误详情对话框 -->
    <el-dialog v-model="showErrorDialog" title="错误详情" width="500px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="通知ID">{{ currentError.id }}</el-descriptions-item>
        <el-descriptions-item label="标题">{{ currentError.title }}</el-descriptions-item>
        <el-descriptions-item label="重试次数">{{ currentError.retry_count }}次</el-descriptions-item>
        <el-descriptions-item label="错误信息">
          <div style="color: red">{{ currentError.error_msg }}</div>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  MessageBox, Position, DataLine, CircleCheck, CircleClose, 
  Clock, Search, DataAnalysis, WarnTriangleFilled 
} from '@element-plus/icons-vue'
import { 
  getNotificationLogs, 
  getNotificationStats,
  sendNotification as sendNotificationApi,
  runDailyReport,
  runOverdueCheck
} from '@/api/notification'

const loading = ref(false)
const sending = ref(false)
const showSendDialog = ref(false)
const showErrorDialog = ref(false)

const logs = ref([])
const total = ref(0)
const stats = ref({
  total: 0,
  success: 0,
  failed: 0,
  pending: 0,
  queue_size: 0
})

const queryForm = ref({
  page: 1,
  page_size: 20,
  provider_type: '',
  status: ''
})

const sendForm = ref({
  title: '',
  content: '',
  providers: ['dingtalk', 'internal']
})

const currentError = ref({})

const sceneNames = {
  'stock_warning': '库存预警',
  'borrow_created': '领用通知',
  'return_reminder': '归还提醒',
  'daily_report': '每日报表',
  'weekly_report': '周报推送',
  'system_announcement': '系统公告'
}

const getSceneName = (scene) => {
  return sceneNames[scene] || scene
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

// 加载统计信息
const loadStats = async () => {
  try {
    const res = await getNotificationStats()
    if (res.code === 200) {
      stats.value = res.data
    }
  } catch (error) {
    console.error('加载统计失败:', error)
  }
}

// 加载日志列表
const loadLogs = async () => {
  loading.value = true
  try {
    const res = await getNotificationLogs(queryForm.value)
    if (res.code === 200) {
      logs.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

// 重置查询
const resetQuery = () => {
  queryForm.value = {
    page: 1,
    page_size: 20,
    provider_type: '',
    status: ''
  }
  loadLogs()
}

// 发送通知
const sendNotification = async () => {
  if (!sendForm.value.title || !sendForm.value.content) {
    ElMessage.warning('请填写标题和内容')
    return
  }
  if (sendForm.value.providers.length === 0) {
    ElMessage.warning('请选择至少一个发送渠道')
    return
  }

  sending.value = true
  try {
    const res = await sendNotificationApi(sendForm.value)
    if (res.code === 200) {
      ElMessage.success('通知已加入发送队列')
      showSendDialog.value = false
      sendForm.value = {
        title: '',
        content: '',
        providers: ['dingtalk', 'internal']
      }
      // 刷新列表
      setTimeout(() => {
        loadLogs()
        loadStats()
      }, 1000)
    } else {
      ElMessage.error(res.message || '发送失败')
    }
  } catch (error) {
    ElMessage.error('发送失败: ' + error.message)
  } finally {
    sending.value = false
  }
}

// 立即执行日报
const handleRunDailyReport = async () => {
  try {
    const res = await runDailyReport()
    if (res.code === 200) {
      ElMessage.success('报表生成任务已启动')
      setTimeout(() => {
        loadLogs()
        loadStats()
      }, 1500)
    }
  } catch (error) {
    ElMessage.error('触发失败: ' + error.message)
  }
}

// 立即执行超期检查
const handleRunOverdueCheck = async () => {
  try {
    const res = await runOverdueCheck()
    if (res.code === 200) {
      ElMessage.success('超期检查任务已启动')
      setTimeout(() => {
        loadLogs()
        loadStats()
      }, 1500)
    }
  } catch (error) {
    ElMessage.error('触发失败: ' + error.message)
  }
}

// 查看错误
const viewError = (row) => {
  currentError.value = row
  showErrorDialog.value = true
}

onMounted(() => {
  loadStats()
  loadLogs()
  
  // 定时刷新统计信息
  setInterval(() => {
    loadStats()
  }, 10000)
})
</script>

<style scoped>
.notification-center {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}

.card-header span {
  display: flex;
  align-items: center;
  gap: 8px;
}

.actions {
  display: flex;
  gap: 10px;
}

:deep(.el-statistic__head) {
  font-size: 14px;
  color: #909399;
}

:deep(.el-statistic__content) {
  font-size: 24px;
  font-weight: 600;
}
</style>
