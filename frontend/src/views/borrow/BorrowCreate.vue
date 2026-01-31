<template>
  <div class="borrow-create">
    <el-card>
      <template #header>
        <span>登记领用</span>
      </template>

      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px">
        <el-form-item label="选择员工" prop="employee_id" required>
          <el-select 
            v-model="formData.employee_id" 
            placeholder="请选择员工" 
            filterable
            @change="checkUnreturnedRecords"
          >
            <el-option
              v-for="emp in employees"
              :key="emp.id"
              :label="`${emp.name} (${emp.employee_no})`"
              :value="emp.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="选择配件" prop="part_id" required>
          <el-select 
            v-model="formData.part_id" 
            placeholder="请选择配件" 
            filterable
            @change="checkUnreturnedRecords"
          >
            <el-option
              v-for="part in parts"
              :key="part.id"
              :label="`${part.name} (${part.part_no}) - 库存: ${part.stock_quantity}`"
              :value="part.id"
            />
          </el-select>
        </el-form-item>

        <!-- 待回收旧件看板 (极致警示版) -->
        <div v-if="unreturnedInfo.has_unreturned" class="old-part-recovery apple-alert-box animate-slide-up">
          <div class="recovery-header">
            <div class="header-left">
              <el-icon :size="22" color="#ff9f0a"><WarningFilled /></el-icon>
              <span class="recovery-title">异常资产警告：检测到未结清旧件</span>
            </div>
            <el-tag type="warning" effect="dark" round size="small">流程挂起中</el-tag>
          </div>
          
          <div class="recovery-body">
            <div class="recovery-stats-grid">
              <div class="stat-node">
                <span class="n-label">累计发放</span>
                <span class="n-val">{{ unreturnedInfo.borrow_quantity }}</span>
              </div>
              <div class="stat-node">
                <span class="n-label">员工持有</span>
                <span class="n-val highlight">{{ unreturnedInfo.unreturned }}</span>
              </div>
              <div class="stat-node">
                <span class="n-label">待销账</span>
                <span class="n-val">{{ unreturnedInfo.unreturned - recoveryForm.return_quantity - recoveryForm.damaged_quantity }}</span>
              </div>
            </div>

            <div class="advanced-recovery-form">
              <div class="input-group-row">
                <div class="field-box">
                  <label>回收入库 (Recycle)</label>
                  <div class="input-with-action">
                    <el-input-number 
                      v-model="recoveryForm.return_quantity" 
                      :min="0" 
                      :max="unreturnedInfo.unreturned - recoveryForm.damaged_quantity"
                      class="apple-giant-input"
                    />
                    <el-button class="apple-pill-btn" @click="fillMax('RETURN')">全量填充</el-button>
                  </div>
                </div>

                <div class="field-box">
                  <label>报废处理 (Scrap)</label>
                  <div class="input-with-action">
                    <el-input-number 
                      v-model="recoveryForm.damaged_quantity" 
                      :min="0" 
                      :max="unreturnedInfo.unreturned - recoveryForm.return_quantity"
                      class="apple-giant-input"
                    />


                    <el-button class="apple-pill-btn" @click="fillMax('SCRAP')">全量填充</el-button>
                  </div>
                </div>
              </div>

              <div class="disposal-actions">
                <el-button 
                  type="primary" 
                  class="apple-sub-btn" 
                  :disabled="recoveryForm.return_quantity + recoveryForm.damaged_quantity === 0"
                  @click="submitRecovery"
                >
                  确 定 处 置
                </el-button>
                <el-button 
                  class="apple-sub-btn secondary" 
                  @click="quickHandleAll('SCRAP')"
                >
                  一 键 全 部 报 废
                </el-button>
              </div>
            </div>
          </div>



          <div class="sync-option">
            <el-checkbox v-model="formData.sync_dispose">
              新领用提交时，同步将剩余旧件自动报废 (Scrap)
            </el-checkbox>
          </div>
        </div>

        <el-divider v-if="unreturnedInfo.has_unreturned" content-position="left">本次新领用</el-divider>

        <el-form-item label="本次领用数量" prop="quantity" required>
          <div class="quantity-input-wrapper">
            <el-input-number v-model="formData.quantity" :min="1" controls-position="right" />
            <span class="unit-text">件</span>
          </div>
        </el-form-item>

        <el-form-item label="领用备注">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="2"
            placeholder="备注领用用途或特殊说明..."
            class="apple-input"
          />
        </el-form-item>

        <div class="submit-guard-notice" v-if="unreturnedInfo.has_unreturned && !formData.sync_dispose">
          <el-icon><InfoFilled /></el-icon>
          <span>请先执行“旧件处置”或勾选“同步处置”，以确保资产流转严谨。</span>
        </div>

        <el-form-item class="submit-section">
          <el-button 
            type="primary" 
            size="large" 
            @click="handleSubmit" 
            :loading="submitting" 
            class="apple-btn-large"
            :disabled="unreturnedInfo.has_unreturned && !formData.sync_dispose"
          >
            {{ unreturnedInfo.has_unreturned && formData.sync_dispose ? '处置并领用' : '确认登记领用' }}
          </el-button>
          <el-button size="large" @click="handleReset" round>重置</el-button>
        </el-form-item>

      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Warning, InfoFilled, WarningFilled } from '@element-plus/icons-vue'
