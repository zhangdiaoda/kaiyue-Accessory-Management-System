<template>
  <div class="notification-config">
    <el-card class="header-card">
      <template #header>
        <div class="card-header">
          <span><el-icon><Bell /></el-icon> 通知系统配置</span>
        </div>
      </template>

      <el-tabs v-model="activeTab" type="border-card">
        <!-- 钉钉配置 -->
        <el-tab-pane label="钉钉通知" name="dingtalk">
          <el-form :model="dingTalkForm" label-width="120px" style="max-width: 600px">
            <el-form-item label="配置名称">
              <el-input v-model="dingTalkForm.config_name" placeholder="如:默认群机器人" />
            </el-form-item>
            
            <el-form-item label="Webhook URL">
              <el-input 
                v-model="dingTalkForm.webhook_url" 
                placeholder="https://oapi.dingtalk.com/robot/send?access_token=..."
                type="textarea"
                :rows="3"
              />
            </el-form-item>
            
            <el-form-item label="加签密钥">
              <el-input 
                v-model="dingTalkForm.secret" 
                placeholder="SECxxxx (可选,增强安全性)"
                type="password"
                show-password
              />
              <div class="form-tip">
                💡 在钉钉群机器人设置中开启"加签"后获取
              </div>
            </el-form-item>
            
            <el-form-item label="@手机号">
              <el-select 
                v-model="dingTalkForm.at_mobiles" 
                multiple 
                filterable 
                allow-create
                placeholder="输入手机号后回车添加"
                style="width: 100%"
              />
              <div class="form-tip">
                可选,填写后通知会@指定人员
              </div>
            </el-form-item>
            
            <el-form-item label="@所有人">
              <el-switch v-model="dingTalkForm.is_at_all" />
            </el-form-item>
            
            <el-form-item label="启用状态">
              <el-switch v-model="dingTalkForm.is_enabled" />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="saveDingTalkConfig" :loading="saving">
                <el-icon><Check /></el-icon> 保存配置
              </el-button>
              <el-button @click="testDingTalk" :loading="testing">
                <el-icon><Promotion /></el-icon> 发送测试消息
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 微信配置 -->
        <el-tab-pane label="微信公众号" name="wechat">
          <el-alert 
            title="微信公众号要求" 
            type="warning" 
            :closable="false"
            style="margin-bottom: 20px"
          >
            <p>1. 必须是<strong>服务号</strong>(订阅号不支持)</p>
            <p>2. 需要用户在公众号内主动订阅后才能发送消息</p>
            <p>3. 每次订阅只能发送1条消息</p>
          </el-alert>
          
          <el-form :model="wechatForm" label-width="120px" style="max-width: 600px">
            <el-form-item label="配置名称">
              <el-input v-model="wechatForm.config_name" placeholder="如:仓储管理公众号" />
            </el-form-item>
            
            <el-form-item label="AppID">
              <el-input v-model="wechatForm.app_id" placeholder="wx..." />
            </el-form-item>
            
            <el-form-item label="AppSecret">
              <el-input 
                v-model="wechatForm.app_secret" 
                placeholder="微信公众平台获取"
                type="password"
                show-password
              />
            </el-form-item>
            
            <el-form-item label="模板ID">
              <el-input 
                v-model="wechatForm.template_id" 
                placeholder="在公众号后台订阅通知中获取"
              />
              <div class="form-tip">
                模板内容格式: thing1(标题) + time2(时间) + thing3(内容)
              </div>
            </el-form-item>
            
            <el-form-item label="启用状态">
              <el-switch v-model="wechatForm.is_enabled" />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="saveWechatConfig" :loading="saving">
                <el-icon><Check /></el-icon> 保存配置
              </el-button>
              <el-button @click="showWechatTest = true">
                <el-icon><Promotion /></el-icon> 发送测试消息
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 通知场景配置 -->
        <el-tab-pane label="通知场景" name="scenes">
          <el-table :data="scenes" style="width: 100%">
            <el-table-column prop="name" label="场景名称" width="200" />
            <el-table-column prop="description" label="说明" />
            <el-table-column label="钉钉" width="80" align="center">
              <template #default="{ row }">
                <el-icon v-if="row.dingtalk" color="green"><Check /></el-icon>
                <el-icon v-else color="gray"><Close /></el-icon>
              </template>
            </el-table-column>
            <el-table-column label="微信" width="80" align="center">
              <template #default="{ row }">
                <el-icon v-if="row.wechat" color="green"><Check /></el-icon>
                <el-icon v-else color="gray"><Close /></el-icon>
              </template>
            </el-table-column>
            <el-table-column label="站内信" width="80" align="center">
              <template #default="{ row }">
                <el-icon v-if="row.internal" color="green"><Check /></el-icon>
                <el-icon v-else color="gray"><Close /></el-icon>
              </template>
            </el-table-column>
          </el-table>
          
          <el-alert 
            title="提示" 
            type="info" 
            :closable="false"
            style="margin-top: 20px"
          >
            通知场景配置已在代码中定义,当前展示的是系统支持的通知场景及其支持的渠道
          </el-alert>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 微信测试对话框 -->
    <el-dialog v-model="showWechatTest" title="测试微信通知" width="500px">
      <el-form label-width="100px">
        <el-form-item label="用户OpenID">
          <el-input v-model="testOpenID" placeholder="需要先在微信公众号获取测试用户的OpenID" />
        </el-form-item>
        <el-alert type="warning" :closable="false" style="margin-bottom: 15px">
          OpenID可通过微信公众平台测试号或用户关注公众号后获取
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="showWechatTest = false">取消</el-button>
        <el-button type="primary" @click="testWechat" :loading="testing">发送测试</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Bell, Check, Close, Promotion } from '@element-plus/icons-vue'
