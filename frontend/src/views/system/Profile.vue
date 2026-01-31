<template>
  <div class="profile-container">
    <el-card class="profile-card glass-effect">
      <template #header>
        <div class="card-header">
          <span class="title">👤 个人设置</span>
          <p class="subtitle">管理您的账号基本信息与安全设置</p>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-position="top"
        class="apple-form"
      >
        <el-row :gutter="40">
          <el-col :span="12">
            <el-form-item label="登录账号">
              <el-input v-model="userStore.userInfo.username" disabled />
              <p class="form-tip">账号名称不可修改</p>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属角色">
              <el-tag>{{ userStore.userInfo.role === 'SUPER_ADMIN' ? '超级管理员' : '库管员' }}</el-tag>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="40">
          <el-col :span="12">
            <el-form-item label="真实姓名" prop="real_name">
              <el-input v-model="formData.real_name" placeholder="请输入您的真实姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="phone">
              <el-input v-model="formData.phone" placeholder="请输入手机号" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider>安全设置</el-divider>

        <el-row :gutter="40">
          <el-col :span="12">
            <el-form-item label="新密码" prop="password">
              <el-input
                v-model="formData.password"
                type="password"
                show-password
                placeholder="留空则不修改密码"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input
                v-model="formData.confirmPassword"
                type="password"
                show-password
                placeholder="请再次输入新密码"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <div class="form-actions">
          <el-button type="primary" :loading="loading" @click="handleSubmit" class="submit-btn">
            保存更改
          </el-button>
          <el-button @click="resetForm">撤销重置</el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { updateProfile, getUserInfo } from '@/api/auth'

const userStore = useUserStore()
const loading = ref(false)
const formRef = ref(null)

const formData = reactive({
  real_name: '',
  phone: '',
  password: '',
  confirmPassword: ''
})

const validatePass2 = (rule, value, callback) => {
  if (formData.password && value !== formData.password) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const rules = {
  real_name: [{ required: true, message: '姓名不能为空', trigger: 'blur' }],
  phone: [{ pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }],
  password: [{ min: 6, message: '密码长度至少为 6 位', trigger: 'blur' }],
  confirmPassword: [{ validator: validatePass2, trigger: 'blur' }]
}

const initData = () => {
  formData.real_name = userStore.userInfo.realName
  formData.phone = userStore.userInfo.phone || ''
  formData.password = ''
  formData.confirmPassword = ''
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate()
  
  loading.value = true
  try {
    const submitData = {
      real_name: formData.real_name,
      phone: formData.phone
    }
    if (formData.password) {
      submitData.password = formData.password
    }

    await updateProfile(submitData)
    ElMessage.success('个人信息更新成功')
    
    // 重新获取用户信息更新 store
    const res = await getUserInfo()
    userStore.setUserInfo(res.data)
    
    formData.password = ''
    formData.confirmPassword = ''
  } catch (error) {
    ElMessage.error(error.message || '更新失败')
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  initData()
}

onMounted(() => {
  initData()
})
</script>

<style scoped>
.profile-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px 0;
}

.profile-card {
  border-radius: 20px;
}

.card-header {
  border-bottom: none;
}

.card-header .title {
  font-size: 24px;
  font-weight: 800;
  color: #1d1d1f;
}

.card-header .subtitle {
  margin: 8px 0 0;
  font-size: 14px;
  color: #86868b;
}

.apple-form {
  padding: 10px 0;
}

.form-tip {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 4px;
}

.form-actions {
  margin-top: 40px;
  display: flex;
  gap: 16px;
}

.submit-btn {
  padding-left: 32px;
  padding-right: 32px;
}

:deep(.el-form-item__label) {
  font-weight: 600 !important;
  color: #1d1d1f !important;
  padding-bottom: 8px !important;
}

:deep(.el-divider__text) {
  background-color: transparent !important;
  font-weight: 600;
  color: #86868b;
}
</style>
