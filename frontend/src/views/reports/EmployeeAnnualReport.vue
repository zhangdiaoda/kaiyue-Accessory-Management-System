<template>
  <div class="employee-annual-report">
    <!-- Header Section -->
    <el-card class="filter-card" shadow="never">
      <div class="header-content">
        <div class="title-group">
          <h2>📊 员工年度领用报表</h2>
          <p class="subtitle">追溯员工全年的领用轨迹与资产状态</p>
        </div>
        <div class="filter-actions">
          <el-form :inline="true" class="search-form">
            <el-form-item label="员工">
              <el-select
                v-model="selectedEmployee"
                placeholder="选择员工"
                filterable
                clearable
                @change="handleSearch"
                style="width: 220px"
              >
                <el-option
                  v-for="emp in employees"
                  :key="emp.id"
                  :label="`${emp.name} (${emp.employee_no})`"
                  :value="emp.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="年份">
              <el-date-picker
                v-model="selectedYear"
                type="year"
                placeholder="选择年份"
                @change="handleSearch"
                style="width: 140px"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="loading" @click="handleSearch">
                <el-icon><Search /></el-icon> 查询
              </el-button>
            </el-form-item>
          </el-form>
          <div class="divider"></div>
          <el-button 
            type="success" 
            plain
            @click="handleExportYear" 
            :disabled="!selectedEmployee || tableData.length === 0"
          >
            <el-icon><Download /></el-icon> 导出年度报表
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Annual Summary Section -->
    <transition name="el-zoom-in-top">
      <div v-if="selectedEmployee && !loading" class="summary-container">
        <el-row :gutter="20">
          <el-col :span="6" v-for="(item, index) in summaryStats" :key="index">
            <div class="summary-card" :style="{ borderLeft: `4px solid ${item.color}` }">
              <div class="summary-val" :style="{ color: item.color }">{{ item.value }}</div>
              <div class="summary-lab">{{ item.label }}</div>
              <el-icon class="bg-icon" :style="{ color: item.color }"><component :is="item.icon" /></el-icon>
            </div>
          </el-col>
        </el-row>
      </div>
    </transition>

    <!-- 12 Months Grid -->
    <div v-if="selectedEmployee && !loading" class="months-grid">
      <el-row :gutter="24">
        <el-col :span="6" v-for="month in 12" :key="month" class="grid-col">
          <div 
            class="month-premium-card" 
            :class="{ 'has-data': getMonthData(month).length > 0 }"
            @click="handleMonthClick(month)"
          >
            <div class="month-label">{{ month }}月</div>
            <div class="month-status">
              <span v-if="getMonthData(month).length > 0" class="status-dot online"></span>
              <span v-else class="status-dot offline"></span>
              {{ getMonthData(month).length > 0 ? '有领用记录' : '无领用记录' }}
            </div>
            
            <div class="month-stats" v-if="getMonthData(month).length > 0">
              <div class="stat-item">
                <span class="lab">领用数</span>
                <span class="val">{{ getMonthTotal(month, 'total_borrow') }}</span>
              </div>
              <div class="stat-item">
                <span class="lab">废品数</span>
                <span class="val" :class="{ 'red': getMonthTotal(month, 'total_damaged') > 0 }">
                  {{ getMonthTotal(month, 'total_damaged') }}
                </span>
              </div>
            </div>
            
            <div class="click-hint">点击查看详情 <el-icon><ArrowRight /></el-icon></div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- Empty State -->
    <div v-else-if="!selectedEmployee" class="empty-splash">
      <el-empty :image-size="200" description="请选择一位员工开始查看年度报表视角" />
    </div>

    <!-- Month Detail Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="`📅 ${selectedYear.getFullYear()}年${currentMonth}月 领用明细`"
      width="900px"
      class="custom-dialog"
      append-to-body
      destroy-on-close
    >
      <div class="dialog-header-actions">
        <div class="month-summary-mini">
          <span>本月领用总计: <strong>{{ getMonthTotal(currentMonth, 'total_borrow') }}</strong></span>
          <span style="margin-left: 20px">损毁总计: <strong class="text-danger">{{ getMonthTotal(currentMonth, 'total_damaged') }}</strong></span>
        </div>
        <el-button type="primary" size="small" @click="handleExportMonth">
          <el-icon><Download /></el-icon> 导出本月明细
        </el-button>
      </div>

      <el-table 
        :data="currentMonthData" 
        border 
        stripe 
        max-height="500px"
        v-loading="loading"
        highlight-current-row
      >
        <el-table-column prop="part_no" label="配件编号" width="130" sortable />
        <el-table-column prop="part_name" label="配件名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="borrow_count" label="领用频次" width="100" align="center" />
        <el-table-column prop="total_borrow" label="领用数量" width="100" align="center" />
        <el-table-column prop="total_return" label="归还数量" width="100" align="center" />
        <el-table-column prop="total_damaged" label="废品数量" width="100" align="center">
          <template #default="{ row }">
            <b :class="{ 'text-danger': row.total_damaged > 0 }">{{ row.total_damaged }}</b>
          </template>
        </el-table-column>
        <el-table-column prop="unreturned" label="待归还" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.unreturned > 0" type="warning" size="small">{{ row.unreturned }}</el-tag>
            <span v-else>0</span>
          </template>
        </el-table-column>
      </el-table>
      
      <div v-if="currentMonthData.length === 0" class="no-data-msg">
        该月份暂无任何领用记录
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getDetailedReport } from '@/api/report'
import { getAllEmployees } from '@/api/employee'

