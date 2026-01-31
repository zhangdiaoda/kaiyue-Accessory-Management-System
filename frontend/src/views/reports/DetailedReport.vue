<template>
  <div class="detailed-report">
    <el-card class="filter-card" shadow="never">
      <div class="header-content">
        <div class="title-group">
          <h2>📊 领用明细透视</h2>
          <p class="subtitle">按月份查看全局领用明细，支持多维度筛选</p>
        </div>
        <div class="filter-actions">
          <el-form :inline="true" class="search-form">
            <el-form-item label="时间范围">
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                style="width: 260px"
                @change="loadData"
              />
            </el-form-item>
            <el-form-item label="员工">
              <el-select v-model="searchForm.employee_id" placeholder="全部员工" clearable filterable style="width: 180px">
                <el-option
                  v-for="emp in employees"
                  :key="emp.id"
                  :label="`${emp.name} (${emp.employee_no})`"
                  :value="emp.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="配件">
              <el-select v-model="searchForm.part_id" placeholder="全部配件" clearable filterable style="width: 180px">
                <el-option
                  v-for="part in parts"
                  :key="part.id"
                  :label="`${part.name} (${part.part_no})`"
                  :value="part.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadData" :loading="loading">查询</el-button>
            </el-form-item>
          </el-form>
          <div class="divider"></div>
          <el-button type="success" plain @click="handleExportAll" :disabled="tableData.length === 0">
            <el-icon><Download /></el-icon> 导出全量报表
          </el-button>
          <el-button type="primary" plain @click="handlePushReport" :disabled="tableData.length === 0">
            <el-icon><Share /></el-icon> 推送报表至钉钉
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Monthly Card Grid -->
    <div v-loading="loading" class="grid-section">
      <el-row :gutter="24" v-if="uniqueMonths.length > 0">
        <el-col :span="6" v-for="month in uniqueMonths" :key="month" class="grid-col">
          <div class="month-premium-card has-data" @click="handleMonthClick(month)">
            <div class="month-label">{{ formatMonthLabel(month) }}</div>
            <div class="month-status">
              <span class="status-dot online"></span>
              {{ getMonthData(month).length }} 条记录明细
            </div>
            
            <div class="month-stats">
              <div class="stat-item">
                <span class="lab">总领用</span>
                <span class="val">{{ getMonthTotal(month, 'total_borrow') }}</span>
              </div>
              <div class="stat-item">
                <span class="lab">总废品</span>
                <span class="val" :class="{ 'red': getMonthTotal(month, 'total_damaged') > 0 }">
                  {{ getMonthTotal(month, 'total_damaged') }}
                </span>
              </div>
            </div>
            
            <div class="click-hint">点击查看详情 <el-icon><ArrowRight /></el-icon></div>
          </div>
        </el-col>
      </el-row>
      <el-empty v-else description="所选范围内暂无明细数据" />
    </div>

    <!-- Detail Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="`📅 ${currentMonth} 领用详情明细`"
      width="1000px"
      class="custom-dialog"
      append-to-body
    >
      <div class="dialog-header-actions">
        <div class="month-summary-mini">
          <span>本周期汇总: <strong>{{ getMonthTotal(currentMonth, 'total_borrow') }}</strong> 领用</span>
          <span style="margin-left: 20px">损毁: <strong class="text-danger">{{ getMonthTotal(currentMonth, 'total_damaged') }}</strong></span>
        </div>
        <el-button type="primary" size="small" @click="handleExportMonth">
          <el-icon><Download /></el-icon> 导出本月明细
        </el-button>
      </div>

      <el-table :data="currentMonthData" border stripe max-height="500px">
        <el-table-column prop="employee_name" label="员工" width="140">
          <template #default="{ row }">
            {{ row.employee_name }} <span class="sub-text">({{ row.employee_no }})</span>
          </template>
        </el-table-column>
        <el-table-column prop="department" label="部门" width="120" />
        <el-table-column prop="part_name" label="配件" min-width="150">
          <template #default="{ row }">
            {{ row.part_name }} <span class="sub-text">({{ row.part_no }})</span>
          </template>
        </el-table-column>
        <el-table-column prop="borrow_count" label="频次" width="80" align="center" />
        <el-table-column prop="total_borrow" label="数量" width="80" align="center" />
        <el-table-column prop="total_return" label="归还" width="80" align="center" />
        <el-table-column prop="total_damaged" label="废品" width="80" align="center">
          <template #default="{ row }">
            <span :class="{ 'text-danger': row.total_damaged > 0 }">{{ row.total_damaged }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="unreturned" label="待归还" width="80" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.unreturned > 0" type="warning" size="small">{{ row.unreturned }}</el-tag>
            <span v-else>0</span>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDetailedReport, pushReport } from '@/api/report'
import { getAllEmployees } from '@/api/employee'
import { getPartList } from '@/api/part'
import { Download, Share, ArrowRight } from '@element-plus/icons-vue'

const loading = ref(false)
const tableData = ref([])
const employees = ref([])
const parts = ref([])
const dialogVisible = ref(false)
const currentMonth = ref(null)

const dateRange = ref([
  new Date(new Date().setMonth(new Date().getMonth() - 5)),
  new Date()
])

const searchForm = reactive({
  employee_id: '',
  part_id: ''
})

const uniqueMonths = computed(() => {
  const months = tableData.value.map(item => item.month)
  return [...new Set(months)].sort().reverse()
})

