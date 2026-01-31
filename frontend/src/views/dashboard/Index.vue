<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="24">
      <el-col :span="6" v-for="(item, index) in stats" :key="index">
        <el-card shadow="never" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" :style="{ backgroundColor: item.color, color: item.textColor }">
              <el-icon :size="28">
                <component :is="item.icon" />
              </el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ item.value }}</div>
              <div class="stat-label">{{ item.label }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 可视化图表 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>📈 月度领用趋势（近6个月）</span>
          </template>
          <div ref="trendChartRef" style="width: 100%; height: 320px"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never">
          <template #header>
            <div class="card-header-flex">
              <span>🚨 废品率 Top 5</span>
              <el-tag size="small" type="danger">风险预警</el-tag>
            </div>
          </template>
          <div ref="damageChartRef" style="width: 100%; height: 320px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <span>📊 员工领用排名 Top 10</span>
              <el-select 
                v-model="selectedPartId" 
                placeholder="全部配件" 
                clearable 
                filterable
                style="width: 200px"
                @change="handlePartChange"
              >
                <el-option label="全部配件" :value="''" />
                <el-option
                  v-for="part in allParts"
                  :key="part.id"
                  :label="part.name"
                  :value="part.id"
                />
              </el-select>
            </div>
          </template>
          <div ref="employeeChartRef" style="width: 100%; height: 320px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>🥧 配件领用占比 Top 8</span>
          </template>
          <div ref="partPieChartRef" style="width: 100%; height: 320px"></div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row style="margin-top: 20px">
      <el-col :span="24">
        <el-card shadow="never">
          <template #header>
            <div class="card-header-flex">
              <span>👤 报废责任详情 (按员工溯源)</span>
              <el-tag type="info">损耗治理核心看板</el-tag>
            </div>
          </template>
          <div class="scrapper-list">
            <el-row :gutter="20">
              <el-col :span="4" v-for="(person, i) in topScrappers" :key="i">
                <div class="scrapper-item glass-effect">
                  <div class="scrapper-rank">#{{ i + 1 }}</div>
                  <div class="scrapper-name">{{ person.name }}</div>
                  <div class="scrapper-dept">{{ person.department }}</div>
                  <div class="scrapper-count">
                    累计报废 <span class="num">{{ person.count }}</span> 件
                  </div>
                  <el-progress :percentage="Math.min(100, (person.count / (topScrappers[0]?.count || 1)) * 100)" :show-text="false" status="exception" stroke-width="4" />
                </div>
              </el-col>
              <el-col :span="24" v-if="topScrappers.length === 0">
                <el-empty description="暂无历史报废责任记录" :image-size="60" />
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>


<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { getDashboardStats } from '@/api/dashboard'
import { getEmployeeReport, getPartReport, getDetailedReport } from '@/api/report'
import { getPartList } from '@/api/part'

const stats = reactive([
  { label: '配件总数', value: 0, icon: 'Box', color: '#f0f4ff', textColor: '#2563eb' },
  { label: '员工总数', value: 0, icon: 'User', color: '#fff1f2', textColor: '#e11d48' },
  { label: '领用中', value: 0, icon: 'DocumentCopy', color: '#f0fdf4', textColor: '#16a34a' },
  { label: '低库存预警', value: 0, icon: 'Warning', color: '#fffbeb', textColor: '#d97706' }
])

const employeeChartRef = ref(null)
const partPieChartRef = ref(null)
const trendChartRef = ref(null)
const damageChartRef = ref(null)

const selectedPartId = ref(null)
const allParts = ref([])
const topScrappers = ref([])

let employeeChart = null
let partPieChart = null
let trendChart = null
let damageChart = null

const loadStats = async () => {
  try {
    const res = await getDashboardStats()
    const data = res.data || {}
    
    // 根据后端API实际返回的字段映射
    stats[0].value = data.base?.part_count || 0
    stats[1].value = data.base?.employee_count || 0
    stats[2].value = data.base?.borrowed_count || 0
    stats[3].value = data.base?.warning_count || 0
    
    topScrappers.value = data.top_scrappers || []
  } catch (error) {
    console.error('加载统计数据失败', error)
  }
}

const initCharts = () => {
  if (!employeeChartRef.value || !partPieChartRef.value || !trendChartRef.value || !damageChartRef.value) {
    console.error('图表DOM元素未找到')
    return
  }

  employeeChart = echarts.init(employeeChartRef.value)
  partPieChart = echarts.init(partPieChartRef.value)
  trendChart = echarts.init(trendChartRef.value)
  damageChart = echarts.init(damageChartRef.value)

  window.addEventListener('resize', handleResize)
}

const handleResize = () => {
  employeeChart?.resize()
  partPieChart?.resize()
  trendChart?.resize()
  damageChart?.resize()
}

const formatDate = (date) => {
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const loadChartData = async () => {
  const endDate = new Date()
  const startDate = new Date()
  startDate.setMonth(startDate.getMonth() - 6)
  
  const params = {
    start_date: formatDate(startDate),
    end_date: formatDate(endDate)
  }

  try {
    await Promise.all([
      loadEmployeeRanking(params),
      loadPartPie(params),
      loadMonthlyTrend(params),
      loadDamageRate(params)
    ])
  } catch (error) {
    console.error('加载图表数据失败', error)
  }
}

const loadEmployeeRanking = async (params) => {
  try {
    // 如果选择了配件，添加筛选参数
    const queryParams = params || {
      start_date: formatDate(new Date(new Date().setMonth(new Date().getMonth() - 6))),
      end_date: formatDate(new Date())
    }
    
    if (selectedPartId.value) {
      queryParams.part_id = selectedPartId.value
    }

    const res = await getEmployeeReport(queryParams)
    const data = (res.data || []).slice(0, 10)

    // 获取选中配件的名称
    const selectedPartName = selectedPartId.value 
      ? allParts.value.find(p => p.id === selectedPartId.value)?.name || '选中配件'
      : '全部配件'

    const option = {
      backgroundColor: 'transparent',
      tooltip: { 
        trigger: 'axis', 
        axisPointer: { type: 'shadow' },
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        borderColor: '#eee',
        borderWidth: 1,
        textStyle: { color: '#1d1d1f' },
        formatter: (params) => {
          const item = params[0]
          return `<div style="font-weight:600;margin-bottom:4px">${item.name}</div>
                  <div style="color:#86868b">领用总数：<span style="color:#0071e3;font-weight:700">${item.value}</span></div>`
        }
      },
      grid: { left: '3%', right: '10%', bottom: '3%', containLabel: true },
      xAxis: { 
        type: 'value', 
        name: '领用数量',
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { type: 'dashed', color: '#f5f5f7' } }
      },
      yAxis: {
        type: 'category',
        data: data.map(item => item.employee_name).reverse(),
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: { color: '#86868b', fontSize: 12 }
      },
      series: [{
        name: '领用总数',
        type: 'bar',
        barWidth: '50%',
        data: data.map(item => item.total_borrow).reverse(),
        itemStyle: {
          borderRadius: [0, 6, 6, 0],
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: '#00c6ff' },
            { offset: 1, color: '#0072ff' }
          ])
        },
        label: { show: true, position: 'right', color: '#1d1d1f', fontWeight: 'bold' }
      }]
    }
    employeeChart.setOption(option)
  } catch (error) {
    console.error('加载员工排名失败', error)
  }
}

