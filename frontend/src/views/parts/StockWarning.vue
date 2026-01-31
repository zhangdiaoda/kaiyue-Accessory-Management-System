<template>
  <div class="stock-warning">
    <div class="header-section">
      <h2 class="page-title">🚨 库存预警中心</h2>
      <el-tag v-if="warningParts.length > 0" type="danger" effect="dark" class="count-badge">
        {{ warningParts.length }} 项需关注
      </el-tag>
    </div>

    <el-alert
      v-if="warningParts.length === 0"
      title="库存状态极佳"
      type="success"
      description="目前所有配件库存均处于安全范围内，无需补充。"
      :closable="false"
      show-icon
      class="glass-effect alert-success"
    />

    <!-- 卡片栅格 -->
    <el-row :gutter="20" v-else class="warning-grid">
      <el-col 
        v-for="part in warningParts" 
        :key="part.id" 
        :xs="24" :sm="12" :md="8" :lg="6"
        class="card-col"
      >
        <el-card shadow="never" class="warning-card glass-effect animate-slide-up">
          <div class="card-body">
            <div class="part-header">
              <div class="icon-box" :class="urgencyClass(part)">
                <el-icon :size="24"><Warning /></el-icon>
              </div>
              <div class="part-info">
                <div class="part-name">{{ part.name }}</div>
                <div class="part-no">{{ part.part_no }}</div>
              </div>
            </div>

            <div class="stock-stats">
              <div class="stat-item">
                <div class="stat-label">当前库存</div>
                <div class="stat-value text-danger">{{ part.stock_quantity }} <small>{{ part.unit }}</small></div>
              </div>
              <div class="stat-divider"></div>
              <div class="stat-item">
                <div class="stat-label">预警阈值</div>
                <div class="stat-value text-muted">{{ part.warning_threshold }}</div>
              </div>
            </div>

            <div class="progress-container">
              <div class="progress-info">
                <span>缺口：{{ part.warning_threshold - part.stock_quantity }} {{ part.unit }}</span>
                <span>{{ Math.round((part.stock_quantity / part.warning_threshold) * 100) }}%</span>
              </div>
              <el-progress 
                :percentage="Math.min(100, Math.round((part.stock_quantity / part.warning_threshold) * 100))" 
                :status="part.stock_quantity === 0 ? 'exception' : 'warning'"
                :stroke-width="8"
                :show-text="false"
              />
            </div>

            <div class="card-actions">
              <el-button type="primary" class="apple-btn full-width" @click="handleEdit(part)">
                立即补充库存
              </el-button>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 补充库存对话框 -->
    <el-dialog v-model="showDialog" title="补给物资" width="450px" class="apple-dialog">
      <div class="supply-dialog-content">
        <div class="supply-header">
          <div class="supply-item-name">{{ currentPart?.name }}</div>
          <div class="supply-item-spec">{{ currentPart?.specification || '标准规格' }}</div>
        </div>
        
        <el-form :model="formData" label-position="top" class="apple-form">
          <el-form-item label="本次补充数量" required>
            <el-input-number 
              v-model="formData.add_quantity" 
              :min="1" 
              style="width: 100%"
              controls-position="right"
            />
          </el-form-item>
          
          <div class="supply-summary">
            <div class="summary-line">
              <span>当前结余</span>
              <span>{{ currentPart?.stock_quantity }}</span>
            </div>
            <div class="summary-line">
              <span>新增入库</span>
              <span class="text-primary">+{{ formData.add_quantity }}</span>
            </div>
            <el-divider />
            <div class="summary-line total">
              <span>入库后预计总额</span>
              <span class="text-success">{{ (currentPart?.stock_quantity || 0) + formData.add_quantity }}</span>
            </div>
          </div>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">确认入库</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Warning } from '@element-plus/icons-vue'
import { getLowStockParts, updatePart } from '@/api/part'

const warningParts = ref([])
const showDialog = ref(false)
const currentPart = ref(null)
const loading = ref(false)

const formData = reactive({
  add_quantity: 1
})

const urgencyClass = (part) => {
  if (part.stock_quantity === 0) return 'is-critical'
  if (part.stock_quantity < part.warning_threshold * 0.3) return 'is-high'
  return 'is-medium'
}

const loadWarningParts = async () => {
  try {
    const res = await getLowStockParts()
    warningParts.value = res.data || []
  } catch (error) {
    ElMessage.error('加载预警数据失败')
  }
}

const handleEdit = (row) => {
  currentPart.value = row
  formData.add_quantity = Math.max(1, row.warning_threshold - row.stock_quantity)
  showDialog.value = true
}

const handleSubmit = async () => {
  loading.value = true
  try {
    const newQuantity = currentPart.value.stock_quantity + formData.add_quantity
    await updatePart(currentPart.value.id, {
      ...currentPart.value,
      stock_quantity: newQuantity
    })
    ElMessage.success('库存补充成功')
    showDialog.value = false
    loadWarningParts()
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadWarningParts()
})
</script>

<style scoped>
.stock-warning {
  min-height: 100%;
}

.header-section {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 800;
  margin: 0;
  color: #1d1d1f;
}

.alert-success {
  border-radius: 16px;
  padding: 20px;
}

.warning-grid {
  margin-top: 10px;
}

.card-col {
  margin-bottom: 20px;
}

.warning-card {
  border-radius: 20px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(0,0,0,0.05);
}

.warning-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 24px rgba(0,0,0,0.08);
}

.card-body {
  padding: 4px;
}

.part-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.icon-box {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.is-critical { background-color: #ff375f15; color: #ff375f; }
.is-high { background-color: #ff9f0a15; color: #ff9f0a; }
.is-medium { background-color: #ffcc0015; color: #ffcc00; }

.part-name {
  font-weight: 700;
  color: #1d1d1f;
  font-size: 16px;
  margin-bottom: 2px;
}

.part-no {
  font-size: 12px;
  color: #86868b;
  font-family: monospace;
}

.stock-stats {
  display: flex;
  justify-content: space-around;
  align-items: center;
  background: rgba(0,0,0,0.02);
  padding: 12px;
  border-radius: 12px;
  margin-bottom: 20px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  font-size: 11px;
  color: #86868b;
  margin-bottom: 4px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.stat-value {
  font-size: 18px;
  font-weight: 800;
}

.stat-value small {
  font-size: 12px;
  font-weight: 500;
}

.stat-divider {
  width: 1px;
  height: 24px;
  background: rgba(0,0,0,0.05);
}

.progress-container {
  margin-bottom: 24px;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #86868b;
  margin-bottom: 8px;
}

.full-width {
  width: 100%;
}

/* 对话框细节 */
.supply-header {
  margin-bottom: 24px;
  text-align: center;
}

.supply-item-name {
  font-size: 20px;
  font-weight: 800;
  color: #1d1d1f;
}

.supply-item-spec {
  color: #86868b;
  font-size: 13px;
  margin-top: 4px;
}

.supply-summary {
  background: #fbfbfd;
  padding: 20px;
  border-radius: 16px;
  margin-top: 20px;
}

.summary-line {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  color: #1d1d1f;
  margin-bottom: 8px;
}

.summary-line.total {
  font-weight: 700;
  font-size: 16px;
}

.animate-slide-up {
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.text-danger { color: #ff375f; }
.text-primary { color: #0071e3; }
.text-success { color: #32d74b; }
.text-muted { color: #86868b; }
</style>
