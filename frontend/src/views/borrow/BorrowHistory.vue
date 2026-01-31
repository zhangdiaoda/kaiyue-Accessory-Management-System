<template>
  <div class="borrow-history">
    <el-card>
      <template #header>
        <span>领用记录</span>
      </template>

      <el-form :model="searchForm" :inline="true" class="search-form">
        <el-form-item>
          <el-select v-model="searchForm.status" placeholder="全部状态" clearable>
            <el-option label="已领用" value="BORROWED" />
            <el-option label="部分归还" value="PARTIAL_RETURNED" />
            <el-option label="已归还" value="RETURNED" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
        </el-form-item>
      </el-form>

      <!-- 移动端卡片布局 -->
      <div v-if="isMobile" class="borrow-cards">
        <div v-for="item in tableData" :key="item.id" class="borrow-card">
          <div class="card-header">
            <div class="header-left">
              <div class="employee-name">{{ item.employee_name }}</div>
              <div class="record-no">{{ item.record_no }}</div>
            </div>
            <el-tag
              :type="item.status === 'RETURNED' ? 'success' : item.status === 'PARTIAL_RETURNED' ? 'warning' : 'primary'"
              effect="light"
              size="small"
              round
            >
              {{ getStatusText(item.status) }}
            </el-tag>
          </div>
          
          <div class="card-body">
            <div class="card-row">
              <span class="label">配件</span>
              <span class="value">{{ item.part_name }}</span>
            </div>
            <div class="card-row">
              <span class="label">配件编号</span>
              <span class="value">{{ item.part_no }}</span>
            </div>
            <div class="card-row">
              <span class="label">领用 / 回收 / 报废</span>
              <span class="value">
                {{ item.borrow_quantity }} / {{ item.return_quantity }} / 
                <span :class="{ 'text-danger': item.damaged_quantity > 0 }">{{ item.damaged_quantity }}</span>
              </span>
            </div>
            <div class="card-row">
              <span class="label">领用时间</span>
              <span class="value">{{ item.borrow_time }}</span>
            </div>
            <div class="card-row" v-if="item.remark">
              <span class="label">备注</span>
              <span class="value">{{ item.remark }}</span>
            </div>
          </div>
          
          <div class="card-footer" v-if="item.status !== 'RETURNED'">
            <el-button type="primary" size="small" @click="handleReturn(item)" class="action-btn">
              处置旧件
            </el-button>
          </div>
        </div>
        
        <el-empty v-if="tableData.length === 0" description="暂无领用记录" :image-size="80" />
      </div>

      <!-- 桌面端表格布局 -->
      <el-table v-else :data="tableData" border stripe>
        <el-table-column prop="record_no" label="记录编号" width="150" />
        <el-table-column prop="employee_name" label="员工" width="100" />
        <el-table-column prop="part_name" label="配件" width="150" />
        <el-table-column prop="part_no" label="配件编号" width="120" />
        <el-table-column prop="borrow_quantity" label="领用" width="80" align="center" />
        <el-table-column prop="return_quantity" label="已回收" width="80" align="center" />
        <el-table-column prop="damaged_quantity" label="已报废" width="80" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.damaged_quantity > 0" type="danger" effect="plain" class="apple-tag-mini">
              {{ row.damaged_quantity }}
            </el-tag>
            <span v-else class="text-muted">0</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 'RETURNED' ? 'success' : row.status === 'PARTIAL_RETURNED' ? 'warning' : 'primary'"
              effect="light"
              round
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="borrow_time" label="领用时间" width="180" />
        <el-table-column prop="return_time" label="最后处置时间" width="180">
          <template #default="{ row }">
            {{ row.return_time || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button 
              v-if="row.status !== 'RETURNED'"
              link 
              type="primary" 
              @click="handleReturn(row)"
              class="apple-link-btn"
            >
              处置旧件
            </el-button>
          </template>
        </el-table-column>

      </el-table>

      <!-- 处置对话框 -->
      <el-dialog 
        v-model="returnVisible" 
        title="旧件处置登记" 
        width="420px" 
        class="apple-dialog"
        append-to-body
      >
        <div class="return-form-container">
          <div class="return-info-card glass-effect">
            <div class="info-row">
              <span class="info-label">持有员工</span>
              <span class="info-value">{{ selectedRow?.employee_name }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">待处理旧件</span>
              <span class="info-value">{{ selectedRow?.part_name }}</span>
            </div>
            <div class="info-row highlight">
              <span class="info-label">待回收总数</span>
              <span class="info-value text-danger">{{ getUnreturnedCount(selectedRow) }} <small>件</small></span>
            </div>
          </div>

          <el-form :model="returnForm" label-position="top" class="apple-form">
            <el-form-item label="处置选项" style="margin-bottom: 24px">
              <el-radio-group v-model="disposeType" @change="handleTypeChange" class="apple-radio-group">
                <el-radio-button label="RECYCLE">回收入库</el-radio-button>
                <el-radio-button label="SCRAP">报废处理</el-radio-button>
                <el-radio-button label="MANUAL">手动分配</el-radio-button>
              </el-radio-group>
            </el-form-item>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="回收入库量">
                  <el-input-number 
                    v-model="returnForm.return_quantity" 
                    :min="0" 
                    :max="getUnreturnedCount(selectedRow) - returnForm.damaged_quantity" 
                    :disabled="disposeType !== 'MANUAL'"
                    class="apple-giant-input"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="报废处理量">
                  <el-input-number 
                    v-model="returnForm.damaged_quantity" 
                    :min="0" 
                    :max="getUnreturnedCount(selectedRow) - returnForm.return_quantity" 
                    :disabled="disposeType !== 'MANUAL'"
                    class="apple-giant-input"
                  />
                </el-form-item>
              </el-col>
            </el-row>

          </el-form>

        </div>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="returnVisible = false" round>暂不处理</el-button>
            <el-button type="primary" @click="submitReturn" :loading="returning" class="apple-btn" round>
              确认处置
            </el-button>
          </div>
        </template>
      </el-dialog>




      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadData"
        @current-change="loadData"
        style="margin-top: 20px"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { useResponsive } from '@/composables/useResponsive'
import { getBorrowRecordList, returnBorrowRecord } from '@/api/borrow'

const { isMobile } = useResponsive()

const tableData = ref([])
const returnVisible = ref(false)
const returning = ref(false)
const selectedRow = ref(null)
const disposeType = ref('RECYCLE')
const returnForm = reactive({
  return_quantity: 0,
  damaged_quantity: 0
})

const handleTypeChange = (val) => {
  const total = getUnreturnedCount(selectedRow.value)
  if (val === 'RECYCLE') {
    returnForm.return_quantity = total
    returnForm.damaged_quantity = 0
  } else if (val === 'SCRAP') {
    returnForm.return_quantity = 0
    returnForm.damaged_quantity = total
  }
}

const getUnreturnedCount = (row) => {
  if (!row) return 0
  return row.borrow_quantity - row.return_quantity - row.damaged_quantity
}

const handleReturn = (row) => {
  selectedRow.value = row
  disposeType.value = 'RECYCLE'
  returnForm.return_quantity = getUnreturnedCount(row)
  returnForm.damaged_quantity = 0
  returnVisible.value = true
}

const submitReturn = async () => {
  returning.value = true
  try {
    await returnBorrowRecord(selectedRow.value.id, returnForm)
    ElMessage.success('归还登记成功')
    returnVisible.value = false
    loadData()
  } catch (error) {
    console.error('归还失败', error)
  } finally {
    returning.value = false
  }
}

const searchForm = reactive({
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadData = async () => {
  try {
    const res = await getBorrowRecordList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: searchForm.status
    })
    tableData.value = res.data.records || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载数据失败')
  }
}

const getStatusText = (status) => {
  const map = {
    'BORROWED': '待回收',
    'PARTIAL_RETURNED': '部分回收',
    'RETURNED': '已结清'
  }
  return map[status] || status
}

loadData()
</script>

<style scoped>
.borrow-history {
  height: 100%;
}

.search-form {
  margin-bottom: 20px;
}

.apple-tag-mini {
  border-radius: 4px;
  font-weight: 700;
  scale: 0.9;
}

.apple-link-btn {
  font-weight: 600;
  letter-spacing: 0.5px;
}

.return-info-card {
  padding: 20px;
  border-radius: 16px;
  margin-bottom: 24px;
  border: 1px solid rgba(0,0,0,0.05);
  background: rgba(0,0,0,0.02);
}

.info-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  font-size: 14px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-label {
  color: #86868b;
}

.info-value {
  color: #1d1d1f;
  font-weight: 600;
}

.info-row.highlight {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid rgba(0,0,0,0.05);
}

.info-row.highlight .info-value {
  font-size: 18px;
}

.info-row.highlight .info-value small {
  font-size: 12px;
  font-weight: 400;
}

/* iOS Segmented Control Style */
:deep(.apple-radio-group .el-radio-button__inner) {
  background: transparent;
  border: none;
  border-radius: 6px;
  box-shadow: none !important;
  color: #1d1d1f;
  font-weight: 500;
  padding: 8px 16px;
  transition: all 0.3s;
}

:deep(.apple-radio-group .el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background: white;
  color: #000;
  box-shadow: 0 2px 8px rgba(0,0,0,0.12) !important;
  font-weight: 600;
}

:deep(.apple-radio-group) {
  background: #7676801f;
  border-radius: 8px;
  padding: 2px;
  display: inline-flex;
}

/* iOS Stepper Style */
.apple-giant-input {
  width: 100%;
}

:deep(.apple-giant-input .el-input__wrapper) {
  height: 44px;
  border-radius: 12px !important;
  font-size: 18px;
  font-weight: 800;
  background-color: #f2f2f7;
  box-shadow: none !important;
  padding-left: 45px;
  padding-right: 45px;
  transition: all 0.2s;
}

:deep(.apple-giant-input .el-input-number__decrease),
:deep(.apple-giant-input .el-input-number__increase) {
  width: 36px;
  height: 36px;
  top: 4px;
  background: white;
  border-radius: 10px;
  border: none;
  box-shadow: 0 2px 5px rgba(0,0,0,0.05);
  color: #0071e3;
  z-index: 2;
  transition: all 0.2s;
}

:deep(.apple-giant-input .el-input-number__decrease:hover),
:deep(.apple-giant-input .el-input-number__increase:hover) {
  color: #0071e3;
  background: #fff;
  transform: scale(1.05);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

:deep(.apple-giant-input .el-input-number__decrease) {
  left: 4px;
  border-right: none;
}

:deep(.apple-giant-input .el-input-number__increase) {
  right: 4px;
  border-left: none;
}

:deep(.apple-giant-input .el-input__inner) {
  text-align: center;
  color: #1d1d1f;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.text-danger {
  color: #ff375f;
}

/* ========== 移动端卡片布局样式 ========== */
.borrow-cards {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.borrow-card {
  background: white;
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  border: 1px solid rgba(0,0,0,0.05);
  transition: transform 0.2s, box-shadow 0.2s;
}

.borrow-card:active {
  transform: scale(0.98);
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
}

.borrow-card .card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(0,0,0,0.05);
}

.borrow-card .header-left {
  flex: 1;
}

.borrow-card .employee-name {
  font-size: 16px;
  font-weight: 700;
  color: #1d1d1f;
  margin-bottom: 4px;
}

.borrow-card .record-no {
  font-size: 12px;
  color: #86868b;
  font-family: monospace;
}

.borrow-card .card-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.borrow-card .card-row {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  align-items: center;
}

.borrow-card .card-row .label {
  color: #86868b;
  font-size: 13px;
}

.borrow-card .card-row .value {
  color: #1d1d1f;
  font-weight: 600;
  text-align: right;
  flex: 1;
  margin-left: 12px;
}

.borrow-card .card-footer {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid rgba(0,0,0,0.05);
}

.borrow-card .action-btn {
  width: 100%;
  height: 44px;
  font-weight: 600;
}

/* 移动端响应式 */
@media (max-width: 768px) {
  .search-form {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  
  .search-form :deep(.el-form-item) {
    margin-right: 0;
    margin-bottom: 0;
    width: 100%;
  }
  
  .search-form :deep(.el-select) {
    width: 100%;
  }
  
  /* 弹窗表单双列改单列 */
  .return-form-container :deep(.el-row) {
    flex-direction: column;
  }
  
  .return-form-container :deep(.el-col) {
    max-width: 100%;
    flex: 0 0 100%;
  }
  
  /* 单选组垂直布局 */
  :deep(.apple-radio-group) {
    flex-direction: column;
    width: 100%;
  }
  
  :deep(.apple-radio-group .el-radio-button) {
    width: 100%;
  }
  
  :deep(.apple-radio-group .el-radio-button__inner) {
    width: 100%;
  }
  
  .dialog-footer {
    flex-direction: column;
    gap: 8px;
  }
  
  .dialog-footer .el-button {
    width: 100%;
  }
}
</style>


