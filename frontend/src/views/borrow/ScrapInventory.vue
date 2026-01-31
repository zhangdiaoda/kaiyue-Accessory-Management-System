<template>
  <div class="inventory-view">
    <div class="header-section">
      <h2 class="page-title">🗑️ 废品管理中心</h2>
      <el-tag type="danger" effect="dark" class="count-badge">历史已累计报废 {{ totalCount }} 项记录</el-tag>
    </div>

    <el-tabs v-model="activeTab" class="apple-tabs">
      <el-tab-pane label="卡片面板" name="card">
        <el-row :gutter="20" v-loading="loading">
          <el-col 
            v-for="item in inventory" 
            :key="item.id" 
            :xs="24" :sm="12" :md="8" :lg="6"
            class="card-col"
          >
            <el-card shadow="never" class="inventory-card scrap-style glass-effect">
              <div class="card-main">
                <div class="part-info">
                  <div class="part-name">{{ item.part_name }}</div>
                  <div class="part-no text-muted">{{ item.part_no }}</div>
                </div>
                <div class="stock-value">
                  <span class="num">{{ item.quantity }}</span>
                  <small>件</small>
                </div>
              </div>
              <el-divider />
              <div class="trace-info">
                <div class="trace-row">
                  <span class="label">责任员工:</span>
                  <span class="value">{{ item.employee_name }}</span>
                </div>
                <div class="trace-row">
                  <span class="label">归属部门:</span>
                  <span class="value">{{ item.department }}</span>
                </div>
                <div class="trace-row">
                  <span class="label">报废时间:</span>
                  <span class="value">{{ new Date(item.updated_at).toLocaleDateString() }}</span>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-empty v-if="!loading && inventory.length === 0" description="废品库目前为空" />
        </el-row>
      </el-tab-pane>

      <el-tab-pane label="详尽清单" name="table">
        <el-table :data="inventory" border stripe v-loading="loading" class="apple-table">
          <el-table-column prop="part_name" label="废品配件" />
          <el-table-column prop="part_no" label="规格编号" width="120" />
          <el-table-column prop="quantity" label="报废数量" width="100" align="center">
            <template #default="{ row }">
              <span class="text-danger font-bold">{{ row.quantity }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="employee_name" label="主要责任人" width="120" />
          <el-table-column prop="department" label="责任部门" width="150" />
          <el-table-column prop="updated_at" label="报废处理日" width="180">
            <template #default="{ row }">
              {{ new Date(row.updated_at).toLocaleString() }}
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getScrapInventory } from '@/api/borrow'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const activeTab = ref('card')
const inventory = ref([])

const totalCount = computed(() => inventory.value.reduce((sum, item) => sum + item.quantity, 0))

const loadData = async () => {
  loading.value = true
  try {
    const res = await getScrapInventory()
    inventory.value = res.data || []
  } catch (error) {
    ElMessage.error('加载废品库失败')
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.inventory-view { min-height: 100%; }

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

.card-col { margin-bottom: 20px; }

.inventory-card {
  border-radius: 18px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(0,0,0,0.05);
}

.scrap-style:hover {
  transform: translateY(-5px);
  border-color: rgba(255, 55, 95, 0.2);
  box-shadow: 0 10px 20px rgba(255, 55, 95, 0.05);
}

.card-main {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.part-name {
  font-weight: 700;
  font-size: 16px;
  color: #1d1d1f;
  margin-bottom: 4px;
}

.part-no { font-size: 12px; font-family: monospace; }

.stock-value { text-align: right; }
.stock-value .num { font-size: 24px; font-weight: 800; color: #ff375f; }
.stock-value small { margin-left: 2px; color: #86868b; }

.trace-info { font-size: 13px; }
.trace-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 6px;
}

.trace-row .label { color: #86868b; }
.trace-row .value { color: #1d1d1f; font-weight: 500; }

.font-bold { font-weight: 700; }
.text-danger { color: #ff375f; }
</style>
