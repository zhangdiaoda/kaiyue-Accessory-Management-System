<template>
  <div class="employee-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>员工管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加员工
          </el-button>
          <div class="header-actions">
            <el-upload
              class="upload-demo"
              action="#"
              :show-file-list="false"
              :auto-upload="false"
              :on-change="handleImport"
              accept=".xlsx,.xls"
            >
              <el-button type="success" plain>
                <el-icon><Upload /></el-icon>
                批量导入
              </el-button>
            </el-upload>
            <el-button type="info" plain @click="downloadTemplate">
              <el-icon><Download /></el-icon>
              下载模板
            </el-button>
            <el-button type="warning" plain @click="handleExport">
              <el-icon><Download /></el-icon>
              导出员工
            </el-button>
          </div>
        </div>
      </template>

      <el-form :inline="true" class="search-form">
        <el-form-item>
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索员工姓名或工号"
            clearable
            @keyup.enter="loadData"
          />
        </el-form-item>
        <el-form-item>
          <el-select v-model="searchForm.status" placeholder="全部状态" clearable>
            <el-option label="在职" :value="1" />
            <el-option label="离职" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 移动端卡片布局 -->
      <div v-if="isMobile" class="employee-cards">
        <div v-for="item in tableData" :key="item.id" class="employee-card">
          <div class="card-header">
            <div class="employee-info">
              <div class="employee-name">{{ item.name }}</div>
              <div class="employee-no">{{ item.employee_no }}</div>
            </div>
            <el-tag :type="item.status === 1 ? 'success' : 'info'" size="small">
              {{ item.status === 1 ? '在职' : '离职' }}
            </el-tag>
          </div>
          
          <div class="card-body">
            <div class="card-row" v-if="item.department">
              <span class="label">部门</span>
              <span class="value">{{ item.department }}</span>
            </div>
            <div class="card-row" v-if="item.position">
              <span class="label">岗位</span>
              <span class="value">{{ item.position }}</span>
            </div>
            <div class="card-row" v-if="item.phone">
              <span class="label">手机号</span>
              <span class="value">{{ item.phone }}</span>
            </div>
            <div class="card-row">
              <span class="label">创建时间</span>
              <span class="value">{{ item.created_at }}</span>
            </div>
          </div>
          
          <div class="card-footer">
            <el-button type="primary" size="small" @click="handleEdit(item)" style="flex: 1;">
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(item)" plain style="flex: 1;">
              删除
            </el-button>
          </div>
        </div>
        
        <el-empty v-if="tableData.length === 0" description="暂无员工数据" :image-size="80" />
      </div>

      <!-- 桌面端表格布局 -->
      <el-table v-else :data="tableData" border stripe>
        <el-table-column prop="employee_no" label="工号" width="120" />
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="department" label="部门" width="150" />
        <el-table-column prop="position" label="岗位" width="150" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '在职' : '离职' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadData"
        @current-change="loadData"
        style="margin-top: 20px"
      />
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="showDialog" :title="dialogTitle" width="600px">
      <el-form :model="formData" label-width="100px">
        <el-form-item label="员工工号" required>
          <el-input v-model="formData.employee_no" placeholder="如: EMP001" />
        </el-form-item>
        <el-form-item label="姓名" required>
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item label="部门">
          <el-input v-model="formData.department" placeholder="如: 机加工车间" />
        </el-form-item>
        <el-form-item label="岗位">
          <el-input v-model="formData.position" placeholder="如: 操作工" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="formData.phone" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio :label="1">在职</el-radio>
            <el-radio :label="0">离职</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Upload, Download } from '@element-plus/icons-vue'
import { useResponsive } from '@/composables/useResponsive'
import { getEmployeeList, createEmployee, updateEmployee, deleteEmployee, importEmployees, downloadTemplateUrl, exportEmployeesUrl } from '@/api/employee'
import { useUserStore } from '@/store/user'
const userStore = useUserStore()
const { isMobile } = useResponsive()