const currentMonthData = computed(() => {
  if (!currentMonth.value) return []
  return tableData.value.filter(item => item.month === currentMonth.value)
})

const formatMonthLabel = (m) => {
  const [y, mm] = m.split('-')
  return `${y}年${parseInt(mm)}月`
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      start_date: formatDate(dateRange.value[0]),
      end_date: formatDate(dateRange.value[1]),
      employee_id: searchForm.employee_id,
      part_id: searchForm.part_id
    }

    const res = await getDetailedReport(params)
    tableData.value = res.data || []
  } catch (error) {
    ElMessage.error('加载详情报表失败')
  } finally {
    loading.value = false
  }
}

const loadEmployees = async () => {
  const res = await getAllEmployees()
  employees.value = res.data || []
}

const loadParts = async () => {
  const res = await getPartList({ page: 1, pageSize: 1000 })
  parts.value = res.data.records || []
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const handleMonthClick = (month) => {
  currentMonth.value = month
  dialogVisible.value = true
}

const exportToCsv = (filename, headers, data) => {
  const csvContent = [
    headers.join(','),
    ...data.map(row => row.join(','))
  ].join('\n')

  const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', filename)
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  ElMessage.success('报表导出成功')
}

const handleExportAll = () => {
  const headers = ['月份', '员工工号', '员工姓名', '部门', '配件编号', '配件名称', '领用次数', '领用总数', '归还总数', '损毁总数', '未归还']
  const data = tableData.value.map(row => [
    `="${row.month}"`,
    row.employee_no,
    row.employee_name,
    row.department,
    row.part_no,
    row.part_name,
    row.borrow_count,
    row.total_borrow,
    row.total_return,
    row.total_damaged,
    row.unreturned
  ])
  exportToCsv(`全局明细报表_${formatDate(new Date())}.csv`, headers, data)
}

const handleExportMonth = () => {
  const headers = ['员工姓名', '部门', '配件名称', '领用次数', '领用数量', '归还数量', '损毁数量', '未归还']
  const data = currentMonthData.value.map(row => [
    row.employee_name,
    row.department,
    row.part_name,
    row.borrow_count,
    row.total_borrow,
    row.total_return,
    row.total_damaged,
    row.unreturned
  ])
  exportToCsv(`明细报表_${currentMonth.value}.csv`, headers, data)
}

const handlePushReport = () => {
  ElMessageBox.confirm('确定要根据当前筛选条件推送业务明细至钉钉吗？', '确认推送', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    try {
      const params = {
        dimension: 'detail',
        start_date: formatDate(dateRange.value[0]),
        end_date: formatDate(dateRange.value[1]),
      }
      const res = await pushReport(params)
      if (res.code === 200) {
        ElMessage.success('推送已启动')
      }
    } catch (error) {
      ElMessage.error('推送失败: ' + error.message)
    }
  })
}

const getMonthData = (month) => tableData.value.filter(item => item.month === month)

const getMonthTotal = (month, field) => {
  const data = getMonthData(month)
  return data.reduce((sum, item) => sum + (item[field] || 0), 0)
}

onMounted(() => {
  loadEmployees()
  loadParts()
  loadData()
})
</script>

<style scoped>
.detailed-report {
  max-width: 1400px;
  margin: 0 auto;
}

.filter-card {
  margin-bottom: 24px;
  border-radius: 12px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-group h2 { margin: 0; font-size: 20px; }
.subtitle { margin: 4px 0 0; font-size: 13px; color: #94a3b8; }

.divider { width: 1px; height: 24px; background: #e2e8f0; margin: 0 20px; }

.grid-section { min-height: 400px; padding: 10px 0; }

.grid-col { margin-bottom: 24px; }

.month-premium-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  border: 1px solid #f1f5f9;
  cursor: pointer;
  transition: all 0.3s ease;
  min-height: 160px;
  display: flex;
  flex-direction: column;
}

.month-premium-card:hover {
  box-shadow: 0 10px 25px rgba(0,0,0,0.06);
  transform: translateY(-5px);
  border-color: #3b82f6;
}

.month-label { font-size: 20px; font-weight: 800; color: #1e293b; margin-bottom: 8px; }

.month-status {
  font-size: 12px;
  color: #64748b;
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 15px;
}

.status-dot { width: 8px; height: 8px; border-radius: 50%; background: #10b981; }

.month-stats {
  margin-top: auto;
  background: #f8fafc;
  padding: 10px;
  border-radius: 10px;
  display: flex;
  justify-content: space-around;
}

.stat-item { text-align: center; }
.stat-item .lab { font-size: 11px; color: #94a3b8; display: block; }
.stat-item .val { font-size: 15px; font-weight: 700; color: #334155; }
.stat-item .val.red { color: #ef4444; }

.click-hint {
  position: absolute;
  top: 20px;
  right: 20px;
  font-size: 11px;
  color: #3b82f6;
  font-weight: 600;
  opacity: 0;
  transition: 0.3s;
}

.month-premium-card:hover .click-hint { opacity: 1; }

/* Dialog */
.dialog-header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f8fafc;
  padding: 12px 20px;
  border-radius: 8px;
  margin-bottom: 15px;
}

.month-summary-mini { font-size: 14px; color: #475569; }
.text-danger { color: #ef4444; }
.sub-text { font-size: 11px; color: #94a3b8; margin-left: 4px; }

:deep(.custom-dialog) { border-radius: 12px; }
:deep(.custom-dialog .el-dialog__header) { border-bottom: 1px solid #f1f5f9; padding: 20px; }
</style>