const loadPartPie = async (params) => {
  try {
    const res = await getPartReport(params)
    const data = (res.data || []).slice(0, 8)

    const option = {
      tooltip: { 
        trigger: 'item', 
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        textStyle: { color: '#1d1d1f' }
      },
      legend: { orient: 'vertical', left: 'left', itemGap: 15, textStyle: { color: '#86868b' } },
      series: [{
        name: '领用数量',
        type: 'pie',
        radius: ['55%', '80%'],
        center: ['60%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 12, borderColor: '#fff', borderWidth: 2 },
        label: { show: false, position: 'center' },
        emphasis: {
          label: { show: true, fontSize: 20, fontWeight: 'bold' }
        },
        data: data.map((item, idx) => ({ 
          value: item.total_borrow, 
          name: item.part_name,
          itemStyle: {
            color: [
              '#0071e3', '#32d74b', '#ff9f0a', '#ff375f', '#bf5af2', '#64d2ff', '#ffd60a', '#ff453a'
            ][idx % 8]
          }
        }))
      }]
    }
    partPieChart.setOption(option)
  } catch (error) {
    console.error('加载配件占比失败', error)
  }
}

const loadMonthlyTrend = async (params) => {
  try {
    const res = await getDetailedReport(params)
    const data = res.data || []

    const monthlyData = {}
    data.forEach(item => {
      if (!monthlyData[item.month]) {
        monthlyData[item.month] = { borrow: 0, damaged: 0 }
      }
      monthlyData[item.month].borrow += item.total_borrow
      monthlyData[item.month].damaged += item.total_damaged
    })

    const months = Object.keys(monthlyData).sort()
    const borrowData = months.map(m => monthlyData[m].borrow)
    const damagedData = months.map(m => monthlyData[m].damaged)

    const option = {
      tooltip: { 
        trigger: 'axis',
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        axisPointer: { lineStyle: { color: '#f5f5f7' } }
      },
      legend: { icon: 'circle', right: 10, textStyle: { color: '#86868b' } },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { 
        type: 'category', 
        boundaryGap: false, 
        data: months,
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: { color: '#86868b' }
      },
      yAxis: { 
        type: 'value',
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { color: '#f5f5f7' } }
      },
      series: [
        {
          name: '领用数量',
          type: 'line',
          smooth: true,
          showSymbol: false,
          data: borrowData,
          lineStyle: { width: 4, color: '#0071e3' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(0, 113, 227, 0.15)' },
              { offset: 1, color: 'rgba(0, 113, 227, 0)' }
            ])
          }
        },
        {
          name: '报废数量',
          type: 'line',
          smooth: true,
          showSymbol: false,
          data: damagedData,
          lineStyle: { width: 4, color: '#ff375f' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(255, 55, 95, 0.1)' },
              { offset: 1, color: 'rgba(255, 55, 95, 0)' }
            ])
          }
        }
      ]
    }
    trendChart.setOption(option)
  } catch (error) {
    console.error('加载月度趋势失败', error)
  }
}

