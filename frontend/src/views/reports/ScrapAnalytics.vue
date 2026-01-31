<template>
  <div class="scrap-analytics">
    <!-- Header Section -->
    <el-card class="filter-card glass-effect" shadow="never">
      <div class="header-content">
        <div class="title-group">
          <h2>♻️ 废品损毁分析看板</h2>
          <p class="subtitle">实时监控资产损毁率，分析废品产生的经济影响</p>
        </div>
        <div class="filter-actions">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 280px"
            @change="loadAllData"
          />
          <el-button type="primary" plain @click="loadAllData" :loading="loading" style="margin-left: 15px">
            <el-icon><Refresh /></el-icon> 刷新数据
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Essential Metrics -->
    <el-row :gutter="20" class="metrics-row">
      <el-col :span="6" v-for="(stat, index) in summaryStats" :key="index">
        <div class="metric-card shadow-sm">
          <div class="m-icon" :style="{ backgroundColor: stat.bg, color: stat.color }">
            <el-icon :size="24"><component :is="stat.icon" /></el-icon>
          </div>
          <div class="m-info">
            <div class="m-label">{{ stat.label }}</div>
            <div class="m-value">{{ stat.value }}<small v-if="stat.unit">{{ stat.unit }}</small></div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Main Charts -->
    <el-row :gutter="20" style="margin-top: 24px">
      <el-col :span="16">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>📈 废品产生月度趋势</span>
              <el-tag size="small" type="info">柱状趋势图</el-tag>
            </div>
          </template>
          <div ref="trendChartRef" style="height: 380px"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>🚨 废品率 Top 6 配件</span>
              <el-tag size="small" type="danger">风险排行</el-tag>
            </div>
          </template>
          <div ref="partRankChartRef" style="height: 380px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 24px">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>👥 员工废品责任分布 (Top 8)</span>
            </div>
          </template>
          <div ref="employeePieRef" style="height: 350px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>📦 重点监控部门废品率</span>
            </div>
          </template>
          <div ref="deptRadarRef" style="height: 350px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row style="margin-top: 24px">
      <el-col :span="24">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>⚠️ 具体报废人责任排行 (全周期)</span>
              <el-tag type="danger">重点管控对象</el-tag>
            </div>
          </template>
          <div class="scrapper-grid">
            <div v-for="(person, i) in sortedTopScrappers" :key="i" class="scrapper-mini-card glass-effect">
              <div class="s-rank">TOP {{ i + 1 }}</div>
              <div class="s-name">{{ person.employee_name }}</div>
              <div class="s-dept">{{ person.department }}</div>
              <div class="s-val">{{ person.total_damaged }} <small>件</small></div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>


<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { getDetailedReport, getPartReport, getEmployeeReport, getDepartmentReport } from '@/api/report'

const loading = ref(false)
const dateRange = ref([
  new Date(new Date().setMonth(new Date().getMonth() - 6)),
  new Date()
])

// Refs for ECharts
const trendChartRef = ref(null)
const partRankChartRef = ref(null)
const employeePieRef = ref(null)
const deptRadarRef = ref(null)

let charts = []

const rawDetailedData = ref([])
const rawPartData = ref([])
const rawEmployeeData = ref([])
const rawDeptData = ref([])

const sortedTopScrappers = computed(() => {
  return [...rawEmployeeData.value]
    .sort((a, b) => b.total_damaged - a.total_damaged)
    .slice(0, 10)
})

const summaryStats = computed(() => {
  const totalScrap = rawDetailedData.value.reduce((sum, i) => sum + i.total_damaged, 0)
  const totalBorrow = rawDetailedData.value.reduce((sum, i) => sum + i.total_borrow, 0)
  const scrapRate = totalBorrow > 0 ? ((totalScrap / totalBorrow) * 100).toFixed(2) : 0
  
  return [
    { label: '累计废品总数', value: totalScrap, unit: '件', icon: 'Delete', color: '#f43f5e', bg: '#fff1f2' },
    { label: '全周期废品率', value: scrapRate, unit: '%', icon: 'TrendCharts', color: '#ef4444', bg: '#fef2f2' },
    { label: '受影响配件种类', value: rawPartData.value.filter(p => p.total_damaged > 0).length, unit: '类', icon: 'Box', color: '#f59e0b', bg: '#fffbeb' },
    { label: '人均废品贡献量', value: rawEmployeeData.value.length > 0 ? (totalScrap / rawEmployeeData.value.length).toFixed(1) : 0, unit: '件/人', icon: 'User', color: '#0ea5e9', bg: '#f0f9ff' }
  ]
})

const loadAllData = async () => {
  loading.value = true
  const params = {
    start_date: formatDate(dateRange.value[0]),
    end_date: formatDate(dateRange.value[1])
  }

  try {
    const [detailed, part, employee, dept] = await Promise.all([
      getDetailedReport(params),
      getPartReport(params),
      getEmployeeReport(params),
      getDepartmentReport(params)
    ])

    rawDetailedData.value = detailed.data || []
    rawPartData.value = part.data || []
    rawEmployeeData.value = employee.data || []
    rawDeptData.value = dept.data || []

    await nextTick()
    renderCharts()
  } catch (error) {
    ElMessage.error('废品分析数据加载失败')
  } finally {
    loading.value = false
  }
}

