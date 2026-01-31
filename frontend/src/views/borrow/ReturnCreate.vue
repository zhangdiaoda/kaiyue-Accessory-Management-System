<template>
  <div class="return-create">
    <el-card>
      <template #header>
        <span>登记归还</span>
      </template>

      <el-form :model="searchForm" :inline="true" class="search-form">
        <el-form-item>
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索员工或配件"
            clearable
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadData">查询</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" border>
        <el-table-column prop="record_no" label="记录编号" width="150" />
        <el-table-column prop="employee_name" label="员工" width="100" />
        <el-table-column prop="part_name" label="配件" width="150" />
        <el-table-column prop="borrow_quantity" label="领用数量" width="100" />
        <el-table-column label="已归还" width="100">
          <template #default="{ row }">
            {{ row.return_quantity + row.damaged_quantity }} / {{ row.borrow_quantity }}
          </template>
        </el-table-column>
        <el-table-column label="未归还" width="100">
          <template #default="{ row }">
            <el-tag type="warning">
              {{ row.borrow_quantity - row.return_quantity - row.damaged_quantity }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="borrow_time" label="领用时间" width="180" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              @click="handleReturn(row)"
              :disabled="row.status === 'RETURNED'"
            >
              {{ row.status === 'RETURNED' ? '已完成' : '归还' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="loadData"
        style="margin-top: 20px"
      />
    </el-card>

    <!-- 归还对话框 -->
    <el-dialog v-model="showDialog" title="登记归还" width="500px">
      <el-form :model="returnForm" label-width="120px">
        <el-form-item label="领用记录">
          <div>{{ currentRecord?.record_no }}</div>
        </el-form-item>
        <el-form-item label="员工">
          <div>{{ currentRecord?.employee_name }}</div>
        </el-form-item>
        <el-form-item label="配件">
          <div>{{ currentRecord?.part_name }}</div>
        </el-form-item>
        <el-form-item label="领用数量">
          <div>{{ currentRecord?.borrow_quantity }}</div>
        </el-form-item>
        <el-form-item label="未归还数量">
          <el-tag type="warning">
            {{ unreturned }}
          </el-tag>
        </el-form-item>
        <el-form-item label="正常归还数量" required>
          <el-input-number
            v-model="returnForm.return_quantity"
            :min="0"
            :max="unreturned"
            @change="handleReturnChange"
          />
        </el-form-item>
        <el-form-item label="损毁数量" required>
          <el-input-number
            v-model="returnForm.damaged_quantity"
            :min="0"
            :max="unreturned - returnForm.return_quantity"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getBorrowRecordList, returnBorrowRecord } from '@/api/borrow'

const tableData = ref([])
const showDialog = ref(false)
const currentRecord = ref(null)
const loading = ref(false)

const searchForm = reactive({
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const returnForm = reactive({
  return_quantity: 0,
  damaged_quantity: 0
})

const unreturned = computed(() => {
  if (!currentRecord.value) return 0
  return currentRecord.value.borrow_quantity - currentRecord.value.return_quantity - currentRecord.value.damaged_quantity
})

const loadData = async () => {
  try {
    const res = await getBorrowRecordList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: 'BORROWED,PARTIAL_RETURNED'
    })
    tableData.value = res.data.records || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载数据失败')
  }
}

const handleReturn = (row) => {
  currentRecord.value = row
  returnForm.return_quantity = unreturned.value
  returnForm.damaged_quantity = 0
  showDialog.value = true
}

const handleReturnChange = () => {
  const maxDamaged = unreturned.value - returnForm.return_quantity
  if (returnForm.damaged_quantity > maxDamaged) {
    returnForm.damaged_quantity = maxDamaged
  }
}

const handleSubmit = async () => {
  if (returnForm.return_quantity + returnForm.damaged_quantity === 0) {
    ElMessage.warning('归还数量和损毁数量不能都为0')
    return
  }

  loading.value = true
  try {
    await returnBorrowRecord(currentRecord.value.id, returnForm)
    ElMessage.success('归还登记成功')
    showDialog.value = false
    loadData()
  } catch (error) {
    ElMessage.error(error.message || '归还登记失败')
  } finally {
    loading.value = false
  }
}

loadData()
</script>

<style scoped>
.return-create {
  height: 100%;
}

.search-form {
  margin-bottom: 20px;
}
</style>