import { 
  getNotificationConfig, 
  updateNotificationConfig,
  testDingTalk as testDingTalkApi,
  testWechat as testWechatApi
} from '@/api/notification'

const activeTab = ref('dingtalk')
const saving = ref(false)
const testing = ref(false)
const showWechatTest = ref(false)
const testOpenID = ref('')

// 钉钉配置表单
const dingTalkForm = ref({
  config_name: '',
  webhook_url: '',
  secret: '',
  at_mobiles: [],
  is_at_all: false,
  is_enabled: true
})

// 微信配置表单
const wechatForm = ref({
  config_name: '',
  app_id: '',
  app_secret: '',
  template_id: '',
  is_enabled: true
})

// 通知场景
const scenes = ref([
  { name: '库存预警', description: '配件库存低于阈值时触发', dingtalk: true, wechat: true, internal: true },
  { name: '领用通知', description: '创建领用记录时触发', dingtalk: true, wechat: true, internal: true },
  { name: '归还提醒', description: '未归还7天后触发', dingtalk: true, wechat: true, internal: false },
  { name: '每日报表', description: '每日8:00自动推送', dingtalk: true, wechat: false, internal: false },
  { name: '周报推送', description: '每周一9:00自动推送', dingtalk: true, wechat: false, internal: false },
  { name: '系统公告', description: '管理员发布公告时触发', dingtalk: true, wechat: true, internal: true }
])

// 加载配置
const loadConfigs = async () => {
  try {
    // 加载钉钉配置
    const dtRes = await getNotificationConfig('dingtalk')
    if (dtRes.code === 200 && dtRes.data) {
      const config = dtRes.data
      dingTalkForm.value.config_name = config.config_name
      dingTalkForm.value.is_enabled = config.is_enabled
      if (config.config_data) {
        const data = JSON.parse(config.config_data)
        dingTalkForm.value.webhook_url = data.webhook_url || ''
        dingTalkForm.value.secret = data.secret || ''
        dingTalkForm.value.at_mobiles = data.at_mobiles || []
        dingTalkForm.value.is_at_all = data.is_at_all || false
      }
    }

    // 加载微信配置
    const wxRes = await getNotificationConfig('wechat')
    if (wxRes.code === 200 && wxRes.data) {
      const config = wxRes.data
      wechatForm.value.config_name = config.config_name
      wechatForm.value.is_enabled = config.is_enabled
      if (config.config_data) {
        const data = JSON.parse(config.config_data)
        wechatForm.value.app_id = data.app_id || ''
        wechatForm.value.app_secret = data.app_secret || ''
        wechatForm.value.template_id = data.template_id || ''
      }
    }
  } catch (error) {
    console.error('加载配置失败:', error)
  }
}

