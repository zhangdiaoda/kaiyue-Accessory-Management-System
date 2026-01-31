<template>
  <div class="inbound-history">
    <div class="page-header">
      <h2 class="page-title">入库明细流水</h2>
      <div class="search-box glass-effect">
        <el-icon class="search-icon"><Search /></el-icon>
        <input 
          v-model="keyword" 
          placeholder="搜索配件名称 / 编号..." 
          class="clean-input"
          @keyup.enter="loadData" 
        />
        <button v-if="keyword" class="clear-btn" @click="resetSearch">
          <el-icon><CircleCloseFilled /></el-icon>
        </button>
      </div>
    </div>

    <el-card class="apple-card">
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="created_at" label="入库时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="配件信息" min-width="200">
          <template #default="{ row }">
            <div class="part-info">
              <span class="part-no">{{ row.part.part_no }}</span>
              <span class="part-name">{{ row.part.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="入库数量" width="120">
          <template #default="{ row }">
            <span class="qty-tag">+{{ row.quantity }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="operator" label="操作人" width="120" />
        <el-table-column prop="batch_no" label="批次/单号" width="150" show-overflow-tooltip />
        <el-table-column prop="remark" label="备注" min-width="200" show-overflow-tooltip />
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadData"
          @current-change="loadData"
          background
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Search, CircleCloseFilled } from '@element-plus/icons-vue'
import { getInboundList } from '@/api/inbound'
import dayjs from 'dayjs'

const loading = ref(false)
const tableData = ref([])
const keyword = ref('')

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await getInboundList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      part_name: keyword.value
    })
    tableData.value = res.data.records || []
    pagination.total = res.data.total || 0
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  keyword.value = ''
  pagination.page = 1
  loadData()
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.inbound-history {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1d1d1f;
  margin: 0;
}

.search-box {
  display: flex;
  align-items: center;
  background: white;
  border-radius: 12px;
  padding: 8px 12px;
  width: 280px;
  border: 1px solid rgba(0,0,0,0.05);
  transition: all 0.2s;
}

.search-box:focus-within {
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

.search-icon { color: #86868b; }
.clear-btn { border: none; background: none; color: #86868b; cursor: pointer; padding: 0; display: flex; }

.apple-card {
  border-radius: 16px;
  border: 1px solid rgba(0,0,0,0.05);
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.part-info {
  display: flex;
  flex-direction: column;
}

.part-no {
  font-size: 12px;
  color: #86868b;
  font-family: monospace;
}

.part-name {
  font-weight: 500;
  color: #1d1d1f;
}

.qty-tag {
  color: #34c759;
  font-weight: 600;
  background: rgba(52, 199, 89, 0.1);
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 13px;
}
</style>