import { createBorrowRecord, checkUnreturned } from '@/api/borrow'
import { returnBorrowRecord } from '@/api/borrow'
import { getAllEmployees } from '@/api/employee'
import { getPartList } from '@/api/part'

const formRef = ref(null)
const employees = ref([])
const parts = ref([])
const submitting = ref(false)

const formData = reactive({
  employee_id: '',
  part_id: '',
  quantity: 1,
  return_quantity: 0,
  return_damaged: 0,
  remark: '',
  sync_dispose: true // 默认开启同步报废
})

const unreturnedInfo = reactive({
  has_unreturned: false,
  record_id: null,
  borrow_quantity: 0,
  returned_quantity: 0,
  damaged_quantity: 0,
  unreturned: 0,
  borrow_time: ''
})

const recoveryForm = reactive({
  return_quantity: 0,
  damaged_quantity: 0
})

// 检查未归还记录并更新看板
const checkUnreturnedRecords = async () => {
  if (!formData.employee_id || !formData.part_id) return

  try {
    const res = await checkUnreturned({
      employee_id: formData.employee_id,
      part_id: formData.part_id
    })
    Object.assign(unreturnedInfo, res.data)
    // 重置处置表单
    recoveryForm.return_quantity = 0
    recoveryForm.damaged_quantity = 0
  } catch (error) {
    console.error('检查旧件失败', error)
  }
}

// 自动填充最大配额
const fillMax = (type) => {
  const remaining = unreturnedInfo.unreturned
  if (type === 'RETURN') {
    recoveryForm.return_quantity = remaining - recoveryForm.damaged_quantity
  } else {
    recoveryForm.damaged_quantity = remaining - recoveryForm.return_quantity
  }
}

// 提及指定的处置数量
const submitRecovery = async () => {
  if (recoveryForm.return_quantity + recoveryForm.damaged_quantity === 0) return

  try {
    submitting.value = true
    await returnBorrowRecord(unreturnedInfo.record_id, {
      return_quantity: recoveryForm.return_quantity,
      damaged_quantity: recoveryForm.damaged_quantity
    })
    ElMessage.success('处置成功')
    await checkUnreturnedRecords()
  } catch (error) {
    ElMessage.error('处置失败')
  } finally {
    submitting.value = false
  }
}

// 一键全部处置 (如全部报废)
const quickHandleAll = async (type) => {
  const isScrap = type === 'SCRAP'
  const title = isScrap ? '确认全部报废' : '确认全部回收'
  
  try {
    await ElMessageBox.confirm(`确认将持有的 ${unreturnedInfo.unreturned} 件旧件全部处理吗？`, title, {
      type: isScrap ? 'error' : 'success',
      roundButton: true
    })
    
    submitting.value = true
    await returnBorrowRecord(unreturnedInfo.record_id, {
      return_quantity: isScrap ? 0 : unreturnedInfo.unreturned,
      damaged_quantity: isScrap ? unreturnedInfo.unreturned : 0
    })
    ElMessage.success('已全部处置')
    await checkUnreturnedRecords()
  } catch (error) {
    // skip cancel
  } finally {
    submitting.value = false
  }
}