const loading = ref(false)
const employees = ref([])
const selectedEmployee = ref(null)
const selectedYear = ref(new Date())
const tableData = ref([])
const currentMonth = ref(null)
const dialogVisible = ref(false)

const summaryStats = computed(() => {
  const totalBorrow = tableData.value.reduce((sum, item) => sum + item.total_borrow, 0)
  const totalReturn = tableData.value.reduce((sum, item) => sum + item.total_return, 0)
  const totalDamaged = tableData.value.reduce((sum, item) => sum + item.total_damaged, 0)
  const unreturned = tableData.value.reduce((sum, item) => sum + item.unreturned, 0)
  
  return [
    { label: '全年领用总额', value: totalBorrow, color: '#4facfe', icon: 'Box' },
    { label: '累计归还记录', value: totalReturn, color: '#67C23A', icon: 'Checked' },
    { label: '全年废品总计', value: totalDamaged, color: '#F56C6C', icon: 'Warning' },
    { label: '当前未归还数', value: unreturned, color: '#E6A23C', icon: 'Timer' }
  ]
})

const currentMonthData = computed(() => {
  if (!currentMonth.value) return []
  const monthStr = `${selectedYear.value.getFullYear()}-${String(currentMonth.value).padStart(2, '0')}`
  return tableData.value.filter(item => item.month === monthStr)
})

const loadEmployees = async () => {
  try {
    const res = await getAllEmployees()
    employees.value = res.data || []
  } catch (error) {
    console.error('加载员工列表失败')
  }
}