const renderCharts = () => {
  // 1. Trend Chart
  const trendChart = echarts.init(trendChartRef.value)
  const monthlyScrap = {}
  rawDetailedData.value.forEach(i => {
    monthlyScrap[i.month] = (monthlyScrap[i.month] || 0) + i.total_damaged
  })
  const months = Object.keys(monthlyScrap).sort()
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: months, axisLabel: { color: '#86868b' } },
    yAxis: { type: 'value', axisLabel: { color: '#86868b' } },
    series: [{
      name: '废品数量',
      type: 'bar',
      barWidth: '40%',
      data: months.map(m => monthlyScrap[m]),
      itemStyle: { 
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: '#f43f5e' },
          { offset: 1, color: '#fb7185' }
        ]),
        borderRadius: [6, 6, 0, 0]
      }
    }]
  })

  // 2. Part Rank
  const partChart = echarts.init(partRankChartRef.value)
  const sortedParts = [...rawPartData.value]
    .filter(p => p.damage_rate > 0)
    .sort((a, b) => b.damage_rate - a.damage_rate)
    .slice(0, 6)
  partChart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c}%' },
    xAxis: { type: 'value', max: 100 },
    yAxis: { type: 'category', data: sortedParts.map(p => p.part_name).reverse() },
    series: [{
      type: 'bar',
      data: sortedParts.map(p => p.damage_rate).reverse(),
      label: { show: true, position: 'right', formatter: '{c}%' },
      itemStyle: { color: '#f43f5e' }
    }]
  })

  // 3. Employee Pie
  const employeeChart = echarts.init(employeePieRef.value)
  const sortedEmp = [...rawEmployeeData.value]
    .sort((a, b) => b.total_damaged - a.total_damaged)
    .slice(0, 8)
  employeeChart.setOption({
    tooltip: { trigger: 'item' },
    legend: { bottom: '5%', left: 'center' },
    series: [{
      name: '废品件数',
      type: 'pie',
      radius: ['45%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
      label: { show: false, position: 'center' },
      emphasis: { label: { show: true, fontSize: '18', fontWeight: 'bold' } },
      data: sortedEmp.map(e => ({ value: e.total_damaged, name: e.employee_name }))
    }]
  })

  // 4. Dept Radar
  const radarChart = echarts.init(deptRadarRef.value)
  const depts = rawDeptData.value.slice(0, 5)
  radarChart.setOption({
    radar: {
      indicator: depts.map(d => ({ name: d.department, max: Math.max(...rawDeptData.value.map(x => x.total_damaged)) || 10 }))
    },
    series: [{
      type: 'radar',
      data: [{ value: depts.map(d => d.total_damaged), name: '废品分布', areaStyle: { color: 'rgba(244, 63, 94, 0.3)' }, lineStyle: { color: '#f43f5e' } }]
    }]
  })

  charts = [trendChart, partChart, employeeChart, radarChart]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const handleResize = () => {
  charts.forEach(c => c.resize())
}

onMounted(() => {
  loadAllData()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  charts.forEach(c => c.dispose())
})
</script>

<style scoped>
.scrap-analytics {
  max-width: 1400px;
  margin: 0 auto;
}

.filter-card {
  margin-bottom: 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-group h2 { margin: 0; font-size: 20px; }
.subtitle { margin: 4px 0 0; font-size: 13px; color: #94a3b8; }

.metrics-row {
  margin-top: 10px;
}

.metric-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
}

.m-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.m-info {
  display: flex;
  flex-direction: column;
}

.m-label {
  font-size: 13px;
  color: #86868b;
  font-weight: 500;
}

.m-value {
  font-size: 26px;
  font-weight: 800;
  color: #1d1d1f;
  margin-top: 2px;
}

.m-value small {
  font-size: 13px;
  font-weight: 500;
  margin-left: 4px;
  color: #94a3b8;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 15px;
}

:deep(.el-card__header) {
  padding: 18px 24px;
  border-bottom: 1px solid rgba(0,0,0,0.05);
}

.shadow-sm {
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

.scrapper-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

.scrapper-mini-card {
  flex: 1;
  min-width: 150px;
  padding: 15px;
  border-radius: 12px;
  border: 1px solid rgba(0,0,0,0.03);
  text-align: center;
}

.s-rank { font-size: 10px; font-weight: 800; color: #f43f5e; margin-bottom: 5px; }
.s-name { font-size: 14px; font-weight: 700; color: #1d1d1f; }
.s-dept { font-size: 11px; color: #86868b; margin-bottom: 8px; }
.s-val { font-size: 18px; font-weight: 800; color: #1d1d1f; }
.s-val small { font-size: 12px; font-weight: 400; color: #86868b; }
</style>