const loadEmployees = async () => {
  try {
    const res = await getAllEmployees()
    employees.value = res.data || []
  } catch (error) {
    ElMessage.error('加载员工列表失败')
  }
}

const loadParts = async () => {
  try {
    const res = await getPartList({ page: 1, pageSize: 1000 })
    parts.value = (res.data.records || []).filter(p => p.stock_quantity > 0)
  } catch (error) {
    ElMessage.error('加载配件列表失败')
  }
}

const handleSubmit = async () => {
  await formRef.value.validate()

  submitting.value = true
  try {
    // 如果勾选了同步处置且有旧件，先执行静默报废
    if (unreturnedInfo.has_unreturned && formData.sync_dispose) {
      await returnBorrowRecord(unreturnedInfo.record_id, {
        return_quantity: 0,
        damaged_quantity: unreturnedInfo.unreturned
      })
    }

    // 执行新领用登记
    await createBorrowRecord({
      employee_id: formData.employee_id,
      part_id: formData.part_id,
      borrow_quantity: formData.quantity,
      remark: formData.remark
    })

    ElMessage.success(unreturnedInfo.has_unreturned && formData.sync_dispose ? '已处置旧件并成功登记领用' : '领用登记成功')
    handleReset()
    loadParts()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
  Object.assign(formData, {
    employee_id: '',
    part_id: '',
    quantity: 1,
    sync_dispose: true,
    remark: ''
  })
  Object.assign(unreturnedInfo, {
    has_unreturned: false,
    record_id: null,
    borrow_quantity: 0,
    returned_quantity: 0,
    damaged_quantity: 0,
    unreturned: 0,
    borrow_time: ''
  })
}


onMounted(() => {
  loadEmployees()
  loadParts()
})
</script>

<style scoped>
.apple-alert-box {
  padding: 24px;
  border-radius: 24px;
  margin-bottom: 30px;
  border: 1px solid #ff9f0a;
  background: #fffcf0;
  box-shadow: 0 12px 40px rgba(255, 159, 10, 0.08);
  position: relative;
  overflow: hidden;
}

.apple-alert-box::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 6px;
  background: #ff9f0a;
}

.recovery-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.recovery-title {
  font-weight: 800;
  color: #1d1d1f;
  font-size: 18px;
  letter-spacing: -0.01em;
}

.recovery-stats-grid {
  display: flex;
  background: white;
  padding: 20px;
  border-radius: 18px;
  margin-bottom: 24px;
  border: 1px solid rgba(0,0,0,0.03);
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
}

.stat-node {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  border-right: 1px solid #f2f2f7;
}

.stat-node:last-child { border-right: none; }

.n-label {
  font-size: 11px;
  color: #86868b;
  font-weight: 600;
  text-transform: uppercase;
  margin-bottom: 4px;
}

.n-val {
  font-size: 24px;
  font-weight: 800;
  color: #1d1d1f;
}

.n-val.highlight { color: #ff9f0a; }

.advanced-recovery-form {
  background: rgba(255, 255, 255, 0.5);
  backdrop-filter: blur(20px);
  padding: 28px;
  border-radius: 20px;
  border: 1px solid white;
}

.input-group-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
  margin-bottom: 24px;
}

.field-box label {
  display: block;
  font-size: 14px;
  font-weight: 700;
  color: #1d1d1f;
  margin-bottom: 12px;
}

.input-with-action {
  display: flex;
  align-items: center;
  gap: 12px;
}

.apple-giant-input {
  width: 160px;
}

:deep(.apple-giant-input .el-input__wrapper) {
  height: 48px;
  border-radius: 14px !important;
  font-size: 20px;
  font-weight: 800;
  background-color: #f2f2f7;
  box-shadow: none !important;
  padding-left: 45px;
  padding-right: 45px;
  transition: all 0.2s;
}