const handleSearch = async () => {
  if (!selectedEmployee.value) return
  
  loading.value = true
  try {
    const year = selectedYear.value.getFullYear()
    const params = {
      employee_id: selectedEmployee.value,
      start_date: `${year}-01-01`,
      end_date: `${year}-12-31`
    }
    const res = await getDetailedReport(params)
    tableData.value = res.data || []
  } catch (error) {
    ElMessage.error('报表数据加载失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

const getMonthData = (month) => {
  const year = selectedYear.value.getFullYear()
  const monthStr = `${year}-${String(month).padStart(2, '0')}`
  return tableData.value.filter(item => item.month === monthStr)
}

const getMonthTotal = (month, field) => {
  if (!month) return 0
  const data = getMonthData(month)
  return data.reduce((sum, item) => sum + (item[field] || 0), 0)
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
  ElMessage.success('导出数据文件生成完成')
}

const handleExportYear = () => {
  const employeeName = employees.value.find(e => e.id === selectedEmployee.value)?.name || '未知员工'
  const headers = ['月份', '配件编号', '配件名称', '领用频次', '领用数量', '归还数量', '损毁数量', '待归还']
  const data = tableData.value.map(row => [
    `="${row.month}"`, // 强制 Excel 识别为文本
    row.part_no,
    row.part_name,
    row.borrow_count,
    row.total_borrow,
    row.total_return,
    row.total_damaged,
    row.unreturned
  ])
  exportToCsv(`${employeeName}_${selectedYear.value.getFullYear()}年度汇总报表.csv`, headers, data)
}

const handleExportMonth = () => {
  const employeeName = employees.value.find(e => e.id === selectedEmployee.value)?.name || '未知员工'
  const headers = ['配件编号', '配件名称', '领用频次', '领用数量', '归还数量', '损毁数量', '待归还']
  const data = currentMonthData.value.map(row => [
    row.part_no,
    row.part_name,
    row.borrow_count,
    row.total_borrow,
    row.total_return,
    row.total_damaged,
    row.unreturned
  ])
  exportToCsv(`${employeeName}_${selectedYear.value.getFullYear()}年${currentMonth.value}月明细报表.csv`, headers, data)
}

onMounted(() => {
  loadEmployees()
})
</script>

<style scoped>
.employee-annual-report {
  padding: 0;
  max-width: 1400px;
  margin: 0 auto;
}

.filter-card {
  margin-bottom: 24px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.05);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
}

.title-group h2 {
  margin: 0;
  font-size: 22px;
  color: #2c3e50;
}

.subtitle {
  margin: 4px 0 0;
  font-size: 13px;
  color: #94a3b8;
}

.filter-actions {
  display: flex;
  align-items: center;
}

.search-form {
  margin-bottom: 0;
}

.search-form :deep(.el-form-item) {
  margin-bottom: 0;
  margin-right: 15px;
}

.divider {
  width: 1px;
  height: 24px;
  background: #e2e8f0;
  margin: 0 20px;
}

/* Summary Section */
.summary-container {
  margin-bottom: 24px;
}

.summary-card {
  background: white;
  padding: 20px;
  border-radius: 12px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0,0,0,0.03);
  transition: transform 0.3s;
}

.summary-card:hover {
  transform: translateY(-3px);
}

.summary-val {
  font-size: 28px;
  font-weight: 800;
  line-height: 1;
}

.summary-lab {
  margin-top: 8px;
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

.bg-icon {
  position: absolute;
  right: -10px;
  bottom: -10px;
  font-size: 60px;
  opacity: 0.1;
  transform: rotate(-15deg);
}

/* 12 Months Grid */
.months-grid {
  margin-bottom: 40px;
}

.grid-col {
  margin-bottom: 24px;
}

.month-premium-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  border: 1px solid #f1f5f9;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  position: relative;
  min-height: 180px;
  display: flex;
  flex-direction: column;
}

.month-premium-card:hover {
  box-shadow: 0 10px 30px rgba(0,0,0,0.08);
  border-color: #3b82f6;
  transform: scale(1.02);
}

.month-label {
  font-size: 24px;
  font-weight: 800;
  color: #1e293b;
  margin-bottom: 8px;
}

.month-status {
  font-size: 12px;
  color: #64748b;
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 20px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.online { background: #10b981; box-shadow: 0 0 8px rgba(16,185,129,0.5); }
.status-dot.offline { background: #cbd5e1; }

.month-stats {
  margin-top: auto;
  background: #f8fafc;
  padding: 12px;
  border-radius: 12px;
  display: flex;
  justify-content: space-around;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-item .lab {
  font-size: 11px;
  color: #94a3b8;
  margin-bottom: 2px;
}

.stat-item .val {
  font-size: 16px;
  font-weight: 700;
  color: #334155;
}

.stat-item .val.red { color: #ef4444; }

.click-hint {
  position: absolute;
  top: 24px;
  right: 24px;
  font-size: 11px;
  color: #3b82f6;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 2px;
  opacity: 0;
  transition: 0.3s;
}

.month-premium-card:hover .click-hint {
  opacity: 1;
}

/* Empty Splash */
.empty-splash {
  margin-top: 60px;
}

/* Dialog Styles */
.dialog-header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f8fafc;
  padding: 15px 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.month-summary-mini {
  font-size: 15px;
  color: #475569;
}

.text-danger {
  color: #ef4444;
}

.no-data-msg {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
  font-style: italic;
}

:deep(.custom-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.custom-dialog .el-dialog__header) {
  margin-right: 0;
  padding: 20px 25px;
  border-bottom: 1px solid #f1f5f9;
}

:deep(.custom-dialog .el-dialog__title) {
  font-weight: 700;
  font-size: 18px;
}
</style>
