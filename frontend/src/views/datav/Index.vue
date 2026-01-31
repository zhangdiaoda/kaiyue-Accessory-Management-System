<template>
  <div class="datav-container">
    <!-- 背景层 -->
    <div class="bg-grid"></div>
    <div class="bg-particles"></div>

    <!-- 顶部标题栏 -->
    <div class="datav-header">
      <div class="header-deco-left"></div>
      <div class="header-center">
        <h1 class="title">智 慧 仓 储 数 据 驾 驶 舱</h1>
        <div class="sub-title">SMART WAREHOUSE DATA CENTER</div>
      </div>
      <div class="header-right-info">
        <span class="time">{{ currentTime }}</span>
        <span class="weather">晴 24°C</span>
        <el-button link class="exit-btn" @click="$router.push('/dashboard')">
          <el-icon><SwitchButton /></el-icon> 退出
        </el-button>
      </div>
    </div>

    <!-- 主体内容区 -->
    <div class="datav-main">
      <!-- 左翼：资产概况 -->
      <div class="datav-col left-col">
        <div class="dv-card animate-slide-right box-shadow-glow">
          <div class="card-header">
            <span class="card-title">库存结构 (Top 5分类)</span>
            <div class="deco-line"></div>
          </div>
          <div class="card-body" ref="pieChartRef"></div>
        </div>

        <div class="dv-card animate-slide-right delay-1 box-shadow-glow">
          <div class="card-header">
            <span class="card-title">实时预警雷达</span>
            <div class="deco-line"></div>
            <el-tag v-if="lowStockList.length > 0" type="danger" effect="dark" size="small" class="blink-tag">
              {{ lowStockList.length }} 项异常
            </el-tag>
          </div>
          <div class="card-body scroll-body">
            <div v-if="lowStockList.length > 0" class="warning-list">
              <div v-for="(item, index) in lowStockList" :key="item.id" class="warning-item">
                <span class="w-index">{{ String(index + 1).padStart(2, '0') }}</span>
                <div class="w-info">
                  <div class="w-name">{{ item.name }}</div>
                  <div class="w-spec">{{ item.part_no }}</div>
                </div>
                <div class="w-stat">
                  <div class="w-val">库存 {{ item.stock_quantity }}</div>
                  <div class="w-threshold">阈值 {{ item.warning_threshold }}</div>
                </div>
              </div>
            </div>
            <div v-else class="empty-placeholder">暂无库存预警</div>
          </div>
        </div>
      </div>

      <!-- 中控：核心指标 -->
      <div class="datav-col center-col">
        <!-- 数字翻牌器 -->
        <div class="kpi-board animate-zoom-in">
          <div class="kpi-item glow-border">
            <div class="label">资产库存总值 (CNY)</div>
            <div class="value-box neon-blue">
              <span class="num">{{ formatNumber(kpiData.totalValue) }}</span>
            </div>
          </div>
          <div class="kpi-item glow-border">
            <div class="label">本月领用人次 (Total)</div>
            <div class="value-box neon-orange">
              <span class="num">{{ formatNumber(kpiData.monthlyBorrowCount) }}</span>
            </div>
          </div>
        </div>
        
        <!-- 中央动态图：实时作业墙 -->
        <div class="center-live-box glow-border animate-zoom-in delay-2">
          <div class="live-header">
            <span class="live-title">LIVE OPERATIONS FEED</span>
            <div class="live-status">
              <span class="status-dot"></span> 实时监控中
            </div>
          </div>
          <div class="live-list-container">
            <div class="live-list" :class="{ 'scrolling': liveRecords.length > 8 }">
               <div v-for="(record, idx) in liveRecords" :key="record.id" class="live-item">
                 <div class="live-time">{{ formatTime(record.borrow_time) }}</div>
                 <div class="live-content">
                   <span class="highlight-user">{{ record.employee_name }}</span>
                   <span class="action-type" :class="record.status === 'RETURNED' ? 'return' : 'borrow'">
                     {{ record.status === 'RETURNED' ? '归还' : '领用' }}
                   </span>
                   <span class="highlight-part">{{ record.part_name }}</span>
                   <span class="quantity">x{{ record.borrow_quantity }}</span>
                 </div>
               </div>
            </div>
            <div v-if="liveRecords.length === 0" class="empty-live">暂无今日作业记录</div>
          </div>
        </div>
      </div>

      <!-- 右翼：效能分析 -->
      <div class="datav-col right-col">
        <div class="dv-card animate-slide-left box-shadow-glow">
          <div class="card-header">
            <span class="card-title">个人领用排行 (本月Top 5)</span>
            <div class="deco-line"></div>
          </div>
          <div class="card-body" ref="barChartRef"></div>
        </div>

        <div class="dv-card animate-slide-left delay-1 box-shadow-glow">
          <div class="card-header">
            <span class="card-title">物资流转趋势 (12个月)</span>
            <div class="deco-line"></div>
          </div>
          <div class="card-body" ref="lineChartRef"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import dayjs from 'dayjs'