:deep(.apple-giant-input .el-input-number__decrease),
:deep(.apple-giant-input .el-input-number__increase) {
  width: 40px;
  height: 40px;
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

.apple-pill-btn {
  height: 48px;
  border-radius: 24px;
  padding: 0 20px;
  font-weight: 600;
  background: #f2f2f7;
  border: none;
  color: #0071e3;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.apple-pill-btn:hover {
  background: #e5e5ea;
  transform: scale(1.02);
}

.disposal-actions {
  display: flex;
  gap: 15px;
  padding-top: 24px;
  border-top: 1px solid rgba(0,0,0,0.05);
}

.apple-sub-btn {
  flex: 1;
  height: 52px;
  border-radius: 14px;
  font-weight: 700;
  font-size: 15px;
  letter-spacing: 1px;
}

.apple-sub-btn.secondary {
  background: white;
  border: 1.5px solid #f2f2f7;
  color: #1d1d1f;
}

.apple-sub-btn.secondary:hover {
  background: #f2f2f7;
}

.apple-btn-small {
  flex: 1;
  border-radius: 10px;
  font-weight: 600;
}

.sync-option {
  padding-top: 15px;
  border-top: 1px dashed rgba(0,0,0,0.05);
}

.quantity-input-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.unit-text {
  color: #86868b;
  font-weight: 500;
}

.submit-guard-notice {
  background: #fff4e5;
  color: #ff9f0a;
  padding: 12px 16px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 20px;
}

.submit-section {
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid rgba(0,0,0,0.05);
}

.apple-btn-large {
  min-width: 160px;
  height: 48px;
  border-radius: 12px;
  font-weight: 700;
}

.animate-slide-up {
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(15px); }
  to { opacity: 1; transform: translateY(0); }
}

.text-warning { color: #ff9f0a; }

:deep(.el-form-item__label) {
  font-weight: 600 !important;
  color: #1d1d1f !important;
}

/* ========== 移动端响应式 ========== */
@media (max-width: 768px) {
  /* 警告盒子移动端优化 */
  .apple-alert-box {
    padding: 16px;
    border-radius: 16px;
  }
  
  .recovery-title {
    font-size: 16px;
  }
  
  /* 统计网格改为单列 */
  .recovery-stats-grid {
    flex-direction: column;
    gap: 12px;
    padding: 16px;
  }
  
  .stat-node {
    border-right: none;
    border-bottom: 1px solid #f2f2f7;
    padding-bottom: 12px;
  }
  
  .stat-node:last-child {
    border-bottom: none;
    padding-bottom: 0;
  }
  
  /* 输入框组改为单列 */
  .input-group-row {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .advanced-recovery-form {
    padding: 16px;
  }
  
  /* Input 适应手机 */
  .input-with-action {
    flex-direction: column;
    align-items: stretch;
  }
  
  .apple-giant-input {
    width: 100%;
  }
  
  .apple-pill-btn {
    width: 100%;
  }
  
  /* 处置按钮垂直布局 */
  .disposal-actions {
    flex-direction: column;
    gap: 12px;
  }
  
  .apple-sub-btn {
    width: 100%;
    font-size: 14px;
    letter-spacing: 0.5px;
  }
  
  /* 表单最小触摸尺寸 */
  :deep(.el-select),
  :deep(.el-input-number) {
    width: 100%;
  }
  
  :deep(.el-select .el-input__wrapper),
  :deep(.el-input-number .el-input__wrapper) {
    min-height: 48px;
  }
  
  /* 提交区域 */
  .submit-section :deep(.el-form-item__content) {
    display: flex !important;
    flex-direction: column;
    gap: 12px;
  }
  
  .apple-btn-large {
    width: 100%;
    min-width: 100%;
  }
  
  .submit-section .el-button {
    width: 100%;
  }
  
  /* 数量输入器 */
  .quantity-input-wrapper {
    flex-direction: column;
    align-items: stretch;
  }
  
  .quantity-input-wrapper :deep(.el-input-number) {
    width: 100%;
  }
}

/* 平板优化 */
@media (min-width: 768px) and (max-width: 1024px) {
  .input-group-row {
    grid-template-columns: 1fr;
  }
}
</style>

