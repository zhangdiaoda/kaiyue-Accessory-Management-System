<template>
  <div class="message-center">
    <el-card shadow="never" class="glass-effect message-card">
      <template #header>
        <div class="card-header">
          <div class="title-info">
            <span>📥 站内信消息中心</span>
            <el-tag type="danger" v-if="unreadCount > 0">{{ unreadCount }} 条未读</el-tag>
          </div>
          <div class="header-actions">
            <el-button type="primary" @click="handleCompose" class="apple-btn">
              <el-icon><EditPen /></el-icon> 撰写消息
            </el-button>
            <el-button type="default" plain @click="loadMessages" :loading="loading">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </div>
        </div>
      </template>

      <div class="message-layout" v-loading="loading">
        <!-- 消息列表 -->
        <div class="message-list">
          <div 
            v-for="msg in messages" 
            :key="msg.id" 
            class="message-item"
            :class="{ 'is-active': selectedMsg?.id === msg.id, 'is-unread': isIncoming(msg) && !msg.is_read }"
            @click="selectMessage(msg)"
          >
            <div class="msg-main">
              <div class="msg-header-row">
                <el-tag size="small" :type="isIncoming(msg) ? 'success' : 'info'" class="direction-tag">
                  {{ isIncoming(msg) ? '收到' : '已发' }}
                </el-tag>
                <div class="msg-title truncate">{{ msg.title }}</div>
              </div>
              <div class="msg-meta-row">
                <span class="msg-sender truncate">{{ isIncoming(msg) ? `发自: ${msg.sender_name || '系统'}` : `至: ${msg.receiver_name || '全体'}` }}</span>
                <span class="msg-time">{{ formatTime(msg.created_at) }}</span>
              </div>
            </div>
            <div class="msg-dot" v-if="isIncoming(msg) && !msg.is_read"></div>
          </div>
          
          <el-empty v-if="!loading && messages.length === 0" description="暂无消息" />
        </div>

        <!-- 消息详情 -->
        <div class="message-detail">
          <div v-if="selectedMsg" class="detail-content animate-fade-in">
            <div class="detail-header">
              <h3>{{ selectedMsg.title }}</h3>
              <el-tag :type="selectedMsg.receiver_id === 0 ? 'warning' : 'primary'" size="small">
                {{ selectedMsg.receiver_id === 0 ? '公共广播' : '对等消息' }}
              </el-tag>
            </div>
            <div class="detail-meta-group">
              <div class="meta-item">
                <span class="label">发件人：</span>
                <span class="value">{{ selectedMsg.sender_name || '系统' }}</span>
              </div>
              <div class="meta-item">
                <span class="label">收件人：</span>
                <span class="value">{{ selectedMsg.receiver_name }}</span>
              </div>
              <div class="meta-item">
                <span class="label">发送时间：</span>
                <span class="value">{{ new Date(selectedMsg.created_at).toLocaleString() }}</span>
              </div>
            </div>
            <el-divider />
            <div class="detail-body">{{ selectedMsg.content }}</div>
          </div>
          <div v-else class="empty-detail">
            <el-icon :size="48"><Postcard /></el-icon>
            <p>请选择一条消息查看详情</p>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 撰写消息对话框 -->
    <el-dialog
      v-model="composeVisible"
      title="撰写新消息"
      width="550px"
      append-to-body
      class="apple-dialog"
    >
      <el-form :model="composeForm" label-position="top">
        <el-form-item label="接收人" required>
          <el-select
            v-model="composeForm.receiver_id"
            placeholder="请选择接收人"
            style="width: 100%"
            filterable
          >
            <el-option :value="0" label="📢 全体成员（系统广播）" />
            <el-option
              v-for="user in userList"
              :key="user.id"
              :label="`${user.real_name} (@${user.username})`"
              :value="user.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="消息标题" required>
          <el-input v-model="composeForm.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="消息内容" required>
          <el-input
            v-model="composeForm.content"
            type="textarea"
            :rows="6"
            placeholder="请输入正文内容..."
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="composeVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCompose" :loading="sending">确认发送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, EditPen, Postcard } from '@element-plus/icons-vue'
import { useResponsive } from '@/composables/useResponsive'
import { getMyMessages, markMessageRead, sendMessage } from '@/api/message'
import { getAllUsers } from '@/api/auth'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()
const { isMobile } = useResponsive()
const loading = ref(false)
const sending = ref(false)
const messages = ref([])
const selectedMsg = ref(null)
const composeVisible = ref(false)
const userList = ref([])

const composeForm = reactive({
  receiver_id: null,
  title: '',
  content: ''
})

const unreadCount = computed(() => 
  messages.value.filter(m => isIncoming(m) && !m.is_read).length
)

const isIncoming = (msg) => {
  // 注意：如果 receiver_id 是 0 且 sender 不是自己，也算收到的广播
  return msg.receiver_id === userStore.userInfo.id || (msg.receiver_id === 0 && msg.sender_id !== userStore.userInfo.id)
}

const formatTime = (timeStr) => {
  const date = new Date(timeStr)
  const now = new Date()
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }
  return date.toLocaleDateString([], { month: '2-digit', day: '2-digit' })
}

