<template>
  <div class="part-list">
    <!-- Header & Search Actions -->
    <div class="page-header">
      <div class="header-content">
        <h2 class="page-title">配件总库</h2>
        <span class="page-subtitle">管理所有精密部件与耗材库存</span>
      </div>
      <div class="action-group">
        <div class="search-box glass-effect">
          <el-icon class="search-icon"><Search /></el-icon>
          <input 
            v-model="searchForm.keyword" 
            placeholder="搜索配件名称 / 编号..." 
            class="clean-input"
            @keyup.enter="loadData" 
          />
          <button v-if="searchForm.keyword" class="clear-btn" @click="resetSearch">
            <el-icon><CircleCloseFilled /></el-icon>
          </button>
        </div>
        <el-button-group style="margin-right: 12px">
          <el-button size="large" class="apple-btn" @click="importDialog = true">
            <el-icon class="el-icon--left"><Upload /></el-icon> 导入
          </el-button>
          <el-button size="large" class="apple-btn" @click="handleExport">
            <el-icon class="el-icon--left"><Download /></el-icon> 导出
          </el-button>
        </el-button-group>
        <el-button type="primary" size="large" class="apple-primary-btn" @click="showDialog = true">
          <el-icon class="el-icon--left"><Plus /></el-icon>
          入库登记
        </el-button>
      </div>
    </div>

    <!-- Asset Grid -->
    <el-row :gutter="20" class="asset-grid">
      <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" v-for="item in tableData" :key="item.id">
        <div class="asset-card animate-fade-up">
          <!-- Status Badge -->
          <div class="status-badge" :class="getStockStatusClass(item)">
            {{ getStockStatusText(item) }}
          </div>
          
          <!-- Card Body -->
          <div class="card-main">
            <div class="asset-icon-placeholder" :style="{ background: getIconColor(item.id) }">
              <span class="icon-text">{{ item.name.charAt(0) }}</span>
            </div>
            <div class="asset-info">
              <div class="code-tag">{{ item.part_no }}</div>
              <h3 class="asset-name" :title="item.name">{{ item.name }}</h3>
              <div class="asset-spec">{{ item.specification || '标准规格' }}</div>
            </div>
          </div>

          <!-- Inventory Visual -->
          <div class="inventory-section">
            <div class="inv-labels">
              <span class="label">当前库存</span>
              <span class="val">{{ item.stock_quantity }} <small>{{ item.unit }}</small></span>
            </div>
            <el-progress 
              :percentage="Math.min(100, (item.stock_quantity / (item.warning_threshold * 3)) * 100)" 
              :status="item.stock_quantity < item.warning_threshold ? 'exception' : 'success'"
              :stroke-width="6"
              :show-text="false"
              class="inv-progress"
            />
          </div>

          <!-- Footer Actions -->
          <div class="card-footer">
            <div class="price-tag">¥ {{ item.price || '-' }}</div>
            <div class="actions">
              <el-button link class="icon-btn" @click="handleRestock(item)" title="补货入库">
                <el-icon><Box /></el-icon>
              </el-button>
              <el-button link class="icon-btn" @click="handleEdit(item)">
                <el-icon><Edit /></el-icon>
              </el-button>
              <el-button link class="icon-btn danger" @click="handleDelete(item)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Pagination -->
    <div class="pagination-bar glass-effect">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[12, 24, 48, 96]"
        layout="prev, pager, next"
        @size-change="loadData"
        @current-change="loadData"
        background
      />
      <span class="total-text">共 {{ pagination.total }} 项资产</span>
    </div>

    <!-- 新增/编辑对话框 (保持原逻辑) -->
    <el-dialog
      v-model="showDialog"
      :title="dialogTitle"
      width="600px"
      class="apple-dialog"
      append-to-body
    >
      <el-form :model="formData" label-position="top" class="apple-form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="配件编号" required>
              <el-input v-model="formData.part_no" placeholder="如: P-2024-001" class="apple-input" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="配件名称" required>
              <el-input v-model="formData.name" class="apple-input" />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="规格型号">
              <el-input v-model="formData.specification" class="apple-input" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="计量单位">
              <el-input v-model="formData.unit" placeholder="件/套/个" class="apple-input" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="初始库存" required>
              <el-input-number v-model="formData.stock_quantity" :min="0" controls-position="right" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="预警阈值" required>
              <el-input-number v-model="formData.warning_threshold" :min="0" controls-position="right" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="参考单价">
              <el-input-number v-model="formData.price" :min="0" :precision="2" controls-position="right" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="备注说明">
          <el-input v-model="formData.remark" type="textarea" :rows="3" class="apple-input" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showDialog = false" round>取消</el-button>
          <el-button type="primary" @click="handleSubmit" round class="apple-primary-btn">确认入库</el-button>
        </div>
      </template>
    </el-dialog>
    <!-- 导入弹窗 -->
    <el-dialog
      v-model="importDialog"
      title="批量导入配件"
      width="500px"
      class="apple-dialog"
      append-to-body
    >
      <div class="import-content">
        <div class="step-box">
          <div class="step-title">第一步：下载模版</div>
          <el-button @click="downloadTemplate" icon="Document">下载标准Excel模版</el-button>
        </div>
        <div class="step-box">
          <div class="step-title">第二步：上传文件</div>
          <el-upload
            ref="uploadRef"
            drag
            action="#"
            :auto-upload="false"
            :limit="1"
            :on-change="handleFileChange"
            accept=".xlsx"
          >
            <el-icon class="el-icon--upload"><Upload /></el-icon>
            <div class="el-upload__text">
              拖拽文件到此处或 <em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                仅支持 .xlsx 文件，请确保表头与模版一致
              </div>
            </template>
          </el-upload>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="importDialog = false" round>取消</el-button>
          <el-button type="primary" :loading="importLoading" @click="submitImport" round class="apple-primary-btn">
            开始导入
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 补货弹窗 -->
    <el-dialog
      v-model="restockDialog"
      title="库存补货与入库"
      width="500px"
      class="apple-dialog"
      append-to-body
    >
      <el-form :model="restockForm" label-position="top" class="apple-form">
        <el-alert
          :title="`正在为 [${restockForm.part_name}] 补充库存`"
          type="info"
          show-icon
          :closable="false"
          style="margin-bottom: 20px"
        />
        <el-form-item label="入库数量" required>
          <el-input-number 
            v-model="restockForm.quantity" 
            :min="1" 
            controls-position="right" 
            style="width: 100%" 
            size="large"
          />
        </el-form-item>
        <el-form-item label="批次号 / 采购单号">
          <el-input v-model="restockForm.batch_no" placeholder="自动生成，可修改" class="apple-input" />
        </el-form-item>
        <el-form-item label="备注说明">
          <el-input 
            v-model="restockForm.remark" 
            type="textarea" 
            :rows="3" 
            placeholder="例如：供应商、采购来源等"
            class="apple-input" 
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="restockDialog = false" round>取消</el-button>
          <el-button type="primary" :loading="restockLoading" @click="submitRestock" round class="apple-primary-btn">
            确认入库
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Edit, Delete, CircleCloseFilled, Download, Upload, Document, Box } from '@element-plus/icons-vue'
import { getPartList, createPart, updatePart, deletePart, getTemplateUrl, getExportUrl, importParts } from '@/api/part'
import { restockPart } from '@/api/inbound'
import dayjs from 'dayjs'

