<template>
  <div class="inventory-view">
    <div class="header-section">
      <h2 class="page-title">📦 旧件回收仓库</h2>
      <el-tag type="info" effect="dark" class="count-badge">当前已回收 {{ totalCount }} 项旧件</el-tag>
    </div>

    <el-tabs v-model="activeTab" class="apple-tabs">
      <el-tab-pane label="卡片视图" name="card">
        <el-row :gutter="20" v-loading="loading">
          <el-col 
            v-for="item in inventory" 
            :key="item.id" 
            :xs="24" :sm="12" :md="8" :lg="6"
            class="card-col"
          >
            <el-card shadow="never" class="inventory-card glass-effect">
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
                  <span class="label">回收来源:</span>
                  <span class="value">{{ item.employee_name }}</span>
                </div>
                <div class="trace-row">
                  <span class="label">所属部门:</span>
                  <span class="value">{{ item.department }}</span>
                </div>
                <div class="trace-row">
                  <span class="label">更新时间:</span>
                  <span class="value">{{ new Date(item.updated_at).toLocaleDateString() }}</span>
                </div>
              </div>
            </el-card>
          </el-col>
          <el-empty v-if="!loading && inventory.length === 0" description="旧件库暂无积压" />
        </el-row>
      </el-tab-pane>

      <el-tab-pane label="列表视图" name="table">
        <el-table :data="inventory" border stripe v-loading="loading" class="apple-table">
          <el-table-column prop="part_name" label="配件名称" />
          <el-table-column prop="part_no" label="编号" width="120" />
          <el-table-column prop="quantity" label="待处置数量" width="100" align="center" />
          <el-table-column prop="employee_name" label="回收来源" width="120" />
          <el-table-column prop="department" label="部门" width="150" />
          <el-table-column prop="updated_at" label="最后操作" width="180">
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
import { getOldInventory } from '@/api/borrow'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const activeTab = ref('card')
const inventory = ref([])

const totalCount = computed(() => inventory.value.reduce((sum, item) => sum + item.quantity, 0))

const loadData = async () => {
  loading.value = true
  try {
    const res = await getOldInventory()
    inventory.value = res.data || []
  } catch (error) {
    ElMessage.error('加载旧件库失败')
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
  transition: transform 0.3s;
  border: 1px solid rgba(0,0,0,0.05);
}

.inventory-card:hover { transform: translateY(-5px); }

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
.stock-value .num { font-size: 24px; font-weight: 800; color: #0071e3; }
.stock-value small { margin-left: 2px; color: #86868b; }

.trace-info { font-size: 13px; }
.trace-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 6px;
}

.trace-row .label { color: #86868b; }
.trace-row .value { color: #1d1d1f; font-weight: 500; }

.animate-slide-up { animation: slideUp 0.4s ease-out; }
@keyframes slideUp {
  from { opacity: 0; transform: translateY(15px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