const loadMessages = async () => {
  loading.value = true
  try {
    const res = await getMyMessages()
    messages.value = res.data || []
    if (selectedMsg.value) {
      const updated = messages.value.find(m => m.id === selectedMsg.value.id)
      if (updated) selectedMsg.value = updated
    }
  } catch (error) {
    ElMessage.error('获取消息失败')
  } finally {
    loading.value = false
  }
}

const selectMessage = async (msg) => {
  selectedMsg.value = msg
  if (isIncoming(msg) && !msg.is_read) {
    try {
      await markMessageRead(msg.id)
      msg.is_read = true
    } catch (error) {
      console.error('标记阅读失败')
    }
  }
}

const handleCompose = async () => {
  composeVisible.value = true
  if (userList.value.length === 0) {
    try {
      const res = await getAllUsers()
      // 过滤掉当前用户自己，除非想自发自收测试
      userList.value = (res.data || []).filter(u => u.username !== userStore.userInfo.username)
    } catch (error) {
      ElMessage.error('获取联系人失败')
    }
  }
}

const submitCompose = async () => {
  if (composeForm.receiver_id === null || !composeForm.title || !composeForm.content) {
    ElMessage.warning('请填写完整发信信息')
    return
  }
  
  sending.value = true
  try {
    await sendMessage(composeForm)
    ElMessage.success('消息已投递')
    composeVisible.value = false
    Object.assign(composeForm, { receiver_id: null, title: '', content: '' })
    loadMessages()
  } catch (error) {
    ElMessage.error('发送失败')
  } finally {
    sending.value = false
  }
}

onMounted(() => {
  loadMessages()
})
</script>

<style scoped>
.message-center {
  max-width: 1100px;
  margin: 0 auto;
}

.message-card {
  border-radius: 20px;
  height: calc(100vh - 160px);
  display: flex;
  flex-direction: column;
}

:deep(.el-card__body) {
  flex: 1;
  padding: 0;
  overflow: hidden;
}

.message-layout {
  display: flex;
  height: 100%;
}

.message-list {
  width: 320px;
  border-right: 1px solid rgba(0,0,0,0.05);
  overflow-y: auto;
  background: rgba(0,0,0,0.01);
}

.message-item {
  padding: 16px 20px;
  border-bottom: 1px solid rgba(0,0,0,0.03);
  cursor: pointer;
  position: relative;
  transition: all 0.2s;
}

.message-item:hover {
  background: white;
}

.message-item.is-active {
  background: white;
  box-shadow: inset 4px 0 0 #0071e3;
}

.msg-header-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.direction-tag {
  flex-shrink: 0;
  font-size: 10px;
  padding: 0 4px;
  height: 18px;
  line-height: 18px;
}

.msg-title {
  font-weight: 600;
  color: #1d1d1f;
  font-size: 14px;
  flex: 1;
}

.message-item.is-unread .msg-title {
  font-weight: 700;
}

.msg-meta-row {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #86868b;
}

.msg-sender {
  max-width: 140px;
}

.truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.msg-dot {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  width: 8px;
  height: 8px;
  background: #ff375f;
  border-radius: 50%;
}

.message-detail {
  flex: 1;
  padding: 40px;
  overflow-y: auto;
  background: white;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.detail-content h3 {
  font-size: 24px;
  font-weight: 800;
  margin: 0;
  color: #1d1d1f;
}

.detail-meta-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 24px;
}

.meta-item {
  font-size: 13px;
  display: flex;
  align-items: baseline;
}

.meta-item .label {
  color: #86868b;
  width: 80px;
}

.meta-item .value {
  color: #1d1d1f;
  font-weight: 500;
}

.detail-body {
  line-height: 1.8;
  color: #4b4b4d;
  white-space: pre-wrap;
  font-size: 15px;
}

.empty-detail {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #b0b0b5;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.title-info {
  display: flex;
  align-items: center;
  gap: 12px;
  font-weight: 700;
}

.animate-fade-in {
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

:deep(.el-form-item__label) {
  font-weight: 600 !important;
  color: #1d1d1f !important;
  padding-bottom: 8px !important;
}

/* 移动端响应式 */
@media (max-width: 768px) {
  .message-card {
    height: calc(100vh - 120px);
  }
  
  .card-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .header-actions {
    flex-direction: column;
    width: 100%;
  }
  
  .header-actions :deep(.el-button) {
    width: 100%;
  }
  
  .message-layout {
    flex-direction: column;
  }
  
  .message-list {
    width: 100%;
    max-height: 50vh;
    border-right: none;
    border-bottom: 2px solid rgba(0,0,0,0.08);
  }
  
  .message-item {
    padding: 14px 16px;
  }
  
  .msg-title {
    font-size: 13px;
  }
  
  .message-detail {
    padding: 20px 16px;
  }
  
  .detail-content h3 {
    font-size: 18px;
  }
  
  .detail-meta-group {
    font-size: 12px;
  }
  
  .detail-body {
    font-size: 14px;
  }
}
</style>