const loadDamageRate = async (params) => {
  try {
    const res = await getPartReport(params)
    const data = (res.data || [])
      .filter(item => item.damage_rate > 0)
      .sort((a, b) => b.damage_rate - a.damage_rate)
      .slice(0, 5)

    const option = {
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        data: data.map(item => item.part_name),
        axisLine: { show: false },
        axisTick: { show: false }
      },
      yAxis: { 
        type: 'value', 
        name: '报废率(%)',
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { color: '#f5f5f7' } }
      },
      series: [{
        name: '报废率',
        type: 'bar',
        barWidth: '40%',
        data: data.map(item => item.damage_rate || 0),
        itemStyle: {
          borderRadius: [6, 6, 0, 0],
          color: params => {
            const value = params.value
            if (value > 10) return '#ff375f'
            if (value > 5) return '#ff9f0a'
            return '#32d74b'
          }
        },
        label: { show: true, position: 'top', formatter: '{c}%', color: '#1d1d1f' }
      }]
    }
    damageChart.setOption(option)


  } catch (error) {
    console.error('加载损毁率失败', error)
  }
}

const handlePartChange = () => {
  // 切换配件时重新加载员工排名图表
  const endDate = new Date()
  const startDate = new Date()
  startDate.setMonth(startDate.getMonth() - 6)
  
  loadEmployeeRanking({
    start_date: formatDate(startDate),
    end_date: formatDate(endDate)
  })
}

const loadParts = async () => {
  try {
    const res = await getPartList({ page: 1, pageSize: 1000 })
    allParts.value = res.data.records || []
  } catch (error) {
    console.error('加载配件列表失败', error)
  }
}

onMounted(async () => {
  loadStats()
  loadParts()
  await nextTick()
  initCharts()
  loadChartData()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  employeeChart?.dispose()
  partPieChart?.dispose()
  trendChart?.dispose()
  damageChart?.dispose()
})
</script>

<style scoped>
.dashboard {
  height: 100%;
  overflow-y: auto;
  padding-bottom: 20px;
}

.stat-card {
  cursor: pointer;
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: 800;
  color: #1d1d1f;
  letter-spacing: -0.02em;
}

.stat-label {
  font-size: 13px;
  color: #86868b;
  margin-top: 2px;
  font-weight: 500;
}

.scrapper-item {
  padding: 20px;
  border-radius: 16px;
  text-align: center;
  border: 1px solid rgba(0,0,0,0.03);
  transition: all 0.3s;
}

.scrapper-item:hover {
  transform: translateY(-5px);
  background: white;
  box-shadow: 0 10px 30px rgba(0,0,0,0.05);
}

.scrapper-rank {
  font-size: 11px;
  font-weight: 800;
  color: #ff375f;
  margin-bottom: 8px;
  letter-spacing: 0.05em;
}

.scrapper-name {
  font-size: 16px;
  font-weight: 700;
  color: #1d1d1f;
}

.scrapper-dept {
  font-size: 11px;
  color: #86868b;
  margin-bottom: 12px;
}

.scrapper-count {
  font-size: 12px;
  color: #86868b;
  margin-bottom: 8px;
}

.scrapper-count .num {
  font-size: 18px;
  font-weight: 800;
  color: #1d1d1f;
}

/* ========== 移动端优化 ========== */
@media (max-width: 768px) {
  .dashboard {
    padding: 12px;
  }

  /* 统计卡片列数调整为单列 */
  :deep(.el-col) {
    max-width: 100%;
    flex: 0 0 100%;
  }

  .stat-card {
    margin-bottom: 12px;
  }

  .stat-value {
    font-size: 28px;
  }

  /* 图表容器高度调整 */
  .el-card {
    margin-bottom: 16px;
  }

  /* 报废人员卡片 */
  .scrapper-list :deep(.el-row) {
    flex-direction: column;
  }

  .scrapper-list :deep(.el-col) {
    margin-bottom: 12px;
  }

  /* 配件选择器全宽 */
  .el-select {
    width: 100% !important;
  }
}

/* 平板优化 */
@media (min-width: 768px) and (max-width: 1024px) {
  :deep(.el-col-6) {
    max-width: 50%;
    flex: 0 0 50%;
  }

  :deep(.el-col-16),
  :deep(.el-col-8),
  :deep(.el-col-12) {
    max-width: 100%;
    flex: 0 0 100%;
    margin-bottom: 16px;
  }

  .scrapper-list :deep(.el-col-4) {
    max-width: 50%;
    flex: 0 0 50%;
  }
}
</style>

