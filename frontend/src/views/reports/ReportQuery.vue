<template>
  <div class="report-query">
    <!-- Filter Section -->
    <el-card class="filter-card" shadow="never">
      <div class="header-content">
        <div class="title-group">
          <h2>🔍 维度统计查询</h2>
          <p class="subtitle">从配件、员工、部门三个视角洞察仓储动态</p>
        </div>
        <div class="filter-actions">
          <el-form :inline="true" class="search-form">
            <el-form-item label="统计维度">
              <el-select v-model="dimension" @change="loadData" style="width: 140px">
                <el-option label="配件视角" value="part" />
                <el-option label="员工视角" value="employee" />
                <el-option label="部门视角" value="department" />
              </el-select>
            </el-form-item>
            <el-form-item label="日期范围">
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始"
                end-placeholder="结束"
                style="width: 240px"
                @change="loadData"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadData" :loading="loading">查询</el-button>
            </el-form-item>
          </el-form>
          <div class="divider"></div>
          <el-button type="success" plain @click="handleExport" :disabled="tableData.length === 0">
            <el-icon><Download /></el-icon> 导出结果
          </el-button>
          <el-button type="primary" plain @click="handlePushReport" :disabled="tableData.length === 0">
            <el-icon><Share /></el-icon> 推送报表
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Results Card Wall -->
    <div v-loading="loading" class="card-wall-section">
      <el-row :gutter="20" v-if="tableData.length > 0">
        <el-col :span="6" v-for="(item, index) in tableData" :key="index" class="grid-col">
          <div class="perspective-card" @click="handleItemClick(item)">
            <div class="card-top">
              <div class="main-label">
                {{ getDisplayName(item) }}
                <span class="sub-id" v-if="getDisplayId(item)">({{ getDisplayId(item) }})</span>
              </div>
              <el-tag size="small" :type="getTagType()">{{ getDimensionLabel() }}</el-tag>
            </div>
            
            <div class="card-stats-grid">
              <div class="stat-item">
                <span class="l">领用总数</span>
                <span class="v">{{ item.total_borrow }}</span>
              </div>
              <div class="stat-item">
                <span class="l">领用频次</span>
                <span class="v">{{ item.borrow_count }}</span>
              </div>
              <div class="stat-item">
                <span class="l">废品总计</span>
                <span class="v" :class="{ 'red': item.total_damaged > 0 }">{{ item.total_damaged }}</span>
              </div>
              <div class="stat-item" v-if="dimension === 'part'">
                <span class="l">废品率</span>
                <span class="v" :class="{ 'red': item.damage_rate > 10 }">{{ item.damage_rate || 0 }}%</span>
              </div>
              <div class="stat-item" v-if="dimension === 'department'">
                <span class="l">覆盖员工</span>
                <span class="v">{{ item.employee_count }}人</span>
              </div>
            </div>
            
            <div class="card-footer">
              <span>点击分析趋势</span>
              <el-icon><DataLine /></el-icon>
            </div>
          </div>
        </el-col>
      </el-row>
      <el-empty v-else description="查询期间内无相关数据，请尝试调整筛选条件" />
    </div>

    <!-- Detail Analysis Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="`📈 ${currentEntityName} - 趋势分析视图`"
      width="800px"
      append-to-body
      class="custom-dialog"
    >
      <div class="dialog-summary">
        <div class="sum-box">
          <span class="sum-l">查询周期内总领用</span>
          <span class="sum-v">{{ currentItem?.total_borrow }}</span>
        </div>
        <div class="sum-box">
          <span class="sum-l">查询周期内废品总计</span>
          <span class="sum-v red">{{ currentItem?.total_damaged }}</span>
        </div>
      </div>
      
      <p style="color: #64748b; font-size: 13px; margin: 20px 0 10px;">
        💡 提示：如需更详尽的单笔原始记录，请前往<b>“明细查询”</b>频道查看。
      </p>
      
      <div v-if="dimension === 'part' && currentItem" class="extra-info-table">
         <el-descriptions :column="2" border>
            <el-descriptions-item label="配件编号">{{ currentItem.part_no }}</el-descriptions-item>
            <el-descriptions-item label="库存健康度">
              <el-tag :type="currentItem.damage_rate > 10 ? 'danger' : 'success'">
                {{ currentItem.damage_rate > 10 ? '低' : '高' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="归还总数">{{ currentItem.total_return }}</el-descriptions-item>
            <el-descriptions-item label="待归还数">{{ currentItem.total_borrow - currentItem.total_return - currentItem.total_damaged }}</el-descriptions-item>
         </el-descriptions>
      </div>

      <template #footer>
        <el-button @click="dialogVisible = false">关闭视图</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getPartReport, getEmployeeReport, getDepartmentReport, pushReport } from '@/api/report'
import { Download, Share, DataLine } from '@element-plus/icons-vue'

const dimension = ref('part')
const tableData = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const currentItem = ref(null)

const dateRange = ref([
  new Date(new Date().setMonth(new Date().getMonth() - 1)),
  new Date()
])

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      start_date: formatDate(dateRange.value[0]),
      end_date: formatDate(dateRange.value[1])
    }

    let res
    if (dimension.value === 'part') {
      res = await getPartReport(params)
    } else if (dimension.value === 'employee') {
      res = await getEmployeeReport(params)
    } else {
      res = await getDepartmentReport(params)
    }

    tableData.value = res.data || []
  } catch (error) {
    ElMessage.error('查询结果加载失败')
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const getDisplayName = (item) => {
  if (dimension.value === 'part') return item.part_name
  if (dimension.value === 'employee') return item.employee_name
  return item.department
}

const getDisplayId = (item) => {
  if (dimension.value === 'part') return item.part_no
  if (dimension.value === 'employee') return item.employee_no
  return null
}

const getDimensionLabel = () => {
  const m = { part: '配件', employee: '员工', department: '部门' }
  return m[dimension.value]
}

const getTagType = () => {
  const m = { part: '', employee: 'success', department: 'warning' }
  return m[dimension.value]
}

const currentEntityName = computed(() => {
  if (!currentItem.value) return ''
  return getDisplayName(currentItem.value)
})

const handleItemClick = (item) => {
  currentItem.value = item
  dialogVisible.value = true
}

const handleExport = () => {
  // ... existing export logic ...
  if (tableData.value.length === 0) return

  let headers = []
  let data = []

  if (dimension.value === 'part') {
    headers = ['配件编号', '配件名称', '领用次数', '领用总数', '归还总数', '损毁总数', '损毁率(%)']
    data = tableData.value.map(row => [
      row.part_no,
      row.part_name,
      row.borrow_count,
      row.total_borrow,
      row.total_return,
      row.total_damaged,
      row.damage_rate || 0
    ])
  } else if (dimension.value === 'employee') {
    headers = ['工号', '姓名', '部门', '领用次数', '领用总数', '损毁总数']
    data = tableData.value.map(row => [
      row.employee_no,
      row.employee_name,
      row.department,
      row.borrow_count,
      row.total_borrow,
      row.total_damaged
    ])
  } else {
    headers = ['部门', '员工数', '领用次数', '领用总数', '损毁总数']
    data = tableData.value.map(row => [
      row.department,
      row.employee_count,
      row.borrow_count,
      row.total_borrow,
      row.total_damaged
    ])
  }

  const csvContent = [
    headers.join(','),
    ...data.map(row => row.join(','))
  ].join('\n')

  const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.setAttribute('href', url)
  link.setAttribute('download', `维度查询_${dimension.value}_${formatDate(new Date())}.csv`)
  link.click()
  ElMessage.success('报表导出成功')
}

const handlePushReport = () => {
  const dimLabel = { part: '配件视角', employee: '员工视角', department: '部门视角' }[dimension.value]
  ElMessageBox.confirm(`确定要推送当前的 [${dimLabel}] 统计排行至钉钉吗？`, '确认推送', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    try {
      const params = {
        dimension: dimension.value,
        start_date: formatDate(dateRange.value[0]),
        end_date: formatDate(dateRange.value[1]),
      }
      const res = await pushReport(params)
      if (res.code === 200) {
        ElMessage.success('报表已推送')
      }
    } catch (error) {
      ElMessage.error('推送失败: ' + error.message)
    }
  })
}