import * as echarts from 'echarts'
import { SwitchButton } from '@element-plus/icons-vue'
import { getDashboardStats } from '@/api/dashboard'
import { getPartList, getLowStockParts } from '@/api/part'
import { getEmployeeReport } from '@/api/report'
import { getBorrowRecordList } from '@/api/borrow'

// 时间
const currentTime = ref(dayjs().format('HH:mm:ss'))
let timer = null

// Refs
const pieChartRef = ref(null)
const barChartRef = ref(null)
const lineChartRef = ref(null)

// Data
const kpiData = reactive({
  totalValue: 0,
  monthlyBorrowCount: 0
})
const lowStockList = ref([])
const liveRecords = ref([])

// Charts Instances
let pieChart = null
let barChart = null
let lineChart = null

const formatNumber = (num) => {
  return num ? num.toLocaleString() : '0'
}

const formatTime = (time) => {
  return dayjs(time).format('MM-DD HH:mm')
}

const initCharts = () => {
  pieChart = echarts.init(pieChartRef.value)
  barChart = echarts.init(barChartRef.value)
  lineChart = echarts.init(lineChartRef.value)

  window.addEventListener('resize', resizeCharts)
}

const resizeCharts = () => {
  pieChart?.resize()
  barChart?.resize()
  lineChart?.resize()
}

const loadData = async () => {
  try {
    // 1. Dashboard Stats (KPI + Trend + Category)
    const statsRes = await getDashboardStats()
    const statsData = statsRes.data || {}
    
    // KPI
    kpiData.totalValue = statsData.base?.total_stock || 0
    kpiData.monthlyBorrowCount = statsData.base?.monthly_borrow || 0

    // Trends
    const trends = statsData.monthly_trend || []
    renderLineChart(trends)

    // Category Pie
    const catStats = statsData.category_stats || []
    renderPieChart(catStats)

    // 2. Low Stock
    const lowStockRes = await getLowStockParts()
    lowStockList.value = lowStockRes.data || []

    // 3. Employee Ranking (Top 5 THIS MONTH)
    const startOfMonth = dayjs().startOf('month').format('YYYY-MM-DD')
    const endOfMonth = dayjs().endOf('month').format('YYYY-MM-DD')
    
    const empRes = await getEmployeeReport({ 
      start_date: startOfMonth, 
      end_date: endOfMonth 
    })
    
    const empData = empRes.data || []
    // Sort and take top 5
    const topEmp = empData
      .sort((a, b) => b.total_borrow - a.total_borrow)
      .slice(0, 5)
    
    renderBarChart(topEmp)
    
    // 4. Live Records (New!)
    const recordRes = await getBorrowRecordList({ 
      page: 1, 
      pageSize: 20 
    })
    liveRecords.value = recordRes.data.records || []

  } catch (error) {
    console.error('DataV load failed', error)
  }
}

