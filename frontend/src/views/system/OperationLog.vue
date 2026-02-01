<template>
  <div class="operation-log">
    <el-card>
      <template #header>
        <div class="card-header">
          <div style="display: flex; align-items: center; gap: 10px;">
            <span>📝 操作日志</span>
            <el-tag v-if="logSize !== null" type="info" size="small">
              当前大小: {{ logSize.size_mb }} MB
            </el-tag>
          </div>
          <div>
            <el-button type="warning" size="small" @click="configVisible = true">⚙️ 清理配置</el-button>
            <el-button type="danger" size="small" @click="handleClearLogs">🗑️ 清空日志</el-button>
          </div>
        </div>
      </template>

      <!-- 筛选表单 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="输入用户名" clearable />
        </el-form-item>
        <el-form-item label="模块">
          <el-select v-model="searchForm.module" placeholder="选择模块" clearable>
            <el-option label="配件管理" value="配件管理" />
            <el-option label="领用管理" value="领用管理" />
            <el-option label="员工管理" value="员工管理" />
            <el-option label="报表管理" value="报表管理" />
            <el-option label="系统管理" value="系统管理" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作类型">
          <el-select v-model="searchForm.operation" placeholder="选择操作" clearable>
            <el-option label="创建" value="创建" />
            <el-option label="更新" value="更新" />
            <el-option label="删除" value="删除" />
            <el-option label="导入" value="导入" />
            <el-option label="导出" value="导出" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="选择状态" clearable>
            <el-option label="成功" value="SUCCESS" />
            <el-option label="失败" value="FAILED" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DDTHH:mm:ss+08:00"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 日志表格 -->
      <el-table :data="tableData" border stripe v-loading="loading">
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column prop="real_name" label="操作人" width="100" />
        <el-table-column prop="module" label="模块" width="100" />
        <el-table-column prop="operation" label="操作" width="80" />
        <el-table-column prop="description" label="描述" min-width="150" />
        <el-table-column prop="ip_address" label="IP地址" width="130" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'SUCCESS' ? 'success' : 'danger'">
              {{ row.status === 'SUCCESS' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="execution_time" label="耗时(ms)" width="100" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleViewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.pageNum"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadData"
        @current-change="loadData"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="操作日志详情" width="800px">
      <el-descriptions :column="2" border v-if="currentLog">
        <el-descriptions-item label="操作人">{{ currentLog.real_name }} ({{ currentLog.username }})</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ currentLog.created_at }}</el-descriptions-item>
        <el-descriptions-item label="模块">{{ currentLog.module }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">{{ currentLog.operation }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentLog.ip_address }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentLog.status === 'SUCCESS' ? 'success' : 'danger'">
            {{ currentLog.status === 'SUCCESS' ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="请求方法">{{ currentLog.request_method }}</el-descriptions-item>
        <el-descriptions-item label="耗时">{{ currentLog.execution_time }}ms</el-descriptions-item>
        <el-descriptions-item label="请求URL" :span="2">{{ currentLog.request_url }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentLog.description }}</el-descriptions-item>
        <el-descriptions-item label="请求参数" :span="2">
          <pre style="max-height: 200px; overflow: auto;">{{ currentLog.request_params || '无' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="User-Agent" :span="2">{{ currentLog.user_agent }}</el-descriptions-item>
        <el-descriptions-item label="错误信息" :span="2" v-if="currentLog.error_message">
          <el-text type="danger">{{ currentLog.error_message }}</el-text>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 清理配置对话框 -->
    <el-dialog v-model="configVisible" title="⚠️ 日志清理配置" width="600px">
      <el-form :model="cleanupConfig" label-width="120px">
        <el-form-item label="启用自动清理">
          <el-switch v-model="cleanupConfig.enabled" />
        </el-form-item>
        <el-form-item label="清理周期">
          <el-select v-model="cleanupConfig.schedule" :disabled="!cleanupConfig.enabled">
            <el-option label="每天" value="daily" />
            <el-option label="每周" value="weekly" />
            <el-option label="每月" value="monthly" />
          </el-select>
        </el-form-item>
        <el-form-item label="大小阈值(MB)">
          <el-input-number 
            v-model="cleanupConfig.size_threshold" 
            :min="100" 
            :max="10000" 
            :step="100"
            :disabled="!cleanupConfig.enabled"
          />
          <div style="font-size: 12px; color: #999; margin-top: 5px;">
            超过此大小将自动清理旧日志
          </div>
        </el-form-item>
        <el-form-item label="保留天数">
          <el-input-number 
            v-model="cleanupConfig.days_to_keep" 
            :min="7" 
            :max="365" 
            :disabled="!cleanupConfig.enabled"
          />
          <div style="font-size: 12px; color: #999; margin-top: 5px;">
            自动删除超过此天数的日志
          </div>
        </el-form-item>
        <el-form-item label="最小保留条数">
          <el-input-number 
            v-model="cleanupConfig.keep_count" 
            :min="1000" 
            :max="100000" 
            :step="1000"
            :disabled="!cleanupConfig.enabled"
          />
          <div style="font-size: 12px; color: #999; margin-top: 5px;">
            按大小清理时至少保留的记录数
          </div>
        </el-form-item>
        <el-form-item label="最后清理时间" v-if="cleanupConfig.last_cleanup_at">
          <el-text>{{ cleanupConfig.last_cleanup_at }}</el-text>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveConfig">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getOperationLogs, getLogDetail, clearLogs, getLogSize, getCleanupConfig, updateCleanupConfig } from '@/api/audit'

const loading = ref(false)
const tableData = ref([])
const detailVisible = ref(false)
const configVisible = ref(false)
const currentLog = ref(null)
const dateRange = ref([])
const logSize = ref(null)

const searchForm = reactive({
  username: '',
  module: '',
  operation: '',
  status: ''
})

const pagination = reactive({
  pageNum: 1,
  pageSize: 20,
  total: 0
})

const cleanupConfig = reactive({
  enabled: false,
  schedule: 'weekly',
  size_threshold: 1024,
  days_to_keep: 30,
  keep_count: 10000,
  last_cleanup_at: ''
})

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      ...searchForm,
      page_num: pagination.pageNum,
      page_size: pagination.pageSize
    }

    // 添加时间范围
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_time = dateRange.value[0]
      params.end_time = dateRange.value[1]
    }

    const res = await getOperationLogs(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    ElMessage.error('加载日志失败')
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  Object.assign(searchForm, {
    username: '',
    module: '',
    operation: '',
    status: ''
  })
  dateRange.value = []
  pagination.pageNum = 1
  loadData()
}

const handleViewDetail = async (row) => {
  try {
    const res = await getLogDetail(row.id)
    if (res.code === 200) {
      currentLog.value = res.data
      detailVisible.value = true
    }
  } catch (error) {
    ElMessage.error('获取详情失败')
  }
}

// 加载日志大小
const loadLogSize = async () => {
  try {
    const res = await getLogSize()
    if (res.code === 200) {
      logSize.value = res.data
    }
  } catch (error) {
    console.error('获取日志大小失败', error)
  }
}

// 加载清理配置
const loadCleanupConfig = async () => {
  try {
    const res = await getCleanupConfig()
    if (res.code === 200) {
      Object.assign(cleanupConfig, res.data)
    }
  } catch (error) {
    console.error('获取清理配置失败', error)
  }
}

// 清空日志
const handleClearLogs = async () => {
  try {
    await ElMessageBox.confirm(
      '此操作将清空所有操作日志，是否继续？',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const res = await clearLogs()
    if (res.code === 200) {
      ElMessage.success(`已清空 ${res.data.deleted_count} 条日志`)
      loadData()
      loadLogSize()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清空失败')
    }
  }
}

// 保存清理配置
const handleSaveConfig = async () => {
  try {
    const res = await updateCleanupConfig(cleanupConfig)
    if (res.code === 200) {
      ElMessage.success('配置已保存')
      configVisible.value = false
    }
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

onMounted(() => {
  loadData()
  loadLogSize()
  loadCleanupConfig()
})
</script>

<style scoped>
.operation-log {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

pre {
  background: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.5;
}
</style>