// 保存钉钉配置
const saveDingTalkConfig = async () => {
  if (!dingTalkForm.value.webhook_url) {
    ElMessage.warning('请填写Webhook URL')
    return
  }

  saving.value = true
  try {
    const configData = {
      webhook_url: dingTalkForm.value.webhook_url,
      secret: dingTalkForm.value.secret,
      at_mobiles: dingTalkForm.value.at_mobiles,
      is_at_all: dingTalkForm.value.is_at_all
    }

    const res = await updateNotificationConfig({
      provider_type: 'dingtalk',
      config_name: dingTalkForm.value.config_name || '钉钉机器人',
      config_data: JSON.stringify(configData),
      is_enabled: dingTalkForm.value.is_enabled
    })

    if (res.code === 200) {
      ElMessage.success('保存成功')
    } else {
      ElMessage.error(res.message || '保存失败')
    }
  } catch (error) {
    ElMessage.error('保存失败: ' + error.message)
  } finally {
    saving.value = false
  }
}

// 保存微信配置
const saveWechatConfig = async () => {
  if (!wechatForm.value.app_id || !wechatForm.value.app_secret) {
    ElMessage.warning('请填写AppID和AppSecret')
    return
  }

  saving.value = true
  try {
    const configData = {
      app_id: wechatForm.value.app_id,
      app_secret: wechatForm.value.app_secret,
      template_id: wechatForm.value.template_id
    }

    const res = await updateNotificationConfig({
      provider_type: 'wechat',
      config_name: wechatForm.value.config_name || '微信公众号',
      config_data: JSON.stringify(configData),
      is_enabled: wechatForm.value.is_enabled
    })

    if (res.code === 200) {
      ElMessage.success('保存成功')
    } else {
      ElMessage.error(res.message || '保存失败')
    }
  } catch (error) {
    ElMessage.error('保存失败: ' + error.message)
  } finally {
    saving.value = false
  }
}

// 测试钉钉
const testDingTalk = async () => {
  if (!dingTalkForm.value.webhook_url) {
    ElMessage.warning('请先填写并保存Webhook URL')
    return
  }

  testing.value = true
  try {
    const res = await testDingTalkApi({
      webhook_url: dingTalkForm.value.webhook_url,
      secret: dingTalkForm.value.secret,
      at_mobiles: dingTalkForm.value.at_mobiles,
      is_at_all: dingTalkForm.value.is_at_all
    })

    if (res.code === 200) {
      ElMessage.success('测试消息发送成功,请检查钉钉群')
    } else {
      ElMessage.error(res.message || '发送失败')
    }
  } catch (error) {
    ElMessage.error('发送失败: ' + error.message)
  } finally {
    testing.value = false
  }
}

// 测试微信
const testWechat = async () => {
  if (!testOpenID.value) {
    ElMessage.warning('请填写测试用户的OpenID')
    return
  }

  testing.value = true
  try {
    const res = await testWechatApi({
      app_id: wechatForm.value.app_id,
      app_secret: wechatForm.value.app_secret,
      template_id: wechatForm.value.template_id,
      openid: testOpenID.value
    })

    if (res.code === 200) {
      ElMessage.success('测试消息发送成功,请检查微信')
      showWechatTest.value = false
    } else {
      ElMessage.error(res.message || '发送失败')
    }
  } catch (error) {
    ElMessage.error('发送失败: ' + error.message)
  } finally {
    testing.value = false
  }
}

onMounted(() => {
  loadConfigs()
})
</script>

<style scoped>
.notification-config {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

:deep(.el-tabs__content) {
  padding: 20px;
}

:deep(.el-alert p) {
  margin: 5px 0;
}
</style>