const renderPieChart = (data) => {
  // Fallback
  if (!data || data.length === 0) {
    data = [{name: '无数据', value: 0}]
  }

  const option = {
    tooltip: { trigger: 'item' },
    legend: {
      bottom: '0%',
      left: 'center',
      textStyle: { color: '#fff' }
    },
    series: [
      {
        name: '库存分布',
        type: 'pie',
        radius: ['45%', '70%'],
        center: ['50%', '45%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 8,
          borderColor: '#101a37',
          borderWidth: 2
        },
        label: { show: false },
        labelLine: { show: false },
        data: data,
        color: ['#00f2ff', '#0071e3', '#ff9f0a', '#ff375f', '#34c759', '#8e44ad', '#2c3e50']
      }
    ]
  }
  pieChart.setOption(option)
}

const renderBarChart = (data) => {
  const option = {
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    grid: { left: '3%', right: '8%', bottom: '3%', top: '10%', containLabel: true },
    xAxis: {
      type: 'value',
      axisLabel: { color: '#89a3c9' },
      splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)' } }
    },
    yAxis: {
      type: 'category',
      data: data.map(i => i.employee_name),
      axisLabel: { color: '#fff', fontSize: 13 },
      axisLine: { show: false },
      axisTick: { show: false }
    },
    series: [
      {
        name: '个人领用量',
        type: 'bar',
        data: data.map(i => i.total_borrow),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: 'rgba(0, 113, 227, 0.6)' },
            { offset: 1, color: '#00f2ff' }
          ]),
          borderRadius: [0, 100, 100, 0]
        },
        barWidth: 16,
        label: { show: true, position: 'right', color: '#fff' }
      }
    ]
  }
  barChart.setOption(option)
}

const renderLineChart = (data) => {
  // If no data, use empty arrays
  const months = data ? data.map(i => i.month) : []
  const borrowData = data ? data.map(i => i.borrow) : []
  const returnData = data ? data.map(i => i.return) : []

  const option = {
    tooltip: { trigger: 'axis', backgroundColor: 'rgba(0,0,0,0.8)', borderColor: '#00f2ff', textStyle: { color: '#fff' } },
    legend: { textStyle: { color: '#fff' }, top: 0, right: 10 },
    grid: { left: '3%', right: '4%', bottom: '3%', top: '15%', containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: months,
      axisLabel: { color: '#89a3c9' },
      axisLine: { lineStyle: { color: 'rgba(255,255,255,0.1)' } }
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#89a3c9' },
      splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)' } }
    },
    series: [
      {
        name: '借出',
        type: 'line',
        smooth: true,
        showSymbol: false,
        data: borrowData,
        lineStyle: { color: '#00f2ff', width: 3 },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(0, 242, 255, 0.2)' },
            { offset: 1, color: 'rgba(0, 242, 255, 0)' }
          ])
        },
        itemStyle: { color: '#00f2ff' }
      },
      {
        name: '归还',
        type: 'line',
        smooth: true,
        showSymbol: false,
        data: returnData,
        lineStyle: { color: '#34c759', width: 3 },
        itemStyle: { color: '#34c759' }
      }
    ]
  }
  lineChart.setOption(option)
}

onMounted(() => {
  initCharts()
  loadData()
  timer = setInterval(() => {
    currentTime.value = dayjs().format('HH:mm:ss')
    // 可选：定时刷新数据
    // loadData() 
  }, 1000)
})

onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
  pieChart?.dispose()
  barChart?.dispose()
  lineChart?.dispose()
  window.removeEventListener('resize', resizeCharts)
})
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700;900&display=swap');

.datav-container {
  width: 100vw;
  height: 100vh;
  background: #0b0c15;
  color: #fff;
  overflow: hidden;
  font-family: "PingFang SC", "Microsoft YaHei", sans-serif;
  display: flex;
  flex-direction: column;
  position: relative;
}

/* Background Effects */
.bg-grid {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: 
    linear-gradient(rgba(0, 242, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 242, 255, 0.03) 1px, transparent 1px);
  background-size: 30px 30px;
  z-index: 0;
  pointer-events: none;
}

.bg-particles {
  position: absolute;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle at 50% 50%, rgba(0, 113, 227, 0.1) 0%, transparent 60%);
  z-index: 0;
  animation: pulse 10s infinite;
}

@keyframes pulse {
  0% { opacity: 0.5; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.1); }
  100% { opacity: 0.5; transform: scale(1); }
}

