<template>
  <div class="config-manage">
    <el-card shadow="never" class="glass-effect config-card">
      <template #header>
        <div class="card-header">
          <span>⚙️ 系统高级参数管理</span>
          <el-button type="primary" @click="handleSave" :loading="saving">保存全局设置</el-button>
        </div>
      </template>

      <el-form :model="configData" label-position="top">
        <el-row :gutter="40">
          <el-col :span="12">
            <el-divider content-position="left">🚀 核心业务参数</el-divider>
            <el-form-item label="默认库存预警阈值">
              <el-input-number v-model="configData.default_warning_threshold" :min="1" style="width: 100%" />
              <span class="field-tips">新建配件时自动填入的缺省报警值。</span>
            </el-form-item>
            
            <el-form-item label="默认分页条数">
              <el-input-number v-model="configData.page_size" :min="5" :max="100" style="width: 100%" />
              <span class="field-tips">全系统列表页面的默认数据展示行数。</span>
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-divider content-position="left">🤖 钉钉机器人集成</el-divider>
            <el-form-item label="Webhook 地址">
              <el-input 
                v-model="configData.dingtalk_webhook" 
                placeholder="https://oapi.dingtalk.com/robot/send?access_token=..." 
              />
              <div class="action-bar">
                <el-button link type="primary" @click="handleTestWebhook" :loading="testing">
                  <el-icon><Connection /></el-icon> 测试连接
                </el-button>
              </div>
            </el-form-item>

            <el-form-item label="定时推送时间">
              <el-time-picker
                v-model="configData.push_time"
                format="HH:mm"
                value-format="HH:mm"
                placeholder="选择每日推送时间"
                style="width: 100%"
              />
              <span class="field-tips">每日定时向钉钉群发送库存报表。</span>
            </el-form-item>
          </el-col>
        </el-row>

        <el-alert
          title="品牌说明"
          type="info"
          description="系统名称、Logo 及版权信息已移至「品牌与公告设置」进行统一化管理，此处仅保留技术性参数。"
          show-icon
          :closable="false"
          style="margin-top: 20px"
        />
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Connection } from '@element-plus/icons-vue'
import { testWebhook } from '@/api/dingtalk'
import { getSystemConfig, updateSystemConfig } from '@/api/system'

const saving = ref(false)
const testing = ref(false)

const configData = reactive({
  default_warning_threshold: 10,
  dingtalk_webhook: '',
  push_time: '09:00',
  page_size: 10
})

const loadConfig = async () => {
  try {
    const res = await getSystemConfig()
    if (res.data) {
       // 将字符串类型的数字转回数字类型
       if (res.data.default_warning_threshold) configData.default_warning_threshold = Number(res.data.default_warning_threshold)
       if (res.data.page_size) configData.page_size = Number(res.data.page_size)
       if (res.data.dingtalk_webhook) configData.dingtalk_webhook = res.data.dingtalk_webhook
       if (res.data.push_time) configData.push_time = res.data.push_time
    }
  } catch (error) {
    ElMessage.error('加载系统配置失败')
  }
}

const handleTestWebhook = async () => {
  if (!configData.dingtalk_webhook) {
    ElMessage.warning('请输入 Webhook 地址')
    return
  }

  testing.value = true
  try {
    await testWebhook({ webhook_url: configData.dingtalk_webhook })
    ElMessage.success('✅ 钉钉连接测试成功！')
  } catch (error) {
    ElMessage.error('❌ 连接失败，请检查 Webhook 有效性')
  } finally {
    testing.value = false
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    // 转换为后端需要的键值对格式
    const payload = {
       default_warning_threshold: String(configData.default_warning_threshold),
       page_size: String(configData.page_size),
       dingtalk_webhook: configData.dingtalk_webhook,
       push_time: configData.push_time
    }
    await updateSystemConfig(payload)
    ElMessage.success('系统参数已更新')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadConfig()
})
</script>

<style scoped>
.config-manage {
  max-width: 1000px;
  margin: 0 auto;
}

.config-card {
  border-radius: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 700;
  font-size: 16px;
}

.field-tips {
  font-size: 12px;
  color: #86868b;
  margin-top: 4px;
  display: block;
}

.action-bar {
  margin-top: 8px;
}

:deep(.el-form-item__label) {
  font-weight: 600 !important;
  color: #1d1d1f !important;
  padding-bottom: 8px !important;
}

:deep(.el-divider__text) {
  background-color: #fbfbfd;
  font-weight: 700;
  color: #1d1d1f;
}

/* 移动端响应式 */
@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .card-header :deep(.el-button) {
    width: 100%;
  }
  
  :deep(.el-row) {
    flex-direction: column !important;
  }
  
  :deep(.el-col) {
    max-width: 100% !important;
    flex: 0 0 100% !important;
  }
  
  :deep(.el-input-number) {
    width: 100% !important;
  }
}
</style>
