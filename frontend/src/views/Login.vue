<template>
  <div class="login-wrapper">
    <!-- Animated Background -->
    <div class="gradient-bg">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
      <div class="blob blob-3"></div>
    </div>

    <!-- Login Card -->
    <div class="login-container">
      <el-card class="login-card glass-effect animate-slide-up" shadow="never">
        <div class="login-header">
          <div class="logo-circle">
            <span v-if="brandConfig.brand_logo" class="brand-logo-text">{{ brandConfig.brand_logo }}</span>
            <el-icon v-else :size="32" color="#0071e3"><Box /></el-icon>
          </div>
          <h2 class="login-title">{{ brandConfig.system_name || '配件仓储管理系统' }}</h2>
          <p class="login-subtitle">{{ brandConfig.login_subtitle || '请使用您的管理员账号登录以继续' }}</p>
        </div>

        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          size="large"
          class="apple-form"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="用户名"
              :prefix-icon="User"
              class="apple-input"
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="登录密码"
              :prefix-icon="Lock"
              show-password
              class="apple-input"
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          
          <div class="login-options">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <el-link type="primary" underline="never">忘记密码？</el-link>
          </div>


          <el-form-item class="submit-item">
            <el-button
              type="primary"
              class="login-btn"
              :loading="loading"
              @click="handleLogin"
            >
              登录系统
            </el-button>
          </el-form-item>
        </el-form>

        <div class="login-footer">
          <div class="footer-links">
            <span>{{ brandConfig.copyright || brandConfig.company_name || '© 2026 WHMS Admin' }}</span>
            <span class="dot">·</span>
            <span>隐私策略</span>
          </div>
        </div>
      </el-card>
    </div>

  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Box } from '@element-plus/icons-vue'
import { login } from '@/api/auth'
import { getBrandingConfig } from '@/api/system'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()
const loginFormRef = ref(null)
const loading = ref(false)
const rememberMe = ref(false)
const brandConfig = ref({})

const loadBranding = async () => {
  try {
    const res = await getBrandingConfig()
    brandConfig.value = res.data || {}
  } catch (error) {
    console.error('加载品牌信息失败', error)
  }
}

onMounted(() => {
  loadBranding()
})


const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const res = await login(loginForm)
      userStore.setToken(res.data.token)
      userStore.setUserInfo(res.data)
      
      ElMessage.success({
        message: '欢迎回来，' + (res.data.real_name || res.data.username),
        duration: 2000
      })
      router.push('/dashboard')
    } catch (error) {
      ElMessage.error(error.message || '登录失败，请检查您的凭据')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-wrapper {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Fluid Background blobs */
.gradient-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: #fbfbfd;
  z-index: -1;
}

.blob {
  position: absolute;
  filter: blur(80px);
  border-radius: 50%;
  opacity: 0.6;
  animation: move 20s infinite alternate;
}

.blob-1 {
  width: 600px;
  height: 600px;
  background: #0071e3;
  left: -100px;
  top: -100px;
}

.blob-2 {
  width: 500px;
  height: 500px;
  background: #f43f5e;
  right: -50px;
  bottom: -50px;
  animation-delay: -5s;
}

.blob-3 {
  width: 400px;
  height: 400px;
  background: #fbbf24;
  left: 40%;
  top: 30%;
  animation-delay: -10s;
}

@keyframes move {
  from { transform: translate(0, 0) scale(1.0); }
  to { transform: translate(100px, 50px) scale(1.1); }
}

.login-container {
  width: 100%;
  max-width: 420px;
  padding: 20px;
  z-index: 1;
}

.login-card {
  padding: 40px 30px;
  border-radius: 28px !important;
  background: rgba(255, 255, 255, 0.72) !important;
  backdrop-filter: blur(20px) !important;
  border: 1px solid rgba(255, 255, 255, 0.3) !important;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.1) !important;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-circle {
  width: 64px;
  height: 64px;
  background: white;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.brand-logo-text {
  font-size: 28px;
  line-height: 1;
}


.login-title {
  font-size: 26px;
  font-weight: 800;
  color: #1d1d1f;
  letter-spacing: -0.01em;
  margin: 0;
}

.login-subtitle {
  font-size: 14px;
  color: #86868b;
  margin: 8px 0 0;
}

.apple-form {
  padding: 0 5px;
}

.apple-input :deep(.el-input__wrapper) {
  background: rgba(0, 0, 0, 0.03) !important;
  border: 1px solid transparent !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.apple-input :deep(.el-input__wrapper.is-focus) {
  background: white !important;
  border-color: #0071e3 !important;
  box-shadow: 0 0 0 4px rgba(0, 113, 227, 0.1) !important;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: -10px 0 25px;
}

.submit-item {
  margin-top: 10px;
}

.login-btn {
  width: 100%;
  height: 50px !important;
  font-size: 16px !important;
  font-weight: 700 !important;
  border-radius: 14px !important;
  background-color: #0071e3 !important;
  border: none !important;
  transition: all 0.3s !important;
}

.login-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 113, 227, 0.3);
}

.login-footer {
  text-align: center;
  margin-top: 40px;
}

.tips {
  font-size: 12px;
  color: #86868b;
  margin-bottom: 20px;
}

.footer-links {
  font-size: 11px;
  color: #b0b0b5;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.animate-slide-up {
  animation: slideUp 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

/* ========== 移动端优化 ========== */
@media (max-width: 768px) {
  .blob {
    filter: blur(60px);
  }

  .blob-1,
  .blob-2,
  .blob-3 {
    width: 300px;
    height: 300px;
  }

  .login-container {
    padding: 16px;
    max-width: 100%;
  }

  .login-card {
    padding: 32px 24px !important;
  }

  .logo-circle {
    width: 56px;
    height: 56px;
  }

  .brand-logo-text {
    font-size: 24px;
  }

  .login-title {
    font-size: 22px;
  }

  .login-subtitle {
    font-size: 13px;
  }

  .login-header {
    margin-bottom: 30px;
  }

  .login-options {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .footer-links {
    flex-direction: column;
    gap: 4px;
  }

  .dot {
    display: none;
  }
}

/* 横屏移动设备 */
@media (max-width: 926px) and (orientation: landscape) {
  .login-card {
    max-height: 85vh;
    overflow-y: auto;
  }

  .login-header {
    margin-bottom: 20px;
  }

  .blob {
    display: none;
  }
}

/* 小屏手机 */
@media (max-width: 375px) {
  .login-card {
    padding: 24px 20px !important;
  }

  .login-title {
    font-size: 20px;
  }

  .login-btn {
    height: 46px !important;
    font-size: 15px !important;
  }
}
</style>
