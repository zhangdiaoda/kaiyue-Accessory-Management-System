<template>
  <div class="backup-container">
    <div class="page-header">
      <h2 class="page-title">数据库备份与恢复</h2>
      <div class="header-actions">
        <!-- 立即备份按钮 -->
        <el-button type="primary" :loading="backuping" @click="handleRunBackup" class="apple-btn">
          <el-icon class="el-icon--left"><VideoPlay /></el-icon> 立即全量备份
        </el-button>
      </div>
    </div>

    <!-- 配置卡片 -->
    <el-card class="apple-card config-card">
      <template #header>
        <div class="card-header">
          <span>备份策略配置</span>
          <el-button type="primary" link @click="handleSaveConfig">保存设置</el-button>
        </div>
      </template>
      
      <el-form label-position="top" class="config-form">
        <el-row :gutter="20">
          <el-col :span="12">
             <el-form-item label="备份存储路径 (绝对路径)">
               <el-input 
                 v-model="config.backup_path" 
                 placeholder="例如: D:\Backups"
                 prefix-icon="Folder"
               />
               <div class="form-tip">请确保后端服务 (Go) 对该路径有写入权限</div>
             </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="自动备份周期 (Cron 表达式)">
               <el-input 
                 v-model="config.schedule" 
                 placeholder="例如: 0 0 2 * * * (每天凌晨2点)" 
                 prefix-icon="Timer"
               >
                 <template #append>
                   <el-popover placement="bottom" :width="300" trigger="click">
                     <template #reference>
                       <el-button>常用</el-button>
                     </template>
                     <div class="cron-presets">
                       <el-button size="small" @click="config.schedule = '0 0 2 * * *'">每天凌晨2点</el-button>
                       <el-button size="small" @click="config.schedule = '0 0 */6 * * *'">每6小时</el-button>
                       <el-button size="small" @click="config.schedule = '0 0 0 * * 1'">每周一凌晨</el-button>
                       <el-button size="small" @click="config.schedule = '0 */30 * * * *'">每30分钟 (测试)</el-button>
                     </div>
                   </el-popover>
                 </template>
               </el-input>
               <div class="form-tip">
                 支持标准的 Cron 表达式与秒级扩展 
                 <span v-if="config.next_run" class="next-run-tip">
                   (下次自动执行: {{ config.next_run }})
                 </span>
               </div>
             </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 列表卡片 -->
    <el-card class="apple-card list-card">
      <template #header>
        <div class="card-header">
          <span>历史备份文件</span>
          <el-button circle icon="Refresh" @click="loadList"></el-button>
        </div>
      </template>

      <el-table :data="fileList" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="文件名" min-width="200">
           <template #default="{ row }">
             <div class="file-name">
               <el-icon><Document /></el-icon>
               <span>{{ row.name }}</span>
               <el-tag v-if="isRecent(row.time)" size="small" type="success" effect="plain" class="ml-2">New</el-tag>
             </div>
           </template>
        </el-table-column>
        <el-table-column prop="time" label="备份时间" width="180" />
        <el-table-column prop="size" label="文件大小" width="120">
          <template #default="{ row }">
            {{ formatSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column prop="path" label="存储路径" min-width="250" show-overflow-tooltip />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button size="small" icon="Download" @click="handleDownload(row)">下载</el-button>
              <el-button size="small" type="warning" icon="RefreshLeft" @click="handleRestore(row)">恢复</el-button>
              <el-button size="small" type="danger" icon="Delete" @click="handleDelete(row)">删除</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { VideoPlay, Folder, Timer, Document, Refresh, Download, RefreshLeft, Delete } from '@element-plus/icons-vue'
import { getBackupConfig, updateBackupConfig, runBackup, getBackupList, deleteBackup, restoreBackup, getDownloadUrl } from '@/api/backup'
import dayjs from 'dayjs'

const config = ref({
  backup_path: '',
  schedule: ''
})

const fileList = ref([])
const loading = ref(false)
const backuping = ref(false)

const loadConfig = async () => {
  try {
    const res = await getBackupConfig()
    config.value = res.data
  } catch (err) {
    console.error(err)
  }
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await getBackupList()
    fileList.value = res.data || []
  } finally {
    loading.value = false
  }
}

const handleSaveConfig = async () => {
  try {
    await updateBackupConfig(config.value)
    ElMessage.success('备份策略已更新')
    loadList() // 可能会因为路径变更导致列表变化
  } catch (err) {
    // error handled by request interceptor
  }
}

const handleRunBackup = async () => {
  try {
    backuping.value = true
    await runBackup()
    ElMessage.success('备份任务已后台启动，请稍后刷新列表')
    // 延迟 2s 刷新列表
    setTimeout(() => {
      loadList()
      backuping.value = false
    }, 2000)
  } catch (err) {
    backuping.value = false
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要永久删除该备份文件吗？', '警告', {
      type: 'warning',
      confirmButtonText: '删除',
      confirmButtonClass: 'el-button--danger'
    })
    await deleteBackup(row.name)
    ElMessage.success('删除成功')
    loadList()
  } catch (err) {
    if (err !== 'cancel') console.error(err)
  }
}

const handleRestore = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要将数据库恢复到 [${row.time}] 的状态吗？\n警告：当前所有新数据将被覆盖！此操作不可逆！`, 
      '高风险操作', 
      {
        type: 'error',
        confirmButtonText: '确定恢复',
        confirmButtonClass: 'el-button--danger',
        cancelButtonText: '取消'
      }
    )
    
    // 二次确认
    await ElMessageBox.prompt('请输入 CONFIRM 以确认恢复操作', '最终确认', {
      confirmButtonText: '执行',
      inputPattern: /^CONFIRM$/,
      inputErrorMessage: '输入不正确'
    })

    const loadingInstance = ElMessageBox.confirm('恢复中，请勿关闭页面...', '执行中', { showConfirmButton: false, showCancelButton: false })
    
    try {
      await restoreBackup(row.name)
      ElMessageBox.close()
      ElMessageBox.alert('数据库已成功恢复，系统将重新加载', '恢复成功', {
        callback: () => window.location.reload()
      })
    } catch (err) {
      ElMessageBox.close()
    }
  } catch (err) {
    if (err !== 'cancel') console.error(err)
  }
}

const handleDownload = (row) => {
  const url = getDownloadUrl(row.name)
  window.open(url, '_blank')
}

const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const isRecent = (timeStr) => {
  return dayjs().diff(dayjs(timeStr), 'hour') < 24
}

onMounted(() => {
  loadConfig()
  loadList()
})
</script>

<style scoped>
.backup-container {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1d1d1f;
  margin: 0;
}

.apple-card {
  border-radius: 12px;
  border: 1px solid rgba(0,0,0,0.05);
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
  margin-bottom: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
  color: #1d1d1f;
}

.form-tip {
  font-size: 12px;
  color: #86868b;
  margin-top: 6px;
  line-height: 1.4;
}

.cron-presets {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.next-run-tip {
  color: #34c759;
  font-weight: 500;
  margin-left: 5px;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.ml-2 { margin-left: 8px; }

</style>
