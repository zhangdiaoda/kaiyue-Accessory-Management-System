<template>
  <div class="data-analytics">
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>数据分析</span>
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                style="width: 240px"
                @change="loadData"
              />
            </div>
          </template>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 员工领用排名 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>📊 员工领用排名（Top 10）</span>
          </template>
          <div ref="employeeChartRef" style="width: 100%; height: 350px"></div>
        </el-card>
      </el-col>

      <!-- 配件领用占比 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>🥧 配件领用占比（Top 10）</span>
          </template>
          <div ref="partPieChartRef" style="width: 100%; height: 350px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 月度领用趋势 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>📈 月度领用趋势</span>
          </template>
          <div ref="trendChartRef" style="width: 100%; height: 350px"></div>
        </el-card>
      </el-col>

      <!-- 损毁率分析 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>⚠️ 配件损毁率（Top 10）</span>
          </template>
          <div ref="damageChartRef" style="width: 100%; height: 350px"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { getEmployeeReport, getPartReport, getDetailedReport } from '@/api/report'

const dateRange = ref([
  new Date(new Date().setMonth(new Date().getMonth() - 3)),
  new Date()
])

const employeeChartRef = ref(null)
const partPieChartRef = ref(null)
const trendChartRef = ref(null)
const damageChartRef = ref(null)

let employeeChart = null
let partPieChart = null
let trendChart = null
let damageChart = null

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 初始化图表
const initCharts = () => {
  if (!employeeChartRef.value || !partPieChartRef.value || !trendChartRef.value || !damageChartRef.value) {
    console.error('图表DOM元素未找到')
    return
  }

  employeeChart = echarts.init(employeeChartRef.value)
  partPieChart = echarts.init(partPieChartRef.value)
  trendChart = echarts.init(trendChartRef.value)
  damageChart = echarts.init(damageChartRef.value)

  // 响应式
  window.addEventListener('resize', handleResize)
}

const handleResize = () => {
  employeeChart?.resize()
  partPieChart?.resize()
  trendChart?.resize()
  damageChart?.resize()
}

// 加载所有数据
const loadData = async () => {
  const params = {
    start_date: formatDate(dateRange.value[0]),
    end_date: formatDate(dateRange.value[1])
  }

  try {
    await Promise.all([
      loadEmployeeRanking(params),
      loadPartPie(params),
      loadMonthlyTrend(params),
      loadDamageRate(params)
    ])
  } catch (error) {
    console.error('加载数据失败', error)
    ElMessage.error('加载数据失败')
  }
}

// 1. 员工领用排名
const loadEmployeeRanking = async (params) => {
  try {
    const res = await getEmployeeReport(params)
    const data = (res.data || []).slice(0, 10)

    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'shadow' }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'value',
        name: '领用数量'
      },
      yAxis: {
        type: 'category',
        data: data.map(item => item.employee_name).reverse(),
        axisLabel: {
          interval: 0,
          rotate: 0
        }
      },
      series: [
        {
          name: '领用总数',
          type: 'bar',
          data: data.map(item => item.total_borrow).reverse(),
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
              { offset: 0, color: '#83bff6' },
              { offset: 1, color: '#188df0' }
            ])
          },
          label: {
            show: true,
            position: 'right'
          }
        }
      ]
    }

    employeeChart.setOption(option)
  } catch (error) {
    console.error('加载员工排名失败', error)
  }
}

// 2. 配件领用占比
const loadPartPie = async (params) => {
  try {
    const res = await getPartReport(params)
    const data = (res.data || []).slice(0, 10)

    const option = {
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 'left',
        data: data.map(item => item.part_name)
      },
      series: [
        {
          name: '领用数量',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: true,
            formatter: '{b}: {d}%'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 16,
              fontWeight: 'bold'
            }
          },
          data: data.map(item => ({
            value: item.total_borrow,
            name: item.part_name
          }))
        }
      ]
    }

    partPieChart.setOption(option)
  } catch (error) {
    console.error('加载配件占比失败', error)
  }
}

// 3. 月度领用趋势
const loadMonthlyTrend = async (params) => {
  try {
    const res = await getDetailedReport(params)
    const data = res.data || []

    // 按月份汇总
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
        trigger: 'axis'
      },
      legend: {
        data: ['领用数量', '损毁数量']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: months
      },
      yAxis: {
        type: 'value',
        name: '数量'
      },
      series: [
        {
          name: '领用数量',
          type: 'line',
          smooth: true,
          data: borrowData,
          itemStyle: { color: '#5470c6' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(84, 112, 198, 0.3)' },
              { offset: 1, color: 'rgba(84, 112, 198, 0)' }
            ])
          }
        },
        {
          name: '损毁数量',
          type: 'line',
          smooth: true,
          data: damagedData,
          itemStyle: { color: '#ee6666' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(238, 102, 102, 0.3)' },
              { offset: 1, color: 'rgba(238, 102, 102, 0)' }
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

// 4. 损毁率分析
const loadDamageRate = async (params) => {
  try {
    const res = await getPartReport(params)
    const data = (res.data || [])
      .filter(item => item.damage_rate > 0)
      .sort((a, b) => b.damage_rate - a.damage_rate)
      .slice(0, 10)

    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'shadow' }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: data.map(item => item.part_name),
        axisLabel: {
          interval: 0,
          rotate: 30
        }
      },
      yAxis: {
        type: 'value',
        name: '损毁率(%)'
      },
      series: [
        {
          name: '损毁率',
          type: 'bar',
          data: data.map(item => item.damage_rate || 0),
          itemStyle: {
            color: params => {
              const value = params.value
              if (value > 10) return '#e74c3c'
              if (value > 5) return '#f39c12'
              return '#27ae60'
            }
          },
          label: {
            show: true,
            position: 'top',
            formatter: '{c}%'
          }
        }
      ]
    }

    damageChart.setOption(option)
  } catch (error) {
    console.error('加载损毁率失败', error)
  }
}

onMounted(async () => {
  await nextTick()
  initCharts()
  loadData()
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
.data-analytics {
  height: 100%;
  overflow-y: auto;
  padding-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