const tableData = ref([])
const showDialog = ref(false)
const dialogTitle = ref('添加员工')

const searchForm = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive({
  id: null,
  employee_no: '',
  name: '',
  department: '',
  position: '',
  phone: '',
  status: 1
})

const loadData = async () => {
  try {
    const res = await getEmployeeList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchForm.keyword,
      status: searchForm.status
    })
    tableData.value = res.data.records || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载数据失败')
  }
}

const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  pagination.page = 1
  loadData()
}

const handleAdd = () => {
  dialogTitle.value = '添加员工'
  resetForm()
  showDialog.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑员工'
  Object.assign(formData, row)
  showDialog.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这个员工吗？', '提示', { type: 'warning' })
    await deleteEmployee(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

const downloadTemplate = () => {
  window.open(`${import.meta.env.VITE_API_URL || ''}${downloadTemplateUrl}?token=${userStore.token}`)
}

const handleExport = () => {
  window.open(`${import.meta.env.VITE_API_URL || ''}${exportEmployeesUrl}?token=${userStore.token}`)
}

const handleImport = async (file) => {
  const formData = new FormData()
  formData.append('file', file.raw)
  try {
    const res = await importEmployees(formData)
    if (res.code === 200) {
      ElMessage.success(res.message)
      loadData()
    } else {
      ElMessage.error(res.message)
    }
  } catch (error) {
    ElMessage.error('导入失败')
  }
}

const handleSubmit = async () => {
  try {
    if (formData.id) {
      await updateEmployee(formData.id, formData)
      ElMessage.success('更新成功')
    } else {
      await createEmployee(formData)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    loadData()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const resetForm = () => {
  formData.id = null
  formData.employee_no = ''
  formData.name = ''
  formData.department = ''
  formData.position = ''
  formData.phone = ''
  formData.status = 1
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.employee-manage {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.search-form {
  margin-bottom: 20px;
}

/* 员工卡片样式 */
.employee-cards {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.employee-card {
  background: white;
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  border: 1px solid rgba(0,0,0,0.05);
  transition: all 0.2s;
}

.employee-card:active {
  transform: scale(0.98);
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
}

.employee-card .card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(0,0,0,0.05);
}

.employee-card .employee-info {
  flex: 1;
}

.employee-card .employee-name {
  font-size: 16px;
  font-weight: 700;
  color: #1d1d1f;
  margin-bottom: 4px;
}

.employee-card .employee-no {
  font-size: 12px;
  color: #86868b;
  font-family: monospace;
}

.employee-card .card-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 12px;
}

.employee-card .card-row {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.employee-card .card-row .label {
  color: #86868b;
  font-size: 13px;
}

.employee-card .card-row .value {
  color: #1d1d1f;
  font-weight: 600;
  text-align: right;
  flex: 1;
  margin-left: 12px;
}

.employee-card .card-footer {
  display: flex;
  gap: 12px;
  padding-top: 12px;
  border-top: 1px solid rgba(0,0,0,0.05);
}

/* 移动端响应式 */
@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .card-header > span {
    margin-bottom: 8px;
  }
  
  .card-header :deep(.el-button) {
    width: 100%;
    margin-bottom: 8px;
  }
  
  .header-actions {
    flex-direction: column;
    width: 100%;
  }
  
  .header-actions :deep(.el-button),
  .header-actions :deep(.el-upload) {
    width: 100%;
  }
  
  .header-actions :deep(.upload-demo) {
    width: 100%;
  }
  
  .header-actions :deep(.upload-demo .el-button) {
    width: 100%;
  }
  
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
  
  .search-form :deep(.el-input),
  .search-form :deep(.el-select) {
    width: 100%;
  }
  
  .search-form :deep(.el-form-item:last-child) {
    display: flex;
    gap: 8px;
  }
  
  .search-form :deep(.el-form-item:last-child .el-form-item__content) {
    display: flex;
    gap: 8px;
    width: 100%;
  }
  
  .search-form :deep(.el-form-item:last-child .el-button) {
    flex: 1;
  }
}
</style>