const tableData = ref([])
const showDialog = ref(false)
const dialogTitle = ref('添加配件')

const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 12, // Grid 布局每页显示更多
  total: 0
})

const formData = reactive({
  id: null,
  part_no: '',
  name: '',
  category_id: 1,
  specification: '',
  unit: '件',
  stock_quantity: 0,
  warning_threshold: 10,
  price: 0,
  remark: ''
})

// 随机色板用于图标背景
const iconColors = ['#0071e3', '#ff9f0a', '#30b0c7', '#ff375f', '#bf5af2', '#34c759']
const getIconColor = (id) => iconColors[id % iconColors.length]

const getStockStatusClass = (item) => {
  if (item.stock_quantity === 0) return 'status-out'
  if (item.stock_quantity < item.warning_threshold) return 'status-low'
  return 'status-ok'
}

const getStockStatusText = (item) => {
  if (item.stock_quantity === 0) return '缺货'
  if (item.stock_quantity < item.warning_threshold) return '低库存'
  return '充足'
}

const loadData = async () => {
  try {
    const res = await getPartList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchForm.keyword
    })
    tableData.value = res.data.records || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载列表失败')
  }
}

const resetSearch = () => {
  searchForm.keyword = ''
  pagination.page = 1
  loadData()
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑配件信息'
  Object.assign(formData, row)
  showDialog.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要移除配件 "${row.name}" 吗？`, '删除确认', {
      type: 'warning',
      confirmButtonText: '移除',
      cancelButtonText: '取消',
      confirmButtonClass: 'el-button--danger'
    })
    await deletePart(row.id)
    ElMessage.success('已移除')
    loadData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('移除失败')
  }
}

const handleSubmit = async () => {
  try {
    if (formData.id) {
      await updatePart(formData.id, formData)
      ElMessage.success('更新成功')
    } else {
      await createPart(formData)
      ElMessage.success('入库成功')
    }
    showDialog.value = false
    loadData()
    resetForm()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const resetForm = () => {
  formData.id = null
  formData.part_no = ''
  formData.name = ''
  formData.specification = ''
  formData.unit = '件'
  formData.stock_quantity = 0
  formData.warning_threshold = 10
  formData.price = 0
  formData.remark = ''
}

// 导入导出逻辑
const importDialog = ref(false)
const importFile = ref(null)
const importLoading = ref(false)
const uploadRef = ref(null)

const handleExport = () => {
  window.open(getExportUrl(), '_blank')
}

const downloadTemplate = () => {
  window.open(getTemplateUrl(), '_blank')
}

const handleFileChange = (file) => {
  importFile.value = file.raw
}

const submitImport = async () => {
  if (!importFile.value) {
    ElMessage.warning('请选择需要导入的Excel文件')
    return
  }
  
  importLoading.value = true
  try {
    const fd = new FormData()
    fd.append('file', importFile.value)
    
    const res = await importParts(fd)
    ElMessage.success(res.message || '导入成功')
    importDialog.value = false
    importFile.value = null
    if (uploadRef.value) uploadRef.value.clearFiles()
    loadData()
  } catch (err) {
    ElMessage.error('导入失败，请检查文件格式')
  } finally {
    importLoading.value = false
  }
}


// 补货逻辑
const restockDialog = ref(false)
const restockLoading = ref(false)
const restockForm = reactive({
  part_id: null,
  part_name: '',
  quantity: 1,
  batch_no: '',
  remark: ''
})

const handleRestock = (row) => {
  restockForm.part_id = row.id
  restockForm.part_name = row.name
  restockForm.quantity = 1
  restockForm.batch_no = 'B' + dayjs().format('YYYYMMDDHHmmss')
  restockForm.remark = ''
  restockDialog.value = true
}

const submitRestock = async () => {
  if (restockForm.quantity <= 0) {
    ElMessage.warning('补货数量必须大于0')
    return
  }
  
  restockLoading.value = true
  try {
    await restockPart({
      part_id: restockForm.part_id,
      quantity: restockForm.quantity,
      batch_no: restockForm.batch_no,
      remark: restockForm.remark
    })
    ElMessage.success('补货成功')
    restockDialog.value = false
    loadData()
  } catch (err) {
    ElMessage.error('操作失败')
  } finally {
    restockLoading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.part-list {
  padding-bottom: 80px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 30px;
  flex-wrap: wrap;
  gap: 20px;
}

.page-title {
  font-size: 28px;
  font-weight: 800;
  color: #1d1d1f;
  margin: 0 0 4px 0;
  letter-spacing: -0.5px;
}

.page-subtitle {
  font-size: 14px;
  color: #86868b;
}

.action-group {
  display: flex;
  gap: 16px;
  align-items: center;
}

.search-box {
  display: flex;
  align-items: center;
  background: white;
  border-radius: 12px;
  padding: 8px 12px;
  width: 280px;
  transition: all 0.2s;
  border: 1px solid rgba(0,0,0,0.03);
}

.search-box:focus-within {
  width: 320px;
  box-shadow: 0 0 0 2px rgba(0,113,227,0.2);
}

.clean-input {
  border: none;
  background: transparent;
  outline: none;
  font-size: 14px;
  color: #1d1d1f;
  flex: 1;
  margin-left: 8px;
}

.search-icon { color: #86868b; font-size: 16px; }
.clear-btn {
  background: none;
  border: none;
  color: #86868b;
  cursor: pointer;
  padding: 0;
  display: flex;
}
.clear-btn:hover { color: #1d1d1f; }

.asset-card {
  background: white;
  border-radius: 20px;
  padding: 20px;
  margin-bottom: 20px;
  position: relative;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(0,0,0,0.03);
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
}

.asset-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 16px 32px rgba(0,0,0,0.08);
  z-index: 2;
}

.status-badge {
  position: absolute;
  top: 15px;
  right: 15px;
  font-size: 10px;
  font-weight: 700;
  padding: 4px 10px;
  border-radius: 12px;
}
.status-ok { background: #e8f5e9; color: #2e7d32; }
.status-low { background: #fff3e0; color: #ef6c00; }
.status-out { background: #ffebee; color: #c62828; }

.card-main {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.asset-icon-placeholder {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-shrink: 0;
}

.icon-text {
  font-size: 24px;
  font-weight: 800;
  color: white;
  text-transform: uppercase;
}

.asset-info {
  flex: 1;
  overflow: hidden;
}

.code-tag {
  font-size: 10px;
  color: #86868b;
  font-family: monospace;
  background: #f5f5f7;
  padding: 2px 6px;
  border-radius: 4px;
  display: inline-block;
  margin-bottom: 4px;
}

.asset-name {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 700;
  color: #1d1d1f;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.asset-spec {
  font-size: 12px;
  color: #86868b;
}

.inventory-section {
  background: #f5f5f7;
  border-radius: 12px;
  padding: 12px;
  margin-bottom: 16px;
}

.inv-labels {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  margin-bottom: 6px;
}

.inv-labels .val {
  font-weight: 700;
  color: #1d1d1f;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px dashed rgba(0,0,0,0.05);
  padding-top: 12px;
}

.price-tag {
  font-size: 14px;
  font-weight: 700;
  color: #1d1d1f;
}

.actions {
  display: flex;
  gap: 8px;
}

.icon-btn {
  padding: 6px;
  height: auto;
  color: #86868b;
}

.icon-btn:hover { color: #0071e3; background: #f2f2f7; border-radius: 6px; }
.icon-btn.danger:hover { color: #ff375f; background: #fff1f2; }

.pagination-bar {
  position: fixed;
  bottom: 0;
  right: 0;
  left: 260px; /* Sidebar width */
  background: rgba(255,255,255,0.8);
  backdrop-filter: blur(20px);
  padding: 16px 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid rgba(0,0,0,0.05);
  z-index: 10;
}

.total-text {
  font-size: 13px;
  color: #86868b;
}

/* Animations */
.animate-fade-up {
  animation: fadeUp 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  opacity: 0;
  transform: translateY(10px);
}

@keyframes fadeUp {
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 768px) {
  .pagination-bar { left: 0; padding: 16px 20px; }
  .action-group { width: 100%; }
  .search-box { width: 100%; }
}
</style>