onMounted(loadData)
</script>

<style scoped>
.report-query {
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

.card-wall-section { min-height: 400px; }
.grid-col { margin-bottom: 20px; }

.perspective-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  border: 1px solid #f1f5f9;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  display: flex;
  flex-direction: column;
  min-height: 200px;
}

.perspective-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 30px rgba(0,0,0,0.07);
  border-color: #3b82f6;
}

.card-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.main-label {
  font-size: 17px;
  font-weight: 800;
  color: #1e293b;
  line-height: 1.3;
}

.sub-id {
  display: block;
  font-size: 11px;
  color: #94a3b8;
  font-weight: 500;
  margin-top: 2px;
}

.card-stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
  margin-bottom: 20px;
}

.stat-item {
  display: flex;
  flex-direction: column;
}

.stat-item .l { font-size: 11px; color: #94a3b8; margin-bottom: 2px; }
.stat-item .v { font-size: 16px; font-weight: 700; color: #334155; }
.stat-item .v.red { color: #ef4444; }

.card-footer {
  margin-top: auto;
  border-top: 1px solid #f8fafc;
  padding-top: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #3b82f6;
  font-weight: 600;
}

/* Dialog */
.dialog-summary {
  display: flex;
  gap: 20px;
  background: #f8fafc;
  padding: 20px;
  border-radius: 12px;
}

.sum-box {
  flex: 1;
  text-align: center;
}

.sum-l { display: block; font-size: 12px; color: #64748b; margin-bottom: 8px; }
.sum-v { font-size: 24px; font-weight: 800; color: #1e293b; }
.sum-v.red { color: #ef4444; }

.extra-info-table { margin-top: 20px; }

:deep(.custom-dialog) { border-radius: 16px; }
</style>
