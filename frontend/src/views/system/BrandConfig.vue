<template>
  <div class="brand-config">
    <el-tabs v-model="activeTab" class="apple-tabs">
      <!-- 品牌设置 -->
      <el-tab-pane label="🏢 品牌与公司信息" name="brand">
        <el-card shadow="never" class="glass-effect config-card">
          <template #header>
            <div class="card-header">
              <span>品牌个性化设置</span>
              <el-button type="primary" @click="saveBrandConfig" :loading="saving">保存配置</el-button>
            </div>
          </template>
          
          <el-form :model="brandForm" label-position="top">
            <el-row :gutter="40">
              <el-col :span="12">
                <el-form-item label="软件/系统名称 (Software Name)">
                  <el-input v-model="brandForm.system_name" placeholder="例如：配件仓储管理系统" />
                  <span class="field-tips">展示在侧边栏顶部、登录页主标题和浏览器标签页。</span>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="公司/单位名称 (Company Name)">
                  <el-input v-model="brandForm.company_name" placeholder="例如：某某机械加工有限公司" />
                  <span class="field-tips">展示在顶栏中间标题、登录页页脚等品牌背书位置。</span>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="40">
              <el-col :span="12">
                <el-form-item label="企业 Logo (Emoji 或 URL)">
                  <el-input v-model="brandForm.brand_logo" placeholder="支持输入 Emoji 或图片 URL" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="系统版权声明 (Copyright)">
                  <el-input v-model="brandForm.copyright" placeholder="例如：© 2026 XX科技 版权所有" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="40">
              <el-col :span="12">
                <el-form-item label="登录页副标题 (Login Subtitle)">
                  <el-input v-model="brandForm.login_subtitle" placeholder="登录界面的辅助文字" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 公告管理 -->
      <el-tab-pane label="📣 公告发布系统" name="announcement">
        <el-card shadow="never" class="glass-effect">
          <template #header>
            <div class="card-header">
              <span>系统公告列表</span>
              <el-button type="primary" @click="handleAddAnnouncement">发布新公告</el-button>
            </div>
          </template>

          <el-table :data="announcements" v-loading="loading">
            <el-table-column prop="title" label="标题" min-width="150" />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag :type="row.type === 'POPUP' ? 'danger' : 'warning'">
                  {{ row.type === 'POPUP' ? '弹窗' : '滚动' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-switch 
                  v-model="row.status" 
                  :active-value="1" 
                  :inactive-value="0"
                  @change="handleStatusChange(row)"
                />
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="发布时间" width="180">
              <template #default="{ row }">
                {{ new Date(row.created_at).toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" align="center">
              <template #default="{ row }">
                <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
                <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 公告编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editForm.id ? '编辑公告' : '发布新公告'"
      width="500px"
      class="apple-dialog"
    >
      <el-form :model="editForm" label-position="top">
        <el-form-item label="公告标题" required>
          <el-input v-model="editForm.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="公告内容" required>
          <el-input 
            v-model="editForm.content" 
            type="textarea" 
            :rows="4" 
            placeholder="请输入公告正文" 
          />
        </el-form-item>
        <el-form-item label="展示类型">
          <el-radio-group v-model="editForm.type">
            <el-radio-button label="SCROLL">顶部滚动条</el-radio-button>
            <el-radio-button label="POPUP">强提醒弹窗</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAnnouncement" :loading="submitting">确认发布</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getSystemConfig, updateSystemConfig, getAnnouncements, createAnnouncement, updateAnnouncement, deleteAnnouncement } from '@/api/system'

const activeTab = ref('brand')
const loading = ref(false)
const saving = ref(false)

// 品牌配置
const brandForm = reactive({
  system_name: '',
  company_name: '',
  brand_logo: '',
  copyright: '',
  login_subtitle: ''
})

const loadBrandConfig = async () => {
  try {
    const res = await getSystemConfig()
    // 兼容旧数据 brand_name -> system_name
    if (res.data.brand_name && !res.data.system_name) {
       res.data.system_name = res.data.brand_name
    }
    Object.assign(brandForm, res.data)
  } catch (error) {
    ElMessage.error('加载系统配置失败')
  }
}

const saveBrandConfig = async () => {
  saving.value = true
  try {
    await updateSystemConfig(brandForm)
    ElMessage.success('品牌配置已生效')
    setTimeout(() => location.reload(), 800)
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 公告管理部分保持不变...
const announcements = ref([])
const dialogVisible = ref(false)
const submitting = ref(false)
const editForm = reactive({
  id: null,
  title: '',
  content: '',
  type: 'SCROLL',
  status: 1
})

const loadAnnouncements = async () => {
  loading.value = true
  try {
    const res = await getAnnouncements()
    announcements.value = res.data || []
  } catch (error) {
    ElMessage.error('加载公告失败')
  } finally {
    loading.value = false
  }
}

const handleAddAnnouncement = () => {
  Object.assign(editForm, { id: null, title: '', content: '', type: 'SCROLL', status: 1 })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  Object.assign(editForm, row)
  dialogVisible.value = true
}

const handleStatusChange = async (row) => {
  try {
    await updateAnnouncement(row.id, row)
    ElMessage.success('状态已更新')
  } catch (error) {
    ElMessage.error('操作失败')
    row.status = row.status === 1 ? 0 : 1
  }
}

const submitAnnouncement = async () => {
  if (!editForm.title || !editForm.content) {
    ElMessage.warning('请填写完整信息')
    return
  }
  submitting.value = true
  try {
    if (editForm.id) {
      await updateAnnouncement(editForm.id, editForm)
    } else {
      await createAnnouncement(editForm)
    }
    ElMessage.success('发布成功')
    dialogVisible.value = false
    loadAnnouncements()
  } catch (error) {
    ElMessage.error('发布失败')
  } finally {
    submitting.value = false
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该公告吗？', '警告', { type: 'warning' }).then(async () => {
    await deleteAnnouncement(row.id)
    ElMessage.success('已删除')
    loadAnnouncements()
  })
}

onMounted(() => {
  loadBrandConfig()
  loadAnnouncements()
})
</script>

<style scoped>
.brand-config {
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

:deep(.el-form-item__label) {
  font-weight: 600 !important;
  color: #1d1d1f !important;
  padding-bottom: 8px !important;
}

:deep(.el-tabs__header) {
  margin-bottom: 24px;
}

:deep(.el-tabs__item) {
  font-size: 16px;
  font-weight: 600;
  height: 50px;
}
</style>