/* Header */
.datav-header {
  height: 80px;
  background: url('https://img.alicdn.com/tfs/TB1Typew.z1gK0jSZLeXXb9kVXa-1920-80.png') no-repeat center bottom;
  background-size: 100% 100%;
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-shrink: 0;
  z-index: 10;
}

.title {
  font-family: 'Orbitron', sans-serif;
  font-size: 32px;
  font-weight: 900;
  letter-spacing: 4px;
  background: linear-gradient(180deg, #ffffff 0%, #76c0ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
  text-align: center;
  text-shadow: 0 0 20px rgba(0, 140, 255, 0.6);
}

.sub-title {
  font-size: 10px;
  color: #4facfe;
  text-align: center;
  letter-spacing: 6px;
  margin-top: 4px;
  opacity: 0.8;
  text-transform: uppercase;
}

.header-right-info {
  position: absolute;
  right: 30px;
  top: 25px;
  display: flex;
  align-items: center;
  gap: 20px;
  font-family: 'Orbitron', sans-serif;
}

.exit-btn {
  color: #89a3c9;
  font-size: 14px;
}
.exit-btn:hover { color: #fff; }

.time {
  font-size: 20px;
  color: #00f2ff;
  font-weight: 700;
}

.weather {
  font-size: 16px;
  color: #89a3c9;
}

/* Main Layout */
.datav-main {
  flex: 1;
  padding: 20px;
  display: flex;
  gap: 20px;
  min-height: 0;
  position: relative;
  z-index: 1;
}

.datav-col {
  display: flex;
  flex-direction: column;
  gap: 20px;
  height: 100%;
}

.left-col, .right-col {
  flex: 3;
}

.center-col {
  flex: 4;
}

/* Cards */
.dv-card {
  background: rgba(11, 19, 41, 0.7);
  border: 1px solid rgba(48, 114, 246, 0.2);
  flex: 1;
  padding: 15px;
  display: flex;
  flex-direction: column;
  position: relative;
  min-height: 0;
  backdrop-filter: blur(10px);
}

.box-shadow-glow {
  box-shadow: 0 0 15px rgba(0, 113, 227, 0.1);
}

.dv-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 15px;
  height: 15px;
  border-top: 2px solid #00f2ff;
  border-left: 2px solid #00f2ff;
}

.dv-card::after {
  content: '';
  position: absolute;
  bottom: 0;
  right: 0;
  width: 15px;
  height: 15px;
  border-bottom: 2px solid #00f2ff;
  border-right: 2px solid #00f2ff;
}

.card-header {
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-bottom: 12px;
  position: relative;
}

.card-title {
  font-size: 16px;
  font-weight: 700;
  color: #e0e6ed;
  padding-left: 10px;
  border-left: 4px solid #00f2ff;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.deco-line {
  height: 1px;
  width: 100%;
  background: linear-gradient(90deg, rgba(0, 242, 255, 0.6) 0%, transparent 100%);
  margin-top: 10px;
}

.card-body {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.scroll-body {
  overflow-y: hidden; /* Hide default scroll */
  position: relative;
}

/* Warning List */
.warning-item {
  display: flex;
  align-items: center;
  background: rgba(43, 22, 30, 0.6);
  margin-bottom: 10px;
  padding: 10px;
  border-radius: 4px;
  border-left: 3px solid #ff375f;
  transition: all 0.3s;
  animation: fadeIn 0.5s;
}

.warning-item:hover {
  background: rgba(255, 55, 95, 0.15);
  transform: translateX(5px);
}

.w-index {
  font-family: 'Orbitron';
  color: #ff375f;
  font-weight: 700;
  margin-right: 12px;
  font-size: 16px;
  opacity: 0.5;
}

.w-info { flex: 1; }

.w-name { font-size: 14px; font-weight: 600; color: #fff; margin-bottom: 2px; }
.w-spec { font-size: 12px; color: #89a3c9; }

.w-stat { text-align: right; }
.w-val { color: #ff375f; font-weight: 700; }
.w-threshold { color: #89a3c9; font-size: 10px; }

.blink-tag {
  position: absolute;
  right: 0;
  top: 0;
  animation: blink 2s infinite;
}

@keyframes blink {
  0% { opacity: 1; box-shadow: 0 0 10px rgba(255, 55, 95, 0.5); }
  50% { opacity: 0.5; }
  100% { opacity: 1; box-shadow: 0 0 10px rgba(255, 55, 95, 0.5); }
}

/* KPI Board */
.kpi-board {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 20px;
}

.kpi-item {
  flex: 1;
  background: rgba(11, 19, 41, 0.8);
  padding: 20px;
  text-align: center;
  border: 1px solid rgba(0, 242, 255, 0.1);
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: relative;
}

.glow-border {
  box-shadow: 0 0 15px rgba(0, 242, 255, 0.1);
  transition: all 0.3s;
}
.glow-border:hover {
  box-shadow: 0 0 25px rgba(0, 242, 255, 0.2);
  border-color: rgba(0, 242, 255, 0.4);
}

.kpi-item .label {
  font-size: 14px;
  color: #89a3c9;
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.kpi-item .value-box {
  background: rgba(0, 0, 0, 0.3);
  padding: 10px;
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

.kpi-item .num {
  font-family: 'Orbitron', sans-serif;
  font-size: 42px;
  font-weight: 700;
  letter-spacing: 2px;
}

.neon-blue { color: #00f2ff; text-shadow: 0 0 15px rgba(0, 242, 255, 0.6); }
.neon-orange { color: #ff9f0a; text-shadow: 0 0 15px rgba(255, 159, 10, 0.6); }

/* Center Live Box */
.center-live-box {
  flex: 1;
  background: rgba(11, 19, 41, 0.9);
  border: 1px solid rgba(0, 242, 255, 0.2);
  padding: 20px;
  display: flex;
  flex-direction: column;
  position: relative;
  backdrop-filter: blur(20px);
}

.live-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding-bottom: 15px;
}

.live-title {
  font-family: 'Orbitron';
  font-weight: 700;
  color: #fff;
  font-size: 18px;
  letter-spacing: 2px;
}

.live-status {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #00f2ff;
  gap: 6px;
}

.status-dot {
  width: 8px;
  height: 8px;
  background: #00f2ff;
  border-radius: 50%;
  box-shadow: 0 0 10px #00f2ff;
  animation: blink 2s infinite;
}

.live-list-container {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.live-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.live-item {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.03);
  padding: 12px 16px;
  border-radius: 4px;
  border-left: 2px solid transparent;
  transition: all 0.2s;
}

.live-item:hover {
  background: rgba(255, 255, 255, 0.08);
  border-left-color: #00f2ff;
}

.live-time {
  width: 90px;
  font-family: 'Orbitron';
  color: #89a3c9;
  font-size: 12px;
}

.live-content {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #e0e6ed;
  font-size: 14px;
}

.highlight-user { color: #fff; font-weight: 600; }
.highlight-part { color: #00f2ff; }

.action-type {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}
.action-type.borrow { background: rgba(0, 242, 255, 0.15); color: #00f2ff; }
.action-type.return { background: rgba(52, 199, 89, 0.15); color: #34c759; }

.quantity {
  margin-left: auto;
  font-family: 'Orbitron';
  color: #ff9f0a;
  font-weight: 700;
}

/* Animations */
.animate-slide-right { animation: slideRight 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) forwards; opacity: 0; }
.animate-slide-left { animation: slideLeft 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) forwards; opacity: 0; }
.animate-zoom-in { animation: zoomIn 0.8s cubic-bezier(0.2, 0.8, 0.2, 1) forwards; opacity: 0; }
.delay-1 { animation-delay: 0.1s; }
.delay-2 { animation-delay: 0.2s; }

@keyframes slideRight { from { transform: translateX(-30px); opacity: 0; } to { transform: translateX(0); opacity: 1; } }
@keyframes slideLeft { from { transform: translateX(30px); opacity: 0; } to { transform: translateX(0); opacity: 1; } }
@keyframes zoomIn { from { transform: scale(0.95); opacity: 0; } to { transform: scale(1); opacity: 1; } }

</style>
