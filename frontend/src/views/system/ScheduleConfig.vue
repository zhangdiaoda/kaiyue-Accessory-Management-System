<template>
  <div class="schedule-config">
    <el-card>
      <template #header>
        <div class="card-header">
          <span><el-icon><Timer /></el-icon> 定时任务配置</span>
        </div>
      </template>

      <el-alert 
        title="Cron 表达式说明" 
        type="info" 
        :closable="false"
        style="margin-bottom: 20px">
        <p>格式：<code>秒 分 时 日 月 周</code></p>
        <p>示例：<code>0 30 8 * * *</code> 表示每天早上 8:30:00</p>
        <p>在线工具：<a href="https://crontab.guru/" target="_blank">https://crontab.guru/</a></p>
      </el-alert>

      <el-table :data="schedules" stripe v-loading="loading">
        <el-table-column prop="config_key" label="任务" width="250">
          <template #default="{ row }">
            <strong>{{ getTaskName(row.config_key) }}</strong>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="说明" />
        <el-table-column prop="config_value" label="Cron 表达式" width="200">
          <template #default="{ row }">
            <el-tag type="primary">{{ row.config_value }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="中文说明" width="250">
          <template #default="{ row }">
            {{ parseCron(row.config_value) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">
              <el-icon><Edit /></el-icon> 编辑
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog v-model="showDialog" title="编辑定时任务" width="600px">
      <el-form label-width="120px">
        <el-form-item label="任务名称">
          <el-input :value="getTaskName(currentConfig.config_key)" disabled />
        </el-form-item>
        <el-form-item label="任务说明">
          <el-input :value="currentConfig.description" disabled />
        </el-form-item>
        <el-form-item label="Cron 表达式" required>
          <el-input 
            v-model="currentConfig.config_value" 
            placeholder="0 30 8 * * *"
            @change="validateCron"
          />
          <div class="form-tip">
            {{ cronDescription }}
          </div>
        </el-form-item>
        <el-form-item label="快捷设置">
          <el-select v-model="quickSet" placeholder="选择快捷时间" @change="applyCronSet" style="width: 100%">
            <el-option label="每天 早上 6:00" value="0 0 6 * * *" />
            <el-option label="每天 早上 7:00" value="0 0 7 * * *" />
            <el-option label="每天 早上 8:00" value="0 0 8 * * *" />
            <el-option label="每天 早上 8:30" value="0 30 8 * * *" />
            <el-option label="每天 早上 9:00" value="0 0 9 * * *" />
            <el-option label="每天 中午 12:00" value="0 0 12 * * *" />
            <el-option label="每天 下午 18:00" value="0 0 18 * * *" />
            <el-option label="每周一 早上 9:00" value="0 0 9 * * 1" />
            <el-option label="每月1号 早上 9:00" value="0 0 9 1 * *" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Timer, Edit } from '@element-plus/icons-vue'
import { getScheduleConfigs, updateScheduleConfig } from '@/api/notification'

const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const schedules = ref([])
const quickSet = ref('')
const cronDescription = ref('')

const currentConfig = reactive({
  config_key: '',
  config_value: '',
  description: ''
})

const taskNames = {
  'schedule_daily_report_cron': '每日报表推送',
  'schedule_overdue_check_cron': '超期未归还检查',
  'schedule_weekly_report_cron': '周报推送',
  'schedule_monthly_report_cron': '月报推送'
}

const getTaskName = (key) => {
  return taskNames[key] || key
}

const parseCron = (cron) => {
  if (!cron) return ''
  const parts = cron.split(' ')
  if (parts.length < 6) return '格式错误'

  const [sec, min, hour, day, month, week] = parts
  let desc = ''

  if (week !== '*') {
    const weekNames = ['日', '一', '二', '三', '四', '五', '六']
    desc = `每周${weekNames[week]} `
  } else if (day !== '*') {
    desc = `每月${day}号 `
  } else {
    desc = '每天 '
  }

  desc += `${hour.padStart(2, '0')}:${min.padStart(2, '0')}`
  return desc
}

const loadConfigs = async () => {
  loading.value = true
  try {
    const res = await getScheduleConfigs()
    if (res.code === 200) {
      schedules.value = res.data || []
    }
  } catch (error) {
    ElMessage.error('加载配置失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row) => {
  Object.assign(currentConfig, row)
  cronDescription.value = parseCron(row.config_value)
  quickSet.value = ''
  showDialog.value = true
}

const applyCronSet = (value) => {
  currentConfig.config_value = value
  cronDescription.value = parseCron(value)
}

const validateCron = () => {
  cronDescription.value = parseCron(currentConfig.config_value)
}

const handleSubmit = async () => {
  if (!currentConfig.config_value) {
    ElMessage.warning('请输入 Cron 表达式')
    return
  }

  // 简单验证
  const parts = currentConfig.config_value.split(' ')
  if (parts.length !== 6) {
    ElMessage.error('Cron 表达式格式错误，应为 6 个字段（秒 分 时 日 月 周）')
    return
  }

  saving.value = true
  try {
    const res = await updateScheduleConfig({
      config_key: currentConfig.config_key,
      config_value: currentConfig.config_value
    })

    if (res.code === 200) {
      ElMessage.success('配置已更新并立即生效')
      showDialog.value = false
      loadConfigs()
    } else {
      ElMessage.error(res.message || '保存失败')
    }
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadConfigs()
})
</script>

<style scoped>
.schedule-config {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
}

.form-tip {
  font-size: 12px;
  color: #409eff;
  margin-top: 5px;
  font-weight: 500;
}

:deep(.el-alert p) {
  margin: 5px 0;
}

:deep(.el-alert code) {
  background: #f2f6fc;
  padding: 2px 6px;
  border-radius: 3px;
  color: #409eff;
}

:deep(.el-alert a) {
  color: #409eff;
  text-decoration: none;
}

:deep(.el-alert a:hover) {
  text-decoration: underline;
}
</style>
